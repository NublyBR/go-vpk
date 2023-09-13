package vpk

import "testing"

func TestFilenameInvalidUnix(t *testing.T) {
	bad := [][]string{
		// Path / File / Extension
		{" ", " ", " "},         // ""
		{"/", " ", " "},         // "/"
		{"/etc", "passwd", " "}, // "/etc/passwd"
		{"..", "hello", "txt"},  // "../hello.txt"
		{"dir", "\x00", ""},     // "dir/\0"
		{"", "hello", "t/x/t"},  // "hello.t/x/t"
		{"", "he/llo", "txt"},   // "hel/lo.txt"
	}

	for _, b := range bad {
		e := &entry{path: b[0], file: b[1], ext: b[2]}

		if e.FilenameSafeUnix() {
			t.Errorf("%q should not be a valid filename", e.Filename())
		}
	}
}

func TestFilenameInvalidWindows(t *testing.T) {
	bad := [][]string{
		// Path / File / Extension
		{" ", " ", " "},         // ""
		{"/", " ", " "},         // "/"
		{"c:", "hello", "txt"},  // "c:/hello.txt"
		{"..", "hello", "txt"},  // "../hello.txt"
		{"dir", "\x00", " "},    // "dir/\0"
		{"dir", "\x1f", " "},    // "dir/\x1f"
		{" ", "con", " "},       // "con"
		{" ", "com1", " "},      // "com1"
		{" ", "?", " "},         // "?"
		{" ", "dot.", " "},      // "dot."
		{" ", "space ", " "},    // "space "
		{" ", "hello", "t/x/t"}, // "hello.t/x/t"
		{" ", "he/llo", "txt"},  // "hell/o.txt"
	}

	for _, b := range bad {
		e := &entry{path: b[0], file: b[1], ext: b[2]}

		if e.FilenameSafeWindows() {
			t.Errorf("%q should not be a valid filename", e.Filename())
		}
	}
}

func TestFilenameValidUnix(t *testing.T) {
	good := [][]string{
		// Path / File / Extension
		{" ", "hello", "txt"},             // "hello.txt"
		{"one/two/three", "hello", "txt"}, // "one/two/three/hello.txt"
		{" ", "con", " "},                 // "con"
		{" ", "com1", " "},                // "com1"
	}

	for _, b := range good {
		e := &entry{path: b[0], file: b[1], ext: b[2]}

		if !e.FilenameSafeUnix() {
			t.Errorf("%q should be a valid filename", e.Filename())
		}
	}
}

func TestFilenameValidWindows(t *testing.T) {
	good := [][]string{
		// Path / File / Extension
		{" ", "hello", "txt"},             // "hello.txt"
		{"one/two/three", "hello", "txt"}, // "one/two/three/hello.txt"
		{" ", "con", "txt"},               // "con.txt"
	}

	for _, b := range good {
		e := &entry{path: b[0], file: b[1], ext: b[2]}

		if !e.FilenameSafeWindows() {
			t.Errorf("%q should be a valid filename", e.Filename())
		}
	}
}
