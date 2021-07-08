package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"../models"
)

func CleanDisinfection(w http.ResponseWriter, r *http.Request) {

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

	db := dbConn()

	selDB, err := db.Query("Select * from Equipment")
	if err != nil {
		panic(err.Error())
	}

	p := models.Equipment{}

	res := []models.Equipment{}

	for selDB.Next() {
		var Id int
		var Name string

		err = selDB.Scan(&Id, &Name)
		if err != nil {
			panic(err.Error())
		}
		p.Id = Id
		p.Name = Name

		res = append(res, p)
	}

	defer db.Close()
	/*_______________________________________________________*/

	/*_______________________________________________________*/

	data := map[string]interface{}{
		"User":   user,
		"Equipo": res,
	}

	tmpl.ExecuteTemplate(w, "CleanDisinfection", data)
}

func GetCleanFilter(w http.ResponseWriter, r *http.Request) {
	db := dbConn()

	selDB, err := db.Query(`
        SELECT
		   id,
		   name
        FROM Clean_filter
    `)

	if err != nil {
		panic(err.Error())
	}

	res := []models.Clean_filter{}
	p := models.Clean_filter{}

	for selDB.Next() {
		var Id int
		var Name string

		err = selDB.Scan(&Id, &Name)

		if err != nil {
			panic(err.Error())
		}

		p.Id = Id
		p.Name = Name

		res = append(res, p)
	}
	defer db.Close()
	json.NewEncoder(w).Encode(res)
}

func GetClean_corrective_action(w http.ResponseWriter, r *http.Request) {
	db := dbConn()

	selDB, err := db.Query(`
        SELECT
		   id,
		   description,
		   bottom,
		   top
        FROM Clean_Corrective_Action
    `)

	if err != nil {
		panic(err.Error())
	}

	res := []models.Clean_corrective_action{}
	p := models.Clean_corrective_action{}

	for selDB.Next() {
		var Id int
		var Description string
		var Bottom float32
		var Top float32

		err = selDB.Scan(&Id, &Description, &Bottom, &Top)

		if err != nil {
			panic(err.Error())
		}

		p.Id = Id
		p.Description = Description
		p.Bottom = Bottom
		p.Top = Top

		res = append(res, p)
	}
	defer db.Close()
	json.NewEncoder(w).Encode(res)
}

func GetWashingStage(w http.ResponseWriter, r *http.Request) {
	db := dbConn()

	selDB, err := db.Query(`
        SELECT
		   id,
		   name
        FROM Washing_stage
    `)

	if err != nil {
		panic(err.Error())
	}
	res := []models.Washing_stage{}
	p := models.Washing_stage{}

	for selDB.Next() {
		var Id int
		var Name string

		err = selDB.Scan(&Id, &Name)

		if err != nil {
			panic(err.Error())
		}

		p.Id = Id
		p.Name = Name

		res = append(res, p)
	}
	defer db.Close()
	json.NewEncoder(w).Encode(res)
}

