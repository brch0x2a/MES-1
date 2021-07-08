package controllers

import (
	//"bytes"
	//"log"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"time"

	//"image/jpeg"
	"encoding/json"

	"../models"
	//	_ "github.com/go-sql-driver/mysql"
)

func ConsolidatedAMIS(w http.ResponseWriter, r *http.Request) {

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
	tmpl.ExecuteTemplate(w, "ConsolidatedAMIS", user)
}

func GetAMISV00(w http.ResponseWriter, r *http.Request) {

	area := r.URL.Query().Get("area")

	idate := r.URL.Query().Get("idate")
	fdate := r.URL.Query().Get("fdate")

	//var goodProduct float32
	var amis []models.AMIS_holder

	//goodProduct = GetGoodVolumeThousandByLineXRange("1", idate, fdate)
	//amis = calcAMISbyRange("1", idate, fdate)

	//==========================================================
	fmt.Printf("_____________________________________________\n")
	fmt.Printf("		AMIS\n")
	fmt.Printf("_____________________________________________\n")
	/*
		fmt.Printf(">>GoodProduct:\t\t%f\n", goodProduct)
		fmt.Printf(">>TotalTime:\t\t%f\n", amis.TotalTime)
		fmt.Printf(">>OT:\t\t%f\n", amis.OT)
		fmt.Printf(">>LT:\t\t%f\n", amis.LT)
		fmt.Printf(">>OEE:\t\t%f\n", amis.OEE)
		fmt.Printf(">>MP:\t\t%f\n", amis.MP)
		fmt.Printf(">>SP:\t\t%f\n", amis.SP)
		fmt.Printf(">>AU:\t\t%f\n", amis.AU)
	*/

	Line := getLinesByArea(area)

	fmt.Println("area: ", area)

	for _, line := range Line {

		lineID := strconv.Itoa(line.Id)

		fmt.Println(lineID)

		amis = append(amis, calcAMISbyRange(lineID, idate, fdate))
	}

	fmt.Printf("_____________________________________________\n")
	fmt.Printf("		END AMIS\n")
	fmt.Printf("_____________________________________________\n")
	//==========================================================

	json.NewEncoder(w).Encode(amis)
}

