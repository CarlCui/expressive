# Type inference

To eliminate some boilerplate, expressive allows implicit type inference during certain situations.

# Variable declaration

During variable declaration, you don't need to specify the type of the declared variable. In that case, the type of that variable is inferred from the right hand side expression.

Example:
```
let foo = 1 + 4; // foo will be int
```
