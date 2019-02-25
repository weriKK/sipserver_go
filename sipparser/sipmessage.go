package sipparser

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// SipMessage contains a SIP request or response
type SipMessage struct {
	startLine string
	version   string

	// Request
	method     string
	requestURI string

	// Response
	isResponse    bool
	statusCode    int
	statusMessage string

	headers   map[string][]string
	body      []string
	separator string
}

func (sm SipMessage) StartLine() string {
	return sm.startLine
}

func (sm *SipMessage) SetStartLine(line string) {
	sm.startLine = line
}

func (sm *SipMessage) SetHeader(name string, val string) {
	sm.headers[name][0] = val
}

func (sm *SipMessage) AddHeader(name string, val string) {
	sm.headers[name] = append(sm.headers[name], val)
}

func (sm SipMessage) Header(name string) string {
	if val, ok := sm.headers[name]; ok {
		return val[0]
	}

	return ""
}

func (sm SipMessage) IsResponse() bool {
	return sm.isResponse
}

func (sm SipMessage) Method() string {
	return sm.method
}

func (sm SipMessage) FromTag() string {
	return parseTagFromHdr(sm.Header("From"))
}

func (sm SipMessage) ToTag() string {
	return parseTagFromHdr(sm.Header("To"))
}

func (sm SipMessage) ContactURI() string {
	return parseURIFromHdr(sm.Header("Contact"))
}

func (sm SipMessage) ToURI() string {
	return parseURIFromHdr(sm.Header("To"))
}

func (sm SipMessage) ContactHostPort() string {
	return parseHostPortFromURI(sm.ContactURI())
}

func (sm SipMessage) ViaHostPort() string {
	via := sm.Header("Via")
	viaParts := strings.SplitN(via, " ", 2)
	host := viaParts[1]
	hostEndIdx := strings.Index(host, ";")
	host = host[:hostEndIdx]
	return host
}

func (sm SipMessage) Expires() (int, error) {
	expires := sm.Header("Expires")
	if expires != "" {
		return strconv.Atoi(expires)
	}

	return 0, errors.New("No Expires hdr in Sip Message")
}

func (sm SipMessage) String() string {
	msg := sm.startLine + sm.separator

	for fieldName, fieldValues := range sm.headers {
		msg += fmt.Sprintf("%s:", fieldName)
		msg += strings.Join(fieldValues, ",")
		msg += sm.separator
	}

	if 0 < len(sm.body) {
		body := strings.Join(sm.body, sm.separator)
		msg += sm.separator + body
	}

	return msg
}

func parseHostPortFromURI(URI string) string {
	hostPortStartIdx := strings.Index(URI, "@")
	hostPortEndIdx := strings.Index(URI, ";")
	if hostPortEndIdx == -1 {
		return URI[hostPortStartIdx+1:]
	}

	return URI[hostPortStartIdx+1 : hostPortEndIdx]
}

func parseURIFromHdr(hdrLine string) string {
	uriStartIdx := strings.LastIndex(hdrLine, "<")
	uriEndIdx := strings.LastIndex(hdrLine, ">")

	return hdrLine[uriStartIdx+1 : uriEndIdx]
}

func parseTagFromHdr(hdrLine string) string {
	uriEndIdx := strings.LastIndex(hdrLine, ">")

	if uriEndIdx+1+5 < len(hdrLine) && hdrLine[uriEndIdx+1:uriEndIdx+6] == ";tag=" {
		return hdrLine[uriEndIdx+6:]
	}

	return ""
}
