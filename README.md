# dywoqlang

dywoqlang is a language-interpreter written in Go, a concept of 'high level assembly'.

# Syntax

The language uses `.dl` as file-extension.

Simple main program:
```dl
add int32 (int32, int32): # int32 is the function returning type, while (int32, int32) are types for parameters
	# Registers param1 and param2 are automatically generated
	mov li1, param1 # We assign parameter one to i1
	mov li2, param2 # We assign parameter two to i2
	add li2, li1    # We add i1 to i2
	ret li2         # And return i2 

# The difference between lN and N registers is
# that the first are local to function, not affecting the global ones: iN. 
# In this example, liN and iN registers are needed to work with integers.

main int32:   
	# Move values to the function as arguments
	mov param1, 10
	mov param2, 10 

	# And call add
	mov i1, call add

	# Outputting a register i1. If any of iN registers are empty,
	# it will simply output <empty>. 
	write i1
	ret 0 
```

If you want to initialize global variables, you need to do this:
```dl
return_code int8 10
hi_dywoqlang string "Hi, dywoqlang!"
```

To use them, you need to explicitly put their names into []:
```dl
main int32:
	mov li1, [return_code]
	mov lstr1, [hi_dywoqlang]

	write li1
	write lstr1
```
