package genc

const prelude = `
typedef unsigned char u8;
typedef unsigned short int u16;
typedef unsigned int u32;
typedef unsigned long int u64;
typedef __uint128_t u128;

typedef signed char i8;
typedef signed short int i16;
typedef signed int i32;
typedef signed long int i64;
typedef __int128_t i128;

typedef float f32;
typedef double f64;
typedef __float128 f128;

typedef u8 bool;
typedef u32 rune;

#define nil 0
#define true 1
#define false 0

`

func (g *Builder) prelude() {
	g.puts(prelude)
}
