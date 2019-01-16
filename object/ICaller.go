package object

import (
	"time"
)

type ICaller interface {
	BuildReq() RawReq
	Call(req []byte, timeoutNS time.Duration) ([]byte, error)
	CheckResp(rawReq RawReq, rawResp RawResp) *CallResult
}
