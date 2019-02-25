package transport

import (
	"fmt"
	"log"
	"math/rand"
	"net"
	"strconv"
	"strings"

	"github.com/weriKK/sipserver/sipparser"
)

func SendSipMessage(sipMsg sipparser.SipMessage, destAddr string) {

	localAddr, err := net.ResolveUDPAddr("udp", ":5080")
	destAddrPort := strings.SplitN(destAddr, ":", 2)
	port, _ := strconv.Atoi(destAddrPort[1])
	remoteAddr := net.UDPAddr{
		IP:   net.ParseIP(destAddrPort[0]),
		Port: port,
	}
	conn, err := net.DialUDP("udp", localAddr, &remoteAddr)
	//conn, err := net.Dial("udp", destAddr)
	if err != nil {
		log.Printf("Error creating UDP connection to %v: %v", destAddr, err)
	}
	defer conn.Close()

	log.Printf("\n\n\n\n----------------------------------------------------------------\n")
	log.Printf("OUTGOING Sip Message. %s -> %s\n-----\n", conn.LocalAddr().String(), conn.RemoteAddr().String())
	log.Printf("%s\n", sipMsg.String())

	_, err = conn.Write([]byte(sipMsg.String()))
	if err != nil {
		log.Printf("Error sending UDP packet to %v: %v", destAddr, err)
	}

}

func SendResponse(sipMsg sipparser.SipMessage, statusCode int) {
	switch statusCode {
	case 200:
		send200(sipMsg)
	}
}

func send200(sipMsg sipparser.SipMessage) {

	sipMsg.SetStartLine("SIP/2.0 200 OK")

	if sipMsg.ToTag() == "" {
		tag := fmt.Sprintf(";tag=%s", randomString(10))
		sipMsg.SetHeader("To", sipMsg.Header("To")+tag)
	}

	SendSipMessage(sipMsg, sipMsg.ViaHostPort())
}

func randomString(len int) string {
	values := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := make([]byte, len)
	for i := 0; i < len; i++ {
		bytes[i] = byte(values[rand.Intn(62)])
	}
	return string(bytes)
}
