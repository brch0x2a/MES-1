package controllers

import (
	//"bytes"
	//"fmt"
	"fmt"
	"log"
	"math"
	"net/http"
	"strconv"
	"time"

	//"image/jpeg"
	"encoding/json"

	"../models"
	_ "github.com/go-sql-driver/mysql"
)

/*
Obtener el calculo del oee, filtrando por linea, turno y fecha del dia a el que aplica
si por parametro el turno es igual a 4 se obtiene el oee de los tres turno que es
equivalente a tener el oee de ese dia
*/
func calcOEEProjection(line string, turn string, date_event string) float32 {

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
            Pl.turn = ? AND Pl.date_planning = ? AND Pl.id_line = ?  
    `
	hoursQ := 6
	minutesQ := 0

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
	if turn == "4" {
		minutesQ = 1800 - 1
		votSTM = `
        SELECT
            Pl.produced,
            Pl.nominal_speed,
            Pr.box_amount
        FROM
            Planning Pl
        INNER JOIN Presentation Pr ON
            Pl.id_presentation = Pr.id
        WHERE
            COALESCE(4, Pl.turn) = ?
            AND Pl.date_planning = ?
            AND Pl.id_line = ?
    `
	} else {
		minutesQ = (hoursQ+8)*60 - 1
	}

	db := dbConn()

	var vot_calc float32
	var mpl_calc float32
	var ot_calc float32
	var pdl_calc float32
	var lt_calc float32
	var oee float32

	//nId := r.URL.Query().Get("id")
	selDB, err := db.Query(votSTM, turn, date_event, line)

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

	/*-------Manufactory Performance Lose-----------*/

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
        AND ExL.date_event  BETWEEN DATE_ADD(DATE(?), INTERVAL `+strconv.Itoa(hoursQ)+` HOUR) 
        AND DATE_ADD(DATE(?), INTERVAL `+strconv.Itoa(minutesQ)+` MINUTE) 
    `, line, date_event, date_event)

	if err_mpl != nil {
		panic(err_mpl.Error())
	}

	p_mpl := models.MPL{}

	res_mpl := []models.MPL{}
	for selDB_mpl.Next() {

		var Minutes float32

		err = selDB_mpl.Scan(&Minutes)
		if err != nil {
			panic(err.Error())
		}
		p_mpl.Minutes = Minutes

		mpl_calc += Minutes
		//log.Printf("Minutes: %f\n", Minutes)

		res_mpl = append(res_mpl, p_mpl)
	}
	defer db.Close()

	dtime := getTimeDifference(line, turn, date_event).DMinutes

	mpl_calc += dtime

	fmt.Printf("\nDTimeDay: %f\n", dtime)

	//log.Printf("MPL Min: %f\n", mpl_calc)
	/*-------------OT--------------------------------*/
	ot_calc = vot_calc + mpl_calc/60

	/*-------Process Driven Losses-----------*/
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
        AND ExL.date_event  BETWEEN DATE_ADD(DATE(?), INTERVAL `+strconv.Itoa(hoursQ)+` HOUR) 
        AND DATE_ADD(DATE(?), INTERVAL `+strconv.Itoa(minutesQ)+` MINUTE) 
    `, line, date_event, date_event)

	if err_pdl != nil {
		panic(err_pdl.Error())
	}

	p_pdl := models.PDL{}

	res_pdl := []models.PDL{}
	for selDB_pdl.Next() {

		var Minutes float32

		err = selDB_pdl.Scan(&Minutes)
		if err != nil {
			panic(err.Error())
		}
		p_pdl.Minutes = Minutes

		pdl_calc += Minutes
		//log.Printf("Minutes: %f\n", Minutes)

		res_pdl = append(res_pdl, p_pdl)
	}
	defer db.Close()

	lt_calc = ot_calc + pdl_calc/60

	oee = vot_calc / lt_calc * 100
	/*----------------------------------------------
	log.Printf("OT: %f\n", ot_calc)
	log.Printf("PDL: %f\n", pdl_calc)
	log.Printf("LT: %f\n", lt_calc)

	log.Printf("\nOEE! %f", oee)
	-----------------------------------*/

	if math.IsNaN(float64(oee)) {
		//log.Println("Null OEE")
		oee = -1
	}

	return oee
}

/*
____________________________________________________________________________________________________________________________________________________________________
Refactory Section
____________________________________________________________________________________________________________________________________________________________________
*/
//Se asume que se trataa como sumatoria cuando son rangos de varios turnos
func getPlannedDataByDay(line, turn, date_event string) models.Planning_holder {
	planeedData := models.Planning_holder{}
	var stm string
	var PlannedSum int
	var ProducedSum int
	var Nominal_speedSum float32
	var Box_amountSum int

	if turn == "4" {
		stm = `SELECT
		Pl.date_planning,
		Pl.turn, 
		Pl.id_line,
		Pr.name,
		Pl.version,
		Pl.planned,
		Pl.produced,
		Pl.nominal_speed,
		Pr.box_amount
		FROM
			Planning Pl
		INNER JOIN Presentation Pr ON
			Pl.id_presentation = Pr.id
		WHERE
		COALESCE(4, Pl.turn) = ?	
		AND Pl.date_planning = ?
		AND Pl.id_line = ?`
	} else {
		stm = `SELECT
		Pl.date_planning,
		Pl.turn, 
		Pl.id_line,
		Pr.name,
		Pl.version,
		Pl.planned,
		Pl.produced,
		Pl.nominal_speed,
		Pr.box_amount
		FROM
			Planning Pl
		INNER JOIN Presentation Pr ON
			Pl.id_presentation = Pr.id
		WHERE
			Pl.turn = ? AND Pl.date_planning = ? AND Pl.id_line = ?`
	}

	db := dbConn()

	selDB, err := db.Query(stm, turn, date_event, line)

	if err != nil {
		panic(err.Error())
	}

	for selDB.Next() {

		var Date_planning string
		var Turn int
		var Line string
		var Presentation string
		var Version int
		var Planned int
		var Produced int
		var Nominal_speed float32
		var Box_amount int

		err = selDB.Scan(&Date_planning, &Turn, &Line, &Presentation,
			&Version, &Planned, &Produced, &Nominal_speed, &Box_amount)

		if err != nil {
			panic(err.Error())
		}

		planeedData.Date_planning = Date_planning
		planeedData.Turn = Turn
		planeedData.Line = Line
		planeedData.Presentation = Presentation
		planeedData.Version = Version

		PlannedSum += Planned
		ProducedSum += Produced

		Nominal_speedSum += Nominal_speed
		Box_amountSum += Box_amount

	}

	planeedData.Planned = PlannedSum
	planeedData.Produced = ProducedSum
	planeedData.Nominal_speed = Nominal_speedSum
	planeedData.Box_amount = Box_amountSum
	defer db.Close()
	return planeedData
}

func getMPLByDay(line, turn, date string) float32 { //MPL = Manufactory Perfomance Losses

	db := dbConn()
	var Minutes float32
	hoursQ := 6
	minutesQ := 0
	var mplCalc float32

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
	if turn == "4" {
		minutesQ = 1800 - 1
	} else {
		minutesQ = (hoursQ+8)*60 - 1
	}

	/*-------Manufactory Performance Lose-----------*/

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
        AND ExL.date_event  BETWEEN DATE_ADD(DATE(?), INTERVAL `+strconv.Itoa(hoursQ)+` HOUR) 
        AND DATE_ADD(DATE(?), INTERVAL `+strconv.Itoa(minutesQ)+` MINUTE) 
    `, line, date, date)

	if err_mpl != nil {
		panic(err_mpl.Error())
	}

	for selDB_mpl.Next() {

		err := selDB_mpl.Scan(&Minutes)
		if err != nil {
			panic(err.Error())
		}

		mplCalc += Minutes
		//log.Printf("Minutes: %f\n", Minutes)

	}
	defer db.Close()

	return mplCalc
}

