package controllers

import (
	//"bytes"
	//"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"

	//"image/jpeg"
	"encoding/json"

	"../models"
	_ "github.com/go-sql-driver/mysql"
)

func FillSubHeaderPChemical(w http.ResponseWriter, r *http.Request) {

	db := dbConn()

	selDB, err := db.Query("Select * from Header where id=1")

	if err != nil {
		panic(err.Error())
	}

	p := models.Header{}

	for selDB.Next() {
		var Id int
		var Name string
		var Cod_doc string
		var Revision_date string
		var Next_revision_date string
		var Revision_no int

		err = selDB.Scan(&Id, &Name, &Cod_doc, &Revision_date, &Next_revision_date, &Revision_no)

		if err != nil {
			panic(err.Error())
		}

		p.Id = Id
		p.Name = Name
		p.Cod_doc = Cod_doc
		p.Revision_date = Revision_date
		p.Next_revision_date = Next_revision_date
		p.Revision_no = Revision_no

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

	data := map[string]interface{}{
		"Header": p,
		"User":   user,
	}

	tmpl.ExecuteTemplate(w, "FillSubHeaderPChemical", data)
	defer db.Close()
}

func InsertSubPChemical(w http.ResponseWriter, r *http.Request) {

	SubHeader := models.Physicochemical_subheader_holder{}
	var LastId int

	db := dbConn()
	if r.Method == "POST" {
		header := r.FormValue("header")
		product := r.FormValue("product")
		area := r.FormValue("area")
		fryma := r.FormValue("fryma")

		insForm, err := db.Prepare(`
		INSERT INTO 
		Physicochemical_subheader(id_header, id_product, 
			id_area, boula_fryma) 
		VALUES(?, ?, ?, ?)`)

		if err != nil {
			panic(err.Error())
		}

		res, err := insForm.Exec(header, product, area, fryma)
		if err != nil {
			log.Fatal("Cannot run insert statement", err)
		}

		id, _ := res.LastInsertId()

		fmt.Printf("Inserted row: %d", id)

		LastId = int(id)
	}
	defer db.Close()

	fmt.Printf("Last Id: %d", LastId)

	selDB_sub, err := db.Query(`
				SELECT
				P.id,
				P.id_header,
				Pr.name,
				A.Name,
				P.boula_fryma
			FROM
				Physicochemical_subheader P
			INNER JOIN Product Pr ON
				Pr.id = P.id_product
			INNER JOIN Area A ON
				A.Id = P.id_area
			WHERE
				P.id = ?
	`, LastId)

	if err != nil {
		panic(err.Error())
	}

	for selDB_sub.Next() {
		var Id int
		var Boula_fryma string

		var Product string
		var Area string
		var Header int

		err = selDB_sub.Scan(&Id, &Header, &Product,
			&Area, &Boula_fryma)

		if err != nil {
			panic(err.Error())
		}

		SubHeader.Id = Id
		SubHeader.Boula_fryma = Boula_fryma
		SubHeader.Product = Product
		SubHeader.Area = Area
		SubHeader.Header = Header

	}
	log.Printf("Area: %s", SubHeader.Area)
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

	data := map[string]interface{}{
		"SubHeader": SubHeader,
		"User":      user,
	}

	tmpl.ExecuteTemplate(w, "FillPChemicalGeneral", data)
}

func FillPChemicalGeneral(w http.ResponseWriter, r *http.Request) {

}

func GetSensorial_analysis_scale(w http.ResponseWriter, r *http.Request) {
	db := dbConn()

	selDB, err := db.Query("Select * from Sensorial_analysis_scale")

	if err != nil {
		panic(err.Error())
	}

	p := models.Sensorial_analysis_scale{}
	res := []models.Sensorial_analysis_scale{}

	for selDB.Next() {
		var Id int
		var Description string

		err = selDB.Scan(&Id, &Description)
		if err != nil {
			panic(err.Error())
		}

		p.Id = Id
		p.Description = Description

		res = append(res, p)
	}

	json.NewEncoder(w).Encode(res)

	defer db.Close()
}

func GetDesition_catalog(w http.ResponseWriter, r *http.Request) {
	db := dbConn()

	selDB, err := db.Query("Select * from Desition_catalog")

	if err != nil {
		panic(err.Error())
	}

	p := models.Desition_catalog{}
	res := []models.Desition_catalog{}

	for selDB.Next() {
		var Id int
		var Description string

		err = selDB.Scan(&Id, &Description)
		if err != nil {
			panic(err.Error())
		}

		p.Id = Id
		p.Description = Description

		res = append(res, p)
	}

	json.NewEncoder(w).Encode(res)

	defer db.Close()
}

func GetPChemicalGeneral(w http.ResponseWriter, r *http.Request) {
	db := dbConn()

	selDB, err := db.Query(`
	SELECT
    id,
    batch,
    tank,
    presentation_density,
    date_reg,
    brix,
    ph,
    ph_pcc,
    chloride,
    consistency_plummet,
    consistency_homogenizer,
    appearance,
    color,
    aroma,
    taste,
    id_analyst,
    id_desition,
    id_subheader
FROM
	Physicochemical_general`)

	if err != nil {
		panic(err.Error())
	}

	p := models.Physicochemical_general{}
	res := []models.Physicochemical_general{}

	for selDB.Next() {
		var Id int
		var Batch int

		var tank string

		var Presentation_density float32
		var Date string

		/*Analisis Fisico-Quimico*/
		var Brix float32
		var Ph float32
		var Ph_pcc float32
		var Chloride float32
		var Consistency_plummet float32
		var Consistency_homogenizer float32

		/*Analisis Sensorial*/
		var Appearance int
		var Color int
		var Aroma int
		var Taste int

		var Analyst int
		var Desition int
		var Subheader int

		err = selDB.Scan(&Id, &Batch, &tank, &Presentation_density,
			&Date, &Brix, &Ph, &Ph_pcc, &Chloride, &Consistency_plummet,
			&Consistency_homogenizer, &Appearance, &Color, &Aroma,
			&Taste, &Analyst, &Desition, &Subheader)

		if err != nil {
			panic(err.Error())
		}

		p.Id = Id
		p.Batch = Batch
		p.Tank = tank
		p.Presentation_density = Presentation_density
		p.Date = Date
		p.Brix = Brix
		p.Ph = Ph
		p.Ph_pcc = Ph_pcc
		p.Chloride = Chloride
		p.Consistency_plummet = Consistency_plummet
		p.Consistency_homogenizer = Consistency_homogenizer
		p.Appearance = Appearance
		p.Color = Color
		p.Aroma = Aroma
		p.Taste = Taste
		p.Analyst = Analyst
		p.Desition = Desition
		p.Subheader = Subheader

		res = append(res, p)
	}

	json.NewEncoder(w).Encode(res)

	defer db.Close()

}

func InsertPChemicalGeneral(w http.ResponseWriter, r *http.Request) {

	SubHeader := models.Physicochemical_subheader_holder{}

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
		tank := r.FormValue("tank")
		density := r.FormValue("density")
		date := r.FormValue("vdate")

		/*Analisis Fisico-Quimico*/
		brix := r.FormValue("brix")
		ph := r.FormValue("ph")
		ph_pcc := r.FormValue("ph_pcc")
		chloride := r.FormValue("chloride")
		consistency_plummet := r.FormValue("plummet")
		consistency_homogenizer := r.FormValue("homogenizer")

		/*Analisis Sensorial*/
		appearance := r.FormValue("appearance")
		color := r.FormValue("color")
		aroma := r.FormValue("aroma")
		taste := r.FormValue("taste")

		analyst := user.UserId
		desition := 4
		subheader := r.FormValue("head")

		log.Printf("sub id %s\n", subheader)
		SubId, _ = strconv.Atoi(subheader)

		insForm, err := db.Prepare(`
		INSERT INTO Physicochemical_general(
			batch,
			tank,
			presentation_density,
			date_reg,
			brix,
			ph,
			ph_pcc,
			chloride,
			consistency_plummet,
			consistency_homogenizer,
			appearance,
			color,
			aroma,
			taste,
			id_analyst,
			id_desition,
			id_subheader
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
			?,
			?,
			?,
			?
		)
		`)

		if err != nil {
			panic(err.Error())
		}

		insForm.Exec(
			batch, tank, density, date, brix, ph, ph_pcc,
			chloride, consistency_plummet, consistency_homogenizer, appearance, color,
			aroma, taste, analyst, desition, subheader)

	}
	defer db.Close()

	selDB_sub, err := db.Query(`
				SELECT
				P.id,
				P.id_header,
				Pr.name,
				A.Name,
				P.boula_fryma
			FROM
				Physicochemical_subheader P
			INNER JOIN Product Pr ON
				Pr.id = P.id_product
			INNER JOIN Area A ON
				A.Id = P.id_area
			WHERE
				P.id = ?
	`, SubId)

	if err != nil {
		panic(err.Error())
	}

	for selDB_sub.Next() {
		var Id int
		var Boula_fryma string

		var Product string
		var Area string
		var Header int

		err = selDB_sub.Scan(&Id, &Header, &Product,
			&Area, &Boula_fryma)

		if err != nil {
			panic(err.Error())
		}

		SubHeader.Id = Id
		SubHeader.Boula_fryma = Boula_fryma
		SubHeader.Product = Product
		SubHeader.Area = Area
		SubHeader.Header = Header

	}
	log.Printf("Area: %s", SubHeader.Area)

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
	defer db.Close()
	http.Redirect(w, r, "/fillSubHeaderPChemical", 301)

}

