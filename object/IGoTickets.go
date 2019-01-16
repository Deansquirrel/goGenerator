package object

type IGoTickets interface {
	Take()
	Return()
	Active() bool
	Total() uint32
	Remainder() uint32
}
