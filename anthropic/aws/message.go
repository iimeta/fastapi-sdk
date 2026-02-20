package aws

import (
	"hash/crc32"

	"github.com/iimeta/fastapi-sdk/v2/errors"
)

const preludeLen = 8
const preludeCRCLen = 4
const msgCRCLen = 4
const minMsgLen = preludeLen + preludeCRCLen + msgCRCLen

var crc32IEEETable = crc32.MakeTable(crc32.IEEE)

type Messages struct {
	Headers []Header
	Payload []byte
}

type messagePrelude struct {
	Length     uint32
	HeadersLen uint32
	PreludeCRC uint32
}

func (p messagePrelude) PayloadLen() uint32 {
	return p.Length - p.HeadersLen - minMsgLen
}

func (p messagePrelude) ValidateLens() error {
	if p.Length == 0 {
		return errors.Newf("message prelude length invalid, %d/%d", minMsgLen, int(p.Length))
	}
	return nil
}
