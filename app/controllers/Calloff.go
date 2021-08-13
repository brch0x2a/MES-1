package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	//"../models"
	"github.com/AdrianRb95/MES/app/models"

	"github.com/gorilla/websocket"
)

func RequestCalloff(w http.ResponseWriter, r *http.Request) {

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

	db := dbConn()

	selDB, err := db.Query("Select * from Material")
	if err != nil {
		panic(err.Error())
	}

	p := models.Material{}

	res := []models.Material{}

	for selDB.Next() {
		var Id int
		var Cod_material int
		var Material_name string

		err = selDB.Scan(&Id, &Cod_material, &Material_name)
		if err != nil {
			panic(err.Error())
		}
		p.Id = Id
		p.Cod_material = Cod_material
		p.Material_name = Material_name
		res = append(res, p)
	}

	defer db.Close()

	data := map[string]interface{}{
		"User":     user,
		"Material": res,
	}

	tmpl.ExecuteTemplate(w, "RequestCalloff", data)
}

func SelectRegisterBodega(w http.ResponseWriter, r *http.Request) {

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

	tmpl.ExecuteTemplate(w, "SelectRegisterBodega", user)
}

func InsertRequestCalloff(w http.ResponseWriter, r *http.Request) {

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

	db := dbConn()
	if r.Method == "POST" {
		id_line := user.Line
		operator := user.UserId
		id_material := r.FormValue("material")
		amount := r.FormValue("amount")
		comment := r.FormValue("comment")

		insForm, err := db.Prepare("insert into Call_off(id_line, id_operator, id_material, amount, comment) values (?, ?, ?, ?, ?)")
		if err != nil {
			panic(err.Error())
		}
		insForm.Exec(id_line, operator, id_material, amount, comment)
	}
	defer db.Close()

	http.Redirect(w, r, "/requestCalloff", 301)
}

func GestionCalloff(w http.ResponseWriter, r *http.Request) {

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

	tmpl.ExecuteTemplate(w, "GestionCalloff", user)
}

/*Obtiene los registros de call off ordenados por fecha reciente*/
func GetCalloffV00(w http.ResponseWriter, r *http.Request) {
	db := dbConn()

	selDB, err := db.Query(`
	SELECT 
		U.profile_picture,
		C.id,
		L.name,
		U.nick_name,
		U.fname,
		U.lname,
		C.request_date,
		C.close_date,
		M.cod_material,
		M.material_name,
		C.amount,
		C.comment,
		T.name,
		T.color
	FROM
		Call_off C
			INNER JOIN
		Line L ON L.id = C.id_line
			INNER JOIN
		User_table U ON U.id = C.id_operator
			INNER JOIN
		Material M ON M.id = C.id_material
			INNER JOIN
		Transaction_state T ON T.id = C.id_state
        where DATE(C.request_date) = DATE(CURRENT_TIMESTAMP)
	ORDER BY C.request_date DESC
    `)

	if err != nil {
		panic(err.Error())
	}

	p := models.Calloff_holder{}

	res := []models.Calloff_holder{}
	for selDB.Next() {
		var Profile_picture string
		var Id int
		var Linea string
		var NickName string
		var Fname string
		var Lname string
		var RequestDate string
		var CloseDate sql.NullString //para poder hacer el query de un null sin que se caiga por null
		var CodMaterial int
		var MaterialName string
		var Amount int
		var Comment string
		var State string
		var StateColor string

		err = selDB.Scan(&Profile_picture,
			&Id, &Linea, &NickName, &Fname, &Lname,
			&RequestDate, &CloseDate, &CodMaterial,
			&MaterialName, &Amount, &Comment, &State, &StateColor)

		if err != nil {
			panic(err.Error())
		}
		p.Profile_picture = Profile_picture
		p.Id = Id
		p.Linea = Linea
		p.NickName = NickName
		p.Fname = Fname
		p.Lname = Lname
		p.RequestDate = RequestDate
		p.CloseDate = CloseDate
		p.CodMaterial = CodMaterial
		p.MaterialName = MaterialName
		p.Amount = Amount
		p.Comment = Comment
		p.State = State
		p.StateColor = StateColor

		res = append(res, p)
	}
	defer db.Close()
	json.NewEncoder(w).Encode(res)
}

