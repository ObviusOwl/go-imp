# IMP Assembly and instructions

Currently there exist two data types the VM can work with. Integers and strings.

In assembly integers are represented with the digits 0-9 and an optional minus sign.

Strings are represented with escaped double quotes following the standard rules 
any modern programing languages uses.

Boolean values are represented as integers. This works well with the jump instructions.
A false is encoded as 0 and true is anything else.

The assembly file format is intended to work with common markdown code block syntax 
hightlighter understanding the NASM assembly (for x86). Also it is pretty easy to read.

There is one instruction per line. Each line starts with whitespace or the name of an 
instruction. As the number or parameters is limited to 0 or 1, there may be a single 
argument depending on the instruction. After the argument there can be whitespace or 
a comment. 

Comments are started with a semicolon `;` and go until the end of the line.

## Original instruction set

This instruction set covers the original VM specification. The assembly mnemonics
however are new.

### Control flow instructions

`lab l` Implemented as the `Label` type. Sets a label for jump instructions. A label
is a unique arbitrary integer (separate from addresses and indices). The label is a 
no-op instruction.

`jmp l` Implemented as the `Jump` type. Uncionditionally jump to the label l.

`jnz l` Implemented as the `JumpNonZero` type. Jump to label l if the poped stack 
content is not 0.

`jez l` Implemented as the `JumpZero` type. Jump to label l if the poped stack 
content is exacclty 0.

`stp l` Implemented as the `Stop` type. Stop the execution imediately.

### Arithmetic instructions

`add` Implemented as the `Add` type. Pop two integers from the stack *add* them and 
push the result to the stack.

`min` Implemented as the `Minus` type. Pop two integers from the stack *subtract* 
them and push the result to the stack.

`div` Implemented as the `Div` type. Pop two integers from the stack *divide* them 
and push the result to the stack.

`mul` Implemented as the `Mult` type. Pop two integers from the stack *multiply* 
them and push the result to the stack.

### Logic instructions 

`eql` Implemented as the `Equal` type. Pop two values from the stack and push integer 
1 if they are equal and 0 if they differ.

`gtt` Implemented as the `Greater` type. Pop two integers from the stack and push integer 
1 if the first is greater and 0 otherwise.

`ltt` Implemented as the `Lesser` type. Pop two integers from the stack and push integer 
1 if the first is lesser and 0 otherwise.

### Data move instructions

`psh n` Implemented as the `PushInt` type. Push the number n on top of the stack.

`stm a` Implemented as the `StoreMemory` type. Pop a value from the stack and store 
it in memory on address a.

`ldm a` Implemented as the `LoadMemory` type. Load the value from memory address a
and push it to the stack.

`out a` Implemented as the `Output` type. Print the content of the memory at address a
on the console as a line.

## Strings extension

This extension adds instructions to work with strings.

`str s` Implemented as the `PushStr` type. Push the given string on the stack.

`len` Implemented as the `LengthStr` type. Pop a string from the stack and push the 
length of the string as an integer to the stack.

`cat` Implemented as the `ConcatStr` type. Pop two strings from the stack, concatenate
them and push the result on the stack.

`fmt s` Implemented as the `FormatStr` type. Interpret the given string as 
[Go formatting syntax](https://pkg.go.dev/fmt). Pops as many value sfrom the stack 
as there are unescaped %-signs. Pushes the formatted string on the stack.
