package data

import (
	"net/url"
	"strings"
	"time"
)

// The resume is assembled from the data the landing page already renders:
// partner cards plus their i18n descriptions become the work history, the
// education cards become the education section, the technology catalog becomes
// the skills section, and the project copy becomes the products section.
//
// Nothing here restates that content. If a role is edited in messages/*.json,
// the PDF, the JSON Resume document and the JSON-LD all change with it.

// Resume is the canonical, locale-resolved view of the whole profile.
type Resume struct {
	Locale    string
	Basics    Basics
	Work      []Position
	Education []EducationEntry
	Projects  []Project
	Skills    []TechnologyCategory
	Languages []Language
	Generated time.Time
}

// Basics holds identity and contact details.
type Basics struct {
	Name        string
	Headline    string
	Summary     string
	Email       string
	Phone       string
	PhoneRaw    string // E.164, for structured data
	PhoneURL    string // WhatsApp click-to-chat
	City        string
	Country     string
	Region      string
	URL         string
	LocationURL string // the city on Google Maps
	Profiles    []Profile
}

// Location renders the human-readable "City, Country" line.
func (b Basics) Location() string { return b.City + ", " + b.Country }

// Profile is an external presence such as LinkedIn.
type Profile struct {
	Network  string
	Username string
	URL      string
	Display  string
}

// Position is one role at one company.
type Position struct {
	Slug       string
	Company    string
	CompanyURL string
	DisplayURL string
	Industry   string
	Role       string
	Period     string // localized, as shown on the site
	StartISO   string // machine-readable, YYYY-MM
	EndISO     string // empty while current
	Current    bool
	Summary    string
	Highlights []string
}

// MainStackSize is how many technologies the PDF names for a product. The full
// list still reaches machines through the JSON and Markdown documents.
const MainStackSize = 6

// MainStack returns the defining technologies, in declared order.
func (p Project) MainStack() []string {
	if len(p.Stack) <= MainStackSize {
		return p.Stack
	}
	return p.Stack[:MainStackSize]
}

// Project is a product built and owned end-to-end.
//
// URL is empty while a product is still behind the launch gate, so the JSON
// and Markdown documents never advertise a link that does not resolve yet.
// HomeURL always carries the destination, letting the PDF show where a
// product will live alongside its "Coming Soon" label.
type Project struct {
	Slug        string
	Name        string
	Tagline     string
	Description string
	URL         string
	HomeURL     string
	DisplayURL  string
	Kind        string
	Features    []string
	Stack       []string
	Live        bool
}

// Language is a spoken language, its fluency label and a 0-5 rating so it can
// be rendered with the same bar as the technology skills.
type Language struct {
	Name    string
	Fluency string
	Rating  float64
	Code    string // ISO country code, used to pick the flag drawn beside it
}

// workTimeline fixes the chronological order and supplies machine-readable
// dates. The site groups partners by industry and stores only localized period
// strings ("Maio 2021"), neither of which a resume can use directly.
var workTimeline = []struct {
	Slug     string
	Suffix   string // i18n key suffix: "" for role/period/desc, "1"/"2" for split roles
	StartISO string
	EndISO   string
}{
	{"ey", "", "2026-03", ""},
	{"stagwell", "1", "2024-12", "2026-02"},
	{"stagwell", "2", "2024-04", "2024-11"},
	{"morsum", "", "2023-05", "2024-04"},
	{"dashfi", "", "2021-10", "2023-05"},
	{"cerc", "", "2021-05", "2021-10"},
	{"begrowth", "", "2020-05", "2021-05"},
}

// projectOrder lists the products in the order the landing page shows them.
var projectOrder = []struct {
	Slug       string
	Key        string // i18n key under landing.projects
	URL        string
	DisplayURL string
	Kind       string
	// Gated products stay unlinked until LaunchDeadline, matching the site.
	Gated    bool
	Features []string
	// Stack lists what each product is actually built with. Technology names
	// are the same in every locale, so they are data rather than translations.
	Stack []string
}{
	{"aiomover", "aiomover", AiomoverURL, "aiomover.com", "web", false,
		[]string{"fast", "secure", "simple"},
		[]string{"Python", "FastAPI", "PostgreSQL", "Next.js", "TypeScript", "Google Cloud",
			"SQLAlchemy", "Tailwind CSS", "Terraform", "Docker", "Stripe"}},
	{"questao", "questao", QuestaoURL, "questao.com.br", "web", true,
		[]string{"ai", "gamified", "community"},
		[]string{"Python", "FastAPI", "PostgreSQL", "Next.js", "TypeScript", "Vertex AI",
			"Tailwind CSS", "Terraform", "Google Cloud", "Stripe"}},
	{"harpa-crista", "harpaCrista", HarpaCristaURL, "Google Play", "mobile", true,
		[]string{"search", "offline", "audio"},
		[]string{"Kotlin", "Jetpack Compose", "Room", "Hilt", "Android",
			"Material 3", "Cloudflare Workers", "Python"}},
}

