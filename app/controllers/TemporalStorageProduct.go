package controllers

import (
	"encoding/json"
	"net/http"

	"../models"
)

func TemporalStorageProduct(w http.ResponseWriter, r *http.Request) {

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

	tmpl.ExecuteTemplate(w, "TemporalStorageProduct", user)
}

func InsertTemporalStorageProduct(w http.ResponseWriter, r *http.Request) {

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

		bache := r.FormValue("bache")
		product := r.FormValue("product")
		tank := r.FormValue("tank")

		comment := r.FormValue("comment")

		smt := `
		insert into TemporalStorageProduct(
			id_line,
			id_product,
			id_autor,
			bache, 
			tank,
			Observation
			)
		values(?, ?, ?, ?, ?, ?)
		`

		insForm, err := db.Prepare(smt)

		if err != nil {
			panic(err.Error())
		}

		insForm.Exec(
			line,
			product,
			autor,
			bache,
			tank,
			comment)

		defer db.Close()
	}

	http.Redirect(w, r, "/temporalStorageProduct", 301)

}

func GetTemporalStorageProductBy(w http.ResponseWriter, r *http.Request) {
	db := dbConn()

	line := r.FormValue("line")
	init := r.FormValue("init")
	end := r.FormValue("end")

	selDB, err := db.Query(`
	SELECT 
		C.id,
		L.name,
		P.name,
		C.date_reg,
		U.profile_picture,
		U.fname,
		U.lname,
		C.bache,
		C.tank,
		C.Observation
	FROM
		TemporalStorageProduct C
			INNER JOIN
		Line L ON L.id = C.id_line
			INNER JOIN
		Product P ON P.id = C.id_product
			INNER JOIN
		User_table U ON U.id = C.id_autor
	WHERE
			C.id_line = ?
		AND C.date_reg BETWEEN ? AND ?
		order by C.date_reg desc;
    `, line, init, end)

	if err != nil {
		panic(err.Error())
	}

	res := []models.TemporalStorageProduct{}
	p := models.TemporalStorageProduct{}

	for selDB.Next() {
		var Id int
		var Line string
		var Product string
		var Date string
		var Profile string
		var Fname string
		var Lname string
		var Bache int
		var Tank int
		var Observation string

		err = selDB.Scan(
			&Id,
			&Line,
			&Product,
			&Date,
			&Profile,
			&Fname,
			&Lname,
			&Bache,
			&Tank,
			&Observation)

		if err != nil {
			panic(err.Error())
		}

		p.Id = Id
		p.Line = Line
		p.Product = Product
		p.Date = Date
		p.Profile = Profile
		p.Fname = Fname
		p.Lname = Lname
		p.Bache = Bache
		p.Tank = Tank
		p.Observation = Observation

		res = append(res, p)
	}
	defer db.Close()
	json.NewEncoder(w).Encode(res)
}

func ConsolidatedTemporalStorageProduct(w http.ResponseWriter, r *http.Request) {

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

	tmpl.ExecuteTemplate(w, "ConsolidatedTemporalStorageProduct", user)
}
