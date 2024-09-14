static void
ku_print_proc_args(ku_proc_args args) {
	for (uint i = 0; i < args.len; i += 1) {
		ku_str s = args.ptr[i];
		ku_print(s);
		ku_print(ku_static_string((const u8*)(u8"\n"), 1));
	}
}
