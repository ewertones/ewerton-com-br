package data

import (
	"encoding/json"
	"os"
	"strings"
	"testing"
	"time"

	"ewerton-go/internal/i18n"
)

// The resume is assembled by joining Go data to i18n keys. A renamed or missing
// key would not fail the build, it would silently ship a resume with a blank
// role or an empty section, so these tests check the join itself.

func store(t *testing.T) *i18n.Store {
	t.Helper()
	en, err := os.ReadFile("../../messages/en.json")
	if err != nil {
		t.Fatalf("read en.json: %v", err)
	}
	ptBR, err := os.ReadFile("../../messages/pt-br.json")
	if err != nil {
		t.Fatalf("read pt-br.json: %v", err)
	}
	s, err := i18n.New(en, ptBR)
	if err != nil {
		t.Fatalf("init i18n: %v", err)
	}
	return s
}

func build(t *testing.T, locale string) Resume {
	t.Helper()
	return BuildResume(locale, store(t).TFunc(locale), time.Date(2026, 9, 9, 0, 0, 0, 0, time.UTC))
}

func TestBuildResumeIsCompleteInEveryLocale(t *testing.T) {
	for _, locale := range []string{"en", "pt-br"} {
		t.Run(locale, func(t *testing.T) {
			r := build(t, locale)

			if r.Basics.Headline == "" || r.Basics.Summary == "" {
				t.Error("basics headline or summary is empty")
			}
			if len(r.Work) != len(workTimeline) {
				t.Fatalf("got %d positions, want %d", len(r.Work), len(workTimeline))
			}
			for _, w := range r.Work {
				if w.Role == "" || w.Period == "" || w.Summary == "" {
					t.Errorf("%s: role/period/summary incomplete: %+v", w.Slug, w)
				}
				if len(w.Highlights) == 0 {
					t.Errorf("%s: no highlights parsed from description", w.Slug)
				}
				if w.Company == "" || w.CompanyURL == "" || w.Industry == "" {
					t.Errorf("%s: company details incomplete", w.Slug)
				}
				if w.StartISO == "" {
					t.Errorf("%s: missing machine-readable start date", w.Slug)
				}
			}
			if len(r.Education) == 0 || len(r.Skills) == 0 || len(r.Projects) == 0 {
				t.Error("education, skills or projects section is empty")
			}
			for _, p := range r.Projects {
				if p.Name == "" || p.Tagline == "" || p.Description == "" {
					t.Errorf("%s: project copy incomplete", p.Slug)
				}
				if len(p.Stack) == 0 {
					t.Errorf("%s: project has no stack", p.Slug)
				}
			}
			for _, l := range r.Languages {
				if l.Name == "" || l.Fluency == "" || l.Code == "" {
					t.Errorf("language incomplete: %+v", l)
				}
			}
		})
	}
}

// Every partner the site offers a description for must appear in the timeline,
// otherwise a role shown on the page would be missing from the resume.
func TestWorkTimelineCoversEveryDescribedPartner(t *testing.T) {
	inTimeline := map[string]bool{}
	for _, e := range workTimeline {
		inTimeline[e.Slug] = true
	}
	for _, group := range PartnerGroups {
		for _, p := range group.Partners {
			if p.HasDescription && !inTimeline[p.Slug] {
				t.Errorf("partner %q has a description on the site but no timeline entry", p.Slug)
			}
		}
	}
}

func TestWorkIsReverseChronological(t *testing.T) {
	r := build(t, "en")
	for i := 1; i < len(r.Work); i++ {
		if r.Work[i].StartISO > r.Work[i-1].StartISO {
			t.Errorf("position %d (%s) starts after the one before it", i, r.Work[i].Slug)
		}
	}
	if !r.Work[0].Current {
		t.Error("the most recent position should be the current one")
	}
}

func TestWorkByCompanyMergesRepeatedEmployers(t *testing.T) {
	r := build(t, "en")
	groups := r.WorkByCompany()

	var stagwell *CompanyGroup
	for i := range groups {
		if groups[i].Company == "Stagwell" {
			stagwell = &groups[i]
		}
	}
	if stagwell == nil {
		t.Fatal("Stagwell group not found")
	}
	if len(stagwell.Roles) != 2 {
		t.Fatalf("got %d Stagwell roles, want 2", len(stagwell.Roles))
	}
	// The tenure must span both roles: the older role's start to the newer
	// role's end.
	if want := "Apr 2024 - Feb 2026"; stagwell.Period != want {
		t.Errorf("tenure = %q, want %q", stagwell.Period, want)
	}
}

func TestSplitDescriptionSeparatesLeadAndBullets(t *testing.T) {
	summary, highlights := splitDescription(
		"Lead paragraph.\n\nKey achievements:\n- First thing.\n- Second thing.")

	if summary != "Lead paragraph." {
		t.Errorf("summary = %q, want %q", summary, "Lead paragraph.")
	}
	if len(highlights) != 2 || highlights[0] != "First thing." {
		t.Errorf("highlights = %#v", highlights)
	}
}

// Unreleased products must not advertise a link before the launch deadline,
// matching the gate the landing page applies.
func TestProjectsStayUnlinkedUntilLaunch(t *testing.T) {
	tf := store(t).TFunc("en")

	before := buildProjects(tf, LaunchDeadline.Add(-time.Hour))
	after := buildProjects(tf, LaunchDeadline)

	for i, p := range before {
		if !projectOrder[i].Gated {
			continue
		}
		if p.URL != "" {
			t.Errorf("%s: linked before launch", p.Slug)
		}
		if p.HomeURL == "" || p.DisplayURL == "" {
			t.Errorf("%s: destination should still be known", p.Slug)
		}
		if after[i].URL == "" {
			t.Errorf("%s: still unlinked after launch", p.Slug)
		}
	}
}

func TestMainStackIsBounded(t *testing.T) {
	for _, p := range build(t, "en").Projects {
		if got := len(p.MainStack()); got > MainStackSize {
			t.Errorf("%s: MainStack has %d entries, want at most %d", p.Slug, got, MainStackSize)
		}
	}
}

// The two locale files must expose the same key set, or a translated page
// silently falls back to English mid-document.
func TestLocaleFilesHaveMatchingKeys(t *testing.T) {
	keys := func(path string) map[string]bool {
		raw, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("read %s: %v", path, err)
		}
		var tree map[string]any
		if err := json.Unmarshal(raw, &tree); err != nil {
			t.Fatalf("parse %s: %v", path, err)
		}
		out := map[string]bool{}
		var walk func(prefix string, node map[string]any)
		walk = func(prefix string, node map[string]any) {
			for k, v := range node {
				path := strings.TrimPrefix(prefix+"."+k, ".")
				if sub, ok := v.(map[string]any); ok {
					walk(path, sub)
					continue
				}
				out[path] = true
			}
		}
		walk("", tree)
		return out
	}

	en := keys("../../messages/en.json")
	ptBR := keys("../../messages/pt-br.json")
	for k := range en {
		if !ptBR[k] {
			t.Errorf("pt-br.json is missing %q", k)
		}
	}
	for k := range ptBR {
		if !en[k] {
			t.Errorf("en.json is missing %q", k)
		}
	}
}
