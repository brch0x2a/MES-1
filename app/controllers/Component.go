package controllers

import (
	"encoding/json"
	"log"
	"net/http"

	//"../models"
	"github.com/AdrianRb95/MES/app/models"
)

func Component(w http.ResponseWriter, r *http.Request) {

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

	res := GetComponentList()

	tmpl.ExecuteTemplate(w, "Component", res)

}

func GetComponentList() []models.Component_holder {

	db := dbConn()

	selDB, err := db.Query(`
	SELECT 
    C.id, M.name, C.name, C.description, C.photo
	FROM
		Component C
			INNER JOIN
		Machine M ON M.id = C.id_machine;
	`)
	if err != nil {
		panic(err.Error())
	}

	p := models.Component_holder{}

	res := []models.Component_holder{}

	for selDB.Next() {
		var Id int
		var Name string
		var Description string
		var Photo string
		var Machine string

		err = selDB.Scan(&Id, &Machine, &Name, &Description, &Photo)
		if err != nil {
			panic(err.Error())
		}
		p.Id = Id
		p.Name = Name
		p.Description = Description
		p.Photo = Photo
		p.Machine = Machine

		res = append(res, p)
	}

	defer db.Close()

	return res
}

func GetComponent(w http.ResponseWriter, r *http.Request) {
	db := dbConn()

	nId := r.URL.Query().Get("id")
	selDB, err := db.Query(`
        SELECT
           *
        FROM Component

        WHERE
            id = ?
    `, nId)

	if err != nil {
		panic(err.Error())
	}

	p := models.Component{}

	for selDB.Next() {
		var Id int
		var Name string
		var Description string
		var Photo string
		var Id_machine int

		err = selDB.Scan(&Id, &Name, &Description, &Photo, &Id_machine)
		if err != nil {
			panic(err.Error())
		}
		p.Id = Id
		p.Name = Name
	}
	defer db.Close()
	json.NewEncoder(w).Encode(p)
}

func NewComponent(w http.ResponseWriter, r *http.Request) {
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
	tmpl.ExecuteTemplate(w, "NewComponent", nil)
}

func InsertComponent(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	if r.Method == "POST" {
		Name := r.FormValue("name")
		Description := r.FormValue("description")
		filePath := FileHandlerEngine(r, "photo", "public/images/AM/Components", "component-*.png")
		machine := r.FormValue("machine")

		insForm, err := db.Prepare(`
		INSERT INTO Component(name, description, photo, id_machine) 
			VALUES(?, ?, ?, ?)`)

		if err != nil {
			panic(err.Error())
		}
		insForm.Exec(Name, Description, filePath, machine)
		log.Println("INSERT Component: Name: " + Name)
	}
	defer db.Close()
	http.Redirect(w, r, "/Component", 301)
}

func UpdateComponent(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	if r.Method == "POST" {
		Id := r.FormValue("pid")
		Name := r.FormValue("name")
		Description := r.FormValue("description")
		filePath := FileHandlerEngine(r, "photo", "public/images/AM/Components", "component-*.png")

		machine := r.FormValue("machine")

		log.Println("Component: " + Name)

		if filePath != "public/images/ul_white.png" {
			insForm, err := db.Prepare(`
				UPDATE Component 
				SET name=?, description=?, photo=?, id_machine=?
				WHERE id=?
			`)
			if err != nil {
				panic(err.Error())
			}
			insForm.Exec(Name, Description, filePath, machine, Id)
		} else {
			insForm, err := db.Prepare(`
				UPDATE Component 
				SET name=?, description=?, id_machine =?
				WHERE id=?
			`)
			if err != nil {
				panic(err.Error())
			}
			insForm.Exec(Name, Description, machine, Id)

		}

	}
	defer db.Close()
	http.Redirect(w, r, "/Component", 301)

}

func DeleteComponent(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	p := r.URL.Query().Get("id")
	delForm, err := db.Prepare("DELETE FROM Component WHERE id=?")
	if err != nil {
		panic(err.Error())
	}
	delForm.Exec(p)
	log.Println("DELETE Component")
	defer db.Close()
	http.Redirect(w, r, "/Component", 301)
}

func GetComponentByMachineList(nId string) []models.Component {

	db := dbConn()

	selDB, err := db.Query(`
		SELECT 
			*
		FROM
			Component
		WHERE
			id_machine = ?;
	`, nId)

	if err != nil {
		panic(err.Error())
	}

	p := models.Component{}

	res := []models.Component{}
	for selDB.Next() {
		var Id int
		var Name string
		var Description string
		var Photo string
		var Id_machine int

		err = selDB.Scan(&Id, &Name, &Description, &Photo, &Id_machine)
		if err != nil {
			panic(err.Error())
		}
		p.Id = Id
		p.Name = Name
		p.Description = Description
		p.Photo = Photo
		p.Id_machine = Id_machine

		res = append(res, p)
	}
	defer db.Close()

	return res
}

func GetComponentByMachine(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	res := GetComponentByMachineList(id)

	json.NewEncoder(w).Encode(res)
}
