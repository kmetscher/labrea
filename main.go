package main

import (
    "fmt"
    "flag"
    "net"
)

func main() {
    protoFlag := flag.String("protocol", "ssh", "protocol to tarpit (defaults to ssh)")
    addressFlag := flag.String("address", "0.0.0.0", "bind address (defaults to 0.0.0.0)")
    portFlag := flag.Int("port", 22, "bind port (defaults to 22)")
    delayFlag := flag.Int("delay", 10, "seconds to delay writes to connection (defaults to 10)")
    jitterFlag := flag.Int("jitter", 5, "plus-minus range to vary delays within (defaults to 5)")
    flag.Parse()

    address := net.ParseIP(*addressFlag)
    channel := make(chan *Message)

    if *protoFlag == "http" {
        go Http(address, *portFlag, *delayFlag, *jitterFlag, channel)
    } else {
        go Ssh(address, *portFlag, *delayFlag, *jitterFlag, channel)
    }

    for {
        message := <-channel
        fmt.Println(message.content)
    }

}
