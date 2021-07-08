package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"../models"
)

func ProcessTemperatureControlStep(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "cookie-name")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	user := GetUser(session)

	SubHeader := models.SalsitasControl_subheader_holder{}

	db := dbConn()

	selDB_sub, err := db.Query(`
				SELECT
				SC.id,
				SC.date_reg,
				SC.id_header,
				SC.turn,
				L.name,
				P.name,
				P.weight_value,
				P.weight_unit,
				P.error_rate,
				UC.fname,
				UO.fname
			FROM
				Salsitas_statistical_process SC
			INNER JOIN Line L ON
				SC.id_line = L.id
			INNER JOIN Presentation P ON
				P.id = SC.id_presentation
			INNER JOIN User_table UC ON
				UC.id = SC.id_coordinator
			INNER JOIN User_table UO ON
				SC.id_operator = UO.id
			where SC.id = ?
	`, user.WeightSubId)

	if err != nil {
		panic(err.Error())
	}

	for selDB_sub.Next() {
		var Id int
		var Date string
		var Turn int
		var Line string
		var Presentation string
		var Pvalue float32
		var Punit string
		var Perror float32
		var Coordinator string
		var Operator string
		var Header int

		err = selDB_sub.Scan(&Id, &Date, &Header, &Turn,
			&Line, &Presentation, &Pvalue, &Punit, &Perror,
			&Coordinator, &Operator)

		if err != nil {
			panic(err.Error())
		}

		SubHeader.Id = Id
		SubHeader.Date = Date
		SubHeader.Turn = Turn
		SubHeader.Line = Line
		SubHeader.Presentation = Presentation
		SubHeader.Pvalue = Pvalue
		SubHeader.Punit = Punit
		SubHeader.Perror = Perror //rate
		SubHeader.Coordinator = Coordinator
		SubHeader.Operator = Operator
		SubHeader.Header = Header

	}

	session.Values["weigthId"] = SubHeader.Id

	if auth := user.Authenticated; !auth {
		session.AddFlash("You don't have access!")
		http.Redirect(w, r, "/forbidden", http.StatusFound)
		return
	}

	err = session.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"SubHeader": SubHeader,
		"User":      user,
	}

	tmpl.ExecuteTemplate(w, "FillTemperatureControl", data)
}

func InsertTemperatureControl(w http.ResponseWriter, r *http.Request) {
	SubHeader := models.SalsitasControl_subheader_holder{}

	var SubId int
	session, err := store.Get(r, "cookie-name")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	user := GetUser(session)

	db := dbConn()

	if r.Method == "POST" {
		batch := r.FormValue("batch")
		psi := r.FormValue("psi")
		exchange_temperature := r.FormValue("exchange_temperature")
		hopper_temperature := r.FormValue("hopper_temperature")
		fill_temperature := r.FormValue("fill_temperature")
		observation := r.FormValue("observation")

		fmt.Println("Observation ", observation)

		SubId, _ = strconv.Atoi(r.FormValue("sub"))

		selDB_sub, err := db.Query(`
		SELECT
		SC.id,
		SC.date_reg,
		SC.id_header,
		SC.turn,
		L.name,
		P.name,
		P.weight_value,
		P.weight_unit,
		P.error_rate,
		UC.fname,
		UO.fname
			FROM
				Salsitas_statistical_process SC
			INNER JOIN Line L ON
				SC.id_line = L.id
			INNER JOIN Presentation P ON
				P.id = SC.id_presentation
			INNER JOIN User_table UC ON
				UC.id = SC.id_coordinator
			INNER JOIN User_table UO ON
				SC.id_operator = UO.id
			where SC.id = ?
		`, SubId)

		if err != nil {
			panic(err.Error())
		}

		for selDB_sub.Next() {
			var Id int
			var Date string
			var Turn int
			var Line string
			var Presentation string
			var Pvalue float32
			var Punit string
			var Perror float32
			var Coordinator string
			var Operator string
			var Header int

			err = selDB_sub.Scan(&Id, &Date, &Header, &Turn,
				&Line, &Presentation, &Pvalue, &Punit, &Perror,
				&Coordinator, &Operator)

			if err != nil {
				panic(err.Error())
			}

			SubHeader.Id = Id
			SubHeader.Date = Date
			SubHeader.Turn = Turn
			SubHeader.Line = Line
			SubHeader.Presentation = Presentation
			SubHeader.Pvalue = Pvalue
			SubHeader.Punit = Punit
			SubHeader.Perror = Perror //rate
			SubHeader.Coordinator = Coordinator
			SubHeader.Operator = Operator
			SubHeader.Header = Header

		}

		insForm, err := db.Prepare(
			`INSERT INTO Process_temperature_control(
				batch,
				hopper_psi,
				exchange_temperature,
				hopper_temperature,
				fill_temperature,
				id_sub_header,
				observation
			)
			VALUES(
				?,
				?,
				?,
				?,
				?,
				?,
				?
			)`)

		if err != nil {
			panic(err.Error())
		}

		insForm.Exec(batch, psi, exchange_temperature,
			hopper_temperature, fill_temperature, SubId, observation)

		defer db.Close()

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
			"SubHeader": SubHeader,
			"User":      user,
		}

		tmpl.ExecuteTemplate(w, "FillTemperatureControl", data)
	}

}

