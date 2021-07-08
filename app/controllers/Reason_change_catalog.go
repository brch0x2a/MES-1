package controllers

import (
	"encoding/json"
	"net/http"

	"../models"
)

func GetReason_change(w http.ResponseWriter, r *http.Request) {
	db := dbConn()

	selDB, err := db.Query("SELECT * FROM Reason_change_catalog")

	if err != nil {
		panic(err.Error())
	}

	p := models.Reason_change_catalog{}

	res := []models.Reason_change_catalog{}

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
