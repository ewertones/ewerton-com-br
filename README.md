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
│   └── images/                 # Optimized webp/svg photography & partner logos
├── templates/                  # Base layouts and complex component partials (HTML)
└── internal/
    ├── data/                   # Strictly-typed Go structs mimicking the old TypeScript data models
    ├── handler/                # Main router handling template rendering, localization, and redirects
    └── i18n/                   # JSON dictionary utility for managing deep recursive translation trees
```

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
# Behind the scenes this runs: docker build -t ewerton-go . & docker run -p 3000:3000
```

## 🛠 Features Implemented

- **Pixel-Perfect Port:** HTML semantic markup was mirrored component-by-component exactly from the original Next.js codebase.
- **Client-Side Vanilla JS:** Features like the Typewriter loop, Dark/Light mode theme persistence, real-time client-side skill searching, and company role modal handlers run completely in vanilla JavaScript IIFEs, dropping megabytes of React overhead.
- **Advanced Dynamic CSS Context:** Merged disparate React CSS Modules into one tightly formatted and cacheable native CSS object layout logic tracking Light/Dark variables successfully natively instead of relying on heavy CSS-in-JS dependencies like Tailwind.
- **Dynamic Extensible i18n:** Replicated identical user-facing paths (`/en` `<->` `/pt-br`) dynamically mapped off JSON definitions.
