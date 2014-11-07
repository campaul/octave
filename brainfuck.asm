LJMP trampoline NZP

code:
	BYTES "++++++++[>++++[>++>+++>+++>+<<<<-]>+>+>->>+[<]<-]>>.>---.+++++++..+++.>>.<-.<.+++.------.--------.>>+.>++."
	;BYTES ".>.>.>.>.>"
	BYTE 0x0

trampoline:
	LJMP start NZP

data:
	;BYTES "fart\n"
    BYTE 0x0
    BYTE 0x0
    BYTE 0x0
    BYTE 0x0
    BYTE 0x0
    BYTE 0x0
    BYTE 0x0
    BYTE 0x0
    BYTE 0x0
    BYTE 0x0
    BYTE 0x0
    BYTE 0x0
    BYTE 0x0
    BYTE 0x0
    BYTE 0x0
    BYTE 0x0
    BYTE 0x0
    BYTE 0x0
    BYTE 0x0
    BYTE 0x0
    BYTE 0x0
    BYTE 0x0

start:
	LAA R2, R1, code
	LAA R2, R3, data

; -1 is 0xFF
; R2 is program counter
; R3 is data pointer
; 0x3E = >  Increment data pointer
; 0x3C = <  Decrement data pointer
; 0x2B = +  Increment byte at R3
; 0x2D = -  Decrement byte at R3
; 0x2E = .  OUT [R3], 1
; 0x2C = ,  IN [R3], 1
; 0x5b = [  
; 0x5d = ]

bf:
	XOR R2, R2
	LOAD R2, R1
	MOV R2, R0

	LOADI 0xC2 ; -3E
	ADD R0, R2
	LJMP handleIncData Z

	LOADI 0xC4 ; -3C
	ADD R0, R2
	LJMP handleDecData Z

	LOADI 0xD2 ; -2E
	ADD R0, R2
	LJMP handlePrint Z

	LOADI 0xD5 ; -2B
	ADD R0, R2
	LJMP handleInc Z

	LOADI 0xD3 ; -2B
	ADD R0, R2
	LJMP handleDec Z

	LOADI 0xA5 ; -5B
	ADD R0, R2
	LJMP handleOpenBrace Z

	LOADI 0xA3 ; -5D
	ADD R0, R2
	LJMP handleCloseBrace Z

	XOR R0, R0 ; NUL byte -> exit
	ADD R0, R2
	LJMP exit Z

	LJMP incProg NZP

exit:
	HALT

incProg:
	LOADI 0x01
	ADD R1, R0

	LJMP bf NZP

handleIncData:
	LOADI 0x01
	ADD R3, R0

	LJMP incProg NZP

handleDecData:
	LOADI 0xFF
	ADD R3, R0

	LJMP incProg NZP

handleInc:
	XOR R2, R2
	LOAD R2, R3

	OUT R1, 0 ; Push program counter

	MOV R1, R0
	LOADI 0x01
	ADD R0, R1
	STORE R2, R3

	IN R1, 0 ; Pop program counter

	LJMP incProg NZP

handleDec:
	XOR R2, R2
	LOAD R2, R3

	OUT R1, 0 ; Push program counter

	MOV R1, R0
	LOADI 0xFF
	ADD R0, R1
	STORE R2, R3

	IN R1, 0 ; Pop program counter

	LJMP incProg NZP

handleOpenBrace:
	XOR R2, R2
	LOAD R2, R3
	AND R0, R0
	LJMP incProg NP

	OUT R3, 0
	; [R3] == 0
	; Find next matching ]
	LOADI 0x01
	MOV R3, R0

	openBraceLoop:
		LOADI 0x01
		ADD R1, R0

		XOR R2, R2
		LOAD R2, R1
		MOV R2, R0

		; R2 = [R1]

		LOADI 0xA5 ; -5B
		ADD R0, R2
		LJMP openBraceInc Z

		LOADI 0xA3 ; -5D
		ADD R0, R2
		LJMP openBraceDec Z

		LJMP openBraceLoop NZP

	openBraceInc:
		LOADI 0x01
		ADD R3, R0

		LJMP openBraceLoop NZP

	openBraceDec:
		LOADI 0xFF
		ADD R3, R0
		LJMP exitBraceLoop Z

		LJMP openBraceLoop NZP

	exitBraceLoop:
		IN R3, 0
		LJMP incProgB NZP

handleCloseBrace:
	XOR R2, R2
	LOAD R2, R3
	AND R0, R0
	LJMP incProgB Z

	LJMP exit NZP

handlePrint:
	XOR R2, R2
	LOAD R2, R3
	OUT R0, 1

	LJMP incProgB NZP

incProgB:
	LOADI 0x01
	ADD R1, R0

	LJMP bf NZP

