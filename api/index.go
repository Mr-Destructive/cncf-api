package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/mr-destructive/cncf-landscape-api/data"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	db := data.GetDb()
	params := r.URL.Query()
	var filter string
	var args []interface{}

	if name := params.Get("name"); name != "" {
		filter += "name LIKE ?"
		args = append(args, "%"+name+"%")
	} else if category := params.Get("category"); category != "" {
		filter += "category LIKE ?"
		args = append(args, "%"+category+"%")
	} else if subcategory := params.Get("subcategory"); subcategory != "" {
		filter += "subcategory LIKE ?"
		args = append(args, "%"+subcategory+"%")
	}

	registryData, err := data.GetRegistry(db, filter, args...)
	if err != nil {
		log.Fatal(err)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(registryData)
	defer db.Close()
}
