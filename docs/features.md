# Additional features

## Embedded assembly file

The Go programming language and tooling is known for producing a single statically 
linked binary containing all library dependencies. This simplifies deployment and 
distribution. 

The binary produced by this project can replicate itself while embedding a given 
asembly file into the binary. When launched, this new binary will automatically 
run the embedded IMP program. Additionally the binary will make selected command 
line arguments available from within the IMP assembly program.

This functionality was inspired by the two Stackoverflow posts 
[how to add and use binary data to compiled executable](https://stackoverflow.com/questions/32437714/how-to-add-and-use-binary-data-to-compiled-executable)
and 
[accessing data appended to an elf binary](https://stackoverflow.com/questions/5660792/accessing-data-appended-to-an-elf-binary)

The solutions was only tested on Linux and probably ony works there anyway.
However with some more work, a more platform independent version can be crafted.

One can append binary data to ELF binaries without breaking them. To embed an assembly 
file, the application opens the currently running binary (itself) in read only mode  
by following the link in `/proc/self/exe`. The content is copied into a new file.
Then the assembly text file is copied, appending to the file. The size in bytes of 
the assembly is recorded and is also appended to the new file. Finally a magic 
string is appended.

When the binary is launched, the application checks if the OS is Linux and then 
opens the file `/proc/self/exe` for reading. If there is no magic string at the 
end of the file, the application runs normally. If there is one, the assembly 
file size is read and the assembly file is loaded using the size to seek to the 
correcy offset. 

After the embedded assembly has been loaded, the command line argment parser is programmed 
according to the metadata at the top of the assembly file. The arguments are parsed 
and inserted at the specified VM memory addresses. Finally the VM is started with 
the program.

To be CPU independend the size of the assembly file is written in network byte order
(big endian). Also a fixed width of 8 Byte is used.

The Header in the assembly file is a special comment syntax to specify the name,
data type, target address and default value of the parameters.

```nasm
; normal comments are ignored
; @param who 5 str "world!"
; @param foo 6 int 5
ldm 5           ; load memory at address 5 on the stack 
str "hello "    ; load string on the stack
cat             ; conctatenate 
stm 1           ; store top of stack in memory at addr 1
out 1           ; print the content of memory at addr 1
```

The header comment must come before any instruction and must contain only one
`@param` stanza per line. The fields are separated by space and are all mandatory.
The format is as follows:

```
@param paramName address dataType defaultValue
```

The parameter name is the name of the paramer on the command line. It is case sensitive.
The address is the target address where to store the value. The type must be exactly one 
of the both lower case values `str` or `int`. The default value must be of the type 
declared in the type field. The same parsing rules apply as for string literals in the 
assembly (double quoted and escaped string).

To build the executable run (skip go build if already compiled):

```sh
go build -o imp main.go
./imp asm -f hello.asm -embed imp-hello
```

Then run the new executable:

```sh
./imp-hello -who 'world !'
./imp-hello -who 'me !'
./imp-hello -who 'all !'
```

Which produces following output:

```
hello world !
hello me !
hello all !
```