func InsertCleanDisinfection(w http.ResponseWriter, r *http.Request) {

	session, err := store.Get(r, "cookie-name")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	user := GetUser(session)

	if r.Method == "POST" {
		db := dbConn()

		date_init := r.FormValue("init")
		washing_stage := r.FormValue("stage")
		equipment := r.FormValue("equipment")

		washed_type := r.FormValue("washed")

		detergent := r.FormValue("detergent")
		disinfectand := r.FormValue("disinfectand")
		chemical := r.FormValue("chemical")
		foam := r.FormValue("foam")
		spray := r.FormValue("spray")

		filter := r.FormValue("filter")
		water_ph := r.FormValue("ph")

		visual_inspection := r.FormValue("visual")
		microbiology := r.FormValue("micro")

		atp := r.FormValue("atp")
		corrective_action := r.FormValue("corrective_action")
		new_atp := r.FormValue("new_atp")

		allergen_state := r.FormValue("allergen")
		maintenance := r.FormValue("maintenance")

		comment := r.FormValue("comment")

		// fmt.Printf("stage: %s\n", washing_stage)
		// fmt.Printf("equipment: %s\n", equipment)
		// fmt.Printf("washed_type: %s\n", washed_type)
		// fmt.Printf("detergent: %s\n", detergent)
		// fmt.Printf("disinfectand: %s\n", disinfectand)
		// fmt.Printf("chemical: %s\n", chemical)
		// fmt.Printf("foam: %s\n", foam)
		// fmt.Printf("spray: %s\n", spray)
		// fmt.Printf("filter: %s\n", filter)
		// fmt.Printf("water_ph: %s\n", water_ph)
		// fmt.Printf("visual_inspection: %s\n", visual_inspection)
		// fmt.Printf("microbiology: %s\n", microbiology)
		// fmt.Printf("atp: %s\n", atp)
		// fmt.Printf("corrective_action: %s\n", corrective_action)
		// fmt.Printf("new_atp: %s\n", new_atp)
		// fmt.Printf("aller_state: %s\n", allergen_state)
		// fmt.Printf("maintenance: %s\n", maintenance)
		// fmt.Printf("commnet: %s\n", comment)
		// fmt.Printf("User: %s\n", user.UserId)

		smt := `
		insert into Clean_Disinfection(
			id_line,
			date_init,
			id_washing_stage,
			id_equiment,
			whashed_type,
			detergent,
			disinfectand,
			chemical,
			foam,
			spray,
			id_filter,
			water_ph,
			visual_inspection,
			microbiology,
			atp,
			id_corrective_action,
			new_atp,
			id_allergen_state,
			maintenance,
			comment, 
			id_autor,
			id_approver,
			id_state
			)
		values(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		`

		insForm, err := db.Prepare(smt)

		if err != nil {
			panic(err.Error())
		}
		insForm.Exec(
			user.Line,
			date_init,
			washing_stage,
			equipment,
			washed_type,
			detergent,
			disinfectand,
			chemical,
			foam,
			spray,
			filter,
			water_ph,
			visual_inspection,
			microbiology,
			atp,
			corrective_action,
			new_atp,
			allergen_state,
			maintenance,
			comment,
			user.UserId,
			5,
			4)

		defer db.Close()
	}

	http.Redirect(w, r, "/cleanDisinfection", 301)

}
func ConsolidatedCleanDisinfection(w http.ResponseWriter, r *http.Request) {

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

	tmpl.ExecuteTemplate(w, "ConsolidatedCleanDisinfection", user)
}

