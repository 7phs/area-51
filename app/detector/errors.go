package detector

import "fmt"

var (
	_ error = CSVLessFields{}
	_ error = CSVOutOfIndex{}
	_ error = FeaturesInvalidFormat("")
	_ error = FeaturesLessIndex{}
	_ error = FeaturesOutOfIndex{}
	_ error = Float64Convert{}
)

type CSVLessFields struct {
	index int
}

func ErrCSVLessFields(index int) CSVLessFields {
	return CSVLessFields{
		index: index,
	}
}

func (e CSVLessFields) Error() string {
	return fmt.Sprintf("csv line has less fields than expected: %d", e.index)
}

type CSVOutOfIndex struct {
	index int
}

func ErrCSVOutOfIndex(index int) CSVOutOfIndex {
	return CSVOutOfIndex{
		index: index,
	}
}

func (e CSVOutOfIndex) Error() string {
	return fmt.Sprintf("csv line has more fields than expected: %d", e.index)
}

type FeaturesInvalidFormat string

func ErrFeaturesInvalidFormat() FeaturesInvalidFormat {
	return FeaturesInvalidFormat("")
}

func (e FeaturesInvalidFormat) Error() string {
	return "feature invalid format"
}

type FeaturesLessIndex struct {
	index int
}

func ErrFeaturesLessIndex(index int) FeaturesLessIndex {
	return FeaturesLessIndex{
		index: index,
	}
}

func (e FeaturesLessIndex) Error() string {
	return fmt.Sprintf("feature line has less fields than expected: %d", e.index)
}

type FeaturesOutOfIndex struct {
	index int
}

func ErrFeaturesOutOfIndex(index int) FeaturesLessIndex {
	return FeaturesLessIndex{
		index: index,
	}
}

func (e FeaturesOutOfIndex) Error() string {
	return fmt.Sprintf("feature line has more fields than expected: %d", e.index)
}

type Float64Convert struct {
	err error
}

func ErrFloat64Convert(err error) Float64Convert {
	return Float64Convert{
		err: err,
	}
}

func (e Float64Convert) Error() string {
	return fmt.Sprintf("failed to convert feature value to float64: '%s'", e.err)
}

func (e Float64Convert) Unwrap() error {
	return e.err
}
