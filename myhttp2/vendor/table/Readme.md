# HTTP Server Hello World 2

## Usage

This is not an app, it's merely a demonstration of completion of homework.
If you really insist on checking it yourself, download the whole folder and place the executable under the `myhttp2` folder then start it.

## Tasks

Implement following features in a web server.

>- Static file service support
>- Simple web api access
>- Form and table rendering
>- Prompt for unknown paths

## Demonstrations

With the server with default argument running, running following curl command should be able to see similar outputs as follow:

### Static File Service Support

```shell
$ curl "localhost:8080/static/index.html"
<!DOCTYPE html>
<html lang="en">

<head>
  <title>Static File</title>
  <meta charset="UTF-8">
</head>

<body>
  <h1>This is a boring static file</h1>
</body>

</html>
```

#### Explanation

The default static folder and path prefix is "static", the index.html is under myhttp2/static, and a the response of `GET` request onto the `static/index.html` path should contains the content of `index.html`

### Simple web api access

```shell
$ curl "localhost:8080/testjs"
{"Status":true,"Message":"get ok"}
$ curl -X POST "localhost:8080/testjs"
{"Status":true,"Message":"post ok"}
```

#### Explanation

This is a simple api that echoes the request method. Since no other details are mentioned in the task of homework, maybe this is enough for demonstration.

### Form and table rendering

```shell
$ curl -d "foo=bar&baz=gak" -X POST "localhost:8080/formtable"
<!DOCTYPE html>
<html lang="en">

<head>
  <title>The Table</title>
  <meta charset="UTF-8">
  <link rel="stylesheet" href="https://maxcdn.bootstrapcdn.com/bootstrap/4.0.0-alpha.6/css/bootstrap.min.css" integrity="sha384-rwoIResjU2yc3z8GV/NPeZWAv56rSmLldC3R/AZzGRnGxQQKnKkoFVhFQhNUwEyJ" crossorigin="anonymous">
</head>

<body>
  <table class="table">
    <thead>
      <tr>
        <th>Key</th>
        <th>Value</th>
      </tr>
    </thead>
    <tbody>

        <tr>
            <td>baz</td>
            <td>gak</td>
        </tr>

        <tr>
            <td>foo</td>
            <td>bar</td>
        </tr>

    </tbody>
  </table>

  <script src="https://code.jquery.com/jquery-3.1.1.slim.min.js" integrity="sha384-A7FZj7v+d/sdmMqp/nOQwliLvUsJfDHW+k9Omg/a/EheAdgtzNs3hpfag6Ed950n" crossorigin="anonymous"></script>
  <script src="https://cdnjs.cloudflare.com/ajax/libs/tether/1.4.0/js/tether.min.js" integrity="sha384-DztdAPBWPRXSA/3eYEEUWrWCy7G5KFbe8fFjk5JAIxUYHKkDx6Qin1DkWx51bBrb" crossorigin="anonymous"></script>
  <script src="https://maxcdn.bootstrapcdn.com/bootstrap/4.0.0-alpha.6/js/bootstrap.min.js" integrity="sha384-vBWWzlZJ8ea9aCX4pEW3rVHjgjt7zpkNpZk+02D9phzyeVkE+jo0ieGizqPLForn" crossorigin="anonymous"></script>
</body>

</html>
```

#### Explanation

I pass the whole content of the form to the template engine, and the template generates a table that contains the keys and values of the form. Although I used `POST` in this example, sending data with `GET` is ok too. I include `bootstrap` in the html code. So you should have a clearer view when using the browser, and attatch the arguments to the original path.

### Prompt for unknown paths

```
$ curl -v "localhost:8080/someunknownpaths"
*   Trying ::1...
* TCP_NODELAY set
* Connection failed
* connect to ::1 port 8080 failed: Connection refused
*   Trying 127.0.0.1...
* TCP_NODELAY set
* Connected to localhost (127.0.0.1) port 8080 (#0)
> GET /someunknownpaths HTTP/1.1
> Host: localhost:8080
> User-Agent: curl/7.56.1
> Accept: */*
>
< HTTP/1.1 500 Internal Server Error
< Date: Thu, 16 Nov 2017 12:44:16 GMT
< Content-Length: 77
< Content-Type: text/plain; charset=utf-8
<
PANIC: Error: unknown path: /someunknownpaths, maybe it's under development?
* Connection #0 to host localhost left intact
```

And except the example paths for homework, all other paths should returns `500 Internal Server Error` as request.

## Implementation Details

Well, this homework is merely stacking all APIs of `gorilla`,`negroni` and `net/http`. So not much to say.
