package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"../models"
)

func Work_order(w http.ResponseWriter, r *http.Request) {
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
	tmpl.ExecuteTemplate(w, "Work_order", nil)
}

func RelativeWork_order(w http.ResponseWriter, r *http.Request) {
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
	tmpl.ExecuteTemplate(w, "RelativeWork_order", user)
}

func InsertWork_order(w http.ResponseWriter, r *http.Request) {

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
		mech := r.FormValue("mech")
		line := r.FormValue("line")
		description := r.FormValue("description")

		smt := `
		insert into Work_Order(id_job, description, planned_init, planned_end, id_mecanic, id_line)
		values(?, ?, ?, ?, ?, ?)
		`

		insForm, err := db.Prepare(smt)

		if err != nil {
			panic(err.Error())
		}
		insForm.Exec(job, description, init, end, mech, line)

		defer db.Close()
	}

	w.Header().Set("Server", "MES")
	w.WriteHeader(200)
}

func GetMetaWork_order(w http.ResponseWriter, r *http.Request) {

	db := dbConn()

	dateJob := r.URL.Query().Get("date")

	selDB, err := db.Query(`
	SELECT 
		W.id,
		J.name,
		J.color,
		W.description,
		W.planned_init,
		W.planned_end,
		M.fname,
		M.lname,
		T.name,
		L.name
	FROM
	    Work_Order W
			INNER JOIN
		Job_catalog J ON J.id = W.id_job
			INNER JOIN
		User_table M ON M.id = W.id_mecanic
			INNER JOIN
		Line L ON L.id = W.id_line
			INNER JOIN
		Transaction_state T ON T.id = W.id_state
	WHERE
		DATE(W.planned_init) = ?

	`, dateJob)
	if err != nil {
		panic(err.Error())
	}

	p := models.MetaWork_order{}

	res := []models.MetaWork_order{}

	for selDB.Next() {
		var Id int
		var Job string
		var Color string
		var Description string
		var Init string
		var End string
		var Fname string
		var Lname string
		var State string
		var Line string

		err = selDB.Scan(&Id, &Job, &Color,
			&Description, &Init, &End, &Fname, &Lname, &State, &Line)
		if err != nil {
			panic(err.Error())
		}

		p.Id = Id
		p.Job = Job
		p.Color = Color
		p.Description = Description
		p.Init = Init
		p.End = End
		p.Fname = Fname
		p.Lname = Lname
		p.State = State
		p.Line = Line

		res = append(res, p)
	}
	defer db.Close()
	json.NewEncoder(w).Encode(res)
}

func GetMetaWork_orderE(w http.ResponseWriter, r *http.Request) {

	nId := r.URL.Query().Get("id")

	p := metaWOE(nId)

	json.NewEncoder(w).Encode(p)
}

func metaWOE(jobId string) models.MetaWork_order {

	db := dbConn()

	selDB, err := db.Query(`
	SELECT 
		W.id,
		J.name,
		J.color,
		W.description,
		W.planned_init,
		W.planned_end,
		M.nick_name,
		M.fname,
		M.lname,
		T.name,
		L.name,
		W.photo_before,
		W.photo_after,
		W.phase
	FROM
	    Work_Order W
			INNER JOIN
		Job_catalog J ON J.id = W.id_job
			INNER JOIN
		User_table M ON M.id = W.id_mecanic
			INNER JOIN
		Line L ON L.id = W.id_line
			INNER JOIN
		Transaction_state T ON T.id = W.id_state
	WHERE
		W.id = ?

	`, jobId)
	if err != nil {
		panic(err.Error())
	}

	p := models.MetaWork_order{}

	for selDB.Next() {
		var Id int
		var Job string
		var Color string
		var Description string
		var Init string
		var End string
		var NickName string
		var Fname string
		var Lname string
		var State string
		var Line string
		var Photo_before string
		var Photo_after string
		var Phase string

		err = selDB.Scan(&Id, &Job, &Color,
			&Description, &Init, &End, &NickName, &Fname,
			&Lname, &State, &Line, &Photo_before,
			&Photo_after, &Phase)
		if err != nil {
			panic(err.Error())
		}

		p.Id = Id
		p.Job = Job
		p.Color = Color
		p.Description = Description
		p.Init = Init
		p.End = End
		p.NickName = NickName
		p.Fname = Fname
		p.Lname = Lname
		p.State = State
		p.Line = Line
		p.Photo_before = Photo_before
		p.Photo_after = Photo_after
		p.Phase = Phase

	}
	defer db.Close()

	return p
}

