package ast

// <FunctionDeclaration> = "fn" <Name> <Signature>
//
// <Name> = <Identifier>
//
// <Signature> = <FunctionSignature>
type FunctionDeclaration struct {
	Signature FunctionSignature
	Name      Identifier
}

// <FunctionDefinition> = <Head> <Body>
//
// <Head> = <FunctionDeclaration>
//
// <Body> = <BlockStatement>
type FunctionDefinition struct {
	Head FunctionDeclaration

	Body BlockStatement
}

// <FunctionSignature> = <Parameters> [ "=>" ( <Result> | "never" ) ]
//
// <Parameters> = "(" { <FieldDefinition> "," } ")"
type FunctionSignature struct {
	// Equals nil if there are no parameters in signature
	Params []FieldDefinition

	// Equals nil if function returns nothing or never returns
	Result TypeSpecifier

	// Equals true if function never returns
	Never bool
}

// <FieldDefinition> = <Name> ":" <TypeSpecifier>
//
// <Name> = <Identifier>
type FieldDefinition struct {
	Name Identifier
	Type TypeSpecifier
}

