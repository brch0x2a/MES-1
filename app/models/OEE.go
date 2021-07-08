package models

type OEE struct {
}

type DayOEE struct {
	Day string
	OEE float32
}

type WeekOEE struct {
	Line    string
	DayInfo []DayOEE
	Total   float32
}

type RelativeOEE struct {
	OEE      float32
	Dtime    float32
	StopTime float32
}

type OEEProjection struct {
	Line        string
	CurrentCase float32
	WorseCase   float32
	AverageCase float32
	BestCase    float32
}

type OEESaturation struct {
	Line      string
	Diff      float32
	Allocated float32
}
