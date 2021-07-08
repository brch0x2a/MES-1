package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"../models"
)

const (
	layoutISO = "2006-01-02"
	layoutUS  = "January 2, 2006"
)

func GetQACurrentMonitor(w http.ResponseWriter, r *http.Request) {
	line := r.URL.Query().Get("line")

	qaMonitor := models.QAMonitor{}

	date := time.Now().Format(layoutISO)

	turn := GetRelativeCurrentTurn(time.Now())

	qaMonitor.WeightCounter = GetWeightCurrentCounter(line, date, turn)
	qaMonitor.ProcessTemperature = GetProcessTempatureCurrentCounter(line, date, turn)
	qaMonitor.JawTeflonState = GetJawTeflonStateCurrentCounter(line, date, turn)
	qaMonitor.SealVerification = GetSealVerificationCurrentCounter(line, date, turn)
	qaMonitor.CRQS = GetCRQSCurrentCounter(line, date, turn)
	qaMonitor.BatchCounter = GetBatchCurrentCounter(line, date, turn)
	qaMonitor.AllergenCounter = GetAllergenCurrentCounter(line, date, turn)
	qaMonitor.CleanDisinfection = GetCleanDisinfectionCounter(line, date, turn)
	qaMonitor.CodingVerification = GetCodingVerificationCounter(line, date, turn)

	json.NewEncoder(w).Encode(qaMonitor)
}

func GetQARelativeMonitorWeekUtilitation(w http.ResponseWriter, r *http.Request) {
	line := r.URL.Query().Get("line")
	// date := r.URL.Query().Get("date")

	startDate := r.URL.Query().Get("init")
	endDate := r.URL.Query().Get("end")

	layout := "2006-01-02"

	tinit, err := time.Parse(layout, startDate)
	tend, err := time.Parse(layout, endDate)

	if err != nil {
		fmt.Println(err)
	}

	diff := int(tend.Sub(tinit).Hours() / 24)

	fmt.Println("Time Range:\t", diff)

	turn := 4

	qaMonitor := models.QAMonitor{}
	lineList := []models.QAMonitor{}

	for i := 0; i < diff; i++ {
		date := tinit.Format(layout)

		qaMonitor.WeightCounter = GetWeightCurrentCounter(line, date, turn)
		qaMonitor.ProcessTemperature = GetProcessTempatureCurrentCounter(line, date, turn)
		qaMonitor.JawTeflonState = GetJawTeflonStateCurrentCounter(line, date, turn)
		qaMonitor.SealVerification = GetSealVerificationCurrentCounter(line, date, turn)
		qaMonitor.CRQS = GetCRQSCurrentCounter(line, date, turn)
		qaMonitor.BatchCounter = GetBatchCurrentCounter(line, date, turn)
		qaMonitor.AllergenCounter = GetAllergenCurrentCounter(line, date, turn)

		lineList = append(lineList, qaMonitor)

		tinit = tinit.AddDate(0, 0, 1)
	}

	json.NewEncoder(w).Encode(lineList)
}

func GetQARelativeMonitor(w http.ResponseWriter, r *http.Request) {
	line := r.URL.Query().Get("line")
	date := r.URL.Query().Get("date")
	turn, _ := strconv.Atoi(r.URL.Query().Get("turn"))

	qaMonitor := models.QAMonitor{}

	qaMonitor.WeightCounter = GetWeightCurrentCounter(line, date, turn)
	qaMonitor.ProcessTemperature = GetProcessTempatureCurrentCounter(line, date, turn)
	qaMonitor.JawTeflonState = GetJawTeflonStateCurrentCounter(line, date, turn)
	qaMonitor.SealVerification = GetSealVerificationCurrentCounter(line, date, turn)
	qaMonitor.CRQS = GetCRQSCurrentCounter(line, date, turn)
	qaMonitor.BatchCounter = GetBatchCurrentCounter(line, date, turn)
	qaMonitor.AllergenCounter = GetAllergenCurrentCounter(line, date, turn)
	qaMonitor.CleanDisinfection = GetCleanDisinfectionCounter(line, date, turn)
	qaMonitor.CodingVerification = GetCodingVerificationCounter(line, date, turn)

	json.NewEncoder(w).Encode(qaMonitor)
}