func GetMetaWork_orderList(w http.ResponseWriter, r *http.Request) {

	init := r.URL.Query().Get("init")
	mech := r.URL.Query().Get("mech")

	p := metaWOEList(init, mech)

	json.NewEncoder(w).Encode(p)
}

func metaWOEList(init string, mech string) []models.MetaWork_order {

	db := dbConn()

	selDB, err := db.Query(`
	SELECT 
		W.id,
		J.name,
		J.color,
		W.description,
		W.planned_init,
		W.planned_end,
		M.nick_name,
		M.fname,
		M.lname,
		T.name,
		L.name,
		W.photo_before,
		W.photo_after,
		W.phase
	FROM
	    Work_Order W
			INNER JOIN
		Job_catalog J ON J.id = W.id_job
			INNER JOIN
		User_table M ON M.id = W.id_mecanic
			INNER JOIN
		Line L ON L.id = W.id_line
			INNER JOIN
		Transaction_state T ON T.id = W.id_state
	WHERE
		W.planned_init BETWEEN DATE_ADD(DATE(?),
		INTERVAL 6 HOUR) AND DATE_ADD(DATE(?),
		INTERVAL 1799 MINUTE)
		AND W.id_mecanic = ?

	`, init, init, mech)
	if err != nil {
		panic(err.Error())
	}

	p := models.MetaWork_order{}
	list := []models.MetaWork_order{}

	for selDB.Next() {
		var Id int
		var Job string
		var Color string
		var Description string
		var Init string
		var End string
		var NickName string
		var Fname string
		var Lname string
		var State string
		var Line string
		var Photo_before string
		var Photo_after string
		var Phase string

		err = selDB.Scan(&Id, &Job, &Color,
			&Description, &Init, &End, &NickName, &Fname,
			&Lname, &State, &Line, &Photo_before,
			&Photo_after, &Phase)
		if err != nil {
			panic(err.Error())
		}

		p.Id = Id
		p.Job = Job
		p.Color = Color
		p.Description = Description
		p.Init = Init
		p.End = End
		p.NickName = NickName
		p.Fname = Fname
		p.Lname = Lname
		p.State = State
		p.Line = Line
		p.Photo_before = Photo_before
		p.Photo_after = Photo_after
		p.Phase = Phase

		list = append(list, p)

	}
	defer db.Close()

	return list
}

func Mantenimiento(w http.ResponseWriter, r *http.Request) {

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

	tmpl.ExecuteTemplate(w, "Mantenimiento", user)
}

func ExecuteWorkOrder(w http.ResponseWriter, r *http.Request) {

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

	tmpl.ExecuteTemplate(w, "ExecuteWorkOrder", user)
}

