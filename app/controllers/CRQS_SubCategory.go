package controllers

import (
	"encoding/json"
	"net/http"

	//"../models"

	"github.com/AdrianRb95/MES/app/models"
)

func GetCRQS_SubCategoryBy(w http.ResponseWriter, r *http.Request) {
	db := dbConn()

	nId := r.URL.Query().Get("id")
	selDB, err := db.Query("SELECT * FROM CRQS_SubCategory WHERE id_category = ?", nId)

	if err != nil {
		panic(err.Error())
	}

	p := models.CRQS_SubCategory{}

	res := []models.CRQS_SubCategory{}

	for selDB.Next() {
		var Id int
		var Name string
		var Category int

		err = selDB.Scan(&Id, &Name, &Category)
		if err != nil {
			panic(err.Error())
		}
		p.Id = Id
		p.Name = Name
		p.Category = Category

		res = append(res, p)
	}
	defer db.Close()
	json.NewEncoder(w).Encode(res)
}
