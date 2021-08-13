package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	//"../models"
	"github.com/AdrianRb95/MES/app/models"
)

func AM_instance(w http.ResponseWriter, r *http.Request) {

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

	res := GetAM_JobList()

	tmpl.ExecuteTemplate(w, "AM_instance", res)

}

func AMStats(w http.ResponseWriter, r *http.Request) {

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

	res := GetAM_JobList()

	tmpl.ExecuteTemplate(w, "AMStats", res)
}

func getStatusJobsBy(init, end, lineId string) (int, int, int) {
	open := 0
	inProgress := 0
	closed := 0

	db := dbConn()

	selDB, err := db.Query(`
	SELECT 
		T.name, COUNT(AM.id_state)
	FROM
		AM_instance AM
			INNER JOIN
		Line L ON L.id = AM.id_line
			INNER JOIN
		Area A ON L.id_area = A.id
			INNER JOIN
		Transaction_state T ON T.id = AM.id_state
			INNER JOIN
		AM_Job J ON J.id = AM.id_job
			INNER JOIN
		Component C ON C.id = J.id_component
			INNER JOIN
		LILA_Point LILA ON LILA.id = J.id_lila
	WHERE
		DATE(AM.planned_init) BETWEEN ? AND ?
			AND L.id = ?
	GROUP BY AM.id_state
	`, init, end, lineId)

	if err != nil {
		panic(err.Error())
	}

	for selDB.Next() {

		var status int
		var statusName string

		err = selDB.Scan(&statusName, &status)

		if err != nil {
			panic(err.Error())
		}

		switch statusName {
		case "Abierto":
			open = status
		case "Por aprobar":
			inProgress = status
		case "Cerrado":
			closed = status
		}

	}

	defer db.Close()

	return open, inProgress, closed
}

func GetAM_Stats(w http.ResponseWriter, r *http.Request) {

	init := r.FormValue("init")
	end := r.FormValue("end")

	data := getAMStats(init, end)

	json.NewEncoder(w).Encode(data)
}

func getAMStats(init string, end string) []models.AMStat {

	db := dbConn()

	selDB, err := db.Query(`
	SELECT 
		L.id, L.name, COUNT(L.id)
	FROM
		AM_instance AM
			INNER JOIN
		Line L ON L.id = AM.id_line
			INNER JOIN
		Area A ON L.id_area = A.id
			INNER JOIN
		Transaction_state T ON T.id = AM.id_state
			INNER JOIN
		AM_Job J ON J.id = AM.id_job
			INNER JOIN
		Component C ON C.id = J.id_component
			INNER JOIN
		LILA_Point LILA ON LILA.id = J.id_lila
	WHERE
		DATE(AM.planned_init) BETWEEN ? AND ?
	GROUP BY L.id
	`, init, end)

	if err != nil {
		panic(err.Error())
	}

	p := models.AMStat{}

	res := []models.AMStat{}

	for selDB.Next() {
		var Id int
		var Line string
		var TotalJobs int
		var OpenJobs int
		var InProgressJobs int
		var ClosedJobs int

		err = selDB.Scan(&Id, &Line, &TotalJobs)

		if err != nil {
			panic(err.Error())
		}

		p.Id = Id
		p.Line = Line
		p.TotalJobs = TotalJobs

		OpenJobs, InProgressJobs, ClosedJobs = getStatusJobsBy(init, end, strconv.Itoa(Id))

		p.OpenJobs = OpenJobs
		p.InProgressJobs = InProgressJobs
		p.ClosedJobs = ClosedJobs

		fmt.Printf("[%d] %s\tTotalJobs: %d | Open: %d | Progress: %d | Closed: %d\n",
			Id, Line, TotalJobs, OpenJobs, InProgressJobs, ClosedJobs)

		res = append(res, p)
	}

	defer db.Close()

	return res
}

