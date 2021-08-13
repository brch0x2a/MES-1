package controllers

import (
	//"bytes"
	"log"
	"net/http"

	//"image/jpeg"
	"encoding/json"
	"strconv"

	//"../models"

	"github.com/AdrianRb95/MES/app/models"

	_ "github.com/go-sql-driver/mysql"
)

func GetEvent(w http.ResponseWriter, r *http.Request) {
	db := dbConn()

	nId := r.URL.Query().Get("id")
	selDB, err := db.Query("SELECT * FROM Event WHERE id_branch=?", nId)

	if err != nil {
		panic(err.Error())
	}

	p := models.Event{}

	res := []models.Event{}
	for selDB.Next() {
		var Id int
		var Description string
		var Id_branch int

		err = selDB.Scan(&Id, &Description, &Id_branch)
		if err != nil {
			panic(err.Error())
		}
		p.Id = Id
		p.Description = Description
		p.Id_branch = Id_branch

		res = append(res, p)
	}
	defer db.Close()
	json.NewEncoder(w).Encode(res)
}

func GetEventBy(w http.ResponseWriter, r *http.Request) {
	db := dbConn()

	nId := r.URL.Query().Get("id")
	selDB, err := db.Query(`
        SELECT
            E.id,
            E.description,
            S.color
        FROM Event E
        INNER JOIN Branch B ON
            E.id_branch = B.id
        INNER JOIN Sub_classification S ON
            S.id = B.id_sub_classification
        WHERE
            E.id = ?
    `, nId)

	if err != nil {
		panic(err.Error())
	}

	p := models.Event_mini_holder{}

	res := []models.Event_mini_holder{}
	for selDB.Next() {
		var Id int
		var Description string

		var Color string

		err = selDB.Scan(&Id, &Description, &Color)
		if err != nil {
			panic(err.Error())
		}
		p.Id = Id
		p.Description = Description
		p.Color = Color

		res = append(res, p)
	}
	defer db.Close()
	json.NewEncoder(w).Encode(res)
}

func GetEventFilterV00(w http.ResponseWriter, r *http.Request) {
	db := dbConn()

	nDate := r.URL.Query().Get("date")
	nTurn := r.URL.Query().Get("turn")
	nLine := r.URL.Query().Get("line")

	hoursQ := 6
	minutesQ := 0

	switch nTurn {
	case "1":
		{
			hoursQ = 6
		}
	case "2":
		{
			hoursQ = 14
		}
	case "3":
		{
			hoursQ = 22
		}

	}
	if nTurn == "4" {
		minutesQ = 1800 - 1
	} else {
		minutesQ = (hoursQ+8)*60 - 1
	}

	selDB, err := db.Query(`
	SELECT
		X.id,
		S.description,
		S.color,
		B.description,
		E.description,
		X.date_event,
		U.nick_name,
		U.fname,
		U.lname,
		U.profile_picture,
		X.turn,
		X.minutes,
		X.note
	FROM
		EventXLine X
	INNER JOIN User_table U ON
		X.id_user = U.id
	INNER JOIN Event E ON
		X.id_event = E.id
	INNER JOIN Branch B ON
		E.id_branch = B.id
	INNER JOIN Sub_classification S ON
		B.id_sub_classification = S.id
			WHERE
		X.id_line = ? 
		AND X.date_event  BETWEEN DATE_ADD(DATE(?), INTERVAL `+strconv.Itoa(hoursQ)+` HOUR) 
        AND DATE_ADD(DATE(?), INTERVAL `+strconv.Itoa(minutesQ)+` MINUTE) 
		order by X.date_event
		`, nLine, nDate, nDate)

	if err != nil {
		panic(err.Error())
	}

	p := models.EventXLineV00{}

	res := []models.EventXLineV00{}
	for selDB.Next() {
		var Id int
		var Sub string
		var Color string
		var Branch string
		var Event string
		var Date_event string
		var Nick_name string
		var Fname string
		var Lname string
		var Profile_picture string
		var Turn int
		var Minutes int
		var Note string

		err = selDB.Scan(&Id, &Sub, &Color, &Branch, &Event,
			&Date_event, &Nick_name, &Fname, &Lname,
			&Profile_picture, &Turn, &Minutes, &Note)
		if err != nil {
			panic(err.Error())
		}

		p.Id = Id
		p.Sub = Sub
		p.Color = Color
		p.Branch = Branch
		p.Event = Event
		p.Minutes = Minutes
		p.Note = Note
		p.Date_event = Date_event
		p.Nick_name = Nick_name
		p.Fname = Fname
		p.Lname = Lname
		p.Profile_picture = Profile_picture
		p.Turn, _ = strconv.Atoi(nTurn)

		res = append(res, p)
	}
	defer db.Close()
	json.NewEncoder(w).Encode(res)
}

