package controllers

import (
	//"bytes"

	"math"
	"net/http"

	//"image/jpeg"
	"encoding/json"

	//"../models"
	"github.com/AdrianRb95/MES/app/models"

	_ "github.com/go-sql-driver/mysql"
)

func GetDTime(w http.ResponseWriter, r *http.Request) {
	turn := r.URL.Query().Get("turn")
	line := r.URL.Query().Get("line")
	date_event := r.URL.Query().Get("date")
	var dtime models.Dtime

	dtime = getTimeDifference(line, turn, date_event)

	json.NewEncoder(w).Encode(dtime)
}

func getTimeDifference(line, turn, date string) models.Dtime {

	plannedData := getPlannedDataByDay(line, turn, date)
	var goodProduct float32
	var loadingTime float32
	var plantStoppage float32
	var stoppages float32
	var operationalTime float32
	var teoricalProduction float32
	var tNeto_teoricalProduction float32
	var dHours float32
	var dMinutes float32
	var mpl, pdl float32
	var dtime models.Dtime
	var turnTime float32

	goodProduct = float32(plannedData.Produced) * float32(plannedData.Box_amount) * .001
	plantStoppage = getPlantStoppageByDay(line, turn, date)
	mpl = getMPLByDay(line, turn, date)
	pdl = getPDLByDay(line, turn, date)

	stoppages = mpl + pdl

	if turn == "4" {
		turnTime = 8 * 3
	} else {
		turnTime = 8
	}

	loadingTime = turnTime - plantStoppage/60

	operationalTime = loadingTime - stoppages/60

	teoricalProduction = operationalTime * 60 * plannedData.Nominal_speed * .001

	tNeto_teoricalProduction = goodProduct * 1000 / plannedData.Nominal_speed / 60

	dHours = operationalTime - tNeto_teoricalProduction
	dMinutes = dHours * 60

	// fmt.Printf("goodProduct =  %d *  %d / 1000", plannedData.Produced, plannedData.Box_amount)

	// fmt.Printf("\n\n____________________________________________________\nTimeDifference\n__________________________________________________\n")
	// fmt.Printf("plantStoppage:\t\t\t%f\n", plantStoppage)
	// fmt.Printf("MPL:\t\t\t\t%f\n", mpl)
	// fmt.Printf("PDL:\t\t\t\t%f\n", pdl)
	// fmt.Printf("goodProduct:\t\t\t%f\n", goodProduct)
	// fmt.Printf("loadingTime:\t\t\t%f\n", loadingTime)
	// fmt.Printf("operationalTime:\t\t%f\n", operationalTime)
	// fmt.Printf("teoricalProduction:\t\t%f\n", teoricalProduction)
	// fmt.Printf("tneto/teoricalValue\t\t%f\n", tNeto_teoricalProduction)
	// fmt.Printf("dHours:\t\t\t\t%f\n", dHours)
	// fmt.Printf("dMinutes:\t\t\t%f\n", dMinutes)

	// fmt.Printf("______________________________________________________________________________\n\n")

	dtime.OperationalTime = operationalTime * 60
	dtime.TeoricalProduction = teoricalProduction
	dtime.DMinutes = dMinutes
	dtime.TNeto_teoricalProduction = tNeto_teoricalProduction * 60

	if math.IsNaN(float64(dtime.DMinutes)) {
		//log.Println("Null OEE")
		dtime.DMinutes = 0
	}

	return dtime
}
