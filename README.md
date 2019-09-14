Micro
---

Micro is a Go implementation of the language covered in the book [Crafting A Compiler With C](https://www.pearson.com/us/higher-education/program/Fischer-Crafting-a-Compiler-with-C/PGM263627.html).

It has the following features/limitations:

- The only supported type is an integer, which is represented as a literal containing a string of digits.
- Identifiers can be assigned, which have a max of 32 characters. Identifiers must start with a letter, and may be followed by letters, digits, and underscores.
- Comments start with `--` and run until the end of the line.
- Assignments are done in the format `identifier := expression;`.
- Expressions are constructed with a combination of identifiers, literals, and the operators `+` and `-`, parenthesis are also allowed.
- I/O can be done with `read` and `write`.
- The program must begin with `begin` and end with `end`.
- Each statement must end with a semicolon.
- Tokens may not run across line bounderies.
