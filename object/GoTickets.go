package object

import (
	"fmt"
	"github.com/Deansquirrel/goGenerator/interface"
	"github.com/kataras/iris/core/errors"
)

type goTickets struct {
	total    uint32
	ticketCh chan struct{}
	active   bool
}

func NewGoTickets(total uint32) (_interface.IGoTickets, error) {
	gt := goTickets{}
	if !gt.init(total) {
		errMsg := fmt.Sprintf("The goroutine ticket pool can not be initialized!(total=%d)\n", total)
		return nil, errors.New(errMsg)
	}
	return &gt, nil
}

func (gt *goTickets) init(total uint32) bool {
	if gt.active {
		return false
	}
	if total == 0 {
		return false
	}
	ch := make(chan struct{}, total)
	for i := 0; i < int(total); i++ {
		ch <- struct{}{}
	}
	gt.ticketCh = ch
	gt.total = total
	gt.active = true
	return true
}

func (gt *goTickets) Take() {
	<-gt.ticketCh
}

func (gt *goTickets) Return() {
	gt.ticketCh <- struct{}{}
}

func (gt *goTickets) Active() bool {
	return gt.active
}

func (gt *goTickets) Total() uint32 {
	return gt.total
}

func (gt *goTickets) Remainder() uint32 {
	return uint32(len(gt.ticketCh))
}
