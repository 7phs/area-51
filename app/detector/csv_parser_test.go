package detector

import (
	data_stream "github.com/7phs/area-51/app/data-stream"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseCSV(t *testing.T) {
	buf := []byte(`id,feature,city,sport,size
100000,"[0.13103569, 0.046841323000000004, 0.19260901]",lisbon,scuba diving,S
100001,"[0.45330983, 0.34414744, 0.38385785]",vilnius,volleyball,M
100002,"[0.6264719999999999, 0.65428007, 0.09229988]",madrid,pool,XXL
`)
	expected := []DataRecord{
		{
			line: []byte(`100000,"[0.13103569, 0.046841323000000004, 0.19260901]",lisbon,scuba diving,S`),
			Key:  []byte("lisbon,scuba diving,S"),
			FeaturesF64: [3]float64{
				0.13103569,
				0.046841323000000004,
				0.19260901,
			},
		},
		{
			line: []byte(`100001,"[0.45330983, 0.34414744, 0.38385785]",vilnius,volleyball,M`),
			Key:  []byte("vilnius,volleyball,M"),
			FeaturesF64: [3]float64{
				0.45330983,
				0.34414744,
				0.38385785,
			},
		},
		{
			line: []byte(`100002,"[0.6264719999999999, 0.65428007, 0.09229988]",madrid,pool,XXL`),
			Key:  []byte("madrid,pool,XXL"),
			FeaturesF64: [3]float64{
				0.6264719999999999,
				0.65428007,
				0.09229988,
			},
		},
	}

	actual := make([]DataRecord, 0, 3)

	actualFirstLine, prevIndex := parseCSV(',', true, true, 0, nil, data_stream.DataBuffer(buf), func(record DataRecord) {
		actual = append(actual, record)
	})

	assert.False(t, actualFirstLine)
	assert.Equal(t, -1, prevIndex)
	assert.Equal(t, expected, actual)
}

func TestParseCSVEOF(t *testing.T) {
	buf := []byte(`id,feature,city,sport,size
100000,"[0.13103569, 0.046841323000000004, 0.19260901]",lisbon,scuba diving,S
100001,"[0.45330983, 0.34414744, 0.38385785]",vilnius,volleyball,M
100002,"[0.6264719999999999, 0.65428007, 0.09229988]",madrid,pool,XXL`)
	expected := []DataRecord{
		{
			line: []byte(`100000,"[0.13103569, 0.046841323000000004, 0.19260901]",lisbon,scuba diving,S`),
			Key:  []byte("lisbon,scuba diving,S"),
			FeaturesF64: [3]float64{
				0.13103569,
				0.046841323000000004,
				0.19260901,
			},
		},
		{
			line: []byte(`100001,"[0.45330983, 0.34414744, 0.38385785]",vilnius,volleyball,M`),
			Key:  []byte("vilnius,volleyball,M"),
			FeaturesF64: [3]float64{
				0.45330983,
				0.34414744,
				0.38385785,
			},
		},
		{
			line: []byte(`100002,"[0.6264719999999999, 0.65428007, 0.09229988]",madrid,pool,XXL`),
			Key:  []byte("madrid,pool,XXL"),
			FeaturesF64: [3]float64{
				0.6264719999999999,
				0.65428007,
				0.09229988,
			},
		},
	}

	actual := make([]DataRecord, 0, 3)

	actualFirstLine, prevIndex := parseCSV(',', true, true, 0, nil, data_stream.EOFBuffer(buf), func(record DataRecord) {
		actual = append(actual, record)
	})

	assert.False(t, actualFirstLine)
	assert.Equal(t, -1, prevIndex)
	assert.Equal(t, expected, actual)
}

func TestParseCSVWithPreBuf(t *testing.T) {
	preIndex := 78
	preBuf := []byte(`100000,"[0.13103569, 0.046841323000000004, 0.19260901]",lisbon,scuba diving,S
100002,"[0.6264719999999999, 0.65428007`)
	buf := []byte(`, 0.09229988]",madrid,pool,XXL
100003,"[0.45330983, 0.34414744, 0.38385785]",vilnius,volleyball,M
`)

	expected := []DataRecord{
		{
			line: []byte(`100002,"[0.6264719999999999, 0.65428007, 0.09229988]",madrid,pool,XXL`),
			Key:  []byte("madrid,pool,XXL"),
			FeaturesF64: [3]float64{
				0.6264719999999999,
				0.65428007,
				0.09229988,
			},
		},
		{
			line: []byte(`100003,"[0.45330983, 0.34414744, 0.38385785]",vilnius,volleyball,M`),
			Key:  []byte("vilnius,volleyball,M"),
			FeaturesF64: [3]float64{
				0.45330983,
				0.34414744,
				0.38385785,
			},
		},
	}

	actual := make([]DataRecord, 0, 3)

	actualFirstLine, prevIndex := parseCSV(',', true, false, preIndex, data_stream.DataBuffer(preBuf), data_stream.DataBuffer(buf), func(record DataRecord) {
		actual = append(actual, record)
	})

	assert.False(t, actualFirstLine)
	assert.Equal(t, -1, prevIndex)
	assert.Equal(t, expected, actual)
}
