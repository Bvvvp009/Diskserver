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
	if r.Method == http.MethodPost {
		username := r.FormValue("username")
		password := r.FormValue("password")

		if password == Users[username] {
			token, err := auth.GenerateJWT(username)
			if err != nil {
				http.Error(w, "Error generating token", http.StatusInternalServerError)
				return
			}

			http.SetCookie(w, &http.Cookie{
				Name:    "token",
				Value:   token,
				Expires: time.Now().Add(30 * time.Minute),
			})

			http.Redirect(w, r, "/", http.StatusSeeOther)
		} else {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		}
	} else {
		tmpl := template.Must(template.ParseFiles("../../static/login.html"))
		tmpl.Execute(w, nil)
	}
}