/*
	Obtiene los registros de call off ordenados por fecha reciente y filtrados por rango
	de fechas y id_usuario y linea, esto para ver solo las transacciones por personas, es decir
	las solicitudes indiviales
*/
func GetCalloffV01(w http.ResponseWriter, r *http.Request) {

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
	iDate := r.URL.Query().Get("idate")
	fDate := r.URL.Query().Get("fdate")

	db := dbConn()

	selDB, err := db.Query(`
	SELECT 
		U.profile_picture,
		C.id,
		L.name,
		U.nick_name,
		U.fname,
		U.lname,
		C.request_date,
		C.close_date,
		M.cod_material,
		M.material_name,
		C.amount,
		C.comment,
		T.name,
		T.color
	FROM
		Call_off C
			INNER JOIN
		Line L ON L.id = C.id_line
			INNER JOIN
		User_table U ON U.id = C.id_operator
			INNER JOIN
		Material M ON M.id = C.id_material
			INNER JOIN
		Transaction_state T ON T.id = C.id_state
	WHERE
		C.id_operator = ? AND C.id_line = ?
			AND ? <= DATE(C.request_date)
			AND DATE(C.request_date) <= ?
	ORDER BY request_date DESC

    `, user.UserId, user.Line, iDate, fDate)

	if err != nil {
		panic(err.Error())
	}

	p := models.Calloff_holder{}

	res := []models.Calloff_holder{}
	for selDB.Next() {
		var Profile_picture string

		var Id int
		var Linea string
		var NickName string
		var Fname string
		var Lname string
		var RequestDate string
		var CloseDate sql.NullString //para poder hacer el query de un null sin que se caiga por null
		var CodMaterial int
		var MaterialName string
		var Amount int
		var Comment string
		var State string
		var StateColor string

		err = selDB.Scan(&Profile_picture,
			&Id, &Linea, &NickName, &Fname, &Lname,
			&RequestDate, &CloseDate, &CodMaterial,
			&MaterialName, &Amount, &Comment, &State, &StateColor)

		if err != nil {
			panic(err.Error())
		}
		p.Profile_picture = Profile_picture
		p.Id = Id
		p.Linea = Linea
		p.NickName = NickName
		p.Fname = Fname
		p.Lname = Lname
		p.RequestDate = RequestDate
		p.CloseDate = CloseDate
		p.CodMaterial = CodMaterial
		p.MaterialName = MaterialName
		p.Amount = Amount
		p.Comment = Comment
		p.State = State
		p.StateColor = StateColor

		res = append(res, p)
	}
	defer db.Close()
	json.NewEncoder(w).Encode(res)
}

/*
	Obtiene los registros de call off ordenados por fecha reciente y filtrados por rango
	de fechas
*/
func GetCalloffV02(w http.ResponseWriter, r *http.Request) {

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
	iDate := r.URL.Query().Get("idate")
	fDate := r.URL.Query().Get("fdate")

	db := dbConn()

	selDB, err := db.Query(`
	SELECT 
		U.profile_picture,
		C.id,
		L.name,
		U.nick_name,
		U.fname,
		U.lname,
		C.request_date,
		C.close_date,
		M.cod_material,
		M.material_name,
		C.amount,
		C.comment,
		T.name,
		T.color
	FROM
		Call_off C
			INNER JOIN
		Line L ON L.id = C.id_line
			INNER JOIN
		User_table U ON U.id = C.id_operator
			INNER JOIN
		Material M ON M.id = C.id_material
			INNER JOIN
		Transaction_state T ON T.id = C.id_state
	WHERE
			? <= DATE(C.request_date)
			AND DATE(C.request_date) <= ?
	ORDER BY request_date DESC

    `, iDate, fDate)

	if err != nil {
		panic(err.Error())
	}

	p := models.Calloff_holder{}

	res := []models.Calloff_holder{}
	for selDB.Next() {
		var Profile_picture string

		var Id int
		var Linea string
		var NickName string
		var Fname string
		var Lname string
		var RequestDate string
		var CloseDate sql.NullString //para poder hacer el query de un null sin que se caiga por null
		var CodMaterial int
		var MaterialName string
		var Amount int
		var Comment string
		var State string
		var StateColor string

		err = selDB.Scan(&Profile_picture,
			&Id, &Linea, &NickName, &Fname, &Lname,
			&RequestDate, &CloseDate, &CodMaterial,
			&MaterialName, &Amount, &Comment, &State, &StateColor)

		if err != nil {
			panic(err.Error())
		}
		p.Profile_picture = Profile_picture
		p.Id = Id
		p.Linea = Linea
		p.NickName = NickName
		p.Fname = Fname
		p.Lname = Lname
		p.RequestDate = RequestDate
		p.CloseDate = CloseDate
		p.CodMaterial = CodMaterial
		p.MaterialName = MaterialName
		p.Amount = Amount
		p.Comment = Comment
		p.State = State
		p.StateColor = StateColor

		res = append(res, p)
	}
	defer db.Close()
	json.NewEncoder(w).Encode(res)
}
func GetCalloffEV00(w http.ResponseWriter, r *http.Request) {
	uid := r.URL.Query().Get("uid")

	db := dbConn()

	selDB, err := db.Query(`
	SELECT 
		U.profile_picture,
		C.id,
		L.name,
		U.nick_name,
		U.fname,
		U.lname,
		C.request_date,
		C.close_date,
		M.cod_material,
		M.material_name,
		C.amount,
		C.comment,
		T.name,
		T.color
	FROM
		Call_off C
			INNER JOIN
		Line L ON L.id = C.id_line
			INNER JOIN
		User_table U ON U.id = C.id_operator
			INNER JOIN
		Material M ON M.id = C.id_material
			INNER JOIN
		Transaction_state T ON T.id = C.id_state
	where C.id = ?
    `, uid)

	if err != nil {
		panic(err.Error())
	}

	p := models.Calloff_holder{}

	for selDB.Next() {
		var Profile_picture string

		var Id int
		var Linea string
		var NickName string
		var Fname string
		var Lname string
		var RequestDate string
		var CloseDate sql.NullString //para poder hacer el query de un null sin que se caiga por null
		var CodMaterial int
		var MaterialName string
		var Amount int
		var Comment string
		var State string
		var StateColor string

		err = selDB.Scan(&Profile_picture,
			&Id, &Linea, &NickName, &Fname, &Lname,
			&RequestDate, &CloseDate, &CodMaterial,
			&MaterialName, &Amount, &Comment, &State, &StateColor)

		if err != nil {
			panic(err.Error())
		}
		p.Profile_picture = Profile_picture
		p.Id = Id
		p.Linea = Linea
		p.NickName = NickName
		p.Fname = Fname
		p.Lname = Lname
		p.RequestDate = RequestDate
		p.CloseDate = CloseDate
		p.CodMaterial = CodMaterial
		p.MaterialName = MaterialName
		p.Amount = Amount
		p.Comment = Comment
		p.State = State
		p.StateColor = StateColor
	}
	defer db.Close()
	json.NewEncoder(w).Encode(p)
}

