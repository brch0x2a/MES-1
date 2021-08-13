package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	//"../models"

	"github.com/AdrianRb95/MES/app/models"

	"github.com/gorilla/websocket"
)

func SelectRegisterTag(w http.ResponseWriter, r *http.Request) {

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

	tmpl.ExecuteTemplate(w, "SelectRegisterTag", user)
}

func BlueTag(w http.ResponseWriter, r *http.Request) {

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

	tmpl.ExecuteTemplate(w, "BlueTag", data)
}

func InsertBlueTag(w http.ResponseWriter, r *http.Request) {

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

		if r.FormValue("line") != "" {
			id_line, _ = strconv.ParseInt(r.FormValue("line"), 10, 64)
		}

		priority := "3"
		equipment := r.FormValue("equipo")
		anomaly := r.FormValue("anomalia")

		qa := "0"
		cost := "0"
		prod := "0"
		mortal := "0"
		deliver := "0"
		safe := "0"
		ustate := "3"

		smt := `
		insert into Tag(
			id_line, id_operator, id_type, id_priority, id_equiment,
			descriptionAnomly, benefit_qa, benefit_cost, benefit_productivity, 
			benefit_moral, benefit_deliver, benefit_safety, is_area_affect, 
			is_tag_before, improvement_user, improvement_description, id_state, close_date)
		values (
			?, ?, ?, ?, ?,
			?, ?, ?, ?,
			?, ?, ?, ?,
			?, ?, ?, ?, current_timestamp)`

		if user.Usertype > 1 {
			ustate = "1"
			smt = `
			insert into Tag(
				id_line, id_operator, id_type, id_priority, id_equiment,
				descriptionAnomly, benefit_qa, benefit_cost, benefit_productivity, 
				benefit_moral, benefit_deliver, benefit_safety, is_area_affect, 
				is_tag_before, improvement_user, improvement_description, id_state, close_date)
			values (
				?, ?, ?, ?, ?,
				?, ?, ?, ?,
				?, ?, ?, ?,
				?, ?, ?, ?, null)`
		}

		if r.FormValue("qa") != "" {
			qa = r.FormValue("qa")
		}

		if r.FormValue("costo") != "" {
			cost = r.FormValue("costo")
		}

		if r.FormValue("productividad") != "" {
			prod = r.FormValue("productividad")
		}

		if r.FormValue("mortal") != "" {
			mortal = r.FormValue("mortal")
		}

		if r.FormValue("entrega") != "" {
			deliver = r.FormValue("entrega")
		}

		if r.FormValue("seguridad") != "" {
			safe = r.FormValue("seguridad")
		}

		affect := r.FormValue("affect")

		before := r.FormValue("before")

		improve := r.FormValue("improve")

		if improve != "" {
			ustate = "3"
		} else {
			ustate = "1"
		}

		fmt.Printf("\nAnomaly: %s\tImprove %s\t State: %s", anomaly, improve, ustate)

		autor := r.FormValue("autor")

		insForm, err := db.Prepare(smt)

		if err != nil {
			panic(err.Error())
		}
		insForm.Exec(id_line, operator, 1, priority, equipment, anomaly,
			qa, cost, prod, mortal, deliver, safe,
			affect, before, autor, improve, ustate)
		defer db.Close()
	}

	http.Redirect(w, r, "/blueTag", 301)
}

func GestionSolicitudBoletas(w http.ResponseWriter, r *http.Request) {

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

	tmpl.ExecuteTemplate(w, "GestionSolicitudBoletas", user)
}

