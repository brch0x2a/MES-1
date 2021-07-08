package models

type AMIS_holder struct {
	Line       string
	GoodVolume float32

	TotalTime float32

	LegalLosses float32

	PDL        float32
	PDLSegment []TimeSegment

	MPL        float32
	MPLSegment []TimeSegment

	UCL        float32
	UCLSegment []TimeSegment

	AT       float32 //Aveable time
	ALT      float32 //Aveable Loading Time
	OT       float32
	LT       float32
	IdleTime float32
	VOT      float32
	OEE      float32
	MP       float32
	SP       float32
	AU       float32
	UCU      float32 //unconstraint capacity utilization
	CCU      float32 //Constraint capacity Utilization

	Breakdowns         int
	GeneralChangeovers int
	FormatChangeovers  int
	ProductChangeovers int

	GeneralChangeoversTime float32
	FormatChangeoversTime  float32
	ProductChangeoversTime float32

	OR float32
}

type TimeSegment struct {
	Line            string
	LineId          int
	SubCategoryName string
	Color           string
	TotalMinutes    float32
}