func getPDLByDay(line, turn, date string) float32 { //MPL = Manufactory Perfomance Losses

	db := dbConn()
	var Minutes float32
	hoursQ := 6
	minutesQ := 0
	var mplCalc float32

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
	if turn == "4" {
		minutesQ = 1800 - 1
	} else {
		minutesQ = (hoursQ+8)*60 - 1
	}

	/*-------------Process Driven Losses-----------*/
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
        L.id = 2  
        AND ExL.id_line = ?
        AND ExL.date_event  BETWEEN DATE_ADD(DATE(?), INTERVAL `+strconv.Itoa(hoursQ)+` HOUR) 
        AND DATE_ADD(DATE(?), INTERVAL `+strconv.Itoa(minutesQ)+` MINUTE) 
    `, line, date, date)

	if err_mpl != nil {
		panic(err_mpl.Error())
	}

	for selDB_mpl.Next() {

		err := selDB_mpl.Scan(&Minutes)
		if err != nil {
			panic(err.Error())
		}

		mplCalc += Minutes
		//log.Printf("Minutes: %f\n", Minutes)

	}
	defer db.Close()

	return mplCalc
}

func getPlantStoppageByDay(line, turn, date string) float32 { //

	/*-------Plant Stoppage-----------*/
	db := dbConn()
	var Minutes float32
	hoursQ := 6
	minutesQ := 0
	var mplCalc float32

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
	if turn == "4" {
		minutesQ = 1800 - 1
	} else {
		minutesQ = (hoursQ+8)*60 - 1
	}

	/*-------------Process Driven Losses-----------*/
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
        L.id = 3  
        AND ExL.id_line = ?
        AND ExL.date_event  BETWEEN DATE_ADD(DATE(?), INTERVAL `+strconv.Itoa(hoursQ)+` HOUR) 
        AND DATE_ADD(DATE(?), INTERVAL `+strconv.Itoa(minutesQ)+` MINUTE) 
    `, line, date, date)

	if err_mpl != nil {
		panic(err_mpl.Error())
	}

	for selDB_mpl.Next() {

		err := selDB_mpl.Scan(&Minutes)
		if err != nil {
			panic(err.Error())
		}

		mplCalc += Minutes
		//log.Printf("Minutes: %f\n", Minutes)

	}
	defer db.Close()

	return mplCalc
}

/*
____________________________________________________________________________________________________________________________________________________________________
______________________________________________________________________________________________________________________________________________________________________
*/

