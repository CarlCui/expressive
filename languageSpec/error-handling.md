# Error handling


## Throwable


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

`let val, exception = try? parseInt(input);`

`exception` will capture the possible exception in the result of the function. Otherwise `null`.

or

`let val = try? parseInt(input);`

If there is an exception throwed in the function, the result is simply `null`.