// BuildResume resolves every section for one locale. t is the translation
// function for that locale, as returned by i18n.Store.TFunc.
func BuildResume(locale string, t func(string) string, now time.Time) Resume {
	eduKey := "en"
	if strings.HasPrefix(locale, "pt") {
		eduKey = "pt"
	}

	return Resume{
		Locale: locale,
		Basics: Basics{
			Name:        FullName,
			Headline:    t("resume.headline"),
			Summary:     t("landing.about.description"),
			Email:       Email,
			Phone:       PhoneDisplay,
			PhoneRaw:    PhoneE164,
			PhoneURL:    WhatsAppLink(),
			City:        City,
			Region:      Region,
			Country:     t("resume.country"),
			URL:         SiteURL,
			LocationURL: MapsLink(City + ", " + t("resume.country")),
			Profiles: []Profile{
				{Network: "LinkedIn", Username: LinkedInHandle, URL: LinkedInURL, Display: "linkedin.com/in/" + LinkedInHandle},
				{Network: "GitHub", Username: GitHubHandle, URL: GitHubURL, Display: "github.com/" + GitHubHandle},
				{Network: "Calendly", Username: LinkedInHandle, URL: CalendlyURL, Display: t("resume.bookMeeting")},
			},
		},
		Work:      buildWork(t),
		Education: EducationData[eduKey],
		Projects:  buildProjects(t, now),
		Skills:    TechnologyCategories,
		Languages: []Language{
			{Name: t("resume.languages.portuguese"), Fluency: t("resume.languages.native"), Rating: 5, Code: "BR"},
			{Name: t("resume.languages.english"), Fluency: t("resume.languages.professional"), Rating: 5, Code: "US"},
			{Name: t("resume.languages.french"), Fluency: t("resume.levels.intermediate"), Rating: 3, Code: "FR"},
			{Name: t("resume.languages.spanish"), Fluency: t("resume.levels.intermediate"), Rating: 3, Code: "ES"},
		},
		Generated: now,
	}
}

// buildWork joins the partner catalog to its localized role descriptions.
func buildWork(t func(string) string) []Position {
	index := map[string]Partner{}
	industry := map[string]string{}
	for _, group := range PartnerGroups {
		for _, p := range group.Partners {
			index[p.Slug] = p
			industry[p.Slug] = t("landing.partners.categories." + group.Slug)
		}
	}

	out := make([]Position, 0, len(workTimeline))
	for _, entry := range workTimeline {
		partner, ok := index[entry.Slug]
		if !ok {
			continue
		}
		base := "landing.partners.descriptions." + entry.Slug + "."
		summary, highlights := splitDescription(t(base + "desc" + entry.Suffix))

		out = append(out, Position{
			Slug:       entry.Slug,
			Company:    partner.Name,
			CompanyURL: partner.URL,
			DisplayURL: partner.DisplayURL,
			Industry:   industry[entry.Slug],
			Role:       t(base + "role" + entry.Suffix),
			Period:     t(base + "period" + entry.Suffix),
			StartISO:   entry.StartISO,
			EndISO:     entry.EndISO,
			Current:    entry.EndISO == "",
			Summary:    summary,
			Highlights: highlights,
		})
	}
	return out
}

// splitDescription turns a site description into a lead paragraph plus bullet
// points. Descriptions look like "paragraph\n\n- bullet\n- bullet", sometimes
// with a "Key achievements:" label in between, which the bullet list makes
// redundant once it is rendered as a list.
func splitDescription(desc string) (summary string, highlights []string) {
	var lead []string
	for _, raw := range strings.Split(desc, "\n") {
		line := strings.TrimSpace(raw)
		switch {
		case line == "":
			continue
		case strings.HasPrefix(line, "- "):
			highlights = append(highlights, strings.TrimSpace(strings.TrimPrefix(line, "- ")))
		case len(highlights) == 0:
			lead = append(lead, line)
		default:
			// Stray prose after the bullets: keep it rather than lose it.
			highlights = append(highlights, line)
		}
	}
	if n := len(lead); n > 0 && strings.HasSuffix(lead[n-1], ":") {
		lead = lead[:n-1]
	}
	return strings.Join(lead, " "), highlights
}

