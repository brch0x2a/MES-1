package controllers

import (
	//"bytes"
	//"fmt"
	//"log"
	"net/http"
	//"strconv"
	//"image/jpeg"
	"encoding/json"

	"../models"
	_ "github.com/go-sql-driver/mysql"
)

func GetHeaderBy(w http.ResponseWriter, r *http.Request) {

	db := dbConn()

	nId := r.URL.Query().Get("id")

	selDB, err := db.Query("Select * from Header where id=?", nId)

	if err != nil {
		panic(err.Error())
	}

	p := models.Header{}

	for selDB.Next() {
		var Id int
		var Name string
		var Cod_doc string
		var Revision_date string
		var Next_revision_date string
		var Revision_no int

		err = selDB.Scan(&Id, &Name, &Cod_doc, &Revision_date, &Next_revision_date, &Revision_no)

		if err != nil {
			panic(err.Error())
		}

		p.Id = Id
		p.Name = Name
		p.Cod_doc = Cod_doc
		p.Revision_date = Revision_date
		p.Next_revision_date = Next_revision_date
		p.Revision_no = Revision_no

	}

	json.NewEncoder(w).Encode(p)

	defer db.Close()
}

func GetHeaderObjBy(nId int) models.Header {

	db := dbConn()

	selDB, err := db.Query("Select * from Header where id=?", nId)

	if err != nil {
		panic(err.Error())
	}

	p := models.Header{}

	for selDB.Next() {
		var Id int
		var Name string
		var Cod_doc string
		var Revision_date string
		var Next_revision_date string
		var Revision_no int

		err = selDB.Scan(&Id, &Name, &Cod_doc, &Revision_date, &Next_revision_date, &Revision_no)

		if err != nil {
			panic(err.Error())
		}

		p.Id = Id
		p.Name = Name
		p.Cod_doc = Cod_doc
		p.Revision_date = Revision_date
		p.Next_revision_date = Next_revision_date
		p.Revision_no = Revision_no

	}

	defer db.Close()
	return p
}
