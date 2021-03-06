package resource

import (
	"bytes"
	"strconv"

	"github.com/kseleven/haul/pkg/util"
)

type Header struct {
	ID      uint16 `json:"id"`
	QR      byte   `json:"qr"`
	Opcode  Opcode `json:"opcode"`
	AA      byte   `json:"aa"`
	TC      byte   `json:"tc"`
	RD      byte   `json:"rd"`
	RA      byte   `json:"ra"`
	Z       byte   `json:"z"`
	Rcode   Rcode  `json:"rcode"`
	QdCount uint16 `json:"qdCount"`
	AnCount uint16 `json:"anCount"`
	NsCount uint16 `json:"nsCount"`
	ArCount uint16 `json:"arCount"`
}

func (h Header) Encode() []byte {
	var buf bytes.Buffer
	buf.Write(util.Uint16ToBytes(h.ID))
	buf.WriteByte(h.QR<<7 | byte(h.Opcode)<<3 | h.AA<<2 | h.TC<<1 | h.RD)
	buf.WriteByte(h.RA<<7 | h.Z<<4 | byte(h.Rcode))
	buf.Write(util.Uint16ToBytes(h.QdCount))
	buf.Write(util.Uint16ToBytes(h.AnCount))
	buf.Write(util.Uint16ToBytes(h.NsCount))
	buf.Write(util.Uint16ToBytes(h.ArCount))
	return buf.Bytes()
}

func (h Header) Decode(data []byte) Header {
	var buf bytes.Buffer
	buf.Grow(len(data))
	buf.Write(data)

	h.ID = util.BytesToUint16(buf.Next(2))
	oneByte := util.BytesToUint16(buf.Next(1))
	h.QR = byte(oneByte >> 7)
	h.Opcode = Opcode(oneByte >> 3 & 0xf)
	h.AA = byte(oneByte >> 2 & 0x1)
	h.TC = byte(oneByte >> 1 & 0x1)
	h.RD = byte(oneByte & 0x1)

	oneByte = util.BytesToUint16(buf.Next(1))
	h.RA = byte(oneByte >> 7)
	h.Z = byte(oneByte >> 4 & 0x7)
	h.Rcode = Rcode(oneByte & 0xf)

	h.QdCount = util.BytesToUint16(buf.Next(2))
	h.AnCount = util.BytesToUint16(buf.Next(2))
	h.NsCount = util.BytesToUint16(buf.Next(2))
	h.ArCount = util.BytesToUint16(buf.Next(2))

	return h
}

func (h Header) String() string {
	var buf bytes.Buffer
	buf.WriteString("Header Section:\n")
	buf.WriteString("opcode: ")
	buf.WriteString(h.Opcode.String())
	buf.WriteString(", status: ")
	buf.WriteString(h.Rcode.String())
	buf.WriteString(", id: ")
	buf.WriteString(strconv.FormatUint(uint64(h.ID), 10))
	buf.WriteString("\n")
	buf.WriteString("Flags:")
	buf.WriteString(" Query: ")
	buf.WriteString(strconv.FormatUint(uint64(h.QdCount), 10))
	buf.WriteString(", Answer: ")
	buf.WriteString(strconv.FormatUint(uint64(h.AnCount), 10))
	buf.WriteString(", Authority: ")
	buf.WriteString(strconv.FormatUint(uint64(h.NsCount), 10))
	buf.WriteString(", Additional: ")
	buf.WriteString(strconv.FormatUint(uint64(h.ArCount), 10))
	return buf.String()
}

const (
	QrQuery    byte = 0
	QrResponse byte = 1
)

type Opcode byte

func (o Opcode) String() string {
	switch o {
	case QrcodeQuery:
		return "query"
	case QrcodeIQuery:
		return "iQuery"
	case QrcodeStatus:
		return "status"
	case QrcodeTest:
		return "test"
	default:
		return "unknown opcode"
	}
}

const (
	QrcodeQuery  Opcode = 0
	QrcodeIQuery Opcode = 1
	QrcodeStatus Opcode = 2
	QrcodeTest   Opcode = 15
)

type Rcode byte

func (r Rcode) String() string {
	switch r {
	case RcodeNoError:
		return "noError"
	case RcodeFormatError:
		return "formatError"
	case RcodeServerFailure:
		return "serverFailure"
	case RcodeNameError:
		return "nameError"
	case RcodeNotImplemented:
		return "notImplemented"
	case RcodeRefused:
		return "refused"
	default:
		return "unknown rcode"
	}
}

const (
	RcodeNoError        Rcode = 0
	RcodeFormatError    Rcode = 1
	RcodeServerFailure  Rcode = 2
	RcodeNameError      Rcode = 3
	RcodeNotImplemented Rcode = 4
	RcodeRefused        Rcode = 5
)
