# Solid and pass values

From language semantic standpoint not all forms of expressions can be used
interchangeably in different circumstances.

To illustarate the above statement, consider the following example:

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

Code above is semantically correct (although a bit silly) and contains several
expressions (in order we encounter them):

```ku
// from function init()
b.bar    // 1 - field select
"mark"   // 2 - string literal
b.a.foo  // 3 - field select via chain 
153      // 4 - integer literal
b        // 5 - symbol usage

// from function example()
b        // 6
init()   // 7  - function call
b.&      // 8  - take variable address
b.bar    // 9
"hello"  // 10
p.@      // 11 - indirect on pointer variable
init()   // 12
p.a      // 13 - indirect field select
b.a      // 14
```

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

Both `p` and `b.&` have the same type `*B` Could we write statement `b.& = p;`?
Also "no". Even the thought about such operation is strange! We cannot assign
a value to an address obtained from a variable. Value of expression `b.&` could
be passed somewhere, but cannot be assigned to. Just like function call result!
