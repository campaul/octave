package main

import (
    "fmt"
    "os"
)

func main() {
    fmt.Println("Initializing Octave CPU...")

    file, err := os.Open(os.Args[1])

    if err != nil {
        return
    }

    defer file.Close()

    stat, err := file.Stat()

    if err != nil {
        return
    }

    instructions := make([]uint8, stat.Size())
    _, err = file.Read(instructions)

    if err != nil {
        return
    }

    for _, i := range instructions {
        execute(decode(i))
    }
}

type CPU struct {
    memory [65536]uint8
    r0 uint8
    r1 uint8
    r2 uint8
    r3 uint8
    pc uint16
}

type instruction func(uint8)

func decode(i uint8) (instruction, uint8) {
    switch i >> 5 {
        case 0:
            return mem, i
        case 1:
            return loadi, i
        case 2:
            if (i << 3) < 25 {
                return stack, i
            } else {
                return inte, i
            }
        case 3:
            return jmp, i
        case 4:
            return math, i
        case 5:
            return logic, i
        case 6:
            return in, i
        case 7:
            return out, i
    }

    return illegal, i
}

func execute(inst_func instruction, i uint8) {
    inst_func(i)
}

func mem(i uint8) {
    fmt.Println("mem")
}

func loadi(i uint8) {
    fmt.Println("loadi")
}

func stack(i uint8) {
    fmt.Println("stack")
}

func inte(i uint8) {
    fmt.Println("inte")
}

func jmp(i uint8) {
    fmt.Println("jmp")
}

func math(i uint8) {
    fmt.Println("math")
}

func logic(i uint8) {
    fmt.Println("logic")
}

func in(i uint8) {
    fmt.Println("in")
}

func out(i uint8) {
    fmt.Println("out")
}

func illegal(i uint8) {
    fmt.Println("illegal")
}
