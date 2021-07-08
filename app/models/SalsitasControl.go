package models

type SalsitasControl_subheader struct {
	Id           int
	Date         string
	Turn         int
	Line         int
	Presentation int
	Coordinator  int
	Operator     int
	Header       int
}

type SalsitasControl_subheader_holder struct {
	Id           int
	Date         string
	Turn         int
	Line         string
	Presentation string
	Pvalue       float32
	Punit        string
	Perror       float32
	Coordinator  string
	Operator     string
	Header       int
}

type Consolidated_weight_salsitas struct {
	Profile_picture string
	Date        string
	Turn        int
	Line        string
	Pname       string
	Pvalue      float32
	Punit       string
	Prate       float32
	Coordinator string
	Clname      string
	Operator    string
	Olname      string
	V1          float32
	V2          float32
	V3          float32
	V4          float32
	V5          float32
}

func (c Consolidated_weight_salsitas) HasDeviation(value float32) bool {
	lv := c.Pvalue - (c.Pvalue * c.Prate)
	lr := c.Pvalue + (c.Pvalue * c.Prate)

	if value > lr || value < lv {
		return true
	}
	return false
}
