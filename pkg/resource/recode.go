package resource

import (
	"bytes"
	"net"
	"strconv"

	"github.com/kseleven/haul/pkg/util"
)

type Recode struct {
	Name     string `json:"name"`     // variable bits
	Type     QType  `json:"type"`     // 16 bits
	Class    QClass `json:"class"`    // 16 bits
	TTL      uint32 `json:"ttl"`      // 32 bits
	RdLength uint16 `json:"rdLength"` // 16 bits
	Rdata    string `json:"rdata"`    // variable bits
}

func (r Recode) Decode(data, src []byte) (Recode, []byte) {
	name, usedLength := util.ParseDomainName(data, src)
	r.Name = name
	var buf bytes.Buffer
	data = data[usedLength:]
	buf.Grow(len(data))
	buf.Write(data)
	r.Type = QType(util.BytesToUint16(buf.Next(2)))
	r.Class = QClass(util.BytesToUint16(buf.Next(2)))
	r.TTL = util.BytesToUint32(buf.Next(4))
	r.RdLength = util.BytesToUint16(buf.Next(2))
	r.Rdata = readRdata(buf.Next(int(r.RdLength)), r.Type, src)
	return r, buf.Bytes()
}

func readRdata(data []byte, qType QType, src []byte) string {
	switch qType {
	case QTypeA, QTypeAAAA:
		return net.IP(data).String()
	case QTypeCNAME, QTypeMX, QTypeTXT, QTypeNS:
		name, _ := util.ParseDomainName(data, src)
		return name
	case QTypeSOA:
		return util.ParseSOARdata(data, src)
	default:
		return "unknown type:" + qType.String()
	}
}

func (r Recode) String() string {
	var buf bytes.Buffer
	buf.WriteString(r.Name)
	buf.WriteString("\t\t")
	buf.WriteString(strconv.FormatUint(uint64(r.TTL), 10))
	buf.WriteString("\t")
	buf.WriteString(r.Class.String())
	buf.WriteString("\t")
	buf.WriteString(r.Type.String())
	buf.WriteString("\t")
	buf.WriteString(r.Rdata)
	return buf.String()
}
