SYS_EXIT   = 0x3c

.section .text

.global coven_linux_syscall_exit

coven_linux_syscall_exit:
	mov $SYS_EXIT, %rax 
	syscall
