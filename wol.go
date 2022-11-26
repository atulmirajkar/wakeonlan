package main

import (
	"bytes"
	"encoding/hex"
	"flag"
	"fmt"
	"net"
	"regexp"
)

// error handling start
const (
	MacAddrFormatWrong int = iota
	MacAddrLengthWrong
	HostNameFormatWrong
	PortFormatWrong
)

type WolError struct {
	kind int
}

func (e WolError) Error() string {
	switch e.kind {
	case MacAddrFormatWrong:
		return "Mac address contains chars other than hex"
	case MacAddrLengthWrong:
		return "Mac address length wrong"
	case HostNameFormatWrong:
		return "Host name is in the wrong format"
	case PortFormatWrong:
		return "Port in the wrong format"
	}
	return "Unknown error"
}

var (
	MacAddrFormatWrongErr  = WolError{kind: MacAddrFormatWrong}
	MacAddrLengthWrongErr  = WolError{kind: MacAddrLengthWrong}
	HostNameFormatWrongErr = WolError{kind: HostNameFormatWrong}
	PortFormatWrongErr     = WolError{kind: PortFormatWrong}
)

//error handling end

/*
https://en.wikipedia.org/wiki/Wake-on-LAN#:~:text=The%20magic%20packet%20is%20a,a%20total%20of%20102%20bytes.
payload contains 6 bytes of all 255
ff ff ff ff ff ff
followed by 16 repetitions of the target computer's mac address
Mac Address format: 48 bit, 6 bytes, 12 hex
*/
func main() {
	host := flag.String("h", "", "Enter hostname")
	port := flag.String("p", "", "Enter port")
	macAddr := flag.String("m", "", "Enter mac addr in hex without the semicolons")

	flag.Parse()
	err := validateDestination(*host, *port, *macAddr)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = sendMagicPacket(*host, *port, *macAddr)
	if err != nil {
		fmt.Println(err)
	}
}

func sendMagicPacket(host string, port string, macAddr string) error {
	hostPort := net.JoinHostPort(host, port)
	payLoad, err := createPayload(macAddr)
	if err != nil {
		return err
	}

	conn, err := net.Dial("udp", hostPort)
	if err != nil {
		return err
	}

	_, err = conn.Write(payLoad)
	if err != nil {
		return err
	}
	conn.Close()
	return nil
}

func validateDestination(host string, port string, macAddr string) error {
	if !isHostNameValid(host) {
		return HostNameFormatWrongErr
	}
	if !isPortValid(port) {
		return PortFormatWrongErr
	}

	if len(macAddr) < 12 {
		return MacAddrLengthWrongErr
	}
	if !isMacAddrValid(macAddr) {
		return MacAddrFormatWrongErr
	}
	return nil
}

func createPayload(macAddr string) ([]byte, error) {
	var b bytes.Buffer
	firstPart, err := hex.DecodeString("FFFFFFFFFFFF")
	if err != nil {
		return nil, err
	}
	b.Write(firstPart)
	for i := 0; i < 16; i++ {
		mac, _ := hex.DecodeString(macAddr)
		b.Write(mac)
	}
	return b.Bytes(), nil
}

func isMacAddrValid(macAddr string) bool {
	matched, _ := regexp.MatchString("^([0-9A-Fa-f]{12})$", macAddr)
	return matched
}

func isHostNameValid(hostName string) bool {
	const pattern = "^(([a-zA-Z0-9]|[a-zA-Z0-9][a-zA-Z0-9\\-]*[a-zA-Z0-9])\\.)*([A-Za-z0-9]|[A-Za-z0-9][A-Za-z0-9\\-]*[A-Za-z0-9])$"
	matched, _ := regexp.MatchString(pattern, hostName)
	if !matched {
		const ipPattern = "^(([a-zA-Z0-9]|[a-zA-Z0-9][a-zA-Z0-9\\-]*[a-zA-Z0-9])\\.)*([A-Za-z0-9]|[A-Za-z0-9][A-Za-z0-9\\-]*[A-Za-z0-9])$"
		matched, _ = regexp.MatchString(ipPattern, hostName)
	}
	return matched
}

func isPortValid(port string) bool {
	const pattern = "^([0-9]{4})$"
	matched, _ := regexp.MatchString(pattern, port)
	return matched
}
