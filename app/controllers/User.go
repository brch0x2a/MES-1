package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"../models"
	_ "github.com/go-sql-driver/mysql"
)

func getUserType(Nick_name string) int {
	var level int

	db := dbConn()

	selDB, err := db.Query(`
		SELECT
			P.level
		FROM
			User_table U
		INNER JOIN Privilege P ON
			P.id = U.id_privilege
		WHERE
			U.nick_name = ?;
	`, Nick_name)

	if err != nil {
		panic(err.Error())
	}

	for selDB.Next() {

		selDB.Scan(&level)

	}

	defer db.Close()
	return level
}

func GetUserTypeAccesor(w http.ResponseWriter, r *http.Request) {

	nick_name := r.FormValue("username")

	metaUser := getUserType(nick_name)

	json.NewEncoder(w).Encode(metaUser)
}

func GetUserHolderBy(id int) models.User_holder {
	db := dbConn()

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
	INNER JOIN Factory F ON
		F.id = U.id_factory
	INNER JOIN Privilege P ON
		P.id = U.id_privilege
	WHERE
		U.id = ?
	`, id)

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

		selDB.Scan(&Id, &Nick_name, &Fname, &Lname, &Privilege, &Factory)

		p.Id = Id
		p.Nick_name = Nick_name
		p.Fname = Fname
		p.Lname = Lname
		p.Privilege = Privilege
		p.Factory = Factory
	}

	defer db.Close()
	return p
}

func GetMetaUser(w http.ResponseWriter, r *http.Request) {

	metaUser := models.MetaUser{}
	presentation := models.Presentation{}
	line := models.Line{}
	userHolder := models.User_holder{}

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

	presentation = GetObjectPresentationBy(int(user.Presentation))
	line = GetLineObjectBy(int(user.Line))

	userHolder = GetUserHolderBy(user.UserId)

	metaUser.Line = line.Name
	metaUser.LineID = strconv.Itoa(line.Id)
	metaUser.Presentation = presentation.Name
	metaUser.Profile_picture = user.Profile_picture

	metaUser.Nick_name = userHolder.Nick_name
	metaUser.Fname = userHolder.Fname
	metaUser.Lname = userHolder.Lname
	metaUser.ProductPhoto = presentation.Photo

	json.NewEncoder(w).Encode(metaUser)
}
