package genc

const entrypoint = `
_Noreturn void
ku_start(void *stack_start) {
	ku_str args_buf[MAX_PROC_ARGS];
	ku_proc_args args;
	args.ptr = args_buf;
	args.len = MAX_PROC_ARGS;

	ku_read_proc_args(&args, stack_start);
	ku_main_main();
	ku_syscall_exit(0);
}

// program entrypoint
extern _Noreturn void
_start();

/*
At entrypoint function stack pointer is given to us by OS loader.
Transfer the pointer to function "start" as a first argument in "rdi" register.
*/
__asm__ (
	".global _start\n"
	"_start:\n\t"
	"mov %rsp, %rdi\n\t" 
	"call ku_start"
);

`

func (g *Builder) entrypoint() {
	g.puts(entrypoint)
}
