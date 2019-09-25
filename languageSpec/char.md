# Character

## Intro

A character in expressive is a valid utf-8 character. It can be seen as a byte array, but should not be treated as an array under normal circumstances. The user can cast it to a byte array to access its elements specifically.

## Character literal

In expressive, a literal surrunded by single quotes is treated as a character literal by default.

## LLVM implementation

### Variable

A character is an array of _int8_. So in llvm, a character is really just a small string.

### Constant literal

A character constant is implemented as a string constant.