func SetWorkOrderPhase(w http.ResponseWriter, r *http.Request) {
	status := 401
	state := "1"

	phaseName := []string{"notification_time", "wait_time", "diagnostic_time",
		"stock_time", "repair_time", "test_time", "delivery_time"}

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
		pass := r.FormValue("pass")
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
			UPDATE Work_Order 
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

			fmt.Printf("\njob: %s\t|\tphase: %s\t|\ttime: %s\n",
				uid, phase, utime)

			status = 200

		} else if uphase == "6" {
			grand = !AuthMecanic(user, pass, uid)
		} else {
			grand = AuthMecanic(user, pass, uid)
		}

		fmt.Printf("Granted %d\n", grand)

		if grand == true {

			if index > 1 && index < 5 {
				state = "2"
				fmt.Printf("stateFlag %s\n", state)
			} else if index > 5 {
				state = "3"
				index = 6
			}

			fmt.Println("index[%d]\tcurrentName: %s", index, phaseName[index])

			smt := `
			UPDATE Work_Order 
			SET 
				phase = ?,
				id_state= ?,
				` + phaseName[index] + `=?
			WHERE
				id = ?
			`

			insForm, err := db.Prepare(smt)

			if err != nil {
				panic(err.Error())
			}
			insForm.Exec(phase, state, utime, uid)

			fmt.Printf("job: %s\t|\tphase: %s\t|\ttime: %s\t|state: %s\n",
				uid, phase, utime, state)
			status = 200
		}

		defer db.Close()

	}
	fmt.Printf("Status:  %d\n", status)

	w.Header().Set("Server", "MES")
	w.WriteHeader(status)
}

func SetWorkOrderNote(w http.ResponseWriter, r *http.Request) {

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

		fmt.Printf("job: %s\t|\tNote: %s\n",
			uid, note)

		smt := `
		UPDATE Work_Order 
		SET 
			note=?
		WHERE
			id = ?
		`

		insForm, err := db.Prepare(smt)

		if err != nil {
			panic(err.Error())
		}
		insForm.Exec(note, uid)

		defer db.Close()
	}

	w.Header().Set("Server", "MES")
	w.WriteHeader(200)
}

func SetWorkOrderPhotoBefore(w http.ResponseWriter, r *http.Request) {

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
		filePath := FileHandlerEngine(r, "capture_before", "public/images/workOrders/", "before-*.png")

		fmt.Printf("job: %s\t|\timgPath: %s\n",
			uid, filePath)

		smt := `
		UPDATE Work_Order 
		SET 
			photo_before=?	
		WHERE
			id = ?
		`

		insForm, err := db.Prepare(smt)

		if err != nil {
			panic(err.Error())
		}
		insForm.Exec(filePath, uid)

		defer db.Close()
	}

	http.Redirect(w, r, "/executeWorkOrder", 301)
}

func SetWorkOrderPhotoAfter(w http.ResponseWriter, r *http.Request) {

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
		filePath := FileHandlerEngine(r, "capture_after", "public/images/workOrders/", "after-*.png")

		fmt.Printf("job: %s\t|\timgPath: %s\n",
			uid, filePath)

		smt := `
		UPDATE Work_Order 
		SET 
			photo_after=?
		WHERE
			id = ?
		`

		insForm, err := db.Prepare(smt)

		if err != nil {
			panic(err.Error())
		}
		insForm.Exec(filePath, uid)

		defer db.Close()
	}

	http.Redirect(w, r, "/executeWorkOrder", 301)
}

func Work_orderLog(w http.ResponseWriter, r *http.Request) {
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
	tmpl.ExecuteTemplate(w, "Work_orderLog", nil)
}