/*
Obtener el calculo del oee, filtrando por linea, turno y fecha del dia a el que aplica
si por parametro el turno es igual a 4 se obtiene el oee de los tres turno que es
equivalente a tener el oee de ese dia
*/
func calcOEE(line string, turn string, date_event string) float32 {

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
            Pl.turn = ? AND Pl.date_planning = ? AND Pl.id_line = ?  
    `
	hoursQ := 6
	minutesQ := 0

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
	if turn == "4" {
		minutesQ = 1800 - 1
		votSTM = `
        SELECT
            Pl.produced,
            Pl.nominal_speed,
            Pr.box_amount
        FROM
            Planning Pl
        INNER JOIN Presentation Pr ON
            Pl.id_presentation = Pr.id
        WHERE
            COALESCE(4, Pl.turn) = ?
            AND Pl.date_planning = ?
            AND Pl.id_line = ?
    `
	} else {
		minutesQ = (hoursQ+8)*60 - 1
	}

	db := dbConn()

	var vot_calc float32
	var mpl_calc float32
	var ot_calc float32
	var pdl_calc float32
	var lt_calc float32
	var oee float32

	var speedLoss float32
	var badBoxes float32
	var goodBoxes float32
	var standar float32

	var qaDetectTime float32
	// var qualityDefect float32

	//nId := r.URL.Query().Get("id")

	selDB, err := db.Query(votSTM, turn, date_event, line)

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

		goodBoxes += GoodVolume
		standar += NominalSpeed / Box_amount * 60

		badBoxes = 8*standar - goodBoxes

		qaDetectTime += badBoxes / NominalSpeed
		//log.Printf("VOT: %.6f\n", calc)

		calc = p.Get() / 60
		vot_calc += calc

		//fmt.Printf("vot: %f\n\n", vot_calc)

		res = append(res, calc)
	}
	defer db.Close()

	/*-------Manufactory Performance Lose-----------*/

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
        AND ExL.date_event  BETWEEN DATE_ADD(DATE(?), INTERVAL `+strconv.Itoa(hoursQ)+` HOUR) 
        AND DATE_ADD(DATE(?), INTERVAL `+strconv.Itoa(minutesQ)+` MINUTE) 
    `, line, date_event, date_event)

	if err_mpl != nil {
		panic(err_mpl.Error())
	}

	p_mpl := models.MPL{}

	res_mpl := []models.MPL{}
	for selDB_mpl.Next() {

		var Minutes float32

		err = selDB_mpl.Scan(&Minutes)
		if err != nil {
			panic(err.Error())
		}
		p_mpl.Minutes = Minutes

		mpl_calc += Minutes
		//log.Printf("Minutes: %f\n", Minutes)

		res_mpl = append(res_mpl, p_mpl)
	}
	defer db.Close()

	//log.Printf("MPL Min: %f\n", mpl_calc)

	mpl_calc += speedLoss + qaDetectTime

	speedLoss = lt_calc - (mpl_calc+pdl_calc)/60 - vot_calc

	// cajasMalas  = 8 horas * Standar (Cjs/Hora)-Cajas buenas
	// standar  = nomalSpeed / unidBox * 60min
	//Quality Defect Time =(Volume produced NC ) / Nominal Speed

	fmt.Println("\n\n\t--------------------\n")
	fmt.Println("\n\n\t------TURN:%s---LINE:%s-----------\n", turn, line)

	fmt.Println("SpeedLoss-->", speedLoss)
	fmt.Println("badBoxes-->", badBoxes)
	fmt.Println("qualityDefect-->", qaDetectTime)

	/*-------------OT--------------------------------*/
	ot_calc = vot_calc + mpl_calc/60

	/*-------Process Driven Losses-----------*/
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
        AND ExL.date_event  BETWEEN DATE_ADD(DATE(?), INTERVAL `+strconv.Itoa(hoursQ)+` HOUR) 
        AND DATE_ADD(DATE(?), INTERVAL `+strconv.Itoa(minutesQ)+` MINUTE) 
    `, line, date_event, date_event)

	if err_pdl != nil {
		panic(err_pdl.Error())
	}

	p_pdl := models.PDL{}

	res_pdl := []models.PDL{}
	for selDB_pdl.Next() {

		var Minutes float32

		err = selDB_pdl.Scan(&Minutes)
		if err != nil {
			panic(err.Error())
		}
		p_pdl.Minutes = Minutes

		pdl_calc += Minutes
		//log.Printf("Minutes: %f\n", Minutes)

		res_pdl = append(res_pdl, p_pdl)
	}
	defer db.Close()

	lt_calc = ot_calc + pdl_calc/60

	oee = vot_calc / lt_calc * 100
	/*----------------------------------------------
	log.Printf("OT: %f\n", ot_calc)
	log.Printf("PDL: %f\n", pdl_calc)
	log.Printf("LT: %f\n", lt_calc)

	log.Printf("\nOEE! %f", oee)
	-----------------------------------*/

	if math.IsNaN(float64(oee)) {
		//log.Println("Null OEE")
		oee = -1
	}

	return oee
}

func WeekRange(year, week int) (start, end time.Time) {
	start = WeekStart(year, week)
	end = start.AddDate(0, 0, 6)
	return
}

func WeekStart(year, week int) time.Time {
	// Start from the middle of the year:
	t := time.Date(year, 7, 1, 0, 0, 0, 0, time.UTC)

	// Roll back to Monday:
	if wd := t.Weekday(); wd == time.Sunday {
		t = t.AddDate(0, 0, -6)
	} else {
		t = t.AddDate(0, 0, -int(wd)+1)
	}

	// Difference in weeks:
	_, w := t.ISOWeek()
	t = t.AddDate(0, 0, (week-w)*7)

	return t
}

func CalcOEEbyRange(line, startDate, endDate string) float32 {

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

	var vot_calc float32
	var mpl_calc float32
	var ot_calc float32
	var pdl_calc float32
	var lt_calc float32
	var oee float32

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

	/*-------Manufactory Performance Lose-----------*/

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
		INTERVAL 6 HOUR)
    `, line, startDate, endDate)

	if err_mpl != nil {
		panic(err_mpl.Error())
	}

	p_mpl := models.MPL{}

	res_mpl := []models.MPL{}
	for selDB_mpl.Next() {

		var Minutes float32

		err = selDB_mpl.Scan(&Minutes)
		if err != nil {
			panic(err.Error())
		}
		p_mpl.Minutes = Minutes

		mpl_calc += Minutes
		//log.Printf("Minutes: %f\n", Minutes)

		res_mpl = append(res_mpl, p_mpl)
	}
	defer db.Close()

	//log.Printf("MPL Min: %f\n", mpl_calc)
	/*-------------OT--------------------------------*/
	ot_calc = vot_calc + mpl_calc/60

	/*-------Process Driven Losses-----------*/
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

	p_pdl := models.PDL{}

	res_pdl := []models.PDL{}
	for selDB_pdl.Next() {

		var Minutes float32

		err = selDB_pdl.Scan(&Minutes)
		if err != nil {
			panic(err.Error())
		}
		p_pdl.Minutes = Minutes

		pdl_calc += Minutes
		//log.Printf("Minutes: %f\n", Minutes)

		res_pdl = append(res_pdl, p_pdl)
	}
	defer db.Close()

	lt_calc = ot_calc + pdl_calc/60

	oee = vot_calc / lt_calc * 100
	/*----------------------------------------------
	log.Printf("OT: %f\n", ot_calc)
	log.Printf("PDL: %f\n", pdl_calc)
	log.Printf("LT: %f\n", lt_calc)

	log.Printf("\nOEE! %f", oee)
	-----------------------------------*/
	log.Printf("\nBefore>>OEE! %f\n\n", oee)

	if math.IsNaN(float64(oee)) {
		//log.Println("Null OEE")
		oee = -1
	}

	return oee

}

