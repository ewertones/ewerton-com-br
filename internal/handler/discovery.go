package handler

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"strings"
	"time"

	"ewerton-go/internal/data"
)

// Discovery endpoints. Crawlers, search engines and AI agents all need a
// documented way in; the site previously shipped no robots.txt or sitemap at
// all, so none of them had one.

// robotsTXT welcomes AI crawlers by name. Several of them only honour an
// explicit user-agent group, so listing them beats relying on the wildcard.
const robotsTXT = `# ewerton.com.br — personal site and resume of ` + data.FullName + `
# Humans, search engines and AI agents are all welcome.
#
# Structured resume (JSON Resume schema): /resume.json
# Plain-text resume for language models:   /resume.md
# Site guide for language models:          /llms.txt

User-agent: *
Allow: /

# AI crawlers and assistants, welcomed explicitly. Several of these ignore the
# wildcard group and only read a group naming them directly.
User-agent: GPTBot
User-agent: OAI-SearchBot
User-agent: ChatGPT-User
User-agent: ClaudeBot
User-agent: Claude-User
User-agent: Claude-SearchBot
User-agent: anthropic-ai
User-agent: PerplexityBot
User-agent: Perplexity-User
User-agent: Google-Extended
User-agent: Applebot
User-agent: Applebot-Extended
User-agent: Amazonbot
User-agent: meta-externalagent
User-agent: FacebookBot
User-agent: CCBot
User-agent: cohere-ai
User-agent: Bytespider
User-agent: Diffbot
User-agent: YouBot
User-agent: Timpibot
Allow: /

Sitemap: ` + data.SiteOrigin + `/sitemap.xml
`

// Robots serves robots.txt.
func (h *ResumeHandler) Robots(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("Cache-Control", "public, max-age=86400")
	fmt.Fprint(w, robotsTXT)
}

// ---------------------------------------------------------------------------
// Sitemap
// ---------------------------------------------------------------------------

type urlSet struct {
	XMLName xml.Name     `xml:"urlset"`
	NS      string       `xml:"xmlns,attr"`
	XHTML   string       `xml:"xmlns:xhtml,attr"`
	URLs    []sitemapURL `xml:"url"`
}

type sitemapURL struct {
	Loc        string          `xml:"loc"`
	LastMod    string          `xml:"lastmod,omitempty"`
	ChangeFreq string          `xml:"changefreq,omitempty"`
	Priority   string          `xml:"priority,omitempty"`
	Links      []alternateLink `xml:"xhtml:link,omitempty"`
}

type alternateLink struct {
	Rel      string `xml:"rel,attr"`
	HrefLang string `xml:"hreflang,attr"`
	Href     string `xml:"href,attr"`
}

// Sitemap serves sitemap.xml, including hreflang alternates for both locales
// and every resume format.
func (h *ResumeHandler) Sitemap(w http.ResponseWriter, r *http.Request) {
	lastMod := time.Now().Format("2006-01-02")

	alternates := make([]alternateLink, 0, len(Locales)+1)
	for _, l := range Locales {
		alternates = append(alternates, alternateLink{"alternate", l, data.SiteOrigin + "/" + l})
	}
	alternates = append(alternates, alternateLink{"alternate", "x-default", data.SiteOrigin + "/en"})

	set := urlSet{
		NS:    "http://www.sitemaps.org/schemas/sitemap/0.9",
		XHTML: "http://www.w3.org/1999/xhtml",
	}
	for _, l := range Locales {
		set.URLs = append(set.URLs, sitemapURL{
			Loc:        data.SiteOrigin + "/" + l,
			LastMod:    lastMod,
			ChangeFreq: "monthly",
			Priority:   "1.0",
			Links:      alternates,
		})
	}
	for _, path := range []string{data.ResumePath, data.ResumeJSONPath, data.ResumeMDPath, "/llms.txt"} {
		set.URLs = append(set.URLs, sitemapURL{
			Loc:        data.SiteOrigin + path,
			LastMod:    lastMod,
			ChangeFreq: "monthly",
			Priority:   "0.8",
		})
	}

	w.Header().Set("Content-Type", "application/xml; charset=utf-8")
	w.Header().Set("Cache-Control", "public, max-age=3600")
	fmt.Fprint(w, xml.Header)
	enc := xml.NewEncoder(w)
	enc.Indent("", "  ")
	if err := enc.Encode(set); err != nil {
		http.Error(w, "sitemap error", http.StatusInternalServerError)
	}
}

