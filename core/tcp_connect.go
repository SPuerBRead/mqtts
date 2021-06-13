package core

import (
	"crypto/tls"
	"encoding/hex"
	"fmt"
	"io"
	"net"
	"strconv"
	"time"
)

var connectTimeout int = 5

func ConnectWithSingleProbePack(ip string, port int, probePack string) (error, []byte) {
	var bufferLens int = 10240
	var tmpBufferLens int = 256
	buf := make([]byte, 0, bufferLens)
	tmp := make([]byte, tmpBufferLens)
	var target = ip + ":" + strconv.Itoa(port)
	conn, connectErr := net.DialTimeout("tcp", target, time.Duration(connectTimeout)*time.Second)
	if connectErr != nil {
		return connectErr, buf
	}
	defer conn.Close()
	setReadDeadlineErr := conn.SetReadDeadline(time.Now().Add(time.Duration(connectTimeout) * time.Second))
	if setReadDeadlineErr != nil {
		return setReadDeadlineErr, buf
	}
	setWriteDeadlineErr := conn.SetWriteDeadline(time.Now().Add(time.Duration(connectTimeout) * time.Second))
	if setWriteDeadlineErr != nil {
		return setWriteDeadlineErr, buf
	}
	hexProbePack := fmt.Sprintf("%X", probePack)
	decoded, decodeProbePackErr := hex.DecodeString(hexProbePack)
	if decodeProbePackErr != nil {
		return decodeProbePackErr, buf
	}
	_, writeErr := conn.Write(decoded)
	if writeErr != nil {
		return writeErr, buf
	}
	for {
		n, err := conn.Read(tmp)
		if err != nil && err != io.EOF {
			break
		}
		if n == 0 {
			break
		}
		if bufferLens-len(buf) > len(tmp) {
			buf = append(buf, tmp[:n]...)
		} else {
			break
		}
	}
	return nil, buf
}

func ConnectWithSingleProbePackTCPTLS(ip string, port int, probePack string) (error, []byte) {
	var bufferLens int = 10240
	var tmpBufferLens int = 256
	buf := make([]byte, 0, bufferLens)
	tmp := make([]byte, tmpBufferLens)
	var target = ip + ":" + strconv.Itoa(port)
	conn, connectErr := net.DialTimeout("tcp", target, time.Duration(connectTimeout)*time.Second)
	if connectErr != nil {
		return connectErr, buf
	}
	defer conn.Close()
	setReadDeadlineErr := conn.SetReadDeadline(time.Now().Add(time.Duration(connectTimeout) * time.Second))
	if setReadDeadlineErr != nil {
		return setReadDeadlineErr, buf
	}
	setWriteDeadlineErr := conn.SetWriteDeadline(time.Now().Add(time.Duration(connectTimeout) * time.Second))
	if setWriteDeadlineErr != nil {
		return setWriteDeadlineErr, buf
	}
	conf := &tls.Config{
		InsecureSkipVerify: true,
	}
	tlsConn := tls.Client(conn, conf)
	handshakeErr := tlsConn.Handshake()
	if handshakeErr != nil {
		return handshakeErr, buf
	}
	hexProbePack := fmt.Sprintf("%X", probePack)
	decoded, decodeProbePackErr := hex.DecodeString(hexProbePack)
	if decodeProbePackErr != nil {
		return decodeProbePackErr, buf
	}
	_, writeErr := tlsConn.Write(decoded)
	if writeErr != nil {
		return writeErr, buf
	}
	for {
		n, err := tlsConn.Read(tmp)
		if err != nil && err != io.EOF {
			break
		}
		if n == 0 {
			break
		}
		if bufferLens-len(buf) > len(tmp) {
			buf = append(buf, tmp[:n]...)
		} else {
			break
		}
	}
	return nil, buf
}
