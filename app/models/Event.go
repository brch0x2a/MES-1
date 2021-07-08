package models

type Event struct{
	Id int
	Description string
	Id_branch int
}

type Event_holder struct{
	Id int
	LTC string //Line time classification
	Sub string //Sub classification
	Branch string 
	Description string
	Color string
}

type Event_mini_holder struct{
	Id int
	Description string
	//Id_branch int
	Color string

}