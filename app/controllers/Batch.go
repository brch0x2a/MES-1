package controllers

import (
	"encoding/json"
	"net/http"

	//"../models"
	"github.com/AdrianRb95/MES/app/models"
)

func ReportBatch(w http.ResponseWriter, r *http.Request) {

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

	tmpl.ExecuteTemplate(w, "ReportBatch", user)
}

func InsertBatch(w http.ResponseWriter, r *http.Request) {

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

	if r.Method == "POST" {
		db := dbConn()
		batch := r.FormValue("batch")

		smt := `
		insert into Batch(id_operator, id_line, id_presentation, batch_count)
		values (?, ?, ?, ?);
		`

		insForm, err := db.Prepare(smt)

		if err != nil {
			panic(err.Error())
		}
		insForm.Exec(user.UserId, user.Line, user.Presentation, batch)

		defer db.Close()
	}

	w.Header().Set("Server", "MES")
	w.WriteHeader(200)
}

func GetBatchBy(w http.ResponseWriter, r *http.Request) {

	iDate := r.URL.Query().Get("idate")
	fDate := r.URL.Query().Get("fdate")
	area := r.URL.Query().Get("area")

	db := dbConn()

	selDB, err := db.Query(`
	SELECT 
		B.id,
		B.change_date,
		U.profile_picture,
		U.nick_name,
		U.fname,
		U.lname,
		L.name,
		P.name,
		B.batch_count
	FROM
		Batch B
			INNER JOIN
		User_table U ON U.id = B.id_operator
			INNER JOIN
		Line L ON L.id = B.id_line
			INNER JOIN
		Area A ON A.id = L.id_area
			INNER JOIN
		Presentation P ON P.id = B.id_presentation
	WHERE
		A.id = ?
		AND DATE(B.change_date) BETWEEN ? AND ?
	`, area, iDate, fDate)

	if err != nil {
		panic(err.Error())
	}

	p := models.BatchHolder{}

	res := []models.BatchHolder{}
	for selDB.Next() {
		var Id int
		var Profile_picture string

		var User string
		var Fname string
		var Lname string
		var ChangeDate string
		var Line string
		var Presentation string
		var BatchCount string

		err = selDB.Scan(&Id, &ChangeDate, &Profile_picture, &User, &Fname,
			&Lname,
			&Line, &Presentation,
			&BatchCount)

		if err != nil {
			panic(err.Error())
		}

		p.Id = Id
		p.Profile_picture = Profile_picture
		p.User = User
		p.Fname = Fname
		p.Lname = Lname
		p.ChangeDate = ChangeDate
		p.Line = Line
		p.Presentation = Presentation
		p.BatchCount = BatchCount

		res = append(res, p)

	}

	defer db.Close()

	json.NewEncoder(w).Encode(res)
}

func ConsolidatedBatch(w http.ResponseWriter, r *http.Request) {

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

	tmpl.ExecuteTemplate(w, "ConsolidatedBatch", user)
}
