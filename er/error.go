package er

import "io"

type Error interface {
	error

	Kind() Kind
	Format(w io.Writer, format string) error
}

// Kind indicates compiler error kind
type Kind uint32

const (
	// Zero value of Kind. Should not be used explicitly.
	//
	// Mostly a trick to detect places where Kind is left unspecified.
	// Upon encountering this value most functions and methods will
	// (and should) panic.
	empty Kind = iota

	UnexpectedToken
	OrphanedAtrBlock

	UnknownImportOrigin
	EmptyImportString
	BadImportString
	SelfImport
	SameStringImport
	ImportCycle
	EmptyUnit

	UnresolvedSymbol
	NotFoundInUnit
	NotExportedSymbol
	MultipleSymbolDeclarations
	MultipleSymbolDefinitions
	MultipleArgumentDeclarations
	NoExplicitInit
	ConditionTypeMismatch
	NotAssignableExpression
	NotCallableSymbol
	TooManyArguments
	NotEnoughArguments
	NotTypeSymbol

	OperatorNotDefined
)
