/*
 *Controlladora encargada de agrupar los historiales de los datos
 *
 */
package controllers

import (
	//"bytes"
	//"log"
	"net/http"
	//"image/jpeg"
	"encoding/json"

	//"../models"

	"github.com/AdrianRb95/MES/app/models"
	_ "github.com/go-sql-driver/mysql"
)

func Consolidado(w http.ResponseWriter, r *http.Request) {
	db := dbConn()

	selDB, err := db.Query(`
        SELECT
            ExL.id,
            ExL.date_event,
            ExL.turn,
            L.name,
            S.description,
            B.description,
            E.description,
            ExL.minutes,
            ExL.note,
            S.color
        FROM
            EventXLine ExL
        INNER JOIN Event E ON
            E.id = ExL.id_event
        INNER JOIN Line L ON
            ExL.id_line = L.id
        INNER JOIN Branch B ON
            E.id_branch = B.id
        INNER JOIN Sub_classification S ON
            S.id = B.id_sub_classification
    `)
	if err != nil {
		panic(err.Error())
	}

	p := models.EventXLine_holder{}

	res := []models.EventXLine_holder{}
	for selDB.Next() {
		var Id int
		var Sub string
		var Branch string
		var Event string
		var Line string
		var Date_event string
		var Turn int
		var Minutes int
		var Note string
		var Color string

		err = selDB.Scan(&Id, &Date_event, &Turn, &Line,
			&Sub, &Branch, &Event,
			&Minutes, &Note, &Color)

		if err != nil {
			panic(err.Error())
		}
		p.Id = Id
		p.Sub = Sub
		p.Branch = Branch
		p.Event = Event
		p.Line = Line
		p.Date_event = Date_event
		p.Turn = Turn
		p.Minutes = Minutes
		p.Note = Note
		p.Color = Color

		res = append(res, p)
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
	tmpl.ExecuteTemplate(w, "Consolidado", res)
	defer db.Close()

}

func HistoricPlanning(w http.ResponseWriter, r *http.Request) {

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

	tmpl.ExecuteTemplate(w, "HistoricPlanning", user)
}

func GetHistoricPlanning(w http.ResponseWriter, r *http.Request) {
	db := dbConn()

	date_event := r.URL.Query().Get("date")
	line := r.URL.Query().Get("line")

	selDB, err := db.Query(`
	SELECT 
		P.id,
		L.name,
		P.turn,
		Pr.name,
		P.version,
		P.nominal_speed,
		P.planned,
		P.produced,
		Pt.photo
	FROM
		Planning P
			INNER JOIN
		Presentation Pr ON Pr.id = P.id_presentation
			INNER JOIN
		Product Pt ON Pt.id = Pr.id_product
			INNER JOIN
		Line L ON L.id = P.id_line
    WHERE
		P.date_planning = ? AND P.id_line = ?
		order by P.turn
    `, date_event, line)

	if err != nil {
		panic(err.Error())
	}

	p := models.Planning_holder{}

	res := []models.Planning_holder{}
	for selDB.Next() {
		var Id int
		// var Date_planning string
		var Turn int
		var Line string
		var Presentation string
		var Version int
		var Planned int
		var Produced int
		var Nominal_speed float32
		// var Box_amount int
		var Photo string

		err = selDB.Scan(&Id, &Line, &Turn, &Presentation, &Version,
			&Nominal_speed, &Planned, &Produced, &Photo)

		if err != nil {
			panic(err.Error())
		}

		p.Id = Id
		p.Planned = Planned
		p.Produced = Produced
		p.Turn = Turn

		p.Nominal_speed = Nominal_speed
		p.Version = Version
		p.Presentation = Presentation
		p.Line = Line

		p.Photo = Photo

		res = append(res, p)
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

	json.NewEncoder(w).Encode(res)
	defer db.Close()

}
