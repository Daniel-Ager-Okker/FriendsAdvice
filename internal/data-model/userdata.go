package model

type Data struct {
	ID            uint
	Pupil         string
	Establishment string
	Class         ClassType
	Letter        LetterType
}

func (ud *Data) HashKey() uint {
	return ud.ID
}

type LetterType string

const (
	LetterA LetterType = "A"
	LetterB LetterType = "B"
	LetterV LetterType = "V"
)

type ClassType uint8

const (
	One    ClassType = 1
	Two    ClassType = 2
	Three  ClassType = 3
	Four   ClassType = 4
	Five   ClassType = 5
	Six    ClassType = 6
	Seven  ClassType = 7
	Eight  ClassType = 8
	Nine   ClassType = 9
	Ten    ClassType = 10
	Eleven ClassType = 11
)
