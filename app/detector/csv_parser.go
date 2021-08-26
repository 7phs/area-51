package detector

import (
	"log"
	"time"

	data_stream "github.com/7phs/area-51/app/data-stream"
)

// TODO: refactor to struct + methods
func parseCSV(delimiter byte, skipHeader, firstLine bool, prevIndex int, prevBuffer, buffer data_stream.Buffer, fn func(record DataRecord)) (bool, int) {
	buf := buffer.Bytes()

	if len(buf) == 0 {
		return firstLine, 0
	}

	var (
		prevBuf []byte
		prev    = 0
	)

	if prevBuffer != nil {
		prevBuf = prevBuffer.Bytes()
	}

	for index, c := range buf {
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
			// TODO: check if it makes sense to create temporary buffer, but not reuse pre-buf
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

	if (buffer.Command() == data_stream.NewData || buffer.Command() == data_stream.Data) && buf[len(buf)-1] != '\n' {
		return firstLine, prev
	}

	if (buffer.Command() == data_stream.EOF && prev < len(buf)) || prev < len(buf)-2 {
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

func parseCSVLineN(delimiter byte, line []byte, N int, fn func(index int, v []byte) error) (int, error) {
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

		if N >= 0 && index+1 >= N {
			break
		}
	}

	if err := fn(index, line[prev:]); err != nil {
		return index, err
	}

	return index, nil
}
