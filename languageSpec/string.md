# String

## Intro

A string is an array of expressive characters, meaning that the user does not need to parse UTF-8 characters specifically with string operations. Strings are immutable, as in all other programming languages.

## Production

_stringConstLiteral_ := `"` _stringWord_ `"`

_stringWord_ := [^"\t] | `\"`

_stringInterpolation_ := `` ` `` ( _exprInterpolation_ | _stringWord_ )* `` ` ``

_exprInterpolation_ := `$` (_noWhitespace_) `{` _expr_ `}`

## creation

`let aString = "123";`

## concat

`"123" + "456"`

## interpolation

```
let aString = `${variable} + 3`;
```
