# Ewerton Portfolio (Go Migration)

This repository contains the professional portfolio website migrated from Next.js (React/TypeScript) to a **pure Go standard library** application.

## 🎯 Architecture Goals

The main objectives behind this rewrite are:

1. **Performance & Lightweight Memory:** Drop the heavy Node.js runtime and Next.js hydration bundle footprint. Transitioning to Go results in an incredibly fast site with a tiny server memory footprint.
2. **Zero External Dependencies:** Built entirely using Go's standard library (`net/http`, `html/template`, `embed`).
3. **Single Binary:** Utilizing `//go:embed`, all assets (Vanilla CSS, Javascript, Images, SVGs, and i18n JSONs) are compiled directly into a single static executable.
4. **Efficient Containerization:** A clean two-stage Alpine Dockerfile drops the final production image weight drastically, making cloud deployments (like Google Cloud Run) extremely resource and cost-efficient.

## 📁 Project Structure

```text
ewerton-go/
├── main.go                     # Entrypoint configuring Go embed file systems & HTTP listener
├── Makefile                    # Developer aliases (dev, build, and docker directives)
├── Dockerfile                  # Alpine Two-Stage Docker runner
├── messages/                   # i18n JSON translations (en, pt-br)
├── static/                     # All static assets
│   ├── css/                    # Pure CSS architecture heavily utilizing fluid OKLCH coloring
│   ├── js/                     # Vanilla Javascript (Typewriter, Themes, Search, Modals)
│   └── images/                 # Optimized webp/svg photography & partner logos, plus
│                               # derived rasters (see scripts/generate_assets.py)
├── scripts/                    # One-off asset generation, kept for reproducibility
├── templates/                  # Base layouts and complex component partials (HTML)
└── internal/
    ├── data/                   # Strictly-typed Go structs mimicking the old TypeScript data models
    │                           # plus resume.go, the canonical profile model
    ├── handler/                # Main router handling template rendering, localization, and redirects
    ├── i18n/                   # JSON dictionary utility for managing deep recursive translation trees
    ├── pdf/                    # Minimal PDF 1.7 writer (stdlib only): text, vectors, gradients,
    │                           # JPEG, links, outlines, XMP metadata, embedded attachments
    └── resume/                 # Renders the profile model as PDF, JSON Resume, Markdown and JSON-LD
```

## 📄 The Live Resume

The resume is **generated on request from the same data the landing page renders**.
There is no PDF checked into the repository.

Work history comes from the partner cards joined to their `messages/*.json`
descriptions, education from `data.EducationData`, skills and their ratings from
`data.TechnologyCategories`, and the products from the project copy. Edit a role
in `messages/en.json` and the PDF, the JSON and the structured data all change
with it; they cannot drift apart.

### Formats

| Route | What it is |
|---|---|
| `/resume.pdf`, `/{locale}/resume.pdf` | Designed PDF, generated per request. `?download=1` forces a save dialog. |
| `/resume.json`, `/{locale}/resume.json` | [JSON Resume](https://jsonresume.org) schema v1.0.0 |
| `/resume.md`, `/{locale}/resume.md` | Markdown rendering |
| `/llms.txt`, `/llms-full.txt` | [llmstxt.org](https://llmstxt.org) site guide, and the whole resume inline |
| `/robots.txt`, `/sitemap.xml` | Crawler discovery, with AI user-agents welcomed by name |

Each artifact is rendered once per locale per day, cached in memory with a
SHA-256 `ETag`, and served through `http.ServeContent` so repeat visits get a
`304`.

### Derived assets

Three images cannot be produced by the Go binary: the portrait is authored as WebP
(which PDF cannot embed), the logo is a 275-path Illustrator SVG, and the Open Graph
cards need a text renderer. `scripts/generate_assets.py` renders all of them so the
committed files are reproducible rather than mystery binaries:

```bash
python scripts/generate_assets.py
```

### Why it is built this way

The previous `resume-ewerton.pdf` was a 1.38 MB Photoshop export: a flat JPEG in
a PDF wrapper, with zero embedded fonts. No applicant tracking system, AI
crawler or screen reader could read a single word of it. The generated document
is real text — extractable, searchable, linkable — and about 90 KB.

- **Pure Go, still zero dependencies.** `internal/pdf` writes PDF 1.7 by hand,
  using the base-14 Helvetica faces so nothing needs embedding.
- **Machine-readable inside and out.** The JSON Resume document is embedded
  *inside* the PDF as an associated file (`/AFRelationship /Source`), so the
  structured data travels with the document.
- **Accessible.** The document is *tagged*: every run of text sits in a marked-content
  sequence collected into a `/StructTreeRoot`, with `H1`/`H2`/`H3`/`P` roles, alternative
  text on the portrait, and all decoration (rules, tints, stars, page furniture) declared
  `/Artifact` so assistive technology skips it. Plus document `/Lang`, a real title,
  section bookmarks, live link annotations and an XMP metadata packet.
- **Structured data.** The landing page embeds a generated schema.org `@graph`
  (`ProfilePage` + `Person` with employment history, credentials, `knowsAbout`
  and the products), built from the same model as the PDF.

## 🚀 Running Locally

Development requires [Go](https://go.dev/doc/install) (1.20+) and standard CLI tools (`make`/`docker`).

### 1. Dev Server

The `dev` target spins up the local testing environment directly from source files.

```bash
make dev
```

**Access the app at:** `http://localhost:3000`

### 2. Standard Binaries

Generates a fully stripped, standalone executable tightly bundled avoiding any external disk references.

```bash
make build
./ewerton-go
```

### 3. Docker Environment

To guarantee environment parity with the pipeline production releases (Google Cloud Run), build the sandbox image.

```bash
make docker
# Behind the scenes this runs: docker build -t ewerton-go . & docker run -p 3000:3000 -d
```

## 🛠 Features Implemented

- **Pixel-Perfect Port:** HTML semantic markup was mirrored component-by-component exactly from the original Next.js codebase.
- **Client-Side Vanilla JS:** Features like the Typewriter loop, Dark/Light mode theme persistence, real-time client-side skill searching, and company role modal handlers run completely in vanilla JavaScript IIFEs, dropping megabytes of React overhead.
- **Advanced Dynamic CSS Context:** Merged disparate React CSS Modules into one tightly formatted and cacheable native CSS object layout logic tracking Light/Dark variables successfully natively instead of relying on heavy CSS-in-JS dependencies like Tailwind.
- **Dynamic Extensible i18n:** Replicated identical user-facing paths (`/en` `<->` `/pt-br`) dynamically mapped off JSON definitions.
- **Runtime-Generated Resume:** A hand-rolled PDF writer renders a designed, fully text-based resume from live site data, alongside JSON Resume, Markdown and schema.org representations. See [The Live Resume](#-the-live-resume).
