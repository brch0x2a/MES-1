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

func Factory(w http.ResponseWriter, r *http.Request) {
	db := dbConn()

	selDB, err := db.Query("Select * from Factory")
	if err != nil {
		panic(err.Error())
	}

	p := models.Factory{}

	res := []models.Factory{}
	for selDB.Next() {
		var Id int
		var Name string

		err = selDB.Scan(&Id, &Name)
		if err != nil {
			panic(err.Error())
		}
		p.Id = Id

		p.Name = Name
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
	tmpl.ExecuteTemplate(w, "Factory", res)
	defer db.Close()
}

func GetFactory(w http.ResponseWriter, r *http.Request) {

	db := dbConn()

	selDB, err := db.Query("Select * from Factory")
	if err != nil {
		panic(err.Error())
	}

	p := models.Factory{}

	res := []models.Factory{}
	for selDB.Next() {
		var Id int
		var Name string

		err = selDB.Scan(&Id, &Name)
		if err != nil {
			panic(err.Error())
		}
		p.Id = Id

		p.Name = Name
		res = append(res, p)
	}
	defer db.Close()
	json.NewEncoder(w).Encode(res)
}

func NewFactory(w http.ResponseWriter, r *http.Request) {
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

	tmpl.ExecuteTemplate(w, "NewFactory", user)
}

func InsertFactory(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	if r.Method == "POST" {
		Name := r.FormValue("name")

		insForm, err := db.Prepare("INSERT INTO Factory(name) VALUES(?)")
		if err != nil {
			panic(err.Error())
		}
		insForm.Exec(Name)
		log.Println("INSERT: Name: " + Name)
	}
	defer db.Close()
	http.Redirect(w, r, "/Factory", 301)
}

func EditFactory(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	nId := r.URL.Query().Get("id")

	selDB, err := db.Query("SELECT * FROM Factory WHERE id=?", nId)
	if err != nil {
		panic(err.Error())
	}
	p := models.Factory{}

	for selDB.Next() {
		var Id int

		var Name string

		err = selDB.Scan(&Id, &Name)
		if err != nil {
			panic(err.Error())
		}
		p.Id = Id

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

	tmpl.ExecuteTemplate(w, "EditFactory", p)
	defer db.Close()
}

func UpdateFactory(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	if r.Method == "POST" {
		Name := r.FormValue("name")

		Id := r.FormValue("uid")
		insForm, err := db.Prepare("UPDATE Factory SET name=? WHERE id=?")
		if err != nil {
			panic(err.Error())
		}
		insForm.Exec(Name, Id)
		log.Println("UPDATE: Name: " + Name)
	}
	defer db.Close()
	http.Redirect(w, r, "/Factory", 301)

}

func DeleteFactory(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	p := r.URL.Query().Get("id")
	delForm, err := db.Prepare("DELETE FROM Factory WHERE id=?")
	if err != nil {
		panic(err.Error())
	}
	delForm.Exec(p)
	log.Println("DELETE")
	defer db.Close()
	http.Redirect(w, r, "/Factory", 301)
}
