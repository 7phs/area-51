package detector

import (
	"log"
	"time"
)

// TODO: refactor to struct + methods
func parseCSV(delimiter byte, skipHeader, firstLine bool, prevIndex int, prevBuf, buf []byte, fn func(record DataRecord)) (bool, int) {
	var (
		prev  = 0
		index = 0
		c     byte
	)

	for index, c = range buf {
		if c != '\n' {
			// TODO: check maximum length of line
			continue
		}

		if firstLine {
			firstLine = false

			if skipHeader {
				prev = index + 1
				continue
			}
		}

		v := buf[prev:index]
		if prevBuf != nil && prevIndex >= 0 {
			i := copy(prevBuf[0:], prevBuf[prevIndex:])
			ln := len(v)
			copy(prevBuf[i:], v)

			v = prevBuf[:i+ln]
		}

		rec, err := parseDataRecord(delimiter, v)
		if err != nil {
			log.Println(time.Now(), "failed to parse csv line, skip it: ", err)
			continue
		}

		fn(rec)

		prevBuf = nil
		prev = index + 1
	}

	if buf[len(buf)-1] != '\n' {
		return firstLine, prev
	}

	if prev < len(buf)-2 {
		rec, err := parseDataRecord(delimiter, buf[prev:])
		switch {
		case err != nil:
			log.Println(time.Now(), "failed to parse csv line, skip it: ", err)
		default:
			fn(rec)
		}
	}

	return firstLine, -1
}

func parseCSVLine(delimiter byte, line []byte, fn func(index int, v []byte) error) (int, error) {
	prev := 0
	index := 0
	quote := false
	skip := false

	for i, c := range line {
		if skip {
			skip = false
			continue
		}

		switch c {
		case '"':
			quote = !quote
		case '\\':
			skip = true
			continue
		}
		if quote {
			continue
		}

		if c != delimiter {
			continue
		}

		if err := fn(index, line[prev:i]); err != nil {
			return index, err
		}

		prev = i + 1
		index++
	}

	if err := fn(index, line[prev:]); err != nil {
		return index, err
	}

	return index, nil
}
