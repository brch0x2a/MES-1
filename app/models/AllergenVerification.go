package models

type AllergenVerification_holder struct {
	Profile_picture		 string

	Id int
	Date_reg string
	Turn int
	Line string
	Presentation string
	Nname string
	Fname string
	Lname string

	Batch string 
	Lote string
	Coil_weight float32

	Lote_plug string
	Lote_top string
	Lote_gallon string
		 
	Front_laminate string
	Allergen_laminate string

	Milk string
	Soya string
	Gluten string
	Egg string

}