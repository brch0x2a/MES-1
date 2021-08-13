package controllers

import (
	"encoding/json"
	"log"
	"net/http"

	//"../models"

	"github.com/AdrianRb95/MES/app/models"
)

func Job_catalog(w http.ResponseWriter, r *http.Request) {
	db := dbConn()

	selDB, err := db.Query("Select * from Job_catalog")
	if err != nil {
		panic(err.Error())
	}

	p := models.Job_catalog{}

	res := []models.Job_catalog{}

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
	tmpl.ExecuteTemplate(w, "Job_catalog", res)
	defer db.Close()

}

func NewJob_catalog(w http.ResponseWriter, r *http.Request) {
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
	tmpl.ExecuteTemplate(w, "NewJob_catalog", nil)
}

func EditJob_catalog(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	nId := r.URL.Query().Get("id")
	selDB, err := db.Query("SELECT * FROM Job_catalog WHERE id=?", nId)
	if err != nil {
		panic(err.Error())
	}
	p := models.Job_catalog{}

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
	tmpl.ExecuteTemplate(w, "EditJob_catalog", p)
	defer db.Close()
}

func InsertJob_catalog(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	if r.Method == "POST" {
		Name := r.FormValue("name")
		Color := r.FormValue("vcolor")
		insForm, err := db.Prepare("INSERT INTO Job_catalog(name, color) VALUES(?,?)")
		if err != nil {
			panic(err.Error())
		}
		insForm.Exec(Name, Color)
		log.Println("INSERT Job_catalog: Name: " + Name + " | Color: " + Color)
	}
	defer db.Close()
	http.Redirect(w, r, "/job_catalog", 301)
}

func UpdateJob_catalog(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	if r.Method == "POST" {
		Name := r.FormValue("name")
		Cod := r.FormValue("vcolor")
		Id := r.FormValue("uid")
		insForm, err := db.Prepare("UPDATE Job_catalog SET name=?, color=? WHERE id=?")
		if err != nil {
			panic(err.Error())
		}
		insForm.Exec(Name, Cod, Id)
		log.Println("UPDATE Job_catalog: Name: " + Name + " | Color: " + Cod)
	}
	defer db.Close()
	http.Redirect(w, r, "/job_catalog", 301)
}

func DeleteJob_catalog(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	p := r.URL.Query().Get("id")
	delForm, err := db.Prepare("DELETE FROM Job_catalog WHERE id=?")
	if err != nil {
		panic(err.Error())
	}
	delForm.Exec(p)
	log.Println("DELETE Job_catalog")
	defer db.Close()
	http.Redirect(w, r, "/job_catalog", 301)
}

func GetJob_catalog(w http.ResponseWriter, r *http.Request) {

	db := dbConn()

	selDB, err := db.Query("Select * from Job_catalog")
	if err != nil {
		panic(err.Error())
	}

	p := models.Job_catalog{}

	res := []models.Job_catalog{}
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
	json.NewEncoder(w).Encode(res)
}

func GetJob_catalogObjectE(nId string) models.Job_catalog {
	db := dbConn()

	selDB, err := db.Query("SELECT * FROM Job_catalog WHERE id=?", nId)

	if err != nil {
		panic(err.Error())
	}

	p := models.Job_catalog{}

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

func GetJob_catalogE(w http.ResponseWriter, r *http.Request) {

	nId := r.URL.Query().Get("id")

	res := GetJob_catalogObjectE(nId)

	json.NewEncoder(w).Encode(res)
}
