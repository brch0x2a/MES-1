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

func GetSub(w http.ResponseWriter, r *http.Request) {
	db := dbConn()

	nId := r.URL.Query().Get("id")
	selDB, err := db.Query("SELECT * FROM Sub_classification WHERE id_LTC=?", nId)

	if err != nil {
		panic(err.Error())
	}

	p := models.Sub_classification{}

	res := []models.Sub_classification{}
	for selDB.Next() {
		var Id int
		var Description string
		var Color string
		var Id_LTC int

		err = selDB.Scan(&Id, &Description, &Color, &Id_LTC)
		if err != nil {
			panic(err.Error())
		}
		p.Id = Id
		p.Description = Description
		p.Color = Color
		p.Id_LTC = Id_LTC

		res = append(res, p)
	}
	defer db.Close()
	json.NewEncoder(w).Encode(res)
}

func InsertSub(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	if r.Method == "POST" {
		Id_LTC := r.FormValue("LTC")
		Description := r.FormValue("description")
		Color := r.FormValue("vcolor")

		insForm, err := db.Prepare("INSERT INTO Sub_classification(description, color, id_LTC) VALUES(?, ?, ?)")
		if err != nil {
			panic(err.Error())
		}
		insForm.Exec(Description, Color, Id_LTC)
		log.Println("INSERT: Name: " + Description + " color: " + Color)
	}
	defer db.Close()
	http.Redirect(w, r, "/SubClassification", 301)
}

func GetSubE(w http.ResponseWriter, r *http.Request) {
	db := dbConn()

	nId := r.URL.Query().Get("id")
	selDB, err := db.Query("SELECT * FROM Sub_classification WHERE id=?", nId)

	if err != nil {
		panic(err.Error())
	}

	p := models.Sub_classification{}

	res := []models.Sub_classification{}
	for selDB.Next() {
		var Id int
		var Description string
		var Color string
		var Id_LTC int

		err = selDB.Scan(&Id, &Description, &Color, &Id_LTC)
		if err != nil {
			panic(err.Error())
		}
		p.Id = Id
		p.Description = Description
		p.Color = Color
		p.Id_LTC = Id_LTC

		res = append(res, p)
	}
	defer db.Close()
	json.NewEncoder(w).Encode(res)
}

func SubClassification(w http.ResponseWriter, r *http.Request) {
	db := dbConn()

	selDB, err := db.Query(`
            SELECT
            S.id,
            S.description,
            S.color,
            LTC.description
        FROM
            Sub_classification S
        INNER JOIN LineTimeClassification LTC ON
            LTC.id = S.id_LTC
        ORDER BY
            S.id_LTC
    `)
	if err != nil {
		panic(err.Error())
	}

	p := models.Sub_classification_holder{}

	res := []models.Sub_classification_holder{}
	for selDB.Next() {
		var Id int
		var Description string
		var Color string
		var LTC string

		err = selDB.Scan(&Id, &Description, &Color, &LTC)
		if err != nil {
			panic(err.Error())
		}
		p.Id = Id
		p.Description = Description
		p.Color = Color
		p.LTC = LTC

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
	tmpl.ExecuteTemplate(w, "SubClassification", res)
	defer db.Close()
}

func NewSubClassification(w http.ResponseWriter, r *http.Request) {
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

	tmpl.ExecuteTemplate(w, "NewSubClassification", user)
}

func EditSubClassification(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	nId := r.URL.Query().Get("id")

	selDB, err := db.Query("SELECT * FROM Sub_classification WHERE id=?", nId)
	if err != nil {
		panic(err.Error())
	}
	p := models.Sub_classification{}

	for selDB.Next() {
		var Id int
		var Description string
		var Color string
		var Id_LTC int

		err = selDB.Scan(&Id, &Description, &Color, &Id_LTC)
		if err != nil {
			panic(err.Error())
		}
		p.Id = Id

		p.Description = Description
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

	tmpl.ExecuteTemplate(w, "EditSubClassification", p)
	defer db.Close()
}

func UpdateSubClassification(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	if r.Method == "POST" {
		Description := r.FormValue("description")
		LTC := r.FormValue("LTC")
		Color := r.FormValue("vcolor")

		Id := r.FormValue("uid")

		insForm, err := db.Prepare("UPDATE Sub_classification SET description=?, color=?, id_LTC=? WHERE id=?")
		if err != nil {
			panic(err.Error())
		}
		insForm.Exec(Description, Color, LTC, Id)
		log.Println("UPDATE: Name: " + Description)
	}
	defer db.Close()
	http.Redirect(w, r, "/SubClassification", 301)

}

func DeleteSubClassification(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	p := r.URL.Query().Get("id")
	delForm, err := db.Prepare("DELETE FROM Sub_classification WHERE id=?")
	if err != nil {
		panic(err.Error())
	}
	delForm.Exec(p)
	log.Println("DELETE")
	defer db.Close()
	http.Redirect(w, r, "/SubClassification", 301)
}

func getAllSubClassification() []models.Sub_classification_holder {
	db := dbConn()

	selDB, err := db.Query(`
        SELECT
            S.id,
            S.description,
            S.color,
            LTC.description
        FROM
            Sub_classification S
        INNER JOIN LineTimeClassification LTC ON
            LTC.id = S.id_LTC
        ORDER BY
            S.id_LTC
    `)
	if err != nil {
		panic(err.Error())
	}

	p := models.Sub_classification_holder{}

	res := []models.Sub_classification_holder{}
	for selDB.Next() {
		var Id int
		var Description string
		var Color string
		var LTC string

		err = selDB.Scan(&Id, &Description, &Color, &LTC)
		if err != nil {
			panic(err.Error())
		}
		p.Id = Id
		p.Description = Description
		p.Color = Color
		p.LTC = LTC

		res = append(res, p)
	}

	defer db.Close()
	return res

}
