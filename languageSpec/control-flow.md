# Control flow

## if

_ifStmt_ := `if` _ifCondition_ _blockStmt_ (_elseIfStmt_)* (_elseStmt_)?

_ifCondition_ := `(` _expr_ `)`

_blockStmt_ := `{` _stmts_ `}`

_elseIfStmt_ := `else` `if` _ifCondition_ _blockStmt_

_elseStmt_ := `else` _blockStmt_



## switch

_switchStmt_ := `switch` `(` _expr_ `)` `{` _switchCase_* _defaultCase_? `}`

_switchCase_ := `case` _expr_ `:` _stmt_* _breakStmt_?

_breakStmt_ := `break` `;`

_defaultCase_ := `default` `:` _stmt_*


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

_whileStmt_ := `while` `(` _expr_ `)` _blockStmt_ .

```
while (i < 4) {

}
```

## for

_forStmt_ := `for` `(` _forExpr_ | _forInExpr_ `)` _blockStmt_

_forExpr_ := _forInitialization_? `;` _forConditionExpr_? `;` _forIterationStmt_?

_forInitialization_ := _assignmentStmt_ | _declarationStmt_

_forConditionExpr_ := _expr_

_forIterationStmt_ := _stmt_

_forInExpr_ := _identifier_ (`,` _identifier_)? `in` _expr_

### normal for

```
let i: int;
for (i = 0; i < 5; i ++) {

}
```

or

```
for (let i = 0; i < 5; i ++) {

}
```

### for in

```
for (ele in someList) {

}
```

```
for (ele, i in someList) {

}
```

## break

_breakStmt_ := `break` `;`

A break statement must be inside of a for, while or switch block. Upon called, it will branch out to the end of the most inner for, while or switch block.

e.g.

```
for (let i = 0; i < 4; i ++) { // first loop
    for (let j = 0; j < 5; j ++) { // second loop
        break; // this break will reach (1), not (2)

        // some other code


        // (1)
    }


    // some other code

    //(2)
}
```

