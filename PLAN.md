# Plan: Rewrite Portfolio from Next.js to Go

## Context

Rewrite the production-level single-page portfolio at `./ewerton` (Next.js 16, React 19, next-intl, CSS Modules) into a pure Go application at `./ewerton-go`. The Go version must be visually and functionally identical: same layout, colors (OKLch palette), light/dark theme, i18n (EN/PT-BR), animations, responsive design, SEO, and Docker deployment.

---

## Architecture Decisions

| Decision | Choice | Why |
|---|---|---|
| Framework | `net/http` stdlib (Go 1.26+) | Zero deps, pattern-based routing is sufficient |
| Templating | `html/template` with partials | Built-in, secure, supports `{{define}}`/`{{template}}` |
| CSS | Single hand-written CSS file | No Node.js build step; merge globals.css + all .module.css |
| Icons | Inline SVG as Go template definitions | No extra HTTP request, works with `currentColor` |
| Assets | `go:embed` in binary | Single binary, no filesystem deps |
| i18n | Load existing JSON files at startup | Reuse `messages/en.json` and `messages/pt-br.json` as-is |
| Data | Go structs in `internal/data/` | Type safety, compiled into binary |
| Client JS | Single vanilla JS file (no bundler) | Only needs: typewriter, theme, language switcher, search, modals |
| Docker | 2-stage Alpine build | ~15MB image vs ~150MB Node.js |
| Font | Google Fonts CDN `<link>` | Same as Next.js, no self-hosting |

---

## Project Structure

```
ewerton-go/
  go.mod
  main.go                           # Entry point, routing, embed directives
  Dockerfile
  Makefile

  internal/
    handler/
      page.go                       # Main page handler (locale routing, template render)
      redirects.go                  # /linkedin, /cv redirects
    i18n/
      i18n.go                       # JSON loading, deep merge, dot-path accessor
    data/
      technologies.go               # 17 categories, 100+ skills with ratings/keywords
      education.go                  # Education entries per locale
      partners.go                   # Partner groups and companies
      constants.go                  # URLs, paths, GA ID

  templates/
    layout.html                     # Base HTML skeleton (head, theme script, GA, SEO)
    page.html                       # Main page assembling all partials
    partials/
      navbar.html                   # Sticky nav: logo, language switcher, theme toggle
      hero.html                     # Typewriter greeting, photo, CTAs
      projects.html                 # 2 project cards
      technologies.html             # Search bar + 17 category grids with star ratings
      education.html                # 3 education cards
      partners.html                 # 5 partner groups + company modal
      closing_cta.html              # Calendly CTA
      footer.html                   # Copyright, Brazil flag animation
      modals.html                   # PDF viewer + company info modal shells
    components/
      icons.html                    # ~30 lucide SVG icon definitions as named templates
      star_rating.html              # Star rating partial
      button.html                   # Reusable button partial

  static/
    css/
      globals.css                   # Merged: reset + CSS vars + all component/page styles
    js/
      app.js                        # Typewriter, theme toggle, lang switcher, search, modals
    images/                         # Copy from Next.js public/ (logo, photos, partners, projects)
    documents/
      resume-ewerton.pdf

  messages/
    en.json                         # Copy as-is from Next.js
    pt-br.json                      # Copy as-is from Next.js
```

---

## Implementation Phases

### Phase 1: Foundation
1. `go mod init ewerton-go`
2. Copy static assets from `ewerton/public/` to `ewerton-go/static/images/` and `static/documents/`
3. Copy `messages/en.json` and `messages/pt-br.json` as-is
4. **`internal/i18n/i18n.go`** -- Load both JSON files via `go:embed`, parse into `map[string]any`, implement `deepMerge` (PT-BR over EN fallback), implement `Get(locale, dotPath) string`
5. **`internal/data/constants.go`** -- LinkedIn URL, Aiomover URL, resume path, Calendly URL, GA ID
6. **`internal/data/technologies.go`** -- Port all 17 categories and 100+ skills from `src/data/technologies.ts`
7. **`internal/data/education.go`** -- Port education entries (locale-keyed) from `src/data/education.ts`
8. **`internal/data/partners.go`** -- Port partner groups from page.tsx `PARTNER_GROUPS`

### Phase 2: CSS
1. Create `static/css/globals.css`:
   - Minimal CSS reset (replaces Tailwind preflight): `box-sizing`, `margin: 0`, etc.
   - Copy `:root` and `[data-theme="dark"]` custom properties verbatim from `src/app/globals.css`
   - Copy base element styles (body, h1-h6, a, `.glass`) from `src/app/globals.css`
   - Append all component `.module.css` files (Navbar, Footer, Button, ThemeToggle, LanguageSwitcher, CompanyModal, StarRating)
   - Append `page.module.css` (922 lines) -- class names are already globally unique, no conflicts
   - Remove `@import "tailwindcss"` line (no Tailwind utilities are used in JSX, only CSS Modules)

### Phase 3: Templates
1. **`templates/layout.html`** -- Full `<head>`:
   - Inline theme init script (reads localStorage, sets `data-theme` before paint)
   - Google Fonts `<link>` for Inter
   - CSS `<link>` to `/static/css/globals.css?v={{.Version}}`
   - SEO: `<title>`, OG tags, hrefLang alternates, canonical, JSON-LD Person schema
   - Google Analytics script
   - `{{block "content" .}}{{end}}` in `<body>`
