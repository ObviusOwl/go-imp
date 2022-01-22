# Go implementation of the IMP language

This project implements the IMP mini programing language using the Go language.
The IMP language is a strongly and statically typed miniature programming language 
with reduced syntax and feature set intended for learning the workings of compilers. 

This project exists for the purpose of learning, it may or may not get updated. 

Implemented:

- A stack based VM
- An assembly parser for the VM
- A lexer for the IMP language

TODO:

- Types to represent the AST
- A recursive parser
- A type checker and inference
- VM code generator

Documentation:

- [Design of the VM with Go types](./docs/vm.md)
- [The assembly language](./docs/asm.md)
- [Additional features](./docs/features.md)

## Getting started

Compile the Go code:

```sh
go build -o imp main.go
```

Create a file `hello.asm`:

```nasm
psh 3
stm 10

lab 1
ldm 10
fmt "Hello World! #%d"
stm 11
out 11
ldm 10
psh -1
add
stm 10
ldm 10
jnz 1
```

Now run the VM with the asm file:

```
./imp asm -f hello.asm
```

The output is from a simple loop printing hello world.

```
Hello World! #3
Hello World! #2
Hello World! #1
```

The lexer sub-command is present only as a debugging tool. The lexers output types 
are meant as input to the parser. To test the lexer, create a file `fib.imp`

```
a := 0; b := 1; c := 1
while count < 5 {
    c = a + b
    print c
    a = b
    b = c
    count = count + (-1)
}
```

And run the lexer:

```
./imp lex -f fib.imp
```

The output is the string represnetation of the tokens:

```
<identifier 'a'>
<op ':='>
<int '0'>
<identifier 'b'>
<op ':='>
<int '1'>
<identifier 'c'>
<op ':='>
<int '1'>
<keyword 'while'>
<identifier 'count'>
<op '<'>
<int '5'>
```

## Development

Build the main binary

```sh
go build -o imp main.go
```

Run the unit tests:

```sh
go test ./...
```

Build the documentation as HTML and PDF with sphinx:

```sh
cd docs
make html
make latexpdf
```
