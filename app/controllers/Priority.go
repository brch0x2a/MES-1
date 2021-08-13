package controllers

import (
	"encoding/json"
	"net/http"

	//"../models"

	"github.com/AdrianRb95/MES/app/models"
)

func GetPriorityBy(w http.ResponseWriter, r *http.Request) {
	db := dbConn()

	nId := r.URL.Query().Get("id")
	selDB, err := db.Query(`
        SELECT
           *
        FROM Priority

        WHERE
            id = ?
    `, nId)

	if err != nil {
		panic(err.Error())
	}

	p := models.Priority{}

	res := []models.Priority{}
	for selDB.Next() {
		var Id int
		var Name string
		var Description string

		err = selDB.Scan(&Id, &Name, &Description)
		if err != nil {
			panic(err.Error())
		}
		p.Id = Id
		p.Name = Name
		p.Description = Description

		res = append(res, p)
	}
	defer db.Close()
	json.NewEncoder(w).Encode(res)
}
