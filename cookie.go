package goKLC

type Cookie interface {
	Create(name string, value string, duration int, path string)
	GetName() string
	GetValue() string
	GetDuration() int
	GetPath() string
}
