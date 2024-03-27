# `glox` Interpreter
A golang implementation of the Lox language from *Crafting Interpreters* by Robert Nystrom.

**Why Go?**
I don't care for Java, wanted to learn Go, and didn't want to copy code snippets directly from the book.

# Test
Use environment variables to update or clean snaps.
`UPDATE_SNAPS=true` or `UPDATE_SNAPS=clean` respectively.
```
go test -v ./...
```

# Build
```bash
go build
```

# Run
```bash
./glox examples/fib.lox
```
