package controllers

import (
	"encoding/json"
	"net/http"

	"../models"
)

func PackingControl(w http.ResponseWriter, r *http.Request) {

	session, err := store.Get(r, "cookie-name")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	user := GetUser(session)

	if auth := user.Authenticated; !auth {
		session.AddFlash("You don't have access!")
		err = session.Save(r, w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/forbidden", http.StatusFound)
		return

	}

	/*_______________________________________________________*/

	tmpl.ExecuteTemplate(w, "PackingControl", user)
}

func InsertPackingControl(w http.ResponseWriter, r *http.Request) {

	session, err := store.Get(r, "cookie-name")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	user := GetUser(session)

	if r.Method == "POST" {
		db := dbConn()

		line := user.Line
		presentation := user.Presentation
		autor := user.UserId

		lote := r.FormValue("lote")
		autoclave := r.FormValue("autoclave")
		pallet := r.FormValue("pallet")
		box := r.FormValue("box")

		comment := r.FormValue("comment")

		smt := `
		insert into PackingControl(
			id_line,
			id_presentation,
			id_autor,
			lote, 
			autoclave,
			pallet,
			box,
			Observation
			)
		values(?, ?, ?, ?, ?, ?, ?, ?)
		`

		insForm, err := db.Prepare(smt)

		if err != nil {
			panic(err.Error())
		}
		insForm.Exec(
			line,
			presentation,
			autor,
			lote,
			autoclave,
			pallet,
			box,
			comment)

		defer db.Close()
	}

	http.Redirect(w, r, "/packingControl", 301)

}

func GetPackingControlBy(w http.ResponseWriter, r *http.Request) {
	db := dbConn()

	line := r.FormValue("line")
	init := r.FormValue("init")
	end := r.FormValue("end")

	selDB, err := db.Query(`
	SELECT 
		C.id,
		L.name,
		T.photo,
		P.name,
		C.date_reg,
		U.profile_picture,
		U.fname,
		U.lname,
		C.lote,
		C.autoclave,
		C.pallet,
		C.box,
		C.Observation
	FROM
		PackingControl C
			INNER JOIN
		Line L ON L.id = C.id_line
			INNER JOIN
		Presentation P ON P.id = C.id_presentation
			INNER JOIN
		Product T ON T.id = P.id_product
			INNER JOIN
		User_table U ON U.id = C.id_autor
		WHERE
			C.id_line = ?
        AND C.date_reg BETWEEN ? AND ?;
    `, line, init, end)

	if err != nil {
		panic(err.Error())
	}

	res := []models.PackingControl{}
	p := models.PackingControl{}

	for selDB.Next() {
		var Id int
		var Line string
		var PPhoto string
		var Presentation string
		var Date string
		var Profile string
		var Fname string
		var Lname string
		var Lote int
		var Autoclave int
		var Pallet int
		var Box int
		var Observation string

		err = selDB.Scan(
			&Id,
			&Line,
			&PPhoto,
			&Presentation,
			&Date,
			&Profile,
			&Fname,
			&Lname,
			&Lote,
			&Autoclave,
			&Pallet,
			&Box,
			&Observation)

		if err != nil {
			panic(err.Error())
		}

		p.Id = Id
		p.Line = Line
		p.PPhoto = PPhoto
		p.Presentation = Presentation
		p.Date = Date
		p.Profile = Profile
		p.Fname = Fname
		p.Lname = Lname
		p.Lote = Lote
		p.Autoclave = Autoclave
		p.Pallet = Pallet
		p.Box = Box
		p.Observation = Observation

		res = append(res, p)
	}
	defer db.Close()
	json.NewEncoder(w).Encode(res)
}

func ConsolidatedPackingControl(w http.ResponseWriter, r *http.Request) {

	session, err := store.Get(r, "cookie-name")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	user := GetUser(session)

	if auth := user.Authenticated; !auth {
		session.AddFlash("You don't have access!")
		err = session.Save(r, w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/forbidden", http.StatusFound)
		return

	}

	/*_______________________________________________________*/

	tmpl.ExecuteTemplate(w, "ConsolidatedPackingControl", user)
}