func CalcOEEMetaProjectionbyRange(line, startDate, endDate string) (float32, float32) {

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

	var vot_calc float32
	var mpl_calc float32
	var ot_calc float32
	var pdl_calc float32
	var lt_calc float32
	var oee float32

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

	/*-------Manufactory Performance Lose-----------*/

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
		INTERVAL 6 HOUR)
    `, line, startDate, endDate)

	if err_mpl != nil {
		panic(err_mpl.Error())
	}

	p_mpl := models.MPL{}

	res_mpl := []models.MPL{}
	for selDB_mpl.Next() {

		var Minutes float32

		err = selDB_mpl.Scan(&Minutes)
		if err != nil {
			panic(err.Error())
		}
		p_mpl.Minutes = Minutes

		mpl_calc += Minutes
		//log.Printf("Minutes: %f\n", Minutes)

		res_mpl = append(res_mpl, p_mpl)
	}
	defer db.Close()

	//log.Printf("MPL Min: %f\n", mpl_calc)
	/*-------------OT--------------------------------*/
	ot_calc = vot_calc + mpl_calc/60

	/*-------Process Driven Losses-----------*/
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

	p_pdl := models.PDL{}

	res_pdl := []models.PDL{}
	for selDB_pdl.Next() {

		var Minutes float32

		err = selDB_pdl.Scan(&Minutes)
		if err != nil {
			panic(err.Error())
		}
		p_pdl.Minutes = Minutes

		pdl_calc += Minutes
		//log.Printf("Minutes: %f\n", Minutes)

		res_pdl = append(res_pdl, p_pdl)
	}
	defer db.Close()

	layout := "2006-01-02"

	tStart, erra := time.Parse(layout, startDate)
	tEnd, errb := time.Parse(layout, endDate)

	if erra != nil {
		fmt.Println(erra)
	}

	if errb != nil {
		fmt.Println(errb)
	}

	fmt.Println(tStart)

	fmt.Println(tEnd)

	diff := tEnd.Sub(tStart)
	days := diff.Hours() / 24

	// fmt.Println("Days diference: %d", days)

	// startDate = day.String()
	// endDate = day.AddDate(0, 0, 6).String()
	var dtime float32
	var allocatedLosses float32

	dtime = 0
	allocatedLosses = 0

	day := tStart

	for index := 0; index < int(days); index++ {

		dtime += getTimeDifference(line, "4", day.String()).DMinutes

		//___________________________________________________________

		//fmt.Printf("\n\nday %s\t\toee: %f", day.String(), calcOEE(line, week, day.String()))
		day = day.AddDate(0, 0, 1)
		// fmt.Printf("\nDtime %f\n", dtime)

	}

	pdl_calc += dtime

	lt_calc = ot_calc + pdl_calc/60

	oee = vot_calc / lt_calc * 100

	// log.Printf("OT: %f\n", ot_calc)
	// log.Printf("PDL: %f\n", pdl_calc)
	// log.Printf("LT: %f\n", lt_calc)
	allocatedLosses = pdl_calc + mpl_calc

	fmt.Printf("Dtime: %f\t Allocated: %f",
		dtime, allocatedLosses)

	// log.Printf("\nOEE! %f", oee)

	// log.Printf("\nBefore>>OEE! %f\n\n", oee)

	if math.IsNaN(float64(oee)) {
		//log.Println("Null OEE")
		oee = -1
	}
	// fmt.Println("\n\n")

	return dtime, allocatedLosses

}

func CalcOEEProjectionbyRange(line, startDate, endDate string, deviation float32) float32 {

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

	var vot_calc float32
	var mpl_calc float32
	var ot_calc float32
	var pdl_calc float32
	var lt_calc float32
	var oee float32

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

	/*-------Manufactory Performance Lose-----------*/

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
		INTERVAL 6 HOUR)
    `, line, startDate, endDate)

	if err_mpl != nil {
		panic(err_mpl.Error())
	}

	p_mpl := models.MPL{}

	res_mpl := []models.MPL{}
	for selDB_mpl.Next() {

		var Minutes float32

		err = selDB_mpl.Scan(&Minutes)
		if err != nil {
			panic(err.Error())
		}
		p_mpl.Minutes = Minutes

		mpl_calc += Minutes
		//log.Printf("Minutes: %f\n", Minutes)

		res_mpl = append(res_mpl, p_mpl)
	}
	defer db.Close()

	//log.Printf("MPL Min: %f\n", mpl_calc)
	/*-------------OT--------------------------------*/
	ot_calc = vot_calc + mpl_calc/60

	/*-------Process Driven Losses-----------*/
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

	p_pdl := models.PDL{}

	res_pdl := []models.PDL{}
	for selDB_pdl.Next() {

		var Minutes float32

		err = selDB_pdl.Scan(&Minutes)
		if err != nil {
			panic(err.Error())
		}
		p_pdl.Minutes = Minutes

		pdl_calc += Minutes
		//log.Printf("Minutes: %f\n", Minutes)

		res_pdl = append(res_pdl, p_pdl)
	}
	defer db.Close()

	layout := "2006-01-02"

	tStart, erra := time.Parse(layout, startDate)
	tEnd, errb := time.Parse(layout, endDate)

	if erra != nil {
		fmt.Println(erra)
	}

	if errb != nil {
		fmt.Println(errb)
	}

	fmt.Println(tStart)

	fmt.Println(tEnd)

	diff := tEnd.Sub(tStart)
	days := diff.Hours() / 24

	// fmt.Println("Days diference: %d", days)

	// startDate = day.String()
	// endDate = day.AddDate(0, 0, 6).String()
	var dtime float32
	var allocatedLosses float32

	dtime = 0
	allocatedLosses = 0

	day := tStart

	for index := 0; index < int(days); index++ {

		dtime += getTimeDifference(line, "4", day.String()).DMinutes

		//___________________________________________________________

		//fmt.Printf("\n\nday %s\t\toee: %f", day.String(), calcOEE(line, week, day.String()))
		day = day.AddDate(0, 0, 1)
		// fmt.Printf("\nDtime %f\n", dtime)

	}

	pdl_calc += dtime * deviation

	lt_calc = ot_calc + pdl_calc/60

	oee = vot_calc / lt_calc * 100

	// log.Printf("OT: %f\n", ot_calc)
	// log.Printf("PDL: %f\n", pdl_calc)
	// log.Printf("LT: %f\n", lt_calc)
	allocatedLosses = pdl_calc + mpl_calc

	fmt.Printf("%f\t%f",
		dtime, allocatedLosses)

	// log.Printf("\nOEE! %f", oee)

	// log.Printf("\nBefore>>OEE! %f\n\n", oee)

	if math.IsNaN(float64(oee)) {
		//log.Println("Null OEE")
		oee = -1
	}
	// fmt.Println("\n\n")

	return oee

}

