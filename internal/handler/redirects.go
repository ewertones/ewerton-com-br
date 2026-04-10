package handler

import (
	"net/http"

	"ewerton-go/internal/data"
)

// LinkedInRedirect sends the user to the LinkedIn profile.
func LinkedInRedirect(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, data.LinkedInURL, http.StatusFound)
}

// CVRedirect sends the user to the resume PDF.
func CVRedirect(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, data.ResumePath, http.StatusFound)
}
