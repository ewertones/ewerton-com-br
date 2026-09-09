package main

import (
	"crypto/sha256"
	"embed"
	"encoding/base64"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
	"sync"

	"ewerton-go/internal/handler"
	"ewerton-go/internal/i18n"
)

//go:embed templates
var templatesFS embed.FS

//go:embed static
var staticFS embed.FS

// staticCache makes browsers revalidate every static asset instead of trusting
// a long max-age.
//
// The assets are embedded in the binary, so their content only changes on
// deploy, but an "immutable" year-long cache needs a URL that changes with it.
// Revalidation keeps the URLs clean; the ETag keeps it cheap, since an unchanged
// file answers with a 304 and no body.
func staticCache(staticFS fs.FS, next http.Handler) http.Handler {
	var etags sync.Map // path -> ETag

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "public, no-cache")

		name := strings.TrimPrefix(path.Clean(r.URL.Path), "/static/")
		if tag, ok := etags.Load(name); ok {
			w.Header().Set("ETag", tag.(string))
		} else if body, err := fs.ReadFile(staticFS, name); err == nil {
			sum := sha256.Sum256(body)
			tag := `"` + base64.RawURLEncoding.EncodeToString(sum[:12]) + `"`
			etags.Store(name, tag)
			w.Header().Set("ETag", tag)
		}

		next.ServeHTTP(w, r)
	})
}

//go:embed messages/en.json messages/pt-br.json
var messagesFS embed.FS

const version = "1"

func main() {
	// Load i18n
	enJSON, err := messagesFS.ReadFile("messages/en.json")
	if err != nil {
		log.Fatal("read en.json:", err)
	}
	ptBRJSON, err := messagesFS.ReadFile("messages/pt-br.json")
	if err != nil {
		log.Fatal("read pt-br.json:", err)
	}

	i18nStore, err := i18n.New(enJSON, ptBRJSON)
	if err != nil {
		log.Fatal("init i18n:", err)
	}

	// Static file server
	staticSub, err := fs.Sub(staticFS, "static")
	if err != nil {
		log.Fatal("static fs:", err)
	}

	// Resume artifacts (PDF, JSON Resume, Markdown, schema.org graph) are all
	// generated from the same data the page renders.
	resumeHandler := handler.NewResumeHandler(staticSub, i18nStore)

	// Build page handler
	pageHandler, err := handler.NewPageHandler(templatesFS, staticSub, i18nStore, version, resumeHandler)
	if err != nil {
		log.Fatal("parse templates:", err)
	}

	mux := http.NewServeMux()

	// Root redirect -> /en
	mux.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/en", http.StatusFound)
	})

	// Convenience redirects
	mux.HandleFunc("GET /linkedin", handler.LinkedInRedirect)
	mux.HandleFunc("GET /cv", handler.CVRedirect)

	// Generated resume, in every format. The locale-prefixed variants let the
	// language switcher and the sitemap address a specific translation.
	mux.HandleFunc("GET /resume.pdf", resumeHandler.PDF)
	mux.HandleFunc("GET /resume.json", resumeHandler.JSON)
	mux.HandleFunc("GET /resume.md", resumeHandler.Markdown)
	// Registered per locale rather than with a {locale} wildcard, which would
	// collide with the /static/ tree.
	for _, locale := range handler.Locales {
		mux.HandleFunc("GET /"+locale+"/resume.pdf", resumeHandler.PDF)
		mux.HandleFunc("GET /"+locale+"/resume.json", resumeHandler.JSON)
		mux.HandleFunc("GET /"+locale+"/resume.md", resumeHandler.Markdown)
		mux.HandleFunc("GET /"+locale+"/llms.txt", resumeHandler.LLMs)
		mux.HandleFunc("GET /"+locale+"/llms-full.txt", resumeHandler.LLMsFull)
	}

	// Crawler and LLM discovery.
	mux.HandleFunc("GET /robots.txt", resumeHandler.Robots)
	mux.HandleFunc("GET /sitemap.xml", resumeHandler.Sitemap)
	mux.HandleFunc("GET /llms.txt", resumeHandler.LLMs)
	mux.HandleFunc("GET /llms-full.txt", resumeHandler.LLMsFull)

	// Static assets
	fileServer := http.StripPrefix("/static/", http.FileServer(http.FS(staticSub)))
	mux.Handle("GET /static/", staticCache(staticSub, fileServer))

	// Locale pages: /en and /pt-br
	mux.Handle("GET /en", pageHandler)
	mux.Handle("GET /pt-br", pageHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	log.Printf("listening on :%s", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatal(err)
	}
}