func calcAMISbyRange(line, startDate, endDate string) models.AMIS_holder {

	layout := "2006-01-02"

	tinit, err := time.Parse(layout, startDate)
	tend, err := time.Parse(layout, endDate)

	if err != nil {
		fmt.Println(err)
	}

	diff := float32(tend.Sub(tinit).Hours() / 24)

	var amis_holder models.AMIS_holder

	var vot float32
	var mpl float32
	var pdl float32
	var oline models.Line

	amis_holder.GoodVolume = calcVolumeProduced(line, startDate, endDate)

	amis_holder.TotalTime = calcTotalTime(8, 3, diff, 1) * 60

	vot = calcVOTbyRange(line, startDate, endDate)
	amis_holder.LegalLosses = calcLegalLossesByRange(line, startDate, endDate)
	amis_holder.UCL = calcUnutilisedCapacityLosses(line, startDate, endDate)
	amis_holder.UCLSegment = calcUCLSegmentbyRange(line, startDate, endDate)

	mpl = calcMPLbyRange(line, startDate, endDate)
	amis_holder.MPLSegment = calcMPLSegmentbyRange(line, startDate, endDate)

	pdl = calcPDLbyRange(line, startDate, endDate)

	amis_holder.PDLSegment = calcPDLSegmentbyRange(line, startDate, endDate)

	amis_holder.PDL = pdl
	amis_holder.MPL = mpl

	amis_holder.VOT = vot
	amis_holder.AT = amis_holder.TotalTime - amis_holder.LegalLosses
	amis_holder.ALT = amis_holder.AT - amis_holder.UCL
	amis_holder.OT = vot + mpl/60
	amis_holder.LT = amis_holder.OT + pdl/60
	amis_holder.IdleTime = amis_holder.ALT - amis_holder.LT

	if math.IsNaN(float64((vot / amis_holder.LT * 100))) {
		amis_holder.OEE = 0.0
	} else {
		amis_holder.OEE = vot / amis_holder.LT * 100
	}

	// --------------------------
	if math.IsNaN(float64((vot / amis_holder.OT * 100))) {
		amis_holder.MP = 0.0
	} else {
		amis_holder.MP = vot / amis_holder.OT * 100
	}

	if math.IsNaN(float64((amis_holder.OT / amis_holder.LT * 100))) {
		amis_holder.SP = 0.0
	} else {
		amis_holder.SP = amis_holder.OT / amis_holder.LT * 100
	}

	if math.IsNaN(float64((vot / amis_holder.TotalTime * 100))) {
		amis_holder.AU = 0.0
	} else {
		amis_holder.AU = vot / (amis_holder.TotalTime / 60) * 100
	}

	if math.IsNaN(float64((amis_holder.LT / amis_holder.AT * 100))) {
		amis_holder.UCU = 0.0
	} else {
		amis_holder.UCU = amis_holder.LT / (amis_holder.AT / 60) * 100
	}

	if math.IsNaN(float64((amis_holder.LT / amis_holder.ALT * 100))) {
		amis_holder.CCU = 0.0
	} else {
		amis_holder.CCU = amis_holder.LT / (amis_holder.ALT / 60) * 100
	}

	amis_holder.Breakdowns = calcBreakdownsRange(line, startDate, endDate)

	amis_holder.ProductChangeovers = calcChangeoversRange(line, "79", startDate, endDate)
	amis_holder.FormatChangeovers = calcChangeoversRange(line, "80", startDate, endDate) + calcChangeoversRange(line, "81", startDate, endDate)
	amis_holder.GeneralChangeovers = amis_holder.ProductChangeovers + amis_holder.FormatChangeovers

	amis_holder.ProductChangeoversTime = calcChangeoversTimeRange(line, "79", startDate, endDate) / 60
	amis_holder.FormatChangeoversTime = (calcChangeoversTimeRange(line, "80", startDate, endDate) + calcChangeoversTimeRange(line, "81", startDate, endDate)) / 60
	amis_holder.GeneralChangeoversTime = amis_holder.ProductChangeoversTime + amis_holder.FormatChangeoversTime

	amis_holder.OR = calcORbyRange(line, startDate, endDate)

	lineID, _ := strconv.Atoi(line)

	oline = GetLineObjectBy(lineID)

	amis_holder.Line = oline.Name

	return amis_holder
}

//--------------------------------------------------
func GetGoodVolumeThousandByLineXRange(line, startDate, endDate string) float32 {
	var goodVolumeThounsand float32
	var goodProduct float32
	var plannedData []models.VOT
	var goodVolume float32
	var boxAmount float32

	plannedData = GetPlannedDataByLineXRange(line, startDate, endDate)

	for i := 0; i < len(plannedData); i++ {
		goodVolume += plannedData[i].GoodVolume
		boxAmount += plannedData[i].Box_amount
	}

	goodProduct = goodVolume * boxAmount

	goodVolumeThounsand = goodProduct * .001

	return goodVolumeThounsand
}

func GetPlannedDataByLineXRange(line, startDate, endDate string) []models.VOT {

	votSTM := `
        SELECT
        Pl.produced,
        Pl.nominal_speed,
        Pr.box_amount
        FROM
            Planning Pl
        INNER JOIN Presentation Pr ON
            Pl.id_presentation = Pr.id
        WHERE
		Pl.id_line = ?
        AND Pl.date_planning BETWEEN ? AND ?
    `

	db := dbConn()

	selDB, err := db.Query(votSTM, line, startDate, endDate)

	if err != nil {
		panic(err.Error())
	}

	p := models.VOT{}

	res := []models.VOT{}
	for selDB.Next() {

		var GoodVolume float32
		var NominalSpeed float32
		var Box_amount float32

		err = selDB.Scan(&GoodVolume, &NominalSpeed, &Box_amount)
		if err != nil {
			panic(err.Error())
		}
		p.GoodVolume = GoodVolume
		p.NominalSpeed = NominalSpeed
		p.Box_amount = Box_amount

		res = append(res, p)
	}
	defer db.Close()

	return res
}

