# goa
  Aspect Oriented Programming for Go

### TOC

  [Usage](https://github.com/deferpanic/goa#usage)

  [Examples](https://github.com/deferpanic/goa#examples)

  [What is AOP](https://github.com/deferpanic/goa#what_is_aop)

  [Why](https://github.com/deferpanic/goa#why)

  [FAQ](https://github.com/deferpanic/goa#faq)

  [Goals](https://github.com/deferpanic/goa#goals)

  [Help](https://github.com/deferpanic/goa#help)

  [Todo](https://github.com/deferpanic/goa#todo)

  [Roadmap](https://github.com/deferpanic/goa#roadmap)

### Usage:

Where you might use

```go
  go build
```

simply replace with

```go
  goa
```

### Use Cases
  * error detection && correction
    (ex: force logging of errors on any methods with this declaration)

  * data validation
    (ex: notate that this data was invalid but allow it to continue)

  * i18n
    (ex: translate this to german if accessed in an improper way)

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

### Examples:

The example folder has some examples you might wish to look at.
Essentially we support aspect files w/in a project. If code exists in
your project we should be able to provide aspect coverage for it.

There are a few design decisions that need to be made to support across
projects && into stdlib. Stdlib probably won't come until we move to IR.

To try things out first try running `go build`. Then try running `goa`.

#### Before Main

```go
aspect {
  pointcut: main
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
  pointcut: beforeBob
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
  pointcut: afterSally
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
  pointcut: aroundTom
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

  I think a good goal to have is to make it as proper go as possible. Suggestions welcome.

## What is AOP !??

  [Aspect oriented programming](http://docs.jboss.org/aop/1.1/aspect-framework/userguide/en/html/what.html)

  in short - we are a pre-processor

  tools existing:
    gofmt - rename api defs on the fly
    gofix - ???
    https://blog.golang.org/cover

### Definitions:

  * join point - places you can apply behavior
    -- method call
    ex:
    ```go
      some.stuff()
    ```

    Code will be executed {before, around, after} this call.

    before:
    ```go
      fmt.Println("before")
      some.stuff()
    ```

    after:
    ```go
      some.stuff()
      fmt.Println("before")
    ```

    around:
    ```
      somewrapper(some.stuff())
    ```

    -- method execution (what we are doing right now)
    ```go
      func stuff() {
        fmt.Println("stuff")
      }
    ```

    before:
    ```go
      func stuff() {
        fmt.Println("before")
        fmt.Println("stuff")
      }
    ```

    after:
    ```go
      func stuff() {
        fmt.Println("stuff")
        fmt.Println("after")
      }
    ```

  * pointcut - expression that details where to apply behavior
    -- right now we only explicitly match on function names

    ```go
      pointcut: beforeBob
    ```

  * advice - behavior to apply

  * aspect - a .goa file - file that contains our behavior

### Aspects:

  Aspects are common features that you use everywhere that don't really
have anything at all to do with your domain logic. If you have a user
interface that deals with updating passwords, setting preferences, etc.
logging might be done in the same way as you would log a dog.

  Similariy if you had a http controller that whenever you got a request
you would update a metric counter for that controller but you do this on
each api controller - that really has nothing at all to do with the
controller logic itself. The metric might simply be another aspect that
is commong everywhere.

### PointCut:

  Pointcuts in other languages such as java can commonly use annotations
    -- we currently don't support this as we want to be un-obtrusive as possible
    -- that is - we don't want to modify go source

  All pointcuts are currently defined in the same file. This is
definitely open to discussion on what is best though.

  All pointcuts are currently defined only on functions.

  There is no method overloading in go so currently the last thing in a
pointcut definition will be the method name (which can be a partial
match).

  Note: this 'grammar' if you can call it that sucks - expect it to
change "heavily".

  * explicit method name
    ```go
      call("blah")
    ```

    ```go
      execute("blah")
    ```

  * partial match method name
    ```go
      call(b.*)
    ```

    ```go
      execute(b.*)
    ```

  * function declaration
    ```go
      (w http.ResponseWriter, r *http.Request)
    ```

  * sub-pkg && method name
    ```go
      pkg/blah
    ```

  * sub-pkg && struct && method-name
    ```go
      pkg/struct.b
    ```

  # note - you have to have the AST for this
  * struct && method name
    ```go
      struct.b
    ```


### Advice:

  We currently support the following advice:

  * before
  * after
  * around

### What's up with the Name?
I was going to name this the flaming neckbeard in honor of those who
after seeing this code or hearing about it would have their respective
beards spontaneously combust into flame.

Instead I named it after Goa, India where I went to relax after
GopherCon India back in February and hacked out deprehend. I see it as
an extension of that work.

The name sucks - suggest a new one.

### Why!??!

"Which is our fulltime job, write a program to write a program"
  - rob pike

We came to go to get away from java!! What's the matter with you!?

I agree this concept can and has been abused in the past. However, being
able to do some of the things you can do with this is just way too
conveinent.

I'm definitely not a code purist - to me coding is a tool first and
foremost.

I simply wanted an easy way to attach common bits of code to large
existing codebases without lifting a finger. That is the rationale
behind this.

### Goals

* fast - obviously there will always be greater overhead than just
  running go build but we don't want this to be obscene

* correct - it goes w/out saying this is highly important to be as
  correct as possible w/our code generation

### FAQ

* why Not go generate?

  I don't intend for this codebase to live on regexen forever. It's more
of a POC while the business logic gets sorted out.

* why not go fmt?

  http://research.swtch.com/gofmt

  ```go
  gofmt -r 'bytes.Compare(a, b) == 0 -> bytes.Equal(a, b)'
  ```

* Why not AST?
  I think we want to move all the regexen to AST. This started out as a
POC and I wanted functionality first.

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
way.

  2) The original reason we did this was over at DeferPanic we had many
requests from people wanting to use our code to automatically insert
code in. For existing codebases this was a lot of work. After we made a
[tool](https://github.com/deferpanic/deprehend) that did this code
generation we had requests to make it non-obtrusive - that is - they
didn't want the code inside their codebase - just available to them at
runtime.

  3) I'd like the ability to turn on/off the behavior at will *and* not
have to re-code it for every project. I think this is where AOP really
shines.

* Are you all insane? This is go heresey!!
  :) No, we are practioners of the "get-shit-done" philosophy. Ergo, we
don't care about philosophy of programming, nor do we care about other
armchair concerns. We only care about - how fast can I get this done?

  Our use cases usually entail us having to jump into brand new large
codebases and we want to send 'tracer bullets' out very very fast. This
style of programming allows us to do that.

  Lastly, you don't have to use this if you don't like it. To each their own.

### What You Should Know Before Using

This is *alpha* software. It's more of an idea right now than anything
else.

* Expect the grammars {aspects, pointcuts} to change.

* This is currently *much* slower compared to native go build. Expect that to
  change but right now it's slow.

* Expect the build system to change soon.

* This *might* eat your cat - watch out.

### TODO - shortlist before opening up

  * better name

  * logo

  * need to be able to modify a goroutine
    (need that AST <3 again)
    (eg: panics.. inside goroutines)

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

### TODO

* import vendoring/re-writing

* inner vs. outer cutpoints
  - tha fuck?

* better error handling
  - can do bail outs if parser doesn't emit correctly
  - prob. err on gratitutious

* matching function declarations
  - with arguments
  - with return arguments

* scope - lol
  - this is currently completely stupid

* relative path fix - lol
  - relative paths are hacky

* convert all this exec stuff to native go if possible

* cross file
  -- test?

* annotations??

* Faster

* Better Tested - lulz

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