func GetWeightCurrentCounter(line string, date string, turn int) int {
	counter := 0

	hoursQ := 6
	minutesQ := 0

	// turn := GetRelativeCurrentTurn()

	// date := time.Now().Format(layoutISO)

	switch turn {

	case 1:
		{
			hoursQ = 6
		}
	case 2:
		{
			hoursQ = 14
		}
	case 3:
		{
			hoursQ = 22
		}

	}
	minutesQ = (hoursQ+8)*60 - 1

	if turn == 4 {
		minutesQ = 1800 - 1
	}

	db := dbConn()

	selDB, err := db.Query(`
	SELECT 
    	COUNT(W.id)
	FROM
    	Weight_control W
        INNER JOIN
    Salsitas_statistical_process S ON W.id_sub_header = S.id
	WHERE
		S.id_line = ?
		AND W.date_reg BETWEEN DATE_ADD(DATE(?), INTERVAL `+strconv.Itoa(hoursQ)+` HOUR) 
			AND DATE_ADD(DATE(?), INTERVAL `+strconv.Itoa(minutesQ)+` MINUTE);

		`, line, date, date)

	if err != nil {
		panic(err.Error())
	}

	for selDB.Next() {

		err = selDB.Scan(&counter)

		if err != nil {
			panic(err.Error())
		}

	}
	defer db.Close()

	return counter
}

func GetProcessTempatureCurrentCounter(line string, date string, turn int) int {
	counter := 0

	hoursQ := 6
	minutesQ := 0

	// turn := GetRelativeCurrentTurn()

	// date := time.Now().Format(layoutISO)

	switch turn {

	case 1:
		{
			hoursQ = 6
		}
	case 2:
		{
			hoursQ = 14
		}
	case 3:
		{
			hoursQ = 22
		}

	}
	minutesQ = (hoursQ+8)*60 - 1

	if turn == 4 {
		minutesQ = 1800 - 1
	}

	db := dbConn()

	selDB, err := db.Query(`
	SELECT 
		COUNT(P.id)
	FROM
		Process_temperature_control P
		INNER JOIN
	Salsitas_statistical_process S ON P.id_sub_header = S.id
	WHERE
		S.id_line = ?
		AND P.date_reg BETWEEN DATE_ADD(DATE(?), INTERVAL `+strconv.Itoa(hoursQ)+` HOUR) 
			AND DATE_ADD(DATE(?), INTERVAL `+strconv.Itoa(minutesQ)+` MINUTE);

		`, line, date, date)

	if err != nil {
		panic(err.Error())
	}

	for selDB.Next() {

		err = selDB.Scan(&counter)

		if err != nil {
			panic(err.Error())
		}

	}
	defer db.Close()

	return counter
}

func GetJawTeflonStateCurrentCounter(line string, date string, turn int) int {
	counter := 0

	hoursQ := 6
	minutesQ := 0

	// turn := GetRelativeCurrentTurn()

	// date := time.Now().Format(layoutISO)

	switch turn {

	case 1:
		{
			hoursQ = 6
		}
	case 2:
		{
			hoursQ = 14
		}
	case 3:
		{
			hoursQ = 22
		}

	}
	minutesQ = (hoursQ+8)*60 - 1

	if turn == 4 {
		minutesQ = 1800 - 1
	}

	db := dbConn()

	selDB, err := db.Query(`
	SELECT 
		COUNT(J.id)
	FROM
		Jaw_teflon_ultrasonic_state J
		INNER JOIN
	Salsitas_statistical_process S ON J.id_sub_header = S.id
	WHERE
		S.id_line = ?
		AND J.date_reg BETWEEN DATE_ADD(DATE(?), INTERVAL `+strconv.Itoa(hoursQ)+` HOUR) 
			AND DATE_ADD(DATE(?), INTERVAL `+strconv.Itoa(minutesQ)+` MINUTE);

		`, line, date, date)

	if err != nil {
		panic(err.Error())
	}

	for selDB.Next() {

		err = selDB.Scan(&counter)

		if err != nil {
			panic(err.Error())
		}

	}
	defer db.Close()

	return counter
}

