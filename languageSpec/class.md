# class

# Productions

_classDeclaration_ := `class` _identifier_ _classDeclarationBlock_

_classDeclarationBlock_ := `{` _classMethodDeclarationStmt_* `}`

_classMethodDeclarationStmt_ := _identifier_ `(` _formalParamList_ `)` (`->` _functionReturn_)? `;`

_classDefinition := `class` _identifier_ _classDefinitionBlock_

_classDefinitionBlock_ := `{` _classDefinitionStmt_* `}`

_classDefinitionStmt_ := _classFieldDefinitionStmt_ | _classMethodDefinitionStmt_

_classFieldDefinitionStmt_ := _identifier_ `:` _typeLiteral_ (`=` _expr)? `;`

_classMethodDefinitionStmt_ := (`public` | `private`) _functionDefinition_

## Declaration

A public class must have a declaration in the declaration block, with only public methods included, acting like an interface.

```
export declare {
    class Foo {
        doA();
        doB();
        doC();
    }
}
```

## Implementation / Definition

Then the implementation follows.

```
class Foo {
    public doA() {

    }
    public doB() {

    }
    public doC() {

    }
}
```

## Inheritance

Expressive is not focusing much on inheritance. 

1. It does not support multiple inheritance. 
1. If no constructor is given in the child class, the parent's constructor will be used instead.
1. If the constructor is given in the child class, it has to call `super` with appropriate parameters inside it.

```

class Foo extends Bar {
    constructor(val: int) {
        super(val);
    }
}

```

## Constructor

_constructorDeclarationStmt_ := `constructor` `(` _formalParamList_ `)` `;`

_constructorDefinitionStmt_ := `constructor` `(` _formalParamList_ `)` _functionDefinitionBlock_

Only one constructor is allowed.

## Implementing interfaces

In declaration and/or definition, a class can be specified to be implementing certain interfaces using keyword `implements`:

```
import "./bar1";
import "./bar2";
import "./bar3";

export declare {
    class Foo implements bar1.Bar1, bar2.Bar2, bar3.Bar3 {

    }
}

class Foo implements bar1.Bar1, bar2.Bar2, bar3.Bar3 {

}
```

In declaration, the functions of interfaces do not need to be explicitly re-declared in the declaration block.

## interface-class

In many cases, the naming of interface and the class that implements the interface could be difficult. For example, I want to create an interface called `HomepageModelBuilder`, and the naming is very specific such that I cannot find another appropriate naming for the actual class. There are majorly two ways to work around this issue. One is to add an `I` in front of the naming of the interface. The other one is to append `Imp` to the naming of the actual class. IMO, both solutions are not ideal. Thus, here is what I propose:

When a declaration of a class is specified with keyword `interface-class`, the declaration creates one interface and one class with the same naming.

1. Other class can still implement this interface-class. When this happens, the implementation in this class is ignored.
1. When passing `interface-class` as a parameter, it is considered as an interface.
1. When creating `interface-class`, it is considered as a class.

For example:

in `foo.ex`:

```
export declare {
    interface-class Foo {
        doA();
        doB();
        doC();
    }
}

class Foo {
    doA() {
        print(5);
    }

    doB() {

    }

    doC() {

    }
}
```

in `bar.ex`:

```
export declare {
    class Bar implements Foo {
        // has doA(), doB(), and doC()
    }
}

class Bar {
    doA() {
        print(6);
    }

    doB() {

    }

    doC() {

    }
}

```

in `main.ex`:

```
import "./foo";
import "./bar";

func doSomething() {
    let foo = new Foo();
    let bar = new Bar();

    acceptsFoo(foo); // prints 5
    acceptsFoo (bar); // prints 6
}

func acceptsFoo(foo: Foo) {
    foo.doA();
}

func main() {
    doSomething();
}
```
