package worker

import (
	"context"
	"errors"
	"github.com/Deansquirrel/goGenerator/interface"
	"github.com/Deansquirrel/goGenerator/object"
	"time"
)

const (
	StatusOriginal uint32 = 0
	StatusSarting uint32 = 1
	StatusStarted uint32 = 2
	StatusStoping uint32 = 3
	StatusStoped uint32 = 4
)

type generator struct {
	resultCh chan * object.CallResult
	timeoutNS time.Duration
	durationNS time.Duration
	concurrency uint32
	tickets _interface.IGoTickets
	lps uint32
	ctx context.Context
	cancelFunc context.CancelFunc
	caller _interface.ICaller
	callCount int64
	status uint32
}

func NewGenerator(
	caller _interface.ICaller,
	timeoutNS time.Duration,
	lps uint32,
	durationNS time.Duration,
	resultCh chan * object.CallResult)(generator,error){
	gen := &generator{
		caller:caller,
		timeoutNS:timeoutNS,
		lps:lps,
		durationNS:durationNS,
		resultCh:resultCh,
	}

	if timeoutNS.Seconds() < 1 {
		return generator{},errors.New("超时时间不能小于1s")
	}
	gen.concurrency = uint32(gen.timeoutNS.Seconds() * float64(lps))
	return *gen,nil
}