func InsertAM_instance(w http.ResponseWriter, r *http.Request) {

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

	if r.Method == "POST" {
		db := dbConn()

		job := r.FormValue("job")
		init := r.FormValue("init")
		end := r.FormValue("end")

		lineRaw := r.FormValue("line")
		line := ""

		lineList := strings.Split(lineRaw, "&")

		recurrence := r.FormValue("recurrence")
		period, _ := strconv.Atoi(r.FormValue("period"))

		frequency := 0
		factor := 0
		offset := 0

		switch recurrence {
		case "1":
			fmt.Println("7")
			frequency = 30
			factor = 0
			offset = 1

		case "7":
			fmt.Println("4")
			frequency = 4
			factor = 1
			offset = 0
		case "15":
			fmt.Println("2")
			factor = 2
			frequency = 2
			offset = 0
		case "30":
			frequency = 1
			factor = 4
			offset = 0
			fmt.Println("0")

		}

		// r.ParseForm() // Required if you don't call r.FormValue()
		// line := r.Form["line"]

		fmt.Println("\n\n>> ", line, "\n")

		layout := "2006-01-02T15:04"

		tinit, err := time.Parse(layout, init)
		tend, err := time.Parse(layout, end)

		if err != nil {
			fmt.Println(err)
		}

		fmt.Println("-------P: %s\tF: %s------------",
			period, frequency)
		for p := 0; p < period; p++ {
			for i := 0; i < frequency; i++ {

				for j := 0; j < len(lineList); j++ {
					line = strings.Split(lineList[j], "line=")[1]

					fmt.Printf("\n\njobId[%s] | %s -- %s | line: %s | recurrence: %s",
						job, tinit, tend, line, recurrence)

					smt := `
					insert into AM_instance(id_job, planned_init, planned_end, id_line)
						values(?, ?, ?, ?)
					`

					insForm, err := db.Prepare(smt)

					if err != nil {
						panic(err.Error())
					}
					insForm.Exec(job, tinit, tend, line)

				}

				tinit = tinit.AddDate(0, 0, (7*factor)+offset)
				tend = tend.AddDate(0, 0, (7*factor)+offset)

			}

			// tinit = tinit.AddDate(0, 1, 0)
			// tend = tend.AddDate(0, 1, 0)
		}

		defer db.Close()
	}

	w.Header().Set("Server", "MES")
	w.WriteHeader(200)
}

func GetMetaAM_instance(w http.ResponseWriter, r *http.Request) {

	db := dbConn()

	dateJob := r.URL.Query().Get("date")
	area := r.URL.Query().Get("area")

	selDB, err := db.Query(`
	SELECT 
		AM.id,
		AM.planned_init,
		AM.planned_end,
		L.name,
		L.id,
		J.id,
		LILA.name,
		LILA.color,
		C.name,
		T.name
	FROM
		AM_instance AM
			INNER JOIN
		Line L ON L.id = AM.id_line
			INNER JOIN
		Area A ON L.id_area = A.id
			INNER JOIN
		Transaction_state T ON T.id = AM.id_state
			INNER JOIN
		AM_Job J ON J.id = AM.id_job
			INNER JOIN
		Component C ON C.id = J.id_component
			INNER JOIN
		LILA_Point LILA ON LILA.id = J.id_lila
	WHERE
		DATE(AM.planned_init) = ?
		AND A.Id = ?

	`, dateJob, area)
	if err != nil {
		panic(err.Error())
	}

	p := models.MetaAM_instance{}

	res := []models.MetaAM_instance{}

	for selDB.Next() {
		var Id int
		var Planned_init string
		var Planned_end string
		var Line string
		var LineId int
		var JobId int
		var Lila string
		var LilaColor string
		var Component string
		var State string

		err = selDB.Scan(&Id, &Planned_init, &Planned_end, &Line, &LineId,
			&JobId, &Lila, &LilaColor, &Component, &State)
		if err != nil {
			panic(err.Error())
		}

		p.Id = Id
		p.Planned_init = Planned_init
		p.Planned_end = Planned_end
		p.Line = Line
		p.LineId = LineId
		p.JobId = JobId
		p.Lila = Lila
		p.LilaColor = LilaColor
		p.Component = Component
		p.State = State

		res = append(res, p)
	}
	defer db.Close()
	json.NewEncoder(w).Encode(res)
}

func AM_instanceLog(w http.ResponseWriter, r *http.Request) {
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
	tmpl.ExecuteTemplate(w, "AM_instanceLog", nil)
}