func GetOEE(w http.ResponseWriter, r *http.Request) {
	log.Printf("\n\n\n---------Request-OEE------------\n")

	turn := r.URL.Query().Get("turn")
	line := r.URL.Query().Get("line")
	date_event := r.URL.Query().Get("date")

	oee := calcOEE(line, turn, date_event)

	json.NewEncoder(w).Encode(oee)
}

func getOEEWeekxLine(line, week, year, lineName string) models.WeekOEE {
	//log.Printf("\n\n\n---------Request-----OEE-WEEK-------\n")

	weekNum, _ := strconv.Atoi(week)
	yearNum, _ := strconv.Atoi(year)

	day := WeekStart(yearNum, weekNum)

	weekOEE := models.WeekOEE{}
	todayOEE := models.DayOEE{}
	startDate := ""
	endDate := ""

	//fmt.Printf("\nWeek %d", week)

	startDate = day.String()
	endDate = day.AddDate(0, 0, 6).String()

	for index := 0; index < 7; index++ {

		todayOEE.Day = day.String()
		todayOEE.OEE = calcOEE(line, "4", day.String())

		weekOEE.DayInfo = append(weekOEE.DayInfo, todayOEE)

		//___________________________________________________________

		fmt.Printf("\n\nday %s\t\toee: %f", day.String(), calcOEE(line, week, day.String()))
		day = day.AddDate(0, 0, 1)

	}
	oeeWeek := CalcOEEbyRange(line, startDate, endDate)

	fmt.Printf("\n----------------------------------------------------------\nLine\t|Week\t|StartDate\t\t\t\t|EndDate\t\t\t|OEE\n")
	fmt.Printf("%s\t|%s\t|%s\t\t|%s\t\t|%f\n", lineName, week, startDate, endDate, oeeWeek)
	fmt.Printf("----------------------------------------------------------\n")

	weekOEE.Line = lineName
	weekOEE.Total = oeeWeek

	return weekOEE
}

func GetOEEWeekxLineAPI(w http.ResponseWriter, r *http.Request) {

	line := r.URL.Query().Get("line")
	week := r.URL.Query().Get("week")
	year := r.URL.Query().Get("year")

	weekOEE := models.WeekOEE{}

	lineName := GetLineObjectE(line).Name

	weekOEE = getOEEWeekxLine(line, week, year, lineName)

	json.NewEncoder(w).Encode(weekOEE)
}

func getRelativeOEEDay(line, day string) []models.RelativeOEE {

	relative := models.RelativeOEE{}
	L := []models.RelativeOEE{}
	var turn string

	for index := 1; index <= 4; index++ {

		turn = strconv.Itoa(index)

		relative.OEE = calcOEE(line, turn, day)
		relative.Dtime = getTimeDifference(line, turn, day).DMinutes
		relative.StopTime = getMPLByDay(line, turn, day) + getPDLByDay(line, turn, day)

		L = append(L, relative)

	}
	return L
}

func getRelativeOEEDayProject(line, day string) models.RelativeOEE {
	relative := models.RelativeOEE{}

	return relative
}

func GetRelativeOEEDayData(w http.ResponseWriter, r *http.Request) {

	line := r.URL.Query().Get("line")
	day := r.URL.Query().Get("day")

	Relative := getRelativeOEEDay(line, day)

	// oeeProjection := calcOEEProjection(line, "4", day)

	// fmt.Printf("[%s]-->oeeProjection: %f", day, oeeProjection)

	json.NewEncoder(w).Encode(Relative)
}

