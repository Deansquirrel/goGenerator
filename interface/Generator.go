package _interface

type IGenerator interface {
	Start() bool
	Stop() bool
	Status() uint32
	CallCount() int64
}
