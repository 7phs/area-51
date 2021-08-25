package detector

import "strconv"

type DataRecord struct {
	Id          []byte
	City        []byte
	Sport       []byte
	Size        []byte
	Features    []byte
	FeaturesF64 [3]float64
}

func parseCSV(delimiter byte, line []byte, fn func(delimiter byte, index int, v []byte)) {
	prev := 0
	index := 0
	quote := false
	skip := false

	for i, c := range line {
		// TODO: careful parsing + convert to iterator
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

		fn(delimiter, index, line[prev:i])

		prev = i + 1
		index++
	}

	fn(delimiter, index, line[prev:])

	// TODO: catch error of format. Counter less than expected
}

func parseDataRecord(delimiter byte, line []byte) DataRecord {
	record := DataRecord{}

	parseCSV(delimiter, line, record.assignValue)

	return record
}

func (d *DataRecord) assignValue(delimiter byte, index int, v []byte) {
	switch index {
	case 0:
		d.Id = v

	case 1:
		d.Features = v
		parseCSV(delimiter, v, d.assignFeaturesF64)

	case 2:
		d.City = v
	case 3:
		d.Sport = v
	case 4:
		d.Size = v
	default:
		// TODO: catch error of format
	}
}

func (d *DataRecord) assignFeaturesF64(_ byte, index int, v []byte) {
	if index >= len(d.FeaturesF64) {
		// TODO: catch error of format
		return
	}

	// TODO: catch error
	d.FeaturesF64[index], _ = strconv.ParseFloat(string(v), 64)
}
