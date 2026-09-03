package data

import "time"

const (
	LinkedInURL    = "https://www.linkedin.com/in/ewertones"
	AiomoverURL    = "https://aiomover.com"
	QuestaoURL     = "https://questao.com.br"
	HarpaCristaURL = "https://play.google.com/store/apps/details?id=com.ewertones.harpacrista"
	ResumePath     = "/static/documents/resume-ewerton.pdf"
	CalendlyURL    = "https://calendly.com/ewertones/30min"
	GoogleAnalID   = "G-XXXXXXXXXX"
)

// LaunchDeadline is when questao and Harpa Cristã go live on the site.
// Both links stay disabled until this moment, then unlock automatically.
var LaunchDeadline = time.Date(2026, time.September, 15, 0, 0, 0, 0, time.FixedZone("BRT", -3*60*60))