/*
	Obtiene los registros de boletas ordenados por fecha reciente y filtrados por rango
	de fechas
*/
func GetTagsV00(w http.ResponseWriter, r *http.Request) {

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
	iDate := r.URL.Query().Get("idate")
	fDate := r.URL.Query().Get("fdate")

	db := dbConn()

	selDB, err := db.Query(`
	SELECT
		U.profile_picture,
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
		T.descriptionAnomly,
		if(T.benefit_qa=1, "Si", "No"),
		if(T.benefit_cost=1, "Si", "No"),
		if(T.benefit_productivity=1, "Si", "No"),
		if(T.benefit_moral=1, "Si", "No"),
		if(T.benefit_deliver=1, "Si", "No"),
		if(T.benefit_safety=1, "Si", "No"),
		if(T.is_area_affect=1, "Si", "No"),
		if(T.is_tag_before=1, "Si", "No"),
		I.nick_name,
		I.fname,
		I.lname,
		T.improvement_description,
		TS.name,
		TS.color
	FROM
		Tag T
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
		User_table I ON I.id = T.improvement_user
			INNER JOIN
		Transaction_state TS ON TS.id = T.id_state
		WHERE
			? <= DATE(T.request_date)
			AND DATE(T.request_date) <= ?
		ORDER BY request_date DESC

		`, iDate, fDate)

	if err != nil {
		panic(err.Error())
	}

	p := models.Tag_holder{}

	res := []models.Tag_holder{}
	for selDB.Next() {
		var Profile_picture string

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
		var Anomaly string

		var Qa string
		var Cost string
		var Product string
		var Mortal string
		var Deliver string
		var Safety string
		var Affect string
		var Before string

		var AUser string
		var AFname string
		var ALname string
		var Improvement string
		var State string
		var SColor string

		err = selDB.Scan(&Profile_picture,
			&Id, &Line, &User, &Fname, &Lname,
			&RequestDate, &CloseDate, &Type, &Color,
			&Priority, &Equipment, &Anomaly, &Qa, &Cost,
			&Product, &Mortal, &Deliver, &Safety, &Affect,
			&Before, &AUser, &AFname, &ALname, &Improvement, &State, &SColor)

		if err != nil {
			panic(err.Error())
		}
		p.Profile_picture = Profile_picture

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
		p.Anomaly = Anomaly
		p.Qa = Qa
		p.Cost = Cost
		p.Product = Product
		p.Mortal = Mortal
		p.Deliver = Deliver
		p.Safety = Safety
		p.Affect = Affect
		p.Before = Before
		p.AUser = AUser
		p.AFname = AFname
		p.ALname = ALname
		p.Improvement = Improvement
		p.State = State
		p.SColor = SColor

		res = append(res, p)
	}
	defer db.Close()
	json.NewEncoder(w).Encode(res)
}

/*_______________________________________________________________*/

