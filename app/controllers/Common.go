package controllers

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"../models"
	//_ "github.com/go-sql-driver/mysql"

	"github.com/jung-kurt/gofpdf"
)

func Index(w http.ResponseWriter, r *http.Request) {

	db := dbConn()

	selDB, err := db.Query(`
				           SELECT
				           Pr.id,
				           P.name,
				           Pr.name,
				           Pr.weight_unit,
				           Pr.weight_value,
				           Pr.error_rate,
						   Pr.box_amount,
						   P.photo
				           FROM
				               Presentation Pr
				           INNER JOIN Product P ON
				               P.id = Pr.id_product
					   `)

	if err != nil {
		fmt.Println("oh no!\t):")
		panic(err.Error())
	}

	p := models.Presentation_holder{}

	res := []models.Presentation_holder{}
	for selDB.Next() {
		var Id int
		var Product string
		var Name string
		var Weight_unit string
		var Weight_value float32
		var Error_rate float32
		var Box_amount int
		var Photo string

		err = selDB.Scan(&Id, &Product,
			&Name, &Weight_unit, &Weight_value,
			&Error_rate, &Box_amount, &Photo)

		if err != nil {
			panic(err.Error())
		}

		p.Id = Id
		p.Product = Product
		p.Name = Name
		p.Weight_unit = Weight_unit
		p.Weight_value = Weight_value
		p.Error_rate = Error_rate
		p.Box_amount = Box_amount
		p.Photo = Photo

		res = append(res, p)
	}
	defer db.Close()
	tmpl.ExecuteTemplate(w, "Index", res)

}

func GetRelativeCurrentTurn(date time.Time) int {

	now := date.Hour()
	turn := 1

	if now > 5 && now < 14 {
		turn = 1
	} else if now > 14 && now < 22 {
		turn = 2
	} else {
		turn = 3
	}
	return turn
}

func Manual(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "cookie-name")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	user := GetUser(session)
	weightId := GetVal(session)

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

	fmt.Printf("last weight id: %d \n\n", weightId)
	fmt.Printf("Name: %s | Date: %s | Turn: %d | Line: %d | Presentation: %d\n",
		user.Username, user.Date, user.Turn, user.Line, user.Presentation)

	tmpl.ExecuteTemplate(w, "Manual", user)
}

func ReportFailures(w http.ResponseWriter, r *http.Request) {

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

	tmpl.ExecuteTemplate(w, "ReportFailures", user)
}

func DFOS(w http.ResponseWriter, r *http.Request) {

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

	tmpl.ExecuteTemplate(w, "DFOS", user)
}

func PBIOEEReport(w http.ResponseWriter, r *http.Request) {

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

	tmpl.ExecuteTemplate(w, "PBIOEEReport", user)
}

func PBIEWOReport(w http.ResponseWriter, r *http.Request) {

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

	tmpl.ExecuteTemplate(w, "PBIEWOReport", user)
}

