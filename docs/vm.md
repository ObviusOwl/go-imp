# Virtual Machine

## Stack based VM and hardware CPU

The stack based virtual machine (VM) runs a program of simple instructions generated
from the abstract syntax tree (AST) that is parsed from the IMP source code. 

The VM allows abstraction from the host programing language (the one used to implement 
the compiler) and the platform used. For example the VM code generated from the 
parser implemented in Go may be executed in a webbrowser on a VM implemented in 
javascript. The implementation of the VM is simple and yet powerful. The tradeoff 
is a slower execution, which we ignore in this project.

There exist obvious parallels between a real CPU and the virtual machine. Thus the 
nomenclature is borrowed from the hardware CPU architectures.

This VM implementation in Go uses a program counter (pc), a pointer variable pointing 
to the location in memory of the current (or next) instruction. On a hardware CPU 
this would be a special register with the same purpose. The pointer can be manipulated
to implement (conditional) jumps in the program, which ultimately implements the 
turing completeness of the instruction set. 

The instructions the VM executes are similar to the instructions a CPU uses. 
An instruction has an identifier that specifies which of the pre-defined behaviours
the CPU should execute next. It also contains a list of parameters influencing
parts of this behaviour. The full list of instruction types the CPU can run is 
called the instruction set, which is described in detail in the manual of the CPU.
The VM also has such an instruction set.

A hardware CPU needs the instructions to be encoded as a large blob of bytes. This is not 
easy to reason about and even less to program by hand. Thus for each CPU architecture
there exist one or more standard ways to translate between the binary form the 
CPU can run and the text form a human can read and edit. These text representations 
are called assembly language (short form: asm).

The VM uses the data types the host programing language provides to represent an 
instruction. These types are the same any program implemented in this language can use. 

Unlike a CPU, the instructions are not decoded on the fly, there is no pipelining 
and no fetch phase. With a general purpose programing language available, it is 
easier to handle the instructions as an array of objects that is fully loaded into 
memory when a program is started. 

If the VM instructions are generated on the fly by the same process that also runs
said instructions in it's implementation of the VM, the data types can be instanciated 
directly using the standard means of allocating memory in the host progaming language.
This is the case when the compiler for the IMP language also runs the program, similar 
to how `go run main.go` directly runs the program and `go build main.go` only produces 
the binary program. This is also the case for most unit tests that test the VM 
implementation.

To be able to load a program of instructions from a file, there must be some sort 
of representation that is not the memory of a running process. The VM implementation
must translate this representation into it's native representation. The format of 
this file should be independent of the programing language the VM is implemented in.
There could be a purely binary file representation, however as noted above, this is 
not practical. Instead the VM directly interprets the assembly code that is designed 
for this VM. 

The instructions the VM can execute accept exactly zero or one parameter. This makes 
implementing the instructions and the parser for the assembly language simple. It 
also ensures that the stack based VM does not mutate into a direct interpreter of 
the abstract syntax tree for IMP, which due to the learning nature of this project 
is quite valable.

A hardware CPU uses, depending on it's architecture, registers and the system memory 
to store the temporary results of computations (instructions). For example the 
two operands for an addition are loaded from two registers and the result is written
to a separate or same register. A register is a memory embedded into the CPU. Usually
there are only a handful available.