/*
	Obtiene los registros de boletas
*/
func GetTagsV01(w http.ResponseWriter, r *http.Request) {
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
		}
	*/
	db := dbConn()

	selDB, err := db.Query(`
	SELECT 
		U.profile_picture,
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
		T.descriptionAnomly,
		if(T.benefit_qa=1, "Si", "No"),
		if(T.benefit_cost=1, "Si", "No"),
		if(T.benefit_productivity=1, "Si", "No"),
		if(T.benefit_moral=1, "Si", "No"),
		if(T.benefit_deliver=1, "Si", "No"),
		if(T.benefit_safety=1, "Si", "No"),
		if(T.is_area_affect=1, "Si", "No"),
		if(T.is_tag_before=1, "Si", "No"),
		I.nick_name,
		I.fname,
		I.lname,
		T.improvement_description,
		TS.name,
		TS.color
	FROM
		Tag T
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
		User_table I ON I.id = T.improvement_user
			INNER JOIN
		Transaction_state TS ON TS.id = T.id_state
		where DATE(T.request_date) = DATE(CURRENT_TIMESTAMP)
		ORDER BY request_date DESC

		`)

	if err != nil {
		panic(err.Error())
	}

	p := models.Tag_holder{}

	res := []models.Tag_holder{}
	for selDB.Next() {
		var Profile_picture string

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
		var Anomaly string

		var Qa string
		var Cost string
		var Product string
		var Mortal string
		var Deliver string
		var Safety string
		var Affect string
		var Before string

		var AUser string
		var AFname string
		var ALname string
		var Improvement string
		var State string
		var SColor string

		err = selDB.Scan(&Profile_picture,
			&Id, &Line, &User, &Fname, &Lname,
			&RequestDate, &CloseDate, &Type, &Color,
			&Priority, &Equipment, &Anomaly, &Qa, &Cost,
			&Product, &Mortal, &Deliver, &Safety, &Affect,
			&Before, &AUser, &AFname, &ALname, &Improvement, &State, &SColor)

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
		p.Anomaly = Anomaly
		p.Qa = Qa
		p.Cost = Cost
		p.Product = Product
		p.Mortal = Mortal
		p.Deliver = Deliver
		p.Safety = Safety
		p.Affect = Affect
		p.Before = Before
		p.AUser = AUser
		p.AFname = AFname
		p.ALname = ALname
		p.Improvement = Improvement
		p.State = State
		p.SColor = SColor

		res = append(res, p)
	}
	defer db.Close()
	json.NewEncoder(w).Encode(res)
}
func GetTagsV02(w http.ResponseWriter, r *http.Request) {

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

	line := r.URL.Query().Get("line")
	iDate := r.URL.Query().Get("idate")
	fDate := r.URL.Query().Get("fdate")
	typeTag := r.URL.Query().Get("typeTag")

	fmt.Println("[Tag]line: ", line)

	smtLine := `T.id_line = ?`

	if line == "-1" {
		smtLine = `coalesce(-1, T.id_line)  = ?`
	}

	smtType := `AND T.id_type = ?`

	if typeTag == "-1" {
		smtType = `AND coalesce(-1, T.id_type)  = ?`
	}

	db := dbConn()

	selDB, err := db.Query(`
	SELECT
		U.profile_picture,
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
		T.descriptionAnomly,
		if(T.benefit_qa=1, "Si", "No"),
		if(T.benefit_cost=1, "Si", "No"),
		if(T.benefit_productivity=1, "Si", "No"),
		if(T.benefit_moral=1, "Si", "No"),
		if(T.benefit_deliver=1, "Si", "No"),
		if(T.benefit_safety=1, "Si", "No"),
		if(T.is_area_affect=1, "Si", "No"),
		if(T.is_tag_before=1, "Si", "No"),
		I.nick_name,
		I.fname,
		I.lname,
		T.improvement_description,
		TS.name,
		TS.color
	FROM
		Tag T
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
		User_table I ON I.id = T.improvement_user
			INNER JOIN
		Transaction_state TS ON TS.id = T.id_state
		WHERE
		`+smtLine+`
		`+smtType+`
			AND ? <= DATE(T.request_date)
			AND DATE(T.request_date) <= ?
		ORDER BY request_date DESC

		`, line, typeTag, iDate, fDate)

	if err != nil {
		panic(err.Error())
	}

	p := models.Tag_holder{}

	res := []models.Tag_holder{}
	for selDB.Next() {
		var Profile_picture string

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
		var Anomaly string

		var Qa string
		var Cost string
		var Product string
		var Mortal string
		var Deliver string
		var Safety string
		var Affect string
		var Before string

		var AUser string
		var AFname string
		var ALname string
		var Improvement string
		var State string
		var SColor string

		err = selDB.Scan(&Profile_picture,
			&Id, &Line, &User, &Fname, &Lname,
			&RequestDate, &CloseDate, &Type, &Color,
			&Priority, &Equipment, &Anomaly, &Qa, &Cost,
			&Product, &Mortal, &Deliver, &Safety, &Affect,
			&Before, &AUser, &AFname, &ALname, &Improvement, &State, &SColor)

		if err != nil {
			panic(err.Error())
		}
		p.Profile_picture = Profile_picture

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
		p.Anomaly = Anomaly
		p.Qa = Qa
		p.Cost = Cost
		p.Product = Product
		p.Mortal = Mortal
		p.Deliver = Deliver
		p.Safety = Safety
		p.Affect = Affect
		p.Before = Before
		p.AUser = AUser
		p.AFname = AFname
		p.ALname = ALname
		p.Improvement = Improvement
		p.State = State
		p.SColor = SColor

		res = append(res, p)
	}
	defer db.Close()
	json.NewEncoder(w).Encode(res)
}

/*_______________________________________________________________*/

func RedTag(w http.ResponseWriter, r *http.Request) {

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

	tmpl.ExecuteTemplate(w, "RedTag", data)
}