func GetMetaLogWork_order(w http.ResponseWriter, r *http.Request) {

	db := dbConn()

	iDate := r.URL.Query().Get("idate")
	fDate := r.URL.Query().Get("fdate")

	selDB, err := db.Query(`
	SELECT 
		W.id,
		J.name,
		J.color,
		W.description,
		W.planned_init,
		W.planned_end,
		M.profile_picture,
		M.fname,
		M.lname,
		T.name,
		T.color,
		L.name, 
		
		IFNULL(W.notification_time, "") as ActualBegin,
		
		IFNULL(timediff(W.wait_time, W.notification_time), "") as P1,
		IFNULL(timediff(W.diagnostic_time, W.wait_time), "") as P2,
		IFNULL(timediff(W.stock_time, W.diagnostic_time), "") as P3,
		IFNULL(timediff(W.repair_time, W.stock_time), "") as P4,
		IFNULL(timediff(W.test_time, W.repair_time), "") as P5,
		IFNULL(timediff(W.delivery_time, W.test_time),"")  as P6,
		
		IFNULL(W.test_time, "") as ActualEnd,

		W.phase,
		IFNULL(W.note, ""),
		W.photo_before,
		W.photo_after
	FROM
		Work_Order W
			INNER JOIN
		Job_catalog J ON J.id = W.id_job
			INNER JOIN
		User_table M ON M.id = W.id_mecanic
			INNER JOIN
		Line L ON L.id = W.id_line
			INNER JOIN
		Transaction_state T ON T.id = W.id_state
	WHERE
		DATE(W.planned_init) between ? and ?

	`, iDate, fDate)
	if err != nil {
		panic(err.Error())
	}

	p := models.MetaLogWork_order{}

	res := []models.MetaLogWork_order{}

	for selDB.Next() {
		var Id int
		var Job string
		var Color string
		var Description string
		var Init string
		var End string
		var Profile string
		var Fname string
		var Lname string
		var State string
		var SColor string
		var Line string

		var Wait_time string
		var Diagnostic_time string
		var Stock_time string
		var Repair_time string
		var Test_time string
		var Delivery_time string
		var Phase string
		var Note string
		var Photo_before string
		var Photo_after string

		var ActualBegin string
		var ActualEnd string

		err = selDB.Scan(&Id, &Job, &Color,
			&Description, &Init, &End, &Profile, &Fname,
			&Lname, &State, &SColor, &Line, &ActualBegin, &Wait_time,
			&Diagnostic_time, &Stock_time, &Repair_time, &Test_time,
			&Delivery_time, &ActualEnd, &Phase, &Note, &Photo_before, &Photo_after)
		if err != nil {
			panic(err.Error())
		}

		p.Id = Id
		p.Job = Job
		p.Color = Color
		p.Description = Description
		p.Init = Init
		p.End = End
		p.Profile = Profile
		p.Fname = Fname
		p.Lname = Lname
		p.State = State
		p.SColor = SColor

		p.Line = Line

		p.ActualBegin = ActualBegin

		p.Wait_time = Wait_time
		p.Diagnostic_time = Diagnostic_time
		p.Stock_time = Stock_time
		p.Repair_time = Repair_time
		p.Test_time = Test_time
		p.Delivery_time = Delivery_time

		p.ActualEnd = ActualEnd

		p.Phase = Phase
		p.Note = Note
		p.Photo_before = Photo_before
		p.Photo_after = Photo_after

		res = append(res, p)
	}
	defer db.Close()
	json.NewEncoder(w).Encode(res)
}

func isOperator(Id string) bool {

	auth := false
	level := ""

	db := dbConn()

	selDB, err := db.Query(`
		SELECT 
			P.description
		FROM
			User_table U
				INNER JOIN
			Privilege P ON P.id = U.id_privilege
		WHERE
			U.id = ?
		`, Id)

	if err != nil {
		panic(err.Error())
	}

	for selDB.Next() {

		err = selDB.Scan(&level)
		if err != nil {
			panic(err.Error())
		}
		fmt.Printf("Found!\n")
	}
	defer db.Close()

	if level == "1" {
		auth = true
	}

	return auth
}

