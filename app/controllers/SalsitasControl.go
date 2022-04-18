package controllers

import (
	//"bytes"
	//"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"

	//	"strings"
	"time"
	//"image/jpeg"
	"encoding/json"
	"io/ioutil"

	//"../models"

	"github.com/AdrianRb95/MES/app/models"

	"github.com/gorilla/websocket"
	//_ "github.com/go-sql-driver/mysql"
)

func FillSalsitasControl(w http.ResponseWriter, r *http.Request) {
	db := dbConn()

	selDB, err := db.Query("Select * from Header where id=2")

	if err != nil {
		panic(err.Error())
	}

	p := models.Header{}

	for selDB.Next() {
		var Id int
		var Name string
		var Cod_doc string
		var Revision_date string
		var Next_revision_date string
		var Revision_no int

		err = selDB.Scan(&Id, &Name, &Cod_doc, &Revision_date, &Next_revision_date, &Revision_no)

		if err != nil {
			panic(err.Error())
		}

		p.Id = Id
		p.Name = Name
		p.Cod_doc = Cod_doc
		p.Revision_date = Revision_date
		p.Next_revision_date = Next_revision_date
		p.Revision_no = Revision_no

	}

	/*-----------------Auth--------------------------*/
	session, err := store.Get(r, "cookie-name")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	user := GetUser(session)

	if auth := user.Authenticated; !auth {
		session.AddFlash("Acesso denegado!")
		err = session.Save(r, w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/forbidden", http.StatusFound)
		return
	}

	data := map[string]interface{}{
		"Header": p,
		"User":   user,
	}
	defer db.Close()
	/*-------------------------------------------*/
	tmpl.ExecuteTemplate(w, "FillSalsitasControl", data)
}

//Metodo para llenar el sub encabeado de pesos de 	CALREG022
func InsertSalsitasControl(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "cookie-name")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	user := GetUser(session)

	SubHeader := models.SalsitasControl_subheader_holder{}
	var LastId int

	db := dbConn()
	if r.Method == "POST" {

		date := r.FormValue("vdate")
		turn := r.FormValue("turn")
		line := r.FormValue("line")
		presentation := r.FormValue("presentation")
		coordinador := r.FormValue("coordinador")
		operator := user.UserId
		header := r.FormValue("header")

		insForm, err := db.Prepare(`
		INSERT INTO Salsitas_statistical_process(
			date_reg,
			turn,
			id_line,
			id_presentation,
			id_coordinator,
			id_operator,
			id_header
		)
		VALUES(
			?,
			?,
			?,
			?,
			?,
			?,
			?
		)
		
		`)

		if err != nil {
			panic(err.Error())
		}

		res, err := insForm.Exec(date,
			turn, line, presentation, coordinador,
			operator, header)
		if err != nil {
			log.Fatal("Cannot run insert statement", err)
		}

		id, _ := res.LastInsertId()

		fmt.Printf("Inserted row: %d", id)

		LastId = int(id)
	}
	defer db.Close()

	fmt.Printf("Last Id: %d", LastId)

	selDB_sub, err := db.Query(`
				SELECT
				SC.id,
				SC.date_reg,
				SC.id_header,
				SC.turn,
				L.name,
				P.name,
				P.weight_value,
				P.weight_unit,
				P.error_rate,
				UC.fname,
				UO.fname
			FROM
				Salsitas_statistical_process SC
			INNER JOIN Line L ON
				SC.id_line = L.id
			INNER JOIN Presentation P ON
				P.id = SC.id_presentation
			INNER JOIN User_table UC ON
				UC.id = SC.id_coordinator
			INNER JOIN User_table UO ON
				SC.id_operator = UO.id
			where SC.id = ?
	`, LastId)

	if err != nil {
		panic(err.Error())
	}

	for selDB_sub.Next() {
		var Id int
		var Date string
		var Turn int
		var Line string
		var Presentation string
		var Pvalue float32
		var Punit string
		var Perror float32
		var Coordinator string
		var Operator string
		var Header int

		err = selDB_sub.Scan(&Id, &Date, &Header, &Turn,
			&Line, &Presentation, &Pvalue, &Punit, &Perror,
			&Coordinator, &Operator)

		if err != nil {
			panic(err.Error())
		}

		SubHeader.Id = Id
		SubHeader.Date = Date
		SubHeader.Turn = Turn
		SubHeader.Line = Line
		SubHeader.Presentation = Presentation
		SubHeader.Pvalue = Pvalue
		SubHeader.Punit = Punit
		SubHeader.Perror = Perror //rate
		SubHeader.Coordinator = Coordinator
		SubHeader.Operator = Operator
		SubHeader.Header = Header

	}
	defer db.Close()
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
		"SubHeader": SubHeader,
		"User":      user,
	}

	tmpl.ExecuteTemplate(w, "FillSalsitasWeight", data)
}

