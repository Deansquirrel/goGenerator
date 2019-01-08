package _interface

import (
	"github.com/Deansquirrel/goGenerator/object"
	"time"
)

type ICaller interface {
	BuildReq() object.RawReq
	Call(req []byte,timeoutNS time.Duration)([]byte,error)
	CheckResp(rawReq object.RawReq,rawResp object.RawResp) *object.CallResult
}
