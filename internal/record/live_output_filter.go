package record

import (
	"io"
	"strings"
)

type liveOutputFilter struct {
	tail string
}

func (f *liveOutputFilter) write(w io.Writer, chunk []byte, flush bool) {
	s := f.tail + string(chunk)
	f.tail = ""

	for {
		idx := strings.Index(s, exitCodeMarkerPrefix)
		if idx < 0 {
			break
		}
		// Write the text before marker can include prompt, here -> "$ ".
		if idx > 0 {
			_, _ = io.WriteString(w, s[:idx])
		}

		rest := s[idx:]
		nl := strings.IndexByte(rest, '\n')
		if nl < 0 {
			// Incomplete marker line, we'll keep it for next chunk.
			f.tail = rest
			return
		}
		// Drop full marker line.
		s = rest[nl+1:]
	}

	if flush {
		_, _ = io.WriteString(w, s)
		return
	}

	// Try to keep only a possible partial marker suffix for next chunk
	keep := longestMarkerPrefixSuffix(s)
	if len(s)-keep > 0 {
		_, _ = io.WriteString(w, s[:len(s)-keep])
	}
	f.tail = s[len(s)-keep:]
}

func longestMarkerPrefixSuffix(s string) int {
	max := len(exitCodeMarkerPrefix) - 1
	if max > len(s) {
		max = len(s)
	}
	for k := max; k > 0; k-- {
		if strings.HasSuffix(s, exitCodeMarkerPrefix[:k]) {
			return k
		}
	}
	return 0
}