func GetOEEWeekLineTotal(w http.ResponseWriter, r *http.Request) {

	line := r.URL.Query().Get("line")
	week := r.URL.Query().Get("week")
	year := r.URL.Query().Get("year")

	weekNum, _ := strconv.Atoi(week)
	yearNum, _ := strconv.Atoi(year)

	day := WeekStart(yearNum, weekNum)

	startDate := ""
	endDate := ""

	//fmt.Printf("\nWeek %d", week)

	startDate = day.String()
	endDate = day.AddDate(0, 0, 6).String()

	oeeWeek := CalcOEEbyRange(line, startDate, endDate)

	json.NewEncoder(w).Encode(oeeWeek)
}

func GetOEEWeek(w http.ResponseWriter, r *http.Request) {

	weekLine := models.WeekOEE{}
	area := []models.WeekOEE{}
	nId := r.URL.Query().Get("area")
	//line := r.URL.Query().Get("line")
	week := r.URL.Query().Get("week")
	year := r.URL.Query().Get("year")

	lines := getLinesByArea(nId)

	for i := 0; i < len(lines); i++ {
		//fmt.Printf("\n\n---- %s ----\n", lines[i].Name)
		weekLine = getOEEWeekxLine(strconv.Itoa(lines[i].Id), week, year, lines[i].Name)
		area = append(area, weekLine)
	}

	json.NewEncoder(w).Encode(area)
}

/*_______________________________________________________*/
func FilterOEE(w http.ResponseWriter, r *http.Request) {
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
	tmpl.ExecuteTemplate(w, "FilterOEE", nil)
}

func OEEDemo00(w http.ResponseWriter, r *http.Request) {
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
	tmpl.ExecuteTemplate(w, "OEEDemo00", nil)
}
func OEEDemo01(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "OEEDemo01", nil)
}
func OEEDemo02(w http.ResponseWriter, r *http.Request) {
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
	tmpl.ExecuteTemplate(w, "OEEDemo02", nil)
}

/*______________________Validation_________________________________*/
func ValidationOEE(w http.ResponseWriter, r *http.Request) {
	db := dbConn()

	selDB, err := db.Query(`
            SELECT
            E.id,
            LTC.description,
            S.description,
            B.description,
            E.description,
            S.color
        FROM Event 
            E
        INNER JOIN Branch B ON
            E.id_branch = B.id
        INNER JOIN Sub_classification S ON
            S.id = B.id_sub_classification
        INNER JOIN LineTimeClassification LTC ON
			LTC.id = S.id_LTC
		ORDER BY S.id
    `)
	if err != nil {
		panic(err.Error())
	}

	p := models.Event_holder{}

	res := []models.Event_holder{}
	for selDB.Next() {
		var Id int
		var LTC string //Line time classification
		var Sub string //Sub classification
		var Branch string
		var Description string
		var Color string

		err = selDB.Scan(&Id, &LTC, &Sub,
			&Branch, &Description, &Color)
		if err != nil {
			panic(err.Error())
		}
		p.Id = Id
		p.LTC = LTC
		p.Sub = Sub
		p.Branch = Branch
		p.Description = Description
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

	data := map[string]interface{}{
		"User":  user,
		"Event": res,
	}
	defer db.Close()

	tmpl.ExecuteTemplate(w, "ValidationOEE", data)
}

func EventsTemplate(w http.ResponseWriter, r *http.Request) {

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

	tmpl.ExecuteTemplate(w, "EventsTemplate", user)
}

func ConsolidateOEE(w http.ResponseWriter, r *http.Request) {
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

	tmpl.ExecuteTemplate(w, "ConsolidateOEE", user)
}

func ConsolidateOEEWeek(w http.ResponseWriter, r *http.Request) {
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

	tmpl.ExecuteTemplate(w, "ConsolidateOEEWeek", user)
}

//-----------------------------------------------------------

func CalcOEEAreabyRange(area, startDate, endDate string) float32 {

	votSTM := `
        SELECT
        Pl.produced,
        Pl.nominal_speed,
        Pr.box_amount
        FROM
            Planning Pl
        INNER JOIN Presentation Pr ON
            Pl.id_presentation = Pr.id
            INNER JOIN
		Line L ON L.id = Pl.id_line
			INNER JOIN
		Area A ON A.Id = L.id_area
		WHERE
			A.Id = ?
			AND Pl.date_planning BETWEEN ? AND ?
    `

	db := dbConn()

	var vot_calc float32
	var mpl_calc float32
	var ot_calc float32
	var pdl_calc float32
	var lt_calc float32
	var oee float32

	//nId := r.URL.Query().Get("id")
	selDB, err := db.Query(votSTM, area, startDate, endDate)

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

	/*-------Manufactory Performance Lose-----------*/

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
        INNER JOIN
    Line Li ON Li.id = ExL.id_line
        INNER JOIN
    Area A ON A.Id = Li.id_area
WHERE
    L.id = 1 AND A.Id = ?
        AND ExL.date_event   BETWEEN DATE_ADD(DATE(?),
		INTERVAL 6 HOUR) AND DATE_ADD(DATE(?),
		INTERVAL 6 HOUR)
    `, area, startDate, endDate)

	if err_mpl != nil {
		panic(err_mpl.Error())
	}

	p_mpl := models.MPL{}

	res_mpl := []models.MPL{}
	for selDB_mpl.Next() {

		var Minutes float32

		err = selDB_mpl.Scan(&Minutes)
		if err != nil {
			panic(err.Error())
		}
		p_mpl.Minutes = Minutes

		mpl_calc += Minutes
		//log.Printf("Minutes: %f\n", Minutes)

		res_mpl = append(res_mpl, p_mpl)
	}
	defer db.Close()

	//log.Printf("MPL Min: %f\n", mpl_calc)
	/*-------------OT--------------------------------*/
	ot_calc = vot_calc + mpl_calc/60

	/*-------Process Driven Losses-----------*/
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
        INNER JOIN
    Line Li ON Li.id = ExL.id_line
        INNER JOIN
    Area A ON A.Id = Li.id_area
WHERE
    L.id = 2 AND A.Id = ?
        AND ExL.date_event  BETWEEN DATE_ADD(DATE(?),
		INTERVAL 6 HOUR) AND DATE_ADD(DATE(?),
		INTERVAL 6 HOUR)
    `, area, startDate, endDate)

	if err_pdl != nil {
		panic(err_pdl.Error())
	}

	p_pdl := models.PDL{}

	res_pdl := []models.PDL{}
	for selDB_pdl.Next() {

		var Minutes float32

		err = selDB_pdl.Scan(&Minutes)
		if err != nil {
			panic(err.Error())
		}
		p_pdl.Minutes = Minutes

		pdl_calc += Minutes
		//log.Printf("Minutes: %f\n", Minutes)

		res_pdl = append(res_pdl, p_pdl)
	}
	defer db.Close()

	lt_calc = ot_calc + pdl_calc/60

	oee = vot_calc / lt_calc * 100
	/*----------------------------------------------
	log.Printf("OT: %f\n", ot_calc)
	log.Printf("PDL: %f\n", pdl_calc)
	log.Printf("LT: %f\n", lt_calc)

	log.Printf("\nOEE! %f", oee)
	-----------------------------------*/
	log.Printf("\nBefore>>OEE! %f\n\n", oee)

	if math.IsNaN(float64(oee)) {
		//log.Println("Null OEE")
		oee = -1
	}

	return oee

}

