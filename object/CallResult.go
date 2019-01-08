package object

import "time"

type CallResult struct {
	ID int64
	Req RawReq
	Resp RawResp
	//Code RetCode
	Msg string
	Elapse time.Duration
}
