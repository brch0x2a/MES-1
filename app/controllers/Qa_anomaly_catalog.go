package controllers

import (
	"encoding/json"
	"net/http"

	//"../models"

	"github.com/AdrianRb95/MES/app/models"
)

func GetQa_anomaly_catalog(w http.ResponseWriter, r *http.Request) {
	db := dbConn()

	selDB, err := db.Query("SELECT id, description FROM Qa_anomaly_catalog")

	if err != nil {
		panic(err.Error())
	}

	p := models.Qa_anomaly_catalog{}

	res := []models.Qa_anomaly_catalog{}
	for selDB.Next() {
		var Id int
		var Description string

		err = selDB.Scan(&Id,
			&Description)
		if err != nil {
			panic(err.Error())
		}
		p.Id = Id
		p.Description = Description
		res = append(res, p)
	}
	defer db.Close()
	json.NewEncoder(w).Encode(res)
}
