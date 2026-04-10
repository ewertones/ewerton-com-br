package main

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
	"os"

	"ewerton-go/internal/handler"
	"ewerton-go/internal/i18n"
)

//go:embed templates
var templatesFS embed.FS

//go:embed static
var staticFS embed.FS

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

	// Build page handler
	pageHandler, err := handler.NewPageHandler(templatesFS, i18nStore, version)
	if err != nil {
		log.Fatal("parse templates:", err)
	}

	// Static file server
	staticSub, err := fs.Sub(staticFS, "static")
	if err != nil {
		log.Fatal("static fs:", err)
	}

	mux := http.NewServeMux()

	// Root redirect -> /en
	mux.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/en", http.StatusFound)
	})

	// Convenience redirects
	mux.HandleFunc("GET /linkedin", handler.LinkedInRedirect)
	mux.HandleFunc("GET /cv", handler.CVRedirect)

	// Static assets
	mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.FS(staticSub))))

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
