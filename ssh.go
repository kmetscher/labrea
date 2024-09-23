package main

import (
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

    channel <- &Message{5, fmt.Sprintf("Opened new SSH listener on %s", tcpAddress.String())}

    for {
        connection, err := listener.Accept()
        if err != nil {
            return err
        }

        go func() {
            remoteAddress := connection.RemoteAddr().String()
            channel <- &Message{6, fmt.Sprintf("SSH: New connection from %s", remoteAddress)}
            start := time.Now().Unix()
            for {
                // mrand.Intn can return 0, creating a zero-length byte slice
                err, tar := MakeTar(mrand.Intn(64) + 1, true)
                if err != nil {
                    break
                }
                err = Handle(connection, channel, tar, start, delay, jitter)
                if err != nil {
                    break
                }
            }
            channel <- &Message{6, fmt.Sprintf("Closed connection from %s", remoteAddress)}
        }()
    }
}


