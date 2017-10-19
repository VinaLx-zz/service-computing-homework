package util

import (
	"crypto/md5"
	"fmt"
	"io"
	"os"
	"time"
)

// Overlapped returns true when two time interval overlapped
func Overlapped(s1, e1, s2, e2 time.Time) bool {
	return s1.Before(e2) && e1.After(s2)
}

// PrintfErr to stderr
func PrintfErr(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format, args...)
}

// PrintlnErr to stderr
func PrintlnErr(args ...interface{}) {
	fmt.Fprintln(os.Stderr, args...)
}

// PrettyHash returns the hex string version of md5 hash
func PrettyHash(text string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(text)))
}

type noopWriter struct{}

// NoopWriter returns the writer that do nothing
func NoopWriter() io.Writer {
	return noopWriter{}
}

func (w noopWriter) Write(p []byte) (n int, err error) {
	return len(p), nil
}

const ymdLayout = "2006-01-02"

// YMDParse parse the string in year-month-day format to Time
func YMDParse(str string) (time.Time, error) {
	return time.Parse(ymdLayout, str)
}

// YMDFormat returns the formatted time in year-month-day layout
func YMDFormat(t time.Time) string {
	return t.Format(ymdLayout)
}
