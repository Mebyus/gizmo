// package tt (type tree) contains code and data structures to perfrorm semantic analysis
// of source code. In short it accepts AST (produced by parser) as input and walks it several
// times, gathering information about used symbols and types. If AST is a semantically correct
// program then resulting tree is produced, otherwise a detailed error is returned
package tt
