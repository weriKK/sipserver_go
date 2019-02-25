package sipparser

import (
	"fmt"
	"log"
	"strconv"
	"strings"
)

func parseStartLine(startLine string, sipMsg *SipMessage) error {
	sipMsg.startLine = startLine

	startLineParts := strings.SplitN(startLine, " ", 3)
	if len(startLineParts) != 3 {
		return fmt.Errorf("[SIPPARSER] Failed to parse request start line: %q", startLine)
	}

	if strings.HasPrefix(startLineParts[0], "SIP") {
		// Response
		sipMsg.isResponse = true
		sipMsg.version = startLineParts[0]
		sipMsg.statusCode, _ = strconv.Atoi(startLineParts[1])
		sipMsg.statusMessage = startLineParts[2]
	} else {
		// Request
		sipMsg.isResponse = false
		sipMsg.method = startLineParts[0]
		sipMsg.requestURI = startLineParts[1]
		sipMsg.version = startLineParts[2]
	}

	return nil
}

// ParseMessage parses a sip request string into SipMessage structure
func ParseMessage(msg string) (SipMessage, error) {

	sipMsg := SipMessage{}
	sipMsg.headers = make(map[string][]string)
	sipMsg.body = []string{}
	sipMsg.separator = "\n"

	if msg[len(msg)-2] == '\r' {
		sipMsg.separator = "\r\n"
	}

	msgLines := strings.Split(msg, sipMsg.separator)

	err := parseStartLine(msgLines[0], &sipMsg)
	if err != nil {
		return SipMessage{}, err
	}

	for i := 1; i < len(msgLines); i++ {

		// Body
		if strings.TrimSpace(msgLines[i]) == "" {
			sipMsg.body = msgLines[i:]
			break
		}

		// Headers
		hdrField := strings.SplitN(msgLines[i], ":", 2)
		if len(hdrField) != 2 {
			log.Printf("[SIPPARSER] Failed to split header field to field-name and field-value: %v", msgLines[i])
			continue
		}

		fieldName := strings.TrimSpace(hdrField[0])
		fieldValue := strings.TrimSpace(hdrField[1])
		if _, ok := sipMsg.headers[fieldName]; ok {
			log.Printf("[SIPPARSER] multiple occurances of the same header not implemented: %v", fieldName)
			continue
		}
		sipMsg.headers[fieldName] = append(sipMsg.headers[fieldName], fieldValue)
	}
	return sipMsg, nil
}
