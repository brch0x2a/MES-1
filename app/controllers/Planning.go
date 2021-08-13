package controllers

import (
	//"bytes"
	//	"log"
	//"fmt"
	"fmt"
	"net/http"
	"strconv"

	//"image/jpeg"
	"encoding/json"

	//"../models"

	"github.com/AdrianRb95/MES/app/models"

	_ "github.com/go-sql-driver/mysql"
)

func GetPlanningV00(w http.ResponseWriter, r *http.Request) {
	db := dbConn()

	nDate := r.URL.Query().Get("date")
	nTurn := r.URL.Query().Get("turn")
	nLine := r.URL.Query().Get("line")

	planSTM := `
    SELECT
        Pl.planned,
        Pl.produced
    FROM
        Planning Pl
    WHERE
            Pl.date_planning = ? 
        AND Pl.id_line = ? 
        AND Pl.turn = ?
    `

	if nTurn == "4" {
		planSTM = `
        SELECT
            Pl.planned,
            Pl.produced
        FROM
            Planning Pl
        WHERE
            Pl.date_planning = ? 
            AND Pl.id_line = ? 
            AND COALESCE(4, Pl.turn) = ?
        `
	}

	selDB, err := db.Query(planSTM, nDate, nLine, nTurn)

	if err != nil {
		panic(err.Error())
	}

	p := models.PlanningV00{}

	res := []models.PlanningV00{}
	for selDB.Next() {
		var Planned int
		var Produced int

		err = selDB.Scan(&Planned, &Produced)
		if err != nil {
			panic(err.Error())
		}

		p.Planned += Planned
		p.Produced += Produced

	}
	res = append(res, p)
	defer db.Close()
	json.NewEncoder(w).Encode(res)
}

func GetPlanningV01(w http.ResponseWriter, r *http.Request) {
	db := dbConn()

	nLine := r.URL.Query().Get("line")
	dateInit := r.URL.Query().Get("dinit")
	dateFinal := r.URL.Query().Get("dfinal")

	planSTM := `
    SELECT
        Pl.planned,
        Pl.produced
    FROM
        Planning Pl
    WHERE
            Pl.date_planning = ? 
        AND Pl.id_line = ? 
        AND Pl.turn = ?
    `

	selDB, err := db.Query(planSTM, nLine, dateInit, dateFinal)

	if err != nil {
		panic(err.Error())
	}

	p := models.PlanningV00{}

	res := []models.PlanningV00{}
	for selDB.Next() {
		var Planned int
		var Produced int

		err = selDB.Scan(&Planned, &Produced)
		if err != nil {
			panic(err.Error())
		}

		p.Planned += Planned
		p.Produced += Produced

	}
	res = append(res, p)
	defer db.Close()
	json.NewEncoder(w).Encode(res)
}

func PlanningValidation(w http.ResponseWriter, r *http.Request) {

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

	tmpl.ExecuteTemplate(w, "PlanningValidation", nil)
}

func ReportPlanningbyWeek(w http.ResponseWriter, r *http.Request) {

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

	tmpl.ExecuteTemplate(w, "ReportPlanningbyWeek", user)
}

