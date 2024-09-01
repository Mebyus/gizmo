package genc

const prelude = `
typedef unsigned char u8;
typedef unsigned short int u16;
typedef unsigned int u32;
typedef unsigned long int u64;
typedef __uint128_t u128;

typedef signed char s8;
typedef signed short int s16;
typedef signed int s32;
typedef signed long int s64;
typedef __int128_t s128;

typedef float f32;
typedef double f64;
typedef __float128 f128;

typedef _Bool bool;
typedef u32 rune;

typedef u64 uint;
typedef s64 sint;

typedef uint pint;

#define nil 0
#define true 1
#define false 0

typedef struct {
	const u8 *ptr;
	uint len;
} str;

static str
ku_static_string(const u8 *ptr, uint len) {
    str s;
    s.ptr = ptr;
    s.len = len;
    return s;
}

_Noreturn static void
ku_trap_unreachable() {
	__builtin_trap();
	__builtin_unreachable();
}

static void
ku_must(bool c) {
	if (c) {
		return;
	}

	ku_trap_unreachable();
}

_Noreturn static void
ku_panic_never(u64 pos) {
	ku_trap_unreachable();
}

#define KU_AMD64_LINUX_SYSCALL_WRITE 1

static sint
ku_syscall_write(u32 fd, const void *buf, uint size)
{
    sint ret;
    __asm__ volatile
    (
        "syscall"
		
		// RAX
        : "=a" (ret)

		// RAX
        : "0"(KU_AMD64_LINUX_SYSCALL_WRITE), 
        //  RDI      RSI       RDX
			"D"(fd), "S"(buf), "d"(size)

		// two registers are clobbered after system call
        : "rcx", "r11", 
			"memory"
    );
    return ret;
}

#define KU_AMD64_LINUX_SYSCALL_EXIT 60

_Noreturn static void
ku_syscall_exit(uint c) {
    sint ret;
    __asm__ volatile
    (
        "syscall"
		
		// outputs
		// RAX
        : "=a" (ret)
		
		// inputs
		// RAX
        : "0"(KU_AMD64_LINUX_SYSCALL_EXIT), 
        //  RDI     
			"D"(c)

		// clobbers
		// two registers are clobbered after system call
        : "rcx", "r11"
    );
	__builtin_unreachable();
}

#define KU_LINUX_STDOUT 1

static void
ku_print(str s) {
	if (s.len == 0) {
		return;
	}

	uint i = 0;
	while (i < s.len) {
		sint n = ku_syscall_write(KU_LINUX_STDOUT, s.ptr, s.len);
		if (n < 0) {
			return;
		}
		i += (uint)(n);
	}
}

`

func (g *Builder) prelude() {
	g.puts(prelude)
}
