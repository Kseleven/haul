package service

import (
	"fmt"
	"net"
	"time"

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

	beginTime := time.Now()
	if err := writeMsg(conn, r); err != nil {
		return err
	}
	data := make([]byte, 512)
	length, rAddr, err := conn.ReadFromUDP(data)
	if err != nil {
		return err
	}

	return ReadMsg(data[:length], rAddr, time.Now().Sub(beginTime))
}

func ReadMsg(data []byte, rAddr *net.UDPAddr, spendTime time.Duration) error {
	message := &resource.Message{}
	message.Decode(data)

	fmt.Printf("Server: %s\n", rAddr.String())
	fmt.Printf("Query Time: %s\n", spendTime)
	fmt.Printf("When: %s\n", time.Now().Format(time.RFC3339))
	fmt.Printf("Message Size Recvd: %d \n\n", len(data))
	fmt.Println(message.Header)
	fmt.Println(message.Question)
	if len(message.Answers) > 0 {
		fmt.Println("Answer Section:")
		for _, answer := range message.Answers {
			fmt.Println(answer.String())
		}
	}

	if len(message.Authority) > 0 {
		fmt.Println("Authority Section:")
		for _, authority := range message.Authority {
			fmt.Println(authority.String())
		}
	}

	return nil
}

func writeMsg(conn *net.UDPConn, r resource.Request) error {
	message := resource.NewMessage(r.QName, r.QType)
	if _, err := conn.Write(message.Encode()); err != nil {
		return fmt.Errorf("send message failed:%s", err.Error())
	}

	return nil
}