func GetMetaLogAM_instance(w http.ResponseWriter, r *http.Request) {

	db := dbConn()

	dateInit := r.URL.Query().Get("idate")
	dateEnd := r.URL.Query().Get("fdate")

	selDB, err := db.Query(`
	SELECT 
		AM.id,
		AM.planned_init,
		AM.planned_end,
		L.name,
		J.id,
		J.description,
		M.name,
		C.name,
		C.photo,
		E.name,
		E.photo,
		LILA.name,
		LILA.color,
		O.profile_picture,
		O.nick_name,
		O.fname,
		O.lname,
		IFNULL(AM.job_init, ''),
		IFNULL(AM.job_end, ''),
		IFNULL(AM.minutes_stop, 0),
		IFNULL(AM.minutes_run, 0),
		A.profile_picture,
		A.nick_name,
		A.fname,
		A.lname,
		T.name,
		T.color
	FROM
		AM_instance AM
			INNER JOIN
		Line L ON L.id = AM.id_line
			INNER JOIN
		Transaction_state T ON T.id = AM.id_state
			INNER JOIN
		AM_Job J ON J.id = AM.id_job
			INNER JOIN
		Component C ON C.id = J.id_component
			INNER JOIN
		Machine M ON M.id = C.id_machine
			INNER JOIN
		EPP E ON E.id = J.id_epp
			INNER JOIN
		LILA_Point LILA ON LILA.id = J.id_lila
			INNER JOIN
		User_table O ON O.id = AM.id_operator
			INNER JOIN
		User_table A ON A.id = AM.id_approver
	WHERE
		DATE(AM.planned_init) BETWEEN ? AND ?

	`, dateInit, dateEnd)
	if err != nil {
		panic(err.Error())
	}

	p := models.MetaLogAM_instance{}

	res := []models.MetaLogAM_instance{}

	for selDB.Next() {

		var Id int
		var Planned_init string
		var Planned_end string
		var Line string

		var JobId int
		var Description string

		var Machine string
		var Component string
		var ComponentPhoto string

		var EPP string
		var EPPPhoto string

		var Lila string
		var LilaColor string

		var OperatorProfile string
		var OperatorNickName string
		var OperatorFname string
		var OperatorLname string

		var JobInit string
		var JobEmd string
		var MinutesStop int
		var MinutesRun int

		var ApproverProfile string
		var ApproverNickName string
		var ApproverFname string
		var ApproverLname string
		var State string
		var StateColor string

		err = selDB.Scan(
			&Id,
			&Planned_init,
			&Planned_end,
			&Line,
			&JobId,
			&Description,
			&Machine,
			&Component,
			&ComponentPhoto,
			&EPP,
			&EPPPhoto,
			&Lila,
			&LilaColor,
			&OperatorProfile,
			&OperatorNickName,
			&OperatorFname,
			&OperatorLname,
			&JobInit,
			&JobEmd,
			&MinutesStop,
			&MinutesRun,
			&ApproverProfile,
			&ApproverNickName,
			&ApproverFname,
			&ApproverLname,
			&State,
			&StateColor)

		if err != nil {
			panic(err.Error())
		}

		p.Id = Id
		p.Planned_init = Planned_init
		p.Planned_end = Planned_end

		p.Line = Line
		p.JobId = JobId
		p.Description = Description

		p.Machine = Machine
		p.Component = Component
		p.ComponentPhoto = ComponentPhoto

		p.EPP = EPP
		p.EPPPhoto = EPPPhoto

		p.Lila = Lila
		p.LilaColor = LilaColor

		p.OperatorProfile = OperatorProfile
		p.OperatorNickName = OperatorNickName
		p.OperatorFname = OperatorFname
		p.OperatorLname = OperatorLname

		p.JobInit = JobInit
		p.JobEmd = JobEmd
		p.MinutesStop = MinutesStop
		p.MinutesRun = MinutesRun

		p.ApproverProfile = ApproverProfile
		p.ApproverNickName = ApproverNickName
		p.ApproverFname = ApproverFname
		p.ApproverLname = ApproverLname
		p.State = State
		p.StateColor = StateColor

		res = append(res, p)
	}
	defer db.Close()
	json.NewEncoder(w).Encode(res)
}

