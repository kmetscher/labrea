package main

import (
	"bytes"
	crand "crypto/rand"
	"fmt"
	mrand "math/rand"
	"net"
	"time"
)

func Ssh(address net.IP, port int, delay int, jitter int, channel chan *Message) (error error) {
    tcpAddress := net.TCPAddr{IP: address, Port: port}
    listener, err := net.ListenTCP("tcp", &tcpAddress)
    
    if err != nil {
        return err
    }

    channel <- &Message{5, fmt.Sprintf("Opened new listener on %s", tcpAddress.String())}

    for {
        connection, err := listener.Accept()
        if err != nil {
            return err
        }

        go func() {
            remoteAddress := connection.RemoteAddr().String()
            channel <- &Message{6, fmt.Sprintf("New connection from %s", remoteAddress)}
            start := time.Now().Unix()

            for {
                tar := make([]byte, mrand.Intn(63))
                _, err := crand.Read(tar)
                if err != nil {
                    break
                }
                tar = append(tar, 0x0A)

                n, err := connection.Write(bytes.ToValidUTF8(tar, []byte("")))
                if err != nil {
                    break
                }

                stuck := time.Now().Unix() - start
                channel <- &Message{7, fmt.Sprintf("Wrote %d bytes to %s (stuck for %d seconds)", n, remoteAddress, stuck)}

                if mrand.Intn(1) == 1 {
                    time.Sleep(time.Duration(delay + mrand.Intn(jitter)) * time.Second)
                } else {
                    time.Sleep(time.Duration(delay - mrand.Intn(jitter)) * time.Second)
                }
            }

            channel <- &Message{6, fmt.Sprintf("Closed connection from %s", remoteAddress)}

        }()
    }
}


