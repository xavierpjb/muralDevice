package update

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/spf13/afero"
)

func TestMain(m *testing.M) {
	setupAll()
	code := m.Run()
	teardownAll()
	os.Exit(code)
}

func setupAll() {
	fmt.Println("setup")

}

func teardownAll() {
	fmt.Println("teardown")

}

func TestPostUpdate(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	updateHandler := New(afero.NewMemMapFs())

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(updateHandler.HandleUpdate)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

}
