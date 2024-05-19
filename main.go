package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	initDB()

	http.HandleFunc("/", indexHandler)
	http.HandleFunc("POST /paste", createPasteHandler)
	http.HandleFunc("GET /paste/{id}", readPasteHandler)
	http.HandleFunc("POST /paste/{id}", readPasteHandler)
	http.HandleFunc("GET /raw/{id}", readRawHandler)
	http.HandleFunc("POST /paste/{id}/delete", deletePasteHandler)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	log.Println("Starting server on https://0.0.0.0:443...")
	log.Fatal(http.ListenAndServeTLS(":443", "certs/server.crt", "certs/server.key", nil))
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	pastes, err := loadAllPastes()
	if err != nil {
		http.Error(w, "Unable to load pastes", http.StatusInternalServerError)
		return
	}

	data := struct {
		Pastes []*Paste
	}{Pastes: pastes}

	renderTemplate(w, data, "templates/index.html", "templates/base.html")
}

func readPasteHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid paste ID", http.StatusBadRequest)
		return
	}

	paste, err := loadPasteByID(id)
	if err != nil {
		error := struct {
			Type    string
			Message string
		}{
			Type:    "Not Found",
			Message: "This paste is no longer available.",
		}
		renderTemplate(w, error, "templates/error.html", "templates/base.html")
		return
	}

	if paste.Password != "" {
		password := r.FormValue("password")
		if password == "" {
			renderTemplate(w, struct{ ID int }{ID: paste.ID}, "templates/password.html", "templates/base.html")
			return
		}
		if err := bcrypt.CompareHashAndPassword([]byte(paste.Password), []byte(password)); err != nil {
			http.Error(w, "Invalid password", http.StatusUnauthorized)
			return
		}
	}

	renderTemplate(w, struct{ Paste *Paste }{Paste: paste}, "templates/paste.html", "templates/base.html")
}

func readRawHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid paste ID", http.StatusBadRequest)
		return
	}

	paste, err := loadPasteByID(id)
	if err != nil {
		http.Error(w, "Unable to load paste", http.StatusInternalServerError)
		return
	}

	if paste.Password != "" {
		// Password protected paste is not allowed to be viewed in raw form
		http.Redirect(w, r, fmt.Sprintf("/paste/%d", paste.ID), http.StatusSeeOther)
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Write([]byte(paste.Content))
}

func createPasteHandler(w http.ResponseWriter, r *http.Request) {
	// retrieve form values or GET parameters values
	content := r.FormValue("content")
	title := r.FormValue("title")
	if title == "" {
		title = "Untitled"
	}
	password := r.FormValue("password")

	paste := &Paste{
		Title:    title,
		Content:  content,
		Password: password,
	}

	if err := paste.Save(); err != nil {
		http.Error(w, "Unable to save paste", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/paste/%d", paste.ID), http.StatusSeeOther)
}

func deletePasteHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid paste ID", http.StatusBadRequest)
		return
	}

	paste, err := loadPasteByID(id)
	if err != nil {
		http.Error(w, "Unable to load paste", http.StatusInternalServerError)
		return
	}

	if err := paste.Delete(); err != nil {
		http.Error(w, "Unable to delete paste", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func renderTemplate(w http.ResponseWriter, data interface{}, files ...string) {
	tmpl := template.Must(template.ParseFiles(files...))
	if err := tmpl.ExecuteTemplate(w, "main", data); err != nil {
		http.Error(w, "Unable to render template", http.StatusInternalServerError)
	}
}