func GetSealVerificationCurrentCounter(line string, date string, turn int) int {
	counter := 0

	hoursQ := 6
	minutesQ := 0

	// turn := GetRelativeCurrentTurn()

	// date := time.Now().Format(layoutISO)

	switch turn {

	case 1:
		{
			hoursQ = 6
		}
	case 2:
		{
			hoursQ = 14
		}
	case 3:
		{
			hoursQ = 22
		}

	}
	minutesQ = (hoursQ+8)*60 - 1

	if turn == 4 {
		minutesQ = 1800 - 1
	}

	db := dbConn()

	selDB, err := db.Query(`
	SELECT 
		COUNT(SP.id)
	FROM
		Seal_verification_pneumaticpress SP
		INNER JOIN
	Salsitas_statistical_process S ON SP.id_sub_header = S.id
	WHERE
		S.id_line = ?
		AND SP.date_reg BETWEEN DATE_ADD(DATE(?), INTERVAL `+strconv.Itoa(hoursQ)+` HOUR) 
			AND DATE_ADD(DATE(?), INTERVAL `+strconv.Itoa(minutesQ)+` MINUTE);

		`, line, date, date)

	if err != nil {
		panic(err.Error())
	}

	for selDB.Next() {

		err = selDB.Scan(&counter)

		if err != nil {
			panic(err.Error())
		}

	}
	defer db.Close()

	return counter
}

func GetCRQSCurrentCounter(line string, date string, turn int) int {
	counter := 0

	hoursQ := 6
	minutesQ := 0

	// turn := GetRelativeCurrentTurn()

	// date := time.Now().Format(layoutISO)

	switch turn {

	case 1:
		{
			hoursQ = 6
		}
	case 2:
		{
			hoursQ = 14
		}
	case 3:
		{
			hoursQ = 22
		}

	}
	minutesQ = (hoursQ+8)*60 - 1

	if turn == 4 {
		minutesQ = 1800 - 1
	}

	db := dbConn()

	selDB, err := db.Query(`

	SELECT 
		COUNT(C.id)
	FROM
		CRQS C
		INNER JOIN
		Salsitas_statistical_process S ON C.id_sub_header = S.id
	WHERE
		S.id_line = ?
		AND C.date_reg BETWEEN DATE_ADD(DATE(?), INTERVAL `+strconv.Itoa(hoursQ)+` HOUR) 
			AND DATE_ADD(DATE(?), INTERVAL `+strconv.Itoa(minutesQ)+` MINUTE);

		`, line, date, date)

	if err != nil {
		panic(err.Error())
	}

	for selDB.Next() {

		err = selDB.Scan(&counter)

		if err != nil {
			panic(err.Error())
		}

	}
	defer db.Close()

	return counter
}

func GetAllergenCurrentCounter(line string, date string, turn int) int {
	counter := 0

	hoursQ := 6
	minutesQ := 0

	// turn := GetRelativeCurrentTurn()

	// date := time.Now().Format(layoutISO)

	switch turn {

	case 1:
		{
			hoursQ = 6
		}
	case 2:
		{
			hoursQ = 14
		}
	case 3:
		{
			hoursQ = 22
		}

	}
	minutesQ = (hoursQ+8)*60 - 1

	if turn == 4 {
		minutesQ = 1800 - 1
	}

	db := dbConn()

	selDB, err := db.Query(`
	SELECT 
	COUNT(A.id)
	FROM
		Allergen_Verification A
		INNER JOIN
		Salsitas_statistical_process S ON A.id_sub_header = S.id
	WHERE
	S.id_line = ?
		AND A.date_reg BETWEEN DATE_ADD(DATE(?), INTERVAL `+strconv.Itoa(hoursQ)+` HOUR) 
			AND DATE_ADD(DATE(?), INTERVAL `+strconv.Itoa(minutesQ)+` MINUTE);

		`, line, date, date)

	if err != nil {
		panic(err.Error())
	}

	for selDB.Next() {

		err = selDB.Scan(&counter)

		if err != nil {
			panic(err.Error())
		}

	}
	defer db.Close()

	return counter
}

