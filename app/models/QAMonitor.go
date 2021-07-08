package models

type QAMonitor struct {
	WeightCounter      int
	ProcessTemperature int
	JawTeflonState     int
	SealVerification   int
	CRQS               int
	BatchCounter       int
	AllergenCounter    int
	CleanDisinfection  int
	CodingVerification int
}
