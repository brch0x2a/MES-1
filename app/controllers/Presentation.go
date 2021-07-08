package controllers

import (
	//"bytes"
	"log"
	"net/http"
	"strconv"

	//"image/jpeg"
	"encoding/json"

	"../models"
	_ "github.com/go-sql-driver/mysql"
)

func Presentations(w http.ResponseWriter, r *http.Request) {
	db := dbConn()

	selDB, err := db.Query(`
        SELECT
        Pr.id,
        P.name,
        Pr.name,
        Pr.weight_unit,
        Pr.weight_value,
        Pr.error_rate,
        Pr.box_amount,
        P.photo
        FROM
            Presentation Pr
        INNER JOIN Product P ON
            P.id = Pr.id_product
    `)
	if err != nil {
		panic(err.Error())
	}

	p := models.Presentation_holder{}

	res := []models.Presentation_holder{}
	for selDB.Next() {
		var Id int
		var Product string
		var Name string
		var Weight_unit string
		var Weight_value float32
		var Error_rate float32
		var Box_amount int
		var Photo string

		err = selDB.Scan(&Id, &Product,
			&Name, &Weight_unit, &Weight_value,
			&Error_rate, &Box_amount, &Photo)

		if err != nil {
			panic(err.Error())
		}

		p.Id = Id
		p.Product = Product
		p.Name = Name
		p.Weight_unit = Weight_unit
		p.Weight_value = Weight_value
		p.Error_rate = Error_rate
		p.Box_amount = Box_amount
		p.Photo = Photo

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
	tmpl.ExecuteTemplate(w, "Presentations", res)
	defer db.Close()
}

func GetObjectPresentationBy(nId int) models.Presentation {

	db := dbConn()

	selDB, err := db.Query(`
	SELECT
		Pr.id,
		P.id,
		Pr.name,
		Pr.weight_unit,
		Pr.weight_value,
		Pr.error_rate,
		Pr.box_amount,
		P.photo
	FROM
		Presentation Pr
	INNER JOIN Product P ON
		P.id = Pr.id_product
	WHERE Pr.id=?`, nId)

	if err != nil {
		panic(err.Error())
	}

	p := models.Presentation{}

	for selDB.Next() {
		var Id int
		var Id_product int
		var Name string
		var Weight_unit string
		var Weight_value float32
		var Error_rate float32
		var Box_amount int
		var Photo string

		err = selDB.Scan(&Id, &Id_product,
			&Name, &Weight_unit, &Weight_value, &Error_rate,
			&Box_amount, &Photo)
		if err != nil {
			panic(err.Error())
		}
		p.Id = Id
		p.Id_product = Id_product
		p.Name = Name
		p.Weight_unit = Weight_unit
		p.Weight_value = Weight_value
		p.Error_rate = Error_rate
		p.Box_amount = Box_amount
		p.Photo = Photo
	}

	defer db.Close()
	return p
}

func GetPresentationBy(w http.ResponseWriter, r *http.Request) {

	nId, _ := strconv.Atoi(r.URL.Query().Get("id"))

	res := GetObjectPresentationBy(nId)

	json.NewEncoder(w).Encode(res)

}

func GetPresentations(w http.ResponseWriter, r *http.Request) {
	db := dbConn()

	nId := r.URL.Query().Get("id")
	selDB, err := db.Query(`
	SELECT
		Pr.id,
		P.id,
		Pr.name,
		Pr.weight_unit,
		Pr.weight_value,
		Pr.error_rate,
		Pr.box_amount,
		P.photo
	FROM
		Presentation Pr
	INNER JOIN Product P ON
		P.id = Pr.id_product
	WHERE P.id=?`, nId)

	if err != nil {
		panic(err.Error())
	}

	p := models.Presentation{}

	res := []models.Presentation{}

	for selDB.Next() {
		var Id int
		var Id_product int
		var Name string
		var Weight_unit string
		var Weight_value float32
		var Error_rate float32
		var Box_amount int
		var Photo string

		err = selDB.Scan(&Id, &Id_product,
			&Name, &Weight_unit,
			&Weight_value, &Error_rate,
			&Box_amount, &Photo)
		if err != nil {
			panic(err.Error())
		}
		p.Id = Id
		p.Id_product = Id_product
		p.Name = Name
		p.Weight_unit = Weight_unit
		p.Weight_value = Weight_value
		p.Error_rate = Error_rate
		p.Box_amount = Box_amount
		p.Photo = Photo

		res = append(res, p)
	}
	defer db.Close()
	json.NewEncoder(w).Encode(res)
}

func NewPresentation(w http.ResponseWriter, r *http.Request) {
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
	tmpl.ExecuteTemplate(w, "NewPresentation", nil)
}

func InsertPresentation(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	if r.Method == "POST" {
		Id_product := r.FormValue("product")
		Name := r.FormValue("name")
		Weight_value := r.FormValue("weight_value")
		Weight_unit := r.FormValue("weight_unit")
		Error_rate := r.FormValue("error_rate")
		Box_amount := r.FormValue("box_amount")

		insForm, err := db.Prepare("INSERT INTO Presentation(id_product, name, weight_unit, weight_value, error_rate, box_amount) VALUES(?,?,?,?,?, ?)")
		if err != nil {
			panic(err.Error())
		}
		insForm.Exec(Id_product, Name, Weight_unit, Weight_value, Error_rate, Box_amount)
		log.Println("INSERT: Name: " + Name + " | Sku: ")
	}
	defer db.Close()
	http.Redirect(w, r, "/Presentations", 301)
}

func EditPresentation(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	nId := r.URL.Query().Get("id")

	selDB, err := db.Query("SELECT * FROM Presentation WHERE id=?", nId)
	if err != nil {
		panic(err.Error())
	}
	p := models.Presentation{}

	for selDB.Next() {
		var Id int
		var Id_product int
		var Name string
		var Weight_unit string
		var Weight_value float32
		var Error_rate float32
		var Box_amount int

		err = selDB.Scan(&Id, &Name, &Weight_unit,
			&Weight_value, &Error_rate, &Box_amount,
			&Id_product)

		if err != nil {
			panic(err.Error())
		}

		p.Id = Id
		p.Name = Name
		p.Weight_unit = Weight_unit
		p.Weight_value = Weight_value
		p.Error_rate = Error_rate
		p.Box_amount = Box_amount
		p.Id_product = Id_product
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

	tmpl.ExecuteTemplate(w, "EditPresentation", p)
	defer db.Close()
}

func UpdatePresentation(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	if r.Method == "POST" {
		Product := r.FormValue("product")
		Name := r.FormValue("name")
		Weight_value := r.FormValue("weight_value")
		Weight_unit := r.FormValue("weight_unit")
		Error_rate := r.FormValue("error_rate")
		Box_amount := r.FormValue("box_amount")

		Id := r.FormValue("uid")

		insForm, err := db.Prepare(`
                UPDATE
                    Presentation
                SET
                    id_product = ?,
                    name = ?,
                    weight_value = ?,
                    weight_unit = ?,
                    error_rate = ?,
                    box_amount = ?
                WHERE
                    id = ?
        `)
		if err != nil {
			panic(err.Error())
		}
		insForm.Exec(Product, Name,
			Weight_value, Weight_unit,
			Error_rate, Box_amount,
			Id)
		log.Println("UPDATE: Name: " + Name)
	}
	defer db.Close()
	http.Redirect(w, r, "/Presentations", 301)

}

func DeletePresentation(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	p := r.URL.Query().Get("id")
	delForm, err := db.Prepare("DELETE FROM Presentation WHERE id=?")
	if err != nil {
		panic(err.Error())
	}
	delForm.Exec(p)
	log.Println("DELETE")
	defer db.Close()
	http.Redirect(w, r, "/Presentations", 301)
}
