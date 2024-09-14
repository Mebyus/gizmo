
u8 hex_digit(u8 x) {
	if (x < 10) {
		return (u8)'0' + x;
	}
	return x - 10 + (u8)'a';
} 

void fmt_hex_byte(u8 *buf, u8 b) {
	u8 high = b >> 4;
	u8 low = b & 0xf;
	
	buf[0] = hex_digit(high);
	buf[1] = hex_digit(low);
}

void print_mem(void *p, uint lines) {
	u8 buf[1 << 14];
	uint n = 0;
	u8 *m = (u8*)p;
	for (uint i = 0; i < lines; i += 1) {
		for (uint j = 0; j < 16; j += 1) {
			fmt_hex_byte(buf + n, m[16 * i + j]);
			n += 2;
			buf[n] = ' ';
			n += 1;
		}
		buf[n] = '\n';
		n += 1;
	}

	ku_syscall_write(1, buf, n);
}


void start() {
	// u64 tip = 0;
	// u64 *a = &tip;
    void *p;
    __asm__ volatile (
        "movl %%rsp, %0" : "=r" (p)
    );
	print_mem(p, 10);
	ku_print(ku_static_string((const u8*)(u8"\nHello, world!\n"), 15));
	ku_syscall_exit(0);
}
