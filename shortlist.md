### TODO - shortlist before opening up

* get rid of examples once regex stuff is completely moved over..

* remove the regex stuff
  -- think we can do this by just commenting out the pointcut match &&
the http only pointcut match

  -- looked at this - the go-routine test is the only thing failing if
we do this and it is fialing cause the ast parser will prob. bitch..

  -- maybe punt this to an integration style test for now ??

* probably change the call/around advice in the http example
  -- right now it's just replacing - we want to show true around usage

* get go-routine coverage working differently (w/out regex?)

* break apart large tests into units