func calcUnutilisedCapacityLosses(line, startDate, endDate string) float32 {

	db := dbConn()
	var ucl_calc float32

	selDB, err := db.Query(`
	SELECT 
		ExL.minutes
	FROM
    EventXLine ExL
        INNER JOIN
    Event E ON ExL.id_event = E.id
        INNER JOIN
    Branch B ON E.id_branch = B.id
		INNER JOIN Sub_classification S ON
	B.id_sub_classification = S.id
    where S.id = 22
	AND B.id != 52
    AND ExL.id_line = ?
	AND ExL.date_event  BETWEEN DATE_ADD(DATE(?),
	INTERVAL 6 HOUR) AND DATE_ADD(DATE(?),
	INTERVAL 6 HOUR)
	`, line, startDate, endDate)

	if err != nil {
		panic("calcUnutilisedCapacityLosses: " + err.Error())
	}

	for selDB.Next() {

		var Minutes float32

		err := selDB.Scan(&Minutes)
		if err != nil {
			panic(err.Error())
		}

		ucl_calc += Minutes
	}
	defer db.Close()

	if math.IsNaN(float64(ucl_calc)) {
		ucl_calc = 0.0
	}

	return ucl_calc
}

func calcLegalLossesByRange(line, startDate, endDate string) float32 {
	db := dbConn()
	var legalLosses_calc float32

	selDB, err := db.Query(`
	SELECT 
		ExL.minutes
	FROM
		EventXLine ExL
			INNER JOIN
		Event E ON ExL.id_event = E.id
			INNER JOIN
		Branch B ON E.id_branch = B.id
		where B.id = 52
		AND ExL.id_line = ?
		AND ExL.date_event  BETWEEN DATE_ADD(DATE(?),
		INTERVAL 6 HOUR) AND DATE_ADD(DATE(?),
		INTERVAL 6 HOUR)
	`, line, startDate, endDate)

	if err != nil {
		panic("Calc Legal Losses: " + err.Error())
	}

	for selDB.Next() {

		var Minutes float32

		err := selDB.Scan(&Minutes)
		if err != nil {
			panic(err.Error())
		}

		legalLosses_calc += Minutes
	}
	defer db.Close()

	if math.IsNaN(float64(legalLosses_calc)) {
		legalLosses_calc = 0.0
	}

	return legalLosses_calc
}

func calcTotalTime(hourTurn, numberOfTurn, daysOfWeek, amountWeek float32) float32 {

	return hourTurn * numberOfTurn * daysOfWeek * amountWeek
}

func calcVolumeProduced(line, startDate, endDate string) float32 {
	votSTM := `
        SELECT
        Pl.produced,
        Pl.nominal_speed,
        Pr.box_amount
        FROM
            Planning Pl
        INNER JOIN Presentation Pr ON
            Pl.id_presentation = Pr.id
        WHERE
		Pl.id_line = ?
        AND Pl.date_planning BETWEEN ? AND ?
    `

	db := dbConn()

	var volume float32

	//nId := r.URL.Query().Get("id")
	selDB, err := db.Query(votSTM, line, startDate, endDate)

	if err != nil {
		panic(err.Error())
	}

	for selDB.Next() {

		var GoodVolume float32
		var NominalSpeed float32
		var Box_amount float32
		var calc float32

		err = selDB.Scan(&GoodVolume, &NominalSpeed, &Box_amount)
		if err != nil {
			panic(err.Error())
		}

		calc = (GoodVolume * Box_amount) / 1000

		volume += calc
	}
	defer db.Close()

	if math.IsNaN(float64(volume)) {
		volume = 0.0
	}

	return volume
}

func calcVOTbyRange(line, startDate, endDate string) float32 {
	votSTM := `
        SELECT
        Pl.produced,
        Pl.nominal_speed,
        Pr.box_amount
        FROM
            Planning Pl
        INNER JOIN Presentation Pr ON
            Pl.id_presentation = Pr.id
        WHERE
		Pl.id_line = ?
        AND Pl.date_planning BETWEEN ? AND DATE_ADD(?, INTERVAL - 24 HOUR)
    `

	db := dbConn()

	var vot_calc float32

	//nId := r.URL.Query().Get("id")
	selDB, err := db.Query(votSTM, line, startDate, endDate)

	if err != nil {
		panic(err.Error())
	}

	p := models.VOT{}

	res := []float32{}
	for selDB.Next() {

		var GoodVolume float32
		var NominalSpeed float32
		var Box_amount float32
		var calc float32

		err = selDB.Scan(&GoodVolume, &NominalSpeed, &Box_amount)
		if err != nil {
			panic(err.Error())
		}
		p.GoodVolume = GoodVolume
		p.NominalSpeed = NominalSpeed
		p.Box_amount = Box_amount

		calc = (GoodVolume * Box_amount) / NominalSpeed / 60

		//log.Printf("VOT: %.6f\n", calc)

		calc = p.Get() / 60
		vot_calc += calc

		//fmt.Printf("vot: %f\n\n", vot_calc)

		res = append(res, calc)
	}
	defer db.Close()

	if math.IsNaN(float64(vot_calc)) {
		vot_calc = 0.0
	}

	return vot_calc
}

