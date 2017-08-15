# Control flow

## if

_ifStmt_ := `if` _ifCondition_ _blockStmt_ (_elseIfStmt_)* (_elseStmt_)?

_ifCondition_ := `(` _expr_ `)`

_blockStmt_ := `{` _stmts_ `}`

_elseIfStmt_ := `else` `if` _ifCondition_ _blockStmt_

_elseStmt_ := `else` _blockStmt_



## switch

_switchStmt_ := `switch`


```
switch (foo) {
    case fooIs5():
        break;
    case 3:
        break;
    default:
}
```

```
switch (true) {
    case foo == 5:
        break;
    case foo > 6:
        break;
    default:
}
```

## while

```
while (i < 4) {

} 
```

## for

```
for (i = 0; i < 5; i ++) {

}
```

## for in

```
for i in someList {

}
```