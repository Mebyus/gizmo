# Amd64 ISA

```asm
endbr64  | f3 0f 1e fa

======================
push rax | 50
push rcx | 51
push rdx | 52
push rbx | 53
push rsp | 54
push rbp | 55
push rsi | 56
push rdi | 57
----------------------
push r8  | 41 50
push r9  | 41 51

push r15 | 41 57
======================

ret | c3

syscall | 0f 05

=======================================
mov rdi  => rax  | 48 89 f8
mov esi  => edx  | 89 f2

==============================================
load 32-bit constant into family-0 register
----------------------------------------------
RN  - register number
VAL - constant value (LE)
-----------------------------------------------
          |  opcode  | register |    value    |
-----------------------------------------------
bits      |    21    |     3    |      32     |
hex mask  | xx xx xb |    bbb   | xx xx xx xx |
encoding  | 48 c7 c0 |   ${RN}  |    ${VAL}   |
-----------------------------------------------
mov  0x3c => rax  |  48 c7 c0 3c 00 00 00
mov  0x0  => rax  |  48 c7 c0 00 00 00 00
mov  0x0  => rcx  |  48 c7 c1 00 00 00 00
mov  0x0  => rdi  |  48 c7 c7 00 00 00 00
----------------------------------------------

---------------------------------------
mov 0x0  => r8   | 49 c7 c0 00 00 00 00
mov 0x0  => r9   | 49 c7 c1 00 00 00 00
mov 0x0  => r15  | 49 c7 c7 00 00 00 00
---------------------------------------------

===============================================
load 64-bit constant into register
-----------------------------------------------
RN  - register number
VAL - constant value (LE)
---------------------------------------------------------
          | opcode | register |         value           |
---------------------------------------------------------
bits      |   13   |     3    |           64            |
hex mask  |  xx xb |    bbb   | xx xx xx xx xx xx xx xx |
encoding  |  48 b1 |   ${RN}  |         ${VAL}          |
---------------------------------------------------------
movabs 0x1021324354 => rax  |  48 b8 54 43 32 21 10 00 00 00
movabs 0x1021324354 => rcx  |  48 b9 54 43 32 21 10 00 00 00
movabs 0x1021324354 => rsi  |  48 be 54 43 32 21 10 00 00 00
movabs 0x1021324354 => rdi  |  48 bf 54 43 32 21 10 00 00 00
----------------------------------------------------------

=====================================================
load 32-bit constant into 64-bit wide memory location
denoted by address stored in register (family-0)
-----------------------------------------------------

movq 0x1234 => [rax]  |  48 c7 00 34 12 00 00
movq 0x1234 => [rcx]  |  48 c7 01 34 12 00 00
movq 0x1234 => [rdx]  |  48 c7 02 34 12 00 00
movq 0x1234 => [rax]  |  48 c7 07 34 12 00 00

movq 0xa    => [rdi]  |  48 c7 07 0a 00 00 00
-----------------------------------------------------

================================================================
load 64-bit value stored in register (family-0) into 64-bit wide
memory location denoted by address stored in another register (family-0)
----------------------------------------------------------------

mov rax => [rax]  |  48 89 00
mov rax => [rcx]  |  48 89 01
mov rax => [rdx]  |  48 89 02
mov rax => [rdi]  |  48 89 07

mov rcx => [rax]  |  48 89 08
mov rcx => [rcx]  |  48 89 09
mov rcx => [rdx]  |  48 89 0a
mov rcx => [rdi]  |  48 89 0f

mov rdx => [rax]  |  48 89 10
mov rdx => [rcx]  |  48 89 11
mov rdx => [rdx]  |  48 89 12
mov rdx => [rdi]  |  48 89 17

mov rdi => [rax]  |  48 89 38
mov rdi => [rcx]  |  48 89 39
mov rdi => [rdx]  |  48 89 3a
mov rdi => [rdi]  |  48 89 3f
----------------------------------------------------------------

=================================
family-0 register => number (0x)
---------------------------------
rax | 0
rcx | 1
rdx | 2
rbx | 3
rsp | 4
rbp | 5
rsi | 6
rdi | 7
---------------------------------

==================================
family-1 register => number (0x)
----------------------------------
r8  | 0
r9  | 1
r10 | 2
r11 | 3
r12 | 4
r13 | 5
r14 | 6
r15 | 7
----------------------------------

===================================================
add 32-bit constant value to register
and store the result into (the same) register
---------------------------------------------------

add  rax + 0x10203 => rax  |  48 05 03 02 01 00
add  rcx + 0x10203 => rcx  |  48 81 c1 03 02 01 00
add  rdx + 0x10203 => rdx  |  48 81 c2 03 02 01 00
add  rdi + 0x10203 => rdi  |  48 81 c7 03 02 01 00

add  r8  + 0x10203 => r8   |  49 81 c0 03 02 01 00
add  r15 + 0x10203 => r15  |  49 81 c7 03 02 01 00
---------------------------------------------------

===================================================
subtract 32-bit constant value from register
and store the result into (the same) register
---------------------------------------------------

sub  rax - 0x10203 => rax  |  48 2d 03 02 01 00
sub  rcx - 0x10203 => rcx  |  48 81 e9 03 02 01 00
sub  rdx - 0x10203 => rdx  |  48 81 ea 03 02 01 00
sub  rdi - 0x10203 => rdi  |  48 81 ef 03 02 01 00

sub  r8  - 0x10203 => r8   |  49 81 e8 03 02 01 00
sub  r15 - 0x10203 => r15  |  49 81 ef 03 02 01 00
---------------------------------------------------

====================================================
add 8-bit constant value to register
and store the result into (the same) register

----------------------------------------------------

add  rax + 0x3  => rax  |  48 83 c0 03
----------------------------------------------------

====================================================
subtract 8-bit constant value from register
and store the result into (the same) register

??? constant value encoding is signed
----------------------------------------------------

sub  rax - 0x12 => rax  |  48 83 e8 12
sub  rsp - 0x20 => rsp  |  48 83 ec 20

sub  r8  - 0x02 => r8   |  49 83 e8 02
sub  r15 - 0x01 => r15  |  49 83 ef 01
----------------------------------------------------


=========================================================
perform relative short jump within the same code segment

jump offset is limited to -128 to 127
rip (after jump) = rip (after reading jump instruction) + offset
offset is encoded as 1 byte number (signed, 2's complement)
---------------------------------------------------------

jmp rip + 0x1   |  eb 01
jmp rip - 0xc   |  eb f4
-----------------------------------------------------

=========================================================
perform relative near call within the same code segment

call offset is signed 32-bit number
rip (after call) = rip (after reading call instruction) + offset
offset is encoded as 4 byte number (signed, 2's complement)
-----------------------------------------------------------

call rip + 0x1   |  e8 01 00 00 00
call rip - 0xf   |  e8 f1 ff ff ff
-----------------------------------------------------------

==========================================================
one byte no operation
----------------------------------------------------------

nop  |  90
-----------------------------------------------------------

==========================================================
7 bytes no operation
----------------------------------------------------------

nopl ??? |  0f 1f 80 00 00 00 00
----------------------------------------------------------

==========================================================
perform logical compare of constant 32-bit value and register

-----------------------------------------------------------

test  0x12      &  rcx  |    48 f7 c1 12 00 00 00
test  0x120304  &  rdx  |    48 f7 c2 04 03 12 00
-----------------------------------------------------------

===========================================================
xor 8-bit constant value and register
and store the result into (the same) register
-----------------------------------------------------------

xor  rax ^ 0x11     => rax  |  48 83 f0 11
xor  rcx ^ 0x12     => rcx  |  48 83 f1 12
xor  rdx ^ 0x120304 => rdx  |  48 81 f2 04 03 12 00
-----------------------------------------------------------

===========================================================
xor two 64-bit registers
and store the result into (one of the two) registers
-----------------------------------------------------------
SRN  - source register number
DRN  - destination register number

SET dst = dst ^ src
----------------------------------------------
          |  opcode  |   src  |   dst  |
----------------------------------------------
bits      |    18    |    3   |    3   |
hex mask  | xx xx bb |   bbb  |   bbb  |
encoding  | 48 31 11 | ${SRN} | ${DRN} |
----------------------------------------------
xor  rax ^ rax  => rax  |  48 31 c0
xor  rax ^ rcx  => rax  |  48 31 c8
xor  rax ^ rdx  => rax  |  48 31 d0
xor  rax ^ rbx  => rax  |  48 31 d8

xor  rcx ^ rax  => rcx  |  48 31 c1
xor  rdx ^ rax  => rdx  |  48 31 c2
xor  rbx ^ rax  => rbx  |  48 31 c3
-----------------------------------------------------------

===========================================================
add two 64-bit registers
and store the result into (one of the two) registers
-----------------------------------------------------------
SRN  - source register number
DRN  - destination register number

SET dst = dst + src
----------------------------------------------
          |  opcode  |   src  |   dst  |
----------------------------------------------
bits      |    18    |    3   |    3   |
hex mask  | xx xx bb |   bbb  |   bbb  |
encoding  | 48 01 11 | ${SRN} | ${DRN} |
----------------------------------------------

add  rax + rax => rax  |  48 01 c0
add  rax + rcx => rax  |  48 01 c8
add  rax + rdx => rax  |  48 01 d0
add  rax + rbx => rax  |  48 01 d8

add  rcx + rax => rcx  |  48 01 c1
add  rdx + rax => rdx  |  48 01 c2
add  rbx + rax => rbx  |  48 01 c3

add  rax + rdi => rax  |  48 01 f8
add  rdi + rdi => rdi  |  48 01 ff
------------------------------------------------------------

============================================================
subtract two 64-bit registers
and store the result into (one of the two) registers
------------------------------------------------------------
SRN  - source register number
DRN  - destination register number

SET dst = dst - src
------------------------------------------------------------

sub  rax - rax => rax  |  48 29 c0
sub  rax - rcx => rax  |  48 29 c8
sub  rax - rdx => rax  |  48 29 d0
sub  rax - rbx => rax  |  48 29 d8

sub  rcx - rax => rcx  |  48 29 c1
sub  rdx - rax => rdx  |  48 29 c2
sub  rbx - rax => rbx  |  48 29 c3

sub  rax - rsi => rax  |  48 29 f0
sub  rax - rdi => rax  |  48 29 f8
sub  rdi - rdi => rdi  |  48 01 ff
------------------------------------------------------------

==============================================================
perform bitwise NOT on 64-bit register
and store the result into the same register
--------------------------------------------------------------

not rax => rax  |  48 f7 d0
not rcx => rcx  |  48 f7 d1
not rdx => rdx  |  48 f7 d2
not rdi => rdi  |  48 f7 d7

not r8  => r8   |  49 f7 d0
not r9  => r9   |  49 f7 d1
not r15 => r15  |  49 f7 d7
--------------------------------------------------------------
```