func GetBatchCurrentCounter(line string, date string, turn int) int {
	counter := 0

	hoursQ := 6
	minutesQ := 0

	// turn := GetRelativeCurrentTurn()

	// date := time.Now().Format(layoutISO)

	switch turn {

	case 1:
		{
			hoursQ = 6
		}
	case 2:
		{
			hoursQ = 14
		}
	case 3:
		{
			hoursQ = 22
		}

	}
	minutesQ = (hoursQ+8)*60 - 1

	if turn == 4 {
		minutesQ = 1800 - 1
	}

	db := dbConn()

	selDB, err := db.Query(`

	SELECT 
		COUNT(B.id)
	FROM
		Batch B
	WHERE
		B.id_line = ?
		AND B.change_date BETWEEN DATE_ADD(DATE(?), INTERVAL `+strconv.Itoa(hoursQ)+` HOUR) 
			AND DATE_ADD(DATE(?), INTERVAL `+strconv.Itoa(minutesQ)+` MINUTE);

		`, line, date, date)

	if err != nil {
		panic(err.Error())
	}

	for selDB.Next() {

		err = selDB.Scan(&counter)

		if err != nil {
			panic(err.Error())
		}

	}
	defer db.Close()

	return counter
}

func GetCleanDisinfectionCounter(line string, date string, turn int) int {
	counter := 0

	hoursQ := 6
	minutesQ := 0

	// turn := GetRelativeCurrentTurn()

	// date := time.Now().Format(layoutISO)

	switch turn {

	case 1:
		{
			hoursQ = 6
		}
	case 2:
		{
			hoursQ = 14
		}
	case 3:
		{
			hoursQ = 22
		}

	}
	minutesQ = (hoursQ+8)*60 - 1

	if turn == 4 {
		minutesQ = 1800 - 1
	}

	db := dbConn()

	selDB, err := db.Query(`
	SELECT 
		COUNT(P.id)
	FROM
		Clean_Disinfection P
	WHERE
		P.id_line = ?
		AND P.date_init BETWEEN DATE_ADD(DATE(?), INTERVAL `+strconv.Itoa(hoursQ)+` HOUR) 
			AND DATE_ADD(DATE(?), INTERVAL `+strconv.Itoa(minutesQ)+` MINUTE);

		`, line, date, date)

	if err != nil {
		panic(err.Error())
	}

	for selDB.Next() {

		err = selDB.Scan(&counter)

		if err != nil {
			panic(err.Error())
		}

	}
	defer db.Close()

	return counter
}

func GetCodingVerificationCounter(line string, date string, turn int) int {
	counter := 0

	hoursQ := 6
	minutesQ := 0

	// turn := GetRelativeCurrentTurn()

	// date := time.Now().Format(layoutISO)

	switch turn {

	case 1:
		{
			hoursQ = 6
		}
	case 2:
		{
			hoursQ = 14
		}
	case 3:
		{
			hoursQ = 22
		}

	}
	minutesQ = (hoursQ+8)*60 - 1

	if turn == 4 {
		minutesQ = 1800 - 1
	}

	db := dbConn()

	selDB, err := db.Query(`
	SELECT 
		COUNT(P.id)
	FROM
		Coding_Verification P
	WHERE
		P.id_line = ?
		AND P.date_reg BETWEEN DATE_ADD(DATE(?), INTERVAL `+strconv.Itoa(hoursQ)+` HOUR) 
			AND DATE_ADD(DATE(?), INTERVAL `+strconv.Itoa(minutesQ)+` MINUTE);

		`, line, date, date)

	if err != nil {
		panic(err.Error())
	}

	for selDB.Next() {

		err = selDB.Scan(&counter)

		if err != nil {
			panic(err.Error())
		}

	}
	defer db.Close()

	return counter
}
