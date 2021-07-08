package controllers

import (
	//"bytes"

	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"../models"
	//"image/jpeg"
	//	_ "github.com/go-sql-driver/mysql"
)

func SpecificDashBoard(w http.ResponseWriter, r *http.Request) {

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

	turn := GetRelativeCurrentTurn(time.Now())

	fmt.Println("Turn: %s", turn)

	tmpl.ExecuteTemplate(w, "SpecificDashboard", user)

}

func GeneralDashboard(w http.ResponseWriter, r *http.Request) {

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

	turn := GetRelativeCurrentTurn(time.Now())

	fmt.Println("Turn: %s", turn)

	tmpl.ExecuteTemplate(w, "GeneralDashboard", user)

}

func DashBoardCounterBy(w http.ResponseWriter, r *http.Request) {

	line := r.URL.Query().Get("line")

	date := time.Now().Format(layoutISO)

	blueCount := MttTagCount(line, "1", "1")
	redCount := MttTagCount(line, "2", "1")
	greenCount := SHETagCount(line, "5", "1")
	orangeCount := QATagCount(line, "4", "1")
	oee := calcOEE(line, "4", date)
	wo := WorkOrderCounter(line, "1")

	dashBoardHolder := models.DashBoardHolder{}

	dashBoardHolder.Blue = blueCount
	dashBoardHolder.Red = redCount
	dashBoardHolder.Green = greenCount
	dashBoardHolder.Orange = orangeCount
	dashBoardHolder.OEE = oee
	dashBoardHolder.WO = wo

	json.NewEncoder(w).Encode(dashBoardHolder)
}
