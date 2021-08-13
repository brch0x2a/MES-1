package controllers

/*
Agrupa lo relacionando a boletas verdes
*/

import (
	"database/sql"
	"encoding/json"
	"net/http"

	//"../models"

	"github.com/AdrianRb95/MES/app/models"
)

/*-------------------------------SHE---------------------------------------------*/
func GreenTag(w http.ResponseWriter, r *http.Request) {

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

	/*_______________________________________________________*/
	db := dbConn()

	selDB, err := db.Query("Select * from Equipment")
	if err != nil {
		panic(err.Error())
	}

	p := models.Equipment{}

	res := []models.Equipment{}

	for selDB.Next() {
		var Id int
		var Name string

		err = selDB.Scan(&Id, &Name)
		if err != nil {
			panic(err.Error())
		}
		p.Id = Id
		p.Name = Name

		res = append(res, p)
	}

	defer db.Close()
	/*_______________________________________________________*/

	/*_______________________________________________________*/

	data := map[string]interface{}{
		"User":   user,
		"Equipo": res,
	}

	tmpl.ExecuteTemplate(w, "GreenTag", data)
}

func GetClass_of_event(w http.ResponseWriter, r *http.Request) {
	db := dbConn()

	selDB, err := db.Query("SELECT id, name FROM Class_of_event")

	if err != nil {
		panic(err.Error())
	}

	p := models.Class_of_event{}

	res := []models.Class_of_event{}
	for selDB.Next() {
		var Id int
		var Name string

		err = selDB.Scan(&Id,
			&Name)
		if err != nil {
			panic(err.Error())
		}
		p.Id = Id
		p.Name = Name
		res = append(res, p)
	}
	defer db.Close()
	json.NewEncoder(w).Encode(res)
}

func GetEvent_cause(w http.ResponseWriter, r *http.Request) {
	db := dbConn()

	selDB, err := db.Query("SELECT id, name FROM Event_cause")

	if err != nil {
		panic(err.Error())
	}

	p := models.Class_of_event{}

	res := []models.Class_of_event{}
	for selDB.Next() {
		var Id int
		var Name string

		err = selDB.Scan(&Id,
			&Name)
		if err != nil {
			panic(err.Error())
		}
		p.Id = Id
		p.Name = Name
		res = append(res, p)
	}
	defer db.Close()
	json.NewEncoder(w).Encode(res)
}

func GetFrequency_catalog(w http.ResponseWriter, r *http.Request) {
	db := dbConn()

	selDB, err := db.Query("SELECT id, name FROM Frequency_catalog")

	if err != nil {
		panic(err.Error())
	}

	p := models.Class_of_event{}

	res := []models.Class_of_event{}
	for selDB.Next() {
		var Id int
		var Name string

		err = selDB.Scan(&Id,
			&Name)
		if err != nil {
			panic(err.Error())
		}
		p.Id = Id
		p.Name = Name
		res = append(res, p)
	}
	defer db.Close()
	json.NewEncoder(w).Encode(res)
}

func GetSeverity_catalog(w http.ResponseWriter, r *http.Request) {
	db := dbConn()

	selDB, err := db.Query("SELECT id, name FROM Severity_catalog")

	if err != nil {
		panic(err.Error())
	}

	p := models.Class_of_event{}

	res := []models.Class_of_event{}
	for selDB.Next() {
		var Id int
		var Name string

		err = selDB.Scan(&Id,
			&Name)
		if err != nil {
			panic(err.Error())
		}
		p.Id = Id
		p.Name = Name
		res = append(res, p)
	}
	defer db.Close()
	json.NewEncoder(w).Encode(res)
}

func GetSHE_standard_catalog(w http.ResponseWriter, r *http.Request) {
	db := dbConn()

	selDB, err := db.Query("SELECT id, name FROM SHE_standard_catalog")

	if err != nil {
		panic(err.Error())
	}

	p := models.Class_of_event{}

	res := []models.Class_of_event{}
	for selDB.Next() {
		var Id int
		var Name string

		err = selDB.Scan(&Id,
			&Name)
		if err != nil {
			panic(err.Error())
		}
		p.Id = Id
		p.Name = Name
		res = append(res, p)
	}
	defer db.Close()
	json.NewEncoder(w).Encode(res)
}

func InsertGreenTag(w http.ResponseWriter, r *http.Request) {

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

	if r.Method == "POST" {
		db := dbConn()
		id_line := user.Line
		operator := user.UserId

		priority := r.FormValue("priority")
		equipment := r.FormValue("equipo")
		class_event := r.FormValue("class_event")
		event_cause := r.FormValue("event_cause")
		anomaly := r.FormValue("anomalia")
		frequency := r.FormValue("frequency")
		//--------------------------------
		investigation := "0"
		inWeb := "0"
		//--------------------------------
		severity := r.FormValue("severity")
		correction_action := r.FormValue("correction_action")
		suggestion := r.FormValue("suggestion")
		lesion_description := r.FormValue("lesion_description")
		damage_description := r.FormValue("damage_description")
		standard := r.FormValue("standard")

		ustate := "1"
		smt := `
		INSERT INTO She_tag_green
		(id_line, id_operator, id_type, id_priority, 
		id_equiment, id_class_of_event, id_event_cause, descriptionAnomly, 
		id_frequency, investigation, id_severity, inWeb,
		corrective_action, suggestion, lesion_description, damage_description, 
		id_she_standard, id_state, close_date) 
	   VALUES (?, ?, ?, ?, 
			   ?, ?, ?, ?,
			   ?, ?, ?, ?,
			   ?, ?, ?, ?,
			   ?, ?, CURRENT_TIMESTAMP)`

		//--------------------------------
		if r.FormValue("inWeb") != "" {
			inWeb = r.FormValue("inWeb")
		}
		if r.FormValue("investigation") != "" {
			investigation = r.FormValue("investigation")
		}

		insForm, err := db.Prepare(smt)

		if err != nil {
			panic(err.Error())
		}
		insForm.Exec(
			id_line, operator, 5, priority,
			equipment, class_event, event_cause, anomaly,
			frequency, investigation, severity, inWeb,
			correction_action, suggestion, lesion_description, damage_description,
			standard, ustate)
		defer db.Close()
	}

	http.Redirect(w, r, "/greenTag", 301)
}