func GetMetaLogAM_instanceE(w http.ResponseWriter, r *http.Request) {

	db := dbConn()

	Id := r.URL.Query().Get("pid")

	selDB, err := db.Query(`
	SELECT 
		AM.id,
		AM.planned_init,
		AM.planned_end,
		L.name,
		L.id,
		J.id,
		J.description,
		M.name,
		C.name,
		C.photo,
		E.name,
		E.photo,
		LILA.name,
		LILA.color,
		O.profile_picture,
		O.nick_name,
		O.fname,
		O.lname,
		IFNULL(AM.job_init, ''),
		IFNULL(AM.job_end, ''),
		IFNULL(AM.minutes_stop, 0),
		IFNULL(AM.minutes_run, 0),
		A.profile_picture,
		A.nick_name,
		A.fname,
		A.lname,
		T.name,
		T.color,
		AM.phase
	FROM
		AM_instance AM
			INNER JOIN
		Line L ON L.id = AM.id_line
			INNER JOIN
		Transaction_state T ON T.id = AM.id_state
			INNER JOIN
		AM_Job J ON J.id = AM.id_job
			INNER JOIN
		Component C ON C.id = J.id_component
			INNER JOIN
		Machine M ON M.id = C.id_machine
			INNER JOIN
		EPP E ON E.id = J.id_epp
			INNER JOIN
		LILA_Point LILA ON LILA.id = J.id_lila
			INNER JOIN
		User_table O ON O.id = AM.id_operator
			INNER JOIN
		User_table A ON A.id = AM.id_approver
	WHERE
		AM.id = ?

	`, Id)

	if err != nil {
		panic(err.Error())
	}

	p := models.MetaLogAM_instance{}

	res := []models.MetaLogAM_instance{}

	for selDB.Next() {

		var Id int
		var Planned_init string
		var Planned_end string
		var Line string
		var LineId int

		var JobId int
		var Description string

		var Machine string
		var Component string
		var ComponentPhoto string

		var EPP string
		var EPPPhoto string

		var Lila string
		var LilaColor string

		var OperatorProfile string
		var OperatorNickName string
		var OperatorFname string
		var OperatorLname string

		var JobInit string
		var JobEmd string
		var MinutesStop int
		var MinutesRun int

		var ApproverProfile string
		var ApproverNickName string
		var ApproverFname string
		var ApproverLname string
		var State string
		var StateColor string
		var Phase int

		err = selDB.Scan(
			&Id,
			&Planned_init,
			&Planned_end,
			&Line,
			&LineId,
			&JobId,
			&Description,
			&Machine,
			&Component,
			&ComponentPhoto,
			&EPP,
			&EPPPhoto,
			&Lila,
			&LilaColor,
			&OperatorProfile,
			&OperatorNickName,
			&OperatorFname,
			&OperatorLname,
			&JobInit,
			&JobEmd,
			&MinutesStop,
			&MinutesRun,
			&ApproverProfile,
			&ApproverNickName,
			&ApproverFname,
			&ApproverLname,
			&State,
			&StateColor,
			&Phase)

		if err != nil {
			panic(err.Error())
		}

		p.Id = Id
		p.Planned_init = Planned_init
		p.Planned_end = Planned_end

		p.Line = Line
		p.LineId = LineId

		p.JobId = JobId
		p.Description = Description

		p.Machine = Machine
		p.Component = Component
		p.ComponentPhoto = ComponentPhoto

		p.EPP = EPP
		p.EPPPhoto = EPPPhoto

		p.Lila = Lila
		p.LilaColor = LilaColor

		p.OperatorProfile = OperatorProfile
		p.OperatorNickName = OperatorNickName
		p.OperatorFname = OperatorFname
		p.OperatorLname = OperatorLname

		p.JobInit = JobInit
		p.JobEmd = JobEmd
		p.MinutesStop = MinutesStop
		p.MinutesRun = MinutesRun

		p.ApproverProfile = ApproverProfile
		p.ApproverNickName = ApproverNickName
		p.ApproverFname = ApproverFname
		p.ApproverLname = ApproverLname
		p.State = State
		p.StateColor = StateColor
		p.Phase = Phase

		res = append(res, p)

	}
	defer db.Close()
	json.NewEncoder(w).Encode(res)
}

func DeleteAM_Intance(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	p := r.URL.Query().Get("id")
	delForm, err := db.Prepare("DELETE FROM AM_instance WHERE id=?")
	if err != nil {
		panic(err.Error())
	}
	delForm.Exec(p)
	log.Println("DELETE AM_instance")
	defer db.Close()
}