func WeightSubControlStep(w http.ResponseWriter, r *http.Request) {

	session, err := store.Get(r, "cookie-name")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	user := GetUser(session)

	SubHeader := models.SalsitasControl_subheader_holder{}

	db := dbConn()

	selDB_sub, err := db.Query(`
				SELECT
				SC.id,
				SC.date_reg,
				SC.id_header,
				SC.turn,
				L.name,
				P.name,
				P.weight_value,
				P.weight_unit,
				P.error_rate,
				UC.fname,
				UO.fname
			FROM
				Salsitas_statistical_process SC
			INNER JOIN Line L ON
				SC.id_line = L.id
			INNER JOIN Presentation P ON
				P.id = SC.id_presentation
			INNER JOIN User_table UC ON
				UC.id = SC.id_coordinator
			INNER JOIN User_table UO ON
				SC.id_operator = UO.id
			where SC.id = ?
	`, user.WeightSubId)

	if err != nil {
		panic(err.Error())
	}

	for selDB_sub.Next() {
		var Id int
		var Date string
		var Turn int
		var Line string
		var Presentation string
		var Pvalue float32
		var Punit string
		var Perror float32
		var Coordinator string
		var Operator string
		var Header int

		err = selDB_sub.Scan(&Id, &Date, &Header, &Turn,
			&Line, &Presentation, &Pvalue, &Punit, &Perror,
			&Coordinator, &Operator)

		if err != nil {
			panic(err.Error())
		}

		SubHeader.Id = Id
		SubHeader.Date = Date
		SubHeader.Turn = Turn
		SubHeader.Line = Line
		SubHeader.Presentation = Presentation
		SubHeader.Pvalue = Pvalue
		SubHeader.Punit = Punit
		SubHeader.Perror = Perror //rate
		SubHeader.Coordinator = Coordinator
		SubHeader.Operator = Operator
		SubHeader.Header = Header

	}

	session.Values["weigthId"] = SubHeader.Id

	if auth := user.Authenticated; !auth {
		session.AddFlash("You don't have access!")
		http.Redirect(w, r, "/forbidden", http.StatusFound)
		return
	}

	err = session.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"SubHeader": SubHeader,
		"User":      user,
	}
	defer db.Close()
	tmpl.ExecuteTemplate(w, "FillSalsitasWeight", data)
}

// Sistemas IQ

func CalculaPesoNetoSalsitas(w http.ResponseWriter, r *http.Request) {
	SubHeader := models.SalsitasControl_subheader_holder{}

	//	var SubId int
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

	nLine := r.URL.Query().Get("line")
	dateInit := r.URL.Query().Get("dinit")
	dateFinal := r.URL.Query().Get("dfinal")

	selDB, err := db.Query(`
	SELECT
		Op.profile_picture,
		W.date_reg,
		S.turn,
		L.name,
		P.name,
		P.weight_value,
		P.weight_unit,
		P.error_rate,
		Co.fname,
		Co.lname,
		Op.fname,
		Op.lname,
		W.value1,
		W.value2,
		W.value3,
		W.value4,
		W.value5
	FROM
		Weight_control W
	INNER JOIN Salsitas_statistical_process S ON
		S.id = W.id_sub_header
	INNER JOIN Line L ON
		L.id = S.id_line
	INNER JOIN Presentation P ON
		P.id = S.id_presentation
	INNER JOIN User_table Co ON
		Co.id = S.id_coordinator
	INNER JOIN User_table Op ON
		Op.id = S.id_operator
	WHERE
		L.id = ? AND
		W.date_reg >= ? AND 
		W.date_reg < ?`, nLine, dateInit, dateFinal)

	if err != nil {
		panic(err.Error())
	}

	for selDB.Next() {
		var Id int
		var Date string
		var Turn int
		var Line string
		var Presentation string
		var Pvalue float32
		var Punit string
		var Perror float32
		var Coordinator string
		var Operator string
		var Header int

		err = selDB.Scan(&Id, &Date, &Header, &Turn,
			&Line, &Presentation, &Pvalue, &Punit, &Perror,
			&Coordinator, &Operator)

		if err != nil {
			panic(err.Error())
		}

		SubHeader.Id = Id
		SubHeader.Date = Date
		SubHeader.Turn = Turn
		SubHeader.Line = Line
		SubHeader.Presentation = Presentation
		SubHeader.Pvalue = Pvalue
		SubHeader.Punit = Punit
		SubHeader.Perror = Perror //rate
		SubHeader.Coordinator = Coordinator
		SubHeader.Operator = Operator
		SubHeader.Header = Header

	}

	data := map[string]interface{}{
		"SubHeader": SubHeader,
		"User":      user,
	}

	defer db.Close()
	
	tmpl.ExecuteTemplate(w, "FillSalsitasWeight", data)

} // fin codigo sistemas IQ