func GetCleanDisinfectionBy(w http.ResponseWriter, r *http.Request) {

	init := r.URL.Query().Get("init")
	end := r.URL.Query().Get("end")
	line := r.URL.Query().Get("line")

	db := dbConn()

	selDB, err := db.Query(`
	SELECT 
		CD.id,
		R.profile_picture,
		CD.date_init,
		R.fname,
		R.lname,
		L.name,
		W.name,
		E.name,
		CD.whashed_type,
		CD.detergent,
		CD.disinfectand,
		CD.chemical,
		CD.foam,
		CD.spray,
		CF.name,
		CD.water_ph,
		IF(CD.visual_inspection, 'Si', 'No'),
		IF(CD.microbiology, 'Si', 'No'),
		CD.atp,
		CD.id_corrective_action,
		CD.new_atp,
		T.name,
		CD.maintenance,
		CD.comment,
		A.profile_picture,
		A.fname,
		A.lname,
		IFNULL(CD.date_end, ''),
		CT.name
	FROM
		Clean_Disinfection CD
			INNER JOIN
		User_table R ON R.id = CD.id_autor
			INNER JOIN
		Line L ON L.id = CD.id_line
			INNER JOIN
		Washing_stage W ON W.id = CD.id_washing_stage
			INNER JOIN
		Equipment E ON E.id = CD.id_equiment
			INNER JOIN
		Clean_filter CF ON CF.id = CD.id_filter
			INNER JOIN
		Transaction_state T ON T.id = CD.id_allergen_state
			INNER JOIN
		User_table A ON A.id = CD.id_approver
			INNER JOIN
		Transaction_state CT ON CT.id = CD.id_state
	WHERE
		CD.id_line = ?
        AND CD.date_init BETWEEN ? AND ?
	`, line, init, end)

	if err != nil {
		panic(err.Error())
	}

	p := models.Clean_DisinfectionMeta{}

	res := []models.Clean_DisinfectionMeta{}

	for selDB.Next() {
		var Id int

		var Date_init string

		var RProfile_Picture string
		var RFname string
		var RLname string

		var Line string

		var Washing_stage string

		var Equipment string

		var Washed_type string

		var Detergent float32
		var Disinfectand float32
		var Chemical float32
		var Foam float32
		var Spray float32

		var Filter string
		var Water_ph float32

		var Visual_inspection string
		var Microbiology string

		var Atp float32
		var Corrective_action int
		var New_atp float32

		var Allergen_state string

		var Maintenance string

		var Comment string

		var AProfile_Picture string
		var AFname string
		var ALname string

		var Date_end string

		var State string

		err = selDB.Scan(
			&Id,
			&RProfile_Picture,
			&Date_init,
			&RFname,
			&RLname,
			&Line,
			&Washing_stage,
			&Equipment,
			&Washed_type,
			&Detergent,
			&Disinfectand,
			&Chemical,
			&Foam,
			&Spray,
			&Filter,
			&Water_ph,
			&Visual_inspection,
			&Microbiology,
			&Atp,
			&Corrective_action,
			&New_atp,
			&Allergen_state,
			&Maintenance,
			&Comment,
			&AProfile_Picture,
			&AFname,
			&ALname,
			&Date_end,
			&State)

		if err != nil {
			panic(err.Error())
		}

		p.Id = Id
		p.Date_init = Date_init
		p.RProfile_Picture = RProfile_Picture
		p.RFname = RFname
		p.RLname = RLname
		p.Line = Line

		p.Washing_stage = Washing_stage
		p.Equipment = Equipment
		p.Washed_type = Washed_type

		p.Detergent = Detergent
		p.Disinfectand = Disinfectand
		p.Chemical = Chemical
		p.Foam = Foam
		p.Spray = Spray

		p.Filter = Filter
		p.Water_ph = Water_ph

		p.Visual_inspection = Visual_inspection
		p.Microbiology = Microbiology

		p.Atp = Atp
		p.Corrective_action = Corrective_action
		p.New_atp = New_atp

		p.Allergen_state = Allergen_state

		p.Maintenance = Maintenance

		p.Comment = Comment

		p.AProfile_Picture = AProfile_Picture
		p.AFname = AFname
		p.ALname = ALname

		p.Date_end = Date_end

		p.State = State

		res = append(res, p)

	}

	defer db.Close()

	json.NewEncoder(w).Encode(res)
}

