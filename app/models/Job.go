package models

type Job struct {
	Date         string
	Turn         int64
	Line         int64
	Presentation int64
	Coordinator  int64
	Config       bool
}
