package models

type Presentation struct {
	Id           int
	Id_product   int
	Name         string
	Weight_unit  string
	Weight_value float32
	Error_rate   float32
	Box_amount   int
	Photo        string
}

//

type Presentation_holder struct {
	Id           int
	Product      string
	Name         string
	Weight_unit  string
	Weight_value float32
	Error_rate   float32
	Box_amount   int
	Photo        string
}
