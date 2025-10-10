package main

import (
	"fmt"
	"html/template"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestNewContact(t *testing.T) {
	name := "Ramaido"
	email := "ramaido@fafifu.com"

	contact := newContact(name, email)

	if contact.Name != name {
		t.Errorf("Expected name %s, got %s", name, contact.Name)
	}
	if contact.Email != email {
		t.Errorf("Expected email %s, got %s", email, contact.Email)
	}
}

func TestNewData(t *testing.T) {
	data := newData()

	if len(data.Contacts) != 2 {
		t.Errorf("Expected 2 contacts, got %d", len(data.Contacts))
	}

	if data.Contacts[0].Name != "fafifu" {
		t.Errorf("Expected first contact name 'fafifu', got '%s'", data.Contacts[0].Name)
	}
}

func TestHasEmail(t *testing.T) {
	data := newData()

	tests := []struct {
		email    string
		expected bool
	}{
		{"fafifu@gmail.com", true},
		{"ytta@gmail.com", true},
		{"nonexistent@gmail.com", false},
		{"", false},
	}

	for _, tt := range tests {
		result := data.hasEmail(tt.email)
		if result != tt.expected {
			t.Errorf("hasEmail(%s) = %v, expected %v", tt.email, result, tt.expected)
		}
	}
}

func TestHomeHandler(t *testing.T) {
	// Setup template for testing with mock
	tmpl := template.Must(template.New("index").Parse(`<html><body>Home</body></html>`))
	templates := &Templates{templates: tmpl}
	data := newData()

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	handler := homeHandler(templates, &data)
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status OK, got %v", rec.Code)
	}

	if !strings.Contains(rec.Body.String(), "Home") {
		t.Error("Expected response to contain 'Home'")
	}
}

func TestContactHandlerInvalidMethod(t *testing.T) {
	templates := &Templates{}
	data := newData()

	req := httptest.NewRequest(http.MethodGet, "/contacts", nil)
	rec := httptest.NewRecorder()

	handler := contactHandler(templates, &data)
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected status MethodNotAllowed, got %v", rec.Code)
	}
}

func TestContactHandlerEmptyData(t *testing.T) {
	tmpl := template.Must(template.New("display").Parse(`<div>{{.}}</div>`))
	templates := &Templates{templates: tmpl}
	data := newData()

	form := url.Values{}
	form.Add("name", "")
	form.Add("email", "")

	req := httptest.NewRequest(http.MethodPost, "/contacts", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()

	handler := contactHandler(templates, &data)
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("Expected status BadRequest, got %v", rec.Code)
	}
}

func TestLoggingMiddleware(t *testing.T) {
	called := false

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	})

	middleware := loggingMiddleware(nextHandler)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()

	middleware.ServeHTTP(rec, req)

	if !called {
		t.Error("Expected next handler to be called")
	}

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status OK, got %v", rec.Code)
	}
}

// Benchmark untuk contactHandler
func BenchmarkContactHandler(b *testing.B) {
	templates := &Templates{}
	data := newData()

	form := url.Values{}
	form.Add("name", "Benchmark User")
	form.Add("email", "bench@example.com")

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest(http.MethodPost, "/contacts", strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		rec := httptest.NewRecorder()

		handler := contactHandler(templates, &data)
		handler.ServeHTTP(rec, req)
	}
}

func TestClearContactsHandler(t *testing.T) {
	tmpl := template.Must(template.New("display").Parse(`<div>{{len .Contacts}}</div>`))
	templates := &Templates{templates: tmpl}
	data := newData()

	if len(data.Contacts) == 0 {
		t.Fatal("Expected initial contacts to be non-empty")
	}

	req := httptest.NewRequest(http.MethodPost, "/clear-contacts", nil)
	rec := httptest.NewRecorder()

	handler := clearContactsHandler(templates, &data)
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status OK, got %v", rec.Code)
	}

	if len(data.Contacts) != 0 {
		t.Errorf("Expected contacts cleared, got %d remaining", len(data.Contacts))
	}
}

func TestContactCountHandler(t *testing.T) {
	data := newData()
	handler := contactCountHandler(&data)

	req := httptest.NewRequest(http.MethodGet, "/contact-count", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("Expected status OK, got %v", rec.Code)
	}

	body := rec.Body.String()
	if !strings.Contains(body, `"count":2`) {
		t.Errorf("Expected count=2 in response, got %s", body)
	}
	if !strings.Contains(body, `"maxLimit":10`) {
		t.Errorf("Expected maxLimit=10 in response, got %s", body)
	}
}

func TestContactHandlerAddSuccess(t *testing.T) {
	tmpl := template.Must(template.New("display").Parse(`<div>{{range .Contacts}}{{.Name}} {{end}}</div>`))
	templates := &Templates{templates: tmpl}
	data := newData()

	form := url.Values{}
	form.Add("name", "fafifu")
	form.Add("email", "fafifu@gmail.com")

	req := httptest.NewRequest(http.MethodPost, "/contacts", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()

	handler := contactHandler(templates, &data)
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status OK, got %v", rec.Code)
	}

	if !data.hasEmail("fafifu@gmail.com") {
		t.Error("Expected new contact to be added")
	}
}

func TestContactHandlerMaxLimit(t *testing.T) {
	tmpl := template.Must(template.New("display").Parse(`<div>{{.}}</div>`))
	templates := &Templates{templates: tmpl}

	// Isi data dengan 10 kontak
	data := Data{}
	for i := 0; i < 10; i++ {
		data.Contacts = append(data.Contacts, newContact(
			"User"+string(rune(i+'A')),
			fmt.Sprintf("user%d@example.com", i),
		))
	}

	form := url.Values{}
	form.Add("name", "Fafifu User")
	form.Add("email", "fafifu@gmail.com")

	req := httptest.NewRequest(http.MethodPost, "/contacts", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()

	handler := contactHandler(templates, &data)
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("Expected status BadRequest when limit reached, got %v", rec.Code)
	}
}

func TestDataHandler(t *testing.T) {
	tmpl := template.Must(template.New("data").Parse(`<div>{{len .Contacts}}</div>`))
	templates := &Templates{templates: tmpl}
	data := newData()

	req := httptest.NewRequest(http.MethodGet, "/data", nil)
	rec := httptest.NewRecorder()

	handler := dataHandler(templates, &data)
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status OK, got %v", rec.Code)
	}

	if !strings.Contains(rec.Body.String(), "2") {
		t.Errorf("Expected contact count rendered, got %s", rec.Body.String())
	}
}