func calcMPLbyRange(line, startDate, endDate string) float32 {

	db := dbConn()
	var mpl_calc float32

	/*-------Manufactory Performance Lose-----------*/
	fmt.Printf("MPL")

	selDB_mpl, err_mpl := db.Query(`
    SELECT
		ExL.minutes
		FROM
			EventXLine ExL
		INNER JOIN Event E ON
			ExL.id_event = E.id
		INNER JOIN Branch B ON
			E.id_branch = B.id
		INNER JOIN Sub_classification S ON
			B.id_sub_classification = S.id
		INNER JOIN LineTimeClassification L ON
			S.id_LTC = L.id
    WHERE
        L.id = 1  
        AND ExL.id_line = ?
        AND ExL.date_event   BETWEEN DATE_ADD(DATE(?),
		INTERVAL 6 HOUR) AND DATE_ADD(DATE(?),
		INTERVAL 359 MINUTE)
    `, line, startDate, endDate)

	if err_mpl != nil {
		panic(err_mpl.Error())
	}

	for selDB_mpl.Next() {

		var Minutes float32

		err := selDB_mpl.Scan(&Minutes)
		if err != nil {
			panic(err.Error())
		}

		mpl_calc += Minutes
	}
	defer db.Close()

	if math.IsNaN(float64(mpl_calc)) {
		mpl_calc = 0.0
	}

	return mpl_calc
}

func calcMPLSegmentbyRange(line, startDate, endDate string) []models.
	TimeSegment {

	db := dbConn()
	var segment models.TimeSegment
	var mpl_calc []models.TimeSegment

	/*-------Manufactory Performance Lose-----------*/
	fmt.Printf("MPL")

	selDB_mpl, err_mpl := db.Query(`
	SELECT 
		S.description, S.color, SUM(ExL.minutes)
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
	WHERE
		L.id = 1 AND ExL.id_line = ?
			AND ExL.date_event  BETWEEN DATE_ADD(DATE(?),
			INTERVAL 6 HOUR) AND DATE_ADD(DATE(?),
			INTERVAL 6 HOUR)
	GROUP BY S.id
    `, line, startDate, endDate)

	if err_mpl != nil {
		panic(err_mpl.Error())
	}

	for selDB_mpl.Next() {

		var SubCategoryName string
		var totalMinutes float32
		var Color string

		err := selDB_mpl.Scan(&SubCategoryName, &Color, &totalMinutes)
		if err != nil {
			panic(err.Error())
		}

		segment.SubCategoryName = SubCategoryName
		segment.Color = Color
		segment.TotalMinutes = totalMinutes

		mpl_calc = append(mpl_calc, segment)
	}
	defer db.Close()
	return mpl_calc
}

func calcPDLbyRange(line, startDate, endDate string) float32 {

	db := dbConn()

	var pdl_calc float32

	/*-------Process Driven Losses-----------*/
	fmt.Printf("PDL")

	selDB_pdl, err_pdl := db.Query(`
		SELECT
			ExL.minutes
			FROM
				EventXLine ExL
			INNER JOIN Event E ON
				ExL.id_event = E.id
			INNER JOIN Branch B ON
				E.id_branch = B.id
			INNER JOIN Sub_classification S ON
				B.id_sub_classification = S.id
			INNER JOIN LineTimeClassification L ON
				S.id_LTC = L.id
			WHERE
				L.id = 2 
				AND ExL.id_line = ?
				AND ExL.date_event  BETWEEN DATE_ADD(DATE(?),
				INTERVAL 6 HOUR) AND DATE_ADD(DATE(?),
				INTERVAL 6 HOUR)
	`, line, startDate, endDate)

	if err_pdl != nil {
		panic(err_pdl.Error())
	}

	for selDB_pdl.Next() {

		var Minutes float32

		err := selDB_pdl.Scan(&Minutes)
		if err != nil {
			panic(err.Error())
		}

		pdl_calc += Minutes
	}

	defer db.Close()

	if math.IsNaN(float64(pdl_calc)) {
		pdl_calc = 0.0
	}

	return pdl_calc
}