func GetEventFilterV01(w http.ResponseWriter, r *http.Request) {
	db := dbConn()

	iDate := r.URL.Query().Get("idate")
	fDate := r.URL.Query().Get("fdate")

	selDB, err := db.Query(`
	SELECT
		X.date_event,
		X.turn,
		L.name,
		S.description,
		B.description,
		E.description,
		X.minutes,
		X.note,
		S.color,
		U.nick_name,
		U.fname,
		U.lname,
		U.profile_picture
	FROM
		EventXLine X
	INNER JOIN Event E ON
		X.id_event = E.id
	INNER JOIN User_table U ON
		X.id_user = U.id
	INNER JOIN Branch B ON
		E.id_branch = B.id
	INNER JOIN Sub_classification S ON
		B.id_sub_classification = S.id
	INNER JOIN Line L ON
		X.id_line = L.id
    WHERE
        DATE(X.date_event) >= ?
            AND DATE(X.date_event) <= ?
    `, iDate, fDate)

	if err != nil {
		panic(err.Error())
	}

	p := models.EventXLineV01{}

	res := []models.EventXLineV01{}
	for selDB.Next() {
		var Date_event string
		var Turn int
		var Line string
		var Sub string
		var Branch string
		var Event string
		var Minutes int
		var Note string
		var Color string
		var Nick_name string
		var Fname string
		var Lname string
		var Profile_picture string

		err = selDB.Scan(
			&Date_event, &Turn, &Line,
			&Sub, &Branch, &Event,
			&Minutes, &Note, &Color,
			&Nick_name, &Fname, &Lname, &Profile_picture)

		if err != nil {
			panic(err.Error())
		}

		p.Date_event = Date_event
		p.Turn = Turn
		p.Line = Line
		p.Sub = Sub
		p.Branch = Branch
		p.Event = Event
		p.Minutes = Minutes
		p.Note = Note
		p.Color = Color
		p.Nick_name = Nick_name
		p.Fname = Fname
		p.Lname = Lname
		p.Profile_picture = Profile_picture

		res = append(res, p)
	}
	defer db.Close()
	json.NewEncoder(w).Encode(res)
}

func InsertEvent(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	if r.Method == "POST" {
		Id_branch := r.FormValue("Branch")
		Description := r.FormValue("description")

		insForm, err := db.Prepare("INSERT INTO Event(description, id_branch ) VALUES(?, ?)")
		if err != nil {
			panic(err.Error())
		}
		insForm.Exec(Description, Id_branch)
		log.Println("INSERT: Name: " + Description)
	}
	defer db.Close()
	http.Redirect(w, r, "/Event", 301)
}

func Event(w http.ResponseWriter, r *http.Request) {
	db := dbConn()

	selDB, err := db.Query(`
	SELECT 
		E.id,
		LTC.description,
		S.description,
		B.description,
		E.description,
		S.color
	FROM
		Event E
			INNER JOIN
		Branch B ON E.id_branch = B.id
			INNER JOIN
		Sub_classification S ON S.id = B.id_sub_classification
			INNER JOIN
		LineTimeClassification LTC ON LTC.id = S.id_LTC
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
	tmpl.ExecuteTemplate(w, "Event", res)
	defer db.Close()
}

func NewEvent(w http.ResponseWriter, r *http.Request) {
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

	tmpl.ExecuteTemplate(w, "NewEvent", user)
}

func EditEvent(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	nId := r.URL.Query().Get("id")

	selDB, err := db.Query("SELECT * FROM Event WHERE id=?", nId)
	if err != nil {
		panic(err.Error())
	}
	p := models.Event{}

	for selDB.Next() {
		var Id int
		var Description string
		var Id_branch int

		err = selDB.Scan(&Id, &Description, &Id_branch)
		if err != nil {
			panic(err.Error())
		}
		p.Id = Id
		p.Id_branch = Id_branch
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

	tmpl.ExecuteTemplate(w, "EditEvent", p)
	defer db.Close()
}

func UpdateEvent(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	if r.Method == "POST" {
		Description := r.FormValue("description")
		Branch := r.FormValue("Branch")
		Id := r.FormValue("uid")

		insForm, err := db.Prepare("UPDATE Event SET description=?, id_branch=? WHERE id=?")
		if err != nil {
			panic(err.Error())
		}
		insForm.Exec(Description, Branch, Id)
		log.Println("UPDATE: Name: " + Description)
	}
	defer db.Close()
	http.Redirect(w, r, "/Event", 301)
}

func DeleteEvent(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	p := r.URL.Query().Get("id")
	delForm, err := db.Prepare("DELETE FROM Event WHERE id=?")
	if err != nil {
		panic(err.Error())
	}
	delForm.Exec(p)
	log.Println("DELETE")
	defer db.Close()
	http.Redirect(w, r, "/Event", 301)
}
