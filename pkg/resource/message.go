package resource

import (
	"bytes"

	"github.com/kseleven/haul/pkg/util"
)

type Message struct {
	Header
	Question
	Answers    []Recode
	Authority  []Recode
	Additional []Recode
}

const (
	HeadLength = 12
)

func NewMessage(qName string, qType QType) *Message {
	m := &Message{
		Header: Header{
			ID:      util.GenRequestId(),
			QR:      0,
			Opcode:  0,
			AA:      0,
			TC:      0,
			RD:      1,
			RA:      0,
			Z:       0,
			Rcode:   0,
			QdCount: 1,
			AnCount: 0,
			NsCount: 0,
			ArCount: 0,
		},
		Question: Question{
			QName:  qName,
			QType:  qType,
			QClass: QClassIN,
		},
	}
	return m
}

func (m Message) Encode() []byte {
	var buf bytes.Buffer
	buf.Write(m.Header.Encode())
	buf.Write(m.Question.Encode())
	return buf.Bytes()
}

func (m *Message) Decode(srcData []byte) {
	var buf bytes.Buffer
	buf.Grow(len(srcData))
	buf.Write(srcData)

	m.Header = m.Header.Decode(buf.Next(HeadLength))
	if m.Header.QdCount > 0 {
		question, data := Question{}.Decode(buf.Bytes(), srcData)
		m.Question = question
		buf.Reset()
		buf.Write(data)
	}

	m.Answers = make([]Recode, m.Header.AnCount)
	for i := 0; i < int(m.Header.AnCount); i++ {
		r, data := Recode{}.Decode(buf.Bytes(), srcData)
		m.Answers[i] = r
		buf.Reset()
		buf.Write(data)
	}

	m.Authority = make([]Recode, m.Header.NsCount)
	for i := 0; i < int(m.Header.NsCount); i++ {
		r, data := Recode{}.Decode(buf.Bytes(), srcData)
		m.Authority[i] = r
		buf.Reset()
		buf.Write(data)
	}
}