func CalcOEELineProjectionbyRange(line, start, end string) models.OEEProjection {
	Projection := models.OEEProjection{}

	Deviation := []float32{0.0, 1.0, 0.15, 0.1}

	lineName := GetLineObjectE(line)

	Projection.Line = lineName.Name
	Projection.CurrentCase = CalcOEEProjectionbyRange(line, start, end, Deviation[0])
	// Projection.WorseCase = CalcOEEProjectionbyRange(line, start, end, Deviation[1])
	// Projection.AverageCase = CalcOEEProjectionbyRange(line, start, end, Deviation[2])
	// Projection.BestCase = CalcOEEProjectionbyRange(line, start, end, Deviation[3])

	return Projection
}

func GetOEEProjectionbyRange(w http.ResponseWriter, r *http.Request) {

	start := r.URL.Query().Get("start")
	end := r.URL.Query().Get("end")

	fmt.Printf("init: %s\tend: %s", start, end)

	areas := GetAreaList()
	Plant := []models.OEEProjection{}

	for a := 0; a < len(areas); a++ {

		areaID := strconv.Itoa(areas[a].Id)

		lines := getLinesByArea(areaID)

		for i := 0; i < len(lines); i++ {
			fmt.Printf("\n\n---- %s ----\n", lines[i].Name)
			lineId := strconv.Itoa(lines[i].Id)

			Projection := CalcOEELineProjectionbyRange(lineId, start, end)

			Plant = append(Plant, Projection)
		}

	}

	json.NewEncoder(w).Encode(Plant)
}

func GetOEEMetaProjectionbyRange(w http.ResponseWriter, r *http.Request) {

	start := r.URL.Query().Get("start")
	end := r.URL.Query().Get("end")

	areas := GetAreaList()

	fmt.Println("--\tGetOEEMetaProjectionbyRange\t--")

	fmt.Printf("\ninit: %s\tend: %s\n", start, end)

	oeeSaturation := []models.OEESaturation{}
	saturation := models.OEESaturation{}

	for a := 0; a < len(areas); a++ {

		areaID := strconv.Itoa(areas[a].Id)

		lines := getLinesByArea(areaID)

		for i := 0; i < len(lines); i++ {
			fmt.Printf("\n\n---- %s ----\n", lines[i].Name)

			lineId := strconv.Itoa(lines[i].Id)

			// diff, allocated := CalcOEEMetaProjectionbyRange(lineId, start, end)

			diff := CalcDtimeProjectionbyRange(lineId, start, end)
			allocated := CalcAllocatedLossesProjectionbyRange(lineId, start, end)

			fmt.Printf("%f\t%f",
				diff, allocated)

			saturation.Line = lines[i].Name
			saturation.Diff = diff
			saturation.Allocated = allocated

			oeeSaturation = append(oeeSaturation, saturation)
		}

	}

	json.NewEncoder(w).Encode(oeeSaturation)
}

func OEEYTDProjection(w http.ResponseWriter, r *http.Request) {

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
	tmpl.ExecuteTemplate(w, "OEEYTDProjection", user)
}

