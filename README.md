# dywoqlang

dywoqlang is a language-interpreter written in Go.

# Syntax

The language uses `.dl` as file-extension.

Simple main program:
```dl
main args string[] int64 {
	make int8 result, 30+30 # Create 8 bit integer with name "result"
	stdout result, "\n"     # Output a integer with new line.
	mov ret, 0              # Assign 0 to ret variable
	ret                     # Return 0
}
```
