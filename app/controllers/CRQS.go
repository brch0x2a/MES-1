package controllers

import (
	"encoding/json"
	//"fmt"
	//"io/ioutil"
	//"log"
	"net/http"
	"strconv"

	"../models"
)

func CRQSStep(w http.ResponseWriter, r *http.Request) {
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
	`, user.CRQSId)

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
	defer db.Close()
	tmpl.ExecuteTemplate(w, "CRQSStep", data)
}

func InsertCRQS(w http.ResponseWriter, r *http.Request) {
	SubHeader := models.SalsitasControl_subheader_holder{}

	var SubId int
	var filePath string = ""
	session, err := store.Get(r, "cookie-name")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	user := GetUser(session)

	db := dbConn()

	if r.Method == "POST" {
		codText := r.FormValue("codText")
		level := r.FormValue("level")
		subCRQS := r.FormValue("subCRQS")
		observation := r.FormValue("observation")

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

		filePath = FileHandlerEngine(r, "capture", "public/images/crqs/", "crqs-*.png")

		insForm, err := db.Prepare(
			`INSERT INTO CRQS(
				Codification_img,
				Codification_text,
				simbol_level,
				id_CRQS_SubCategory,
				Observation,
				id_sub_header
			)
			VALUES(
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

		insForm.Exec(filePath, codText, level, subCRQS, observation, SubId)

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

		tmpl.ExecuteTemplate(w, "CRQSStep", data)
	}
}

func GetCRQSV00(w http.ResponseWriter, r *http.Request) {
	db := dbConn()

	line := r.URL.Query().Get("line")
	dateInit := r.URL.Query().Get("dinit")
	dateFinal := r.URL.Query().Get("dfinal")

	selDB, err := db.Query(`
	SELECT 
		U.profile_picture,
		C.id,
		C.date_reg,
		SP.turn,
		L.name,
		P.name,
		U.nick_name,
		U.fname,
		U.lname,
		C.Codification_img,
		C.Codification_text,
		IF(C.simbol_level = 1, 'V', if(C.simbol_level=2, 'A', 'R')),
		S.name,
		C.observation
		FROM
		CRQS C
			INNER JOIN
		CRQS_SubCategory S ON C.id_CRQS_SubCategory = S.id
			INNER JOIN
		Salsitas_statistical_process SP ON SP.id = C.id_sub_header
			INNER JOIN
		Line L ON L.id = SP.id_line
			INNER JOIN
		Presentation P ON P.id = SP.id_presentation
			INNER JOIN
		User_table U ON U.id = SP.id_operator
	WHERE
			L.id = ?
			AND DATE(C.date_reg) >= ?
			AND DATE(C.date_reg) <  ?`, line, dateInit, dateFinal)

	if err != nil {
		panic(err.Error())
	}

	p := models.CRQS_holder{}
	res := []models.CRQS_holder{}

	for selDB.Next() {
		var Profile_picture string
		var Id int
		var Date_reg string
		var Turn int
		var Line string
		var Presentation string
		var Nname string
		var Fname string
		var Lname string
		var Codification_img string
		var Codification_text string
		var Simbol_level string
		var SubCategory string
		var Observation string

		err := selDB.Scan(&Profile_picture, &Id, &Date_reg, &Turn, &Line, &Presentation, &Nname, &Fname, &Lname,
			&Codification_img, &Codification_text, &Simbol_level, &SubCategory, &Observation)

		if err != nil {
			panic(err.Error())
		}
		p.Profile_picture = Profile_picture
		p.Id = Id
		p.Date_reg = Date_reg
		p.Turn = Turn
		p.Line = Line
		p.Presentation = Presentation
		p.Nname = Nname
		p.Fname = Fname
		p.Lname = Lname
		p.Codification_img = Codification_img
		p.Codification_text = Codification_text
		p.Simbol_level = Simbol_level
		p.SubCategory = SubCategory
		p.Observation = Observation

		res = append(res, p)
	}
	defer db.Close()
	json.NewEncoder(w).Encode(res)
}

func ConsolidatedCRQS(w http.ResponseWriter, r *http.Request) {

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
	tmpl.ExecuteTemplate(w, "ConsolidatedCRQS", user)
}
