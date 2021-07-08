package controllers

import (
	"encoding/json"
	"log"
	"net/http"

	"../models"
)

func AM_Job(w http.ResponseWriter, r *http.Request) {

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

	res := GetAM_JobList()

	tmpl.ExecuteTemplate(w, "AM_Job", res)

}

func GetAM_JobList() []models.AM_Job_holder {
	db := dbConn()

	selDB, err := db.Query(`
		SELECT 
			AMJ.id,
			AMJ.description,
			M.name,
			C.name,
			C.photo,
			E.name,
			E.photo,
			L.name,
			L.color
		FROM
			AM_Job AMJ
				INNER JOIN
			Component C ON C.id = AMJ.id_component
				INNER JOIN
			Machine M ON M.id = C.id_machine
				INNER JOIN
			EPP E ON E.id = AMJ.id_epp
				INNER JOIN
			LILA_Point L ON L.id = AMJ.id_lila;
	`)
	if err != nil {
		panic(err.Error())
	}

	p := models.AM_Job_holder{}

	res := []models.AM_Job_holder{}

	for selDB.Next() {
		var Id int
		var Description string
		var Machine string
		var Component string
		var ComponentPhoto string
		var EPP string
		var EPPPhoto string
		var LILA string
		var LILAColor string

		err = selDB.Scan(
			&Id, &Description, &Machine,
			&Component, &ComponentPhoto, &EPP,
			&EPPPhoto, &LILA, &LILAColor)

		if err != nil {
			panic(err.Error())
		}
		p.Id = Id
		p.Description = Description
		p.Machine = Machine
		p.Component = Component
		p.ComponentPhoto = ComponentPhoto
		p.EPP = EPP
		p.EPPPhoto = EPPPhoto
		p.LILA = LILA
		p.LILAColor = LILAColor

		res = append(res, p)
	}

	defer db.Close()

	return res
}

func GetAM_JobCatalog(w http.ResponseWriter, r *http.Request) {
	res := GetAM_JobList()

	json.NewEncoder(w).Encode(res)
}

func GetAM_JobBy(w http.ResponseWriter, r *http.Request) {
	db := dbConn()

	nId := r.URL.Query().Get("id")

	selDB, err := db.Query(`
	SELECT 
		AMJ.id,
		AMJ.description,
		M.name,
		C.name,
		C.photo,
		E.name,
		E.photo,
		L.name,
		L.color
	FROM
		AM_Job AMJ
			INNER JOIN
		Component C ON C.id = AMJ.id_component
			INNER JOIN
		Machine M ON M.id = C.id_machine
			INNER JOIN
		EPP E ON E.id = AMJ.id_epp
			INNER JOIN
		LILA_Point L ON L.id = AMJ.id_lila
    WHERE
            AMJ.id = ?
    `, nId)

	if err != nil {
		panic(err.Error())
	}

	p := models.AM_Job_holder{}

	for selDB.Next() {
		var Id int
		var Description string
		var Machine string
		var Component string
		var ComponentPhoto string
		var EPP string
		var EPPPhoto string
		var LILA string
		var LILAColor string

		err = selDB.Scan(&Id, &Description, &Machine,
			&Component, &ComponentPhoto, &EPP,
			&EPPPhoto, &LILA, &LILAColor)
		if err != nil {
			panic(err.Error())
		}
		p.Id = Id
		p.Description = Description
		p.Machine = Machine
		p.Component = Component
		p.ComponentPhoto = ComponentPhoto
		p.EPP = EPP
		p.EPPPhoto = EPPPhoto
		p.LILA = LILA
		p.LILAColor = LILAColor
	}
	defer db.Close()
	json.NewEncoder(w).Encode(p)
}

func NewAM_Job(w http.ResponseWriter, r *http.Request) {
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

	componentList := GetComponentList()
	eppList := GetEPPList()

	data := map[string]interface{}{
		"EPP":       eppList,
		"Component": componentList,
	}

	tmpl.ExecuteTemplate(w, "NewAM_Job", data)
}

func InsertAM_Job(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	if r.Method == "POST" {
		description := r.FormValue("description")
		component := r.FormValue("component")
		epp := r.FormValue("epp")
		lila := r.FormValue("lila")

		// component := r.FormValue("component")

		insForm, err := db.Prepare(`
		INSERT INTO AM_Job(description, id_component, id_epp, id_lila)
		 VALUES(?, ?, ?, ?)`)

		if err != nil {
			panic(err.Error())
		}
		insForm.Exec(description, component, epp, lila)

		log.Println("\nINSERT AM_Job: Name: " + description)
	}
	defer db.Close()
	http.Redirect(w, r, "/AM_Job", 301)
}

func UpdateAM_Job(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	if r.Method == "POST" {
		Id := r.FormValue("pid")
		Name := r.FormValue("pname")
		log.Println("AM_Job: " + Name)
		insForm, err := db.Prepare("UPDATE AM_Job SET name=? WHERE id=?")
		if err != nil {
			panic(err.Error())
		}
		insForm.Exec(Name, Id)
	}
	defer db.Close()
	http.Redirect(w, r, "/AM_Job", 301)

}

func DeleteAM_Job(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	p := r.URL.Query().Get("id")
	delForm, err := db.Prepare("DELETE FROM AM_Job WHERE id=?")
	if err != nil {
		panic(err.Error())
	}
	delForm.Exec(p)
	log.Println("DELETE AM_Job")
	defer db.Close()
	http.Redirect(w, r, "/AM_Job", 301)
}
