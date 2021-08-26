package detector

import (
	"bytes"
	"io"
	"strconv"

	data_stream "github.com/7phs/area-51/app/data-stream"
)

const (
	dataRecordFieldCount = 3
	featuresCount        = 3
	featuresDelimiter    = ','
)

var (
	newLine = []byte(`
`)
)

type DataRecord struct {
	line []byte

	Key         []byte
	FeaturesF64 [featuresCount]float64
	Command     data_stream.Command
}

func parseDataRecord(delimiter byte, line []byte) (DataRecord, error) {
	record := DataRecord{
		line: line,
	}

	index, err := parseCSVLineN(delimiter, line, 3, record.assignValue)
	if err != nil {
		return record, err
	}
	if index < dataRecordFieldCount-1 {
		return record, ErrCSVLessFields(index)
	}

	return record, err
}

func (d *DataRecord) IsCommand() bool {
	return d.Command != 0
}

func (d *DataRecord) assignValue(index int, v []byte) error {
	switch index {
	case 0:

	case 1:
		if err := d.parseFeatures(v); err != nil {
			return err
		}

	case 2:
		d.Key = v

	default:
		return ErrCSVOutOfIndex(index)
	}

	return nil
}

func (d *DataRecord) parseFeatures(v []byte) error {
	if len(v) < 4 {
		return ErrFeaturesLessIndex(0)
	}
	if v[0] != '"' || v[1] != '[' || v[len(v)-2] != ']' || v[len(v)-1] != '"' {
		return ErrFeaturesInvalidFormat()
	}

	index, err := parseCSVLineN(featuresDelimiter, v[2:len(v)-2], -1, d.assignFeaturesF64)
	if err != nil {
		return err
	}
	if index < featuresCount-1 {
		return ErrFeaturesLessIndex(index)
	}

	return nil
}

func (d *DataRecord) assignFeaturesF64(index int, v []byte) error {
	if index >= featuresCount {
		return ErrFeaturesOutOfIndex(index)
	}

	var err error

	d.FeaturesF64[index], err = strconv.ParseFloat(string(bytes.TrimSpace(v)), 64)
	if err != nil {
		return ErrFloat64Convert(err)
	}

	return nil
}

func (d *DataRecord) Serialize(w io.Writer) {
	_, _ = w.Write(d.line)
	_, _ = w.Write(newLine)
}
