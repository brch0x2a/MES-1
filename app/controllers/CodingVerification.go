package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"../models"
)

func CodingVerification(w http.ResponseWriter, r *http.Request) {

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

	/*_______________________________________________________*/

	/*_______________________________________________________*/

	tmpl.ExecuteTemplate(w, "CodingVerification", user)
}

func InsertCodingVerification(w http.ResponseWriter, r *http.Request) {

	session, err := store.Get(r, "cookie-name")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	user := GetUser(session)

	if r.Method == "POST" {
		verification_type := r.FormValue("verification_type")
		comment := r.FormValue("comment")

		photo := FileHandlerEngine(r, "capture_codification", "public/images/Coding_Verification/", "codification-*.png")

		db := dbConn()

		insForm, err := db.Prepare(
			`
			insert into Coding_Verification (
				id_operator,
				id_line,
				id_presentation,
				photo_coding,
				verification_type,
				comment, 
				id_checker
			)values(
				?, 
				?,
				?, 
				?,
				?,
				?,
				?
			)
			`)

		if err != nil {
			panic(err.Error())
		}

		insForm.Exec(user.UserId, user.Line, user.Presentation,
			photo, verification_type, comment, 5)

		defer db.Close()

	}
	tmpl.ExecuteTemplate(w, "CodingVerification", user)

}
func ConsolidatedCodingVerification(w http.ResponseWriter, r *http.Request) {

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

	tmpl.ExecuteTemplate(w, "ConsolidatedCodingVerification", user)
}

func GetCodingVerificationBy(w http.ResponseWriter, r *http.Request) {

	init := r.URL.Query().Get("init")
	end := r.URL.Query().Get("end")
	line := r.URL.Query().Get("line")

	db := dbConn()

	selDB, err := db.Query(`
	SELECT 
		V.id,
		V.date_reg,
		O.profile_picture,
		O.fname,
		O.lname,
		L.name,
		P.name,
		V.photo_coding,
		V.verification_type,
		V.comment,
		C.profile_picture,
		C.fname,
		C.lname,
		ifnull(V.date_check, "") as date_check,
        T.name
	FROM
		Coding_Verification V
			INNER JOIN
		User_table O ON O.id = V.id_operator
			INNER JOIN
		Line L ON L.id = V.id_line
			INNER JOIN
		Presentation P ON P.id = V.id_presentation
			INNER JOIN
		User_table C ON C.id = V.id_checker
			inner join
		Transaction_state T on T.id = V.id_state
		WHERE
			V.id_line = ?
			AND V.date_reg BETWEEN ? AND ?
	`, line, init, end)

	if err != nil {
		panic(err.Error())
	}

	p := models.CodingVerification{}

	res := []models.CodingVerification{}

	for selDB.Next() {
		var Id int
		var Date_reg string
		var OFname string
		var OAProfile_Picture string
		var OLname string

		var Line string
		var Presentation string
		var Photo_coding string

		var Type string

		var Comment string

		var Date_check string
		var CFname string
		var CAProfile_Picture string
		var CLname string

		var State string

		err = selDB.Scan(
			&Id,
			&Date_reg,
			&OAProfile_Picture,
			&OFname,
			&OLname,
			&Line,
			&Presentation,
			&Photo_coding,
			&Type,
			&Comment,
			&CAProfile_Picture,
			&CFname,
			&CLname,
			&Date_check,
			&State)

		if err != nil {
			panic(err.Error())
		}

		p.Id = Id
		p.Date_reg = Date_reg
		p.OAProfile_Picture = OAProfile_Picture
		p.OFname = OFname
		p.OLname = OLname
		p.Line = Line
		p.Presentation = Presentation
		p.Photo_coding = Photo_coding
		p.Type = Type
		p.Comment = Comment
		p.Date_check = Date_check
		p.CFname = CFname
		p.CAProfile_Picture = CAProfile_Picture
		p.CLname = CLname
		p.State = State

		res = append(res, p)

	}

	defer db.Close()

	json.NewEncoder(w).Encode(res)
}

func SetCodingVerification(w http.ResponseWriter, r *http.Request) {
	status := 401
	db := dbConn()
	if r.Method == "POST" {

		id := r.FormValue("id")
		user := r.FormValue("user")
		pass := r.FormValue("pass")
		state := r.FormValue("state")

		now := time.Now()

		date_end := now.Format("2006-01-02 15:04:05")

		fmt.Println("\n\n____________Id: %s", id)
		fmt.Println("User: %s", user)
		fmt.Println("Pass: %s", pass)
		fmt.Println("Date: %s", date_end)
		fmt.Println("State: %s", state)

		smt := `
			UPDATE Coding_Verification 
			SET 
				id_checker= ?,
				id_state = ?,
				date_check= ?
			WHERE
				id = ?
			`
		insForm, err := db.Prepare(smt)
		if err != nil {
			panic(err.Error())
		}

		userId := GetUserByNickName(user)

		insForm.Exec(userId, state, date_end, id)
		fmt.Printf("approver:%s\tstate:%s\tdate: %s\tjob:%s\n", userId, state, date_end, id)

		status = 200

	}
	defer db.Close()

	w.Header().Set("Server", "MES")
	w.WriteHeader(status)
}
func GetCodingVerificationE(w http.ResponseWriter, r *http.Request) {

	id := r.URL.Query().Get("id")

	db := dbConn()

	selDB, err := db.Query(`
	SELECT 
		V.id,
		V.date_reg,
		O.profile_picture,
		O.fname,
		O.lname,
		L.name,
		P.name,
		V.photo_coding,
		V.verification_type,
		V.comment,
		C.profile_picture,
		C.fname,
		C.lname,
		ifnull(V.date_check, "") as date_check,
        T.name
	FROM
		Coding_Verification V
			INNER JOIN
		User_table O ON O.id = V.id_operator
			INNER JOIN
		Line L ON L.id = V.id_line
			INNER JOIN
		Presentation P ON P.id = V.id_presentation
			INNER JOIN
		User_table C ON C.id = V.id_checker
			inner join
		Transaction_state T on T.id = V.id_state
		WHERE
			V.id = ?
	`, id)

	if err != nil {
		panic(err.Error())
	}

	p := models.CodingVerification{}

	for selDB.Next() {
		var Id int
		var Date_reg string
		var OFname string
		var OAProfile_Picture string
		var OLname string

		var Line string
		var Presentation string
		var Photo_coding string

		var Type string

		var Comment string

		var Date_check string
		var CFname string
		var CAProfile_Picture string
		var CLname string

		var State string

		err = selDB.Scan(
			&Id,
			&Date_reg,
			&OAProfile_Picture,
			&OFname,
			&OLname,
			&Line,
			&Presentation,
			&Photo_coding,
			&Type,
			&Comment,
			&CAProfile_Picture,
			&CFname,
			&CLname,
			&Date_check,
			&State)

		if err != nil {
			panic(err.Error())
		}

		p.Id = Id
		p.Date_reg = Date_reg
		p.OAProfile_Picture = OAProfile_Picture
		p.OFname = OFname
		p.OLname = OLname
		p.Line = Line
		p.Presentation = Presentation
		p.Photo_coding = Photo_coding
		p.Type = Type
		p.Comment = Comment
		p.Date_check = Date_check
		p.CFname = CFname
		p.CAProfile_Picture = CAProfile_Picture
		p.CLname = CLname
		p.State = State

	}

	defer db.Close()

	json.NewEncoder(w).Encode(p)
}
