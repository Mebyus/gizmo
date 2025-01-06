#define PLATFORM_AMD64_LINUX_SYSCALL_WRITE 1

static sint
platform_amd64_linux_syscall_write(uint fd, const void* buf, uint size)
{
    sint ret;
    __asm__ volatile
    (
        "syscall"
		
		// RAX
        : "=a" (ret)

		// RAX
        : "0"(PLATFORM_AMD64_LINUX_SYSCALL_WRITE), 
        //  RDI      RSI       RDX
			"D"(fd), "S"(buf), "d"(size)

		// two registers are clobbered after system call
        : "rcx", "r11", 
			"memory"
    );
    return ret;
}

#define PLATFORM_AMD64_LINUX_SYSCALL_EXIT 60

_Noreturn static void
platform_amd64_linux_syscall_exit(uint c) {
    sint ret;
    __asm__ volatile
    (
        "syscall"
		
		// outputs
		// RAX
        : "=a" (ret)
		
		// inputs
		// RAX
        : "0"(PLATFORM_AMD64_LINUX_SYSCALL_EXIT), 
        //  RDI     
			"D"(c)

		// clobbers
		// two registers are clobbered after system call
        : "rcx", "r11"
    );
	__builtin_unreachable();
}
