# Dataset anomalies :flying_saucer: detector

`area-51` is a tool for detecting analytics in a set of data using [Z-test].
A statistics preferences of reference data uses for scoring raw data.
Source of data is files in CSV format.

## How to use

`area-51` is CLI tool to process data stored in files.

Reference data and raw data should store in two different files.

Raw data splits by `Z-score` and put into two files:
* Clear data is stored to `clean.csv`;
* Record contains data greater than expected deviation is store to `anomalies.csv`.

The tool launches, processes existing files and listens changes of reference and raw files till a user stops it.
It processes changes in data (appending, replacing of file, etc.) in real time.

Log messages inform a user that exising data is completely processed.
A user stop the tool at proper time, when files will not change, or by another reason.

### Run

There are three CLI options to configure the tool:
* `--reference` - path to reference file;
* `--raw` - path to raw file;
* `--output` - path to directory for output file.

There are several options to run a tool.

#### Run as a docker image 

Docker image of `area-51` is stored in Docker hub. Use this command to start:

```bash
docker run -v /test-data:/mnt/test-data 7phs/area-51 --reference /mnt/test-data/in/ref.csv --raw /mnt/test-data/in/raw.csv --output /mnt/test-data/output/
```

#### Run with go tool 

`area-51` requires Go 1.17+ to build it.

```bash
git clone git@github.com:7phs/area-51.git
cd ./area-51
go run ./cmd/server --reference /test-data/in/ref.csv --raw /test-data/in/raw.csv --output /test-data/output/
```

### Makefile

## Description of the project

### Statistics

### Engineering

## TODO

* [X] Check existing file per queue and open and read it immediately
* [X] Parse CSV
* [X] Parse features float array
* [X] Probably send command, or add a feature of handling that file is close on data-stream
* [X] Skip header
* [X] Inter-buffer processing
* [X] Handle error of csv format - just skip record
* [X] Check a record for anomalies (call dummy preference)
* [X] Split output stream into two
* [X] Write record to file
* [X] Reference collector
* [X] Calculate Z-test of reference
* [X] Waiting for reference before calculating Z-test
* [X] Anomalies detector, which uses the reference collector
* [X] Dockerfile and build image
* [X] Push image to hub.docker.com
* [ ] Description of the project
* [ ] Configuration for delimiter, skipping the first line, size of buffer
* [ ] Pool for read buffer
* [ ] Unit-test of happy cases of watcher, data-stream, detector, reference, etc.
* [ ] Logger with levels

## References:

1. [Z-test]: (https://en.wikipedia.org/wiki/Z-test)
2. [Rapid calculation methods for standard deviation](https://en.wikipedia.org/wiki/Standard_deviation#Rapid_calculation_methods)
3. [Comparing means: z and t tests](https://mgimond.github.io/Stats-in-R/z_t_tests.html)
