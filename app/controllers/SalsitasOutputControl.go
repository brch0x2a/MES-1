package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"../models"
)

func SalsitasOutputControl(w http.ResponseWriter, r *http.Request) {

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

	tmpl.ExecuteTemplate(w, "SalsitasOutputControl", user)
}

func InsertSalsitasOutputControl(w http.ResponseWriter, r *http.Request) {

	if r.Method == "POST" {
		db := dbConn()

		line := r.FormValue("line")
		presentation := r.FormValue("presentation")

		output_date := r.FormValue("output_date")
		batch_init := r.FormValue("batch_init")
		batch_end := r.FormValue("batch_end")

		minor := GetUserByNickName(r.FormValue("minor"))
		major := GetUserByNickName(r.FormValue("major"))

		observation := r.FormValue("observation")

		fmt.Printf("line:%s\n", line)
		fmt.Printf("presentation:%s\n", presentation)
		fmt.Printf("output_date:%s\n", output_date)
		fmt.Printf("batch_init:%s\n", batch_init)
		fmt.Printf("batch_end:%s\n", batch_end)
		fmt.Printf("minor:%s\n", minor)
		fmt.Printf("major:%s\n", major)

		fmt.Printf("observation:%s\n", observation)

		smt := `
		insert into SalsitasOutputControl(
			id_line,
			id_presentation,
			output_date,
			batch_init,
			batch_end,
			responsable_minor,
			responsable_major,
			observations
			)
		values(
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?
		)
		`

		insForm, err := db.Prepare(smt)

		if err != nil {
			panic(err.Error())
		}
		insForm.Exec(
			line,
			presentation,
			output_date,
			batch_init,
			batch_end,
			minor,
			major,
			observation)

		defer db.Close()
	}

	http.Redirect(w, r, "/salsitasOutputControl", 301)

}

func ConsolidatedSalsitasOutputControl(w http.ResponseWriter, r *http.Request) {

	session, err := store.Get(r, "cookie-name")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	user := GetUser(session)

	if auth := user.Authenticated; !auth {
		session.AddFlash("Acceso denegado!")
		err = session.Save(r, w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/forbidden", http.StatusFound)
		return
	}

	tmpl.ExecuteTemplate(w, "ConsolidatedSalsitasOutputControl", user)
}

func GetSalsitasOutputControlBy(w http.ResponseWriter, r *http.Request) {

	init := r.URL.Query().Get("init")
	end := r.URL.Query().Get("end")
	line := r.URL.Query().Get("line")

	db := dbConn()

	selDB, err := db.Query(`
	SELECT 
		S.id,
		L.name,
		P.name,
		S.output_date,
		S.batch_init,
		S.batch_end,
		Mi.profile_picture,
		Mi.fname,
		Mi.lname,
		Ma.profile_picture,
		Ma.fname,
		Ma.lname,
		S.observations
	FROM
		SalsitasOutputControl S
			INNER JOIN
		Line L ON L.id = S.id_line
			INNER JOIN
		Presentation P ON P.id = S.id_presentation
			INNER JOIN
		User_table Mi ON Mi.id = S.responsable_minor
			INNER JOIN
		User_table Ma ON Ma.id = S.responsable_major
	WHERE
		S.id_line = ?
			AND S.output_date BETWEEN ? AND ?
	`, line, init, end)

	if err != nil {
		panic(err.Error())
	}

	p := models.SalsitasOutputControl{}

	res := []models.SalsitasOutputControl{}

	for selDB.Next() {
		var Id int

		var Line string
		var Presentation string
		var Output_date string
		var Batch_init string
		var Batch_end string

		var MiProfilePicture string
		var MiFname string
		var MiLname string

		var MaProfilePicture string
		var MaFname string
		var MaLname string

		var Observations string

		err = selDB.Scan(
			&Id,
			&Line,
			&Presentation,
			&Output_date,
			&Batch_init,
			&Batch_end,
			&MiProfilePicture,
			&MiFname,
			&MiLname,
			&MaProfilePicture,
			&MaFname,
			&MaLname,
			&Observations)

		if err != nil {
			panic(err.Error())
		}

		p.Id = Id
		p.Line = Line
		p.Presentation = Presentation
		p.Output_date = Output_date
		p.Batch_init = Batch_init
		p.Batch_end = Batch_end
		p.MiProfilePicture = MiProfilePicture
		p.MiFname = MiFname
		p.MiLname = MiLname
		p.MaProfilePicture = MaProfilePicture
		p.MaFname = MaFname
		p.MaLname = MaLname
		p.Observations = Observations

		res = append(res, p)

	}

	defer db.Close()

	json.NewEncoder(w).Encode(res)

}
