goweave
  Aspect Oriented Programming for Go

[![wercker status](https://app.wercker.com/status/7d7912cec649bc9763736051f64da3fa/s/master "wercker status")](https://app.wercker.com/project/bykey/7d7912cec649bc9763736051f64da3fa)

[![GoDoc](https://godoc.org/github.com/deferpanic/goweave?status.svg)](https://godoc.org/github.com/deferpanic/goweave)


![Weave](http://i.imgur.com/JUUgIuv.png)

### WARNING - Major Hackage

  This really isn't meant to be used by anyone yet - definitely not in a
production environment - you have been warned!

  Many 'design decisions' were not decisions at all - they were simply
the "most simplest thing that could work". Lots of work left to do.

### TOC

  [What is AOP](https://github.com/deferpanic/goweave#what_is_aop)

  [Why](https://github.com/deferpanic/goweave#why)

  [Usage](https://github.com/deferpanic/goweave#usage)

  [Examples](https://github.com/deferpanic/loom#examples)

  [Loom](https://github.com/deferpanic/loom)

  [FAQ](https://github.com/deferpanic/goweave#faq)

  [Info](https://github.com/deferpanic/goweave#info)

  [Reserved Keywords](https://github.com/deferpanic/goweave#reserved-keywords)

  [Differences](https://github.com/deferpanic/goweave#differences)

  [Performance](https://github.com/deferpanic/goweave#performance)

  [Tests](https://github.com/deferpanic/goweave#tests)

  [Help](https://github.com/deferpanic/goweave#help)

  [Todo](https://github.com/deferpanic/goweave#todo)

  [Roadmap](https://github.com/deferpanic/goweave#roadmap)

## What is AOP !??

  [Aspect oriented programming](http://docs.jboss.org/aop/1.1/aspect-framework/userguide/en/html/what.html)

  In short - this is a pre-processor that generates code defined by a
goweave specification file - it's not a proper grammar yet although it
probably should be.

### Why!??!

> "... which is our fulltime job, write a program to write a program"
> rob pike

The critics yell red in the face - "We came to go to get away from enterprise java!!
What the hell is the matter with you!?"

I agree this concept can and has been abused in the past.

However, I'm definitely not a code purist - to me coding is a tool first and
foremost.

I simply wanted an easy way to attach common bits of code to large
existing codebases without having to hack it in each time to each
different codebase. I also needed the ability to do this on large
projects without modifying the source.

Sure I could write method wrappers, I could change code, etc. but I
don't want to be constantly writing/re-writing/re-moving/inserting tons
of code just to check things out.

Automate all the things. For me, this is simply a powertool for deep introspection of go programs.

That is the rationale behind this.

### Existing Tools:

#### [go fmt](https://golang.org/pkg/fmt/)
This is actually used for around advice currently. It allows you to wrap
methods. Having said that - we wish to do more proper around advice than
simply re-writing the function declaration.

  ex:

  ```go
  gofmt -r 'bytes.Compare(a, b) == 0 -> bytes.Equal(a, b)'
  ```

#### [go fix](https://golang.org/cmd/fix/):
This is one hell of an awesome tool. I just think it's a little
too low-level for what we are wanting to do. Remember - one of the
solutions of this tool is to make things as trivial as possible to
insert new functionality.

#### [go cover](https://godoc.org/golang.org/x/tools/cmd/cover):
This is used to provide code coverage and has similar properties
to what we want.

#### [go generate](http://blog.golang.org/generate):
We are generating code but we are looking for more extensive code
generation.

### Usage:

Where you might use

```go
  go build
```

simply replace with

```go
  goweave
```

Note: depending on how you build your go this may or may not work -
patches/pulls to make this more generic for everyone are definitely
welcome.

### Use Cases
  * error detection && correction
    (ex: force logging of errors on any methods with this declaration)

  * data validation
    (ex: notate that this data was invalid but allow it to continue)

  * i18n
    (ex: translate this to esperanto if you can't detect the language)

  * security
    (ex: authenticate this user in each http request)

  * caching
    (ex: cache these variables when exposed in a http request)

  * logging
    (ex: log when this group of users accesses these methods)

  * monitoring
    (ex: ensure that if this channel closes we always alert joebob)

  * metrics
    (ex: cnt the number of times this function is called)

  * tracing
    (ex: print out the value of this variable in a pkg)

  * dealing with legacy code
    (ex: overriding a method/API w/minimum of work)

  * static validation
    (ex: force closing a file if we detect that it hasn't been closed)

### Grammar:

  The aspect 'grammar' if you can call it that is a total piece right
now. It is a little bit of go, a little of json, etc. It is most definitely
not going to stay the same - it will be improved in the future.

  I apologize for giving you the forks to stab your collective eyes out.

  I think a good goal to have is to make it as proper go as possible and
then extend it maybe through comments.

  Suggestions/pull requests/discussions more than welcome!

### Definitions:

  I probably have not defined certain things properly here - open a pull
request if you find something off.

### Join Points

  Places in your code you can apply behavior.

### Aspects:

  A .weave file that contains our behavior.

  Right now we support multiple .weave projects for a project and they
will apply advice recursively over a project.

  The programming theory department says that aspects are common features
that you use everywhere that don't really have anything at all to do with
your domain logic. Logging is a canonical example - most everything you
log does not really have anything to do with all the other stuff you
log.

  Similarly if you had a http controller that whenever you got a request
you would update a metric counter for that controller but you do this on
each api controller - that really has nothing at all to do with the
controller logic itself. The metric might simply be another aspect that
is common everywhere.

  Once again someone might point out why don't you just make a method
and then wrap each call? The point here is that 1) we don't want to
modify code, 2) we might not *know* all the places that happens and
could easily leave something out, 3) we are eternally lazy and would
rather the computer do this for us.

### PointCut:

  An expression that details where to apply behavior.

  Pointcuts in other languages such as java can commonly use annotations
    -- we currently don't support this as we want to be un-obtrusive as possible
    -- that is - we don't want to modify go source

  We support {call, execute, within} pointcut primitives right now:

__call__:

    These happen before, after or wrap around calling a method. The code
is outside of the function.

__execute__:

    These happen before or after executing a method. The code is put
inside the method.

__within__:

    These happen for *every* statement within a function body
declaration.


  All pointcuts are currently defined only on functions. Struct field
members are definitely a future feature we could support although go
generate might do this acceptably already.

  Note: this 'grammar' if you can call it is nowhere close to 'ok' - expect it to
change "heavily".

#### explicit method name
```go
  call(blah())
```

```go
  execute(blah())
```

#### partial match method name - TODO

```go
  call(b.*)
```

```go
  execute(b.*)
```

#### function declaration w/wildcard arguments
```go
  call(http.HandleFunc(d, s))
```

#### wildcard function name w/explicit arguments
      execute((w http.ResponseWriter, r *http.Request))

#### doesn't work yet - TODO

sub-pkg && method name
```go
  execute(pkg/blah())
```

sub-pkg && struct && method-name
```go
  execute(pkg/struct.b())
```

#### struct && method name - TODO
```go
  execute(struct.b())
```

### Advice:

  Behavior to apply.

  * before
  * after
  * around

  Around advice currently only works with call pointcuts.

  We currently support the following advice:

#### call examples:
    ```go
      some.stuff()
    ```

    Code will be executed {before, around, after} this call.

    __call before:__
    ```go
      fmt.Println("before")
      some.stuff()
    ```

    __call after:__
    ```go
      some.stuff()
      fmt.Println("before")
    ```

    __call around:__
    ```
      somewrapper(some.stuff())
    ```

#### execute examples:
    ```go
      func stuff() {
        fmt.Println("stuff")
      }
    ```

    __execute before:__
    ```go
      func stuff() {
        fmt.Println("before")
        fmt.Println("stuff")
      }
    ```

    __execute after:__
    ```go
      func stuff() {
        fmt.Println("stuff")
        fmt.Println("after")
      }
    ```

#### within examples:
    ```go
      func blah() {
        slowCall()
        fastCall()
      }
    ```

    __within before:__
    ```go
      func blah() {
        beforeEach()
        slowCall()
        beforeEach()
        fastcall()
      }
    ```

### Goals

* FAST - obviously there will always be greater overhead than just
  running go build but we don't want this to be obscene - right now it's
  a little obscene

* CORRECT - it goes w/out saying this is highly important to be as
  correct as possible w/our code generation

* NO CODE MODIFICATIONS - my main use cases involve *not* modifying code
  so that is why we initially did not support annotations - I'm not
  opposed to adding these but that's not my intended goal

* create tooling around AO development for go

* maybe move towards compiler extension?

### FAQ

#### Why not do everything via the AST?
  We are moving all the regexen to AST modifications. This started out as a
POC and I wanted functionality first - correctness comes after.

#### Why not modify the AST instead of re-writing the source each time an
  aspect is applied?
  We also want to do this - once again - it was a POC and
correctness/speed comes later. Pull requests welcome.

#### What about IR generation?

#### What about aspects on binary/closed-source?
  This is arguably one of the bigger benefits of AOP (at least for our
purposes) and it's definitely something we intend to support/code for in
the future.

  That's a long ways away but not off the radar/roadmap.

#### Why wouldn't you just code this into your source?
  A couple of reasons.

  1) If you are going to do something like development tracing (eg:
sprinkle some fmt.Println everywhere) you don't want that in your
production code. It's much better to apply it when necessary in your
binary, fix the problem and go - there is no need to code it in, then
hack it back out (and potentially miss some). It's *much* cleaner this
way and it's *much* faster.

  2) The original reason we did this was over at DeferPanic we had many
requests from people wanting to use our code to automatically insert
code in. For existing codebases this was a lot of work. After we made a
[tool](https://github.com/deferpanic/deprehend) that did this code
generation we had requests to make it non-obtrusive - that is - they
didn't want the code inside their codebase - just available to them at
runtime.

  3) I'd like the ability to turn on/off the behavior at will *and* not
have to re-code it for every project. I think this is where this really
shines.

#### Are you all insane? This is go heresey!! Burn them at the stake!

  :) We are practioners of the "get-shit-done" philosophy. Whether this
is considered good or bad practice is not quite a concern for our
usecases.

  We only care about - how fast can I get this done?

  Our use cases usually entail us having to jump into brand new large
codebases and we want to send 'tracer bullets' out very very fast. This
style of programming allows us to do that.

### Info

* builds are currently made in ~/go/src/_weave - I don't think this is
  ideal so am open to further suggestions

* builds are *really* *really* slow right now - there is a lot of file
  I/O we shouldn't be doing - many files we read/write multiple times -
this goes hand in hand w/the text processing - most of this should just
be moved to the AST processing - plenty of cruft laying around as well
that needs to be refactored

### Reserved Keywords

  Right now the only time you'll run into reserved keywords are in the
experimental 'within' and 'get'/'set' advice section although there is an intention to
support a set of keywords that one can use in their aspects.

  * mName
    If you use 'mName' within your within advice it will translate to a
string representation of the joinpoint found by your within pointcut.

  * mAvars
    this is a list of abstract variables used within get/set advice

  We'd appreciate help from the community formulating a more formal
approach for this. Namespacing, the set of keywords supported, etc.

### Differences

  There are many differences between this and other AOP implementations
such as AspectJ.

  Firstly, since go has no vm executing bytecode so we don't weave at the bytecode
level. Currently we weave at source code level at build time.

  Secondly, we don't currently support annotations. I'm open to this in
the future but it wasn't a primary usecase for me so I didn't add them.

### Performance

  Is pretty bad right now. Lots of needless reading/writing of files. Part of
this was cause it started out with regex/line parsing and slowly moved
towards AST manipulating.

  Once most of the file rewriting is moved towards AST modification the
performance should imrpove dramatically but right now it really sucks
and is going to make you cry.

  Lots of further work in this area to do. If you want to help - pulls
are definitely welcome.

  For a point of reference - on a well known web application it takes 12
seconds to build versus 1.59 seconds with just go build.

  We'd very much like help from the community in refactoring some of
these performance problems.

### Tests

  The tests are very brittle right now and are more functional than unit
based. Lots of work here to do. Once most of the file reading/writing
stuff is replaced with AST replacement transformations the tests should be much more specific not to mention must faster.

  I really don't like the fact that the tests are the way they are right
now but just need to ensure certain things work until we refactor it.

  We would appreciate any help from the community in refactoring the
code so the tests aren't big blocks of text - no need for that.

### What You Should Know Before Using

This is *alpha* software - at best. It's more of an idea right now than anything
else.

* Expect the "grammars" {aspects, pointcuts} to change.

* This is currently *much* slower compared to native go build. Expect that to
  change but right now it's slow.

* Expect the build system to change soon. It's slow and crap. Getting
  the latency down is very much an immediate goal.

* This *might* eat your cat - watch out.

### TODO
#### a.k.a. - Known Suckiness

* add ability to add global function advice to pkg

* make everything use ast - no raw txt

* multi-pass parsing
  - this should ideally be a single pass
  - most of the regex/line scanning should be converted to AST node
    replacement

* add support for matching function declaration w/returns

* import vendoring/re-writing
  -- very open to different ways of approaching this - it kinda sucks
right now

* better error handling
  - can do bail outs if parser doesn't emit correctly

* matching function declarations
  - with arguments
  - with return arguments
  - partial function matching

* scope
  - for the regex && line-editing stuff this is completely naive - pulls
    pls

* annotations??

* Faster

* Better Test coverage

### Pointcut Todo
  * create a more proper language

  * match against method receivers
  * match against return arguments
  * match stdlib
  * match 3rd pkgs

### Aspect Todo

  * ability to add functions to global namespace
    -- maybe just need some tests here

### Help

  Want to help? Ideas for helping out:

  * test coverage

  * benchmark coverage

  * documentation

  * sample aspects - aspects [should be shared on the loom](https://github.com/deferpanic/loom)
    - no need to re-invent the wheel

  Need helping visualizing what you are looking at? Check out http://goast.yuroyoro.net/

### Roadmap

#### Grammars
  * better aspect grammar
  * better pointcut grammar

#### Parsing/Speed
  * move from regexen to AST
  * move from AST to IR

#### Extending
  * add support for 3rd party pkgs
  * add support for stdlib

#### Interface Pointcuts
  * be able to define on interface fields
  * be able to define on methods that satisfy interfaces
