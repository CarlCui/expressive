# Character

## Intro

A character in expressive is a valid utf-8 character. It can be seen as a byte array, and it's immutable, like string.

## LLVM implementation

In llvm, an expressive character constant is implemented as a string constant. A variable with character type is a byte pointer.
