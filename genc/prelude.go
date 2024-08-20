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

_Noreturn void
ku_trap_unreachable() {
	__builtin_trap();
	__builtin_unreachable();
}

void
ku_must(bool c) {
	if (c) {
		return;
	}

	ku_trap_unreachable();
}

_Noreturn void
ku_panic_never(u64 pos) {
	ku_trap_unreachable();
}

`

func (g *Builder) prelude() {
	g.puts(prelude)
}
