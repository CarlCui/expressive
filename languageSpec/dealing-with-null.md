# Dealing with null

## Nullable

Every datatype can be nullable, even the primary ones like int.

## Check operator

_nullCheckExpr_ := `check` _nullCheckParamList_ .

_nullCheckParamList_ := _expr_ (`,` _expr_)* .


e.g.

1. `check foo;`
1. `let valid = check foo;`
1. `let allValid = check foo, bar.field1, bar.field2;`
1. `check foo[1]`
1. `check foo[1: 3]`
1. `check foo.map(bar -> parseInt(bar))`

### Use case:

Consider this senario. We want to access a `length` property of `input.handler.values`. Without null-checking, it would potentially causing a run-time error of referencing from null. So we need to check null before accessing the value.

With `check`, we can do this:

```
func doSomething(input: someStruct) {
    let valid = check input.handler.values;

    if (valid && input.handler.values.length > 0) {

    }
}
```

or

```
func doSomething(input: someStruct) {
    let valid = check input.handler.values;

    if (!valid) {
        //report error
    }

    // do business logic
    if (input.handler.values > 0) {

    }
}
```

Meanwhile, with optional chaining in C# and swift, we need to do this:

```
func doSomething(input: someStruct) {
    if (input?.handler?.values?.length > 0) {

    }
}
```

However, in this case, the behaviour is not obvious:

1. If anything in that expression is null, the evaluated value is null.
1. If length can be evaluated, the final value is boolean.

To conclude, the result type is not deterministic at compile time, and we have to allow implicit conversion from `null` to `boollean` in order to make this work in this if condition.

### Checking multiple fields at once

In some cases, we want to access multiple fields of an object, and we need to check null of them.

With optional chaining:

```
func doSomething(input: someStruct) {
    if (input?.handler?.values?.length > 0) {
        if (input?.source?.fileName?.contains("abc")) {
            // business logic
        }
    }
}
```

With `check`:
```
func doSomething(input: someStruct) {
    let valid = check input.handler.values,
                      input.source.filename;
    
    if (!valid) {
        // report error
    }

    // access normally and safely
    if (input.handler.values.length > 0) {
        if (input.source.fileName.contains("abc")) {
            // business logic
        }
    }
}
```