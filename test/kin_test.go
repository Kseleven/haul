package test

import (
	"encoding/hex"
	"github.com/kseleven/haul/pkg/resource"
	"github.com/kseleven/haul/pkg/service"
	"github.com/kseleven/haul/pkg/util"
	"testing"
)

func BenchmarkGenRequestId(b *testing.B) {
	for i := 0; i < 10; i++ {
		b.Log(util.GenRequestId())
	}
}

func TestParseUint16(t *testing.T) {
	var id uint16 = 65535
	t.Logf("id %b", id)
	result := util.Uint16ToBytes(id)
	t.Log(result)
	t.Log(util.BytesToUint16(result))

	h := resource.Header{QR: 1, Opcode: resource.QrcodeIQuery}
	one := h.QR<<3 | byte(h.Opcode>>1)
	t.Log(one)
}

func TestDomainName(t *testing.T) {
	name := "www.google.com"
	result := util.DomainNameToBytes(name)
	t.Log(result)
	t.Log(hex.Dump(result))
}

func TestMessage(t *testing.T) {
	message := resource.Message{}
	message.Header = resource.Header{
		ID:      0x3fc8,
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
	}
	message.Question = resource.Question{
		QName:  "www.google.com",
		QType:  resource.QTypeA,
		QClass: resource.QClassIN,
	}

	encode := message.Encode()
	t.Log("encode", encode)
	t.Log("hex encode", hex.Dump(encode))
}

func TestParseDomainName(t *testing.T) {
	data := []byte{3, 0x77, 0x77, 0x77, 6, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 3, 0x63, 0x6f, 0x6d, 0}
	t.Log(util.ParseDomainName(data, data))
}

func TestIsPointer(t *testing.T) {
	var length byte = 192
	t.Log(util.IsLabelPointer(length))
}

func TestRequest(t *testing.T) {
	r := resource.Request{
		Host:  "114.114.114.114",
		Port:  53,
		QName: ".",
		QType: resource.QTypeNS,
	}
	if err := service.Request(r); err != nil {
		t.Errorf(err.Error())
	}
}