func InsertSalsitasWeight_control(w http.ResponseWriter, r *http.Request) {
	SubHeader := models.SalsitasControl_subheader_holder{}

	var SubId int
	session, err := store.Get(r, "cookie-name")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	user := GetUser(session)

	db := dbConn()

	if r.Method == "POST" {
		v1 := r.FormValue("v1")
		v2 := r.FormValue("v2")
		v3 := r.FormValue("v3")
		v4 := r.FormValue("v4")
		v5 := r.FormValue("v5")
		//date := r.FormValue("vdate")
		unit := r.FormValue("unit")
		SubId, _ = strconv.Atoi(r.FormValue("sub"))

		selDB_sub, err := db.Query(`
		SELECT
		SC.id,
		SC.date_reg,
		SC.id_header,
		SC.turn,
		L.name,
		P.name,
		P.weight_value,
		P.weight_unit,
		P.error_rate,
		UC.fname,
		UO.fname
			FROM
				Salsitas_statistical_process SC
			INNER JOIN Line L ON
				SC.id_line = L.id
			INNER JOIN Presentation P ON
				P.id = SC.id_presentation
			INNER JOIN User_table UC ON
				UC.id = SC.id_coordinator
			INNER JOIN User_table UO ON
				SC.id_operator = UO.id
			where SC.id = ?
		`, SubId)

		if err != nil {
			panic(err.Error())
		}

		for selDB_sub.Next() {
			var Id int
			var Date string
			var Turn int
			var Line string
			var Presentation string
			var Pvalue float32
			var Punit string
			var Perror float32
			var Coordinator string
			var Operator string
			var Header int

			err = selDB_sub.Scan(&Id, &Date, &Header, &Turn,
				&Line, &Presentation, &Pvalue, &Punit, &Perror,
				&Coordinator, &Operator)

			if err != nil {
				panic(err.Error())
			}

			SubHeader.Id = Id
			SubHeader.Date = Date
			SubHeader.Turn = Turn
			SubHeader.Line = Line
			SubHeader.Presentation = Presentation
			SubHeader.Pvalue = Pvalue
			SubHeader.Punit = Punit
			SubHeader.Perror = Perror //rate
			SubHeader.Coordinator = Coordinator
			SubHeader.Operator = Operator
			SubHeader.Header = Header

		}

		//date := strings.Split(SubHeader.Date, "T")[0]
		//date := time.Now()

		insForm, err := db.Prepare(
			`INSERT INTO Weight_control(
				value1,
				value2,
				value3,
				value4,
				value5,
			
				unit,
				id_sub_header
			)
			VALUES(
				?,
				?,
				?,
				?,
				?,
				
				?,
				?
			)`)

		if err != nil {
			panic(err.Error())
		}

		insForm.Exec(v1, v2, v3, v4, v5, unit, SubId)

		defer db.Close()

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
			"SubHeader": SubHeader,
			"User":      user,
		}

		tmpl.ExecuteTemplate(w, "FillSalsitasWeight", data)
	}
}

func ConsolidatedWeightSalsitas(w http.ResponseWriter, r *http.Request) {

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
	tmpl.ExecuteTemplate(w, "ConsolidatedWeightSalsitas", user)
}

