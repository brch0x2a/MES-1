package controllers

import (
	"database/sql"
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"text/template"
	"time"

	//"../models"
	"github.com/AdrianRb95/MES/app/models"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
)

func dbConn() (db *sql.DB) {
	dbDriver := "mysql"
	dbUser := "brch"
	dbPass := "tk2718"
	//dbPass := "wrf!C:w(>7:&"
	dbName := "mes"
	//dbAdress := "10.0.1.36"
	//dbAdress :=  "192.168.1.191"
	dbAdress := "localhost"
	// dbAdress := "192.168.1.97"

	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@tcp("+dbAdress+")/"+dbName+"?parseTime=true")
	if err != nil {
		panic(err.Error())
	}

	return db
}

var tmpl = template.Must(template.ParseGlob("views/*"))

// store will hold all session data
var store *sessions.CookieStore

func Init() {
	authKeyOne := securecookie.GenerateRandomKey(64)
	encryptionKeyOne := securecookie.GenerateRandomKey(32)

	store = sessions.NewCookieStore(
		authKeyOne,
		encryptionKeyOne,
	)

	store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   60 * 480,
		HttpOnly: true,
	}

	gob.Register(models.User{})
	gob.Register(models.Job{})

	log.Printf("Init()")

}

//
func GetCurrentTurn(hour int) int64 {
	var turn int64
	if 6 <= hour && hour < 14 {
		turn = 1
	} else if 14 <= hour && hour < 22 {
		turn = 2
	} else {
		turn = 3
	}
	return turn
}

// login authenticates the user

// login authenticates the user
func Login(w http.ResponseWriter, r *http.Request) {
	var usertype int
	var hash string
	var id int
	var LastId int
	var crqsId int
	var profile_picture string
	var line int64
	var presentation int64
	var coordinador int64

	db := dbConn()

	f, err := os.OpenFile("MES.log",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	defer f.Close()

	log.SetOutput(f)

	username := r.FormValue("username")
	password := r.FormValue("code")
	moment := strings.Split(time.Now().String(), " ") // moment is the string of the current date like 2019/10/08 12:46:03
	date := moment[0]
	hour, _ := strconv.Atoi(strings.Split(moment[1], ":")[0])
	fmt.Printf("\nHour: %s\n", hour)

	//turn, _ := strconv.ParseInt(r.FormValue("turn")[0:], 10, 64)
	turn := GetCurrentTurn(hour)

	if r.FormValue("line") != "" {
		line, _ = strconv.ParseInt(r.FormValue("line")[0:], 10, 64)
	} else {
		fmt.Printf("falto la linea")
		http.Redirect(w, r, "/forbidden", http.StatusFound)
		return
	}

	if r.FormValue("presentation") != "" {
		presentation, _ = strconv.ParseInt(r.FormValue("presentation")[0:], 10, 64)
	} else {
		fmt.Printf("falto la precentacion")
		http.Redirect(w, r, "/forbidden", http.StatusFound)
		return
	}

	if r.FormValue("coordinador") != "" {
		fmt.Printf("falto el coordinador")
		coordinador, _ = strconv.ParseInt(r.FormValue("coordinador")[0:], 10, 64)
	} else {
		http.Redirect(w, r, "/forbidden", http.StatusFound)
		return
	}

	header := 2

	selDB, err := db.Query(`
	SELECT
		U.id,
		U.password,
		U.profile_picture,
		P.level
	FROM
		User_table U
	INNER JOIN Privilege P ON
		U.id_privilege = P.id
	WHERE
		U.nick_name = ?
	`, username)

	if err != nil {
		panic(err.Error())
	}

	for selDB.Next() {

		err = selDB.Scan(&id, &hash, &profile_picture, &usertype)
		if err != nil {
			panic(err.Error())
		}
		log.Printf("Found!\n")
	}
	defer db.Close()

	session, err := store.Get(r, "cookie-name")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	//log.Printf("Pass: %s\tHash: %s", password, hash)
	//if r.FormValue("code") != "code" {

	if presentation < 1 {
		http.Redirect(w, r, "/forbidden", http.StatusFound)
		return
	}

	if !CheckPasswordHash(password, hash) {
		log.Printf("Match!\n")

		if password == "" {
			session.AddFlash("Ingrese una clave")
		}
		session.AddFlash("Clave incorrecta!")
		err = session.Save(r, w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/forbidden", http.StatusFound)
		return
	}

	/*------------------Setup---docs------------------*/
	/*-----------COntrol-estadistico-proceso-----------*/
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
		id, header)
	if err != nil {
		log.Fatal("Cannot run insert statement", err)
	}

	idW, _ := res.LastInsertId()

	LastId = int(idW)

	defer db.Close()

	/*-----------CRQS-ALLERGEN-----------*/
	insForm_crqs, err := db.Prepare(`
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

	res_crqs, err_crqs := insForm_crqs.Exec(date,
		turn, line, presentation, coordinador,
		id, 3)
	if err_crqs != nil {
		log.Fatal("Cannot run insert statement", err_crqs)
	}

	idW_crqs, _ := res_crqs.LastInsertId()

	crqsId = int(idW_crqs)

	defer db.Close()

	/*---------------------------------------*/
	user := &models.User{
		UserId:          id,
		Username:        username,
		Authenticated:   true,
		Usertype:        usertype,
		Date:            date,
		Turn:            turn,
		Line:            line,
		Profile_picture: profile_picture,
		Presentation:    presentation,
		Coordinator:     coordinador,
		Config:          true,
		WeightSubId:     LastId,
		CRQSId:          crqsId,
	}

	session.Values["user"] = user
	//session.Values["weigthId"] = LastId

	err = session.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	lineHolder := GetLineObjectBy(int(line))
	userHolder := GetUserHolderBy(int(id))

	log.Printf("Login()")
	log.Printf("Ingreso a MES: %s | %s | %s |\t>> en la linea %s",
		userHolder.Nick_name, userHolder.Fname, userHolder.Lname, lineHolder.Name)

	if user.Usertype > 1 {
		http.Redirect(w, r, "/generalDashboard", http.StatusFound)
	} else {
		http.Redirect(w, r, "/specificDashBoard", http.StatusFound)
	}

}

// logout revokes authentication for a user
func Logout(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "cookie-name")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	session.Values["user"] = models.User{}
	session.Options.MaxAge = -1

	err = session.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Printf("Logout()")

	http.Redirect(w, r, "/", http.StatusFound)
}

// secret displays the secret message for authorized users
func Secret(w http.ResponseWriter, r *http.Request) {
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

	tmpl.ExecuteTemplate(w, "secret.gohtml", user.Username)
}

func Forbidden(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "cookie-name")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	flashMessages := session.Flashes()
	err = session.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.ExecuteTemplate(w, "Forbidden", flashMessages)
}

// GetUser returns a user from session s
// on error returns an empty user
func GetUser(s *sessions.Session) models.User {
	val := s.Values["user"]
	var user = models.User{}
	user, ok := val.(models.User)
	if !ok {
		return models.User{Authenticated: false}
	}
	return user
}

func GetVal(s *sessions.Session) int {
	val, ok := s.Values["weigthId"]
	if !ok {
		return -1
	}
	return val.(int)
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
