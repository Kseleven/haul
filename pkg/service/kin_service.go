package service

import (
	"encoding/hex"
	"fmt"
	"net"

	"github.com/kseleven/haul/pkg/resource"
)

func Request(r resource.Request) error {
	targetAddr := fmt.Sprintf("%s:%d", r.Host, r.Port)
	addr, err := net.ResolveUDPAddr("udp", targetAddr)
	if err != nil {
		return fmt.Errorf("resolve address %s failed:%s", targetAddr, err.Error())
	}

	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		return fmt.Errorf("dial host %s failed:%s", targetAddr, err.Error())
	}
	defer conn.Close()

	if err := writeMsg(conn, r); err != nil {
		return err
	}

	data := make([]byte, 512)
	length, rAddr, err := conn.ReadFromUDP(data)
	if err != nil {
		return err
	}
	fmt.Println(rAddr, length, hex.Dump(data[:length]))
	ReadMsg(data[:length])
	return nil
}

func ReadMsg(data []byte) {
	message := &resource.Message{}
	message.Decode(data)
	fmt.Printf("head section:%+v\n", message.Header)
	fmt.Printf("question section:%+v\n", message.Question)
	for _, answer := range message.Answers {
		fmt.Printf("answers section:%+v\n", answer)
	}
	for _, authority := range message.Authority {
		fmt.Printf("authority section:%+v\n", authority)
	}
}

func writeMsg(conn *net.UDPConn, r resource.Request) error {
	message := resource.NewMessage(r.QName, r.QType)
	n, err := conn.Write(message.Encode())
	if err != nil {
		return fmt.Errorf("write message failed:%s", err.Error())
	}
	fmt.Println("write length", n)

	return nil
}