func GetMetaLogAM_instanceEList(w http.ResponseWriter, r *http.Request) {

	db := dbConn()

	line := r.URL.Query().Get("line")
	init := r.URL.Query().Get("init")
	turn := r.URL.Query().Get("turn")
	pstate := r.URL.Query().Get("pstate")

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

	minutesQ = (hoursQ+8)*60 - 1

	selDB, err := db.Query(`
	SELECT 
		AM.id,
		AM.planned_init,
		AM.planned_end,
		L.name,
		L.id, 
		J.id,
		J.description,
		M.name,
		C.name,
		C.photo,
		E.name,
		E.photo,
		LILA.name,
		LILA.color,
		O.profile_picture,
		O.nick_name,
		O.fname,
		O.lname,
		IFNULL(AM.job_init, ''),
		IFNULL(AM.job_end, ''),
		IFNULL(AM.minutes_stop, 0),
		IFNULL(AM.minutes_run, 0),
		A.profile_picture,
		A.nick_name,
		A.fname,
		A.lname,
		T.name,
		T.color,
		AM.phase
	FROM
		AM_instance AM
			INNER JOIN
		Line L ON L.id = AM.id_line
			INNER JOIN
		Transaction_state T ON T.id = AM.id_state
			INNER JOIN
		AM_Job J ON J.id = AM.id_job
			INNER JOIN
		Component C ON C.id = J.id_component
			INNER JOIN
		Machine M ON M.id = C.id_machine
			INNER JOIN
		EPP E ON E.id = J.id_epp
			INNER JOIN
		LILA_Point LILA ON LILA.id = J.id_lila
			INNER JOIN
		User_table O ON O.id = AM.id_operator
			INNER JOIN
		User_table A ON A.id = AM.id_approver
	WHERE
		AM.planned_init BETWEEN DATE_ADD(DATE(?),
        INTERVAL `+strconv.Itoa(hoursQ)+`  HOUR) AND DATE_ADD(DATE(?),
        INTERVAL  `+strconv.Itoa(minutesQ)+`  MINUTE)
			AND L.id = ?
			AND AM.id_state = ? 
	`, init, init, line, pstate)

	if err != nil {
		panic(err.Error())
	}

	p := models.MetaLogAM_instance{}

	res := []models.MetaLogAM_instance{}

	for selDB.Next() {

		var Id int
		var Planned_init string
		var Planned_end string
		var Line string
		var LineId int

		var JobId int
		var Description string

		var Machine string
		var Component string
		var ComponentPhoto string

		var EPP string
		var EPPPhoto string

		var Lila string
		var LilaColor string

		var OperatorProfile string
		var OperatorNickName string
		var OperatorFname string
		var OperatorLname string

		var JobInit string
		var JobEmd string
		var MinutesStop int
		var MinutesRun int

		var ApproverProfile string
		var ApproverNickName string
		var ApproverFname string
		var ApproverLname string
		var State string
		var StateColor string
		var Phase int

		err = selDB.Scan(
			&Id,
			&Planned_init,
			&Planned_end,
			&Line,
			&LineId,
			&JobId,
			&Description,
			&Machine,
			&Component,
			&ComponentPhoto,
			&EPP,
			&EPPPhoto,
			&Lila,
			&LilaColor,
			&OperatorProfile,
			&OperatorNickName,
			&OperatorFname,
			&OperatorLname,
			&JobInit,
			&JobEmd,
			&MinutesStop,
			&MinutesRun,
			&ApproverProfile,
			&ApproverNickName,
			&ApproverFname,
			&ApproverLname,
			&State,
			&StateColor,
			&Phase)

		if err != nil {
			panic(err.Error())
		}

		p.Id = Id
		p.Planned_init = Planned_init
		p.Planned_end = Planned_end

		p.Line = Line
		p.LineId = LineId
		p.JobId = JobId
		p.Description = Description

		p.Machine = Machine
		p.Component = Component
		p.ComponentPhoto = ComponentPhoto

		p.EPP = EPP
		p.EPPPhoto = EPPPhoto

		p.Lila = Lila
		p.LilaColor = LilaColor

		p.OperatorProfile = OperatorProfile
		p.OperatorNickName = OperatorNickName
		p.OperatorFname = OperatorFname
		p.OperatorLname = OperatorLname

		p.JobInit = JobInit
		p.JobEmd = JobEmd
		p.MinutesStop = MinutesStop
		p.MinutesRun = MinutesRun

		p.ApproverProfile = ApproverProfile
		p.ApproverNickName = ApproverNickName
		p.ApproverFname = ApproverFname
		p.ApproverLname = ApproverLname
		p.State = State
		p.StateColor = StateColor
		p.Phase = Phase

		res = append(res, p)

	}
	defer db.Close()
	json.NewEncoder(w).Encode(res)
}

func ExecuteAM_instance(w http.ResponseWriter, r *http.Request) {

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

	tmpl.ExecuteTemplate(w, "ExecuteAM_instance", user)

}