func GetWeightSalsitas(w http.ResponseWriter, r *http.Request) {
	db := dbConn()

	nLine := r.URL.Query().Get("line")
	dateInit := r.URL.Query().Get("dinit")
	dateFinal := r.URL.Query().Get("dfinal")

	selDB, err := db.Query(`
	SELECT
		Op.profile_picture,
		W.date_reg,
		S.turn,
		L.name,
		P.name,
		P.weight_value,
		P.weight_unit,
		P.error_rate,
		Co.fname,
		Co.lname,
		Op.fname,
		Op.lname,
		W.value1,
		W.value2,
		W.value3,
		W.value4,
		W.value5
	FROM
		Weight_control W
	INNER JOIN Salsitas_statistical_process S ON
		S.id = W.id_sub_header
	INNER JOIN Line L ON
		L.id = S.id_line
	INNER JOIN Presentation P ON
		P.id = S.id_presentation
	INNER JOIN User_table Co ON
		Co.id = S.id_coordinator
	INNER JOIN User_table Op ON
		Op.id = S.id_operator
	WHERE
		L.id = ? AND
		W.date_reg >= ? AND 
		W.date_reg < ?`, nLine, dateInit, dateFinal)

	if err != nil {
		panic(err.Error())
	}

	p := models.Consolidated_weight_salsitas{}

	res := []models.Consolidated_weight_salsitas{}

	for selDB.Next() {
		var Profile_picture string
		var Date string
		var Turn int
		var Line string
		var Pname string
		var Pvalue float32
		var Punit string
		var Prate float32
		var Coordinator string
		var Clname string
		var Operator string
		var Olname string
		var V1 float32
		var V2 float32
		var V3 float32
		var V4 float32
		var V5 float32

		err = selDB.Scan(&Profile_picture, &Date, &Turn, &Line, &Pname,
			&Pvalue, &Punit, &Prate, &Coordinator, &Clname,
			&Operator, &Olname,
			&V1, &V2, &V3, &V4, &V5)

		p.Profile_picture = Profile_picture
		p.Date = Date
		p.Turn = Turn
		p.Line = Line
		p.Pname = Pname
		p.Pvalue = Pvalue
		p.Punit = Punit
		p.Prate = Prate
		p.Coordinator = Coordinator
		p.Clname = Clname
		p.Operator = Operator
		p.Olname = Olname
		p.V1 = V1
		p.V2 = V2
		p.V3 = V3
		p.V4 = V4
		p.V5 = V5

		res = append(res, p)
	}

	json.NewEncoder(w).Encode(res)
	defer db.Close()
}

func GetWeightAll(w http.ResponseWriter, r *http.Request) {
	db := dbConn()

	date := r.URL.Query().Get("date")

	selDB, err := db.Query(`
	SELECT
		W.date_reg,
		S.turn,
		L.name,
		P.name,
		P.weight_value,
		P.weight_unit,
		P.error_rate,
		Co.fname,
		Co.lname,
		Op.fname,
		Op.lname,
		W.value1,
		W.value2,
		W.value3,
		W.value4,
		W.value5
	FROM
		Weight_control W
	INNER JOIN Salsitas_statistical_process S ON
		S.id = W.id_sub_header
	INNER JOIN Line L ON
		L.id = S.id_line
	INNER JOIN Presentation P ON
		P.id = S.id_presentation
	INNER JOIN User_table Co ON
		Co.id = S.id_coordinator
	INNER JOIN User_table Op ON
		Op.id = S.id_operator
	WHERE
		date(W.date_reg) = ?`, date)

	if err != nil {
		panic(err.Error())
	}

	p := models.Consolidated_weight_salsitas{}

	res := []models.Consolidated_weight_salsitas{}

	for selDB.Next() {
		var Date string
		var Turn int
		var Line string
		var Pname string
		var Pvalue float32
		var Punit string
		var Prate float32
		var Coordinator string
		var Clname string
		var Operator string
		var Olname string
		var V1 float32
		var V2 float32
		var V3 float32
		var V4 float32
		var V5 float32

		err = selDB.Scan(&Date, &Turn, &Line, &Pname,
			&Pvalue, &Punit, &Prate, &Coordinator, &Clname,
			&Operator, &Olname,
			&V1, &V2, &V3, &V4, &V5)

		p.Date = Date
		p.Turn = Turn
		p.Line = Line
		p.Pname = Pname
		p.Pvalue = Pvalue
		p.Punit = Punit
		p.Prate = Prate
		p.Coordinator = Coordinator
		p.Clname = Clname
		p.Operator = Operator
		p.Olname = Olname
		p.V1 = V1
		p.V2 = V2
		p.V3 = V3
		p.V4 = V4
		p.V5 = V5

		res = append(res, p)
	}

	json.NewEncoder(w).Encode(res)
	defer db.Close()
}

func WsEndpointS4(w http.ResponseWriter, r *http.Request) {
	//fmt.Fprintf(w, "Hello World")
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}

	log.Println("Client Connected")

	go WriterS4(ws)
}

