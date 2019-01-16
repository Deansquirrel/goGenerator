package worker

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Deansquirrel/goGenerator/common"
	"github.com/Deansquirrel/goGenerator/object"
	"sync/atomic"
	"time"
)

const (
	StatusOriginal uint32 = 0
	StatusSarting  uint32 = 1
	StatusStarted  uint32 = 2
	StatusStopping uint32 = 3
	StatusStopped  uint32 = 4
)

type generator struct {
	resultCh    chan *object.CallResult
	timeoutNS   time.Duration
	durationNS  time.Duration
	concurrency uint32
	tickets     object.IGoTickets
	lps         uint32
	ctx         context.Context
	cancelFunc  context.CancelFunc
	caller      object.ICaller
	callCount   int64
	status      uint32
}

func NewGenerator(
	caller object.ICaller,
	timeoutNS time.Duration,
	lps uint32,
	durationNS time.Duration,
	resultCh chan *object.CallResult) (generator, error) {
	gen := &generator{
		caller:     caller,
		timeoutNS:  timeoutNS,
		lps:        lps,
		durationNS: durationNS,
		resultCh:   resultCh,
	}

	if timeoutNS.Seconds() < 1 {
		return generator{}, errors.New("超时时间不能小于1s")
	}
	atomic.StoreUint32(&gen.status, StatusOriginal)
	gen.concurrency = uint32(gen.timeoutNS.Seconds() * float64(lps))
	gen.ctx, gen.cancelFunc = context.WithTimeout(
		context.Background(), gen.durationNS)
	return *gen, nil
}

func (gen *generator) Start() {
	gen.callCount = 0
	atomic.StoreUint32(&gen.status, StatusStarted)
}

func (gen *generator) genLoad(throttle <-chan time.Time) {
	for {
		select {
		case <-gen.ctx.Done():
			gen.prepareToStop(gen.ctx.Err())
			return
		default:
		}
		gen.asyncCall()
		if gen.lps > 0 {
			select {
			case <-throttle:
			case <-gen.ctx.Done():
				gen.prepareToStop(gen.ctx.Err())
				return
			}
		}
	}
}

func (gen *generator) prepareToStop(ctxError error) {
	atomic.CompareAndSwapUint32(&gen.status, StatusStarted, StatusStopping)
	close(gen.resultCh)
	atomic.StoreUint32(&gen.status, StatusStopped)
}

func (gen *generator) asyncCall() {
	gen.tickets.Take()
	go func() {
		defer func() {
			p := recover()
			if p != nil {
				err, ok := interface{}(p).(error)
				var errMsg string
				if ok {
					errMsg = fmt.Sprintf("Async Call Panic!(error:%s)", err)
				} else {
					errMsg = fmt.Sprintf("Async Call Panic!(clue:%#v)", p)
				}
				common.PrintAndLog(errMsg)
				result := &object.CallResult{
					ID: -1,
					//Code:
					Msg: errMsg,
				}
				gen.sendResult(result)
			}
			gen.tickets.Return()
		}()

		var rawReq object.RawReq

		var callStatus uint32
		//0-未调用或调用中，1-调用完成，2-调用超时
		timer := time.AfterFunc(gen.timeoutNS, func() {
			if atomic.CompareAndSwapUint32(&callStatus, 0, 2) {
				return
			}
			result := object.CallResult{
				ID:  rawReq.ID,
				Req: rawReq,
				//Code:
				Msg:    fmt.Sprintf("Timeout!(expected: < %v", gen.timeoutNS),
				Elapse: gen.timeoutNS,
			}
			gen.sendResult(&result)
		})
		rawResp := gen.callOne(&rawReq)
		if !atomic.CompareAndSwapUint32(&callStatus, 0, 1) {
			return
		}
		timer.Stop()

		var result *object.CallResult
		if rawResp.Err != nil {
			result = &object.CallResult{
				ID:  rawResp.ID,
				Req: rawReq,
				//Code:
				Msg:    rawResp.Err.Error(),
				Elapse: rawResp.Elapse,
			}
		} else {
			result = gen.caller.CheckResp(rawReq, *rawResp)
			result.Elapse = rawResp.Elapse
		}
		gen.sendResult(result)
	}()
}

func (gen *generator) callOne(rawReq *object.RawReq) *object.RawResp {
	return nil
}

func (gen *generator) sendResult(result *object.CallResult) bool {
	if atomic.LoadUint32(&gen.status) != StatusStarted {
		gen.printIgnoredResult(result, "stopped load generator")
		return false
	}
	select {
	case gen.resultCh <- result:
		return true
	default:
		gen.printIgnoredResult(result, "full result channel")
		return false
	}
}

func (gen *generator) printIgnoredResult(result *object.CallResult, reason string) {
	msg := reason + " - "
	val, err := json.Marshal(&result)
	if err != nil {
		msg = msg + err.Error()
	} else {
		msg = msg + string(val)
	}
	common.PrintAndLog(msg)
}
