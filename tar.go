package main

import (
    "bytes"
    crand "crypto/rand"
)

func MakeTar(length int, newline bool) (err error, tar []byte) {
    tar = make([]byte, length)
    _, e := crand.Read(tar)
    if e != nil {
        return e, []byte("")
    }
    if newline {
        tar[len(tar)-1] = 0x0A
    }
    return nil, bytes.ToValidUTF8(tar, []byte(" ")) 
}
