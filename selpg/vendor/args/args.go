package args

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

// Args is the struct returned by command line argument parser
type Args struct {
	Sources chan *ReadSrc
	Dest    io.Writer
	Filter  PageFilter
	Pager   Pager
}

// Get returns the pointer to the result of parsed args
func Get() *Args {
	if args == nil {
		args = parseArgs()
	}
	return args
}

// PageFilter receive the page number and the page content returns whether
// the page is wanted or not, and whether it should be the last page
type PageFilter func(int, []byte) (bool, bool)

// ReadSrc specifies the name and reader of the source
type ReadSrc struct {
	Reader io.Reader
	Name   string
	Next   chan bool
}

// Pager read string from the reader and combine them into pages
type Pager func(io.Reader) chan []byte

func pageFilter(start, end int) PageFilter {
	return func(n int, content []byte) (contained, isLast bool) {
		if n < start || n > end {
			return
		}
		contained = true
		if n == end {
			isLast = true
		}
		return
	}
}

func pageDelimiterPager(reader io.Reader) chan []byte {
	channel := make(chan []byte)
	coroutine := func() {
		b := bufio.NewReader(reader)
		for {
			s, err := b.ReadBytes(byte('\f'))
			s = append(s, byte('\n'))
			channel <- s
			if err != nil {
				if err == io.EOF {
					close(channel)
					break
				} else {
					panic(fmt.Sprintf(
						"Unexpected error when reading from reader: %s",
						err.Error()))
				}
			}
		}
	}
	go coroutine()
	return channel
}

func fixLinePager(n uint) Pager {
	return func(reader io.Reader) chan []byte {
		channel := make(chan []byte)
		coroutine := func() {
			lineScan := bufio.NewScanner(reader)
			for {
				end := false
				data := make([]byte, 0)
				for i := uint(0); i < n; i++ {
					if lineScan.Scan() {
						line := lineScan.Bytes()
						line = append(line, byte('\n'))
						data = append(data, line...)
					} else {
						end = true
						break
					}
				}
				channel <- data
				if end {
					close(channel)
					break
				}
			}
		}
		go coroutine()
		return channel
	}
}

func argumentError(reason string, usage bool) {
	fmt.Fprintln(os.Stderr, reason)
	if usage {
		flag.Usage()
	}
	os.Exit(2)
}

var args *Args

func pageFilterOrExit(start, end int) PageFilter {
	if start < 0 || end < 0 {
		argumentError("both start and end argument should be positive", true)
	}
	if start > end {
		argumentError("start should be less than end", true)
	}
	return pageFilter(start, end)
}

func getSources() chan *ReadSrc {
	channel := make(chan *ReadSrc)
	next := make(chan bool)

	coroutine := func() {
		if len(flag.Args()) == 0 {
			channel <- &ReadSrc{Reader: os.Stdout}
			close(channel)
			return
		}
		for _, f := range flag.Args() {
			file, err := os.Open(f)
			if err != nil {
				fmt.Fprintln(os.Stderr, err.Error())
				continue
			}
			// equivalent: name := len(flag.Args()) == 1 ? "" : f
			var name string
			if len(flag.Args()) != 1 {
				name = f
			}
			channel <- &ReadSrc{Name: name, Reader: file, Next: next}
			b := <-next // waiting user handing over control
			file.Close()
			if !b {
				break
			}
		}
		close(channel)
	}
	go coroutine()
	return channel
}

func getDestOrExit(dest string) io.Writer {
	dest = strings.TrimSpace(dest)
	if dest == "" {
		return os.Stdout
	}
	spaces, err := regexp.Compile("\\s+")
	if err != nil {
		panic("what spaces regex compile fail??")
	}
	cmds := spaces.Split(dest, -1)
	cmd := exec.Command(cmds[0], cmds[1:]...)
	writer, err := cmd.StdinPipe()
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Start()
	if err != nil {
		argumentError(err.Error(), false)
	}
	return writer
}

func getPagerOrExit(lines uint, pageSep bool) Pager {
	if pageSep {
		return pageDelimiterPager
	}
	return fixLinePager(lines)
}

func parseArgs() *Args {
	start := flag.Int("s", -1, "the starting index of page number (mandatory)")
	end := flag.Int("e", -1, "the ending index of page number (mandatory)")
	pageSeperator := flag.Bool("f", false, "use \\f to seperate pages")
	pageLines := flag.Uint("l", 72, "number of lines for each page")
	dest := flag.String(
		"d", "", "a destination command that receives the output,"+
			" if it's not specified, output is printed to stdout")

	flag.Parse()
	filter := pageFilterOrExit(*start, *end)
	sources := getSources()
	destWriter := getDestOrExit(*dest)
	pager := getPagerOrExit(*pageLines, *pageSeperator)
	return &Args{
		Sources: sources,
		Filter:  filter,
		Dest:    destWriter,
		Pager:   pager,
	}
}
