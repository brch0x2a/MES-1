package models

type LossTree struct {
	OEE    float32
	Losses []TimeSegment
}

type LossTreeArea struct {
	S Sub_classification_holder
	T []TimeSegment
}
