package main

import (
	"fmt"
	mrand "math/rand"
    crand "crypto/rand"
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
                tar := make([]byte, mrand.Intn(64))
                crand.Read(tar)

                n, err := connection.Write(tar)
                if err != nil {
                    break
                }

                stuck := time.Now().Unix() - start
                channel <- &Message{7, fmt.Sprintf("Wrote %d bytes to %s (stuck for %d seconds)", n, remoteAddress, stuck)}

                if mrand.Intn(1) == 1 {
                    time.Sleep(time.Duration(delay + jitter) * time.Second)
                } else {
                    time.Sleep(time.Duration(delay - jitter) * time.Second)
                }
            }

            channel <- &Message{6, fmt.Sprintf("Closed connection from %s", remoteAddress)}

        }()
    }
}