func AuthMecanic(user, pass, jobId string) bool {
	auth := false
	var hash string
	var id int
	db := dbConn()

	p := metaWOE(jobId)

	if p.NickName == user {
		selDB, err := db.Query(`
		SELECT
			U.id,
			U.password
		FROM
			User_table U
		INNER JOIN Privilege P ON
			U.id_privilege = P.id
		WHERE
			U.nick_name = ?
		`, user)

		if err != nil {
			panic(err.Error())
		}

		for selDB.Next() {

			err = selDB.Scan(&id, &hash)
			if err != nil {
				panic(err.Error())
			}
			fmt.Printf("Found!\n")
		}
		defer db.Close()

		fmt.Printf("[%s] %s == %s | pass:  %s==%s\n\n",
			jobId, user, p.NickName, pass, hash)

		if !CheckPasswordHash(pass, hash) {
			fmt.Printf("Match!\n")

			if err != nil {
				auth = false
			}

		} else {
			auth = true
		}

	}

	return auth
}

func AuthAdmin(user, pass string) bool {
	auth := false
	var hash string
	var id int
	db := dbConn()

	selDB, err := db.Query(`
		SELECT
			U.id,
			U.password
		FROM
			User_table U
		INNER JOIN Privilege P ON
			U.id_privilege = P.id
		WHERE
			U.nick_name = ?
		`, user)

	if err != nil {
		panic(err.Error())
	}

	for selDB.Next() {

		err = selDB.Scan(&id, &hash)
		if err != nil {
			panic(err.Error())
		}
		fmt.Printf("Found!\n")
	}
	defer db.Close()

	if !CheckPasswordHash(pass, hash) {
		fmt.Printf("Match!\n")

		if err != nil {
			auth = false
		}

	} else {
		auth = true
	}

	fmt.Printf("Auth: ", auth)

	return auth
}

func GetUserByNickName(user string) string {

	var id string
	db := dbConn()

	selDB, err := db.Query(`
	SELECT 
		id
	FROM
		User_table
	WHERE
		nick_name = ?
	`, user)

	if err != nil {
		panic(err.Error())
	}

	for selDB.Next() {

		err = selDB.Scan(&id)
		if err != nil {
			panic(err.Error())
		}
	}
	defer db.Close()

	return id
}

func calcPlannedTime(mecanic int, day string) float32 {

	db := dbConn()

	selDB, err := db.Query(`

	SELECT 
		IFNULL(TIMEDIFF(W.planned_end, W.planned_init),
				0) AS saturation
	FROM
		mes.Work_Order W
	WHERE
		W.id_mecanic = ?
		AND DATE(W.planned_init) = ?
	`, mecanic, day)

	if err != nil {
		panic(err.Error())
	}

	var timeRes []string

	for selDB.Next() {
		var t string

		err = selDB.Scan(&t)

		if err != nil {
			panic(err.Error())
		}

		timeRes = append(timeRes, t)

		fmt.Printf("[%d]t: %s\n", mecanic, t)
	}
	defer db.Close()

	var workTime float32

	for _, t := range timeRes {

		var splitTime = strings.Split(t, ":")

		h, _ := strconv.ParseFloat(splitTime[0], 32)
		m, _ := strconv.ParseFloat(splitTime[1], 32)
		s, _ := strconv.ParseFloat(splitTime[2], 32)

		hour := float32(h)
		minute := float32(m)
		second := float32(s)

		workTime += hour + minute/60 + second/3600
	}

	return workTime
}

func CalcPlannedWorkSaturation(year, week string) []models.WorkSaturation {

	weekNum, _ := strconv.Atoi(week)
	yearNum, _ := strconv.Atoi(year)

	day := WeekStart(yearNum, weekNum)

	Data := []models.WorkSaturation{}

	Mecanics := GetMecanicsObj()

	for _, m := range Mecanics {

		fmt.Printf("mecanic: %s\n", m.Fname)

		saturation := models.WorkSaturation{}

		saturation.Mecanic = m
		sum := float32(0)

		for index := 0; index < 7; index++ {

			begin := day.String()
			t := calcPlannedTime(m.Id, begin)

			day = day.AddDate(0, 0, 1)
			saturation.Times = append(saturation.Times, t)

			sum += t
		}

		day = WeekStart(yearNum, weekNum)

		saturation.Times = append(saturation.Times, sum)

		Data = append(Data, saturation)
	}

	return Data
}

