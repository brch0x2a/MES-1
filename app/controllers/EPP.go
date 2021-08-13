package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	//"../models"
	"github.com/AdrianRb95/MES/app/models"
)

func EPP(w http.ResponseWriter, r *http.Request) {

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

	res := GetEPPList()

	tmpl.ExecuteTemplate(w, "EPP", res)

}

func GetEPPList() []models.EPP_holder {
	db := dbConn()

	selDB, err := db.Query("Select * from EPP")
	if err != nil {
		panic(err.Error())
	}

	p := models.EPP_holder{}

	res := []models.EPP_holder{}

	for selDB.Next() {
		var Id int
		var Name string
		var Photo string

		err = selDB.Scan(&Id, &Name, &Photo)
		if err != nil {
			panic(err.Error())
		}
		p.Id = Id
		p.Name = Name
		p.Photo = Photo

		res = append(res, p)
	}

	defer db.Close()

	return res
}

func GetEPPCatalog(w http.ResponseWriter, r *http.Request) {
	res := GetEPPList()

	json.NewEncoder(w).Encode(res)
}

func GetEPPobj(id string) models.EPP_holder {
	db := dbConn()

	selDB, err := db.Query(`
        SELECT
           *
        FROM EPP

        WHERE
            id = ?
    `, id)

	if err != nil {
		panic(err.Error())
	}

	p := models.EPP_holder{}

	for selDB.Next() {
		var Id int
		var Name string
		var Photo string

		err = selDB.Scan(&Id, &Name, &Photo)
		if err != nil {
			panic(err.Error())
		}
		p.Id = Id
		p.Name = Name
		p.Photo = Photo
	}
	defer db.Close()

	return p
}

func GetEPPBy(w http.ResponseWriter, r *http.Request) {

	nId := r.URL.Query().Get("id")

	p := GetEPPobj(nId)

	json.NewEncoder(w).Encode(p)
}

func NewEPP(w http.ResponseWriter, r *http.Request) {
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
	tmpl.ExecuteTemplate(w, "NewEPP", nil)
}

func InsertEPP(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	if r.Method == "POST" {
		Name := r.FormValue("name")
		filePath := FileHandlerEngine(r, "photo", "public/images/AM/EPP", "epp-*.png")

		insForm, err := db.Prepare("INSERT INTO EPP(name, photo) VALUES(?, ?)")
		if err != nil {
			panic(err.Error())
		}
		insForm.Exec(Name, filePath)
		log.Println("INSERT EPP: Name: " + Name)
	}
	defer db.Close()
	http.Redirect(w, r, "/EPP", 301)
}

func UpdateEPP(w http.ResponseWriter, r *http.Request) {

	fmt.Println("updateEPP")

	db := dbConn()
	if r.Method == "POST" {
		Id := r.FormValue("pid")
		Name := r.FormValue("name")

		p := GetEPPobj(Id)

		filePath := FileHandlerEngine(r, "photo", "public/images/AM/EPP", "epp-*.png")

		fmt.Printf("\n\n--\tname:%s\tfile:%s\n\n", Name, filePath)

		if filePath == "public/images/ul_white.png" {
			filePath = p.Photo
		}

		fmt.Printf("\n\n--\tname:%s\tfile:%s\n\n", Name, filePath)

		insForm, err := db.Prepare(`
			UPDATE EPP 
			SET name=?,  photo=?
			WHERE id=?
		`)
		if err != nil {
			panic(err.Error())
		}
		insForm.Exec(Name, filePath, Id)

	}
	defer db.Close()
	http.Redirect(w, r, "/EPP", 301)

}

func DeleteEPP(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	p := r.URL.Query().Get("id")
	delForm, err := db.Prepare("DELETE FROM EPP WHERE id=?")
	if err != nil {
		panic(err.Error())
	}
	delForm.Exec(p)
	log.Println("DELETE EPP")
	defer db.Close()
	http.Redirect(w, r, "/EPP", 301)
}
