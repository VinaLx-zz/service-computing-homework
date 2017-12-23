# Goroutine Test

## Task

From [Mr. PML's blog](http://blog.csdn.net/pmlpml/article/details/78850661)

Simulate the scenario in [this blog](https://jersey.github.io/documentation/latest/rx-client.html#d0e5556), both the synchronized and asynchronizd one.

## Demo

usage:

```
$ go get github.com/vinalx/service-computing-homework/goroutine
$ goroutine -h
Usage of ./goroutine:
  -s, --scenario string
    	choose a scenario, 1 for the synchronized scenario, 2 for the asynchronized one
```

And as the blog says, two version of the program approximately run in 5400ms and 730ms

```
$ time ./goroutine -s 1
./goroutine -s 1  0.00s user 0.01s system 0% cpu 5.424 total
$ time ./goroutine -s 2
./goroutine -s 2  0.00s user 0.00s system 0% cpu 0.740 total
```

## A more general solution

Well, "message driven" concurrency model want to increase the performance and reactivity of the program by parallelizing actions which have no data dependency. We can work out a general library that formulate this kind of concurrent scenario but it's not possible in golang (without reflection, but reflection solve this problem while introducing another, which is static type safety).

And good reference would be the [`promise`](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/Promise) in javascript, but it's not the peek of the abstraction. (There is still a more general concept which formulate the control flow in language theory.).

But we probably need to modify the promise library to support both synchonized and asynchronized actions, and I would  strongly recommand the reading of [parallel moand](https://hackage.haskell.org/package/monad-par) **if you have functional programming background**.

So that's it.