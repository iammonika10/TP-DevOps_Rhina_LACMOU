package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// Cas nominal : makeHandlerFunc renvoie le code et le corps en JSON
func TestMakeHandlerFunc(t *testing.T) {
	handler := makeHandlerFunc(func(req *http.Request) (int, any) {
		return http.StatusCreated, "hello"
	})
	rec := httptest.NewRecorder() // un faux "navigateur" qui enregistre la réponse

	handler(rec, httptest.NewRequest("GET", "/", nil))

	if rec.Code != http.StatusCreated {
		t.Errorf("attendu 201, obtenu %d", rec.Code)
	}
	if rec.Header().Get("content-type") != "application/json" {
		t.Errorf("content-type inattendu : %s", rec.Header().Get("content-type"))
	}
	if strings.TrimSpace(rec.Body.String()) != `"hello"` {
		t.Errorf("corps inattendu : %s", rec.Body.String())
	}
}

// Cas limite : une fonction qui plante donne une erreur 500, sans arrêter le serveur
func TestMakeHandlerFuncPanic(t *testing.T) {
	handler := makeHandlerFunc(func(req *http.Request) (int, any) {
		panic("boum")
	})
	rec := httptest.NewRecorder()

	handler(rec, httptest.NewRequest("GET", "/", nil))

	if rec.Code != http.StatusInternalServerError {
		t.Errorf("attendu 500, obtenu %d", rec.Code)
	}
}

// Les routes déclarées dans newApp mènent aux bonnes fonctions
func TestRoutes(t *testing.T) {
	resetDB()
	server := httptest.NewServer(newApp()) // un vrai serveur de test, sur un port libre
	defer server.Close()

	cases := []struct {
		method, path string
		want         int
	}{
		{"GET", "/", http.StatusOK},
		{"GET", "/swagger/", http.StatusOK},
		{"GET", "/api/cats", http.StatusOK},
		{"GET", "/api/cats/id1", http.StatusOK},
		{"DELETE", "/api/cats/inconnu", http.StatusNotFound},
		{"PUT", "/api/cats", http.StatusMethodNotAllowed},
		{"GET", "/api/rien", http.StatusNotFound},
	}

	for _, c := range cases {
		req, _ := http.NewRequest(c.method, server.URL+c.path, nil)
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("%s %s : erreur %v", c.method, c.path, err)
		}
		res.Body.Close()
		if res.StatusCode != c.want {
			t.Errorf("%s %s : attendu %d, obtenu %d", c.method, c.path, c.want, res.StatusCode)
		}
	}
}
