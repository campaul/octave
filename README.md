# octave

An 8-bit CPU designed for education

## ISA

### Jump

Format     | Operation
---------- | -----------------
000 RS NZP | PC <- MEM[RS] NZP

### Load Immediate

Format     | Operation
---------  | -------------
0010 IIII  | R0[HIGH] <- I
0011 IIII  | R0[LOW] <- I

### ALU Operations

Format     | Operation
---------- | -------------
0100 RD RS | RD <- RD + RS
0101 RD RS | RD <- RS
0110 RD RS | RD <- RD & RS
0111 RD RS | RD <- RD ^ RS

### Memory Operations

Format     | Operation
---------- | -----------------
1000 RX RY | R0 <- MEM[RX, RY]
1001 RX RY | MEM[RX, RY] <- R0

### Stack Operations

Format     | Operation
---------- | -----------
10100000   | add8
10100001   | sub8
10100010   | mul8
10100011   | div8
10100100   | mod8
10100101   | neg8
10100110   | and8
10100111   | or8
10101000   | xor8
10101001   | not8
10101010   | add16
10101011   | sub16
10101100   | mul16
10101101   | div16
10101110   | mod16
10101111   | neg16
10110000   | and16
10110001   | or16
10110010   | xor16
10110011   | not16
10110100   | call
10110101   | trap
10110110   | ret
10110111   | iret
10111000   | int0 enable
10111001   | int1 enable
10111010   | int2 enable
10111011   | int3 enable
10111100   | int4 enable
10111101   | int5 enable
10111110   | int6 enable
10111111   | int7 enable

### Device IO

Format     | Operation
---------- | -------------
110 RD DV  | RD <- DEV[DV]
111 DV RS  | DEV[DV] <- RS

## Assembly Hello World

```
; Octave CPU - Hello World

; R2R1 ← string
LAA R2, R1, string

loop:
; R0 ← [R2R1]
LOAD R2, R1

; IF R0 == 0, JMP exit
XOR R3, R3
ADD R3, R0
LRA exit
JMP R0 Z

; PRINT R3
OUT R3, 1

; INC R1
LOADI 0x01
ADD R1, R0

; JMP loop
LRA loop
JMP R0 NZP

exit:
HALT

string:
BYTES "Hello world!\n"
BYTE 0x00
```