func GetPlannedWorkSaturarion(w http.ResponseWriter, r *http.Request) {

	year := r.URL.Query().Get("year")
	week := r.URL.Query().Get("week")

	saturation := CalcPlannedWorkSaturation(year, week)

	json.NewEncoder(w).Encode(saturation)
}

//Actual work order ejecution saturation time

func calcActualTime(mecanic int, day string) float32 {

	db := dbConn()

	selDB, err := db.Query(`

	SELECT 
		IFNULL(TIMEDIFF(W.test_time, W.notification_time),
				0) AS saturation
	FROM
		mes.Work_Order W
	WHERE
		W.id_mecanic = ?
		AND DATE(W.planned_init) = ?
	`, mecanic, day)

	if err != nil {
		panic(err.Error())

	}

	var timeRes []string

	for selDB.Next() {
		var t string

		err = selDB.Scan(&t)

		if err != nil {
			panic(err.Error())
		}

		timeRes = append(timeRes, t)
	}
	defer db.Close()

	var workTime float32

	for _, t := range timeRes {

		var splitTime = strings.Split(t, ":")

		if len(splitTime) == 1 {
			h, _ := strconv.ParseFloat(splitTime[0], 32)
			workTime += float32(h)
		} else {
			h, _ := strconv.ParseFloat(splitTime[0], 32)

			m, _ := strconv.ParseFloat(splitTime[1], 32)
			s, _ := strconv.ParseFloat(splitTime[2], 32)

			hour := float32(h)
			minute := float32(m)
			second := float32(s)

			workTime += hour + minute/60 + second/3600
		}
	}

	return workTime
}

func CalcActualWorkSaturation(year, week string) []models.WorkSaturation {

	weekNum, _ := strconv.Atoi(week)
	yearNum, _ := strconv.Atoi(year)

	day := WeekStart(yearNum, weekNum)

	Data := []models.WorkSaturation{}

	Mecanics := GetMecanicsObj()

	for _, m := range Mecanics {

		saturation := models.WorkSaturation{}

		saturation.Mecanic = m

		for index := 0; index < 7; index++ {

			begin := day.String()
			t := calcActualTime(m.Id, begin)

			day = day.AddDate(0, 0, 1)
			saturation.Times = append(saturation.Times, t)
		}

		sum := float32(0)

		for _, d := range saturation.Times {
			sum += d
		}

		saturation.Times = append(saturation.Times, sum)

		Data = append(Data, saturation)

		day = WeekStart(yearNum, weekNum)
	}

	return Data
}

func GetActualWorkSaturarion(w http.ResponseWriter, r *http.Request) {

	year := r.URL.Query().Get("year")
	week := r.URL.Query().Get("week")

	saturation := CalcActualWorkSaturation(year, week)

	json.NewEncoder(w).Encode(saturation)
}

func WorkOrderSaturation(w http.ResponseWriter, r *http.Request) {
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
	tmpl.ExecuteTemplate(w, "WorkOrderSaturation", user)
}

func DeleteWorkOrder(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	p := r.URL.Query().Get("id")
	delForm, err := db.Prepare("DELETE FROM Work_Order WHERE id=?")
	if err != nil {
		panic(err.Error())
	}
	delForm.Exec(p)
	log.Println("DELETE Work_Order")
	defer db.Close()
}

func WorkOrderCounter(line string, state string) int {
	db := dbConn()

	selDB, err := db.Query(`
	SELECT 
		COUNT(id)
	FROM
		Work_Order
	WHERE
			id_state = ? 
		AND id_line = ?

		`, state, line)

	if err != nil {
		panic(err.Error())
	}

	var res int

	for selDB.Next() {

		err = selDB.Scan(&res)

		if err != nil {
			panic(err.Error())
		}

	}
	defer db.Close()

	return res
}
