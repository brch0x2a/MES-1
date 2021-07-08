package models

type Product struct {
	Id                 int
	Sku                int
	Name               string
	Photo              string
	PSI_bottom         float32
	PSI_top            float32
	Bares_bottom       float32
	Bares_top          float32
	Lung_bottom        float32
	Lung_top           float32
	Interchange_bottom float32
	Interchange_top    float32
	Hopper_bottom      float32
	Hopper_top         float32
	Fill_bottom        float32
	Fill_top           float32
}

//
