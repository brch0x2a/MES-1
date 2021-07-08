package controllers

import (
	//"bytes"
	//"fmt"
	"log"
	"net/http"

	//"image/jpeg"
	"encoding/json"

	"../models"
	// _ "github.com/go-sql-driver/mysql"
)

func Gestion_personal(w http.ResponseWriter, r *http.Request) {
	db := dbConn()

	selDB, err := db.Query(`
			SELECT
			U.id,
			U.nick_name,
			U.fname,
			U.lname,
			U.profile_picture,
			P.description,
			F.name
		FROM
			User_table U
		INNER JOIN Privilege P ON
			U.id_privilege = P.id
		INNER JOIN Factory F ON
			U.id_factory = F.id
	`)
	if err != nil {
		panic(err.Error())
	}

	p := models.User_holder{}

	res := []models.User_holder{}
	for selDB.Next() {
		var Id int
		var Nick_name string
		//var Password  string
		var Fname string
		var Lname string
		var Profile_picture string
		var Privilege string
		var Factory string

		err = selDB.Scan(&Id, &Nick_name, &Fname,
			&Lname, &Profile_picture, &Privilege, &Factory)

		if err != nil {
			panic(err.Error())
		}
		p.Id = Id
		p.Nick_name = Nick_name
		p.Fname = Fname
		p.Lname = Lname
		p.Privilege = Privilege
		p.Profile_picture = Profile_picture
		p.Factory = Factory

		res = append(res, p)
	}

	/*-----------------Auth--------------------------*/
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
	/*-------------------------------------------*/
	tmpl.ExecuteTemplate(w, "Gestion_personal", res)
}

func NewUser(w http.ResponseWriter, r *http.Request) {

	/*-----------------Auth--------------------------*/
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
	/*-------------------------------------------*/
	tmpl.ExecuteTemplate(w, "NewUser", nil)

}

