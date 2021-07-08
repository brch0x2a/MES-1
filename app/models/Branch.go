package models

type Branch struct{
	Id int
	Description string
	Id_sub_classification int
}


type Branch_holder struct{
	Id int
	LTC string//line time classification
	Sub string//Subclassification
	Description string 
	Color string
}


