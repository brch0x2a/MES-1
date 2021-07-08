package models

type VOT struct{
	GoodVolume float32
	NominalSpeed float32
	Box_amount float32 
}

func (vot VOT) Get() float32{
	return (vot.GoodVolume *  vot.Box_amount) / vot.NominalSpeed
}