func calcPDLSegmentbyRange(line, startDate, endDate string) []models.
	TimeSegment {

	db := dbConn()
	var segment models.TimeSegment
	var pdl_calc []models.TimeSegment

	/*-------Manufactory Performance Lose-----------*/
	fmt.Printf("PDL")

	selDB, err_mpl := db.Query(`
	SELECT 
    S.description, S.color, SUM(ExL.minutes)
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
	WHERE
		L.id = 2 AND ExL.id_line = ?
			AND ExL.date_event BETWEEN DATE_ADD(DATE(?),
			INTERVAL 6 HOUR) AND DATE_ADD(DATE(?),
			INTERVAL 6 HOUR)
	GROUP BY S.id
    `, line, startDate, endDate)

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

		pdl_calc = append(pdl_calc, segment)
	}

	defer db.Close()
	return pdl_calc
}

func calcUCLSegmentbyRange(line, startDate, endDate string) []models.
	TimeSegment {

	db := dbConn()
	var segment models.TimeSegment
	var ucl_calc []models.TimeSegment

	/*-------Manufactory Performance Lose-----------*/
	fmt.Printf("UCL")

	selDB, err_mpl := db.Query(`
	SELECT 
    	S.description, S.color, SUM(ExL.minutes)
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
	WHERE
		L.id = 3 AND ExL.id_line = ?
			AND B.id != 52
			AND ExL.date_event BETWEEN DATE_ADD(DATE(?),
			INTERVAL 6 HOUR) AND DATE_ADD(DATE(?),
			INTERVAL 6 HOUR)
	GROUP BY S.id
    `, line, startDate, endDate)

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

		ucl_calc = append(ucl_calc, segment)
	}
	defer db.Close()

	return ucl_calc
}

func calcGeneralSegmentbyRange(area, startDate, endDate string) []models.
	TimeSegment {
	db := dbConn()
	var segment models.TimeSegment
	var general_calc []models.TimeSegment

	/*-------Manufactory Performance Lose-----------*/
	fmt.Printf("General")

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
					INNER JOIN
				Area A ON A.id = Li.id_area
			where A.id = ?
			AND ExL.date_event BETWEEN DATE_ADD(DATE(?),
			INTERVAL 6 HOUR) AND DATE_ADD(DATE(?),
			INTERVAL 6 HOUR)
			AND L.id != 3
			group by S.id
			order by TotalMinutes desc
    `, area, startDate, endDate)

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

		general_calc = append(general_calc, segment)
	}
	defer db.Close()
	return general_calc
}

func calcByLineGeneralSegmentbyRange(line, startDate, endDate string) []models.
	TimeSegment {
	db := dbConn()
	var segment models.TimeSegment
	var general_calc []models.TimeSegment

	/*-------Manufactory Performance Lose-----------*/
	fmt.Printf("General")

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
        	INNER JOIN
		Area A ON A.id = Li.id_area
	where Li.id = ?
	AND ExL.date_event BETWEEN DATE_ADD(DATE(?),
	INTERVAL 6 HOUR) AND DATE_ADD(DATE(?),
	INTERVAL 6 HOUR)
	AND L.id != 3
	group by S.id
	order by TotalMinutes desc
    `, line, startDate, endDate)

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

		general_calc = append(general_calc, segment)
	}
	defer db.Close()
	return general_calc
}

//--------------------------------
func calcBreakdownsRange(line, startDate, endDate string) int {
	number := 0

	db := dbConn()

	selDB, err_mpl := db.Query(`
	SELECT 
    	COUNT(ExL.minutes) AS Breakdowns
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
			INNER JOIN
		Area A ON A.id = Li.id_area
	WHERE
			Li.id = ?
		AND S.id = 13
		AND ExL.date_event BETWEEN ? AND ?
	ORDER BY Breakdowns
    `, line, startDate, endDate)

	if err_mpl != nil {
		panic(err_mpl.Error())
	}

	for selDB.Next() {

		var count int

		err := selDB.Scan(&count)
		if err != nil {
			panic(err.Error())
		}

		number += count
	}
	defer db.Close()

	return number
}

