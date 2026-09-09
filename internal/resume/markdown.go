package resume

import (
	"fmt"
	"strings"

	"ewerton-go/internal/data"
)

// Markdown renders the resume as plain Markdown.
//
// This is the format language models handle best: no layout to reverse
// engineer, no column order to guess, headings that map cleanly onto sections.
// It backs both /resume.md and /llms-full.txt.
func Markdown(r data.Resume, t func(string) string) []byte {
	var b strings.Builder
	b.Grow(8192)

	sec := func(key string) {
		fmt.Fprintf(&b, "\n## %s\n\n", t("resume.sections."+key))
	}

	fmt.Fprintf(&b, "# %s\n\n", r.Basics.Name)
	fmt.Fprintf(&b, "**%s** — %s\n\n", r.Basics.Headline, r.Basics.Location())

	fmt.Fprintf(&b, "- Website: <%s>\n", r.Basics.URL)
	fmt.Fprintf(&b, "- Email: <%s>\n", r.Basics.Email)
	fmt.Fprintf(&b, "- Phone: %s\n", r.Basics.Phone)
	for _, p := range r.Basics.Profiles {
		fmt.Fprintf(&b, "- %s: <%s>\n", p.Network, p.URL)
	}

	sec("summary")
	fmt.Fprintf(&b, "%s\n", r.Basics.Summary)

	sec("experience")
	for _, w := range r.Work {
		fmt.Fprintf(&b, "### %s — %s\n\n", w.Role, w.Company)
		fmt.Fprintf(&b, "*%s · %s · %s*\n\n", w.Period, w.Industry, hostOf(w.CompanyURL))
		if w.Summary != "" {
			fmt.Fprintf(&b, "%s\n\n", w.Summary)
		}
		for _, h := range w.Highlights {
			fmt.Fprintf(&b, "- %s\n", h)
		}
		b.WriteString("\n")
	}

	sec("education")
	for _, e := range r.Education {
		fmt.Fprintf(&b, "### %s — %s\n\n", e.Degree, e.Institution)
		fmt.Fprintf(&b, "*%s · %s*\n\n", e.Period, e.Grade)
		for _, pt := range e.Points {
			fmt.Fprintf(&b, "- %s\n", pt)
		}
		if e.Skills != "" {
			fmt.Fprintf(&b, "- %s\n", e.Skills)
		}
		b.WriteString("\n")
	}

	sec("projects")
	for _, p := range r.Projects {
		fmt.Fprintf(&b, "### %s — %s\n\n", p.Name, p.Tagline)
		if p.URL != "" {
			fmt.Fprintf(&b, "*%s · <%s>*\n\n", p.Kind, p.URL)
		} else {
			fmt.Fprintf(&b, "*%s · %s*\n\n", p.Kind, t("common.comingSoon"))
		}
		fmt.Fprintf(&b, "%s\n\n", p.Description)
		for _, f := range p.Features {
			fmt.Fprintf(&b, "- %s\n", f)
		}
		b.WriteString("\n")
	}

	sec("skills")
	b.WriteString("| " + t("resume.table.category") + " | " + t("resume.table.skill") +
		" | " + t("resume.table.level") + " | " + t("resume.table.rating") + " |\n")
	b.WriteString("| --- | --- | --- | --- |\n")
	for _, cat := range r.Skills {
		category := t("landing.technologies.categories." + cat.Slug)
		for _, s := range cat.Skills {
			fmt.Fprintf(&b, "| %s | %s | %s | %s/5 |\n",
				category, s.Name, t("resume.levels."+levelKey(s.Rating)), trimZero(s.Rating))
		}
	}

	sec("languages")
	for _, lang := range r.Languages {
		fmt.Fprintf(&b, "- **%s** — %s\n", lang.Name, lang.Fluency)
	}

	fmt.Fprintf(&b, "\n---\n\n%s %s · %s · [PDF](%s) · [JSON Resume](%s)\n",
		t("resume.generatedOn"), r.Generated.Format("2006-01-02"), data.SiteOrigin,
		data.SiteOrigin+data.ResumePath, data.SiteOrigin+data.ResumeJSONPath)

	return []byte(b.String())
}
