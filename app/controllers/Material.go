package controllers

import (

	//"bytes"

	"encoding/json"
	"log"
	"net/http"

	//"image/jpeg"

	//"../models"
	"github.com/AdrianRb95/MES/app/models"

	_ "github.com/go-sql-driver/mysql"
)

func Material(w http.ResponseWriter, r *http.Request) {

	db := dbConn()

	selDB, err := db.Query("Select * from Material")
	if err != nil {
		panic(err.Error())
	}

	p := models.Material{}

	res := []models.Material{}

	for selDB.Next() {
		var Id int
		var Cod_material int
		var Material_name string

		err = selDB.Scan(&Id, &Cod_material, &Material_name)
		if err != nil {
			panic(err.Error())
		}
		p.Id = Id
		p.Cod_material = Cod_material
		p.Material_name = Material_name
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
	tmpl.ExecuteTemplate(w, "Material", res)
	defer db.Close()
}

func NewMaterial(w http.ResponseWriter, r *http.Request) {
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
	tmpl.ExecuteTemplate(w, "NewMaterial", nil)
}

func InsertMaterial(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	if r.Method == "POST" {
		Name := r.FormValue("name")
		Cod := r.FormValue("cod")
		insForm, err := db.Prepare("INSERT INTO Material(cod_material, material_name) VALUES(?,?)")
		if err != nil {
			panic(err.Error())
		}
		insForm.Exec(Cod, Name)
		log.Println("INSERT MAterial: Name: " + Name + " | Cod: " + Cod)
	}
	defer db.Close()
	http.Redirect(w, r, "/material", 301)
}

func EditMaterial(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	nId := r.URL.Query().Get("id")
	selDB, err := db.Query("SELECT * FROM Material WHERE id=?", nId)
	if err != nil {
		panic(err.Error())
	}
	p := models.Material{}

	for selDB.Next() {
		var Id int
		var Cod_material int
		var Material_name string

		err = selDB.Scan(&Id, &Cod_material, &Material_name)
		if err != nil {
			panic(err.Error())
		}
		p.Id = Id
		p.Cod_material = Cod_material
		p.Material_name = Material_name

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
	tmpl.ExecuteTemplate(w, "EditMaterial", p)
	defer db.Close()
}

func UpdateMaterial(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	if r.Method == "POST" {
		Name := r.FormValue("name")
		Cod := r.FormValue("cod")
		Id := r.FormValue("uid")
		insForm, err := db.Prepare("UPDATE Material SET material_name=?, cod_material=? WHERE id=?")
		if err != nil {
			panic(err.Error())
		}
		insForm.Exec(Name, Cod, Id)
		log.Println("UPDATE MAterial: Name: " + Name + " | Cod: " + Cod)
	}
	defer db.Close()
	http.Redirect(w, r, "/material", 301)
}

func DeleteMaterial(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	p := r.URL.Query().Get("id")
	delForm, err := db.Prepare("DELETE FROM Material WHERE id=?")
	if err != nil {
		panic(err.Error())
	}
	delForm.Exec(p)
	log.Println("DELETE Material")
	defer db.Close()
	http.Redirect(w, r, "/material", 301)
}

func GetMaterialBy(w http.ResponseWriter, r *http.Request) {
	db := dbConn()

	nId := r.URL.Query().Get("id")
	selDB, err := db.Query(`
        SELECT
           *
        FROM Material

        WHERE
            id = ?
    `, nId)

	if err != nil {
		panic(err.Error())
	}

	p := models.Material{}

	res := []models.Material{}
	for selDB.Next() {
		var Id int
		var Cod_material int
		var Material_name string

		err = selDB.Scan(&Id, &Cod_material, &Material_name)
		if err != nil {
			panic(err.Error())
		}
		p.Id = Id
		p.Cod_material = Cod_material
		p.Material_name = Material_name

		res = append(res, p)
	}
	defer db.Close()
	json.NewEncoder(w).Encode(res)
}