func GetCleanDisinfectionE(w http.ResponseWriter, r *http.Request) {

	id := r.URL.Query().Get("id")

	db := dbConn()

	selDB, err := db.Query(`
	SELECT 
		CD.id,
		R.profile_picture,
		CD.date_init,
		R.fname,
		R.lname,
		L.name,
		W.name,
		E.name,
		CD.whashed_type,
		CD.detergent,
		CD.disinfectand,
		CD.chemical,
		CD.foam,
		CD.spray,
		CF.name,
		CD.water_ph,
		IF(CD.visual_inspection, 'Si', 'No'),
		IF(CD.microbiology, 'Si', 'No'),
		CD.atp,
		CD.id_corrective_action,
		CD.new_atp,
		T.name,
		CD.maintenance,
		CD.comment,
		A.profile_picture,
		A.fname,
		A.lname,
		IFNULL(CD.date_end, ''),
		CT.name
	FROM
		Clean_Disinfection CD
			INNER JOIN
		User_table R ON R.id = CD.id_autor
			INNER JOIN
		Line L ON L.id = CD.id_line
			INNER JOIN
		Washing_stage W ON W.id = CD.id_washing_stage
			INNER JOIN
		Equipment E ON E.id = CD.id_equiment
			INNER JOIN
		Clean_filter CF ON CF.id = CD.id_filter
			INNER JOIN
		Transaction_state T ON T.id = CD.id_allergen_state
			INNER JOIN
		User_table A ON A.id = CD.id_approver
			INNER JOIN
		Transaction_state CT ON CT.id = CD.id_state
	WHERE
		CD.id  = ?
	`, id)

	if err != nil {
		panic(err.Error())
	}

	p := models.Clean_DisinfectionMeta{}

	for selDB.Next() {
		var Id int

		var Date_init string

		var RProfile_Picture string
		var RFname string
		var RLname string

		var Line string

		var Washing_stage string

		var Equipment string

		var Washed_type string

		var Detergent float32
		var Disinfectand float32
		var Chemical float32
		var Foam float32
		var Spray float32

		var Filter string
		var Water_ph float32

		var Visual_inspection string
		var Microbiology string

		var Atp float32
		var Corrective_action int
		var New_atp float32

		var Allergen_state string

		var Maintenance string

		var Comment string

		var AProfile_Picture string
		var AFname string
		var ALname string

		var Date_end string

		var State string

		err = selDB.Scan(
			&Id,
			&RProfile_Picture,
			&Date_init,
			&RFname,
			&RLname,
			&Line,
			&Washing_stage,
			&Equipment,
			&Washed_type,
			&Detergent,
			&Disinfectand,
			&Chemical,
			&Foam,
			&Spray,
			&Filter,
			&Water_ph,
			&Visual_inspection,
			&Microbiology,
			&Atp,
			&Corrective_action,
			&New_atp,
			&Allergen_state,
			&Maintenance,
			&Comment,
			&AProfile_Picture,
			&AFname,
			&ALname,
			&Date_end,
			&State)

		if err != nil {
			panic(err.Error())
		}

		p.Id = Id
		p.Date_init = Date_init
		p.RProfile_Picture = RProfile_Picture
		p.RFname = RFname
		p.RLname = RLname
		p.Line = Line

		p.Washing_stage = Washing_stage
		p.Equipment = Equipment
		p.Washed_type = Washed_type

		p.Detergent = Detergent
		p.Disinfectand = Disinfectand
		p.Chemical = Chemical
		p.Foam = Foam
		p.Spray = Spray

		p.Filter = Filter
		p.Water_ph = Water_ph

		p.Visual_inspection = Visual_inspection
		p.Microbiology = Microbiology

		p.Atp = Atp
		p.Corrective_action = Corrective_action
		p.New_atp = New_atp

		p.Allergen_state = Allergen_state

		p.Maintenance = Maintenance

		p.Comment = Comment

		p.AProfile_Picture = AProfile_Picture
		p.AFname = AFname
		p.ALname = ALname

		p.Date_end = Date_end

		p.State = State

	}

	defer db.Close()

	json.NewEncoder(w).Encode(p)
}

func SetCleanState(w http.ResponseWriter, r *http.Request) {
	status := 401
	db := dbConn()
	if r.Method == "POST" {

		id := r.FormValue("id")
		user := r.FormValue("user")
		pass := r.FormValue("pass")
		state := r.FormValue("state")

		now := time.Now()

		date_end := now.Format("2006-01-02 15:04:05")

		var grand bool

		grand = AuthAdmin(user, pass)

		fmt.Println("\n\n____________Id: %s", id)
		fmt.Println("User: %s", user)
		fmt.Println("Pass: %s", pass)
		fmt.Println("Date: %s", date_end)
		fmt.Println("State: %s", state)

		if grand {

			userId := GetUserByNickName(user)
			fmt.Println("Flag grand:%d", grand)

			smt := `
				UPDATE Clean_Disinfection 
				SET 
					id_approver= ?,
					id_state = ?,
					date_end= ?
				WHERE
					id = ?
				`
			insForm, err := db.Prepare(smt)
			if err != nil {
				panic(err.Error())
			}

			fmt.Printf("approver:%s\tstate:%s\tdate: %s\tjob:%s\n", userId, state, date_end, id)

			insForm.Exec(userId, state, date_end, id)

			status = 200
		}

	}
	defer db.Close()

	w.Header().Set("Server", "MES")
	w.WriteHeader(status)
}
