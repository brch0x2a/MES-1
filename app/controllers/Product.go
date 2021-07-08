package controllers

import (

	//"bytes"

	"log"
	"net/http"

	//"image/jpeg"
	"encoding/json"

	"../models"
	_ "github.com/go-sql-driver/mysql"
)

//-------------Products---------------------
func Products(w http.ResponseWriter, r *http.Request) {

	db := dbConn()

	selDB, err := db.Query(`
        Select 
            id, 
            sku, 
            name,
			photo,
			psi_bottom,
			psi_top,
			bares_bottom,
			bares_top,
			lung_bottom,
			lung_top,
			interchange_bottom,
			interchange_top,
			hopper_bottom,
			hopper_top,
			fill_bottom,
			fill_top
        from Product
    `)
	if err != nil {
		panic(err.Error())
	}

	p := models.Product{}

	res := []models.Product{}
	for selDB.Next() {
		var Id int
		var Sku int
		var Name string
		var Photo string

		var PSI_bottom float32
		var PSI_top float32

		var Bares_bottom float32
		var Bares_top float32

		var Lung_bottom float32
		var Lung_top float32

		var Interchange_bottom float32
		var Interchange_top float32

		var Hopper_bottom float32
		var Hopper_top float32

		var Fill_bottom float32
		var Fill_top float32

		err = selDB.Scan(&Id, &Sku, &Name, &Photo,
			&PSI_bottom, &PSI_top,
			&Bares_bottom, &Bares_top,
			&Lung_bottom, &Lung_top,
			&Interchange_bottom, &Interchange_top,
			&Hopper_bottom, &Hopper_top,
			&Fill_bottom, &Fill_top)

		if err != nil {
			panic(err.Error())
		}
		p.Id = Id
		p.Sku = Sku
		p.Name = Name
		p.Photo = Photo

		p.PSI_bottom = PSI_bottom
		p.PSI_top = PSI_top
		p.Bares_bottom = Bares_bottom
		p.Bares_top = Bares_top
		p.Lung_bottom = Lung_bottom
		p.Lung_top = Lung_top
		p.Interchange_bottom = Interchange_bottom
		p.Interchange_top = Interchange_top
		p.Hopper_bottom = Hopper_bottom
		p.Hopper_top = Hopper_top
		p.Fill_bottom = Fill_bottom
		p.Fill_top = Fill_top

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
	tmpl.ExecuteTemplate(w, "Products", res)
	defer db.Close()

}

func GetProductBy(w http.ResponseWriter, r *http.Request) {
	db := dbConn()

	nID := r.URL.Query().Get("id")

	selDB, err := db.Query(`
        Select 
            id, 
            sku, 
            name,
			photo,
			psi_bottom,
			psi_top,
			bares_bottom,
			bares_top,
			lung_bottom,
			lung_top,
			interchange_bottom,
			interchange_top,
			hopper_bottom,
			hopper_top,
			fill_bottom,
			fill_top
		from Product
		where 
		id = ? 
    `, nID)
	if err != nil {
		panic(err.Error())
	}

	p := models.Product{}

	for selDB.Next() {
		var Id int
		var Sku int
		var Name string
		var Photo string

		var PSI_bottom float32
		var PSI_top float32

		var Bares_bottom float32
		var Bares_top float32

		var Lung_bottom float32
		var Lung_top float32

		var Interchange_bottom float32
		var Interchange_top float32

		var Hopper_bottom float32
		var Hopper_top float32

		var Fill_bottom float32
		var Fill_top float32

		err = selDB.Scan(&Id, &Sku, &Name, &Photo,
			&PSI_bottom, &PSI_top,
			&Bares_bottom, &Bares_top,
			&Lung_bottom, &Lung_top,
			&Interchange_bottom, &Interchange_top,
			&Hopper_bottom, &Hopper_top,
			&Fill_bottom, &Fill_top)
		if err != nil {
			panic(err.Error())
		}

		p.Id = Id
		p.Sku = Sku
		p.Name = Name
		p.Photo = Photo

		p.PSI_bottom = PSI_bottom
		p.PSI_top = PSI_top
		p.Bares_bottom = Bares_bottom
		p.Bares_top = Bares_top
		p.Lung_bottom = Lung_bottom
		p.Lung_top = Lung_top
		p.Interchange_bottom = Interchange_bottom
		p.Interchange_top = Interchange_top
		p.Hopper_bottom = Hopper_bottom
		p.Hopper_top = Hopper_top
		p.Fill_bottom = Fill_bottom
		p.Fill_top = Fill_top

	}

	defer db.Close()
	json.NewEncoder(w).Encode(p)

}

func GetProducts(w http.ResponseWriter, r *http.Request) {
	db := dbConn()

	selDB, err := db.Query(`
        Select 
            id, 
            sku, 
            name,
			photo,
			psi_bottom,
			psi_top,
			bares_bottom,
			bares_top,
			lung_bottom,
			lung_top,
			interchange_bottom,
			interchange_top,
			hopper_bottom,
			hopper_top,
			fill_bottom,
			fill_top
        from Product
    `)
	if err != nil {
		panic(err.Error())
	}

	p := models.Product{}

	res := []models.Product{}
	for selDB.Next() {
		var Id int
		var Sku int
		var Name string
		var Photo string

		var PSI_bottom float32
		var PSI_top float32

		var Bares_bottom float32
		var Bares_top float32

		var Lung_bottom float32
		var Lung_top float32

		var Interchange_bottom float32
		var Interchange_top float32

		var Hopper_bottom float32
		var Hopper_top float32

		var Fill_bottom float32
		var Fill_top float32

		err = selDB.Scan(&Id, &Sku, &Name, &Photo,
			&PSI_bottom, &PSI_top,
			&Bares_bottom, &Bares_top,
			&Lung_bottom, &Lung_top,
			&Interchange_bottom, &Interchange_top,
			&Hopper_bottom, &Hopper_top,
			&Fill_bottom, &Fill_top)
		if err != nil {
			panic(err.Error())
		}

		p.Id = Id
		p.Sku = Sku
		p.Name = Name
		p.Photo = Photo

		p.PSI_bottom = PSI_bottom
		p.PSI_top = PSI_top
		p.Bares_bottom = Bares_bottom
		p.Bares_top = Bares_top
		p.Lung_bottom = Lung_bottom
		p.Lung_top = Lung_top
		p.Interchange_bottom = Interchange_bottom
		p.Interchange_top = Interchange_top
		p.Hopper_bottom = Hopper_bottom
		p.Hopper_top = Hopper_top
		p.Fill_bottom = Fill_bottom
		p.Fill_top = Fill_top

		res = append(res, p)
	}
	defer db.Close()
	json.NewEncoder(w).Encode(res)
}

func EditProduct(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	nId := r.URL.Query().Get("id")
	selDB, err := db.Query(`
            SELECT 
                id, 
                sku,
                name, 
                photo
            FROM
                Product
            WHERE
                id = ?
                `, nId)
	if err != nil {
		panic(err.Error())
	}
	p := models.Product{}

	for selDB.Next() {
		var Id int
		var Sku int
		var Name string
		var Photo string

		err = selDB.Scan(&Id, &Sku, &Name, &Photo)
		if err != nil {
			panic(err.Error())
		}
		p.Id = Id

		p.Sku = Sku
		p.Name = Name
		p.Photo = Photo

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
	tmpl.ExecuteTemplate(w, "EditProduct", p)
	defer db.Close()
}

func UpdateProduct(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	if r.Method == "POST" {
		Name := r.FormValue("name")
		Sku := r.FormValue("sku")

		Id := r.FormValue("uid")

		psi_bottom := r.FormValue("psi_bottom")
		psi_top := r.FormValue("psi_top")

		bares_bottom := r.FormValue("bares_bottom")
		bares_top := r.FormValue("bares_top")

		lung_bottom := r.FormValue("lung_bottom")
		lung_top := r.FormValue("lung_top")

		interchange_bottom := r.FormValue("interchange_bottom")
		interchange_top := r.FormValue("interchange_top")

		hopper_bottom := r.FormValue("hopper_bottom")
		hopper_top := r.FormValue("hopper_top")

		fill_bottom := r.FormValue("fill_bottom")
		fill_top := r.FormValue("fill_top")

		filePath := FileHandlerEngine(r, "photo", "public/images/Product/", "product-*.png")

		if filePath != "public/images/ul_white.png" {
			// filePath = "/public/images/doypack.png"
			insForm, err := db.Prepare(`
			UPDATE Product 
			SET 
				name = ?,
				sku = ?,
				photo = ?,
				psi_bottom = ?,
				psi_top = ?,
				bares_bottom = ?,
				bares_top = ?,
				lung_bottom = ?,
				lung_top = ?,
				interchange_bottom = ?,
				interchange_top = ?,
				hopper_bottom = ?,
				hopper_top = ?,
				fill_bottom = ?,
				fill_top = ?
			WHERE
				id = ?
			`)

			if err != nil {
				panic(err.Error())
			}

			insForm.Exec(Name, Sku, filePath,
				psi_bottom, psi_top,
				bares_bottom, bares_top,
				lung_bottom, lung_top,
				interchange_bottom, interchange_top,
				hopper_bottom, hopper_top,
				fill_bottom, fill_top,
				Id)
		} else {
			insForm, err := db.Prepare(`
			UPDATE Product 
			SET 
				name = ?,
				sku = ?,
				psi_bottom = ?,
				psi_top = ?,
				bares_bottom = ?,
				bares_top = ?,
				lung_bottom = ?,
				lung_top = ?,
				interchange_bottom = ?,
				interchange_top = ?,
				hopper_bottom = ?,
				hopper_top = ?,
				fill_bottom = ?,
				fill_top = ?
			WHERE
				id = ?
			`)

			if err != nil {
				panic(err.Error())
			}

			insForm.Exec(Name, Sku,
				psi_bottom, psi_top,
				bares_bottom, bares_top,
				lung_bottom, lung_top,
				interchange_bottom, interchange_top,
				hopper_bottom, hopper_top,
				fill_bottom, fill_top,
				Id)
		}

	}
	defer db.Close()
	http.Redirect(w, r, "/Products", 301)

}

func NewProduct(w http.ResponseWriter, r *http.Request) {
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
	tmpl.ExecuteTemplate(w, "NewProduct", nil)
}

func InsertProduct(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	if r.Method == "POST" {
		Name := r.FormValue("name")
		Sku := r.FormValue("sku")

		psi_bottom := r.FormValue("psi_bottom")
		psi_top := r.FormValue("psi_top")

		bares_bottom := r.FormValue("bares_bottom")
		bares_top := r.FormValue("bares_top")

		lung_bottom := r.FormValue("lung_bottom")
		lung_top := r.FormValue("lung_top")

		interchange_bottom := r.FormValue("interchange_bottom")
		interchange_top := r.FormValue("interchange_top")

		hopper_bottom := r.FormValue("hopper_bottom")
		hopper_top := r.FormValue("hopper_top")

		fill_bottom := r.FormValue("fill_bottom")
		fill_top := r.FormValue("fill_top")

		filePath := FileHandlerEngine(r, "photo", "public/images/Product/", "product-*.png")

		if filePath == "public/images/ul_white.png" {
			filePath = "/public/images/doypack.png"
		}

		insForm, err := db.Prepare(`
		INSERT INTO Product(
			name,
			sku,
			photo,
			psi_bottom,
			psi_top,
			bares_bottom,
			bares_top,
			lung_bottom,
			lung_top,
			interchange_bottom,
			interchange_top,
			hopper_bottom,
			hopper_top,
			fill_bottom,
			fill_top 
		)
		VALUES(
			?, ?, ?, ?, ?,
			?, ?, ?, ?, ?,
			?, ?, ?, ?, ?)
		
		`)
		if err != nil {
			panic(err.Error())
		}
		insForm.Exec(
			Name, Sku, filePath, psi_bottom, psi_top,
			bares_bottom, bares_top, lung_bottom, lung_top,
			interchange_bottom, interchange_top, hopper_bottom, hopper_top,
			fill_bottom, fill_top)
		log.Println("INSERT: Name: " + Name + " | Sku: " + Sku)
	}
	defer db.Close()
	http.Redirect(w, r, "/Products", 301)
}

func DeleteProduct(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	p := r.URL.Query().Get("id")
	delForm, err := db.Prepare("DELETE FROM Product WHERE id=?")
	if err != nil {
		panic(err.Error())
	}
	delForm.Exec(p)
	log.Println("DELETE")
	defer db.Close()
	http.Redirect(w, r, "/Products", 301)
}
