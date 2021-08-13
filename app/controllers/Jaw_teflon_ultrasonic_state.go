package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	//"../models"

	"github.com/AdrianRb95/MES/app/models"
)

func Jaw_teflon_ultrasonic_StateStep(w http.ResponseWriter, r *http.Request) {
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
	tmpl.ExecuteTemplate(w, "FillJawTeflonUltrasonicState", data)
}

func InsertJawControl(w http.ResponseWriter, r *http.Request) {
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
		j1 := r.FormValue("j1")
		j2 := r.FormValue("j2")
		j3 := r.FormValue("j3")
		j4 := r.FormValue("j4")
		j5 := r.FormValue("j5")
		j6 := r.FormValue("j6")
		j7 := r.FormValue("j7")
		j8 := r.FormValue("j8")
		j9 := r.FormValue("j9")
		j10 := r.FormValue("j10")
		j11 := r.FormValue("j11")
		j12 := r.FormValue("j12")
		clean := r.FormValue("clean")
		good := r.FormValue("good")
		ms := r.FormValue("ms")
		amplitude := r.FormValue("amplitude")
		psi := r.FormValue("psi")

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
			`INSERT INTO Jaw_teflon_ultrasonic_state(
				j1, j2, j3, j4, j5, j6, j7, j8, j9, j10, j11, j12,
				jaw_state, teflon_state,
				ultrasonic_time, ultrasonic_amplitude, ultrasonic_pressure, 
				id_sub_header
			)
			VALUES(
				?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?,
				?, ?,
				?, ?, ?,
				? 
			)`)

		if err != nil {
			panic(err.Error())
		}

		insForm.Exec(j1, j2, j3, j4, j5, j6, j7, j8, j9, j10, j11, j12,
			clean, good,
			ms, amplitude, psi,
			SubId)

		defer db.Close()

		log.Println("------FLAG insert jaw_teflon_ultrasonic_state-----")

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

		tmpl.ExecuteTemplate(w, "FillJawTeflonUltrasonicState", data)
	}

}

func GetJawControlV00(w http.ResponseWriter, r *http.Request) {
	db := dbConn()

	nLine := r.URL.Query().Get("line")
	dateInit := r.URL.Query().Get("dinit")
	dateFinal := r.URL.Query().Get("dfinal")

	selDB, err := db.Query(`
	SELECT 
	U.profile_picture,
    J.date_reg,
    SP.turn,
    L.name,
    P.name,
    U.nick_name,
    U.fname,
    U.lname,
    J.j1,
    J.j2,
    J.j3,
    J.j4,
    J.j5,
    J.j6,
    J.j7,
    J.j8,
    J.j9,
    J.j10,
    J.j11,
    J.j12,
    IF(J.jaw_state = 1, 'Si', 'No'),
    IF(J.teflon_state = 1, 'Si', 'No'),
    ultrasonic_time,
    ultrasonic_amplitude,
    ultrasonic_pressure
	FROM
		Jaw_teflon_ultrasonic_state J
			INNER JOIN
		Salsitas_statistical_process SP ON SP.id = J.id_sub_header
			INNER JOIN
		Line L ON L.id = SP.id_line
			INNER JOIN
		Presentation P ON P.id = SP.id_presentation
			INNER JOIN
		User_table U ON U.id = SP.id_operator
		WHERE
				L.id = ?
				AND DATE(J.date_reg) >= ?
				AND DATE(J.date_reg) <  ?`, nLine, dateInit, dateFinal)

	if err != nil {
		panic(err.Error())
	}

	p := models.Jaw_teflon_ultrasonic_state_holder{}
	res := []models.Jaw_teflon_ultrasonic_state_holder{}

	for selDB.Next() {
		var Profile_picture string
		var Date string
		var Turn int
		var Line string
		var Pname string
		var Nname string
		var Fname string
		var Lname string

		var J1 float32
		var J2 float32
		var J3 float32
		var J4 float32
		var J5 float32
		var J6 float32
		var J7 float32
		var J8 float32
		var J9 float32
		var J10 float32
		var J11 float32
		var J12 float32

		var Jaw_state string
		var Teflon_state string

		var Ultrasonic_time float32
		var Ultrasonic_amplitude float32
		var Ultrasonic_pressure float32

		err := selDB.Scan(&Profile_picture, &Date, &Turn, &Line, &Pname, &Nname, &Fname, &Lname,
			&J1, &J2, &J3, &J4, &J5, &J6, &J7, &J8, &J9, &J10, &J11, &J12,
			&Jaw_state, &Teflon_state,
			&Ultrasonic_time, &Ultrasonic_amplitude, &Ultrasonic_pressure)

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

		p.J1 = J1
		p.J2 = J2
		p.J3 = J3
		p.J4 = J4
		p.J5 = J5
		p.J6 = J6
		p.J7 = J7
		p.J8 = J8
		p.J9 = J9
		p.J10 = J10
		p.J11 = J11
		p.J12 = J12

		p.Jaw_state = Jaw_state
		p.Teflon_state = Teflon_state
		p.Ultrasonic_time = Ultrasonic_time
		p.Ultrasonic_amplitude = Ultrasonic_amplitude
		p.Ultrasonic_pressure = Ultrasonic_pressure

		res = append(res, p)
	}
	defer db.Close()
	json.NewEncoder(w).Encode(res)
}

func ConsolidatedJawControl(w http.ResponseWriter, r *http.Request) {

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
	tmpl.ExecuteTemplate(w, "ConsolidatedJawControl", user)
}
