package main

import (
	"context"
	"html/template"
	"io"
	"log"
	"net/http"
)

type Templates struct {
	templates *template.Template
}

func (t *Templates) Render(w io.Writer, name string, data interface{}, c context.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func newTemplate() *Templates {
	return &Templates{
		templates: template.Must(template.ParseGlob("*.html")),
	}
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Request: %s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

func homeHandler(t *Templates, data *Data) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := t.Render(w, "index", data, r.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func dataHandler(t *Templates, data *Data) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		err := t.Render(w, "data", data, r.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func contactHandler(t *Templates, data *Data) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}

		name := r.FormValue("name")
		email := r.FormValue("email")

		if name == "" || email == "" {
			http.Error(w, "name and email are required", http.StatusBadRequest)
			return
		}

		data.Contacts = append(data.Contacts, newContact(name, email))
		err := t.Render(w, "display", data, r.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

type Count struct {
	Count int
}

type Contact struct {
	Name  string
	Email string
}

func newContact(name, email string) Contact {
	return Contact{Name: name, Email: email}
}

type Contacts = []Contact

type Data struct {
	Contacts Contacts
}

func (d *Data) hasEmail(email string) bool {
	for _, contact := range d.Contacts {
		if contact.Email == email {
			return true
		}
	}
	return false
}

func newData() Data {
	return Data{
		Contacts: []Contact{
			newContact("fafifu", "fafifu@gmail.com"),
			newContact("ytta", "ytta@gmail.com"),
		},
	}
}

func main() {
	templates := newTemplate()

	data := newData()

	// Set up the HTTP server
	mux := http.NewServeMux()
	mux.Handle("/", loggingMiddleware(http.HandlerFunc(homeHandler(templates, &data))))
	mux.Handle("/contacts", loggingMiddleware(http.HandlerFunc(contactHandler(templates, &data))))

	// Start the server
	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
