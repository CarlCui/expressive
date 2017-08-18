# List

## Intro

A list data structure is basically an c array with slicing and other built-in functionalities.

## Creation

To create

`let aList: int[] = make(int[]);`

`let aList: int[] = new int[];`

`let aList: int[23];`

## Quick Initialization

Lists can be easily populated with `[` `]`, like in javascript. However, it is still statically typed. If the type of the list is not provided, it is inferred from the type of the elements. In this situation, the elements must have the same type.

`let aList = [1, 2, 3];`
`let aList: int[] = [];`

## Manipulation

`aList + bList`
`aList + [1, 2, 3]`
`[1] + aList`

## Append

`aList.append(1, 2, 3)`



## functional programming with lists

`aList.map((ele, index) -> ele + 1)`
`aList.filter`
`aList.find`

##

```
func createNewList() {
    return [1, 2, 3];
}
```