func ConsolidatedSHETags(w http.ResponseWriter, r *http.Request) {

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

	tmpl.ExecuteTemplate(w, "ConsolidatedSHETags", user)
}

/*
	Obtiene los registros de boletas de SHE ordenados por fecha reciente y filtrados por rango
	de fechas
*/
func GetSHETagsV00(w http.ResponseWriter, r *http.Request) {
	/*
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
		}*/
	iDate := r.URL.Query().Get("idate")
	fDate := r.URL.Query().Get("fdate")

	db := dbConn()

	selDB, err := db.Query(`
	SELECT 
		T.id,
		L.name,
		U.nick_name,
		U.fname,
		U.lname,
		T.request_date,
		T.close_date,
		TG.name,
		TG.color,
		P.name,
		E.name,
		C.name,
		Ec.name,
		T.descriptionAnomly,
		F.name,
		IF(T.investigation = 1, 'Si', 'No'),
		Sc.name,
		IF(T.inWeb = 1, 'Si', 'No'),
		T.corrective_action,
		T.suggestion,
		T.lesion_description,
		T.damage_description,
		She.name,
		TS.name,
		TS.color
	FROM
		She_tag_green T
			INNER JOIN
		Line L ON L.id = T.id_line
			INNER JOIN
		TypeOfTag TG ON TG.id = T.id_type
			INNER JOIN
		Priority P ON P.id = T.id_priority
			INNER JOIN
		Equipment E ON E.id = T.id_equiment
			INNER JOIN
		User_table U ON U.id = T.id_operator
			INNER JOIN
		Transaction_state TS ON TS.id = T.id_state
			INNER JOIN
		Class_of_event C ON C.id = T.id_class_of_event
			INNER JOIN
		Event_cause Ec ON Ec.id = T.id_event_cause
			INNER JOIN
		Frequency_catalog F ON F.id = T.id_frequency
			INNER JOIN
		Severity_catalog Sc ON Sc.id = T.id_severity
			INNER JOIN
		she_standard_catalog She ON She.id = T.id_she_standard
	WHERE
		? <= DATE(T.request_date)
			AND DATE(T.request_date) <= ?
	ORDER BY request_date DESC

		`, iDate, fDate)

	if err != nil {
		panic(err.Error())
	}

	p := models.SHE_tag_green_holder{}

	res := []models.SHE_tag_green_holder{}
	for selDB.Next() {
		var Id int
		var Line string
		var User string
		var Fname string
		var Lname string
		var RequestDate string
		var CloseDate sql.NullString
		var Type string
		var Color string
		var Priority string
		var Equipment string

		var Class_of_event string
		var Event_cause string
		var Anomaly string
		var Frequency_catalog string
		var Investigation string
		var Severity_catalog string
		var InWeb string
		var Correction_action string
		var Suggestion string
		var Lesion_description string
		var Damage_description string

		var SHE_standard_catalog string
		var State string
		var SColor string

		err = selDB.Scan(
			&Id, &Line, &User, &Fname, &Lname,
			&RequestDate, &CloseDate, &Type, &Color,
			&Priority, &Equipment, &Class_of_event,
			&Event_cause, &Anomaly, &Frequency_catalog,
			&Investigation, &Severity_catalog, &InWeb,
			&Correction_action, &Suggestion, &Lesion_description,
			&Damage_description, &SHE_standard_catalog,
			&State, &SColor)

		if err != nil {
			panic(err.Error())
		}
		p.Id = Id
		p.Line = Line
		p.User = User
		p.Fname = Fname
		p.Lname = Lname
		p.RequestDate = RequestDate
		p.CloseDate = CloseDate
		p.Type = Type
		p.Color = Color
		p.Priority = Priority
		p.Equipment = Equipment
		p.Class_of_event = Class_of_event
		p.Event_cause = Event_cause
		p.Anomaly = Anomaly
		p.Frequency_catalog = Frequency_catalog
		p.Investigation = Investigation
		p.Severity_catalog = Severity_catalog
		p.InWeb = InWeb
		p.Correction_action = Correction_action
		p.Suggestion = Suggestion
		p.Lesion_description = Lesion_description
		p.Damage_description = Damage_description
		p.SHE_standard_catalog = SHE_standard_catalog
		p.State = State
		p.SColor = SColor

		res = append(res, p)
	}
	defer db.Close()
	json.NewEncoder(w).Encode(res)
}