func InsertRedTag(w http.ResponseWriter, r *http.Request) {

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

		if r.FormValue("line") != "" {
			id_line, _ = strconv.ParseInt(r.FormValue("line"), 10, 64)
		}

		operator := user.UserId

		priority := "3"
		equipment := r.FormValue("equipo")
		anomaly := r.FormValue("anomalia")
		qa := "0"
		cost := "0"
		prod := "0"
		mortal := "0"
		deliver := "0"
		safe := "0"

		if r.FormValue("qa") != "" {
			qa = r.FormValue("qa")
		}

		if r.FormValue("costo") != "" {
			cost = r.FormValue("costo")
		}

		if r.FormValue("productividad") != "" {
			prod = r.FormValue("productividad")
		}

		if r.FormValue("mortal") != "" {
			mortal = r.FormValue("mortal")
		}

		if r.FormValue("entrega") != "" {
			deliver = r.FormValue("entrega")
		}

		if r.FormValue("seguridad") != "" {
			safe = r.FormValue("seguridad")
		}

		affect := r.FormValue("affect")

		before := r.FormValue("before")

		improve := r.FormValue("mejora")

		autor := r.FormValue("autor")

		insForm, err := db.Prepare(`
			insert into Tag(
				id_line, id_operator, id_type, id_priority, id_equiment,
				descriptionAnomly, benefit_qa, benefit_cost, benefit_productivity, 
				benefit_moral, benefit_deliver, benefit_safety, is_area_affect, 
				is_tag_before, improvement_user, improvement_description, id_state)
			values (
				?, ?, ?, ?, ?,
				?, ?, ?, ?,
				?, ?, ?, ?,
				?, ?, ?, ?)`)

		if err != nil {
			panic(err.Error())
		}
		insForm.Exec(id_line, operator, 2, priority, equipment, anomaly,
			qa, cost, prod, mortal, deliver, safe,
			affect, before, autor, improve, 1)
		defer db.Close()
	}

	http.Redirect(w, r, "/redTag", 301)
}

/*
	Obtiene una boleta en especifico filtrado por id o numero de transaccion
*/
func GetTagEV00(w http.ResponseWriter, r *http.Request) {

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
	id := r.URL.Query().Get("id")

	db := dbConn()

	selDB, err := db.Query(`
	SELECT 
		U.profile_picture,
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
		T.descriptionAnomly,
		if(T.benefit_qa=1, "Si", "No"),
		if(T.benefit_cost=1, "Si", "No"),
		if(T.benefit_productivity=1, "Si", "No"),
		if(T.benefit_moral=1, "Si", "No"),
		if(T.benefit_deliver=1, "Si", "No"),
		if(T.benefit_safety=1, "Si", "No"),
		if(T.is_area_affect=1, "Si", "No"),
		if(T.is_tag_before=1, "Si", "No"),
		I.nick_name,
		I.fname,
		I.lname,
		T.improvement_description,
		TS.name,
		TS.color
	FROM
		Tag T
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
		User_table I ON I.id = T.improvement_user
			INNER JOIN
		Transaction_state TS ON TS.id = T.id_state
		where T.id = ?

		`, id)

	if err != nil {
		panic(err.Error())
	}

	p := models.Tag_holder{}

	for selDB.Next() {
		var Profile_picture string

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
		var Anomaly string

		var Qa string
		var Cost string
		var Product string
		var Mortal string
		var Deliver string
		var Safety string
		var Affect string
		var Before string

		var AUser string
		var AFname string
		var ALname string
		var Improvement string
		var State string
		var SColor string

		err = selDB.Scan(&Profile_picture,
			&Id, &Line, &User, &Fname, &Lname,
			&RequestDate, &CloseDate, &Type, &Color,
			&Priority, &Equipment, &Anomaly, &Qa, &Cost,
			&Product, &Mortal, &Deliver, &Safety, &Affect,
			&Before, &AUser, &AFname, &ALname, &Improvement, &State, &SColor)

		if err != nil {
			panic(err.Error())
		}
		p.Profile_picture = Profile_picture

		p.Profile_picture = Profile_picture
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
		p.Anomaly = Anomaly
		p.Qa = Qa
		p.Cost = Cost
		p.Product = Product
		p.Mortal = Mortal
		p.Deliver = Deliver
		p.Safety = Safety
		p.Affect = Affect
		p.Before = Before
		p.AUser = AUser
		p.AFname = AFname
		p.ALname = ALname
		p.Improvement = Improvement
		p.State = State
		p.SColor = SColor

	}
	defer db.Close()
	json.NewEncoder(w).Encode(p)
}

