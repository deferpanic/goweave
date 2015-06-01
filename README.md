# goweave
  Aspect Oriented Programming for Go

![Weave](http://i.imgur.com/JUUgIuv.png)

### TOC

  [What is AOP](https://github.com/deferpanic/goweave#what_is_aop)

  [Why](https://github.com/deferpanic/goweave#why)

  [Usage](https://github.com/deferpanic/goweave#usage)

  [Examples](https://github.com/deferpanic/goweave#examples)

  [FAQ](https://github.com/deferpanic/goweave#faq)

  [Help](https://github.com/deferpanic/goweave#help)

  [Todo](https://github.com/deferpanic/goweave#todo)

  [Roadmap](https://github.com/deferpanic/goweave#roadmap)

## What is AOP !??

  [Aspect oriented programming](http://docs.jboss.org/aop/1.1/aspect-framework/userguide/en/html/what.html)

  In short - this is a pre-processor that generates code defined by a
goweave specification file - it's not a proper grammar yet although it
probably should be.

  Existing Tools:
    * go fmt:
      This is actually used for around advice currently. It allows you to wrap
methods. Having said that - we wish to do more proper around advice than
simply re-writing the function declaration.

    * go fix:
      This is one hell of an awesome tool. I just think it's a little
too low-level for what we are wanting to do. Remember - one of the
solutions of this tool is to make things as trivial as possible to
insert new functionality.

    * go cover:
      This is used to provide code coverage and has similar properties
to what we want.

    * go generate:
      We are generating code but we are looking for more extensive code
generation.

### Why!??!

> "... which is our fulltime job, write a program to write a program"
  - rob pike

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

Automate all the things. All this is, is a powertool - I don't care 1% about
programming theory/philosophy.

That is the rationale behind this.

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

### Examples:

The example folder has some examples you might wish to look at.
Essentially we support aspect files w/in a project. If code exists in
your project we should be able to provide aspect coverage for it.

There are a few design decisions that need to be made to support across
projects && into stdlib. Stdlib probably won't come until we move to IR.

To try things out first try running `go build`. Then try running `goweave`.

#### Before Main

```go
aspect {
  pointcut: execute(main)
  advice: {
    before: {
      fmt.Println("before main")
    }
  }
}
```

#### Before Function
```go
aspect {
  pointcut: execute(beforeBob)
  advice: {
    before: {
      fmt.Println("before bob")
    }
  }
}
```

#### After Function
```
aspect {
  pointcut: execute(afterSally)
  advice: {
    after: {
      fmt.Println("after sally")
    }
  }
}
```

### Around Function
```
aspect {
  pointcut: execute(aroundTom)
  advice: {
    before: {
      fmt.Println("before tom")
    }
    after: {
      fmt.Println("after tom")
    }
  }
}
```

### Grammar:

  The 'grammar' if you can call it that is a total piece of shit right
now. It is a little bit of go, a little of json, etc. It is most definitely
not going to stay the same - it will be improved in the future.

  I apologize for giving you the forks to stab your collective eyes out.

  I think a good goal to have is to make it as proper go as possible and
then extend it maybe through comments.

  Suggestions/pull requests more than welcome!

### Definitions:

  I probably have not defined certain things properly here - open a pull
request if you find something off.

  * join point - places in your code you can apply behavior

  * pointcut - expression that details where to apply behavior

  We support both method && call pointcut primitive right now:
    -- __call__
    These happen before, after or wrap around calling a method. The code
is outside of the function.

    -- __execute__
    These happen before or after executing a method. The code is put
inside the method.

    call examples:
    ```go
      some.stuff()
    ```

    Code will be executed {before, around, after} this call.

    call before:
    ```go
      fmt.Println("before")
      some.stuff()
    ```

    call after:
    ```go
      some.stuff()
      fmt.Println("before")
    ```

    call around:
    ```
      somewrapper(some.stuff())
    ```

    execute examples:
    ```go
      func stuff() {
        fmt.Println("stuff")
      }
    ```

    execute before:
    ```go
      func stuff() {
        fmt.Println("before")
        fmt.Println("stuff")
      }
    ```

    execute after:
    ```go
      func stuff() {
        fmt.Println("stuff")
        fmt.Println("after")
      }
    ```

    pointcut examples:
    ```go
      pointcut: execute(beforeBob)
    ```

    explicit function name w/wildcard arguments:
    ```go
      pointcut: http.HandleFunc(d, s)
    ```

    explicit function declaration w/wildcard function name:
    ```go
      pointcut: execute(d(http.ResponseWriter, *http.Request))
    ```

  * advice - behavior to apply
    Behavior can be {before, after, around}. Around advice currently only works
with call pointcuts.

  * aspect - a .weave file that contains our behavior
    Right now we support multiple .weave projects for a project and they
will apply advice recursively over a project.

### Aspects:

  The programming theory department says that aspects are common features
that you use everywhere that don't really have anything at all to do with
your domain logic. Logging is a canonical example - most everything you
log does not really have anything to do with all the other stuff you
log.

  Similarly if you had a http controller that whenever you got a request
you would update a metric counter for that controller but you do this on
each api controller - that really has nothing at all to do with the
controller logic itself. The metric might simply be another aspect that
is commong everywhere.

  Once again someone might point out why don't you just make a method
and then wrap each call? The point here is that 1) we don't want to
modify code, 2) we might not *know* all the places that happens and
could easily leave something out, 3) we are eternally lazy and would
rather the computer do this for us.

### PointCut:

  Pointcuts in other languages such as java can commonly use annotations
    -- we currently don't support this as we want to be un-obtrusive as possible
    -- that is - we don't want to modify go source

  All pointcuts are currently defined only on functions. Struct field
members are definitely a future feature we could support although go
generate might do this acceptably already.

  Note: this 'grammar' if you can call it that sucks - expect it to
change "heavily".

  * explicit method name
    ```go
      call(blah)
    ```

    ```go
      execute(blah)
    ```

  TODO
  * partial match method name
    ```go
      call(b.*)
    ```

    ```go
      execute(b.*)
    ```

  * function declaration w/wildcard arguments
    ```go
      call(http.HandleFunc(d, s))
    ```

  * wildcard function name w/explicit arguments
      execute((w http.ResponseWriter, r *http.Request))
 
  * sub-pkg && method name
    ```go
      execute(pkg/blah)
    ```

  * sub-pkg && struct && method-name
    ```go
      execute(pkg/struct.b)
    ```

  # note - you have to have the AST for this
  * struct && method name
    ```go
      execute(struct.b)
    ```


### Advice:

  We currently support the following advice:

  * before
  * after
  * around

### Goals

* fast - obviously there will always be greater overhead than just
  running go build but we don't want this to be obscene - right now it's
  a little obscene

* correct - it goes w/out saying this is highly important to be as
  correct as possible w/our code generation

* no code modifications - my main use cases involve *not* modifying code
  so that is why we initially did not support annotations - I'm not
  opposed to adding these but that's not my intended goal

### FAQ

* Why not go generate?

* why not go fmt?

  We actually use go fmt code for the around conditions.

  http://research.swtch.com/gofmt

  ```go
  gofmt -r 'bytes.Compare(a, b) == 0 -> bytes.Equal(a, b)'
  ```

* Why not do everything via the AST?
  I think we want to move all the regexen to AST. This started out as a
POC and I wanted functionality first - correctness comes after.

* What about IR generation?
  This is probably the next step in the chain after converting most of
this to AST based processing.

* What about aspects on binary/closed-source?
  This is arguably one of the bigger benefits of AOP (at least for our
purposes) and it's definitely something we intend to support/code for in
the future.

  That's a long ways away but not off the radar/roadmap.

* Why wouldn't you just code this into your source?
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

* Are you all insane? This is go heresey!! Burn them at the stake!
  :) No, we are practioners of the "get-shit-done" philosophy. Ergo, we
don't care about philosophy of programming, nor do we care about other
armchair concerns. We only care about - how fast can I get this done?

  Our use cases usually entail us having to jump into brand new large
codebases and we want to send 'tracer bullets' out very very fast. This
style of programming allows us to do that.

  Lastly, you don't have to use this if you don't like it. To each their own.

### What You Should Know Before Using

This is *alpha* software - at best. It's more of an idea right now than anything
else.

* Expect the "grammars" {aspects, pointcuts} to change.

* This is currently *much* slower compared to native go build. Expect that to
  change but right now it's slow.

* Expect the build system to change soon. It's slow and crap.

* This *might* eat your cat - watch out.

### TODO - shortlist before opening up

  * remove the regex stuff

  * break apart large tests into units

### random thoughts

  * https://groups.google.com/forum/#!topic/golang-nuts/TiRX4HcdZMw
  * https://github.com/skelterjohn/gorf

  * gofmt code 
    - russ or robert?

  * AST re-writing -->
    http://golang.org/src/cmd/gofmt/rewrite.go
    https://github.com/ncw/gotemplate/compare/ast-arguments
    https://github.com/tsuna/gorewrite

  * split up file re-writing into something that we can test easily

  * partial function matching

  * ~/ap/main.go
    https://github.com/tmc/fix/blob/master/fix.go

  * go fix yo' shit
    - https://code.google.com/p/go/source/browse/src/cmd/fix/netdial.go?name=go1
    - https://code.google.com/p/go/source/browse/src/cmd/fix/osopen.go?name=go1
    - https://code.google.com/p/go/source/browse/src/cmd/fix/httpserver.go?name=go1

  * how much do we want to integrate ??
    https://www.godoc.org/github.com/tmc/fix
    (port of rsc's)


### TODO - a.k.a. - Known Suckiness

* make everything use ast

* make everything modify ast not raw text

* multi-pass parsing
  - this should ideally be a single pass
  - most of the regex/line scanning should be converted to AST node
    replacement

* import vendoring/re-writing

* better error handling
  - can do bail outs if parser doesn't emit correctly

* matching function declarations
  - with arguments
  - with return arguments
  - partial function matching

* scope - lol
  - for the regex && line-editing stuff this is completely naive - pulls
    pls

* relative path fix
  - relative paths are super hacky - pulls pls

* annotations??

* Faster

* Better Test coverage - lulz

* make it easy to share advice/aspects through a central site
  -- maybe start off w/just github?

### Help

  Want to help? Ideas for helping out:

    * test coverage
    * benchmark coverage
    * sample aspects - aspects should be shared - no need to re-invent
      the wheel here

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
