# Gizmo

This repo contains compiler, build tools, package manager and some other utilities
to work with Ku programming language. Its syntax and set of features is inspired by
Rust, Go, Zig, Odin and few other modern languages

Currently Ku is transpiled to a very restricted subset of C

## Example snippets

```ku
// Loop over builtin iterator range
fun for_each_loop_example(k: sint) => sint {
    var s: sint = 0;
    for i in range(k) {
        s += i;
    }
    return s;
}

type Stack struct {
    s: [16]s32,
    p: uint,
}

fun make_empty_stack() => Stack {
    let s: Stack = {
        s: dirty,
        p: 0,
    };

    return s;
}

// builtin symbol rv (short for method receiver) resolves to object instance
fun [*Stack] push(x: s32) {
    g.s[s.p] = x;
    g.p += 1;
}

// under the hood it is just a function with one extra argument
fun compiler_expanded_method_push(g: *Stack, x: s32) {
    g.s[g.p] = x;
    g.p += 1;
}
```
