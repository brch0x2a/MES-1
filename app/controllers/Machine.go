package controllers

import (
	"encoding/json"
	"log"
	"net/http"

	"../models"
)

func Machine(w http.ResponseWriter, r *http.Request) {

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

	res := GetMachineList()

	tmpl.ExecuteTemplate(w, "Machine", res)

}

func GetMachineList() []models.Machine {
	db := dbConn()

	selDB, err := db.Query("Select * from Machine")
	if err != nil {
		panic(err.Error())
	}

	p := models.Machine{}

	res := []models.Machine{}

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

	return res
}

func GetMachineCatalog(w http.ResponseWriter, r *http.Request) {
	res := GetMachineList()

	json.NewEncoder(w).Encode(res)
}

func GetMachineBy(w http.ResponseWriter, r *http.Request) {
	db := dbConn()

	nId := r.URL.Query().Get("id")
	selDB, err := db.Query(`
        SELECT
           *
        FROM Machine

        WHERE
            id = ?
    `, nId)

	if err != nil {
		panic(err.Error())
	}

	p := models.Machine{}

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

func NewMachine(w http.ResponseWriter, r *http.Request) {
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
	tmpl.ExecuteTemplate(w, "NewMachine", nil)
}

func InsertMachine(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	if r.Method == "POST" {
		Name := r.FormValue("name")

		insForm, err := db.Prepare("INSERT INTO Machine(name) VALUES(?)")
		if err != nil {
			panic(err.Error())
		}
		insForm.Exec(Name)
		log.Println("INSERT Machine: Name: " + Name)
	}
	defer db.Close()
	http.Redirect(w, r, "/Machine", 301)
}

func UpdateMachine(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	if r.Method == "POST" {
		Id := r.FormValue("pid")
		Name := r.FormValue("pname")
		log.Println("Machine: " + Name)
		insForm, err := db.Prepare("UPDATE Machine SET name=? WHERE id=?")
		if err != nil {
			panic(err.Error())
		}
		insForm.Exec(Name, Id)
	}
	defer db.Close()
	http.Redirect(w, r, "/Machine", 301)

}

func DeleteMachine(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	p := r.URL.Query().Get("id")
	delForm, err := db.Prepare("DELETE FROM Machine WHERE id=?")
	if err != nil {
		panic(err.Error())
	}
	delForm.Exec(p)
	log.Println("DELETE Machine")
	defer db.Close()
	http.Redirect(w, r, "/Machine", 301)
}
