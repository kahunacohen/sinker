# sinker

A modular utility written in golang for syncing files to various backends written as an experiment to understand Golang.

## Background

I did this to explore a new generation of languages such as golang, rust, swift etc. In general these bring several things to the table:

* compile-time errors
* Strong(ish) type system
* native complilation with no run-time
* more predictible run-time error handling
* baked-in concurrent facilities

### Golang

#### Pros
* Enough type safety, but not too much
* Simple enough to write quick command-line tools (good balance)
* Compiles to portable, native code
* Fast
* Batteries included
* Major tools written in it (docker, kubernetes, influx etc.
* Radical simplicity of syntax
  * Easy to learn
  * Easy to read
  * You don't shoot yourself in the foot
* Concurrency as first class citizen
* High demand, high salary

#### Cons
* Error handling can be clumsy
* Conservative (e.g. still no generics)
* handling null pointers and errors is not enforced compile-time.
* Learning it doesn't teach you much new.
