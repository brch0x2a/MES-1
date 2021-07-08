package models

type Clean_filter struct {
	Id   int
	Name string
}

type Washing_stage struct {
	Id   int
	Name string
}

type Clean_corrective_action struct {
	Id          int
	Description string
	Bottom      float32
	Top         float32
}

type Clean_Disinfection struct {
	Id int

	Date_init string
	Date_end  string

	Washing_stage string

	Equipment string
	Comment   string

	Washed_type string

	Detergent    float32
	Disinfectand float32
	Chemical     float32
	Foam         float32
	Spray        float32

	Filter   float32
	Water_ph float32

	Visual_inspection bool
	Microbiology      bool

	Atp               float32
	Corrective_action int
	New_atp           float32

	Allergen_state string

	Maintenance string

	Autor    string
	Approver string
}

type Clean_DisinfectionMeta struct {
	Id int

	Date_init string

	RProfile_Picture string
	RFname           string
	RLname           string

	Line string

	Washing_stage string

	Equipment string

	Washed_type string

	Detergent    float32
	Disinfectand float32
	Chemical     float32
	Foam         float32
	Spray        float32

	Filter   string
	Water_ph float32

	Visual_inspection string
	Microbiology      string

	Atp               float32
	Corrective_action int
	New_atp           float32

	Allergen_state string

	Maintenance string

	Comment string

	AProfile_Picture string
	AFname           string
	ALname           string

	Date_end string

	State string
}
