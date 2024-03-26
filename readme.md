# Gizmo

This repo contains compiler, build tools, package manager and some other utilities
to work with Gizmo programming language. Its syntax and set of features is inspired by Rust, Go, Zig and
few other modern languages

Currently gizmo is transpiled to a very restricted subset of C++

## Example snippets

```gizmo
// Loop over builtin iterator range
fn for_each_loop_example(k: int) => int {
    var s: int = 0;
    for let i in range(k) {
        s += i;
    }
    return s;
}
```

```gizmo
type Stack struct {
    s: [16]i32,
    p: uint,
}

fn make_empty_stack() => Stack {
    let s: Stack = {
        s: dirty,
        p: 0,
    };

    return s;
}

// builtin symbol rv (short for method receiver) resolves to object instance
fn [Stack] push(x: i32) {
    rv.s[rv.p] = x;
    rv.p += 1;
}

// under the hood it is just a function with one extra argument
fn compiler_expanded_method_push(rv: *Stack, x: i32) {
    rv.s[rv.p] = x;
    rv.p += 1;
}
```