func SetupJob(w http.ResponseWriter, r *http.Request) {

	db := dbConn()

	selDB, err := db.Query(`
        SELECT
        Pr.id,
        P.name,
        Pr.name,
        Pr.weight_unit,
        Pr.weight_value,
        Pr.error_rate,
        Pr.box_amount
        FROM
            Presentation Pr
        INNER JOIN Product P ON
            P.id = Pr.id_product
    `)
	if err != nil {
		panic(err.Error())
	}

	p := models.Presentation_holder{}

	res := []models.Presentation_holder{}
	for selDB.Next() {
		var Id int
		var Product string
		var Name string
		var Weight_unit string
		var Weight_value float32
		var Error_rate float32
		var Box_amount int

		err = selDB.Scan(&Id, &Product,
			&Name, &Weight_unit, &Weight_value,
			&Error_rate, &Box_amount)

		if err != nil {
			panic(err.Error())
		}

		p.Id = Id
		p.Product = Product
		p.Name = Name
		p.Weight_unit = Weight_unit
		p.Weight_value = Weight_value
		p.Error_rate = Error_rate
		p.Box_amount = Box_amount

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

	data := map[string]interface{}{
		"User":         user,
		"Presentation": res,
	}

	tmpl.ExecuteTemplate(w, "SetupJob", data)
}

func SetJob(w http.ResponseWriter, r *http.Request) {

	session, err := store.Get(r, "cookie-name")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	userTemp := GetUser(session)

	if r.Method == "POST" {

		date := r.FormValue("vdate")
		turn, _ := strconv.ParseInt(r.FormValue("turn")[0:], 10, 64)

		line, _ := strconv.ParseInt(r.FormValue("line")[0:], 10, 64)

		presentation, _ := strconv.ParseInt(r.FormValue("presentation")[0:], 10, 64)

		coordinador, _ := strconv.ParseInt(r.FormValue("coordinador")[0:], 10, 64)

		user := &models.User{
			UserId:        userTemp.UserId,
			Username:      userTemp.Username,
			Authenticated: userTemp.Authenticated,
			Usertype:      userTemp.Usertype,
			Date:          date,
			Turn:          turn,
			Line:          line,
			Presentation:  presentation,
			Coordinator:   coordinador,
			Config:        true,
		}

		fmt.Printf("Name: %s | Date: %s | Turn: %d | Line: %d | Presentation: %d\n",
			user.Username, user.Date, user.Turn, user.Line, user.Presentation)

		session.Values["user"] = user

		err = session.Save(r, w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	session2, err := store.Get(r, "cookie-name")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	userCurrent := GetUser(session2)

	fmt.Printf("Name: %s | Date: %s | Turn: %d | Line: %d | Presentation: %d\n",
		userCurrent.Username, userCurrent.Date, userCurrent.Turn, userCurrent.Line, userCurrent.Presentation)

	http.Redirect(w, r, "/setupJob", http.StatusFound)
}

/*
func TestPDF(w http.ResponseWriter, r *http.Request) {

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(40, 10, "Hello, world")
	err := pdf.OutputFileAndClose("public/docs/hello.pdf")

	if err != nil {
		panic(err.Error())
	}
}
*/
func strDelimit(str string, sepstr string, sepcount int) string {
	pos := len(str) - sepcount
	for pos > 0 {
		str = str[:pos] + sepstr + str[pos:]
		pos = pos - sepcount
	}
	return str
}

func WeightSigner2(w http.ResponseWriter, r *http.Request) {

	date := r.FormValue("vdate")

	duplicate := isSignDuplicate(date, "2")

	fmt.Println("Is duplicate?\t", duplicate)

	fileName, success := SignTempGen(r)

	if success {
		fmt.Printf("FIrma alias:\t%s", fileName)
	}

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
	//-----------------------------------------------------
	type values struct {
		V1, V2, V3, V4, V5 string
	}
	//countryList := make([]values, 0, 8)
	headerWeight := []string{"Fecha", "Turno", "Linea", "Presentacion", "Valor", "Unidad", "%Error", "Coordinador", "Operador", "V1", "V2", "V3", "V4", "V5"}

	weightData := GetWeightConsolidateByDate(date)

	//---------------------------------------------------
	pdf := gofpdf.New("P", "mm", "A4", "")
	titleStr := "Reporte de Pesos"

	pdf.SetTitle(titleStr, false)
	pdf.SetAuthor("Admin", false)

	pdf.SetHeaderFunc(func() {
		Header := GetHeaderObjBy(2)

		fmt.Println("Name\t", Header.Name)

		//image
		pdf.Image("public/images/logo00.jpg", 10, 10, 0, 0, false, "", 0, "")

		// Arial bold 15

		// Calculate width of title and position
		//wd := pdf.GetStringWidth(titleStr) + 6
		//pdf.SetX((210 - wd) / 2)
		w := []float64{60, 40, 30}
		wSum := 0.0
		for _, v := range w {
			wSum += v
		}
		left := (210 - wSum) / 2
		pdf.SetX(left)

		// Colors of frame, background and text
		pdf.SetDrawColor(0, 0, 0)
		pdf.SetFillColor(255, 255, 255)
		pdf.SetTextColor(0, 0, 0)

		// Thickness of frame (1 mm)
		pdf.SetLineWidth(0.3)
		pdf.SetFont("Arial", "B", 9)
		// Title
		pdf.CellFormat(w[0], 9, "Unilever Planta Belen Costa Rica", "1", 0, "C", false, 0, "")
		pdf.CellFormat(61, 9, "Cod Documento: "+Header.Cod_doc, "1", 0, "C", true, 0, "")
		pdf.SetX(161)
		pdf.CellFormat(0, 9, "Revision No: "+strconv.Itoa(Header.Revision_no), "1", 0, "C", false, 0, "")
		pdf.Ln(-1)
		pdf.SetX(left)
		//------------
		pdf.CellFormat(160, 9, Header.Name, "1", 0, "C", false, 0, "")
		pdf.Ln(-1)
		pdf.SetX(left)
		//-------------
		pdf.CellFormat(80, 9, "Fecha de Ultima Revision: "+Header.Revision_date, "1", 0, "C", false, 0, "")
		pdf.CellFormat(80, 9, "Fecha de Proxima Revision: "+Header.Next_revision_date, "1", 0, "C", false, 0, "")

		// Line break
		pdf.Ln(10)
	})
	pdf.SetFooterFunc(func() {

		owner := GetUserHolderBy(user.UserId)

		fmt.Println("owner:\t", owner.Fname)

		//sign_name := "upload-613055199"

		// Position at 1.5 cm from bottom
		pdf.SetY(-30)
		pdf.CellFormat(0, 9, "Analista: "+owner.Fname+"  "+owner.Lname, "", 0, "R", false, 0, "")

		pdf.Image(fileName, 175, 250, 25, 25, false, "", 0, "")
		pdf.Ln(-1)
		pdf.CellFormat(0, 9, "Fecha: "+strings.Split(time.Now().String(), " ")[0], "", 0, "R", false, 0, "")

		pdf.SetY(-15)
		// Arial italic 8
		pdf.SetFont("Arial", "I", 8)
		// Text color in gray
		pdf.SetTextColor(128, 128, 128)
		// Page number
		pdf.CellFormat(0, 10, fmt.Sprintf("Pagina %d", pdf.PageNo()),
			"", 0, "C", false, 0, "")
	})
	chapterTitle := func(chapNum int, titleStr string) {
		// 	// Arial 12
		pdf.SetTextColor(255, 255, 255)
		pdf.SetFont("Arial", "", 12)
		// Background color
		pdf.SetFillColor(0, 128, 255)
		// Title
		pdf.CellFormat(0, 6, fmt.Sprintf("%s", titleStr),
			"", 1, "C", true, 0, "")
		// Line break
		pdf.Ln(4)
	}

	restoration := func() {
		// Color and font restoration
		pdf.SetFillColor(224, 235, 255)
		pdf.SetTextColor(0, 0, 0)
		pdf.SetFont("", "", 0)
	}

	evalCaseDesviation := func(x, y float64, value, pvalue, prate float32) {

		lv := pvalue - (pvalue * prate)
		lr := pvalue + (pvalue * prate)

		var desviation int

		desviation = 0

		if value > lr {
			desviation = 1
		} else if value < lv {
			desviation = 2
		}

		switch desviation {
		case 0:
			pdf.SetFillColor(0, 250, 0) //green

		case 1:
			pdf.SetFillColor(255, 150, 0) //orange

		case 2:
			pdf.SetFillColor(220, 0, 0) //red

		}

		pdf.CellFormat(x, y, fmt.Sprintf("%.2f", value), "LR", 0, "C", true, 0, "")
		restoration()
	}

	fancyTable := func() {
		// Colors, line width and bold font
		pdf.SetFillColor(0, 128, 255)
		pdf.SetTextColor(255, 255, 255)
		pdf.SetDrawColor(0, 128, 255)
		pdf.SetLineWidth(.3)

		pdf.SetFont("Arial", "B", 7)
		// 	Header
		w := []float64{28, 5, 8, 52, 10, 8, 10, 18, 18, 10, 10, 10, 10, 10}

		wSum := 0.0
		for _, v := range w {
			wSum += v
		}

		left := (210 - wSum) / 2
		pdf.SetX(left)

		for j, str := range headerWeight {
			pdf.CellFormat(w[j], 7, str, "1", 0, "C", true, 0, "")
		}
		pdf.Ln(-1)
		restoration()
		// 	Data
		fill := false
		count := 0

		for _, c := range weightData {
			for index := 0; index < 7; index++ {

				pdf.SetX(left)
				pdf.CellFormat(w[0], 6, c.Date, "LR", 0, "C", fill, 0, "")
				pdf.CellFormat(w[1], 6, fmt.Sprintf("%d", c.Turn), "LR", 0, "C", fill, 0, "")
				pdf.CellFormat(w[2], 6, c.Line, "LR", 0, "C", fill, 0, "")
				pdf.CellFormat(w[3], 6, c.Pname, "LR", 0, "C", fill, 0, "")
				pdf.CellFormat(w[4], 6, fmt.Sprintf("%.2f", c.Pvalue), "LR", 0, "C", fill, 0, "")
				pdf.CellFormat(w[5], 6, c.Punit, "LR", 0, "C", fill, 0, "")
				pdf.CellFormat(w[6], 6, fmt.Sprintf("%.2f", c.Prate), "LR", 0, "C", fill, 0, "")
				pdf.CellFormat(w[7], 6, c.Coordinator+" "+c.Clname, "LR", 0, "C", fill, 0, "")
				pdf.CellFormat(w[8], 6, c.Operator+" "+c.Olname, "LR", 0, "C", fill, 0, "")
				//---------------------EVAL----CASE--DEsVIATIOn
				evalCaseDesviation(w[9], 6, c.V1, c.Pvalue, c.Prate)
				evalCaseDesviation(w[10], 6, c.V2, c.Pvalue, c.Prate)
				evalCaseDesviation(w[11], 6, c.V2, c.Pvalue, c.Prate)
				evalCaseDesviation(w[12], 6, c.V4, c.Pvalue, c.Prate)
				evalCaseDesviation(w[13], 6, c.V5, c.Pvalue, c.Prate)

				pdf.Ln(-1)
				fill = !fill
				count++
			}
		}
		fmt.Printf("Count:\t%d", count)
		pdf.SetX(left)
		pdf.CellFormat(wSum, 0, "", "T", 0, "", false, 0, "")
	}

	chapterBody := func() {
		pdf.SetY(60)
		fancyTable()
		//line break
		pdf.Ln(-1)
		// Mention in italics
		pdf.SetFont("", "I", 0)
		pdf.CellFormat(0, 5, "(terminan datos...)", "", 0, "R", false, 0, "")
	}

	printChapter := func(chapNum int, titleStr, fileStr string) {
		pdf.AddPage()
		pdf.Ln(12)
		chapterTitle(chapNum, titleStr)
		chapterBody()
	}

	printChapter(1, "Reporte Pesos", "texto xD...")
	//printChapter(2, "THE PROS AND CONS", example.TextFile("20k_c2.txt"))
	fileStr := "public/docs/report" + date + ".pdf"
	errp := pdf.OutputFileAndClose(fileStr)
	//example.Summary(errp, fileStr)
	if errp != nil {
		panic(errp.Error())
	}

	/*-----------------------*/

	http.Redirect(w, r, fileStr, http.StatusFound)
	// Output:
	// Successfully generated pdf/Fpdf_MultiCell.pdf
}

func FileHandlerEngine(r *http.Request, formReference, folder, patron string) string {
	var filePath string

	/*------------------------------------------------------------------*/
	fmt.Println("File Upload Endpoint Hit")

	// Parse our multipart form, 10 << 20 specifies a maximum
	// upload of 10 MB files.
	r.ParseMultipartForm(10 << 20)
	// FormFile returns the first file for the given key `myFile`
	// it also returns the FileHeader so we can get the Filename,
	// the Header and the size of the file
	file, handler, err := r.FormFile(formReference)
	if err != nil {
		fmt.Println("Error Retrieving the File")
		fmt.Println(err)
		return "public/images/ul_white.png" //default image
	}
	defer file.Close()
	fmt.Printf("Uploaded File: %+v\n", handler.Filename)
	fmt.Printf("File Size: %+v\n", handler.Size)
	fmt.Printf("MIME Header: %+v\n", handler.Header)

	// Create a temporary file within our temp-images directory that follows
	// a particular naming pattern
	tempFile, err := ioutil.TempFile(folder, patron)
	if err != nil {
		fmt.Println(err)
	}

	defer tempFile.Close()

	// read all of the contents of our uploaded file into a
	// byte array
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println(err)
	}
	// write this byte array to our temporary file
	tempFile.Write(fileBytes)
	// return that we have successfully uploaded our file!
	fmt.Printf("Successfully Uploaded File\n")

	/*------------------------------------------------------------------*/

	filePath = tempFile.Name()

	return filePath
}
