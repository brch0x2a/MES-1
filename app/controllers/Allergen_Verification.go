package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"../models"
)

func AllergenVerificationStep(w http.ResponseWriter, r *http.Request) {
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
	tmpl.ExecuteTemplate(w, "AllergenVerificationStep", data)
}

func InsertAllergenVerification(w http.ResponseWriter, r *http.Request) {
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
		reason := r.FormValue("reason")
		batch := r.FormValue("batch")

		coil_weight := r.FormValue("coil_weight")

		plug := r.FormValue("plug")
		top := r.FormValue("top")
		galon := r.FormValue("galon")

		milk := "0"
		soya := "0"
		gluten := "0"
		egg := "0"

		coil := "#"

		if coil != "" {

			coil = r.FormValue("coil")
		}
		fmt.Printf("\n\n\t\t-- Coil: %s\n\n", coil)

		if r.FormValue("milk") != "" {
			milk = r.FormValue("milk")
		}

		if r.FormValue("soya") != "" {
			soya = r.FormValue("soya")
		}

		if r.FormValue("gluten") != "" {
			gluten = r.FormValue("gluten")
		}

		if r.FormValue("egg") != "" {
			egg = r.FormValue("egg")
		}

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

		//date := strings.Split(SubHeader.Date, "T")[0]
		//date := time.Now()

		insForm, err := db.Prepare(
			`
			INSERT INTO Allergen_Verification(
				id_reason_change,
				batch,
				lote,
				coil_weight,
				lote_plug,
				lote_top,
				lote_gallon,
				front_laminate,
				allergen_laminate,
				milk,
				soya,
				gluten,
				egg,
				id_sub_header
			)
			VALUES(
				?,
				?,
				?,
				?,
				?,
				?,
				?,
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

		coverPath := FileHandlerEngine(r, "capture_frontal", "public/images/crqs/", "crqs-*.png")
		backPath := FileHandlerEngine(r, "capture_info", "public/images/crqs/", "crqs-*.png")

		insForm.Exec(reason, batch, coil, coil_weight,
			plug, top, galon,
			coverPath, backPath,
			milk, soya, gluten, egg,
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

		tmpl.ExecuteTemplate(w, "AllergenVerificationStep", data)
	}
}

func GetAllergenVerificationV00(w http.ResponseWriter, r *http.Request) {
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
		C.batch,
		C.lote,
		C.coil_weight,
		C.lote_plug,
		C.lote_top,
		C.lote_gallon,
		C.front_laminate,
		C.allergen_laminate,
		IF(C.milk = 1, "Si", "No"),
		IF(C.soya = 1, "Si", "No"),
		IF(C.gluten = 1, "Si", "No"),
		IF(C.egg = 1, "Si", "No")
	FROM
		Allergen_Verification C
	INNER JOIN Salsitas_statistical_process SP ON
		SP.id = C.id_sub_header
	INNER JOIN Line L ON
		L.id = SP.id_line
	INNER JOIN Presentation P ON
		P.id = SP.id_presentation
	INNER JOIN User_table U ON
		U.id = SP.id_operator
	WHERE
			L.id = ?
			AND DATE(C.date_reg) >= ?
			AND DATE(C.date_reg) <  ?`, line, dateInit, dateFinal)

	if err != nil {
		panic(err.Error())
	}

	p := models.AllergenVerification_holder{}
	res := []models.AllergenVerification_holder{}

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

		var Batch string
		var Lote string
		var Coil_weight float32

		var Lote_plug string
		var Lote_top string
		var Lote_gallon string

		var Front_laminate string
		var Allergen_laminate string

		var Milk string
		var Soya string
		var Gluten string
		var Egg string

		err := selDB.Scan(&Profile_picture, &Id, &Date_reg, &Turn, &Line, &Presentation, &Nname, &Fname, &Lname,
			&Batch, &Lote, &Coil_weight,
			&Lote_plug, &Lote_top, &Lote_gallon,
			&Front_laminate, &Allergen_laminate,
			&Milk, &Soya, &Gluten, &Egg)

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

		p.Batch = Batch
		p.Lote = Lote
		p.Coil_weight = Coil_weight

		p.Lote_plug = Lote_plug
		p.Lote_top = Lote_top
		p.Lote_gallon = Lote_gallon

		p.Front_laminate = Front_laminate
		p.Allergen_laminate = Allergen_laminate

		p.Milk = Milk
		p.Soya = Soya
		p.Gluten = Gluten
		p.Egg = Egg

		res = append(res, p)
	}
	defer db.Close()
	json.NewEncoder(w).Encode(res)
}

func ConsolidatedAllergenVerification(w http.ResponseWriter, r *http.Request) {

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
	tmpl.ExecuteTemplate(w, "ConsolidatedAllergenVerification", user)
}
