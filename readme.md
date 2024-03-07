# Gizmo

This repo contains compiler, build tools, package manager and some other utilities
to work with Gizmo programming language. Its syntax and set of features is inspired by Rust, Go, Zig and
few other modern languages

Currently gizmo is transpiled to a very restricted subset of C++

## Example snippets

```gizmo
fn for_each_loop_example(k: int) => int {
    var s: int = 0;
    for i in range(k) {
        s += i;
    }
    return s;
}
```
