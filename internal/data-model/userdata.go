package model

type Data struct {
	ID            uint
	Pupil         string
	Establishment string
	Subject       SubjectType
	KnowlegdeTest KnowlegdeTestType
	Grade         GradeType
}

func (ud *Data) HashKey() uint {
	return ud.ID
}

type SubjectType string

const (
	Russian     SubjectType = "Russian"
	Mathematics             = "Mathematics"
	Physics                 = "Physics"
	Literature              = "Literature"
	English                 = "English"
	History                 = "History"
	Technology              = "Technology"
)

type KnowlegdeTestType string

const (
	Annual  KnowlegdeTestType = "Annual test"
	Quarter                   = "Quarter test"
	Test                      = "Test"
	Work                      = "Independent work"
	Board                     = "Answer to the board"
)

type GradeType uint8

const (
	Two   GradeType = 2
	Three           = 3
	Four            = 4
	Five            = 5
)
