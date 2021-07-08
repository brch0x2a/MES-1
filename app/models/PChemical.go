package models

type Physicochemical_subheader struct {
	Id          int
	Boula_fryma string
	Date        string
	Observation string
	Product     int
	Area        int
	Header      int
}

type Physicochemical_subheader_holder struct {
	Id          int
	Boula_fryma string
	Product     string
	Area        string
	Header      int
}

type Sensorial_analysis_scale struct {
	Id          int
	Description string
}

type Desition_catalog struct {
	Id          int
	Description string
}

type Physicochemical_general struct {
	Id                   int
	Batch                int
	Tank                 string
	Presentation_density float32
	Date                 string

	/*Analisis Fisico-Quimico*/
	Brix                    float32
	Ph                      float32
	Ph_pcc                  float32
	Chloride                float32
	Consistency_plummet     float32
	Consistency_homogenizer float32

	/*Analisis Sensorial*/
	Appearance int
	Color      int
	Aroma      int
	Taste      int

	Analyst   int
	Desition  int
	Subheader int
}

type Physicochemical_general_holder struct {
	Profile_picture		 string
	Area string

	Analystnn    string
	Analystfname string
	Analystlname string

	Date    string
	Product string

	Fryma float32
	Batch int

	/*Analisis Sensorial*/
	Appearance_level int
	Appearance       string
	Color_level      int
	Color            string
	Aroma_level      int
	Aroma            string
	Taste_level      int
	Taste            string

	/*Analisis Fisico-Quimico*/
	Ph                      float32
	Ph_pcc                  float32
	Chloride                float32
	Brix                    float32
	Consistency_plummet     float32
	Consistency_homogenizer float32
	Presentation_density    float32

	Tank string
}
