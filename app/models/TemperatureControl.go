package models

type TemperatureControl_holder struct {
	Profile_picture      string
	Date                 string
	Turn                 int
	Line                 string
	Pname                string
	Nname                string
	Fname                string
	Lname                string
	Batch                int
	Psi                  float32
	Exchange_temperature float32
	Hopper_temperature   float32
	Fill_temperature     float32

	PSI_bottom         float32
	PSI_top            float32
	Interchange_bottom float32
	Interchange_top    float32
	Hopper_bottom      float32
	Hopper_top         float32
	Fill_bottom        float32
	Fill_top           float32

	Observation string
}
