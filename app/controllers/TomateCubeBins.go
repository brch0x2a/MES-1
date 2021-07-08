package controllers

import (
	"encoding/json"
	"net/http"

	"../models"
)

func TomateCubeBins(w http.ResponseWriter, r *http.Request) {

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

	tmpl.ExecuteTemplate(w, "TomateCubeBins", user)
}

func InsertTomateCubeBins(w http.ResponseWriter, r *http.Request) {

	session, err := store.Get(r, "cookie-name")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	user := GetUser(session)

	if r.Method == "POST" {
		db := dbConn()

		line := user.Line
		autor := user.UserId

		lote := r.FormValue("lote")
		weight := r.FormValue("weight")
		aperture_time := r.FormValue("aperture")

		consumption_time := r.FormValue("consumption")

		sticker := r.FormValue("sticker")

		stickerState := 0

		if sticker == "on" {
			stickerState = 1
		}

		observation := r.FormValue("observation")

		smt := `
		insert into TomateCubeBins(
			id_line,
			id_autor,
			lote, 
			weight,
			aperture_time,
			consumption_time,
			sticker,
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
			autor,
			lote,
			weight,
			aperture_time,
			consumption_time,
			stickerState,
			observation)

		defer db.Close()
	}

	http.Redirect(w, r, "/TomateCubeBins", 301)

}

func GetTomateCubeBinsBy(w http.ResponseWriter, r *http.Request) {
	db := dbConn()

	line := r.FormValue("line")
	init := r.FormValue("init")
	end := r.FormValue("end")

	selDB, err := db.Query(`
	SELECT 
		C.id,
		L.name,
		C.date_reg,
		U.profile_picture,
		U.fname,
		U.lname,
		C.lote,
		C.weight,
		C.aperture_time,
		C.consumption_time,
		IFNULL(TIMEDIFF(C.consumption_time, C.aperture_time),
				'') AS exposure,
		C.sticker,
		C.observation
	FROM
		TomateCubeBins C
			INNER JOIN
		Line L ON L.id = C.id_line
			INNER JOIN
		User_table U ON U.id = C.id_autor
	WHERE
			C.id_line = ?
		AND C.date_reg BETWEEN ? AND ?;
    `, line, init, end)

	if err != nil {
		panic(err.Error())
	}

	res := []models.TomateCubeBins{}
	p := models.TomateCubeBins{}

	for selDB.Next() {
		var Id int
		var Line string
		var Date string
		var Aperture string
		var Consumption string

		var Exposure string

		var Profile string
		var Fname string
		var Lname string
		var Lote int
		var Weight float32

		var Sticker bool

		var Observation string

		err = selDB.Scan(
			&Id,
			&Line,
			&Date,
			&Profile,
			&Fname,
			&Lname,
			&Lote,
			&Weight,
			&Aperture,
			&Consumption,
			&Exposure,
			&Sticker,
			&Observation)

		if err != nil {
			panic(err.Error())
		}

		p.Id = Id
		p.Line = Line
		p.Date = Date
		p.Aperture = Aperture
		p.Consumption = Consumption
		p.Exposure = Exposure
		p.Profile = Profile
		p.Fname = Fname
		p.Lname = Lname
		p.Lote = Lote
		p.Weight = Weight
		p.Sticker = Sticker
		p.Observation = Observation

		res = append(res, p)
	}
	defer db.Close()
	json.NewEncoder(w).Encode(res)
}

func ConsolidatedTomateCubeBins(w http.ResponseWriter, r *http.Request) {

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

	tmpl.ExecuteTemplate(w, "ConsolidatedTomateCubeBins", user)
}
