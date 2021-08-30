package util

import (
	"bytes"
	"crypto/rand"
	"fmt"
)

const (
	Dot      = 46
	RootName = "."
)

func Uint16ToBytes(data uint16) []byte {
	var buf bytes.Buffer
	buf.Grow(2)
	buf.WriteByte(byte(data >> 8))
	buf.WriteByte(byte(data))
	return buf.Bytes()
}

func BytesToUint16(data []byte) uint16 {
	if len(data) != 2 {
		return uint16(data[0])
	}
	return uint16(data[0])<<8 | uint16(data[1])
}

func BytesToUint32(data []byte) uint32 {
	if len(data) < 4 {
		return 0
	}
	return uint32(data[0])<<24 | uint32(data[1])<<16 | uint32(data[2])<<8 | uint32(data[3])
}

func DomainNameToBytes(name string) []byte {
	if name == RootName {
		return []byte{0}
	}
	tail := len(name) - 1
	if tail < 0 {
		return nil
	}
	var label []byte
	var buf bytes.Buffer
	for i, n := range name {
		if n != Dot {
			label = append(label, byte(n))
			if i != tail {
				continue
			}
		}
		buf.WriteByte(byte(len(label)))
		buf.Write(label)
		label = label[:0]
	}

	buf.WriteByte(0)
	return buf.Bytes()
}

func ParseDomainName(data []byte, src []byte) (string, int) {
	var (
		buf         bytes.Buffer
		nameBuf     bytes.Buffer
		labelLength byte
		i           byte
	)

	bufLength := len(data)
	buf.Grow(bufLength)
	buf.Write(data)
	if labelLength = buf.Next(1)[0]; IsLabelPointer(labelLength) {
		offset := buf.Next(1)[0]
		name, _ := ParseDomainName(src[offset:], src)
		return name, 2
	} else if IsRoot(labelLength) {
		return RootName, 1
	}
readBuf:
	for i = 0; i < labelLength; i++ {
		nameBuf.Write(buf.Next(1))
	}
	nameBuf.WriteByte(Dot)
	if labelLength = buf.Next(1)[0]; labelLength == 0 {
		return nameBuf.String(), bufLength - buf.Len()
	} else if IsLabelPointer(labelLength) {
		offset := buf.Next(1)[0]
		name, _ := ParseDomainName(src[offset:], src)
		nameBuf.WriteString(name)
		return nameBuf.String(), bufLength - buf.Len()
	} else {
		goto readBuf
	}
}

func ParseSOARdata(data []byte, src []byte) string {
	SoaName, usedCount := ParseDomainName(data, src)
	data = data[usedCount:]
	subName, usedCount := ParseDomainName(data, src)
	data = data[usedCount:]
	SoaName += " " + subName
	var buf bytes.Buffer
	buf.Grow(len(data))
	buf.Write(data)
	for buf.Len() != 0 {
		SoaName += fmt.Sprintf(" %d", BytesToUint32(buf.Next(4)))
	}
	return SoaName
}

func IsLabelPointer(length byte) bool {
	return length == 0xc0
}

func IsRoot(length byte) bool {
	return length == 0
}

func GenRequestId() uint16 {
	uid := make([]byte, 2)
	n, err := rand.Read(uid)
	if err != nil || n != len(uid) {
		return 0
	}
	return BytesToUint16(uid)
}
