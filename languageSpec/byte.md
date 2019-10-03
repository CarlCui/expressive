# Byte

## Introduction

Byte data type represents a 8-bit value in expressive. It is really just an i8, but can be also seen as a traditional char (as in a ASCII character).

## Byte literal

A byte literal is seen the same as a character literal, surrunded by single quotes. Depending on the context, the compiler will treat it properly whether to be a byte literal or a character literal. If the context is ambiguous, it shall be treated as a character literal first, rather than a byte literal.

```
let a: byte = 'a'; // 'a' means a byte literal
```