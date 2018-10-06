package internal

type Method interface {
	Process() (string, error)
}