// buildProjects resolves the product cards, honouring the same launch gate the
// landing page uses so unreleased products stay unlinked.
func buildProjects(t func(string) string, now time.Time) []Project {
	launched := !now.Before(LaunchDeadline)

	out := make([]Project, 0, len(projectOrder))
	for _, p := range projectOrder {
		base := "landing.projects." + p.Key + "."
		features := make([]string, 0, len(p.Features))
		for _, f := range p.Features {
			features = append(features, t(base+"features."+f))
		}
		live := launched || !p.Gated
		url := p.URL
		if !live {
			url = ""
		}
		out = append(out, Project{
			Slug:        p.Slug,
			Name:        t(base + "title"),
			Tagline:     t(base + "tagline"),
			Description: t(base + "description"),
			URL:         url,
			HomeURL:     p.URL,
			DisplayURL:  p.DisplayURL,
			Kind:        t("landing.projects.groups." + p.Kind),
			Features:    features,
			Stack:       p.Stack,
			Live:        live,
		})
	}
	return out
}

// Keywords collects the skill names, used for PDF metadata and JSON-LD
// knowsAbout. Search and AI indexes lean on these heavily.
func (r Resume) Keywords() []string {
	out := make([]string, 0, 64)
	for _, cat := range r.Skills {
		for _, s := range cat.Skills {
			out = append(out, s.Name)
		}
	}
	return out
}

// TopSkills returns the highest-rated skills across every category, which is
// what a summary line or a metadata keyword list should lead with.
func (r Resume) TopSkills(min float64) []string {
	var out []string
	for _, cat := range r.Skills {
		for _, s := range cat.Skills {
			if s.Rating >= min {
				out = append(out, s.Name)
			}
		}
	}
	return out
}

// CompanyGroup collects every role held at one company, the way LinkedIn shows
// a multi-role tenure under a single employer heading.
type CompanyGroup struct {
	Company    string
	CompanyURL string
	DisplayURL string
	Industry   string
	Period     string // spans the earliest start to the latest end
	Current    bool
	Roles      []Position
}

// WorkByCompany groups consecutive positions at the same employer. The flat
// Work slice stays the source for JSON Resume, where one entry per role is the
// schema's expectation.
func (r Resume) WorkByCompany() []CompanyGroup {
	var out []CompanyGroup
	for _, p := range r.Work {
		if n := len(out); n > 0 && out[n-1].Company == p.Company {
			out[n-1].Roles = append(out[n-1].Roles, p)
			continue
		}
		out = append(out, CompanyGroup{
			Company:    p.Company,
			CompanyURL: p.CompanyURL,
			DisplayURL: p.DisplayURL,
			Industry:   p.Industry,
			Current:    p.Current,
			Roles:      []Position{p},
		})
	}
	// Roles arrive newest first, so the tenure runs from the last role's start
	// to the first role's end.
	for i := range out {
		roles := out[i].Roles
		start, _ := splitRange(roles[len(roles)-1].Period)
		_, end := splitRange(roles[0].Period)
		out[i].Period = start + " - " + end
	}
	return out
}

// splitRange splits a localized period such as "Dec 2024 - Feb 2026".
func splitRange(period string) (start, end string) {
	if a, b, ok := strings.Cut(period, " - "); ok {
		return strings.TrimSpace(a), strings.TrimSpace(b)
	}
	return strings.TrimSpace(period), strings.TrimSpace(period)
}

// WhatsAppLink builds a click-to-chat URL. WhatsApp expects the number in full
// international format with no plus sign, spaces, dashes or leading zeros.
func WhatsAppLink() string {
	return "https://wa.me/" + strings.TrimPrefix(PhoneE164, "+")
}

// MapsLink points at a place using the Google Maps URLs API, which is the
// documented, platform-independent way to deep-link a search.
func MapsLink(query string) string {
	return "https://www.google.com/maps/search/?api=1&query=" + url.QueryEscape(query)
}