func UpdateTag(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	if r.Method == "POST" {
		State := r.FormValue("state")
		Uid := r.FormValue("uid")
		smt := ""

		switch State {
		case "1":
			smt = "UPDATE Tag SET id_state=?, close_date=null WHERE id=?"
		case "2":
			smt = "UPDATE Tag SET id_state=?, close_date=null WHERE id=?"
		case "3":
			smt = "UPDATE Tag SET id_state=?, close_date=current_timestamp WHERE id=?"
		}

		insForm, err := db.Prepare(smt)

		if err != nil {
			panic(err.Error())
		}

		insForm.Exec(State, Uid)
		log.Println("UPDATE Tag: State: " + State + " | Uid: " + Uid)
	}
	defer db.Close()
}

func CloseTag(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	if r.Method == "POST" {
		State := r.FormValue("state")
		Uid := r.FormValue("uid")
		improve := r.FormValue("mejora")
		autor := r.FormValue("autor")

		smt := ""

		smt = `
			UPDATE Tag 
			SET id_state=?, close_date=current_timestamp,
			improvement_description=?, improvement_user=?
			WHERE id=?`

		insForm, err := db.Prepare(smt)

		if err != nil {
			panic(err.Error())
		}

		insForm.Exec(State, improve, autor, Uid)
		log.Println("UPDATE Tag: State: " + State + " | Uid: " + Uid)
	}
	defer db.Close()
}

func SetPriotityTag(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	if r.Method == "POST" {

		Uid := r.FormValue("uid")
		priority := r.FormValue("priority")

		smt := `
			UPDATE Tag 
			SET id_priority=?
			WHERE id=?`

		insForm, err := db.Prepare(smt)

		if err != nil {
			panic(err.Error())
		}

		insForm.Exec(priority, Uid)
		log.Println("UPDATE Tag: priority: " + priority + " | Uid: " + Uid)
	}
	defer db.Close()
}

/*------------------------------------------*/
func GestionCurrentTags(w http.ResponseWriter, r *http.Request) {

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

	tmpl.ExecuteTemplate(w, "GestionCurrentTags", user)
}

/*----------------------------------*/
func WsEndpointTag(w http.ResponseWriter, r *http.Request) {
	//fmt.Fprintf(w, "Hello World")
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}

	//clientResult, _ := readerFilter(ws)

	go WriterByTag(ws)
}

func WriterByTag(conn *websocket.Conn) {

	for {
		fmt.Printf("\n\n----------Tag--Monitor----------------\n")

		ticker := time.NewTicker(time.Second)

		if err := conn.WriteMessage(websocket.TextMessage, []byte("")); err != nil {
			fmt.Println(err)
			return
		}

		for t := range ticker.C {

			msg, err_msg := readerFilter(conn)

			if err_msg != nil {
				fmt.Printf("Client socket close!")
				conn.Close()
				break
			}

			if msg == "on" {
				fmt.Printf("Updating At: %+v\n", t)

				//data = rand.Intn(100)

				response, err := http.Get("http://localhost:3000/getTagsV01")

				if err != nil {
					fmt.Printf("The HTTP request failed with error %s\n", err)
					break
				}

				data, _ := ioutil.ReadAll(response.Body)
				fmt.Println("data")
				if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
					fmt.Println(err)
					break
				}
			} else {
				fmt.Printf("Client socket close!")
				conn.Close()
				break
			}

		}
		defer conn.Close()
	}
}

/*-----------------------------------------------------*/
func ConsolidatedTags(w http.ResponseWriter, r *http.Request) {

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

	tmpl.ExecuteTemplate(w, "ConsolidatedTags", user)
}

/*______________________-------------------------_________________________*/

func OrangeTag(w http.ResponseWriter, r *http.Request) {

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

	tmpl.ExecuteTemplate(w, "OrangeTag", data)
}

func InsertOrangeTag(w http.ResponseWriter, r *http.Request) {

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
		anomaly := r.FormValue("anomalia")
		id_anomaly := r.FormValue("id_anomaly")
		ustate := "1"
		smt := `
		insert into Qa_tag_orange(
			id_line, id_operator, id_type, id_priority, id_equiment, id_anomaly,
			descriptionAnomly, improvement_description, id_state, close_date)
		values (
				?, ?, ?, ?, ?, ?,
				?, ?, ?, current_timestamp)`

		improve := r.FormValue("mejora")

		insForm, err := db.Prepare(smt)

		if err != nil {
			panic(err.Error())
		}
		insForm.Exec(id_line, operator, 4, priority, equipment, id_anomaly,
			anomaly, improve, ustate)
		defer db.Close()
	}

	http.Redirect(w, r, "/orangeTag", 301)
}

