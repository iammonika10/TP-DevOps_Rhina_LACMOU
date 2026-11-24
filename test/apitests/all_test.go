package apitests

import (
	"fmt"
	"net/http"
	"testing"
)

var initCatId string

func init() {
	// Preparation: delete all existing & create a cat
	ids := []string{}
	call("GET", "/cats", nil, nil, &ids)

	for _, id := range ids {
		code := 0
		call("DELETE", "/cats/"+id, nil, &code, nil)
		fmt.Println("DELETE /cats ->", code)
	}

	// Create a single cat into the DB
	call("POST", "/cats", &CatModel{Name: "Toto"}, nil, &initCatId)
}

func TestGetCats(t *testing.T) {

	code := 0
	result := []string{}
	err := call("GET", "/cats", nil, &code, &result)
	if err != nil {
		t.Error("Request error", err)
	}

	fmt.Println("GET /cats ->", code, result)

	if code != http.StatusOK {
		t.Error("We should get code 200, got", code)
	}

	// DELETE fonctionne : init() a tout supprimé puis créé 1 seul chat
	if len(result) != 1 {
		t.Error("We should get one item, got", len(result))
		return
	}

	if result[0] != initCatId {
		t.Error("Listing the IDs, got", result[0])
	}
}

// GET /cats/{id} sur un chat qui existe -> 200 + le chat
func TestGetCat(t *testing.T) {

	code := 0
	cat := CatModel{}
	err := call("GET", "/cats/"+initCatId, nil, &code, &cat)
	if err != nil {
		t.Error("Request error", err)
	}

	fmt.Println("GET /cats/"+initCatId, "->", code, cat)

	if code != http.StatusOK {
		t.Error("We should get code 200, got", code)
	}

	if cat.Name != "Toto" || cat.ID != initCatId {
		t.Error("We should get the cat Toto, got", cat)
	}
}

// GET /cats/{id} sur un id inconnu -> 404
func TestGetCatNotFound(t *testing.T) {

	code := 0
	err := call("GET", "/cats/unknown-id", nil, &code, nil)
	if err != nil {
		t.Error("Request error", err)
	}

	fmt.Println("GET /cats/unknown-id ->", code)

	if code != http.StatusNotFound {
		t.Error("We should get code 404, got", code)
	}
}

// DELETE /cats/{id} : 204, puis 404 au GET suivant, puis 404 au 2e DELETE
func TestDeleteCat(t *testing.T) {

	// Preparation: create a cat to delete
	catId := ""
	err := call("POST", "/cats", &CatModel{Name: "Felix"}, nil, &catId)
	if err != nil {
		t.Error("Request error", err)
		return
	}

	code := 0
	call("DELETE", "/cats/"+catId, nil, &code, nil)
	fmt.Println("DELETE /cats/"+catId, "->", code)
	if code != http.StatusNoContent {
		t.Error("We should get code 204, got", code)
	}

	call("GET", "/cats/"+catId, nil, &code, nil)
	if code != http.StatusNotFound {
		t.Error("After DELETE, GET should give 404, got", code)
	}

	call("DELETE", "/cats/"+catId, nil, &code, nil)
	if code != http.StatusNotFound {
		t.Error("Deleting twice should give 404, got", code)
	}
}

// POST /cats avec un chat valide -> 201 + l'id, puis le chat est bien lisible
func TestCreateCat(t *testing.T) {

	code := 0
	catId := ""
	err := call("POST", "/cats", &CatModel{Name: "Minou", Color: "black"}, &code, &catId)
	if err != nil {
		t.Error("Request error", err)
		return
	}

	fmt.Println("POST /cats ->", code, catId)

	if code != http.StatusCreated {
		t.Error("We should get code 201, got", code)
	}
	if catId == "" {
		t.Error("We should get the new cat id")
		return
	}

	// Vérification : le chat créé est bien enregistré avec les bonnes données
	cat := CatModel{}
	call("GET", "/cats/"+catId, nil, &code, &cat)
	if code != http.StatusOK || cat.Name != "Minou" || cat.Color != "black" {
		t.Error("The created cat should be readable, got", code, cat)
	}

	// Nettoyage : on supprime le chat pour ne pas gêner les autres tests
	call("DELETE", "/cats/"+catId, nil, nil, nil)
}

// POST /cats avec un corps qui n'est pas un chat -> 400
func TestCreateCatInvalidInput(t *testing.T) {

	code := 0
	// call() encode "pas un chat" en JSON : c'est un texte, pas un objet chat
	err := call("POST", "/cats", "pas un chat", &code, nil)
	if err != nil {
		t.Error("Request error", err)
	}

	fmt.Println("POST /cats (invalid) ->", code)

	if code != http.StatusBadRequest {
		t.Error("We should get code 400, got", code)
	}
}

// Un verbe non prévu par la spec -> 405
func TestMethodNotAllowed(t *testing.T) {

	code := 0
	err := call("PUT", "/cats", nil, &code, nil)
	if err != nil {
		t.Error("Request error", err)
	}

	if code != http.StatusMethodNotAllowed {
		t.Error("We should get code 405, got", code)
	}
}
