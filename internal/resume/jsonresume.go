package resume

import (
	"encoding/json"
	"strings"

	"ewerton-go/internal/data"
)

// JSON Resume (jsonresume.org) is the closest thing the industry has to a
// standard machine-readable CV, and it is what LLM tooling and ATS importers
// reach for first. The structs below mirror schema v1.0.0.

const jsonResumeSchema = "https://raw.githubusercontent.com/jsonresume/resume-schema/v1.0.0/schema.json"

type jsonResume struct {
	Schema    string       `json:"$schema"`
	Basics    jrBasics     `json:"basics"`
	Work      []jrWork     `json:"work"`
	Education []jrEdu      `json:"education"`
	Projects  []jrProject  `json:"projects"`
	Skills    []jrSkill    `json:"skills"`
	Languages []jrLanguage `json:"languages"`
	Meta      jrMeta       `json:"meta"`
}

type jrBasics struct {
	Name     string      `json:"name"`
	Label    string      `json:"label"`
	Image    string      `json:"image"`
	Email    string      `json:"email"`
	Phone    string      `json:"phone"`
	URL      string      `json:"url"`
	Summary  string      `json:"summary"`
	Location jrLocation  `json:"location"`
	Profiles []jrProfile `json:"profiles"`
}

type jrLocation struct {
	City        string `json:"city"`
	Region      string `json:"region"`
	CountryCode string `json:"countryCode"`
}

type jrProfile struct {
	Network  string `json:"network"`
	Username string `json:"username"`
	URL      string `json:"url"`
}

type jrWork struct {
	Name       string   `json:"name"`
	Position   string   `json:"position"`
	URL        string   `json:"url,omitempty"`
	StartDate  string   `json:"startDate"`
	EndDate    string   `json:"endDate,omitempty"`
	Summary    string   `json:"summary,omitempty"`
	Highlights []string `json:"highlights,omitempty"`
}

type jrEdu struct {
	Institution string   `json:"institution"`
	URL         string   `json:"url,omitempty"`
	Area        string   `json:"area"`
	StudyType   string   `json:"studyType"`
	StartDate   string   `json:"startDate,omitempty"`
	EndDate     string   `json:"endDate,omitempty"`
	Score       string   `json:"score,omitempty"`
	Courses     []string `json:"courses,omitempty"`
}

type jrProject struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Highlights  []string `json:"highlights,omitempty"`
	URL         string   `json:"url,omitempty"`
	Type        string   `json:"type,omitempty"`
	Roles       []string `json:"roles,omitempty"`
	Keywords    []string `json:"keywords,omitempty"`
}

type jrSkill struct {
	Name     string   `json:"name"`
	Level    string   `json:"level"`
	Keywords []string `json:"keywords,omitempty"`
}

type jrLanguage struct {
	Language string `json:"language"`
	Fluency  string `json:"fluency"`
}

type jrMeta struct {
	Canonical    string `json:"canonical"`
	Version      string `json:"version"`
	LastModified string `json:"lastModified"`
	Language     string `json:"language"`
	Generator    string `json:"generator"`
}

// JSONResume renders the resume as an indented JSON Resume document.
func JSONResume(r data.Resume, t func(string) string) ([]byte, error) {
	doc := jsonResume{
		Schema: jsonResumeSchema,
		Basics: jrBasics{
			Name:    r.Basics.Name,
			Label:   r.Basics.Headline,
			Image:   data.SiteOrigin + "/static/images/photos/ewerton.jpg",
			Email:   r.Basics.Email,
			Phone:   r.Basics.PhoneRaw,
			URL:     r.Basics.URL,
			Summary: r.Basics.Summary,
			Location: jrLocation{
				City:        r.Basics.City,
				Region:      r.Basics.Region,
				CountryCode: "BR",
			},
		},
		Meta: jrMeta{
			Canonical:    data.SiteOrigin + data.ResumeJSONPath,
			Version:      r.Generated.Format("2006.01.02"),
			LastModified: r.Generated.Format("2006-01-02T15:04:05Z07:00"),
			Language:     r.Locale,
			Generator:    "ewerton-go (generated from " + data.SiteOrigin + ")",
		},
	}

	for _, p := range r.Basics.Profiles {
		doc.Basics.Profiles = append(doc.Basics.Profiles, jrProfile{
			Network: p.Network, Username: p.Username, URL: p.URL,
		})
	}

	for _, w := range r.Work {
		doc.Work = append(doc.Work, jrWork{
			Name:       w.Company,
			Position:   w.Role,
			URL:        w.CompanyURL,
			StartDate:  w.StartISO,
			EndDate:    w.EndISO,
			Summary:    w.Summary,
			Highlights: w.Highlights,
		})
	}

	for _, e := range r.Education {
		start, end := splitPeriod(e.Period)
		doc.Education = append(doc.Education, jrEdu{
			Institution: e.Institution,
			URL:         e.WebsiteURL,
			StudyType:   degreeType(e.Degree),
			Area:        degreeArea(e.Degree),
			StartDate:   start,
			EndDate:     end,
			Score:       afterColon(e.Grade),
			Courses:     e.Points,
		})
	}

	for _, p := range r.Projects {
		doc.Projects = append(doc.Projects, jrProject{
			Name:        p.Name,
			Description: strings.TrimSpace(p.Tagline + ". " + p.Description),
			Highlights:  p.Features,
			URL:         p.URL,
			Type:        p.Kind,
			Roles:       []string{t("resume.soleDeveloper")},
			Keywords:    p.Stack,
		})
	}

	for _, cat := range r.Skills {
		category := t("landing.technologies.categories." + cat.Slug)
		for _, s := range cat.Skills {
			doc.Skills = append(doc.Skills, jrSkill{
				Name:     s.Name,
				Level:    t("resume.levels." + levelKey(s.Rating)),
				Keywords: append([]string{category}, s.Keywords...),
			})
		}
	}

	for _, lang := range r.Languages {
		doc.Languages = append(doc.Languages, jrLanguage{Language: lang.Name, Fluency: lang.Fluency})
	}

	return json.MarshalIndent(doc, "", "  ")
}

// levelKey buckets a 0-5 self-assessment into a named proficiency level.
func levelKey(rating float64) string {
	switch {
	case rating >= 4.5:
		return "expert"
	case rating >= 3.5:
		return "advanced"
	case rating >= 2.5:
		return "intermediate"
	default:
		return "familiar"
	}
}

// degreeType returns the qualification, e.g. "Graduate Degree" from
// "Graduate Degree, Computer Software Engineering". It mirrors the template
// helper of the same name in internal/handler.
func degreeType(degree string) string {
	if i := strings.Index(degree, ","); i >= 0 {
		return strings.TrimSpace(degree[:i])
	}
	return degree
}

// degreeArea returns the field of study, the other half of the degree string.
func degreeArea(degree string) string {
	if i := strings.Index(degree, ","); i >= 0 {
		return strings.TrimSpace(degree[i+1:])
	}
	return ""
}

// splitPeriod turns "2021 – 2022" into its two years.
func splitPeriod(period string) (start, end string) {
	for _, sep := range []string{"–", "-", "—"} {
		if a, b, ok := strings.Cut(period, sep); ok {
			return strings.TrimSpace(a), strings.TrimSpace(b)
		}
	}
	return strings.TrimSpace(period), ""
}

// afterColon extracts the value from a "Grade: 99/100" style label.
func afterColon(s string) string {
	if _, v, ok := strings.Cut(s, ":"); ok {
		return strings.TrimSpace(v)
	}
	return strings.TrimSpace(s)
}