// ---------------------------------------------------------------------------
// llms.txt
// ---------------------------------------------------------------------------

// LLMs serves /llms.txt, the llmstxt.org convention: a short, linkable index
// written for a model rather than a browser.
//
// Every string is localized, headings included, so a Portuguese reader does not
// get Portuguese prose filed under English section titles.
func (h *ResumeHandler) LLMs(w http.ResponseWriter, r *http.Request) {
	locale := localeFromRequest(r)
	res := h.model(locale)
	t := h.i18n.TFunc(locale)
	origin := data.SiteOrigin

	label := func(key string) string { return t("resume.llms.labels." + key) }
	section := func(key string) string { return "\n## " + t("resume.llms.sections."+key) + "\n\n" }
	link := func(key, href string) string {
		return fmt.Sprintf("- [%s](%s): %s\n", label(key), href, t("resume.llms.descriptions."+key))
	}

	var b strings.Builder
	fmt.Fprintf(&b, "# %s\n\n", res.Basics.Name)
	fmt.Fprintf(&b, "> %s — %s, %s %s. %s\n\n",
		res.Basics.Name, res.Basics.Headline,
		t("resume.llms.basedIn"), res.Basics.Location(), res.Basics.Summary)
	fmt.Fprintf(&b, "%s\n", t("resume.llms.intro"))

	b.WriteString(section("resume"))
	b.WriteString(link("markdown", origin+data.ResumeMDPath))
	b.WriteString(link("json", origin+data.ResumeJSONPath))
	b.WriteString(link("pdf", origin+data.ResumePath))
	b.WriteString(link("full", origin+"/llms-full.txt"))

	b.WriteString(section("contact"))
	fmt.Fprintf(&b, "- %s: %s\n", label("email"), res.Basics.Email)
	fmt.Fprintf(&b, "- %s: %s\n", label("phone"), res.Basics.Phone)
	for _, p := range res.Basics.Profiles {
		fmt.Fprintf(&b, "- %s: %s\n", p.Network, p.URL)
	}

	b.WriteString(section("site"))
	b.WriteString(link("english", origin+"/en"))
	b.WriteString(link("portuguese", origin+"/pt-br"))

	b.WriteString(section("products"))
	for _, p := range res.Projects {
		if p.URL != "" {
			fmt.Fprintf(&b, "- [%s](%s): %s\n", p.Name, p.URL, p.Tagline)
			continue
		}
		fmt.Fprintf(&b, "- %s (%s): %s\n", p.Name, strings.ToLower(t("common.comingSoon")), p.Tagline)
	}

	b.WriteString(section("optional"))
	fmt.Fprintf(&b, "- [%s](%s/sitemap.xml)\n", label("sitemap"), origin)

	writeText(w, b.String())
}

// LLMsFull serves /llms-full.txt: the entire resume in one fetch, for agents
// that would rather not follow links.
func (h *ResumeHandler) LLMsFull(w http.ResponseWriter, r *http.Request) {
	a, err := h.get("md", localeFromRequest(r))
	if err != nil {
		http.Error(w, "generation failed", http.StatusInternalServerError)
		return
	}
	writeText(w, string(a.body))
}

func writeText(w http.ResponseWriter, body string) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("Cache-Control", "public, no-cache")
	w.Header().Set("X-Robots-Tag", "all")
	fmt.Fprint(w, body)
}
