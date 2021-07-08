package models

type Planning struct {
	Id              int
	Planned         int
	Produced        int
	Turn            int
	Date_planning   string
	Nominal_speed   float32
	Version         int
	Id_presentation int
	Id_line         int
}

type PlanningV00 struct {
	Planned  int
	Produced int
}

type Planning_holder struct {
	Id            int
	Date_planning string
	Turn          int
	Line          string
	Presentation  string
	Version       int
	Planned       int
	Produced      int
	Nominal_speed float32
	Box_amount    int
	Photo         string
}

type Planning_insert struct {
	Id            int
	Date_planning string
	Turn          string
	Line          string
	Presentation  string
	Version       string
	Planned       string
	Produced      string
	Nominal_speed string
}

type ActualPlanning struct {
	Id           int
	Date_reg     string
	Line         string
	Presentation string
	Box          float32
	Profile      string
	Fname        string
	Lname        string
}