```
// from GAS
// source => destination

48 89 c0                mov    %rax,%rax
48 89 c1                mov    %rax,%rcx
48 89 c2                mov    %rax,%rdx
48 89 c8                mov    %rcx,%rax
48 89 c9                mov    %rcx,%rcx
48 89 ca                mov    %rcx,%rdx
4c 89 c0                mov    %r8,%rax
4c 89 c1                mov    %r8,%rcx
4c 89 c8                mov    %r9,%rax
4c 89 c9                mov    %r9,%rcx
4d 89 c0                mov    %r8,%r8
4d 89 c1                mov    %r8,%r9
4d 89 c8                mov    %r9,%r8
4d 89 c9                mov    %r9,%r9
49 89 c0                mov    %rax,%r8
49 89 c1                mov    %rax,%r9
49 89 c2                mov    %rax,%r10
49 89 f0                mov    %rsi,%r8
49 89 f1                mov    %rsi,%r9
49 89 f2                mov    %rsi,%r10
49 89 f8                mov    %rdi,%r8
49 89 f9                mov    %rdi,%r9
49 89 fa                mov    %rdi,%r10


48 85 c0                test   %rax,%rax
48 85 c1                test   %rax,%rcx
48 85 c6                test   %rax,%rsi
48 85 c7                test   %rax,%rdi
48 85 c0                test   %rax,%rax
48 85 c8                test   %rcx,%rax
48 85 f0                test   %rsi,%rax
48 85 f8                test   %rdi,%rax
```
