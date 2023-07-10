# Implementation of the [Lox Programming language](https://craftinginterpreters.com/the-lox-language.html)

[WIP]
---

## Quick Start
### Install latest release binary (only available on macos and linux):
Run the script below in your terminal to install `glox` in your system inside `/usr/local/bin/glox`. After installing the binary, you can run the command `glox` in your terminal and it should start the REPL.
```sh
curl -L -s https://api.github.com/repos/silverhairs/crafting-interpreters/releases/latest \
| grep "browser_download_url.*glox-$(uname -s | tr '[:upper:]' '[:lower:]' | sed 's/darwin/macos/')-$(uname -m)" \
| cut -d '"' -f 4 \
| wget -qi - \
&& chmod +x glox-$(uname -s | tr '[:upper:]' '[:lower:]' | sed 's/darwin/macos/')-$(uname -m) \
&& sudo mv glox-$(uname -s | tr '[:upper:]' '[:lower:]' | sed 's/darwin/macos/')-$(uname -m) /usr/local/bin/glox
```

### Run code
Cloning the repository and runing `main.go` should start the REPL. `main.go` is located in `crafting-interpreters/glox`. You need to have Golang installed in your system for this option.

```sh
go run main.go
```

## Grammar

Production rules:

```txt
    expression -> literal
                | unary
                | binary
                | grouping
                | ternary ;

    literal    -> NUMBER | STRING | boolean | "nil" ;
    unary      -> ( "-" | "!" ) expression ;
    grouping   -> "(" expression ")" ;
    binary     -> expression operator expression ;
    operator   -> "==" | "!=" | "<" | ">" | "<=" | ">="
                | "+" | "-" | "*" | "/" ;
    ternary -> expression "?" expression ":" expression;
```
