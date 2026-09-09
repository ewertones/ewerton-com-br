package resume

import (
	"encoding/json"

	"ewerton-go/internal/data"
)

// JSONLD builds the schema.org graph embedded in the landing page.
//
// It replaces a hand-written four-line Person stub. Because it is generated
// from the same model as the PDF, the structured data a crawler reads and the
// document a human downloads can never drift apart.
func JSONLD(r data.Resume, t func(string) string, locale string) ([]byte, error) {
	origin := data.SiteOrigin
	pageURL := origin + "/" + locale
	personID := origin + "/#person"
	siteID := origin + "/#website"

	// Employment history, expressed as OrganizationRole so each position keeps
	// its own dates rather than collapsing into a list of company names.
	roles := make([]any, 0, len(r.Work))
	for _, w := range r.Work {
		role := map[string]any{
			"@type":    "OrganizationRole",
			"roleName": w.Role,
			"worksFor": map[string]any{
				"@type": "Organization",
				"name":  w.Company,
				"url":   w.CompanyURL,
			},
			"startDate": w.StartISO,
		}
		if w.EndISO != "" {
			role["endDate"] = w.EndISO
		}
		roles = append(roles, role)
	}

	alumni := make([]any, 0, len(r.Education))
	credentials := make([]any, 0, len(r.Education))
	for _, e := range r.Education {
		alumni = append(alumni, map[string]any{
			"@type": "CollegeOrUniversity",
			"name":  e.Institution,
			"url":   e.WebsiteURL,
		})
		credentials = append(credentials, map[string]any{
			"@type":              "EducationalOccupationalCredential",
			"name":               e.Degree,
			"credentialCategory": degreeType(e.Degree),
			"about":              degreeArea(e.Degree),
			"recognizedBy": map[string]any{
				"@type": "CollegeOrUniversity",
				"name":  e.Institution,
				"url":   e.WebsiteURL,
			},
		})
	}

	sameAs := make([]string, 0, len(r.Basics.Profiles))
	for _, p := range r.Basics.Profiles {
		sameAs = append(sameAs, p.URL)
	}

	languages := make([]any, 0, len(r.Languages))
	for _, l := range r.Languages {
		languages = append(languages, map[string]any{"@type": "Language", "name": l.Name})
	}

	person := map[string]any{
		"@type":       "Person",
		"@id":         personID,
		"name":        r.Basics.Name,
		"givenName":   "Ewerton",
		"familyName":  "de Souza",
		"jobTitle":    r.Basics.Headline,
		"description": r.Basics.Summary,
		"url":         origin,
		"image":       origin + "/static/images/photos/ewerton.jpg",
		"email":       "mailto:" + r.Basics.Email,
		"telephone":   r.Basics.PhoneRaw,
		"address": map[string]any{
			"@type":           "PostalAddress",
			"addressLocality": r.Basics.City,
			"addressRegion":   r.Basics.Region,
			"addressCountry":  "BR",
		},
		"sameAs":        sameAs,
		"knowsAbout":    r.Keywords(),
		"knowsLanguage": languages,
		"worksFor":      roles,
		"alumniOf":      alumni,
		"hasCredential": credentials,
		"hasOccupation": map[string]any{
			"@type":                "Occupation",
			"name":                 r.Basics.Headline,
			"occupationalCategory": "15-1252.00", // O*NET: Software Developers
			"skills":               r.TopSkills(4),
		},
		// Point crawlers straight at the machine-readable artifacts.
		"subjectOf": []any{
			map[string]any{
				"@type":               "DigitalDocument",
				"name":                t("common.viewResume"),
				"url":                 origin + data.ResumePath,
				"encodingFormat":      "application/pdf",
				"dateModified":        r.Generated.Format("2006-01-02"),
				"inLanguage":          locale,
				"isAccessibleForFree": true,
			},
			map[string]any{
				"@type":               "DataDownload",
				"name":                "JSON Resume",
				"url":                 origin + data.ResumeJSONPath,
				"encodingFormat":      "application/json",
				"isAccessibleForFree": true,
			},
		},
	}

	graph := []any{
		map[string]any{
			"@type":      "WebSite",
			"@id":        siteID,
			"url":        origin,
			"name":       r.Basics.Name,
			"inLanguage": locale,
			"publisher":  map[string]any{"@id": personID},
		},
		map[string]any{
			"@type":        "ProfilePage",
			"@id":          pageURL + "#page",
			"url":          pageURL,
			"name":         r.Basics.Name + " — " + r.Basics.Headline,
			"description":  r.Basics.Summary,
			"inLanguage":   locale,
			"isPartOf":     map[string]any{"@id": siteID},
			"about":        map[string]any{"@id": personID},
			"mainEntity":   map[string]any{"@id": personID},
			"dateModified": r.Generated.Format("2006-01-02"),
		},
		person,
	}

	// Owned products, so the projects section is machine-legible too.
	for _, p := range r.Projects {
		app := map[string]any{
			"@type":       "SoftwareApplication",
			"name":        p.Name,
			"description": p.Tagline,
			"author":      map[string]any{"@id": personID},
			"inLanguage":  locale,
		}
		if p.URL != "" {
			app["url"] = p.URL
		}
		graph = append(graph, app)
	}

	return json.Marshal(map[string]any{
		"@context": "https://schema.org",
		"@graph":   graph,
	})
}
