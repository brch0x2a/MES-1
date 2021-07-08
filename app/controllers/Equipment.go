package controllers

import (
	"encoding/json"
	"log"
	"net/http"

	"../models"
)

func GetEquipment(w http.ResponseWriter, r *http.Request) {
	db := dbConn()

	selDB, err := db.Query(`
        SELECT
		   id,
		   name
        FROM Equipment
    `)

	if err != nil {
		panic(err.Error())
	}
	res := []models.Equipment{}
	p := models.Equipment{}

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

func GetEquipmentBy(w http.ResponseWriter, r *http.Request) {
	db := dbConn()

	nId := r.URL.Query().Get("id")
	selDB, err := db.Query(`
        SELECT
           *
        FROM Equipment

        WHERE
            id = ?
    `, nId)

	if err != nil {
		panic(err.Error())
	}

	p := models.Equipment{}

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
	defer db.Close()
	json.NewEncoder(w).Encode(p)
}

func Equipment(w http.ResponseWriter, r *http.Request) {

	db := dbConn()

	selDB, err := db.Query("Select * from Equipment")
	if err != nil {
		panic(err.Error())
	}

	p := models.Equipment{}

	res := []models.Equipment{}

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
	tmpl.ExecuteTemplate(w, "Equipment", res)
	defer db.Close()
}

func NewEquipment(w http.ResponseWriter, r *http.Request) {
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
	tmpl.ExecuteTemplate(w, "NewEquipment", nil)
}

func InsertEquipment(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	if r.Method == "POST" {
		Name := r.FormValue("name")

		insForm, err := db.Prepare("INSERT INTO Equipment(name) VALUES(?)")
		if err != nil {
			panic(err.Error())
		}
		insForm.Exec(Name)
		log.Println("INSERT MAterial: Name: " + Name)
	}
	defer db.Close()
	http.Redirect(w, r, "/equipment", 301)
}

func UpdateEquipment(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	if r.Method == "POST" {
		Id := r.FormValue("pid")
		Name := r.FormValue("pname")
		log.Println("Equipment: " + Name)
		insForm, err := db.Prepare("UPDATE Equipment SET name=? WHERE id=?")
		if err != nil {
			panic(err.Error())
		}
		insForm.Exec(Name, Id)
	}
	defer db.Close()
	http.Redirect(w, r, "/equipment", 301)

}

func DeleteEquipment(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	p := r.URL.Query().Get("id")
	delForm, err := db.Prepare("DELETE FROM Equipment WHERE id=?")
	if err != nil {
		panic(err.Error())
	}
	delForm.Exec(p)
	log.Println("DELETE Equipment")
	defer db.Close()
	http.Redirect(w, r, "/equipment", 301)
}
