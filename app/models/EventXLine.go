package models

type EventXLine struct {
	Id         int
	Id_event   int
	Id_line    int
	Date_event string
	Turn       int
	Minutes    int
	Note       string
}

/*
bathes
reportes de despertidicios de laminano, la dos no

monitoriando las preciones

salida d pulmon,

conteo de minutos acitva y en paro, hora y fecha

por turno cuanto tiempo trabajando,

red covertura de mayonesa, frijoles, fibra optica,
M3
*/

type EventXLineV00 struct {
	Id         int
	Sub        string
	Color      string
	Branch     string
	Event      string
	Date_event string
	Nick_name string
	Fname     string
	Lname     string
	Profile_picture string
	Turn       int
	Minutes    int
	Note       string
}

type EventXLineV01 struct {
	Date_event string
	Nick_name string
	Fname     string
	Lname     string
	Profile_picture string
	Turn       int
	Line       string
	Sub        string
	Branch     string
	Event      string
	Minutes    int
	Note       string
	Color      string
}

type EventXLine_holder struct {
	Id         int
	Sub        string
	Branch     string
	Event      string
	Line       string
	Date_event string
	Turn       int
	Minutes    int
	Note       string
	Color      string
}
