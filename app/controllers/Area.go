package controllers

import (
	//"bytes"
	"log"
	"net/http"

	//"image/jpeg"
	"encoding/json"

	"../models"
	//	_ "github.com/go-sql-driver/mysql"
)

func Area(w http.ResponseWriter, r *http.Request) {
	db := dbConn()

	selDB, err := db.Query("select A.id, A.Name, F.name from Area A inner join Factory F where A.id_factory = F.id ")
	if err != nil {
		panic(err.Error())
	}

	p := models.Area_holder{}

	res := []models.Area_holder{}
	for selDB.Next() {
		var Id int
		var Name string
		var Factory_name string

		err = selDB.Scan(&Id, &Name, &Factory_name)
		if err != nil {
			panic(err.Error())
		}
		p.Id = Id
		p.Name = Name
		p.Factory_name = Factory_name
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
	tmpl.ExecuteTemplate(w, "Area", res)
	defer db.Close()
}

func GetAreaList() []models.Area {
	db := dbConn()

	selDB, err := db.Query("SELECT id, Name, id_factory FROM Area")

	if err != nil {
		panic(err.Error())
	}

	p := models.Area{}

	res := []models.Area{}
	for selDB.Next() {
		var Id int
		var Name string
		var Id_factory int

		err = selDB.Scan(&Id, &Name, &Id_factory)
		if err != nil {
			panic(err.Error())
		}
		p.Id = Id
		p.Name = Name
		p.Id_factory = Id_factory

		res = append(res, p)
	}
	defer db.Close()

	return res
}

func GetArea(w http.ResponseWriter, r *http.Request) {

	res := GetAreaList()

	json.NewEncoder(w).Encode(res)
}

func GetAreaBy(w http.ResponseWriter, r *http.Request) {
	db := dbConn()

	nId := r.URL.Query().Get("id")
	selDB, err := db.Query(`
        SELECT
            id, Name, id_factory
        FROM
            Area
        WHERE
            id_factory = ?
    `, nId)

	if err != nil {
		panic(err.Error())
	}

	p := models.Area{}

	res := []models.Area{}
	for selDB.Next() {
		var Id int
		var Name string
		var Id_factory int

		err = selDB.Scan(&Id, &Name, &Id_factory)
		if err != nil {
			panic(err.Error())
		}
		p.Id = Id
		p.Name = Name
		p.Id_factory = Id_factory

		res = append(res, p)
	}
	defer db.Close()
	json.NewEncoder(w).Encode(res)
}

func NewArea(w http.ResponseWriter, r *http.Request) {
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

	tmpl.ExecuteTemplate(w, "NewArea", user)
}

func InsertArea(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	if r.Method == "POST" {
		Name := r.FormValue("name")
		Factory := r.FormValue("factory")

		insForm, err := db.Prepare("INSERT INTO Area(Name, id_factory) VALUES(?, ?)")
		if err != nil {
			panic(err.Error())
		}
		insForm.Exec(Name, Factory)
		log.Println("INSERT: Name: " + Name)
	}
	defer db.Close()
	http.Redirect(w, r, "/Area", 301)
}

func EditArea(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	nId := r.URL.Query().Get("id")

	selDB, err := db.Query("SELECT id, Name, id_factory FROM Area WHERE id=?", nId)
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

	tmpl.ExecuteTemplate(w, "EditArea", p)
	defer db.Close()
}

func UpdateArea(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	if r.Method == "POST" {
		Name := r.FormValue("name")
		Factory := r.FormValue("factory")

		Id := r.FormValue("uid")
		insForm, err := db.Prepare("UPDATE Area SET Name=?, id_factory=? WHERE id=?")
		if err != nil {
			panic(err.Error())
		}
		insForm.Exec(Name, Factory, Id)
		log.Println("UPDATE: Name: " + Name)
	}
	defer db.Close()
	http.Redirect(w, r, "/Area", 301)

}

func DeleteArea(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	p := r.URL.Query().Get("id")
	delForm, err := db.Prepare("DELETE FROM Area WHERE id=?")
	if err != nil {
		panic(err.Error())
	}
	delForm.Exec(p)
	log.Println("DELETE")
	defer db.Close()
	http.Redirect(w, r, "/Area", 301)
}
