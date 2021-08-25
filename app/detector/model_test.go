package detector

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseDataRecord(t *testing.T) {
	line := []byte(`100000,"[0.13103569, 0.046841323000000004, 0.19260901]",lisbon,scuba diving,S`)
	expected := DataRecord{
		Id:       []byte("100000"),
		City:     []byte("lisbon"),
		Sport:    []byte("scuba diving"),
		Size:     []byte("S"),
		Features: []byte(`"[0.13103569, 0.046841323000000004, 0.19260901]"`),
		FeaturesF64: [3]float64{
			0.13103569,
			0.046841323000000004,
			0.19260901,
		},
	}

	actual, err := parseDataRecord(',', line)

	require.NoError(t, err)

	assert.Equal(t, expected, actual)
}

func TestParseDataRecordInvalidFeaturesLine(t *testing.T) {
	line := []byte(`100000,"[0.13103569, 0.046841323000000004, ase0.19260901]",lisbon,scuba diving,S`)

	_, err := parseDataRecord(',', line)

	require.Error(t, err)
}
