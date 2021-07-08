package models

type MetaWork_order struct {
	Id           int
	Job          string
	Color        string
	Description  string
	Init         string
	End          string
	NickName     string
	Fname        string
	Lname        string
	State        string
	Line         string
	Photo_before string
	Photo_after  string
	Phase        string
}

type MetaLogWork_order struct {
	Id           int
	Job          string
	Color        string
	Description  string
	Init         string
	End          string
	Profile      string
	Fname        string
	Lname        string
	State        string
	SColor       string
	Line         string
	Photo_before string
	Photo_after  string

	ActualBegin string

	Wait_time       string
	Diagnostic_time string
	Stock_time      string
	Repair_time     string
	Test_time       string
	Delivery_time   string

	ActualEnd string

	Phase string
	Note  string
}

type WorkSaturation struct {
	Mecanic User_holder
	Times   []float32
}
