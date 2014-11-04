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

    instructions := make([]byte, stat.Size())
    _, err = file.Read(bs)

    if err != nil {
        return
    }

    for _, i := range instructions {
        execute(decode(i))
    }
}

type instruction func(byte)

func decode(i byte) (instruction, byte) {
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

func execute(inst_func instruction, i byte) {
    inst_func(i)
}

func mem(i byte) {
    fmt.Println("mem")
}

func loadi(i byte) {
    fmt.Println("loadi")
}

func stack(i byte) {
    fmt.Println("stack")
}

func inte(i byte) {
    fmt.Println("inte")
}

func jmp(i byte) {
    fmt.Println("jmp")
}

func math(i byte) {
    fmt.Println("math")
}

func logic(i byte) {
    fmt.Println("logic")
}

func in(i byte) {
    fmt.Println("in")
}

func out(i byte) {
    fmt.Println("out")
}

func illegal(i byte) {
    fmt.Println("illegal")
}
