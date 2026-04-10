package handler

import (
	"fmt"
	"html/template"
	"io/fs"
	"math"
	"net/http"
	"strings"

	"ewerton-go/internal/data"
	"ewerton-go/internal/i18n"
)

// PageData is passed to every template render.
type PageData struct {
	Locale          string
	OtherLocale     string
	OtherLocaleName string
	T               func(string) string
	Technologies    []data.TechnologyCategory
	Education       []data.EducationEntry
	PartnerGroups   []data.PartnerGroup
	LinkedInURL     string
	AiomoverURL     string
	ResumePath      string
	CalendlyURL     string
	Version         string
	CSSPath         string
	JSPath          string
}

// PageHandler renders the single-page portfolio.
type PageHandler struct {
	tmpl     *template.Template
	i18nSt   *i18n.Store
	staticFS fs.FS
	version  string
}

// NewPageHandler parses all templates and returns a ready handler.
func NewPageHandler(templatesFS fs.FS, staticFS fs.FS, i18nStore *i18n.Store, version string) (*PageHandler, error) {
	funcMap := template.FuncMap{
		"safeHTML": func(s string) template.HTML { return template.HTML(s) },
		"pct": func(rating float64) float64 {
			return math.Round((rating/5.0)*100*10) / 10
		},
		"join":      strings.Join,
		"lower":     strings.ToLower,
		"contains":  strings.Contains,
		"hasSuffix": strings.HasSuffix,
		"dict": func(values ...any) (map[string]any, error) {
			if len(values)%2 != 0 {
				return nil, fmt.Errorf("dict: odd number of arguments")
			}
			m := make(map[string]any, len(values)/2)
			for i := 0; i < len(values); i += 2 {
				k, ok := values[i].(string)
				if !ok {
					return nil, fmt.Errorf("dict: key must be string")
				}
				m[k] = values[i+1]
			}
			return m, nil
		},
		"keywords": func(kws []string) string {
			return strings.Join(kws, " ")
		},
		"hasDegreeComma": func(degree string) bool {
			return strings.Contains(degree, ",")
		},
		"degreeType": func(degree string) string {
			if idx := strings.Index(degree, ","); idx >= 0 {
				return strings.TrimSpace(degree[:idx])
			}
			return ""
		},
		"degreeName": func(degree string) string {
			if idx := strings.Index(degree, ","); idx >= 0 {
				return strings.TrimSpace(degree[idx+1:])
			}
			return degree
		},
	}

	tmpl, err := template.New("").Funcs(funcMap).ParseFS(
		templatesFS,
		"templates/*.html",
		"templates/partials/*.html",
		"templates/components/*.html",
	)
	if err != nil {
		return nil, err
	}

	return &PageHandler{tmpl: tmpl, i18nSt: i18nStore, staticFS: staticFS, version: version}, nil
}

func (h *PageHandler) resolveAsset(path string) string {
	// e.g. /static/css/globals.css -> static/css/globals.min.css
	if !strings.HasPrefix(path, "/static/") {
		return path
	}

	relPath := strings.TrimPrefix(path, "/static/")
	extIdx := strings.LastIndex(relPath, ".")
	if extIdx == -1 {
		return path
	}

	minPath := relPath[:extIdx] + ".min" + relPath[extIdx:]
	if _, err := fs.Stat(h.staticFS, minPath); err == nil {
		return "/static/" + minPath
	}

	return path
}

func (h *PageHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	locale := strings.TrimPrefix(r.URL.Path, "/")
	locale = strings.Split(locale, "/")[0]

	if locale != "en" && locale != "pt-br" {
		http.Redirect(w, r, "/en", http.StatusFound)
		return
	}

	otherLocale := "pt-br"
	otherLocaleName := "Português"
	if locale == "pt-br" {
		otherLocale = "en"
		otherLocaleName = "English"
	}

	eduKey := "en"
	if strings.HasPrefix(locale, "pt") {
		eduKey = "pt"
	}

	pd := PageData{
		Locale:          locale,
		OtherLocale:     otherLocale,
		OtherLocaleName: otherLocaleName,
		T:               h.i18nSt.TFunc(locale),
		Technologies:    data.TechnologyCategories,
		Education:       data.EducationData[eduKey],
		PartnerGroups:   data.PartnerGroups,
		LinkedInURL:     data.LinkedInURL,
		AiomoverURL:     data.AiomoverURL,
		ResumePath:      data.ResumePath,
		CalendlyURL:     data.CalendlyURL,
		Version:         h.version,
		CSSPath:         h.resolveAsset("/static/css/globals.css"),
		JSPath:          h.resolveAsset("/static/js/app.js"),
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := h.tmpl.ExecuteTemplate(w, "layout.html", pd); err != nil {
		http.Error(w, "template error: "+err.Error(), http.StatusInternalServerError)
	}
}
