# Error handling

## Productions

_tryCatchBlock_ := `try` _blockStmt_ `catch` `(` _identifier_ `:` _typeLiteral_ `)` _blockStmt_

_blockStmt_ := `{` _stmts_ `}`

_tryExpr_ := `try?` _throwableExpr_

_tryAssignmentStmt_ := _identifier_ `catch` _identifier_ `=` _tryExpr_ `;`

## Throwable

When a function could potentially throw an exception, it has to be marked with keyword `throwable`. 

```
func throwable parseInt(input: string) -> int {
    let result = // try parse

    if (result == null) {
        throw SomeException;
    }

    return result;
}
```

## Catching error

### try-catch block
```
try {
    let val = parseInt(input);
} catch (exception: SomeException) {

}
```

### try?

`let val catch exception = try? parseInt(input);`

`exception` will capture the possible exception in the result of the function. Otherwise `null`.

or

`let val = try? parseInt(input);`

If there is an exception throwed in the function, the result is simply `null`.