/*----------------------------------*/
func WsEndpointCalloff(w http.ResponseWriter, r *http.Request) {
	//fmt.Fprintf(w, "Hello World")
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}

	//clientResult, _ := readerFilter(ws)

	go WriterByCalloff(ws)
}

func WriterByCalloff(conn *websocket.Conn) {

	for {
		fmt.Printf("\n\n----------Calloff--Monitor----------------\n")

		ticker := time.NewTicker(time.Second)

		if err := conn.WriteMessage(websocket.TextMessage, []byte("")); err != nil {
			fmt.Println(err)
			return
		}

		for t := range ticker.C {

			msg, err_msg := readerFilter(conn)

			if err_msg != nil {
				fmt.Printf("Client socket close!")
				conn.Close()
				break
			}

			if msg == "on" {
				fmt.Printf("Updating At: %+v\n", t)

				//data = rand.Intn(100)

				response, err := http.Get("http://localhost:3000/getCalloffV00")

				if err != nil {
					fmt.Printf("The HTTP request failed with error %s\n", err)
					break
				}

				data, _ := ioutil.ReadAll(response.Body)
				fmt.Println("data")
				if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
					fmt.Println(err)
					break
				}
			} else {
				fmt.Printf("Client socket close!")
				conn.Close()
				break
			}

		}
		defer conn.Close()
	}
}

func UpdateCalloff(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	if r.Method == "POST" {
		State := r.FormValue("state")
		Uid := r.FormValue("uid")
		smt := ""

		switch State {
		case "1":
			smt = "UPDATE Call_off SET id_state=?, close_date=null WHERE id=?"
		case "2":
			smt = "UPDATE Call_off SET id_state=?, close_date=null WHERE id=?"
		case "3":
			smt = "UPDATE Call_off SET id_state=?, close_date=current_timestamp WHERE id=?"
		}

		insForm, err := db.Prepare(smt)

		if err != nil {
			panic(err.Error())
		}

		insForm.Exec(State, Uid)
		log.Println("UPDATE Calloff: State: " + State + " | Uid: " + Uid)
	}
	defer db.Close()
}

func ConsolidatedCalloff(w http.ResponseWriter, r *http.Request) {

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

	tmpl.ExecuteTemplate(w, "ConsolidatedCalloff", user)
}

func GestionSolicitudCalloff(w http.ResponseWriter, r *http.Request) {

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

	tmpl.ExecuteTemplate(w, "GestionSolicitudCalloff", user)
}
