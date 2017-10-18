package model

import (
	"os"
)

// AgendaDir is the directory contains the agenda persistent files
func AgendaDir() string {
	home, present := os.LookupEnv("HOME")
	if !present {
		home = "."
	}
	return home + "/.agenda/"
}

// EnsureAgendaDir checks the presentation of agenda home directory,
// if not, create it
func EnsureAgendaDir() {
	os.Mkdir(AgendaDir(), 0755)
}

// LoginFile contains the path to the login state file
func LoginFile() string {
	return AgendaDir() + "login"
}

// UserFile path
func UserFile() string {
	return AgendaDir() + "user"
}

// MeetingFile path
func MeetingFile() string {
	return AgendaDir() + "meeting"
}
