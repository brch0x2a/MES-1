package controllers

import (
	"encoding/json"
	"net/http"

	//"../models"
	"github.com/AdrianRb95/MES/app/models"
)

func GetCQRS_Category(w http.ResponseWriter, r *http.Request) {
	db := dbConn()

	selDB, err := db.Query(`
        SELECT
           *
        FROM CRQS_Category
    `)

	if err != nil {
		panic(err.Error())
	}

	p := models.CRQS_Category{}

	res := []models.CRQS_Category{}

	for selDB.Next() {
		var Id int
		var Name string

		err = selDB.Scan(&Id, &Name)
		if err != nil {
			panic(err.Error())
		}
		p.Id = Id
		p.Name = Name

		res = append(res, p)
	}
	defer db.Close()
	json.NewEncoder(w).Encode(res)
}
