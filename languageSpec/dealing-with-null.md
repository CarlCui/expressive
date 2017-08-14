# Dealing with null

## Nullable

Every datatype can be nullable, even the primary ones like int.

## Checking null

_nullCheckExpr_ := `check` _nullCheckParamList_ .

_nullCheckParamList_ := _expr_ (`,` _expr_)* .


e.g.

1. `check foo;`
1. `let valid = check foo;`
1. `let allValid = check foo, bar.field1, bar.field2;`
1. `check foo[1]`
1. `check foo[1: 3]`
1. `check foo.map(bar -> parseInt(bar))`

```
func doSomething(input: someStruct) {
    let valid = check input.handler.values;

    if (!valid) {

    }
}
```

With optional checking:

```
func doSomething(input: someStruct) {
    if (input?.handler?.values?.length > 0) {

    }
}
```

However, in this case, the behaviour is not obvious:

1. If anything in that expression is null, the evaluated value is null.
1. If length can be evaluated, the final value is boolean.