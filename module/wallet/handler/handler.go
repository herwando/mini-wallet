package handler

import (
	"github.com/herwando/mini-wallet/lib/common/writer"
)

var (
	writerWriteStrOK        = writer.WriteStrOK
	writerWriteData         = writer.WriteOK
	writerWriteDataAccepted = writer.WriteAccepted
	writerWriteJSONAPIError = writer.WriteJSONAPIError
)

type BasicResponse struct {
	Data   interface{} `json:"data"`
	Status string      `json:"status"`
}