/*----------------------------------------*/
/*-----------------------------------------------------*/
func ConsolidatedQaTags(w http.ResponseWriter, r *http.Request) {

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

	tmpl.ExecuteTemplate(w, "ConsolidatedQaTags", user)
}

/*
	Obtiene los registros de boletas de Qa ordenados por fecha reciente y filtrados por rango
	de fechas
*/
func GetQaTagsV00(w http.ResponseWriter, r *http.Request) {
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
		U.profile_picture,
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
		T.descriptionAnomly,
		T.improvement_description,
		TS.name,
		TS.color
	FROM
		Qa_tag_orange T
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
		WHERE
			? <= DATE(T.request_date)
			AND DATE(T.request_date) <= ?
		ORDER BY request_date DESC

		`, iDate, fDate)

	if err != nil {
		panic(err.Error())
	}

	p := models.QaTag_holder{}

	res := []models.QaTag_holder{}
	for selDB.Next() {
		var Profile_picture string
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
		var Anomaly string
		var Improvement string
		var State string
		var SColor string

		err = selDB.Scan(&Profile_picture,
			&Id, &Line, &User, &Fname, &Lname,
			&RequestDate, &CloseDate, &Type, &Color,
			&Priority, &Equipment, &Anomaly, &Improvement, &State, &SColor)

		if err != nil {
			panic(err.Error())
		}
		p.Profile_picture = Profile_picture
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
		p.Anomaly = Anomaly
		p.Improvement = Improvement
		p.State = State
		p.SColor = SColor

		res = append(res, p)
	}
	defer db.Close()
	json.NewEncoder(w).Encode(res)
}

func CountOpenTagsBy(w http.ResponseWriter, r *http.Request) {

	line := r.URL.Query().Get("line")

	blueCount := MttTagCount(line, "1", "1")
	redCount := MttTagCount(line, "2", "1")
	greenCount := SHETagCount(line, "5", "1")
	orangeCount := QATagCount(line, "4", "1")

	tagCounter := models.TagCounter{}

	tagCounter.Blue = blueCount
	tagCounter.Red = redCount
	tagCounter.Green = greenCount
	tagCounter.Orange = orangeCount

	json.NewEncoder(w).Encode(tagCounter)
}

func MttTagCount(line string, ttype string, state string) int {

	db := dbConn()

	selDB, err := db.Query(`
	SELECT 
		COUNT(id)
	FROM
		Tag
	WHERE
		id_state = ?
		AND id_type = ?
		AND id_line = ?

		`, state, ttype, line)

	if err != nil {
		panic(err.Error())
	}

	var res int

	for selDB.Next() {

		err = selDB.Scan(&res)

		if err != nil {
			panic(err.Error())
		}

	}
	defer db.Close()

	return res
}

func SHETagCount(line string, ttype string, state string) int {

	db := dbConn()

	selDB, err := db.Query(`
	SELECT 
		COUNT(id)
	FROM
		she_tag_green
	WHERE
		id_state =?
		AND id_type = ?
		AND id_line = ?

		`, state, ttype, line)

	if err != nil {
		panic(err.Error())
	}

	var res int

	for selDB.Next() {

		err = selDB.Scan(&res)

		if err != nil {
			panic(err.Error())
		}

	}
	defer db.Close()

	return res
}

func QATagCount(line string, ttype string, state string) int {

	db := dbConn()

	selDB, err := db.Query(`
	SELECT 
		COUNT(id)
	FROM
		qa_tag_orange
	WHERE
		id_state =?
		AND id_type = ?
		AND id_line = ?

		`, state, ttype, line)

	if err != nil {
		panic(err.Error())
	}

	var res int

	for selDB.Next() {

		err = selDB.Scan(&res)

		if err != nil {
			panic(err.Error())
		}

	}
	defer db.Close()

	return res
}
