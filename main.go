package main

import (
	"log"

	"github.com/weriKK/sipserver/registrar"
	"github.com/weriKK/sipserver/sipparser"
	"github.com/weriKK/sipserver/transport"
)

func main() {
	incoming := make(chan transport.Message)
	go transport.ListenForRequests(incoming)

	for {
		message := <-incoming

		log.Printf("\n\n\n\n----------------------------------------------------------------\n")
		log.Printf("INCOMING Sip Message. %s -> %s\n-----\n", message.RemoteAddr(), message.LocalAddr())

		sipMsg, err := sipparser.ParseMessage(message.String())
		if err != nil {
			// error response
			// log.Printf("Failed to parse sip message: %q\n", message)
			// log.Println(err.Error())
			continue
		}
		log.Printf("%s\n-----\n", sipMsg)

		switch sipMsg.Method() {
		case "REGISTER":
			duration, err := sipMsg.Expires()
			if err != nil {
				log.Println("Invalid REGISTER message: missing Expires header")
				continue
			}

			if 0 < duration {
				registrar.RegisterSubscriber(sipMsg.ToURI(), sipMsg.ContactHostPort())
			} else {
				registrar.DeregisterSubscriber(sipMsg.ToURI())
			}
			transport.SendResponse(sipMsg, 200)
		case "INVITE":
			destAddr, err := registrar.Contact(sipMsg.ToURI())
			if err != nil {
				log.Printf("Called Party not registered: %v", sipMsg.ToURI())
			}

			sipMsg.AddHeader("Record-Route", "<sip:192.168.1.100:5080;lr>")
			// sipMsg.SetHeader("Via", sipMsg.Header("Via")+"_kova")

			transport.SendSipMessage(sipMsg, destAddr)
		default:
			var destAddr string
			var err error

			if sipMsg.IsResponse() {
				destAddr = sipMsg.ViaHostPort()
			} else {
				destAddr, err = registrar.Contact(sipMsg.ToURI())
			}
			if err != nil {
				log.Printf("Called Party not registered: %v", sipMsg.ToURI())
			}

			// sipMsg.AddHeader("Via", sipMsg.Header("Via"))
			// sipMsg.SetHeader("Via", sipMsg.Header("Via")+"_kova")

			transport.SendSipMessage(sipMsg, destAddr)
		}
	}

	//pc.WriteTo([]byte("Hello from server"), addr)

}
