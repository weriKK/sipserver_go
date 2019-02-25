package transport

import (
	"log"
	"net"
)

type Message struct {
	remoteAddr net.Addr
	localAddr  net.Addr
	buffer     []byte
	len        int
}

func (m Message) RemoteAddr() string {
	return m.remoteAddr.String()
}

func (m Message) LocalAddr() string {
	return m.localAddr.String()
}

func (m Message) String() string {
	return string(m.buffer)
}

func ListenForRequests(output chan Message) {

	pc, err := net.ListenPacket("udp", ":5080")
	if err != nil {
		log.Fatal(err)
	}
	defer pc.Close()

	buffer := make([]byte, 65507)

	for {
		n, addr, _ := pc.ReadFrom(buffer)
		if err != nil {
			log.Fatal(err)
		}

		msg := Message{
			remoteAddr: addr,
			localAddr:  pc.LocalAddr(),
			buffer:     buffer[:n],
			len:        n,
		}

		output <- msg
	}
}
