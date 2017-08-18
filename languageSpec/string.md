# String

## Intro 

Strings are immutable, as in all other programming languages. It means that any operation on string will create a new one instead of modifying the existing one.

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