/*______________________Validation_________________________________*/
func ValidationPChemical(w http.ResponseWriter, r *http.Request) {
	db := dbConn()

	selDB, err := db.Query("Select * from Header where id=1")

	if err != nil {
		panic(err.Error())
	}

	p := models.Header{}

	for selDB.Next() {
		var Id int
		var Name string
		var Cod_doc string
		var Revision_date string
		var Next_revision_date string
		var Revision_no int

		err = selDB.Scan(&Id, &Name, &Cod_doc, &Revision_date, &Next_revision_date, &Revision_no)

		if err != nil {
			panic(err.Error())
		}

		p.Id = Id
		p.Name = Name
		p.Cod_doc = Cod_doc
		p.Revision_date = Revision_date
		p.Next_revision_date = Next_revision_date
		p.Revision_no = Revision_no

	}
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

	tmpl.ExecuteTemplate(w, "ValidationPChemical", p)
	defer db.Close()

}

func GetPChemicalV00(w http.ResponseWriter, r *http.Request) {

	area_id := r.URL.Query().Get("area_id")
	iDate := r.URL.Query().Get("idate")
	fDate := r.URL.Query().Get("fdate")

	db := dbConn()

	selDB, err := db.Query(`
		SELECT
			U.profile_picture,
			A.name,
			U.nick_name,
			U.fname,
			U.lname,
			Ph.date_reg,
			P.name,
			Ps.boula_fryma,
			Ph.batch,
			Ph.appearance,
			Sa.description,
			Ph.color,
			Sc.description,
			Ph.aroma,
			Sr.description,
			Ph.taste,
			St.description,
			Ph.ph,
			Ph.ph_pcc,
			Ph.chloride,
			Ph.brix,
			Ph.consistency_plummet,
			Ph.consistency_homogenizer,
			Ph.presentation_density,
			Ph.tank
			#D.description
		FROM
			Physicochemical_general Ph
				INNER JOIN
			Sensorial_analysis_scale Sa ON Sa.id = Ph.appearance
				INNER JOIN
			Sensorial_analysis_scale Sc ON Sc.id = Ph.color
				INNER JOIN
			Sensorial_analysis_scale Sr ON Sr.id = Ph.aroma
				INNER JOIN
			Sensorial_analysis_scale St ON St.id = Ph.taste
				INNER JOIN
			User_table U ON U.id = Ph.id_analyst
				INNER JOIN
			Desition_catalog D ON D.id = Ph.id_desition
				INNER JOIN
			Physicochemical_subheader Ps ON Ps.id = Ph.id_subheader
				INNER JOIN
			Product P ON P.id = Ps.id_product
				INNER JOIN
			Area A ON A.id = Ps.id_area
				WHERE
				? <= DATE(Ph.date_reg)
				AND DATE(Ph.date_reg) <= ?
				AND A.id = ?
				ORDER BY date_reg DESC`, iDate, fDate, area_id)

	if err != nil {
		panic(err.Error())
	}

	p := models.Physicochemical_general_holder{}

	res := []models.Physicochemical_general_holder{}
	for selDB.Next() {
		var Profile_picture string

		var Area string

		var Analystnn string
		var Analystfname string
		var Analystlname string

		var Date string
		var Product string

		var Fryma float32
		var Batch int

		/*Analisis Sensorial*/
		var Appearance_level int
		var Appearance string
		var Color_level int
		var Color string
		var Aroma_level int
		var Aroma string
		var Taste_level int
		var Taste string

		/*Analisis Fisico-Quimico*/
		var Ph float32
		var Ph_pcc float32
		var Chloride float32
		var Brix float32
		var Consistency_plummet float32
		var Consistency_homogenizer float32
		var Presentation_density float32

		var Tank string

		err = selDB.Scan(&Profile_picture, &Area, &Analystnn, &Analystfname, &Analystlname,
			&Date, &Product, &Fryma, &Batch, &Appearance_level, &Appearance,
			&Color_level, &Color, &Aroma_level, &Aroma, &Taste_level,
			&Taste, &Ph, &Ph_pcc, &Chloride, &Brix, &Consistency_plummet,
			&Consistency_homogenizer, &Presentation_density, &Tank)

		if err != nil {
			panic(err.Error())
		}

		p.Profile_picture = Profile_picture
		p.Area = Area
		p.Analystnn = Analystnn
		p.Analystfname = Analystfname
		p.Analystlname = Analystlname
		p.Date = Date
		p.Product = Product
		p.Fryma = Fryma
		p.Batch = Batch

		p.Appearance_level = Appearance_level
		p.Appearance = Appearance
		p.Color_level = Color_level
		p.Color = Color
		p.Aroma_level = Aroma_level
		p.Aroma = Aroma
		p.Taste_level = Taste_level
		p.Taste = Taste

		p.Ph = Ph
		p.Ph_pcc = Ph_pcc
		p.Chloride = Chloride
		p.Brix = Brix
		p.Consistency_plummet = Consistency_plummet
		p.Consistency_homogenizer = Consistency_homogenizer
		p.Presentation_density = Presentation_density
		p.Tank = Tank

		res = append(res, p)
	}
	defer db.Close()
	json.NewEncoder(w).Encode(res)

}

func ConsolidatedPChemical(w http.ResponseWriter, r *http.Request) {

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

	tmpl.ExecuteTemplate(w, "ConsolidatedPChemical", user)
}
