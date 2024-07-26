// package stg (Symbol Type Graph) contains algorithms and data structures
// to perfrorm semantic analysis of source code. In short it accepts AST
// (produced by parser) as input and walks it several times, gathering
// information about used symbols and types. In the end AST is transformed
// to STG which represents source code as a graph with detailed information
// about types and symbols used in source code.
//
// In broad strokes transformation does the following:
//
//	- symbol resolution based on name and lexical scope
//	- type spec resolution
//	- unit-level symbol hoisting
//	- methods to types bindings
//	- symbols and expressions typing
//	- expressions and statements type checking
package stg
