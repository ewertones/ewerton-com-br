package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"testing/fstest"

	"ewerton-go/internal/data"
	"ewerton-go/internal/i18n"
)

func newHandler(t *testing.T) *ResumeHandler {
	t.Helper()
	en, err := os.ReadFile("../../messages/en.json")
	if err != nil {
		t.Fatalf("read en.json: %v", err)
	}
	ptBR, err := os.ReadFile("../../messages/pt-br.json")
	if err != nil {
		t.Fatalf("read pt-br.json: %v", err)
	}
	store, err := i18n.New(en, ptBR)
	if err != nil {
		t.Fatalf("init i18n: %v", err)
	}

	photo, err := os.ReadFile("../../static/images/photos/ewerton.jpg")
	if err != nil {
		t.Fatalf("read portrait: %v", err)
	}
	return NewResumeHandler(fstest.MapFS{
		"images/photos/ewerton.jpg": &fstest.MapFile{Data: photo},
	}, store)
}

func get(t *testing.T, h http.HandlerFunc, path string) *httptest.ResponseRecorder {
	t.Helper()
	rec := httptest.NewRecorder()
	h(rec, httptest.NewRequest(http.MethodGet, path, nil))
	return rec
}

func TestArtifactsAreServedForEveryLocale(t *testing.T) {
	h := newHandler(t)

	cases := []struct {
		name    string
		serve   http.HandlerFunc
		path    string
		mime    string
		leading string
	}{
		{"pdf en", h.PDF, "/en/resume.pdf", "application/pdf", "%PDF-1.7"},
		{"pdf pt", h.PDF, "/pt-br/resume.pdf", "application/pdf", "%PDF-1.7"},
		{"json en", h.JSON, "/en/resume.json", "application/json", "{"},
		{"json pt", h.JSON, "/pt-br/resume.json", "application/json", "{"},
		{"markdown en", h.Markdown, "/en/resume.md", "text/markdown", "# Ewerton"},
		{"markdown pt", h.Markdown, "/pt-br/resume.md", "text/markdown", "# Ewerton"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			rec := get(t, c.serve, c.path)
			if rec.Code != http.StatusOK {
				t.Fatalf("status = %d, want 200", rec.Code)
			}
			if ct := rec.Header().Get("Content-Type"); !strings.HasPrefix(ct, c.mime) {
				t.Errorf("Content-Type = %q, want %q", ct, c.mime)
			}
			if body := rec.Body.String(); !strings.HasPrefix(body, c.leading) {
				t.Errorf("body starts with %q, want %q", body[:min(20, len(body))], c.leading)
			}
			if rec.Header().Get("ETag") == "" {
				t.Error("missing ETag")
			}
		})
	}
}

// The PDF is regenerated per request, so identical requests must produce
// identical bytes or caching is pointless.
func TestRepeatRequestsAreCacheable(t *testing.T) {
	h := newHandler(t)

	first := get(t, h.PDF, "/en/resume.pdf")
	etag := first.Header().Get("ETag")

	second := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/en/resume.pdf", nil)
	req.Header.Set("If-None-Match", etag)
	h.PDF(second, req)

	if second.Code != http.StatusNotModified {
		t.Errorf("status = %d, want 304", second.Code)
	}
	if second.Body.Len() != 0 {
		t.Error("304 response should have an empty body")
	}
}

func TestLocalesProduceDifferentDocuments(t *testing.T) {
	h := newHandler(t)
	en := get(t, h.Markdown, "/en/resume.md").Body.String()
	pt := get(t, h.Markdown, "/pt-br/resume.md").Body.String()

	if en == pt {
		t.Fatal("the two locales produced identical documents")
	}
	if !strings.Contains(pt, "Experiência") {
		t.Error("the Portuguese document is not translated")
	}
	if !strings.Contains(en, "Experience") {
		t.Error("the English document is not in English")
	}
}

func TestDownloadQuerySwitchesDisposition(t *testing.T) {
	h := newHandler(t)

	if got := get(t, h.PDF, "/en/resume.pdf").Header().Get("Content-Disposition"); !strings.HasPrefix(got, "inline") {
		t.Errorf("Content-Disposition = %q, want inline", got)
	}
	if got := get(t, h.PDF, "/en/resume.pdf?download=1").Header().Get("Content-Disposition"); !strings.HasPrefix(got, "attachment") {
		t.Errorf("Content-Disposition = %q, want attachment", got)
	}
}

