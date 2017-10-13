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
	Sources []io.Reader
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

func errorExit(reason string) {
	fmt.Fprintln(os.Stderr, reason)
	flag.Usage()
	os.Exit(2)
}

var args *Args

func pageFilterOrExit(start, end int) PageFilter {
	if start < 0 || end < 0 {
		errorExit("both start and end argument should be positive")
	}
	if start > end {
		errorExit("start should be less than end")
	}
	return pageFilter(start, end)
}

func getSourcesOrExit() []io.Reader {
	// assert flag.Parsed
	if len(flag.Args()) == 0 {
		return []io.Reader{os.Stdin}
	}
	sources := make([]io.Reader, 0, len(flag.Args()))
	for _, f := range flag.Args() {
		file, err := os.Open(f)
		if err != nil {
			errorExit(err.Error())
		}
		sources = append(sources, file)
	}
	return sources
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
	if err != nil {
		errorExit(err.Error())
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
	sources := getSourcesOrExit()
	destWriter := getDestOrExit(*dest)
	pager := getPagerOrExit(*pageLines, *pageSeperator)
	return &Args{
		Sources: sources,
		Filter:  filter,
		Dest:    destWriter,
		Pager:   pager,
	}
}
