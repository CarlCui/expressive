# Syntax for BNF

The syntax used in this language spec is a mix of traditional BNF syntax and various regexp syntaxes.

1. symbol: _symbol_ 
1. keyword: `keyword`
1. production: _symbol_ := expression
1. character set: [ ] (e.g. [A-Z])
1. grouping: ( ) (e.g. (abc) )
1. or: | (e.g. abc|cde means symbol abc or symbol cde)
1. optional: ? (e.g. (abc)? means symbol abc is optional)
1. 0 or more: * (e.g. (abc)* means symbol abc can be 0 or repeating any times)
1. 1 or more: + (similar to * except that the minimal is 1)