func TestJSONResumeIsValidAndPopulated(t *testing.T) {
	h := newHandler(t)

	var doc struct {
		Schema string `json:"$schema"`
		Basics struct {
			Name     string `json:"name"`
			Email    string `json:"email"`
			Profiles []struct {
				Network string `json:"network"`
				URL     string `json:"url"`
			} `json:"profiles"`
		} `json:"basics"`
		Work []struct {
			Name      string `json:"name"`
			Position  string `json:"position"`
			StartDate string `json:"startDate"`
		} `json:"work"`
		Skills   []map[string]any `json:"skills"`
		Projects []struct {
			Keywords []string `json:"keywords"`
		} `json:"projects"`
	}
	if err := json.Unmarshal(get(t, h.JSON, "/en/resume.json").Body.Bytes(), &doc); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}

	if doc.Schema == "" {
		t.Error("missing $schema")
	}
	if doc.Basics.Name != data.FullName || doc.Basics.Email == "" {
		t.Errorf("basics = %+v", doc.Basics)
	}
	if len(doc.Basics.Profiles) < 2 {
		t.Errorf("expected at least LinkedIn and GitHub profiles, got %d", len(doc.Basics.Profiles))
	}
	if len(doc.Work) == 0 || len(doc.Skills) == 0 {
		t.Error("work or skills section is empty")
	}
	for _, w := range doc.Work {
		if w.Name == "" || w.Position == "" || w.StartDate == "" {
			t.Errorf("incomplete work entry: %+v", w)
		}
	}
	for _, p := range doc.Projects {
		if len(p.Keywords) == 0 {
			t.Error("project is missing its technology keywords")
		}
	}
}

func TestJSONLDIsValidGraph(t *testing.T) {
	h := newHandler(t)

	var graph struct {
		Context string           `json:"@context"`
		Graph   []map[string]any `json:"@graph"`
	}
	if err := json.Unmarshal([]byte(h.JSONLD("en")), &graph); err != nil {
		t.Fatalf("invalid JSON-LD: %v", err)
	}
	if graph.Context != "https://schema.org" {
		t.Errorf("@context = %q", graph.Context)
	}

	types := map[string]bool{}
	for _, node := range graph.Graph {
		types[node["@type"].(string)] = true
	}
	for _, want := range []string{"WebSite", "ProfilePage", "Person", "SoftwareApplication"} {
		if !types[want] {
			t.Errorf("graph is missing a %s node", want)
		}
	}
}

func TestDiscoveryEndpoints(t *testing.T) {
	h := newHandler(t)

	robots := get(t, h.Robots, "/robots.txt").Body.String()
	for _, agent := range []string{"GPTBot", "ClaudeBot", "PerplexityBot", "Google-Extended"} {
		if !strings.Contains(robots, agent) {
			t.Errorf("robots.txt does not mention %s", agent)
		}
	}
	if !strings.Contains(robots, "Sitemap: "+data.SiteOrigin+"/sitemap.xml") {
		t.Error("robots.txt does not declare the sitemap")
	}

	sitemap := get(t, h.Sitemap, "/sitemap.xml").Body.String()
	for _, want := range []string{"<urlset", "hreflang", data.ResumePath, data.ResumeJSONPath} {
		if !strings.Contains(sitemap, want) {
			t.Errorf("sitemap is missing %q", want)
		}
	}

	llms := get(t, h.LLMs, "/llms.txt").Body.String()
	if !strings.HasPrefix(llms, "# "+data.FullName) {
		t.Error("llms.txt should open with an H1 naming the person")
	}
	for _, want := range []string{data.ResumeJSONPath, data.ResumeMDPath, "llms-full.txt"} {
		if !strings.Contains(llms, want) {
			t.Errorf("llms.txt is missing %q", want)
		}
	}

	full := get(t, h.LLMsFull, "/llms-full.txt").Body.String()
	if !strings.Contains(full, "## Experience") {
		t.Error("llms-full.txt does not contain the resume body")
	}
}

// An unknown locale must fall back rather than fail.
func TestUnknownLocaleFallsBack(t *testing.T) {
	h := newHandler(t)
	rec := get(t, h.JSON, "/klingon/resume.json")
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "Senior Backend") {
		t.Error("expected the English fallback")
	}
}
