# List

## Intro

A list data structure is an array-like data structure with dynamic length and slicing capabilities.

## Type Literal

See [types](./types.md) for more details.

## Creation

To create

`let aList: int[] = make(int[]);`

## Quick Initialization

Lists can be easily populated with `[` `]`, like in javascript. However, it is still statically typed. If the type of the list is not provided, it is inferred from the type of the elements. In this situation, the elements must have the same type.

`let aList = [1, 2, 3];`
`let aList: int[] = [];`

## Getting length

Represents the length of the content stored in the list

`aList.length()`

## Getting size

Represents the storage size of the current list

`aList.size()`

## Manipulation

```
aList + bList
aList + [1, 2, 3]
[1] + aList
```

## Slicing

```
let aList = [1, 2, 3, 4, 5];
aList[0:2]; // 1, 2
aList[3:4]; // 4

Below are all the same
aList[0:5]; // 1, 2, 3, 4, 5
aList[0:];
aList[0:-1];
aList[:-1];

aList[0:2:5]; // 1, 3, 5
aList[-1]; // 5
aList[0:2] + aList[2:3] // 1, 2, 3
```

## Append

`aList.append(1, 2, 3)`


## functional programming with lists

`aList.map((ele, index) -> ele + 1)`

`aList.filter`

`aList.find`

## Examples

```
func createNewList() {
    return [1, 2, 3];
}
```

## Implementation details


### Header

Under the hood, a list is, unlike in C, a complex data structure that stores more than what the user can access.


```
      | header        | content | termination
      | length | size | content | 0
bytes |   4    |  4   |   ...   | 1
```