func WriterS4(conn *websocket.Conn) {

	for {

		ticker := time.NewTicker(time.Second * 3 / 2)
		//var weight float32

		for t := range ticker.C {
			fmt.Printf("Updating At: %+v\n", t)

			response, err := http.Get("http://192.168.1.191/e-factory/3.4/php/get_weight.php?line=S4")

			if err != nil {
				fmt.Printf("The HTTP request failed with error %s\n", err)
				break
			}

			weight, _ := ioutil.ReadAll(response.Body)
			//fmt.Println("weight->>>--->"+string(weight))

			if err := conn.WriteMessage(websocket.TextMessage, []byte(weight)); err != nil {
				fmt.Println(err)
				break
			}
		}
	}

}

/*-------------------------------*/

func WsEndpointS8(w http.ResponseWriter, r *http.Request) {
	//fmt.Fprintf(w, "Hello World")
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}

	log.Println("Client Connected")

	go WriterS8(ws)
}

func WriterS8(conn *websocket.Conn) {

	for {

		ticker := time.NewTicker(time.Second * 3 / 2)
		//var weight float32

		for t := range ticker.C {
			fmt.Printf("Updating At: %+v\n", t)

			response, err := http.Get("http://192.168.1.191/e-factory/3.4/php/get_weight.php?line=S8")

			if err != nil {
				fmt.Printf("The HTTP request failed with error %s\n", err)
				break
			}

			weight, _ := ioutil.ReadAll(response.Body)
			//fmt.Println("weight->>>--->"+string(weight))

			if err := conn.WriteMessage(websocket.TextMessage, []byte(weight)); err != nil {
				fmt.Println(err)
				break
			}
		}
	}

}

/*---------------------------------------------*/
func WsEndpointWeight(w http.ResponseWriter, r *http.Request) {

	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}

	log.Println("Client Connected")

	clientData, _ := readerFilter(ws)

	line := clientData

	go WriterWeight(ws, line)
}

func WriterWeight(conn *websocket.Conn, line string) {

	for {

		ticker := time.NewTicker(time.Second / 2)

		for t := range ticker.C {

			msg, err_msg := readerFilter(conn)

			if err_msg != nil {
				defer conn.Close()
				break
			}

			if msg == "on" {

				fmt.Printf("Updating At: %+v\n", t)

				response, err := http.Get("http://192.168.1.191/e-factory/3.4/php/get_weight.php?line=" + line)

				if err != nil {
					fmt.Printf("The HTTP request failed with error %s\n", err)
					break
				}

				data, _ := ioutil.ReadAll(response.Body)

				if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
					fmt.Println(err)
					break
				}

			} else {
				fmt.Printf("Client socket close!")
				defer conn.Close()
				break
			}
		}
		defer conn.Close()
	}
}

func WeightVerification(w http.ResponseWriter, r *http.Request) {
	db := dbConn()

	selDB, err := db.Query("Select * from Header where id=2")

	if err != nil {
		panic(err.Error())
	}

	p := models.Header{}

	for selDB.Next() {
		var Id int
		var Name string
		var Cod_doc string
		var Revision_date string
		var Next_revision_date string
		var Revision_no int

		err = selDB.Scan(&Id, &Name, &Cod_doc, &Revision_date, &Next_revision_date, &Revision_no)

		if err != nil {
			panic(err.Error())
		}

		p.Id = Id
		p.Name = Name
		p.Cod_doc = Cod_doc
		p.Revision_date = Revision_date
		p.Next_revision_date = Next_revision_date
		p.Revision_no = Revision_no

	}

	/*-----------------Auth--------------------------*/
	session, err := store.Get(r, "cookie-name")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	user := GetUser(session)

	if auth := user.Authenticated; !auth {
		session.AddFlash("Acesso denegado!")
		err = session.Save(r, w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/forbidden", http.StatusFound)
		return
	}

	data := map[string]interface{}{
		"Header": p,
		"User":   user,
	}

	/*-------------------------------------------*/
	tmpl.ExecuteTemplate(w, "WeightVerification", data)
}

/*funcion que lee el canvas capturado y lo pasa a un archivo temporal*/
func SignTempGen(r *http.Request) (string, bool) {

	success := true
	fileName := ""

	r.ParseMultipartForm(10 << 20)

	file, handler, err := r.FormFile("sign")
	if err != nil {
		fmt.Println("Error Retrieving the File")
		fmt.Println(err)
		success = false
	}
	defer file.Close()
	fmt.Printf("Uploaded File: %+v\n", handler.Filename)
	fmt.Printf("File Size: %+v\n", handler.Size)
	fmt.Printf("MIME Header: %+v\n", handler.Header)

	// Create a temporary file within our temp-images directory that follows
	// a particular naming pattern
	tempFile, err := ioutil.TempFile("temp", "upload-*.png")
	if err != nil {
		fmt.Println(err)
		success = false
	}
	fileName = tempFile.Name()
	fmt.Println("tempfile name: ", fileName)
	defer tempFile.Close()

	// read all of the contents of our uploaded file into a
	// byte array
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println(err)
		success = false
	}
	// write this byte array to our temporary file
	tempFile.Write(fileBytes)

	return fileName, success
}

