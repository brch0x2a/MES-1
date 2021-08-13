package controllers

import (
	//"bytes"
	"log"
	"net/http"

	//"image/jpeg"
	"encoding/json"

	//"../models"

	"github.com/AdrianRb95/MES/app/models"

	_ "github.com/go-sql-driver/mysql"
)

func Line(w http.ResponseWriter, r *http.Request) {
	db := dbConn()

	selDB, err := db.Query("select L.id, L.Name, F.name, A.Name from Line L inner join Area A on L.id_area = A.Id INNER JOIN Factory F on A.id_factory = F.id ")
	if err != nil {
		panic(err.Error())
	}

	p := models.Line_holder{}

	res := []models.Line_holder{}
	for selDB.Next() {
		var Id int
		var Name string
		var Factory_name string
		var Area_name string

		err = selDB.Scan(&Id, &Name, &Factory_name, &Area_name)
		if err != nil {
			panic(err.Error())
		}
		p.Id = Id
		p.Name = Name
		p.Factory_name = Factory_name
		p.Area_name = Area_name
		res = append(res, p)
	}
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
	tmpl.ExecuteTemplate(w, "Line", res)
	defer db.Close()
}

func GetLineObjectBy(nId int) models.Line {
	db := dbConn()

	selDB, err := db.Query("SELECT * FROM Line WHERE id=?", nId)

	if err != nil {
		panic(err.Error())
	}

	p := models.Line{}

	for selDB.Next() {
		var Id int
		var Name string
		var Id_area int

		err = selDB.Scan(&Id, &Name, &Id_area)
		if err != nil {
			panic(err.Error())
		}
		p.Id = Id
		p.Name = Name
		p.Id_area = Id_area
	}
	defer db.Close()

	return p
}

func GetLineObjectE(nId string) models.Line {
	db := dbConn()

	selDB, err := db.Query("SELECT * FROM Line WHERE id=?", nId)

	if err != nil {
		panic(err.Error())
	}

	p := models.Line{}

	for selDB.Next() {
		var Id int
		var Name string
		var Id_area int

		err = selDB.Scan(&Id, &Name, &Id_area)
		if err != nil {
			panic(err.Error())
		}
		p.Id = Id
		p.Name = Name
		p.Id_area = Id_area
	}
	defer db.Close()

	return p
}

func GetLineE(w http.ResponseWriter, r *http.Request) {

	nId := r.URL.Query().Get("line")

	res := GetLineObjectE(nId)

	json.NewEncoder(w).Encode(res)
}

func getLinesByArea(nId string) []models.Line {

	db := dbConn()

	selDB, err := db.Query("SELECT * FROM Line WHERE id_area=?  order by name", nId)

	if err != nil {
		panic(err.Error())
	}

	p := models.Line{}

	res := []models.Line{}
	for selDB.Next() {
		var Id int
		var Name string
		var Id_area int

		err = selDB.Scan(&Id, &Name, &Id_area)
		if err != nil {
			panic(err.Error())
		}
		p.Id = Id
		p.Name = Name
		p.Id_area = Id_area

		res = append(res, p)
	}
	defer db.Close()

	return res
}

func GetLineBy(w http.ResponseWriter, r *http.Request) {

	nId := r.URL.Query().Get("id")

	res := getLinesByArea(nId)

	json.NewEncoder(w).Encode(res)
}

func NewLine(w http.ResponseWriter, r *http.Request) {
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

	tmpl.ExecuteTemplate(w, "NewLine", user)
}

func InsertLine(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	if r.Method == "POST" {
		Name := r.FormValue("name")

		Area := r.FormValue("area")

		insForm, err := db.Prepare("INSERT INTO Line(Name, id_area) VALUES(?, ?)")
		if err != nil {
			panic(err.Error())
		}
		insForm.Exec(Name, Area)
		log.Println("INSERT: Name: " + Name)
	}
	defer db.Close()
	http.Redirect(w, r, "/Line", 301)
}

func EditLine(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	nId := r.URL.Query().Get("id")

	selDB, err := db.Query("SELECT * FROM Line WHERE id=?", nId)
	if err != nil {
		panic(err.Error())
	}
	p := models.Area{}

	for selDB.Next() {
		var Id int
		var Id_factory int
		var Name string

		err = selDB.Scan(&Id, &Name, &Id_factory)
		if err != nil {
			panic(err.Error())
		}
		p.Id = Id
		p.Id_factory = Id_factory
		p.Name = Name
	}
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

	tmpl.ExecuteTemplate(w, "EditLine", p)
	defer db.Close()
}

func UpdateLine(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	if r.Method == "POST" {
		Name := r.FormValue("name")
		Area := r.FormValue("area")

		Id := r.FormValue("uid")

		insForm, err := db.Prepare("UPDATE Line SET name=?, id_area=? WHERE id=?")
		if err != nil {
			panic(err.Error())
		}
		insForm.Exec(Name, Area, Id)
		log.Println("UPDATE: Name: " + Name)
	}
	defer db.Close()
	http.Redirect(w, r, "/Line", 301)

}

func DeleteLine(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	p := r.URL.Query().Get("id")
	delForm, err := db.Prepare("DELETE FROM Line WHERE id=?")
	if err != nil {
		panic(err.Error())
	}
	delForm.Exec(p)
	log.Println("DELETE")
	defer db.Close()
	http.Redirect(w, r, "/Line", 301)
}
