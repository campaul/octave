package main

import (
    "fmt"
    "os"
)

func main() {
    fmt.Println("Initializing Octave CPU...")

    file, err := os.Open("hello.bin")

    if err != nil {
        return
    }

    defer file.Close()

    stat, err := file.Stat()

    if err != nil {
        return
    }

    bs := make([]byte, stat.Size())
    _, err = file.Read(bs)

    if err != nil {
        return
    }

    for _, instruction := range bs {
        decode(instruction)
    }
}

func decode(instruction byte) {
    switch instruction >> 5 {
        case 1:
            fmt.Println("MEM")
        case 2:
            fmt.Println("LOADI")
        case 3:
            fmt.Println("STACK")
        case 4:
            fmt.Println("JMP")
        case 5:
            fmt.Println("MATH")
        case 6:
            fmt.Println("LOGIC")
        case 7:
            fmt.Println("IN")
        case 8:
            fmt.Println("OUT")
    }
}
