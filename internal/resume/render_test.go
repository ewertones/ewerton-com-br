package resume

import (
	"bytes"
	"compress/zlib"
	"encoding/json"
	"os"
	"regexp"
	"strings"
	"testing"
	"time"

	"ewerton-go/internal/data"
	"ewerton-go/internal/i18n"
)

func translator(t *testing.T, locale string) func(string) string {
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
	return store.TFunc(locale)
}

// pdfText pulls the visible text out of a rendered document, which is the same
// thing an ATS or a crawler extracts.
func pdfText(t *testing.T, doc []byte) string {
	t.Helper()
	show := regexp.MustCompile(`\((?:\\.|[^\\()])*\) Tj`)
	var out []string
	for _, m := range regexp.MustCompile(`(?s)stream\n(.*?)\nendstream`).FindAllSubmatch(doc, -1) {
		zr, err := zlib.NewReader(bytes.NewReader(m[1]))
		if err != nil {
			continue
		}
		body := new(bytes.Buffer)
		if _, err := body.ReadFrom(zr); err != nil {
			continue
		}
		for _, run := range show.FindAll(body.Bytes(), -1) {
			raw := run[1:bytes.LastIndexByte(run, ')')]
			raw = bytes.ReplaceAll(raw, []byte(`\(`), []byte("("))
			raw = bytes.ReplaceAll(raw, []byte(`\)`), []byte(")"))
			out = append(out, string(raw))
		}
	}
	return strings.Join(out, " ")
}

func render(t *testing.T, locale string, now time.Time) []byte {
	t.Helper()
	tf := translator(t, locale)
	r := data.BuildResume(locale, tf, now)
	doc, err := RenderPDF(r, tf, Assets{}, nil)
	if err != nil {
		t.Fatalf("RenderPDF: %v", err)
	}
	return doc
}

// The landing page unlocks questao and Harpa Cristã at data.LaunchDeadline. The
// generated documents follow the same gate with no redeploy: before it they
// carry a "Coming Soon" label and the structured data withholds the URL; after
// it the label is gone and the URL is published. The PDF links the destination
// either way, since it is a real address whether or not it has opened.
func TestLaunchGateFlipsGeneratedDocuments(t *testing.T) {
	tf := translator(t, "en")
	dayBefore := data.LaunchDeadline.AddDate(0, 0, -1)
	launchDay := data.LaunchDeadline

	t.Run("before launch", func(t *testing.T) {
		doc := render(t, "en", dayBefore)
		if !strings.Contains(pdfText(t, doc), "Coming Soon") {
			t.Error("expected a Coming Soon label before the deadline")
		}
		raw, err := JSONResume(data.BuildResume("en", tf, dayBefore), tf)
		if err != nil {
			t.Fatal(err)
		}
		if bytes.Contains(raw, []byte(data.QuestaoURL)) {
			t.Error("JSON Resume advertises the URL before launch")
		}
	})

	t.Run("on launch day", func(t *testing.T) {
		doc := render(t, "en", launchDay)
		if strings.Contains(pdfText(t, doc), "Coming Soon") {
			t.Error("Coming Soon label survived the deadline")
		}
		for _, url := range []string{data.QuestaoURL, data.HarpaCristaURL, data.AiomoverURL} {
			if !bytes.Contains(doc, []byte(url)) {
				t.Errorf("PDF is missing a link to %s after launch", url)
			}
		}

		raw, err := JSONResume(data.BuildResume("en", tf, launchDay), tf)
		if err != nil {
			t.Fatal(err)
		}
		var parsed struct {
			Projects []struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"projects"`
		}
		if err := json.Unmarshal(raw, &parsed); err != nil {
			t.Fatal(err)
		}
		for _, p := range parsed.Projects {
			if p.URL == "" {
				t.Errorf("project %q still has no URL after launch", p.Name)
			}
		}
	})
}

func TestRenderedPDFCarriesEverySection(t *testing.T) {
	text := pdfText(t, render(t, "en", time.Date(2026, 9, 9, 0, 0, 0, 0, time.UTC)))

	for _, want := range []string{
		"SUMMARY", "EXPERIENCE", "EDUCATION", "LANGUAGES",
		"PRODUCTS I BUILD", "TECHNOLOGIES & SKILLS",
		data.FullName, data.Email, data.PhoneDisplay,
		"EY", "Stagwell", "Morsum", "DashFi", "CERC", "BeGrowth",
		"aiomover", "Python", "Portuguese",
	} {
		if !strings.Contains(text, want) {
			t.Errorf("rendered PDF is missing %q", want)
		}
	}
}

// Section order is a deliberate editorial choice; assert it so a refactor of
// the render sequence cannot quietly reshuffle the document.
func TestSectionOrder(t *testing.T) {
	text := pdfText(t, render(t, "en", time.Date(2026, 9, 9, 0, 0, 0, 0, time.UTC)))

	want := []string{
		"SUMMARY", "EXPERIENCE", "EDUCATION", "LANGUAGES",
		"PRODUCTS I BUILD", "TECHNOLOGIES & SKILLS",
	}
	at := -1
	for _, section := range want {
		i := strings.Index(text, section)
		if i < 0 {
			t.Fatalf("section %q not found", section)
		}
		if i < at {
			t.Errorf("section %q appears out of order", section)
		}
		at = i
	}
}

// Roles at the same employer must be nested under one heading, not repeated.
func TestStagwellAppearsOnceWithBothRoles(t *testing.T) {
	text := pdfText(t, render(t, "en", time.Date(2026, 9, 9, 0, 0, 0, 0, time.UTC)))

	if n := strings.Count(text, "Stagwell stagwellglobal.com"); n != 1 {
		t.Errorf("Stagwell heading appears %d times, want 1", n)
	}
	for _, role := range []string{"Lead Back-end Engineer", "Senior Data Engineer (Stagwell Marketing Cloud)"} {
		if !strings.Contains(text, role) {
			t.Errorf("missing role %q", role)
		}
	}
}

func TestPortugueseRendersAccentedText(t *testing.T) {
	text := pdfText(t, render(t, "pt-br", time.Date(2026, 9, 9, 0, 0, 0, 0, time.UTC)))

	// The content stream holds WinAnsi bytes, so compare against the encoded form.
	for _, want := range []string{"EXPERIÊNCIA", "FORMAÇÃO ACADÊMICA", "IDIOMAS", "Português"} {
		if !strings.Contains(text, winAnsi(want)) {
			t.Errorf("Portuguese document is missing %q", want)
		}
	}
}

// winAnsi mirrors the encoding the content stream uses, for test comparisons.
func winAnsi(s string) string {
	var b strings.Builder
	for _, r := range s {
		if r < 0x100 {
			b.WriteByte(byte(r))
			continue
		}
		b.WriteRune(r)
	}
	return b.String()
}