func InsertUser(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	if r.Method == "POST" {

		Nick_name := r.FormValue("num")
		Fname := r.FormValue("fname")
		Lname := r.FormValue("lname")
		Privilege := r.FormValue("level")
		Factory := r.FormValue("factory")
		Pass := r.FormValue("pass")

		hash, _ := HashPassword(Pass)
		var filePath string

		filePath = FileHandlerEngine(r, "capture", "public/images/user/", "user-*.png")

		insForm, err := db.Prepare(`
		INSERT INTO User_table(
			nick_name,
			password,
			fname,
			lname,
			profile_picture,
			id_privilege,
			id_factory
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
		insForm.Exec(Nick_name, hash, Fname, Lname, filePath, Privilege, Factory)
	}
	defer db.Close()
	http.Redirect(w, r, "/gestion_personal", 301)
}

func GetPrivilege(w http.ResponseWriter, r *http.Request) {

	db := dbConn()

	selDB, err := db.Query("Select * from Privilege")
	if err != nil {
		panic(err.Error())
	}

	p := models.Privilege{}

	res := []models.Privilege{}
	for selDB.Next() {
		var Id int
		var Level int
		var Description string

		err = selDB.Scan(&Id, &Level, &Description)
		if err != nil {
			panic(err.Error())
		}
		p.Id = Id

		p.Level = Level
		p.Description = Description
		res = append(res, p)
	}
	defer db.Close()
	json.NewEncoder(w).Encode(res)
}

func GetUsers(w http.ResponseWriter, r *http.Request) {

	db := dbConn()

	selDB, err := db.Query("Select id, nick_name, fname, lname from User_table")
	if err != nil {
		panic(err.Error())
	}

	p := models.User_holder{}

	res := []models.User_holder{}
	for selDB.Next() {
		var Id int
		var Nick_name string
		var Fname string
		var Lname string

		err = selDB.Scan(&Id, &Nick_name, &Fname, &Lname)
		if err != nil {
			panic(err.Error())
		}
		p.Id = Id
		p.Nick_name = Nick_name
		p.Fname = Fname
		p.Lname = Lname

		res = append(res, p)
	}
	defer db.Close()
	json.NewEncoder(w).Encode(res)
}

func GetLineManagers(w http.ResponseWriter, r *http.Request) {

	db := dbConn()

	selDB, err := db.Query(`
		SELECT 
		U.id, U.nick_name, U.fname, U.lname
		FROM
			User_table U
				INNER JOIN
			Privilege P ON P.id = U.id_privilege
			Where P.level > 3
	`)
	if err != nil {
		panic(err.Error())
	}

	p := models.User_holder{}

	res := []models.User_holder{}
	for selDB.Next() {
		var Id int
		var Nick_name string
		var Fname string
		var Lname string

		err = selDB.Scan(&Id, &Nick_name, &Fname, &Lname)
		if err != nil {
			panic(err.Error())
		}
		p.Id = Id
		p.Nick_name = Nick_name
		p.Fname = Fname
		p.Lname = Lname

		res = append(res, p)
	}
	defer db.Close()
	json.NewEncoder(w).Encode(res)
}

func GetMecanicsObj() []models.User_holder {

	db := dbConn()

	selDB, err := db.Query(`
		SELECT 
		U.id, U.nick_name, U.fname, U.lname
		FROM
			User_table U
				INNER JOIN
			Privilege P ON P.id = U.id_privilege
			Where P.level = 2
	`)
	if err != nil {
		panic(err.Error())
	}

	p := models.User_holder{}

	res := []models.User_holder{}
	for selDB.Next() {
		var Id int
		var Nick_name string
		var Fname string
		var Lname string

		err = selDB.Scan(&Id, &Nick_name, &Fname, &Lname)
		if err != nil {
			panic(err.Error())
		}
		p.Id = Id
		p.Nick_name = Nick_name
		p.Fname = Fname
		p.Lname = Lname

		res = append(res, p)
	}
	defer db.Close()

	return res

}

func GetMecanics(w http.ResponseWriter, r *http.Request) {

	res := GetMecanicsObj()

	json.NewEncoder(w).Encode(res)
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	p := r.URL.Query().Get("id")
	delForm, err := db.Prepare("DELETE FROM User_table WHERE id=?")
	if err != nil {
		panic(err.Error())
	}
	delForm.Exec(p)
	log.Println("DELETE")
	defer db.Close()
	http.Redirect(w, r, "/gestion_personal", 301)
}

func EditUser(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	nId := r.URL.Query().Get("id")

	selDB, err := db.Query(`
	SELECT
		U.id,
		U.nick_name,
		U.fname,
		U.lname,
		P.description,
		F.name
	FROM
		User_table U
	INNER JOIN Privilege P ON
		U.id_privilege = P.id
	INNER JOIN Factory F ON
		U.id_factory = F.id
	WHERE U.id=?`, nId)

	if err != nil {
		panic(err.Error())
	}
	p := models.User_holder{}

	for selDB.Next() {
		var Id int
		var Nick_name string
		//Password  string
		var Fname string
		var Lname string
		var Privilege string
		var Factory string

		err = selDB.Scan(&Id, &Nick_name, &Fname,
			&Lname, &Privilege, &Factory)

		if err != nil {
			panic(err.Error())
		}

		p.Id = Id
		p.Nick_name = Nick_name
		p.Fname = Fname
		p.Lname = Lname
		p.Privilege = Privilege
		p.Factory = Factory
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

	tmpl.ExecuteTemplate(w, "EditUser", p)
	defer db.Close()
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	if r.Method == "POST" {
		Nick := r.FormValue("num")
		Fname := r.FormValue("fname")
		Lname := r.FormValue("lname")
		Pass := r.FormValue("pass")
		Privilege := r.FormValue("level")
		Factory := r.FormValue("factory")
		var filePath string = ""

		Id := r.FormValue("uid")

		hash, _ := HashPassword(Pass)

		filePath = FileHandlerEngine(r, "capture", "public/images/user/", "user-*.png")

		insForm, err := db.Prepare(`
		UPDATE User_table SET nick_name=?, password=?, 
		fname=?, lname=?, profile_picture=?, id_privilege=?, id_factory=?
		WHERE id=?`)

		if err != nil {
			panic(err.Error())
		}

		insForm.Exec(Nick, hash, Fname,
			Lname, filePath, Privilege, Factory, Id)

	}
	defer db.Close()
	http.Redirect(w, r, "/gestion_personal", 301)
}
