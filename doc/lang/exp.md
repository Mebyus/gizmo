# Expressions

## Chain expressions

When it comes to expression design choices many languages tend to analyze and
structure them based on broad categories of concepts:

- operand
- unary expression (and operator)
- binary expression (and operator)
- function call
- index expression
- ...

These forms of expressions serve as building blocks for constructing other expressions
of arbitrary complexity.

In Ku we decided to separate some operators into a special category. When these
operators are applied to operands they form what language calls **chain operands**.

Let's look at them in example first:

```ku
type A struct {
    foo: u64,
}

type B struct {
    a:   A,
    bar: str,
}

fun init() => B {
    var b: B;

    b.bar = "mark";
    b.a.foo = 153;

    return b;
}

fun example() => A {
    var b: B;
    
    b = init();

    var p: *B = b.&;

    b.bar = "hello";
    p.@ = init();
    p.a = b.a;

    return p.a;
}
```

Below we supplied each expression from the snippet with comments which indicate
chain operands and name some of the operations or expressions:

```ku
// from function init()
b.bar    // 1 (chain) - field select
"mark"   // 2         - string literal
b.a.foo  // 3 (chain) - field select via chain 
153      // 4         - integer literal
b        // 5         - symbol usage

// from function example()
b        // 6
init()   // 7  (chain) - function call
b.&      // 8  (chain) - take variable address
b.bar    // 9  (chain)
"hello"  // 10
p.@      // 11 (chain) - indirect on pointer variable
init()   // 12 (chain)
p.a      // 13 (chain) - indirect field select
b.a      // 14 (chain)
```

Here the complete list of chain operations with examples:

```ku
<C> - marks chain operand
<I> - expression which result can be used as index value

=========================================

1. Field select

<C> => <C> "." <field>

Examples:

    a.foo
    b.foo.a

=========================================

2. Call // TODO: move calls to chain terminators

<C> => <C> "(" <args> ")" // <args> - call arguments, may be empty

```

## Mutation model

From language semantic standpoint not all forms of expressions can be used
interchangeably in different circumstances.

To illustarate the above statement, consider the following example:

// TODO: reference example above

Code above is semantically correct (although a bit silly) and contains several
expressions (in order we encounter them):



Expressions **6** and **7** result values have the same type `B`. Does that mean
we could swap them and still have correct code?

The answer is obviously "no". The following assignment

```ku
init() = b;
```

is nonsense! How to even interpret value assignment to a function call result?
Can we execute such code on a real hardware? Maybe it is possible to patch
together some arcane rules to handle such code, but there is no real value in
doing so.

Thus we conclude that not all expressions are equal, even if they result in
values of the same type.

Does this limitation apply only to function calls? The answer is "no". Let's
look at next example:

```ku
var b: B;
var p: *B;
p = b.&;
```

Both `p` and `b.&` have the same type `*B` Could we write a statement `b.& = p;`?
Also "no". Even the thought about such operation is strange! We cannot assign
a value to an address obtained from a variable. Value of expression `b.&` could
be passed somewhere, but cannot be assigned to. Just like function call result!

Those who come from C++ background probably are already familiar with this
phenomenon. In C++ these two types of expression values are called **lvalue**
and **rvalue**. We will use different names for them, but the concept stays
the same.

## Stored and passed values

In Ku programming language expression result value which we could assign to is
called **stored value**. The name hints to us that this value is stored somewhere
and thus it can be changed. The other type of values is called **passed value**,
meaning that such values are only passed around, not stored, and thus could not
be mutated.

> Where does this concept come from? My uneducated theory is that the need for
> such distinction between expression values comes from high-level nature of
> the language. The need to mutate (change) values is integral to any imperative
> structural high-level language and with it comes the necessity to differ what
> could be mutated and what could not.
>
> On the other hand machine code executed by real processor is free of such
> limitation, because it operates based on different concepts. In hardware
> every value is stored somewhere: either in register or memory. Result of any
> operation is always stored, for CPU to pass something it must be stored
> somewhere first.

Now, when we understand the problem, our next step is to classify expressions
that yield stored values. In Ku the classification is somewhat simple and
restrictive:

```text
<S> - marks s-expression - expression which results in stored value
<I> - expression which result can be used as index value

=========================================

1. Every s-expression starts with symbol:

    <S> => <symbol>

Assignment examples:

    a = 10;
    b = "hello";
    my_var = a;

=========================================

2. Selecting a field from s-expression gives another s-expression:

    <S> => <S> "." <field>

Assignment examples:

    a.foo = 10;
    b.foo.a = "hello";

=========================================

3. If s-expression results in a pointer-type value then dereferencing
such pointer yields s-expression:

    <S> => <S> ".@"

Assignment examples:

    a.@ = 10;
    b.foo.@ = "hello";

=========================================

4. When s-expression results in an indexable value then indexing it
produces s-expression:

    <S> => <S> "[" <I> "]"

Assignment examples:

    a[i+2] = 10;
    b.@.foo[0] = "hello";

=========================================

5. Same goes for array-pointers:

    <S> => <S> ".[" <I> "]"

Assignment examples:

    a.[5] = 10;
    b.foo[0].bar.[i] = "hello";

=========================================
```

The above list is complete. Note that any s-expression must also satisfy type
limitations. For example it is not possible to assign a value to symbol which
denotes a function or select a field which does not exist on a type.

Some forms of expressions which could be used for assignment in C, are prohibited
in Ku for such usage:

```C
// 1.
//
// First example is dereferencing a pointer returned from function call.
//
// You cannot do something like this in Ku:
*(get_ptr()) = 5;

// 2.
int a[4] = {0};
int *p = a;

// Another example includes pointer arithmetic (there is no equivalent
// for such operation in Ku anyway).
//
// This is also not valid in Ku:
*(p + 2) = 9;
```

## Addressable values

Operator `.&` takes an address of the value that preceeds the operator. Just like
assignment address operator cannot be used on arbitrary expression. To take an
address of a value the value must be stored somewhere first. Thus it is somewhat
obvious from definitions that `.&` can only be applied to stored values which
we already discussed previously.