A stack based CPU (or VM) uses instead a stack, which is a 
[very common data type](https://en.wikipedia.org/wiki/Stack_(abstract_data_type)).
In a first step the operands are pushed to the stack. Then, when the instruction 
using the operands is executed, the operands are poped (removed) and the result 
is pushed onto the stack so that the next instruction can use it as it's input.

How many operands each instructions expects to be available on the stack, how 
many values are pushed to the stack by the instruction and what the instruction 
actually computes with the data is part of the instruction set manual. The behaviour 
of each instruction regarding the stack must be implemented in the VM.

Thus all instructions have in common that they pop from the stack, do something 
with the values and push back values to the stack. This is a very convenient 
interface to implement as the logic needed to implement each instruction is kept 
to the very core with a very small and standardized dependency on the environment
they are run in.


## The VM written in Go

As discussed in the previous chapter the interfaces to implement are very uniform
and the actual implementation logic of each instruction can be kept to its intrinsic
details.

Each instruction must implement the `Executer` interface, which consists of the 
single method `Exec()`. This method takes as argument a type implementing the 
`Runner` interface, one implementing the `Stack` interface and one implementing
`Memory`. Any type providing this method can be used as instruction in the VM program. 
This is a very idiomatic Go implementation as the interface can be combined freely with 
any other interface. [Interfaces](https://golangdocs.com/interfaces-in-golang) 
is how Go implements polymorphism. 

The `Program` type is based on a slice of types implementing `Executer`. Thus 
a program is simply a list of instructions, nothing more nothing less.

Each instruction needs access to the stack. Any instruction needs at most two types 
of operations on it: `Push()` and `Pop()`. A type implementing the `Push` method 
is a `Pusher` and any type implementing `Pop()` is a `Poper`. When both methods 
are present the type can be effectively used as a `Stack`. The separation of the 
methods into separate interfaces with a single method allows functions using only
one of the method to be more versatile and composable. Most instructions need 
both methods. 

Using the interfaces as polymorphic abstraction allows any `Stack` implementation 
to be used with any instruction (`Executer`), regardless from where the implementation 
comes. This design allows third party software to load the VM as dependency into 
their own code, implement their own stack and extend the instruction set with their
own instructions. No instruction type must be known ahead or hardwired into the VM.
There is however a base set of `Executer` shipped with the VM.

An instruction must also be able to influence the control flow, in particular what 
instruction is executed next by the VM, even without knowing how the core logic of 
the VM works. This is implementd via the `Runner` interface, which serves a 
control unit.

The `Memory` interface describes how the instruction can access the system memory 
to store more permanent data. Only two methods are needed `Load()` and `Store()`. 
Any type having these methods can do the job, regardless how the data is stored: 
in files, a database over the network or in memory. 

The three interfaces `Runner`, `Stack` and `Memory` provide separation of consers 
and polymorphism, which facilitate composition and thus effective code reuse, 
high maintability, extensibility and flexibility. The three interfaces can also 
be implemented by the same type, as it is with the mocked type used in the unit tests. 

The `Memory` and the `Stack` hold both values of the type `DataValue`, which is 
based on the empty interface. In go the empty interface is implemented by all 
types. It is similar to the `object` type in Java. The data types available to 
by the instructions should not be restricted a priori. For example the instructions
handling strings were added as an extension. They need to store Go strings. 
This can be done by abstracting the storage as array of bytes, but it is much 
more convenient, faster and safer to use Go 
[type assertions](https://golangdocs.com/type-assertions-in-golang) 
and [type switches](https://golangdocs.com/type-switches-in-golang) 
to convert the types back and forth.

Here is a simplified version of the VM implementation:

```go
type Pusher interface {
	Push(item interface{})
}
type Poper interface {
	Pop() (interface{}, error)
}
type Stack interface {
	Pusher
	Poper
}

type MapMemory map[int]DataValue
type DefaultRunner struct {
    program Program
    pc      int
    stack   Stack
}
type Machine struct {
    ctrl DefaultRunner
    mem  MapMemory
}

type Executer interface {
    Exec(vm Runner, st stack.Stack, mem Memory)
}
type Program []Executer

type Memory interface {
    Load(address int) DataValue
    Store(address int, value DataValue)
}
type Runner interface {
    Run(program Program, mem Memory)
    Jump(label Label)
}
type Label int
type DataValue interface{}

func (ctrl *DefaultRunner) Jump(label Label){
    for idx, inst := range ctrl.program {
        if value, ok := inst.(Label); ok && value == label {
            ctrl.pc = idx
        }
    }
}

func (ctrl *DefaultRunner) Run(program Program, mem Memory){
    for ; ctrl.pc < len(ctrl.program); ctrl.pc++ {
        ctrl.program[ctrl.pc].Exec(ctrl, ctrl.stack, mem)
    }
}

func (mem MapMemory) Load(address int) DataValue {
    return mem[address]
}

func (mem MapMemory) Store(address int, value DataValue) {
    mem[address] = value
}
```

Here is a simplified version of three instructions:

```go
type JumpZero Label
func (inst JumpZero) Exec(vm Runner, st stack.Stack, mem Memory){
    item, _ := st.Pop()
    if value, ok := item.(int); ok && value != 0{
        vm.Jump(Label(inst))
    }
}

type Add struct{}
func (inst Add) Exec(vm Runner, st stack.Stack, mem Memory){
	item1, _ := st.Pop()
    if op1, ok := item1.(int); ok{
    	item2, _ := st.Pop()
        if op2, ok := item1.(int); ok{
            st.Push(op1 + op2)
        }
    }
}

type LoadMemory int
func (inst LoadMemory) Exec(vm Runner, st stack.Stack, mem Memory){
	st.Push(mem.Load(int(inst)))
}
```