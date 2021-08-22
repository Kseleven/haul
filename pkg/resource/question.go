package resource

import (
	"bytes"
	"github.com/kseleven/haul/pkg/util"
)

type Question struct {
	QName  string `json:"qName"`  //variable bit
	QType  QType  `json:"qType"`  //16 bit
	QClass QClass `json:"qClass"` //16 bit
}

func (q Question) Encode() []byte {
	var buf bytes.Buffer
	buf.Write(util.DomainNameToBytes(q.QName))
	buf.Write(util.Uint16ToBytes(uint16(q.QType)))
	buf.Write(util.Uint16ToBytes(uint16(q.QClass)))
	return buf.Bytes()
}

func (q Question) Decode(data, src []byte) (Question, []byte) {
	name, usedLength := util.ParseDomainName(data, src)
	q.QName = name
	data = data[usedLength:]
	var buf bytes.Buffer
	buf.Grow(len(data))
	buf.Write(data)
	q.QType = QType(util.BytesToUint16(buf.Next(2)))
	q.QClass = QClass(util.BytesToUint16(buf.Next(2)))
	return q, buf.Bytes()
}

func (q Question) String() string {
	var buf bytes.Buffer
	buf.WriteString("Question Section:\n")
	buf.WriteString(q.QName)
	buf.WriteString("\t\t\t")
	buf.WriteString(q.QClass.String())
	buf.WriteString("\t")
	buf.WriteString(q.QType.String())
	buf.WriteString("\n")
	return buf.String()
}

type QType uint16

func (q QType) String() string {
	for s, qType := range QTypeSet {
		if q == qType {
			return s
		}
	}
	return "unknown qtype"
}

var QTypeSet = map[string]QType{
	"A":     QTypeA,
	"NS":    QTypeNS,
	"MD":    QTypeMD,
	"MF":    QTypeMF,
	"CNAME": QTypeCNAME,
	"SOA":   QTypeSOA,
	"MB":    QTypeMB,
	"NULL":  QTypeNULL,
	"WKS":   QTypeWKS,
	"PTR":   QTypePRT,
	"HINFO": QTypeHINFO,
	"MINFO": QTypeMINFO,
	"MX":    QTypeMX,
	"TXT":   QTypeTXT,
	"AAAA":  QTypeAAAA,
	"AXFR":  QTypeAXFR,
	"MAILB": QTypeMAILB,
	"MAILA": QTypeMAILA,
	"ANY":   QTypeAny,
}

const (
	QTypeA     QType = 0x01
	QTypeNS    QType = 0x02
	QTypeMD    QType = 0x03
	QTypeMF    QType = 0x04
	QTypeCNAME QType = 0x05
	QTypeSOA   QType = 0x06
	QTypeMB    QType = 0x07
	QTypeMG    QType = 0x08
	QTypeMR    QType = 0x09
	QTypeNULL  QType = 0x0a
	QTypeWKS   QType = 0x0b
	QTypePRT   QType = 0x0c
	QTypeHINFO QType = 0x0d
	QTypeMINFO QType = 0x0e
	QTypeMX    QType = 0x0f
	QTypeTXT   QType = 0x10
	QTypeAAAA  QType = 0x1c
	QTypeAXFR  QType = 0xfc
	QTypeMAILB QType = 0xfd
	QTypeMAILA QType = 0xfe
	QTypeAny   QType = 0xff
)

type QClass uint16

func (q QClass) String() string {
	switch q {
	case QClassIN:
		return "IN"
	case QClassCS:
		return "CS"
	case QClassCH:
		return "CH"
	case QClassHS:
		return "HS"
	case QClassAny:
		return "ANY"
	default:
		return "unknown qClass"
	}
}

const (
	QClassIN  QClass = 0x01
	QClassCS  QClass = 0x02
	QClassCH  QClass = 0x03
	QClassHS  QClass = 0x04
	QClassAny QClass = 0xff
)
