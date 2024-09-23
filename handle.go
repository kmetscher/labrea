package main

import (
    "fmt"
    mrand "math/rand"
	"net"
	"time"
)

func Handle(connection net.Conn, channel chan *Message, tar []byte, start int64, delay int, jitter int) (err error) {
    n, err := connection.Write(tar)
    if err != nil {
        return err
    }
    stuck := time.Now().Unix() - start
    channel <- &Message{7, fmt.Sprintf("Wrote %d bytes to %s (stuck for %d seconds)", n, connection.RemoteAddr().String(), stuck)}

    if mrand.Int() % 2 == 1 {
        time.Sleep(time.Duration(delay + mrand.Intn(jitter)) * time.Second)
    } else {
        time.Sleep(time.Duration(delay - mrand.Intn(jitter)) * time.Second)
    }
    return nil
}
