package controllers

import (
	//"bytes"
	//"fmt"

	"encoding/json"
	"net/http"
	"strconv"

	//"image/jpeg"
	"../models"

	_ "github.com/go-sql-driver/mysql"
)

func LossTree(w http.ResponseWriter, r *http.Request) {

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

	tmpl.ExecuteTemplate(w, "LossTree", user)
}

func GetLossTreeData(w http.ResponseWriter, r *http.Request) {
	nId := r.URL.Query().Get("area")
	startDate := r.URL.Query().Get("startDate")
	endDate := r.URL.Query().Get("endDate")

	var lossTree models.LossTree

	lossTree.OEE = CalcOEEAreabyRange(nId, startDate, endDate)
	lossTree.Losses = calcGeneralSegmentbyRange(nId, startDate, endDate)

	json.NewEncoder(w).Encode(lossTree)
}

func GetLossTreeDataGrid(w http.ResponseWriter, r *http.Request) {

	nId := r.URL.Query().Get("area")
	startDate := r.URL.Query().Get("startDate")
	endDate := r.URL.Query().Get("endDate")

	var SubClassification []models.Sub_classification_holder
	var area []models.LossTreeArea
	var lossTreeHolder models.LossTreeArea

	SubClassification = getAllSubClassification()

	for _, s := range SubClassification {
		lossTreeHolder.S = s
		lossTreeHolder.T = getLossBySub(nId, strconv.Itoa(s.Id), startDate, endDate)
		area = append(area, lossTreeHolder)
	}
	json.NewEncoder(w).Encode(area)

}

func GetLossTreeDataGridBy(w http.ResponseWriter, r *http.Request) {

	nId := r.URL.Query().Get("line")
	startDate := r.URL.Query().Get("startDate")
	endDate := r.URL.Query().Get("endDate")
	SubId := r.URL.Query().Get("sub")

	data := getLossBySubAllByLine(SubId, startDate, endDate, nId)

	json.NewEncoder(w).Encode(data)
}

func GetByLineLossTreeDataGridBy(w http.ResponseWriter, r *http.Request) {

	line := r.URL.Query().Get("line")
	startDate := r.URL.Query().Get("startDate")
	endDate := r.URL.Query().Get("endDate")

	data := calcByLineGeneralSegmentbyRange(line, startDate, endDate)

	json.NewEncoder(w).Encode(data)
}

func getLossBySub(area, subClasificationID, startDate, endDate string) []models.TimeSegment {
	var lines []models.Line
	var linesloss []models.TimeSegment

	lines = getLinesByArea(area)

	for _, l := range lines {
		//fmt.Printf(" | %s\n", l.Name)
		loss := getLossBySubLine(subClasificationID, startDate, endDate, l)

		//fmt.Printf("minutes: %f\n", loss.TotalMinutes)
		linesloss = append(linesloss, loss)
	}
	return linesloss
}

func getLossBySubAllByLine(SubClassification, startDate, endDate, line string) []models.TimeSegment {
	db := dbConn()
	var segment models.TimeSegment
	var loss []models.TimeSegment

	selDB, err_mpl := db.Query(`
	SELECT 
    	E.description, S.color, SUM(ExL.minutes) as TotalMinutes
	FROM
		EventXLine ExL
			INNER JOIN
		Event E ON ExL.id_event = E.id
			INNER JOIN
		Branch B ON E.id_branch = B.id
			INNER JOIN
		Sub_classification S ON B.id_sub_classification = S.id
			INNER JOIN
		LineTimeClassification L ON S.id_LTC = L.id
			INNER JOIN
		Line Li ON Li.id = ExL.id_line
	where Li.id = ?
    AND S.id = ?
	AND ExL.date_event BETWEEN DATE_ADD(DATE(?),
	INTERVAL 6 HOUR) AND DATE_ADD(DATE(?),
	INTERVAL 6 HOUR)
	AND L.id != 3
	group by E.id
	order by TotalMinutes desc
`, line, SubClassification, startDate, endDate)

	if err_mpl != nil {
		panic(err_mpl.Error())
	}

	for selDB.Next() {

		var SubCategoryName string
		var totalMinutes float32
		var Color string

		err := selDB.Scan(&SubCategoryName, &Color, &totalMinutes)
		if err != nil {
			panic(err.Error())
		}

		segment.SubCategoryName = SubCategoryName
		segment.Color = Color
		segment.TotalMinutes = totalMinutes

		loss = append(loss, segment)
	}

	defer db.Close()
	return loss
}

func getLossBySubLine(SubClassification, startDate, endDate string, line models.Line) models.TimeSegment {
	db := dbConn()
	var segment models.TimeSegment

	selDB, err_mpl := db.Query(`
	SELECT 
    	S.description, S.color, SUM(ExL.minutes) as TotalMinutes
	FROM
		EventXLine ExL
			INNER JOIN
		Event E ON ExL.id_event = E.id
			INNER JOIN
		Branch B ON E.id_branch = B.id
			INNER JOIN
		Sub_classification S ON B.id_sub_classification = S.id
			INNER JOIN
		LineTimeClassification L ON S.id_LTC = L.id
			INNER JOIN
		Line Li ON Li.id = ExL.id_line

	where Li.id = ?
    AND S.id = ?
	AND ExL.date_event BETWEEN DATE_ADD(DATE(?),
	INTERVAL 6 HOUR) AND DATE_ADD(DATE(?),
	INTERVAL 6 HOUR)
	AND L.id != 3
	group by S.id
	order by TotalMinutes desc
`, line.Id, SubClassification, startDate, endDate)

	if err_mpl != nil {
		panic(err_mpl.Error())
	}

	for selDB.Next() {

		var SubCategoryName string
		var totalMinutes float32
		var Color string

		err := selDB.Scan(&SubCategoryName, &Color, &totalMinutes)
		if err != nil {
			panic(err.Error())
		}

		segment.SubCategoryName = SubCategoryName
		segment.Color = Color
		segment.TotalMinutes = totalMinutes

	}
	segment.Line = line.Name
	segment.LineId = line.Id
	defer db.Close()
	return segment
}
