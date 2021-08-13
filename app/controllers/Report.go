package controllers

import (
	//"bytes"
	"log"
	//"html/template"
	"fmt"
	"net/http"
	"strconv"

	//"image/jpeg"
	//"encoding/json"
	//"../models"

	"github.com/AdrianRb95/MES/app/models"
	_ "github.com/go-sql-driver/mysql"
)

func ReportEvent(w http.ResponseWriter, r *http.Request) {

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

	db := dbConn()

	stm := `
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
	`
	// smt := `
	// SELECT
	// 	E.id,
	// 	LTC.description,
	// 	S.description,
	// 	B.description,
	// 	E.description,
	// 	S.color
	// FROM
	// 	LossesXLine X
	// 		INNER JOIN
	// 	Event E ON E.id = X.id_event
	// 		INNER JOIN
	// 	Branch B ON E.id_branch = B.id
	// 		INNER JOIN
	// 	Sub_classification S ON S.id = B.id_sub_classification
	// 		INNER JOIN
	// 	LineTimeClassification LTC ON LTC.id = S.id_LTC
	// 		INNER JOIN
	// 	Line L ON X.id_line = L.id
	// WHERE
	// 	X.id_line = ?
	// ORDER BY S.id;
	// `

	// selDB, err := db.Query(stm, user.Line)
	selDB, err := db.Query(stm)

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

	data := map[string]interface{}{
		"User":  user,
		"Event": res,
	}

	tmpl.ExecuteTemplate(w, "ReportEvent", data)
	defer db.Close()
}
func ReportPlanning(w http.ResponseWriter, r *http.Request) {

	session, err := store.Get(r, "cookie-name")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	user := GetUser(session)

	if auth := user.Authenticated; !auth {
		session.AddFlash("Acceso denegado!")
		err = session.Save(r, w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/forbidden", http.StatusFound)
		return
	}

	tmpl.ExecuteTemplate(w, "ReportPlanning", user)
}

func SelectRegisterWCM(w http.ResponseWriter, r *http.Request) {

	session, err := store.Get(r, "cookie-name")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	user := GetUser(session)

	if auth := user.Authenticated; !auth {
		session.AddFlash("Acceso denegado!")
		err = session.Save(r, w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/forbidden", http.StatusFound)
		return
	}

	tmpl.ExecuteTemplate(w, "SelectRegisterWCM", user)
}

func InsertReportEventValidation(w http.ResponseWriter, r *http.Request) {

	session, err := store.Get(r, "cookie-name")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	user := GetUser(session)

	if auth := user.Authenticated; !auth {
		session.AddFlash("Acceso denegado!")
		err = session.Save(r, w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/forbidden", http.StatusFound)
		return
	}

	db := dbConn()
	if r.Method == "POST" {
		Id_event := r.FormValue("Event")
		Id_line := r.FormValue("line")
		Date := r.FormValue("vdate")
		//Hour := r.FormValue("vtime")
		Turn := user.Turn
		//log.Println("H: "+ Hour)

		Minutes := r.FormValue("minutes")
		Note := r.FormValue("note")

		insForm, err := db.Prepare(`
		INSERT INTO
		EventXLine(id_event, id_line, date_event, id_user, minutes, note, turn) 
		VALUES(?, ?, ?, ?, ?, ?, ?)`)
		if err != nil {
			panic(err.Error())
		}
		insForm.Exec(Id_event, Id_line, Date, user.UserId, Minutes, Note, Turn)
		log.Println("Minutes: " + Minutes)
	}
	defer db.Close()

	w.Header().Set("Server", "MES")
	w.WriteHeader(200)
}

func InsertReportEvent(w http.ResponseWriter, r *http.Request) {

	session, err := store.Get(r, "cookie-name")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	user := GetUser(session)

	if auth := user.Authenticated; !auth {
		session.AddFlash("Acceso denegado!")
		err = session.Save(r, w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/forbidden", http.StatusFound)
		return
	}

	db := dbConn()
	if r.Method == "POST" {
		Id_event := r.FormValue("Event")
		Id_line := user.Line
		Date := r.FormValue("vdate")
		//Hour := r.FormValue("vtime")
		Turn := user.Turn
		//log.Println("H: "+ Hour)

		Minutes := r.FormValue("minutes")
		Note := r.FormValue("note")

		insForm, err := db.Prepare(`
		INSERT INTO
		EventXLine(id_event, id_line, date_event, id_user, minutes, note, turn) 
		VALUES(?, ?, ?, ?, ?, ?, ?)`)
		if err != nil {
			panic(err.Error())
		}
		insForm.Exec(Id_event, Id_line, Date, user.UserId, Minutes, Note, Turn)
		log.Println("Minutes: " + Minutes)
	}
	defer db.Close()

	http.Redirect(w, r, "/reportEvent", 301)
}

func InsertReportPlanning(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	if r.Method == "POST" {

		Line := r.FormValue("line")
		Date_planning := r.FormValue("date_planning")
		Turn := r.FormValue("turn")
		Presentation := r.FormValue("presentation")
		Version := r.FormValue("version")
		Nominal_speed := r.FormValue("nominal_speed")
		Planned := r.FormValue("planned")
		Produced := r.FormValue("produced")

		insForm, err := db.Prepare("INSERT INTO Planning(id_line, id_presentation, date_planning, version, turn, nominal_speed, produced, planned) VALUES(?, ?, ?, ?, ?, ?, ?, ?)")
		if err != nil {
			panic(err.Error())
		}
		insForm.Exec(Line, Presentation, Date_planning, Version, Turn, Nominal_speed, Produced, Planned)

	}
	defer db.Close()
	http.Redirect(w, r, "/reportPlanning", 301)
}

func InsertReportPlanningByWeek(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	if r.Method == "POST" {

		p := models.Planning_insert{}
		count := 3 * 7
		res := []models.Planning_insert{}
		ns := ""
		turn := 0

		fmt.Printf("Test\t Nominal Speed %s\n", r.FormValue("vs1"))

		for n := 1; n <= count; n++ {
			ns = strconv.Itoa(n)

			p.Line = r.FormValue("line")
			p.Date_planning = r.FormValue("f" + ns)
			p.Turn = r.FormValue("t" + ns)
			p.Presentation = r.FormValue("pt" + ns)
			p.Version = r.FormValue("v" + ns)
			p.Nominal_speed = r.FormValue("vs" + ns)
			p.Planned = r.FormValue("pl" + ns)
			p.Produced = r.FormValue("pr" + ns)

			res = append(res, p)
		}

		fmt.Printf("Line\t|Date\t\t|Turn\t|Presentation\t|V\t|NS\t|Planned\t|Produced\n")
		for m := 0; m < count; m++ {
			turn = m%3 + 1

			fmt.Printf("%s\t|%s\t|%s\t|%s\t\t|%s\t|%s\t|%s\t|%s\n",
				res[m].Line, res[m].Date_planning, strconv.Itoa(turn), res[m].Presentation,
				res[m].Version, res[m].Nominal_speed, res[m].Planned, res[m].Produced)

			insForm, err := db.Prepare("INSERT INTO Planning(id_line, id_presentation, date_planning, version, turn, nominal_speed, produced, planned) VALUES(?, ?, ?, ?, ?, ?, ?, ?)")
			if err != nil {
				panic(err.Error())
			}
			insForm.Exec(res[m].Line, res[m].Presentation, res[m].Date_planning,
				res[m].Version, strconv.Itoa(turn), res[m].Nominal_speed, res[m].Produced, res[m].Planned)
		}

	}
	defer db.Close()

	http.Redirect(w, r, "/reportPlanningbyWeek", 301)
}

func ReportProduced(w http.ResponseWriter, r *http.Request) {

	session, err := store.Get(r, "cookie-name")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	user := GetUser(session)

	if auth := user.Authenticated; !auth {
		session.AddFlash("Acceso denegado!")
		err = session.Save(r, w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/forbidden", http.StatusFound)
		return
	}

	tmpl.ExecuteTemplate(w, "ReportProduced", user)
}
