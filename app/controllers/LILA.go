package controllers

import (
	"encoding/json"
	"net/http"

	"../models"
)

func GetLILAList() []models.LILA {
	db := dbConn()

	selDB, err := db.Query("Select * from LILA_Point")
	if err != nil {
		panic(err.Error())
	}

	p := models.LILA{}

	res := []models.LILA{}

	for selDB.Next() {
		var Id int
		var Name string
		var Color string

		err = selDB.Scan(&Id, &Name, &Color)
		if err != nil {
			panic(err.Error())
		}
		p.Id = Id
		p.Name = Name
		p.Color = Color

		res = append(res, p)
	}

	defer db.Close()

	return res
}

func GetLILACatalog(w http.ResponseWriter, r *http.Request) {
	res := GetLILAList()

	json.NewEncoder(w).Encode(res)
}

func GetLILAobj(id string) models.LILA {
	db := dbConn()

	selDB, err := db.Query(`
        SELECT
           *
        FROM LILA_Point

        WHERE
            id = ?
    `, id)

	if err != nil {
		panic(err.Error())
	}

	p := models.LILA{}

	for selDB.Next() {
		var Id int
		var Name string
		var Color string

		err = selDB.Scan(&Id, &Name, &Color)
		if err != nil {
			panic(err.Error())
		}
		p.Id = Id
		p.Name = Name
		p.Color = Color
	}
	defer db.Close()

	return p
}

func GetLILABy(w http.ResponseWriter, r *http.Request) {

	nId := r.URL.Query().Get("id")

	p := GetLILAobj(nId)

	json.NewEncoder(w).Encode(p)
}
