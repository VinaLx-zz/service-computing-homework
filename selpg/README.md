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

## TODO

Option `-d` seems don't work currently.

## Examples

`TODO`

## Implementation Details

`TODO`