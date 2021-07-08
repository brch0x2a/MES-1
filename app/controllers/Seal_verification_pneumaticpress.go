package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"../models"
)

func FillSealStep(w http.ResponseWriter, r *http.Request) {
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
	defer db.Close()
	tmpl.ExecuteTemplate(w, "FillSealStep", data)
}

func InsertSealControl(w http.ResponseWriter, r *http.Request) {

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
		top := r.FormValue("top")
		lateral_a := r.FormValue("lateral_a")
		lateral_b := r.FormValue("lateral_b")
		bottom_a := r.FormValue("bottom_a")
		bottom_b := r.FormValue("bottom_b")
		yield := r.FormValue("yield")

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
			`insert into Seal_verification_pneumaticpress(
						seal_top, seal_lateral_a, seal_lateral_b,
						seal_bottom_a, seal_bottom_b, 
						yield, id_sub_header) 
			values (?, ?, ?, ?, ?, ?, ?)`)

		if err != nil {
			panic(err.Error())
		}

		insForm.Exec(top, lateral_a, lateral_b, bottom_a, bottom_b, yield,
			SubId)

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

		tmpl.ExecuteTemplate(w, "FillSealStep", data)
	}

}

func GetSealControlV00(w http.ResponseWriter, r *http.Request) {
	db := dbConn()

	nLine := r.URL.Query().Get("line")
	dateInit := r.URL.Query().Get("dinit")
	dateFinal := r.URL.Query().Get("dfinal")

	selDB, err := db.Query(`
	SELECT 
		U.profile_picture,
		S.date_reg,
		SP.turn,
		L.name,
		P.name,
		U.nick_name,
		U.fname,
		U.lname,
		S.seal_top,
		S.seal_lateral_a,
		S.seal_lateral_b,
		S.seal_bottom_a,
		S.seal_bottom_b,
		IF(S.yield = 1, 'Si', 'No')
		FROM
			Seal_verification_pneumaticpress S
				INNER JOIN
			Salsitas_statistical_process SP ON SP.id = S.id_sub_header
				INNER JOIN
			Line L ON L.id = SP.id_line
				INNER JOIN
			Presentation P ON P.id = SP.id_presentation
				INNER JOIN
			User_table U ON U.id = SP.id_operator
		WHERE
				L.id = ?
				AND DATE(S.date_reg) >= ?
				AND DATE(S.date_reg) <  ?`, nLine, dateInit, dateFinal)

	if err != nil {
		panic(err.Error())
	}

	p := models.Seal_verification_pneumaticpress_holder{}
	res := []models.Seal_verification_pneumaticpress_holder{}

	for selDB.Next() {
		var Profile_picture string
		var Date string
		var Turn int
		var Line string
		var Pname string
		var Nname string
		var Fname string
		var Lname string
		var Top float32
		var Lateral_a float32
		var Lateral_b float32
		var Bottom_a float32
		var Bottom_b float32
		var Yield string

		err := selDB.Scan(&Profile_picture, &Date, &Turn, &Line, &Pname, &Nname, &Fname, &Lname,
			&Top, &Lateral_a, &Lateral_b, &Bottom_a, &Bottom_b, &Yield)

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

		p.Top = Top
		p.Lateral_a = Lateral_a
		p.Lateral_b = Lateral_b
		p.Bottom_a = Bottom_a
		p.Bottom_b = Bottom_b
		p.Yield = Yield

		res = append(res, p)
	}
	defer db.Close()
	json.NewEncoder(w).Encode(res)
}

func ConsolidatedSealControl(w http.ResponseWriter, r *http.Request) {

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
	tmpl.ExecuteTemplate(w, "ConsolidatedSealControl", user)
}
