package systemd

import "strings"

// Inspired by https://golang.org/src/net/url/url.go

func shouldEscape(pos int, c byte) bool {
	if 'A' <= c && c <= 'Z' || 'a' <= c && c <= 'z' || '0' <= c && c <= '9' {
		return false
	}
	if c == '_' {
		return false
	}
	if c == '.' {
		if pos > 0 {
			return false
		}
	}
	if c == ' ' {
		return false
	}
	if c == '/' {
		return false
	}

	return true
}

// From http://www.freedesktop.org/software/systemd/man/systemd.unit.html
// Some unit names reflect paths existing in the file system namespace.
// Example: a device unit dev-sda.device refers to a device with the device
// node /dev/sda in the file system namespace. If this applies, a special way
// to escape the path name is used, so that the result is usable as part of a
// filename. Basically, given a path, "/" is replaced by "-" and all other
// characters which are not ASCII alphanumerics are replaced by C-style "\x2d"
// escapes (except that "_" is never replaced and "." is only replaced when it
// would be the first character in the escaped path). The root directory "/"
// is encoded as single dash, while otherwise the initial and ending "/" are
// removed from all paths during transformation. This escaping is reversible.
// Properly escaped paths can be generated using the systemd-escape(1) command.
func Escape(s string) string {
	s = strings.TrimRight(s, " /")
	hexCount, slashCount := 0, 0
	for i := 0; i < len(s); i++ {
		switch c := s[i]; {
		case c == '/':
			slashCount++
		case shouldEscape(i, c):
			hexCount++
		}
	}

	if hexCount+slashCount == 0 {
		return s
	}

	t := make([]byte, len(s)+3*hexCount)
	j := 0
	for i := 0; i < len(s); i++ {
		switch c := s[i]; {
		case c == '/':
			t[j] = '-'
			j++
			break
		case shouldEscape(j, c):
			t[j] = '\\'
			t[j+1] = 'x'
			t[j+2] = "0123456789abcdef"[c>>4]
			t[j+3] = "0123456789abcdef"[c&15]
			j += 4
			break
		default:
			t[j] = s[i]
			j++
		}
	}
	return string(t)
}
