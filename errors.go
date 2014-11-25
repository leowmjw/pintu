package pintu

import (
	"errors"
	"log"
	"net/http"
)

type ErrorResponse struct {
	Title     string
	Message   string
	LoginPath string
}

var (
	ErrInvalidCredentials = errors.New("Invalid credentials")
	ErrAuthServerDown     = errors.New("Authentication server offline")
)

// Error compiles error response
func DefaultError(w http.ResponseWriter, r *http.Request, code int, title string, message string) {
	log.Printf("ErrorPage %d %s %s", code, title, message)

	w.WriteHeader(code)

	templates := GetTemplates()
	templates.ExecuteTemplate(w, "error.html", &ErrorResponse{
		Title:     title,
		Message:   message,
		LoginPath: GetHostPath(r, loginPromptPath),
	})
}

// CustomError checks for error and throws a 500 error
func CustomError(w http.ResponseWriter, r *http.Request, err error) {
	DefaultError(w, r, 500, "Internal Error", err.Error())
}

// Denied throws a 404 page
func Denied(w http.ResponseWriter, r *http.Request) {
	DefaultError(w, r, 403, "Access Denied", "Please login")
}

// NotFound throws a 404 page
func NotFound(w http.ResponseWriter, r *http.Request) {
	DefaultError(w, r, 404, "Not Found", "Page not found")
}
