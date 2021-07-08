package models

import "database/sql"

/*
Agrupa lo relacionando a boletas verdes
*/
type Class_of_event struct {
	Id   int
	Name string
}

type Event_cause struct {
	Id   int
	Name string
}

type Frequency_catalog struct {
	Id   int
	Name string
}

type Severity_catalog struct {
	Id   int
	Name string
}

type SHE_standard_catalog struct {
	Id   int
	Name string
}

type SHE_tag_green_holder struct {
	Id          int
	Line        string
	User        string
	Fname       string
	Lname       string
	RequestDate string
	CloseDate   sql.NullString
	Type        string
	Color       string
	Priority    string
	Equipment   string

	Class_of_event     string
	Event_cause        string
	Anomaly            string
	Frequency_catalog  string
	Investigation      string
	Severity_catalog   string
	InWeb              string
	Correction_action  string
	Suggestion         string
	Lesion_description string
	Damage_description string

	SHE_standard_catalog string
	State                string
	SColor               string
}
