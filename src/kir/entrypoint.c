// program entrypoint
extern _Noreturn void
_start(void);

/*
At entrypoint function stack pointer is given to us by OS loader.
Transfer the pointer to function "xk_start" as a first argument in "rdi" register.
*/
__asm__ volatile (
	".global _start\n"
	"_start:\n\t"
    "xor %rbp, %rbp\n\t"
	"mov %rsp, %rdi\n\t"
	"call start"
);
