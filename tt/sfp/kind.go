package sfp

// Kind indicates special flow point kind
type Kind uint8

const (
	// Zero value of Kind. Should not be used explicitly.
	//
	// Mostly a trick to detect places where Kind is left unspecified.
	empty Kind = iota

	// Indicates execution flow which typically occurs in a block without SEPs.
	// This value is never assigned as kind to SEP. Example:
	//
	//	if x < 0 {
	//		print("reverse x");
	//		x = -x;
	//	}
	//
	// In this case body of if statement branch has execution flow of Pass.
	Pass

	// Indicates execution flow in a block with non-trivial SEP logic.
	// This value is never assigned as kind to SEP. Example:
	//
	//	{
	//		let x: int = a;
	//		if x < 0 {
	//			print("reverse x");
	//			x = -x;
	//			return x;
	//		}
	//
	//		x += 1;
	//		print("x is positive");
	//	}
	//
	// In this example the outer block has Comp flow, because it can either return or pass,
	// depending on runtime value of variable a.
	Comp

	// Marks a SEP which is a return from a function.
	Ret

	// Marks a SEP which is a call to a function that never returns.
	Never

	// An endless for loop without SEPs.
	For

	// Jump to next loop iteration.
	Next

	// Jump to immediate outside end of loop.
	Out
)

var text = [...]string{
	empty: "<nil>",

	Pass: "pass",
}

func (k Kind) String() string {
	return text[k]
}
