" Vim syntax file
" Language:   octavecpu
" Maintainer: The Octave CPU Team
" Version:    $Revision$

if version < 600
  syntax clear
elseif exists("b:current_syntax")
  finish
endif


syn case match

syn keyword octaveInstruction HALT WAITI INTE INTD JMP
syn keyword octaveInstruction LOADI LOADIL LOADIH
syn keyword octaveInstruction ADD DIV
syn keyword octaveInstruction AND XOR
syn keyword octaveInstruction LOAD STORE
syn keyword octaveInstruction LRA LAA
syn keyword octaveInstruction STACKOP
syn keyword octaveInstruction IN OUT
syn keyword octaveInstruction PUSH MOV
syn keyword octaveInstruction NZP

syn keyword octaveRegister R0 R1 R2 R3

syn match   octaveComment /;.*$/
syn match   octaveLabel /[-a-zA-Z$._][-a-zA-Z$._0-9]*:/

syn match   octaveConstant /\d\+/
syn match   octaveConstant /0x\x\x/
syn match   octaveConstant /0x\x/

if version >= 508 || !exists("did_c_syn_inits")
  if version < 508
    let did_c_syn_inits = 1
    command -nargs=+ HiLink hi link <args>
  else
    command -nargs=+ HiLink hi def link <args>
  endif

  HiLink octaveInstruction Keyword
  HiLink octaveRegister Identifier
  HiLink octaveComment Comment
  HiLink octaveLabel Label
  HiLink octaveConstant Number

  delcommand HiLink
endif

let b:current_syntax = "octavecpu"