func WeightSigner(w http.ResponseWriter, r *http.Request) {

	if r.Method == "POST" {

		date := r.FormValue("vdate")

		fmt.Println("Date:\t" + date)

		fileName, success := SignTempGen(r)

		fmt.Printf("FIrma alias:\t%s", fileName)

		if success {
			fmt.Fprintf(w, "Successfully Uploaded File\n")
		}

	}

}

//GetWeightConsolidateByDate
//obtener consolidado de todas lineas de un dia filtrado por fecha
func GetWeightConsolidateByDate(date string) []models.Consolidated_weight_salsitas {
	db := dbConn()

	selDB, err := db.Query(`
	SELECT
		W.date_reg,
		S.turn,
		L.name,
		P.name,
		P.weight_value,
		P.weight_unit,
		P.error_rate,
		Co.fname,
		Co.lname,
		Op.fname,
		Op.lname,
		W.value1,
		W.value2,
		W.value3,
		W.value4,
		W.value5
	FROM
		Weight_control W
	INNER JOIN Salsitas_statistical_process S ON
		S.id = W.id_sub_header
	INNER JOIN Line L ON
		L.id = S.id_line
	INNER JOIN Presentation P ON
		P.id = S.id_presentation
	INNER JOIN User_table Co ON
		Co.id = S.id_coordinator
	INNER JOIN User_table Op ON
		Op.id = S.id_operator
	WHERE
		date(W.date_reg) = ?`, date)

	if err != nil {
		panic(err.Error())
	}

	p := models.Consolidated_weight_salsitas{}
	res := []models.Consolidated_weight_salsitas{}

	for selDB.Next() {
		var Date string
		var Turn int
		var Line string
		var Pname string
		var Pvalue float32
		var Punit string
		var Prate float32
		var Coordinator string
		var Clname string
		var Operator string
		var Olname string
		var V1 float32
		var V2 float32
		var V3 float32
		var V4 float32
		var V5 float32

		err = selDB.Scan(&Date, &Turn, &Line, &Pname,
			&Pvalue, &Punit, &Prate, &Coordinator, &Clname,
			&Operator, &Olname,
			&V1, &V2, &V3, &V4, &V5)

		p.Date = Date
		p.Turn = Turn
		p.Line = Line
		p.Pname = Pname
		p.Pvalue = Pvalue
		p.Punit = Punit
		p.Prate = Prate
		p.Coordinator = Coordinator
		p.Clname = Clname
		p.Operator = Operator
		p.Olname = Olname
		p.V1 = V1
		p.V2 = V2
		p.V3 = V3
		p.V4 = V4
		p.V5 = V5

		res = append(res, p)
	}
	defer db.Close()
	return res
}

func isSignDuplicate(date string, idDoc string) bool {
	db := dbConn()

	count := 0

	selDB, err := db.Query(`
	SELECT
		COUNT(D.id)
	FROM
		DocSign D
	INNER JOIN Header H ON
		H.id = D.id_header
	WHERE
		H.id = ? AND eventDate = ?`, idDoc, date)

	if err != nil {
		panic(err.Error())
	}

	for selDB.Next() {
		err = selDB.Scan(&count)
	}

	fmt.Printf("\nDate:%s\tDoc:%d\tCount:%d", date, idDoc, count)

	defer db.Close()

	if count > 0 {
		return true
	}
	return false
}

func SignProceedVerification(w http.ResponseWriter, r *http.Request) {

	date := r.URL.Query().Get("date")
	idDoc := r.URL.Query().Get("idDoc")

	json.NewEncoder(w).Encode(isSignDuplicate(date, idDoc))
}
