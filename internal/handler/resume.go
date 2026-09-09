package handler

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"html/template"
	"io/fs"
	"net/http"
	"strings"
	"sync"
	"time"

	"ewerton-go/internal/data"
	"ewerton-go/internal/i18n"
	"ewerton-go/internal/resume"
)

// Locales the site publishes. Anything else redirects to the default.
var Locales = []string{"en", "pt-br"}

// ResumeHandler generates and serves every resume artifact: the PDF, the JSON
// Resume document, the Markdown rendering and the schema.org graph the landing
// page embeds.
//
// Artifacts are rendered on demand and memoised. The cache key carries the day
// and the product-launch state, the only two inputs that change without a
// redeploy, so a long-lived instance still serves an accurate document while
// ETags stay stable for a whole day.
type ResumeHandler struct {
	i18n   *i18n.Store
	assets resume.Assets

	mu    sync.Mutex
	cache map[string]*artifact
}

type artifact struct {
	body        []byte
	etag        string
	modified    time.Time
	contentType string
}

// NewResumeHandler loads the artwork from the embedded static files. Missing
// artwork is not fatal: the PDF falls back to a monogram and drops the mark.
func NewResumeHandler(staticFS fs.FS, store *i18n.Store) *ResumeHandler {
	read := func(name string) []byte {
		body, err := fs.ReadFile(staticFS, name)
		if err != nil {
			return nil
		}
		return body
	}
	return &ResumeHandler{
		i18n: store,
		assets: resume.Assets{
			Photo: read("images/photos/ewerton.jpg"),
			Logo:  read("images/logo-mark.jpg"),
		},
		cache: map[string]*artifact{},
	}
}

// normalizeLocale maps anything unrecognised onto the default locale.
func normalizeLocale(locale string) string {
	for _, l := range Locales {
		if locale == l {
			return l
		}
	}
	return "en"
}

// model builds the canonical resume for a locale, pinned to the start of the
// current day so repeated requests produce byte-identical documents.
func (h *ResumeHandler) model(locale string) data.Resume {
	now := time.Now().In(data.LaunchDeadline.Location())
	day := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	return data.BuildResume(locale, h.i18n.TFunc(locale), day)
}

// get returns a memoised artifact, rendering it on first use.
func (h *ResumeHandler) get(kind, locale string) (*artifact, error) {
	r := h.model(locale)
	key := fmt.Sprintf("%s|%s|%s", kind, locale, r.Generated.Format("2006-01-02"))

	h.mu.Lock()
	defer h.mu.Unlock()
	if a, ok := h.cache[key]; ok {
		return a, nil
	}

	t := h.i18n.TFunc(locale)
	var body []byte
	var contentType string
	var err error

	switch kind {
	case "json":
		body, err = resume.JSONResume(r, t)
		contentType = "application/json; charset=utf-8"
	case "md":
		body = resume.Markdown(r, t)
		contentType = "text/markdown; charset=utf-8"
	case "jsonld":
		body, err = resume.JSONLD(r, t, locale)
		contentType = "application/ld+json"
	case "pdf":
		// The JSON document is embedded inside the PDF as an associated file,
		// so the structured data travels with the document itself.
		var attach []byte
		if attach, err = resume.JSONResume(r, t); err == nil {
			body, err = resume.RenderPDF(r, t, h.assets, attach)
		}
		contentType = "application/pdf"
	default:
		return nil, fmt.Errorf("unknown resume artifact %q", kind)
	}
	if err != nil {
		return nil, err
	}

	sum := sha256.Sum256(body)
	a := &artifact{
		body:        body,
		etag:        `"` + base64.RawURLEncoding.EncodeToString(sum[:12]) + `"`,
		modified:    r.Generated,
		contentType: contentType,
	}
	h.cache[key] = a
	return a, nil
}

// JSONLD returns the schema.org graph for a locale, ready to inline in a
// <script type="application/ld+json"> element.
func (h *ResumeHandler) JSONLD(locale string) template.JS {
	a, err := h.get("jsonld", normalizeLocale(locale))
	if err != nil {
		return ""
	}
	return template.JS(a.body)
}

// serve writes an artifact with validators so repeat visits get a 304.
func (h *ResumeHandler) serve(w http.ResponseWriter, r *http.Request, kind, locale, filename string, inline bool) {
	a, err := h.get(kind, normalizeLocale(locale))
	if err != nil {
		http.Error(w, "resume generation failed", http.StatusInternalServerError)
		return
	}

	disposition := "inline"
	if !inline {
		disposition = "attachment"
	}

	w.Header().Set("Content-Type", a.contentType)
	w.Header().Set("Content-Disposition", fmt.Sprintf("%s; filename=%q", disposition, filename))
	w.Header().Set("ETag", a.etag)
	// The document is regenerated from live site data, so it must never be
	// served from a stale copy. no-cache alone is not enough: Chrome will reuse
	// an iframe's in-memory copy without revalidating, which is exactly how the
	// viewer ends up showing yesterday's resume. no-store closes that door.
	w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")
	w.Header().Set("X-Robots-Tag", "all")
	// Advertise the sibling formats so an agent that fetched one can find the rest.
	w.Header().Set("Link", strings.Join([]string{
		fmt.Sprintf(`<%s>; rel="alternate"; type="application/json"`, data.ResumeJSONPath),
		fmt.Sprintf(`<%s>; rel="alternate"; type="text/markdown"`, data.ResumeMDPath),
		fmt.Sprintf(`<%s>; rel="alternate"; type="application/pdf"`, data.ResumePath),
	}, ", "))

	http.ServeContent(w, r, filename, a.modified, bytes.NewReader(a.body))
}

// localeFromRequest reads the locale from the URL prefix (/pt-br/resume.pdf),
// falling back to Accept-Language on the unprefixed routes.
func localeFromRequest(r *http.Request) string {
	if first, _, ok := strings.Cut(strings.TrimPrefix(r.URL.Path, "/"), "/"); ok {
		for _, l := range Locales {
			if first == l {
				return l
			}
		}
	}
	if strings.Contains(strings.ToLower(r.Header.Get("Accept-Language")), "pt") {
		return "pt-br"
	}
	return "en"
}

// filename builds a descriptive download name, e.g. Ewerton-Souza-Resume-en.pdf.
func filename(locale, ext string) string {
	return fmt.Sprintf("%s-Resume-%s.%s", strings.ReplaceAll(data.FullName, " ", "-"), locale, ext)
}

// PDF serves the generated resume PDF.
func (h *ResumeHandler) PDF(w http.ResponseWriter, r *http.Request) {
	locale := localeFromRequest(r)
	// ?download=1 switches the viewer from inline preview to a save dialog.
	inline := r.URL.Query().Get("download") == ""
	h.serve(w, r, "pdf", locale, filename(locale, "pdf"), inline)
}

// JSON serves the JSON Resume document.
func (h *ResumeHandler) JSON(w http.ResponseWriter, r *http.Request) {
	locale := localeFromRequest(r)
	h.serve(w, r, "json", locale, filename(locale, "json"), true)
}

// Markdown serves the Markdown rendering.
func (h *ResumeHandler) Markdown(w http.ResponseWriter, r *http.Request) {
	locale := localeFromRequest(r)
	h.serve(w, r, "md", locale, filename(locale, "md"), true)
}
