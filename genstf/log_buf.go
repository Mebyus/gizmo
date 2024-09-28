package genstf

const logbuf = `
let LOG_BUFFER_SIZE = 1 << 16;
var logbuf: [LOG_BUFFER_SIZE]u8;

`

func (g *Builder) logbuf() {
	g.puts(logbuf)
}