2. **`templates/components/icons.html`** -- ~30 lucide icons as `{{define "icon-xxx"}}` templates with inline SVG paths
3. **`templates/components/star_rating.html`** -- Overlapping filled/unfilled stars with percentage clip
4. **Each partial** in `templates/partials/` -- Port JSX to Go template syntax:
   - Replace `{t("key")}` with `{{call .T "key"}}`
   - Replace `{variable}` with `{{.Variable}}`
   - Replace `className=` with `class=`
   - Replace `{condition && <jsx>}` with `{{if .Condition}}...{{end}}`
   - Replace `.map()` with `{{range .Items}}...{{end}}`
   - Replace `<Image>` with `<img>` (add explicit width/height/alt/loading="lazy")
   - Add `data-*` attributes for JS interactivity (search keywords, company slugs, etc.)
5. **`templates/page.html`** -- Includes all partials in order

### Phase 4: Handlers & Server
1. **`internal/handler/page.go`**:
   - Parse templates with FuncMap (`t`, `seq`, `mul`, `div`, `safeHTML`, etc.)
   - `PageHandler` extracts `{locale}` from URL, validates (en|pt-br), builds `PageData`, renders template
   - `PageData` struct: Locale, OtherLocale, T func, Technologies, Education, PartnerGroups, constants, Version
2. **`internal/handler/redirects.go`** -- `/linkedin` -> LinkedIn, `/cv` -> resume PDF
3. **`main.go`**:
   - `go:embed` directives for `static`, `templates`, `messages`
   - Route registration: `GET /{$}` -> redirect `/en`, `GET /{locale}` -> page handler, `GET /static/` -> file server
   - Listen on `$PORT` (default 3000)

### Phase 5: Client-Side JavaScript
Create `static/js/app.js` as a single IIFE:
1. **Theme toggle** -- Click handler swaps `data-theme`, persists to localStorage, toggles sun/moon icon visibility
2. **Language switcher** -- Dropdown toggle, click-outside-to-close, navigate to `/{newLocale}` on selection
3. **Typewriter effect** -- Port the React `useEffect` state machine: 12 HELLO_WORDS, typing/deleting/pausing states, exact same timing (2000ms pause, 30ms delete, 60ms type)
4. **Skill search** -- `input` listener filters `.techSkillCard` elements by `data-name`/`data-keywords` matching, hides empty categories
5. **PDF viewer modal** -- Show/hide overlay on button click, close on backdrop click
6. **Company info modal** -- Open with matching `data-company-slug` content visible, close on backdrop/button/Escape, body scroll lock

### Phase 6: Docker & Polish
1. **Dockerfile** (2-stage):
   - `golang:1.26-alpine` builder: `go mod download`, `go build` with `-ldflags="-s -w"`
   - `alpine` runner: copy binary only (assets are embedded), non-root user, expose 3000
2. **Makefile**: `dev`, `build`, `docker` targets
3. End-to-end testing:
   - Both locales render correctly (`/en`, `/pt-br`)
   - `/` redirects to `/en`
   - `/linkedin` and `/cv` redirect properly
   - Light/dark theme toggle with no flash on reload
   - Typewriter animation cycles through all greetings
   - Skill search filters correctly
   - All modals open/close properly
   - Responsive layout at 600px and 900px breakpoints
   - All images load (partner logos, photo, project icons)
   - SEO: check `<title>`, OG tags, hrefLang, JSON-LD in page source

---

## Critical Source Files to Port

| Source (Next.js) | Target (Go) | Notes |
|---|---|---|
| `src/app/[locale]/page.tsx` (477 lines) | `templates/partials/*.html` + `static/js/app.js` | JSX -> templates, hooks -> vanilla JS |
| `src/app/globals.css` | `static/css/globals.css` (first section) | Copy CSS vars verbatim |
| `src/app/[locale]/page.module.css` (922 lines) | `static/css/globals.css` (appended) | Copy as-is, class names are unique |
| `src/components/**/*.module.css` | `static/css/globals.css` (appended) | Copy as-is |
| `src/components/**/*.tsx` | `templates/partials/*.html` + `templates/components/*.html` | Port JSX to Go templates |
| `src/data/technologies.ts` | `internal/data/technologies.go` | TS arrays -> Go structs |
| `src/data/education.ts` | `internal/data/education.go` | TS objects -> Go structs |
| `src/i18n/request.ts` | `internal/i18n/i18n.go` | Deep merge + dot-path accessor |
| `src/lib/ThemeProvider.tsx` | Inline `<script>` + `app.js` | React context -> vanilla JS |
| `src/proxy.ts` | `internal/handler/redirects.go` | Middleware -> route handlers |
| `messages/*.json` | `messages/*.json` | Copy as-is |
| `Dockerfile` | `Dockerfile` | Node multi-stage -> Go multi-stage |

---

## Verification

1. `go run .` and open `http://localhost:3000` -- should redirect to `/en`
2. Compare side-by-side with `npm run dev` of the Next.js version at every breakpoint
3. `docker build -t ewerton-go . && docker run -p 3000:3000 ewerton-go`
4. Lighthouse audit for performance, SEO, accessibility parity
5. Test all interactive features: typewriter, theme, language switch, search, modals