func CalcDtimeProjectionbyRange(line, startDate, endDate string) float32 {

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

	var vot_calc float32
	var mpl_calc float32
	var ot_calc float32
	var pdl_calc float32
	var lt_calc float32
	var oee float32

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

	/*-------Manufactory Performance Lose-----------*/

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
		INTERVAL 6 HOUR)
    `, line, startDate, endDate)

	if err_mpl != nil {
		panic(err_mpl.Error())
	}

	p_mpl := models.MPL{}

	res_mpl := []models.MPL{}
	for selDB_mpl.Next() {

		var Minutes float32

		err = selDB_mpl.Scan(&Minutes)
		if err != nil {
			panic(err.Error())
		}
		p_mpl.Minutes = Minutes

		mpl_calc += Minutes
		//log.Printf("Minutes: %f\n", Minutes)

		res_mpl = append(res_mpl, p_mpl)
	}
	defer db.Close()

	//log.Printf("MPL Min: %f\n", mpl_calc)
	/*-------------OT--------------------------------*/
	ot_calc = vot_calc + mpl_calc/60

	/*-------Process Driven Losses-----------*/
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

	p_pdl := models.PDL{}

	res_pdl := []models.PDL{}
	for selDB_pdl.Next() {

		var Minutes float32

		err = selDB_pdl.Scan(&Minutes)
		if err != nil {
			panic(err.Error())
		}
		p_pdl.Minutes = Minutes

		pdl_calc += Minutes
		//log.Printf("Minutes: %f\n", Minutes)

		res_pdl = append(res_pdl, p_pdl)
	}
	defer db.Close()

	layout := "2006-01-02"

	tStart, erra := time.Parse(layout, startDate)
	tEnd, errb := time.Parse(layout, endDate)

	if erra != nil {
		fmt.Println(erra)
	}

	if errb != nil {
		fmt.Println(errb)
	}

	fmt.Println(tStart)

	fmt.Println(tEnd)

	diff := tEnd.Sub(tStart)
	days := diff.Hours() / 24

	// fmt.Println("Days diference: %d", days)

	// startDate = day.String()
	// endDate = day.AddDate(0, 0, 6).String()
	var dtime float32
	// var allocatedLosses float32

	// dtime = 0
	// allocatedLosses = 0

	day := tStart

	for index := 0; index < int(days); index++ {

		dtime += getTimeDifference(line, "4", day.String()).DMinutes

		//___________________________________________________________

		//fmt.Printf("\n\nday %s\t\toee: %f", day.String(), calcOEE(line, week, day.String()))
		day = day.AddDate(0, 0, 1)
		// fmt.Printf("\nFLAG - Dtime %f\n", dtime)

	}

	pdl_calc += dtime

	lt_calc = ot_calc + pdl_calc/60

	oee = vot_calc / lt_calc * 100

	// log.Printf("OT: %f\n", ot_calc)
	// log.Printf("PDL: %f\n", pdl_calc)
	// log.Printf("LT: %f\n", lt_calc)
	// allocatedLosses = pdl_calc + mpl_calc

	// log.Printf("\nOEE! %f", oee)

	// log.Printf("\nBefore>>OEE! %f\n\n", oee)

	if math.IsNaN(float64(oee)) {
		//log.Println("Null OEE")
		oee = -1
	}
	// fmt.Println("\nDTIME: ", dtime)

	return dtime

}

func CalcAllocatedLossesProjectionbyRange(line, startDate, endDate string) float32 {

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

	var vot_calc float32
	var mpl_calc float32
	var ot_calc float32
	var pdl_calc float32
	var lt_calc float32
	var oee float32

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

	/*-------Manufactory Performance Lose-----------*/

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
		INTERVAL 6 HOUR)
    `, line, startDate, endDate)

	if err_mpl != nil {
		panic(err_mpl.Error())
	}

	p_mpl := models.MPL{}

	res_mpl := []models.MPL{}
	for selDB_mpl.Next() {

		var Minutes float32

		err = selDB_mpl.Scan(&Minutes)
		if err != nil {
			panic(err.Error())
		}
		p_mpl.Minutes = Minutes

		mpl_calc += Minutes
		//log.Printf("Minutes: %f\n", Minutes)

		res_mpl = append(res_mpl, p_mpl)
	}
	defer db.Close()

	//log.Printf("MPL Min: %f\n", mpl_calc)
	/*-------------OT--------------------------------*/
	ot_calc = vot_calc + mpl_calc/60

	/*-------Process Driven Losses-----------*/
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

	p_pdl := models.PDL{}

	res_pdl := []models.PDL{}
	for selDB_pdl.Next() {

		var Minutes float32

		err = selDB_pdl.Scan(&Minutes)
		if err != nil {
			panic(err.Error())
		}
		p_pdl.Minutes = Minutes

		pdl_calc += Minutes
		//log.Printf("Minutes: %f\n", Minutes)

		res_pdl = append(res_pdl, p_pdl)
	}
	defer db.Close()

	layout := "2006-01-02"

	tStart, erra := time.Parse(layout, startDate)
	tEnd, errb := time.Parse(layout, endDate)

	if erra != nil {
		fmt.Println(erra)
	}

	if errb != nil {
		fmt.Println(errb)
	}

	fmt.Println(tStart)

	fmt.Println(tEnd)

	diff := tEnd.Sub(tStart)
	days := diff.Hours() / 24

	// fmt.Println("Days diference: %d", days)

	// startDate = day.String()
	// endDate = day.AddDate(0, 0, 6).String()
	var dtime float32
	var allocatedLosses float32

	dtime = 0
	allocatedLosses = 0

	day := tStart

	for index := 0; index < int(days); index++ {

		dtime += getTimeDifference(line, "4", day.String()).DMinutes

		//___________________________________________________________

		//fmt.Printf("\n\nday %s\t\toee: %f", day.String(), calcOEE(line, week, day.String()))
		day = day.AddDate(0, 0, 1)
		// fmt.Printf("\nDtime %f\n", dtime)

	}

	pdl_calc += dtime

	lt_calc = ot_calc + pdl_calc/60

	oee = vot_calc / lt_calc * 100

	// log.Printf("OT: %f\n", ot_calc)
	// log.Printf("PDL: %f\n", pdl_calc)
	// log.Printf("LT: %f\n", lt_calc)
	allocatedLosses = pdl_calc + mpl_calc

	// log.Printf("\nOEE! %f", oee)

	// log.Printf("\nBefore>>OEE! %f\n\n", oee)

	if math.IsNaN(float64(oee)) {
		//log.Println("Null OEE")
		oee = -1
	}
	// fmt.Println("\n\n")

	return allocatedLosses

}
