# HTTP Server hello world

## Usage

```
$ go get -u -v https://github.com/VinaLx/service-computing-homework/myhttp
$ $GOPATH/bin/http -h

Usage of myhttp:
  -host string
    	host of the server listening to (default "localhost")
  -log string
    	file output of log (default "stdout")
  -port uint
    	port of the server listening to (default 8080)
```

## Test

Run server:
```
$ $GOPATH/bin/http
hello world serving at localhost:8080
```

Test with `curl`
```
# at another shell

$ curl -v "127.0.0.1:8080"
* Rebuilt URL to: 127.0.0.1:8080/
*   Trying 127.0.0.1...
* TCP_NODELAY set
* Connected to 127.0.0.1 (127.0.0.1) port 8080 (#0)
> GET / HTTP/1.1
> Host: 127.0.0.1:8080
> User-Agent: curl/7.56.1
> Accept: */*
>
< HTTP/1.1 200 OK
< Date: Thu, 02 Nov 2017 16:33:03 GMT
< Content-Length: 12
< Content-Type: text/plain; charset=utf-8
<
hello world
* Connection #0 to host 127.0.0.1 left intact
```

Test with `ab`

```
$ ab -n 10000 -c 100 http://127.0.0.1:8080/
This is ApacheBench, Version 2.3 <$Revision: 1796539 $>
Copyright 1996 Adam Twiss, Zeus Technology Ltd, http://www.zeustech.net/
Licensed to The Apache Software Foundation, http://www.apache.org/

Benchmarking 127.0.0.1 (be patient)
Completed 1000 requests
Completed 2000 requests
Completed 3000 requests
Completed 4000 requests
Completed 5000 requests
Completed 6000 requests
Completed 7000 requests
Completed 8000 requests
Completed 9000 requests
Completed 10000 requests
Finished 10000 requests


Server Software:
Server Hostname:        127.0.0.1
Server Port:            8080

Document Path:          /
Document Length:        12 bytes

Concurrency Level:      100
Time taken for tests:   0.748 seconds
Complete requests:      10000
Failed requests:        0
Total transferred:      1290000 bytes
HTML transferred:       120000 bytes
Requests per second:    13373.02 [#/sec] (mean)
Time per request:       7.478 [ms] (mean)
Time per request:       0.075 [ms] (mean, across all concurrent requests)
Transfer rate:          1684.69 [Kbytes/sec] received

Connection Times (ms)
              min  mean[+/-sd] median   max
Connect:        1    3   1.1      3       9
Processing:     1    4   1.4      4      13
Waiting:        0    3   1.2      3      13
Total:          3    7   1.8      7      16

Percentage of the requests served within a certain time (ms)
  50%      7
  66%      8
  75%      8
  80%      9
  90%     10
  95%     11
  98%     12
  99%     13
 100%     16 (longest request)
```

## Implementation Details

Wellllll, I didn't use any external framework, the go standard http package is ok for the job since basically there isn't any requirement for this homework.

But there's still something worth mentioning.

Although the hello world task is extreeeeeemely simple, I decouple the logic of program into four parts:

1. Command line argument parsing, into an struct `Args` (Implemented by `flag`)
2. Server logic, which runs the server with given configuration, host, port, log etc. (Implemented by `http.Server`)
3. Multiplexer logic, which dispatch the request by different path (Implement by `http.ServeMux`)
4. the main program which connect things together.

And that's it.

## Last

Err?