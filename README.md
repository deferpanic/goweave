# goa
Aspect Oriented Programming for Go


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
  * data validation
  * i18n
  * security
  * caching
  * logging
  * monitoring
  * metrics
  * tracing

### Examples:

The coffee folder has some examples you might wish to look at.
Essentially we support aspect files w/in a project. If code exists in
your project we should be able to provide aspect coverage for it.

There are a few design decisions that need to be made to support across
projects && into stdlib. Stdlib probably won't come until we move to IR.

#### Before Main

```go
advice execution("func main()") : before {
  fmt.Println("before main")
}
```

#### Before Function
```go
advice execution("func beforeBob()"): before {
  fmt.Println("before bob")
}
```

#### After Function
```
advice execution("func afterSally()"): after {
  fmt.Println("after sally")
}
```

### Around Function
```
advice execution("func aroundTom()"): around {
  fmt.Println("before tom")
  goaProceed()
  fmt.Println("after tom")
}
```

The goaProceed() fake function is stripped and denotes the mark between
before/after code.

### Grammar:

  I know - just like the name - the grammar sucks right now. It shall be
improved in the future.

## What is AOP !??

  [Aspect oriented programming](http://docs.jboss.org/aop/1.1/aspect-framework/userguide/en/html/what.html)

### Aspects:

### Cutpoints:

### Advice:

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
We came to go to get away from java!! I agree this concept can and has
been abused in the past. However, being able to do some of the things
you can do with this is just way too conveinent.

I'm definitely not a code purist - to me coding is a tool first and
foremost.


### FAQ

* why Not go generate?

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

* Did DeferPanic Just Jump the Shark?
  :) No, we are practioners of the "get-shit-done" philosophy. Ergo, we
don't care about philosophy of programming, nor do we care about other
armchair constraints. We only care about - how fast can I get this done?

  Our use cases usually entail us having to jump into brand new large
codebases and we want to send 'tracer bullets' out very very fast. This
practice allows us to do that.

### Future

* arguments ??
  -- no argument support yet - want to write a test/patch?

* scope is completely fucked atm
  - choices are to track track scope or to ast or something else..

* return statements
  -- explicit return
  -- no return

* cross file
* cross package w/in project
* cross package
* affects stdlib

* Better advice
* More seamlessness
* Access to stdlib
* Faster
* Better Tested - lulz
