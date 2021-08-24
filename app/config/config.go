package config

import (
	"encoding/json"
	"flag"
	"path/filepath"
)

type Config struct {
	ReferenceFile string `json:"reference"`
	RawDataFile   string `json:"raw"`
	OutputDir     string `json:"output"`
}

func (c Config) Validate() error {
	// Reference file
	if c.ReferenceFile == "" {
		return ErrEmptyParam("reference")
	}
	if err := isFileValid(c.ReferenceFile); err != nil {
		return err
	}
	if err := isDirValid(filepath.Dir(c.ReferenceFile)); err != nil {
		return err
	}

	// Raw data file
	if c.RawDataFile == "" {
		return ErrEmptyParam("raw")
	}
	if err := isFileValid(c.RawDataFile); err != nil {
		return err
	}
	if err := isDirValid(filepath.Dir(c.RawDataFile)); err != nil {
		return err
	}

	if c.ReferenceFile == c.RawDataFile {
		return ErrEqualPath(c.ReferenceFile, c.RawDataFile)
	}

	// Output path
	if c.OutputDir == "" {
		return ErrEmptyParam("output")
	}
	if err := isDirValid(c.OutputDir); err != nil {
		return err
	}

	return nil
}

func (c Config) Dump() string {
	b, _ := json.Marshal(&c)
	return string(b)
}

func Parse() Config {
	config := Config{}

	flag.StringVar(&config.ReferenceFile, "reference", "", "a path to a reference file")
	flag.StringVar(&config.RawDataFile, "raw", "", "a path to a raw data file")
	flag.StringVar(&config.OutputDir, "output", "", "a directory of output files")

	flag.Parse()

	return config
}
