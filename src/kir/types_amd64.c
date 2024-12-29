// Unsigned integers of fixed size.
typedef unsigned char u8;
typedef unsigned short int u16;
typedef unsigned int u32;
typedef unsigned long int u64;
typedef __uint128_t u128;

// Signed integers of fixed size.
typedef signed char s8;
typedef signed short int s16;
typedef signed int s32;
typedef signed long int s64;
typedef __int128_t s128;

typedef float f32;
typedef double f64;
typedef __float128 f128;

typedef _Bool bool;

// Represents unicode point.
typedef u32 rune;

// Arch dependant unsigned integer (same size as pointer).
typedef u64 uint;

// Arch dependant signed integer (same size as pointer).
typedef s64 sint;

// Arch dependant pointer in raw integer form. Meant for performing
// arbitrary pointer manipulations.
typedef uint pint;

// This should only be used with pointer types.
#define nil 0

// These should only be used with boolean type.
#define true 1
#define false 0

// Macro for type cast with fancy syntax.
#define cast(T, x) (T)(x)

_Noreturn static void
panic_trap(void) {
	__builtin_trap();
	__builtin_unreachable();
}