func SetAMPhase(w http.ResponseWriter, r *http.Request) {

	phaseName := []string{"job_init", "job_end"}

	status := 401
	state := "1"

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

	if r.Method == "POST" {
		db := dbConn()
		uid := r.FormValue("id")
		user := r.FormValue("user")
		// pass := r.FormValue("pass")
		uphase := r.FormValue("phase")

		var grand bool

		n, _ := strconv.ParseInt(uphase, 10, 64)

		index := n

		n++

		phase := strconv.FormatInt(n, 10)

		date := strings.Split(time.Now().String(), " ")[0]
		ptime := strings.Split(time.Now().String(), " ")[1]

		utime := date + " " + strings.Split(ptime, ".")[0]

		if uphase == "0" {
			grand = true

			smt := `
			UPDATE AM_instance 
			SET 
				phase = ?,
				id_state= ?,
				` + phaseName[0] + `=?
			WHERE
				id = ?
			`

			insForm, err := db.Prepare(smt)

			if err != nil {
				panic(err.Error())
			}
			insForm.Exec(phase, state, utime, uid)

			status = 200

		}

		fmt.Printf("\njob: %s\t|\tphase: %s\t|\ttime: %s\n",
			uid, phase, utime)

		userID := GetUserByNickName(user)

		grand = isOperator(userID)

		if grand == true {

			if phase == "2" {
				state = "4"

				if grand == true {

					smt := `
					UPDATE AM_instance 
					SET 
						phase = ?,
						id_state= ?,
						id_approver=?
					WHERE
						id = ?
					`

					insForm, err := db.Prepare(smt)

					if err != nil {
						panic(err.Error())
					}
					insForm.Exec("2", state, userID, uid)

					status = 200
				}

			} else {

				if index > 1 && index < 3 {
					state = "2"
				}

				smt := `
				UPDATE AM_instance 
				SET 
					phase = ?,
					id_state= ?,
					id_operator=?,
					` + phaseName[index] + `=?
				WHERE
					id = ?
				`

				insForm, err := db.Prepare(smt)

				if err != nil {
					panic(err.Error())
				}
				insForm.Exec(phase, state, userID, utime, uid)

				status = 200
			}

		}

		if phase == "3" {
			state = "3"
			fmt.Println("1-Flag")
			if grand == false {
				fmt.Println("2-Flag")

				smt := `
				UPDATE AM_instance 
				SET 
					phase = ?,
					id_state= ?,
					id_approver=?
				WHERE
					id = ?
				`

				insForm, err := db.Prepare(smt)

				if err != nil {
					panic(err.Error())
				}
				insForm.Exec("3", state, userID, uid)

				status = 200
			}

		}

		defer db.Close()

	}
	fmt.Printf("Status:  %d\n", status)

	w.Header().Set("Server", "MES")
	w.WriteHeader(status)
}

func SetAMFlashPhase(w http.ResponseWriter, r *http.Request) {

	status := 200

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

	if r.Method == "POST" {
		db := dbConn()
		uid := r.FormValue("id")
		minutes := r.FormValue("minutes")
		condition := r.FormValue("condition")

		userID := user.UserId

		state := "4"

		smt := `
				UPDATE AM_instance 
				SET 
					phase = ?,
					id_state= ?,
					id_operator=?,
					comment=?,
					minutes_stop=?
				WHERE
					id = ?
					`

		insForm, err := db.Prepare(smt)

		if err != nil {
			panic(err.Error())
		}
		insForm.Exec("2", state, userID, condition, minutes, uid)

		fmt.Println("AM| %s | %s | %s | %s | %s ", state, userID, condition, minutes, uid)

		defer db.Close()
	}
	fmt.Printf("Status:  %d\n", status)

	w.Header().Set("Server", "MES")
	w.WriteHeader(status)
}

func SetAMNote(w http.ResponseWriter, r *http.Request) {

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

	if r.Method == "POST" {
		db := dbConn()
		uid := r.FormValue("id")
		note := r.FormValue("note")
		run := r.FormValue("run")
		stop := r.FormValue("stop")

		fmt.Printf("job: %s\t|\tNote: %s \t|run: %s stop: %s\n",
			uid, note, run, stop)

		smt := `
		UPDATE AM_instance 
		SET 
			comment=?,
			minutes_run=?,
			minutes_stop=?
		WHERE
			id = ?
		`

		insForm, err := db.Prepare(smt)

		if err != nil {
			panic(err.Error())
		}
		insForm.Exec(note, run, stop, uid)

		defer db.Close()
	}

	w.Header().Set("Server", "MES")
	w.WriteHeader(200)
}
