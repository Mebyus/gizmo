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

// bytes chunk
typedef struct {
	u8   *ptr;
	uint  len;
} ku_bc;

typedef struct {
	const u8 *ptr;
	uint len;
} ku_str;

static ku_str
ku_static_string(const u8 *ptr, uint len) {
    ku_str s;
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

static bool
ku_bc_valid(ku_bc c) {
	if (c.ptr == nil && c.len == 0) {
		return true;
	}
	if (c.ptr != nil && c.len != 0) {
		return true;
	}
	return false;
}

static ku_str
ku_str_from_chunk(ku_bc c) {
	ku_must(ku_bc_valid(c));

	ku_str s;
	s.ptr = c.ptr;
	s.len = c.len;
	return s;
}

static ku_str
ku_str_from_cstr_ptr(const u8 *ptr) {
	ku_must(ptr != nil);

	uint i = 0;
	while (ptr[i] != 0) {
		i += 1;
	}

	ku_str s;
	s.ptr = ptr;
	s.len = i;
	return s;
}

#define MAX_PROC_ARGS 64

typedef struct {
	ku_str  *ptr;
	uint  len;
} ku_proc_args;

/*

Layout of memory at stack_start

***            --                                             -- low addresses
argc           -- number of process args                      -- 8 bytes
array_arg_ptr  -- array with argc number of pointers          -- 8 bytes each (times argc)
0              -- array_arg_ptr is terminated by nil pointer  -- 8 bytes
array_env_ptr  -- array with environment entry pointer        -- 8 bytes each (until nil terminator)
0              -- array_env_ptr is terminated by nil pointer  -- 8 bytes
***            --                                             -- high addresses

Each arg_ptr leads to zero-terminated byte array (C string) with argument's string value.
Each env_ptr leads to zero-terminated byte array (C string) with environment entry.

Environment entry has the following form: "<name>=<value>".

*/
static void
ku_read_proc_args(ku_proc_args *args, void *stack_start) {
	uint argc = *(uint*)(stack_start);
	ku_must(argc <= args->len); // we provide only fixed amount of space in buffer to store chunk of strings

	u8 **arg_ptr = (u8**)(stack_start);

	// start index at 1 because we want to skip 8 bytes of argc
	uint i = 1;
	for (; i <= argc; i += 1) {
		u8 *cstr_ptr = arg_ptr[i]; // obtain pointer to C string with arg value
		if (cstr_ptr == nil) {
			// terminate early if we encounter nil pointer
			//
			// this is just a precaution, OS ABI (System V) guarantees
			// that there will be exactly argc number of non-nil valid pointers
			break;
		}

		args->ptr[i-1] = ku_str_from_cstr_ptr(cstr_ptr);
	}

	args->len = i - 1;
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
ku_print(ku_str s) {
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
