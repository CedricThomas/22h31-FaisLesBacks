package model

type noSuchEntityErr struct {
}

func (noSuchEntityErr) Error() string {
	return "unable to find the requested entity"
}

var NoSuchEntity = noSuchEntityErr{}
