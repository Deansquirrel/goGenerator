package object

import "time"

type RawResp struct {
	ID int64
	Resp []byte
	Err error
	Elapse time.Duration
}
