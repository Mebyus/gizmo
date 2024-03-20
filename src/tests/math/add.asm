.section .text

.global coven_add_asm

coven_add_asm:
    add %rax, %rdi
    mov %rdi, %rcx
    ret
