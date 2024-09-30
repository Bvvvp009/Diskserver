package handlers

import (
	"diskserver/internal/auth"
	"html/template"
	"net/http"
	"time"
)

var Users = map[string]string{
	"user1": "password1",
	"user2": "password2",
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	// Check if a valid token already exists
	cookie, err := r.Cookie("token")
	if err == nil {
		// Validate the token
		_, err := auth.ValidateJWT(cookie.Value)
		if err == nil {
			// Token is valid, redirect to homepage
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
	}

	// If it's a POST request (login attempt)
	if r.Method == http.MethodPost {
		username := r.FormValue("username")
		password := r.FormValue("password")

		// Verify the username and password
		if password == Users[username] {
			// Generate a JWT token
			token, err := auth.GenerateJWT(username)
			if err != nil {
				http.Error(w, "Error generating token", http.StatusInternalServerError)
				return
			}

			// Set the token in a cookie
			http.SetCookie(w, &http.Cookie{
				Name:    "token",
				Value:   token,
				Expires: time.Now().Add(300 * time.Minute),
			})

			// Redirect to the homepage
			http.Redirect(w, r, "/", http.StatusSeeOther)
		} else {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		}
	} else {
		// Render the login page if the request is GET and no valid token exists
		tmpl := template.Must(template.ParseFiles("../../static/login.html"))
		tmpl.Execute(w, nil)
	}
}
