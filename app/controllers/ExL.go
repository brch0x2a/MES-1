package controllers

import (
	//"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	//"image/jpeg"
	// "encoding/json"
	"../models"
	_ "github.com/go-sql-driver/mysql"
)

func DeleteExL(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	p := r.URL.Query().Get("id")
	delForm, err := db.Prepare("DELETE FROM EventXLine WHERE id=?")
	if err != nil {
		panic(err.Error())
	}
	delForm.Exec(p)
	log.Println("DELETE")
	defer db.Close()
	//http.Redirect(w, r, "/consolidado", 301)
}

func EditExL(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	nId := r.URL.Query().Get("id")

	selDB, err := db.Query("SELECT * FROM EventXLine WHERE id=?", nId)
	if err != nil {
		panic(err.Error())
	}
	p := models.EventXLine{}

	for selDB.Next() {
		var Id int
		var Id_event int
		var Id_line int
		var Date_event string
		var Turn int
		var Minutes int
		var Note string

		err = selDB.Scan(&Id, &Id_event, &Id_line,
			&Date_event, &Turn, &Minutes, &Note)

		if err != nil {
			panic(err.Error())
		}
		p.Id = Id
		p.Id_event = Id_event
		p.Id_line = Id_line
		p.Date_event = Date_event
		p.Turn = Turn
		p.Minutes = Minutes
		p.Note = Note

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

	tmpl.ExecuteTemplate(w, "EditExL", p)
	defer db.Close()
}

func UpdateExL(w http.ResponseWriter, r *http.Request) {

	db := dbConn()
	if r.Method == "POST" {
		Id_event := r.FormValue("Event")
		Id_line := r.FormValue("line")
		Date := r.FormValue("vdate")
		Turn := r.FormValue("vturn")
		Minutes := r.FormValue("minutes")
		Note := r.FormValue("note")

		Id := r.FormValue("uid")

		insForm, err := db.Prepare(`
            UPDATE EventXLine SET id_event=?, id_line=?, 
            date_event=?, turn=?, minutes=?, note=? 
            WHERE id=?`)

		if err != nil {
			panic(err.Error())
		}

		insForm.Exec(Id_event, Id_line, Date, Turn,
			Minutes, Note, Id)

	}
	defer db.Close()
	http.Redirect(w, r, "/consolidado", 301)

}

func GetExLby(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	nId := r.URL.Query().Get("id")

	selDB, err := db.Query(`
	SELECT 
		id, id_event, id_line, date_event,
		turn, minutes, note
	FROM EventXLine WHERE id=?`, nId)
	if err != nil {
		panic(err.Error())
	}
	p := models.EventXLine{}

	for selDB.Next() {
		var Id int
		var Id_event int
		var Id_line int
		var Date_event string
		var Turn int
		var Minutes int
		var Note string

		err = selDB.Scan(&Id, &Id_event, &Id_line,
			&Date_event, &Turn, &Minutes, &Note)

		if err != nil {
			panic(err.Error())
		}
		p.Id = Id
		p.Id_event = Id_event
		p.Id_line = Id_line
		p.Date_event = Date_event
		p.Turn = Turn
		p.Minutes = Minutes
		p.Note = Note

	}

	defer db.Close()
	json.NewEncoder(w).Encode(p)
}

func UpdateExLPartial(w http.ResponseWriter, r *http.Request) {

	db := dbConn()
	if r.Method == "POST" {
		uid := r.FormValue("uid")
		uminutes := r.FormValue("uminutes")

		insForm, err := db.Prepare(`
            UPDATE EventXLine SET  minutes=?
            WHERE id=?`)

		if err != nil {
			panic(err.Error())
		}
		fmt.Printf("-uminutes: %s\tuid:  -", uminutes, uid)
		insForm.Exec(uminutes, uid)
	}
	defer db.Close()
	//http.Redirect(w, r, "/consolidado", 301)
}