func ConsolidatedTemperatureControl(w http.ResponseWriter, r *http.Request) {

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
	tmpl.ExecuteTemplate(w, "ConsolidatedTemperatureControl", user)
}

func GetTemperatureControlV00(w http.ResponseWriter, r *http.Request) {
	db := dbConn()

	nLine := r.URL.Query().Get("line")
	dateInit := r.URL.Query().Get("dinit")
	dateFinal := r.URL.Query().Get("dfinal")

	selDB, err := db.Query(`
	SELECT 
		U.profile_picture,
		TC.date_reg,
		SP.turn,
		L.name,
		P.name,
		U.nick_name,
		U.fname,
		U.lname,
		TC.batch,
		TC.hopper_psi,
		TC.exchange_temperature,
		TC.hopper_temperature,
		TC.fill_temperature,
		Pt.psi_bottom,
		Pt.psi_top,
		Pt.interchange_bottom,
		Pt.interchange_top,
		Pt.hopper_bottom,
		Pt.hopper_top,
		Pt.fill_bottom,
		Pt.fill_top,
		TC.observation
	FROM
		Process_temperature_control TC
			INNER JOIN
		Salsitas_statistical_process SP ON SP.id = TC.id_sub_header
			INNER JOIN
		Line L ON L.id = SP.id_line
			INNER JOIN
		Presentation P ON P.id = SP.id_presentation
			INNER JOIN
		User_table U ON U.id = SP.id_operator
			INNER JOIN
		Product Pt ON P.id_product = Pt.id
	WHERE
			L.id = ?
			AND DATE(TC.date_reg) >= ?
			AND DATE(TC.date_reg) < ?`, nLine, dateInit, dateFinal)

	if err != nil {
		panic(err.Error())
	}

	p := models.TemperatureControl_holder{}
	res := []models.TemperatureControl_holder{}

	for selDB.Next() {
		var Profile_picture string
		var Date string
		var Turn int
		var Line string
		var Pname string
		var Nname string
		var Fname string
		var Lname string
		var Batch int
		var Psi float32
		var Exchange_temperature float32
		var Hopper_temperature float32
		var Fill_temperature float32
		var PSI_bottom float32
		var PSI_top float32
		var Interchange_bottom float32
		var Interchange_top float32
		var Hopper_bottom float32
		var Hopper_top float32
		var Fill_bottom float32
		var Fill_top float32
		var Observation string

		err := selDB.Scan(&Profile_picture, &Date, &Turn, &Line, &Pname, &Nname,
			&Fname, &Lname, &Batch, &Psi, &Exchange_temperature,
			&Hopper_temperature, &Fill_temperature,
			&PSI_bottom, &PSI_top, &Interchange_bottom, &Interchange_top,
			&Hopper_bottom, &Hopper_top, &Fill_bottom, &Fill_top, &Observation)

		if err != nil {
			panic(err.Error())
		}

		p.Profile_picture = Profile_picture
		p.Date = Date
		p.Turn = Turn
		p.Line = Line
		p.Pname = Pname
		p.Nname = Nname
		p.Fname = Fname
		p.Lname = Lname
		p.Batch = Batch
		p.Psi = Psi
		p.Exchange_temperature = Exchange_temperature
		p.Hopper_temperature = Hopper_temperature
		p.Fill_temperature = Fill_temperature

		p.PSI_bottom = PSI_bottom
		p.PSI_top = PSI_top
		p.Interchange_bottom = Interchange_bottom
		p.Interchange_top = Interchange_top
		p.Hopper_bottom = Hopper_top
		p.Fill_bottom = Fill_bottom
		p.Fill_top = Fill_top

		p.Observation = Observation

		res = append(res, p)
	}
	defer db.Close()
	json.NewEncoder(w).Encode(res)
}
