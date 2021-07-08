package models

type Jaw_teflon_ultrasonic_state_holder struct {
	Profile_picture		 string
	Date  string
	Turn  int
	Line  string
	Pname string
	Nname string
	Fname string
	Lname string

	J1  float32
	J2  float32
	J3  float32
	J4  float32
	J5  float32
	J6  float32
	J7  float32
	J8  float32
	J9  float32
	J10 float32
	J11 float32
	J12 float32

	Jaw_state    string
	Teflon_state string

	Ultrasonic_time      float32
	Ultrasonic_amplitude float32
	Ultrasonic_pressure  float32
}