func EditPlanning(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	nId := r.URL.Query().Get("id")

	selDB, err := db.Query("SELECT * FROM Planning WHERE id=?", nId)
	if err != nil {
		panic(err.Error())
	}
	p := models.Planning{}

	for selDB.Next() {
		var Id int
		var Planned int
		var Produced int
		var Turn int
		var Date_planning string
		var Nominal_speed float32
		var Version int
		var Id_presentation int
		var Id_line int

		err = selDB.Scan(&Id, &Planned, &Produced,
			&Turn, &Nominal_speed, &Id_presentation,
			&Id_line, &Date_planning, &Version)

		if err != nil {
			panic(err.Error())
		}
		p.Id = Id
		p.Planned = Planned
		p.Produced = Produced
		p.Turn = Turn
		p.Date_planning = Date_planning
		p.Nominal_speed = Nominal_speed
		p.Version = Version
		p.Id_presentation = Id_presentation
		p.Id_line = Id_line
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

	tmpl.ExecuteTemplate(w, "EditPlanning", p)
	defer db.Close()
}

func UpdatePlanning(w http.ResponseWriter, r *http.Request) {

	db := dbConn()
	if r.Method == "POST" {

		Id := r.FormValue("uid")
		Line := r.FormValue("line")
		Date := r.FormValue("date_planning")
		Turn := r.FormValue("turn")
		Presentation := r.FormValue("presentation")
		Version := r.FormValue("version")
		Nominal_speed := r.FormValue("nominal_speed")
		Planned := r.FormValue("planned")
		Produced := r.FormValue("produced")

		insForm, err := db.Prepare(`
        UPDATE Planning SET planned=?, produced=?,
        turn=?, nominal_speed=?, id_presentation=?, 
        id_line=?, date_planning=?, version=? 
        WHERE id=?`)

		if err != nil {
			panic(err.Error())
		}

		insForm.Exec(Planned, Produced,
			Turn, Nominal_speed, Presentation,
			Line, Date, Version, Id)

	}

	defer db.Close()

	http.Redirect(w, r, "/planningValidation", 301)

}

func UpdateProducedBox(w http.ResponseWriter, r *http.Request) {

	session, err := store.Get(r, "cookie-name")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	user := GetUser(session)

	Id := r.FormValue("id")
	Box := r.FormValue("box")

	if r.Method == "POST" {

		db := dbConn()
		insForm, err := db.Prepare(`
		UPDATE Planning SET produced = produced + ?
		WHERE id=?`)

		if err != nil {
			panic(err.Error())
		}

		insForm.Exec(Box, Id)

		defer db.Close()

		dbLog := dbConn()
		// LOG
		insFormLog, errLog := dbLog.Prepare(
			`
			insert into ActualProduction (
				id_line,
				id_presentation,
				box,
				id_autor
			)values(
				?, 
				?,
				?,
				?)
			`)

		if errLog != nil {
			panic(errLog.Error())
		}

		insFormLog.Exec(user.Line, user.Presentation,
			Box, user.UserId)

		defer dbLog.Close()

	}

	w.Header().Set("Server", "MES")
	w.WriteHeader(200)
}

func UpdateProducedBoxLog(w http.ResponseWriter, r *http.Request) {

	session, err := store.Get(r, "cookie-name")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	user := GetUser(session)

	fmt.Println("LOG ")

	if r.Method == "POST" {
		db := dbConn()
		Box := r.FormValue("box")
		// LOG
		insForm, err := db.Prepare(
			`
			insert into ActualProduction (
				id_line,
				id_presentation,
				box,
				id_autor
			)values(
				?, 
				?,
				?,
				?)
			`)

		if err != nil {
			panic(err.Error())
		}

		fmt.Println("FLAG - %d\t%d\t%f\t%d", user.Line, user.Presentation, Box, user.UserId)

		insForm.Exec(user.Line, user.Presentation,
			Box, user.UserId)

		defer db.Close()

	}

	w.Header().Set("Server", "MES")
	w.WriteHeader(200)
}

func DeletePlanning(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	p := r.URL.Query().Get("id")
	delForm, err := db.Prepare("DELETE FROM Planning WHERE id=?")
	if err != nil {
		panic(err.Error())
	}
	delForm.Exec(p)

	defer db.Close()
	http.Redirect(w, r, "/planningValidation", 301)
}

func GetPlanningV02(w http.ResponseWriter, r *http.Request) {
	db := dbConn()

	nLine := r.URL.Query().Get("line")
	dateInit := r.URL.Query().Get("dinit")
	dateFinal := r.URL.Query().Get("dfinal")

	planSTM := `
	SELECT 
		Pl.id,
		Pl.date_planning,
		Pl.turn,
		L.name,
		P.name,
		Pl.version,
		Pl.planned,
		Pl.produced,
		Pl.nominal_speed
	FROM
		Planning Pl
			INNER JOIN
		Presentation P ON P.id = Pl.id_presentation
			INNER JOIN
		Line L ON L.id = Pl.id_line
	WHERE
		Pl.id_line = ?
			AND Pl.date_planning BETWEEN ? AND ?

    `

	selDB, err := db.Query(planSTM, nLine, dateInit, dateFinal)

	if err != nil {
		panic(err.Error())
	}

	p := models.Planning_holder{}

	res := []models.Planning_holder{}
	for selDB.Next() {
		var Id int
		var Date_planning string
		var Turn int
		var Line string
		var Presentation string
		var Version int
		var Planned int
		var Produced int
		var Nominal_speed float32

		err = selDB.Scan(&Id, &Date_planning, &Turn,
			&Line, &Presentation, &Version,
			&Planned, &Produced, &Nominal_speed)

		if err != nil {
			panic(err.Error())
		}

		p.Id = Id
		p.Date_planning = Date_planning
		p.Turn = Turn
		p.Line = Line
		p.Presentation = Presentation
		p.Version = Version
		p.Planned = Planned
		p.Produced = Produced
		p.Nominal_speed = Nominal_speed

		res = append(res, p)
	}

	defer db.Close()
	json.NewEncoder(w).Encode(res)
}

func GetActualPlanningBy(w http.ResponseWriter, r *http.Request) {
	db := dbConn()

	line := r.URL.Query().Get("line")
	init := r.URL.Query().Get("init")
	turn := r.URL.Query().Get("turn")

	hoursQ := 6
	minutesQ := 0

	// turn := GetRelativeCurrentTurn()

	// date := time.Now().Format(layoutISO)

	switch turn {

	case "1":
		{
			hoursQ = 6
		}
	case "2":
		{
			hoursQ = 14
		}
	case "3":
		{
			hoursQ = 22
		}

	}
	minutesQ = (hoursQ+8)*60 - 1

	if turn == "4" {
		minutesQ = 1800 - 1
	}

	planSTM := `
	SELECT 
		A.id,
		L.name,
		P.name,
		A.date_reg,
		A.box,
		U.profile_picture,
		U.fname,
		U.lname
	FROM
		ActualProduction A
			INNER JOIN
		Line L ON L.id = A.id_line
			INNER JOIN
		Presentation P ON P.id = A.id_presentation
			INNER JOIN
		User_table U ON U.id = A.id_autor
		Where 
			A.id_line = ? 
			AND A.date_reg BETWEEN DATE_ADD(DATE(?), INTERVAL ` + strconv.Itoa(hoursQ) + ` HOUR) 
			AND DATE_ADD(DATE(?), INTERVAL ` + strconv.Itoa(minutesQ) + ` MINUTE);
    `

	selDB, err := db.Query(planSTM, line, init, init)

	if err != nil {
		panic(err.Error())
	}

	p := models.ActualPlanning{}

	res := []models.ActualPlanning{}
	for selDB.Next() {
		var Id int
		var Date_reg string
		var Line string
		var Presentation string
		var Box float32
		var Profile string
		var Fname string
		var Lname string

		err = selDB.Scan(&Id, &Line, &Presentation,
			&Date_reg, &Box, &Profile,
			&Fname, &Lname)

		if err != nil {
			panic(err.Error())
		}

		p.Id = Id
		p.Date_reg = Date_reg
		p.Line = Line
		p.Presentation = Presentation
		p.Box = Box
		p.Profile = Profile
		p.Fname = Fname
		p.Lname = Lname

		res = append(res, p)
	}

	defer db.Close()
	json.NewEncoder(w).Encode(res)
}
