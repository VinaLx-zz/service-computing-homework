# Database Service

## Dependencies

[negroni](https://github.com/urfave/negroni)

[gorilla/mux](https://github.com/gorilla/mux)

[pflag](https://github.com/ogier/pflag)

[gorm](https://github.com/jinzhu/gorm)

[go-sqlite3](https://github.com/mattn/go-sqlite3)

## Task

Implement a http based db service (storing "User"), using both native "SQL" and ORM (Object Relational Mapping) library. Use `ab` to compare their performance.

## Idea

I give the same `DAO` interface to the upper layer router, and give two implementation to the DAO interface, `database/ormdao` and `database/sqldao` respectively. And use command line argument to select which one to use.

And in the benchmark, both method starts with a empty database and empty table, comparing the time to insert 10000 users into the database. The concurrency is not one of the test, since the test mainly compares the performance of two database service implementation schemes.

## Result

The testing idea is given above. And here's the result.

Native SQL implementation:
```shell
$ ./db
[negroni] listening on localhost:8080

# switch to another terminal
$ ab -n 10000 -pform -q -T 'application/x-www-form-urlencoded' http://127.0.0.1:8080/adduser
This is ApacheBench, Version 2.3 <$Revision: 1796539 $>
Copyright 1996 Adam Twiss, Zeus Technology Ltd, http://www.zeustech.net/
Licensed to The Apache Software Foundation, http://www.apache.org/

Benchmarking 127.0.0.1 (be patient).....done


Server Software:
Server Hostname:        127.0.0.1
Server Port:            8080

Document Path:          /adduser
Document Length:        38 bytes

Concurrency Level:      1
Time taken for tests:   11.769 seconds
Complete requests:      10000
Failed requests:        0
Total transferred:      1550000 bytes
Total body sent:        1990000
HTML transferred:       380000 bytes
Requests per second:    849.70 [#/sec] (mean)
Time per request:       1.177 [ms] (mean)
Time per request:       1.177 [ms] (mean, across all concurrent requests)
Transfer rate:          128.62 [Kbytes/sec] received
                        165.13 kb/s sent
                        293.74 kb/s total

Connection Times (ms)
              min  mean[+/-sd] median   max
Connect:        0    0   0.0      0       1
Processing:     1    1   0.4      1      21
Waiting:        1    1   0.4      1      21
Total:          1    1   0.4      1      21

Percentage of the requests served within a certain time (ms)
  50%      1
  66%      1
  75%      1
  80%      1
  90%      1
  95%      2
  98%      2
  99%      3
 100%     21 (longest request)
```

ORM Library(GORM) implementation:
```
$ ./db --orm
[negroni] listening on localhost:8080

# switch to another terminal
$ ab -n 10000 -pform -q -T 'application/x-www-form-urlencoded' http://127.0.0.1:8080/adduser
This is ApacheBench, Version 2.3 <$Revision: 1796539 $>
Copyright 1996 Adam Twiss, Zeus Technology Ltd, http://www.zeustech.net/
Licensed to The Apache Software Foundation, http://www.apache.org/

Benchmarking 127.0.0.1 (be patient).....done


Server Software:
Server Hostname:        127.0.0.1
Server Port:            8080

Document Path:          /adduser
Document Length:        38 bytes

Concurrency Level:      1
Time taken for tests:   12.795 seconds
Complete requests:      10000
Failed requests:        0
Total transferred:      1550000 bytes
Total body sent:        1990000
HTML transferred:       380000 bytes
Requests per second:    781.58 [#/sec] (mean)
Time per request:       1.279 [ms] (mean)
Time per request:       1.279 [ms] (mean, across all concurrent requests)
Transfer rate:          118.31 [Kbytes/sec] received
                        151.89 kb/s sent
                        270.20 kb/s total

Connection Times (ms)
              min  mean[+/-sd] median   max
Connect:        0    0   0.1      0       4
Processing:     1    1   1.5      1     145
Waiting:        0    1   1.5      1     145
Total:          1    1   1.5      1     145

Percentage of the requests served within a certain time (ms)
  50%      1
  66%      1
  75%      1
  80%      1
  90%      1
  95%      2
  98%      3
  99%      3
 100%    145 (longest request)
```

## Explanation

So as we can see, in this trivial example, the ORM library did have some negative effect on the performance (from request per second and transfer rate), but from my point of view, it's not toooo much. Since using a well implemented ORM interfaces would be a lot easier and maintainable than writing plain SQL, which has bad scalability, bad readability, and maybe go wrong elsewhere.

## Other Functionalities of the User DB Service

Since no emphasize of this part is mentioned in the requirement, so I only implement those which appeared in [the sample code of PML](http://blog.csdn.net/pmlpml/article/details/78602290).

### adding user

```shell
$ curl -d "username=userA&password=somepassword" "localhost:8080/adduser"
{"OK":true,"Data":"add user success"}
```

### querying user by id

```shell
$ curl -d "userid=0" "localhost:8080/getuser"
{"OK":true,"Data":{"ID":0,"Username":"userA","Password":"somepassword","SignUpDate":"2017-11-24T20:41:03.702438+08:00"}}
```
### getting all users


```shell
$ curl "localhost:8080/getallusers"
{"OK":true,"Data":[{"ID":0,"Username":"userA","Password":"somepassword","SignUpDate":"2017-11-24T20:41:03.702438+08:00"}]}
```

## Last

The implementation and abstraction is just ordinary in this code, not really anything worth to mention.

Only after I finish this implementation did I realize I didn't settle the concurrency issue, I was thinking about the web framework helping me handle all the race condition here until I realize the database connection object should be carefully initialized, and ... fine that's it.