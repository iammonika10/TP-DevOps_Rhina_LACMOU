package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// newCatRequest prépare une requête comme si le routeur avait lu {catId} dans l'adresse
func newCatRequest(method, catID string) *http.Request {
	req := httptest.NewRequest(method, "/api/cats/"+catID, nil)
	req.SetPathValue("catId", catID)
	return req
}

// Cas nominal : le chat existe
func TestGetCat(t *testing.T) {
	resetDB()

	code, body := getCat(newCatRequest("GET", "id1"))

	if code != http.StatusOK {
		t.Errorf("attendu 200, obtenu %d", code)
	}
	cat := body.(Cat)
	if cat.Name != "Toto" {
		t.Errorf("attendu Toto, obtenu %s", cat.Name)
	}
}

// Cas limite : le chat n'existe pas
func TestGetCatNotFound(t *testing.T) {
	resetDB()

	code, _ := getCat(newCatRequest("GET", "inconnu"))

	if code != http.StatusNotFound {
		t.Errorf("attendu 404, obtenu %d", code)
	}
}

// Cas nominal : la suppression retire bien le chat
func TestDeleteCat(t *testing.T) {
	resetDB()

	code, body := deleteCat(newCatRequest("DELETE", "id1"))

	if code != http.StatusNoContent {
		t.Errorf("attendu 204, obtenu %d", code)
	}
	if body != nil {
		t.Errorf("attendu aucun corps, obtenu %v", body)
	}
	if _, found := catsDatabase["id1"]; found {
		t.Error("id1 est toujours dans la base après DELETE")
	}
}

// Cas limite : supprimer un chat inconnu ne touche pas à la base
func TestDeleteCatNotFound(t *testing.T) {
	resetDB()

	code, _ := deleteCat(newCatRequest("DELETE", "inconnu"))

	if code != http.StatusNotFound {
		t.Errorf("attendu 404, obtenu %d", code)
	}
	if len(catsDatabase) != 1 {
		t.Errorf("la base devrait encore contenir Toto, elle contient %d chats", len(catsDatabase))
	}
}