func calcChangeoversRange(line, idCat, startDate, endDate string) int {
	number := 0

	db := dbConn()

	selDB, err_mpl := db.Query(`
	SELECT 
		COUNT(ExL.minutes) AS Changeovers
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
			INNER JOIN
		Area A ON A.id = Li.id_area
	WHERE
		E.id = ?
		AND Li.id = ?
		AND ExL.date_event BETWEEN ? AND ?
	ORDER BY Changeovers
    `, idCat, line, startDate, endDate)

	if err_mpl != nil {
		panic(err_mpl.Error())
	}

	for selDB.Next() {

		var count int

		err := selDB.Scan(&count)
		if err != nil {
			panic(err.Error())
		}

		number += count
	}
	defer db.Close()

	return number
}

func calcChangeoversTimeRange(line, idCat, startDate, endDate string) float32 {
	var number float32

	db := dbConn()

	selDB, err_mpl := db.Query(`
	SELECT 
		COALESCE(SUM(ExL.minutes), 0) AS Changeovers
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
			INNER JOIN
		Area A ON A.id = Li.id_area
	WHERE
		E.id = ?
		AND Li.id = ?
		AND ExL.date_event BETWEEN ? AND ?
	ORDER BY Changeovers
    `, idCat, line, startDate, endDate)

	if err_mpl != nil {
		panic(err_mpl.Error())
	}

	for selDB.Next() {

		var count float32

		err := selDB.Scan(&count)
		if err != nil {
			panic(err.Error())
		}

		number += count
	}
	defer db.Close()

	return number
}

func calcORbyRange(line, startDate, endDate string) float32 {
	var diffTotal float32
	var plannedTotal float32
	var producedTotal float32
	var or float32

	db := dbConn()

	selDB, err_mpl := db.Query(`
	SELECT 
    Pl.planned,
    Pl.produced
	FROM
		Planning Pl
			INNER JOIN
		Presentation Pr ON Pl.id_presentation = Pr.id
	WHERE
		Pl.id_line = ?
		AND DATE(Pl.date_planning) BETWEEN ? AND ?
    `, line, startDate, endDate)

	if err_mpl != nil {
		panic(err_mpl.Error())
	}

	for selDB.Next() {

		var planned float32
		var produced float32

		err := selDB.Scan(&planned, &produced)
		if err != nil {
			panic(err.Error())
		}

		plannedTotal += planned
		producedTotal += produced

		diffTotal += produced - planned
	}
	defer db.Close()

	or = (plannedTotal - diffTotal) / plannedTotal

	if math.IsNaN(float64(or)) {
		or = 0.0
	}

	fmt.Printf("\n[%s]\t[Planned:%f | Produced:%f]\t%OR:  %f\n", line, plannedTotal, producedTotal, diffTotal, or)

	return or
}

func CalcORbyTurn(w http.ResponseWriter, r *http.Request) {
	var diffTotal float32
	var plannedTotal float32
	var producedTotal float32
	var or float32

	line := r.URL.Query().Get("line")
	init := r.URL.Query().Get("init")
	turn := r.URL.Query().Get("turn")

	db := dbConn()

	selDB, err_mpl := db.Query(`
	SELECT 
	Pl.planned,
	Pl.produced
	FROM
		Planning Pl
			INNER JOIN
		Presentation Pr ON Pl.id_presentation = Pr.id
	WHERE
		Pl.id_line = ?
		AND Pl.date_planning = ?
		AND Pl.turn = ? `, line, init, turn)

	if err_mpl != nil {
		panic(err_mpl.Error())
	}

	for selDB.Next() {

		var planned float32
		var produced float32

		err := selDB.Scan(&planned, &produced)
		if err != nil {
			panic(err.Error())
		}

		plannedTotal += planned
		producedTotal += produced

		diffTotal += planned - produced
	}
	defer db.Close()

	or = (plannedTotal - diffTotal) / plannedTotal * 100

	if math.IsNaN(float64(or)) {
		or = 0.0
	}

	fmt.Printf("\n[%s]\t[Planned:%f | Produced:%f]\t%OR:  %f\n", line, plannedTotal, producedTotal, diffTotal, or)

	json.NewEncoder(w).Encode(or)
}
