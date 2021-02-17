package mural

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/spf13/afero"
)

// Software contains the properties to send to mural when connecting
type Software struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	UUID        string `json:"uuid"`
	Version     string `json:"version"`
}

// MuralModel contains information relating to the device hosting the mural code
type MuralModel struct {
	Software Software
}

// New Creates a new instance of mural base on the file provided
func New(file afero.File) *MuralModel {
	// read from filesystem then instantiate obj
	defer file.Close()

	byteVals, _ := ioutil.ReadAll(file)

	var mural MuralModel
	var software Software

	json.Unmarshal(byteVals, &software)
	mural.Software = software
	return &mural
}

// GetSoftwareSummary Returns the json of mural model
func (mu MuralModel) GetSoftwareSummary(w http.ResponseWriter, r *http.Request) {
	log.Println("Requested mural info")
	if r.Method != "GET" {
		fmt.Fprintf(w, "This endpoint only supports get")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(mu.Software)

}
