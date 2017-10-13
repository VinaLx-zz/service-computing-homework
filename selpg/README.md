# selpg

[selpg](https://www.ibm.com/developerworks/cn/linux/shell/clutil/index.html) is a self-defined pagewise `cat`

## Build

```shell
go get -v -u github.com/VinaLx/service-computing-homework/selpg
```

## Usage

```shell
$GOPATH/bin/selpg
```

It should output:

```shell
Usage of selpg:
  -d string
    	a destination command that receives the output, if it's not specified, output is printed to stdout
  -e int
    	the ending index of page number (mandatory) (default -1)
  -f	use \f to seperate pages
  -l uint
    	number of lines for each page (default 72)
  -s int
    	the starting index of page number (mandatory) (default -1)
```

## Examples

Assume there are two files, `file1` and `file2`
```shell
$ cat file1
line 1 of file one
line 2 of file one
line 3 of file one
line 4 of file one
line 5 of file one
line 6 of file one
line 7 of file one
line 8 of file one
line 9 of file one
line 10 of file one
$ cat file2
line 1 of file two
line 2 of file two
line 3 of file two
line 4 of file two
line 5 of file two
line 6 of file two
line 7 of file two
line 8 of file two
line 9 of file two
line 10 of file two
```

Read from standard input
```shell
$ ./selpg -s 1 -e 2 -l 3 < file1
line 4 of file one
line 5 of file one
line 6 of file one
line 7 of file one
line 8 of file one
line 9 of file one
```

Read from file1
```shell
$ selpg -s 1 -e 2 -l 3 file1
line 4 of file one
line 5 of file one
line 6 of file one
line 7 of file one
line 8 of file one
line 9 of file one
```

Read from file1 and file2
```shell
$ selpg -s 1 -e 1 -l 3 file1 file2
file1:
line 4 of file one
line 5 of file one
line 6 of file one
file2:
line 4 of file two
line 5 of file two
line 6 of file two
```

Invalid inputs:
```
$ ./selpg -s 3 -e 2 -l 3
start should be less than end
Usage of ./selpg:
...
$ ./selpg -s -3 -e -2 -l 3
both start and end argument should be positive
Usage of ./selpg:
...
$ ./selpg -s 1 -e 2 -l 3 "no such file"
open no such file: no such file or directory
$
```

## Implementation Details

### `flag`

Standard library `flag` is used to parse command line arguments. There are some limitations of `flag`.

First it doesn't seems to support "mutual exclusive" options and generate error message accordingly, so after all arguments are parsed, we must manually validate the arguments we received.

Second, we must provide a default value for all options, so that we can't actually express the idea that some option is unpresent. Take `selpg` as example, I have to give "start" and "end" an default value, in this case -1, since empty value 0 is a valid option here. So that I can't know whether is the user who give the value -1 or it's just the unset default value. It can be done by `flag.Var` and use `nil` as an indicator, but it would add too much complication which the problem doesn't deserve.

### Modularization

Basically the `selpg` is doing following things:
```
for each input file 'f'
  split the content of 'f' into many 'p' with index 'i'
  determine whether we want the page
  if so, print 'p' to some 'destination'
```

And the parameters here is fully determined by the command line arguments, so the main logic of program and the parameters parsing could be simply decoupled. The "actual arguments" expressed by code could be:

```go
type Args struct {
	Sources chan *ReadSrc
	Dest    io.Writer
	Filter  PageFilter
	Pager   Pager
}

type ReadSrc struct {
	Reader io.Reader
	Name   string
	Next   chan bool
}

type PageFilter func(int, []byte) (bool, bool)

type Pager func(io.Reader) chan []byte
```

The `chan` here should be read as "a generator of". We would talk about that in the following section.

### IO

`go` has a bunch of io utilities to use. `io.Reader` and `io.Writer` do a great job in abstracting out the IO implementation details and can be passed around with no harm.

A slight problem is that if the underlying implementation is a opened file, we should guarentee that these files are closed if we run a server or things alike. We should use `io.Closer` to add the ability to close the file later in those cases. But since here the program is merely a cli tool. We can rely on the operating system to close a small amount of unclosed file when the program exit. But since number of files specified as the command line arguments could possibly be very large, so those files should be handle carefully.

### Optimization with coroutine

Although the overall struture of program is simple, there are some subtleties here.

First is we shouldn't begin to pick our desired pages after reading all contents of file since the file size could be very large or even infinite. So coroutine (or generator) is suitable here to generate pages as needed while decoupling the logic of page seperation with the main logic of program.

Another problem is that if there are more than one files to read, we shouldn't open all of them in advance since there's a limit of the maximal amount of file descriptors that can be opened by single process. So what we want to do here is to open a file, read it if no error, close it, and move to the next. To decoupling this logic from the main program, coroutine again comes into rescue.

So this is the reason I use the two `chan`s above.

## Last

Well, I did have some fun on this.