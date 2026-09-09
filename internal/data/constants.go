package data

import "time"

const (
	LinkedInHandle = "ewertones"
	LinkedInURL    = "https://www.linkedin.com/in/ewertones"
	GitHubHandle   = "ewertones"
	GitHubURL      = "https://github.com/ewertones"
	AiomoverURL    = "https://aiomover.com"
	QuestaoURL     = "https://questao.com.br"
	HarpaCristaURL = "https://play.google.com/store/apps/details?id=com.ewertones.harpacrista"
	CalendlyURL    = "https://calendly.com/ewertones/30min"
	GoogleAnalID   = "G-MFH72L1LCE"
)

// Identity and contact details. These are the only resume facts that do not
// already exist somewhere on the landing page.
const (
	FullName     = "Ewerton de Souza"
	Email        = "ewerton@ewerton.com.br"
	PhoneDisplay = "+55 31 9 9854 6167"
	PhoneE164    = "+5531998546167"
	City         = "Belo Horizonte"
	Region       = "MG"
)

// SiteOrigin is the canonical origin, used to build absolute URLs in the
// sitemap, JSON-LD and generated documents.
const (
	SiteOrigin = "https://ewerton.com.br"
	SiteURL    = SiteOrigin
)

// Resume artifact paths. ResumePath is what the "View Resume" button opens; the
// PDF is now generated per request rather than served from disk.
const (
	ResumePath     = "/resume.pdf"
	ResumeJSONPath = "/resume.json"
	ResumeMDPath   = "/resume.md"
)

// LaunchDeadline is when questao and Harpa Cristã go live on the site.
// Both links stay disabled until this moment, then unlock automatically.
var LaunchDeadline = time.Date(2026, time.September, 15, 0, 0, 0, 0, time.FixedZone("BRT", -3*60*60))
