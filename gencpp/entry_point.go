package gencpp

import (
	"bytes"
)

func EntryPoint(buf *bytes.Buffer, name, entryPointLinkName, exitName string) {
	// signature
	buf.WriteString(`extern "C" void`)
	buf.WriteString("\n")
	buf.WriteString(entryPointLinkName)
	buf.WriteString("() {\n")

	// call entry
	buf.WriteString("\t")
	buf.WriteString(name)
	buf.WriteString("();\n")

	// call exit
	buf.WriteString("\t")
	buf.WriteString(exitName)
	buf.WriteString("();\n")
	buf.WriteString("}\n")
}
