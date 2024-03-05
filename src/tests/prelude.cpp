namespace coven {

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

// Platform architecture dependent unsigned integer type
//
// Always has enough bytes to represent size of any memory
// region (or offset in memory region)
typedef u64 uarch;

// Platform architecture dependent signed integer type
//
// Always has the same size as ua type
typedef i64 iarch;

// Type for performing arbitrary arithmetic pointer operations
// or comparisons
typedef uarch uptr;

// Rune represents a unicode code point
typedef u32 rune;

}
