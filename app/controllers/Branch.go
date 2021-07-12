package controllers


import (
    //"bytes"
	"log"
    "net/http"
    //"image/jpeg"
    "encoding/json"
    "../models"
    _"github.com/go-sql-driver/mysql"
)

func GetBranch(w http.ResponseWriter, r *http.Request){
    db := dbConn()

    nId := r.URL.Query().Get("id")
    selDB, err := db.Query("SELECT * FROM Branch WHERE id_sub_classification=?", nId)

    if err != nil{
        panic(err.Error())
    }

    p := models.Branch{}

    res := []models.Branch{}
    for selDB.Next(){
        var Id int
        var Description string
        var Id_sub_classification int

        err = selDB.Scan(&Id, &Description, &Id_sub_classification)
        if err != nil {
            panic(err.Error())
        }
        p.Id    = Id
        p.Description = Description
        p.Id_sub_classification= Id_sub_classification

        res = append(res, p)
    }
    defer db.Close()
    json.NewEncoder(w).Encode(res)
}


func InsertBranch(w http.ResponseWriter, r *http.Request) {
    db := dbConn()
    if r.Method == "POST" {
        Id_sub_classification:= r.FormValue("Sub")
        Description:= r.FormValue("description")
        
        
        insForm, err := db.Prepare("INSERT INTO Branch(description, id_sub_classification ) VALUES(?, ?)")
        if err != nil {
            panic(err.Error())
        }
        insForm.Exec(Description,  Id_sub_classification)
        log.Println("INSERT: Name: " + Description)
    }
    defer db.Close()
    http.Redirect(w, r, "/Branch", 301)
}


func Branch(w http.ResponseWriter, r *http.Request) {
    db := dbConn()

    selDB, err := db.Query(`
            SELECT
            B.id,
            LTC.description,
            S.description,
            B.description,
            S.color
        FROM
            Branch B
        INNER JOIN Sub_classification S ON
            B.id_sub_classification = S.id
        INNER JOIN LineTimeClassification LTC ON
            LTC.id = S.id_LTC
        ORDER BY
            S.id_LTC,
            B.id_sub_classification
    `)
    if err != nil{
        panic(err.Error())
    }

    p := models.Branch_holder{}

    res := []models.Branch_holder{}
    for selDB.Next(){
        var Id int
        var LTC string//line time classification
        var Sub string//Subclassification
        var Description string 
        var Color string

        err = selDB.Scan(&Id, &LTC, &Sub, &Description, &Color)
        if err != nil {
            panic(err.Error())
        }
        p.Id    = Id
        p.LTC   = LTC
        p.Sub   = Sub
        p.Description = Description 
        p.Color = Color

        res = append(res, p)
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
    tmpl.ExecuteTemplate(w, "Branch", res)
    defer db.Close()
}


func NewBranch(w http.ResponseWriter, r *http.Request) {
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
    

    tmpl.ExecuteTemplate(w, "NewBranch", user)
}

func EditBranch(w http.ResponseWriter, r *http.Request) {
    db := dbConn()
    nId := r.URL.Query().Get("id")

    selDB, err := db.Query("SELECT * FROM Branch WHERE id=?", nId)
    if err != nil {
        panic(err.Error())
    }
    p := models.Sub_classification{}

    for selDB.Next() {
        var Id int
        var Description string
        var Id_sub_classification int

        err = selDB.Scan(&Id, &Description, &Id_sub_classification)
        if err != nil {
            panic(err.Error())
        }
        p.Id    = Id
    
        p.Description  = Description
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
    

    tmpl.ExecuteTemplate(w, "EditBranch", p)
    defer db.Close()
}


func UpdateBranch(w http.ResponseWriter, r *http.Request) {
    db := dbConn()
    if r.Method == "POST" {
        Description := r.FormValue("description")
        Sub := r.FormValue("Sub")
        Id := r.FormValue("uid")

        insForm, err := db.Prepare("UPDATE Branch SET description=?, id_sub_classification=? WHERE id=?")
        if err != nil {
            panic(err.Error())
        }
        insForm.Exec(Description,  Sub, Id)
        log.Println("UPDATE: Name: " + Description)
    }
    defer db.Close()
    http.Redirect(w, r, "/Branch", 301)   
}


func DeleteBranch(w http.ResponseWriter, r *http.Request) {
    db := dbConn()
    p := r.URL.Query().Get("id")
    delForm, err := db.Prepare("DELETE FROM Branch WHERE id=?")
    if err != nil {
        panic(err.Error())
    }
    delForm.Exec(p)
    log.Println("DELETE")
    defer db.Close()
    http.Redirect(w, r, "/Branch", 301)
}


