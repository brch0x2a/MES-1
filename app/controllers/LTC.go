package controllers

import (
	//"bytes"
	"fmt"
	"log"
	"net/http"

	//"image/jpeg"
	"encoding/json"

	//"../models"
	"github.com/AdrianRb95/MES/app/models"
)

func LineTimeClassification(w http.ResponseWriter, r *http.Request) {
	db := dbConn()

	selDB, err := db.Query("Select * from LineTimeClassification")
	if err != nil {
		panic(err.Error())
	}

	p := models.LTC{}

	res := []models.LTC{}
	for selDB.Next() {
		var Id int
		var Description string

		err = selDB.Scan(&Id, &Description)
		if err != nil {
			panic(err.Error())
		}
		p.Id = Id

		p.Description = Description
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
	tmpl.ExecuteTemplate(w, "LineTimeClassification", res)
	defer db.Close()
}

func NewLtc(w http.ResponseWriter, r *http.Request) {
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

	tmpl.ExecuteTemplate(w, "NewLtc", user)
}

func Ltc(w http.ResponseWriter, r *http.Request) {
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
	tmpl.ExecuteTemplate(w, "Ltc", nil)
}

func GetLTC(w http.ResponseWriter, r *http.Request) {
	db := dbConn()

	selDB, err := db.Query("Select * from LineTimeClassification")
	if err != nil {
		panic(err.Error())
	}

	p := models.LTC{}

	res := []models.LTC{}
	for selDB.Next() {
		var Id int
		var Description string

		err = selDB.Scan(&Id, &Description)
		if err != nil {
			panic(err.Error())
		}
		p.Id = Id

		p.Description = Description
		res = append(res, p)
	}
	defer db.Close()
	json.NewEncoder(w).Encode(res)
}

func InsertLTC(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	if r.Method == "POST" {
		Description := r.FormValue("description")

		insForm, err := db.Prepare("INSERT INTO LineTimeClassification(description) VALUES(?)")
		if err != nil {
			panic(err.Error())
		}
		insForm.Exec(Description)
		log.Println("INSERT: Name: " + Description)
	}
	defer db.Close()
	http.Redirect(w, r, "/LineTimeClassification", 301)
}

func Editltc(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	nId := r.URL.Query().Get("id")

	selDB, err := db.Query("SELECT * FROM LineTimeClassification WHERE id=?", nId)
	if err != nil {
		panic(err.Error())
	}
	p := models.LTC{}

	for selDB.Next() {
		var Id int

		var Description string

		err = selDB.Scan(&Id, &Description)
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

	tmpl.ExecuteTemplate(w, "Editltc", p)
	defer db.Close()
}

func Updateltc(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	if r.Method == "POST" {
		Name := r.FormValue("description")

		Id := r.FormValue("uid")
		insForm, err := db.Prepare("UPDATE LineTimeClassification SET description=? WHERE id=?")
		if err != nil {
			panic(err.Error())
		}
		insForm.Exec(Name, Id)
		log.Println("UPDATE: Name: " + Name)
	}
	defer db.Close()
	http.Redirect(w, r, "/LineTimeClassification", 301)

}

func Deleteltc(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	p := r.URL.Query().Get("id")
	delForm, err := db.Prepare("DELETE FROM LineTimeClassification WHERE id=?")
	if err != nil {
		panic(err.Error())
	}
	delForm.Exec(p)
	log.Println("DELETE")
	defer db.Close()
	http.Redirect(w, r, "/LineTimeClassification", 301)
}

/*______________________Validation_________________________________*/
func LTCLine(w http.ResponseWriter, r *http.Request) {
	db := dbConn()

	selDB, err := db.Query(`
        SELECT
            E.id,
            LTC.description,
            S.description,
            B.description,
            E.description,
            S.color
        FROM Event 
            E
        INNER JOIN Branch B ON
            E.id_branch = B.id
        INNER JOIN Sub_classification S ON
            S.id = B.id_sub_classification
        INNER JOIN LineTimeClassification LTC ON
			LTC.id = S.id_LTC
		ORDER BY S.id
    `)
	if err != nil {
		panic(err.Error())
	}

	p := models.Event_holder{}

	res := []models.Event_holder{}
	for selDB.Next() {
		var Id int
		var LTC string //Line time classification
		var Sub string //Sub classification
		var Branch string
		var Description string
		var Color string

		err = selDB.Scan(&Id, &LTC, &Sub,
			&Branch, &Description, &Color)
		if err != nil {
			panic(err.Error())
		}
		p.Id = Id
		p.LTC = LTC
		p.Sub = Sub
		p.Branch = Branch
		p.Description = Description
		p.Color = Color

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

	data := map[string]interface{}{
		"User":  user,
		"Event": res,
	}
	defer db.Close()

	tmpl.ExecuteTemplate(w, "LTCLine", data)
}

func GetLTCLineBy(w http.ResponseWriter, r *http.Request) {

	line := r.URL.Query().Get("line")

	db := dbConn()

	selDB, err := db.Query(`
	SELECT
		X.id,
		L.name,
		LTC.description,
		S.description,
		B.description,
		E.description,
		E.id,
		S.color
	FROM 
		LossesXLine X
	Inner join Event E On E.id = X.id_event
	INNER JOIN Branch B ON
		E.id_branch = B.id
	INNER JOIN Sub_classification S ON
		S.id = B.id_sub_classification
	INNER JOIN LineTimeClassification LTC ON
		LTC.id = S.id_LTC
	inner join Line L on X.id_line = L.id
	Where X.id_line = ?
	ORDER BY S.id;
	`, line)
	if err != nil {
		panic(err.Error())
	}

	p := models.LTCLine{}

	res := []models.LTCLine{}
	for selDB.Next() {
		var Id int
		var Line string
		var LTC string //Line time classification
		var Sub string //Sub classification
		var Branch string
		var Event string
		var Code int
		var Color string

		err = selDB.Scan(
			&Id, &Line, &LTC, &Sub, &Branch, &Event, &Code, &Color)
		if err != nil {
			panic(err.Error())
		}
		p.Id = Id

		p.Line = Line
		p.LTC = LTC
		p.Sub = Sub
		p.Branch = Branch
		p.Event = Event
		p.Code = Code
		p.Color = Color

		res = append(res, p)
	}
	defer db.Close()
	json.NewEncoder(w).Encode(res)
}

func InsertLTCLine(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	if r.Method == "POST" {
		line := r.FormValue("Line")
		event := r.FormValue("Event")

		fmt.Println("Line: %d\tEvent:%d", line, event)

		insForm, err := db.Prepare("INSERT INTO LossesXLine(id_line, id_event) VALUES(?, ?)")
		if err != nil {
			panic(err.Error())
		}
		insForm.Exec(line, event)
	}
	defer db.Close()

	w.Header().Set("Server", "MES")
	w.WriteHeader(200)
}

func DeleteLTCLine(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	p := r.URL.Query().Get("id")
	delForm, err := db.Prepare("DELETE FROM LossesXLine WHERE id=?")
	if err != nil {
		panic(err.Error())
	}
	delForm.Exec(p)
	log.Println("DELETE")
	defer db.Close()
	//http.Redirect(w, r, "/consolidado", 301)
}
