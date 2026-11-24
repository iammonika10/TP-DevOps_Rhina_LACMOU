package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// resetDB remet la "base de données" dans un état connu avant chaque test,
// pour que les tests ne dépendent pas les uns des autres.
func resetDB() {
	catsDatabase = map[string]Cat{
		"id1": {Name: "Toto", Color: "Grey", BirthDate: "2023-04-16"},
	}
}

// Cas nominal : une map de 2 chats donne 2 identifiants
func TestListMapKeys(t *testing.T) {
	// 1. Préparer l'entrée
	input := map[string]Cat{"a": {Name: "A"}, "b": {Name: "B"}}

	// 2. Appeler la fonction testée
	keys := listMapKeys(input)

	// 3. Vérifier la sortie
	if len(keys) != 2 {
		t.Errorf("attendu 2 identifiants, obtenu %d", len(keys))
	}
}

// Cas limite : une map vide donne une liste vide (et pas nil, qui deviendrait "null" en JSON)
func TestListMapKeysEmpty(t *testing.T) {
	keys := listMapKeys(map[string]Cat{})

	if keys == nil || len(keys) != 0 {
		t.Errorf("attendu une liste vide [], obtenu %v", keys)
	}
}

func TestListCats(t *testing.T) {
	resetDB()
	req := httptest.NewRequest("GET", "/api/cats", nil)

	code, body := listCats(req)

	if code != http.StatusOK {
		t.Errorf("attendu 200, obtenu %d", code)
	}
	ids := body.([]string)
	if len(ids) != 1 || ids[0] != "id1" {
		t.Errorf("attendu [id1], obtenu %v", ids)
	}
}

// Cas nominal : un JSON valide crée un chat
func TestCreateCat(t *testing.T) {
	resetDB()
	req := httptest.NewRequest("POST", "/api/cats", strings.NewReader(`{"name":"Felix","color":"blue"}`))

	code, body := createCat(req)

	if code != http.StatusCreated {
		t.Errorf("attendu 201, obtenu %d", code)
	}
	newID := body.(string)
	saved, found := catsDatabase[newID]
	if !found {
		t.Fatalf("le chat %s n'a pas été enregistré", newID)
	}
	if saved.Name != "Felix" || saved.ID != newID {
		t.Errorf("chat enregistré incorrect : %+v", saved)
	}
}

// Cas limite : un JSON cassé est refusé et rien n'est enregistré
func TestCreateCatInvalidJSON(t *testing.T) {
	resetDB()
	req := httptest.NewRequest("POST", "/api/cats", strings.NewReader(`{"name": `))

	code, body := createCat(req)

	if code != http.StatusBadRequest {
		t.Errorf("attendu 400, obtenu %d", code)
	}
	if body != "Invalid JSON input" {
		t.Errorf("message inattendu : %v", body)
	}
	if len(catsDatabase) != 1 {
		t.Errorf("la base ne devrait contenir que Toto, elle contient %d chats", len(catsDatabase))
	}
}

// Cas limite : un corps vide est refusé
func TestCreateCatEmptyBody(t *testing.T) {
	resetDB()
	req := httptest.NewRequest("POST", "/api/cats", strings.NewReader(""))

	code, _ := createCat(req)

	if code != http.StatusBadRequest {
		t.Errorf("attendu 400, obtenu %d", code)
	}
}
