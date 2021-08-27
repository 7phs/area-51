# Dataset anomalies :flying_saucer: detector

`area-51` is a tool for detecting analytics in a set of data using [Z-test](https://en.wikipedia.org/wiki/Z-test).
A statistics preferences of reference data use for scoring raw data.
The sources of data are files in CSV format.

## How to use

`area-51` is CLI tool to process data stored in files.

Reference data and raw data should store in two different files.

Raw data is split by `Z-score` and put into two files:
* Clear data is stored to `clean.csv`;
* Record contains data greater than expected deviation is stored to `anomalies.csv`.

The tool launches processes existing files and listens to reference and raw files until a user stops it.
It processes changes in data (appending, replacing a file, etc.) in real-time.

A user can launch the tool even source files are not in place.
The tool will process it if a user puts files on a place identified by CLI options `--reference` and `--raw`.

Log messages inform a user that existing data is completely processed.
A user stops the tool at the proper time, when files will not change, or another reason.

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

## Description of the project

There are two sides of the solution that I would like to describe in detail.

### Statistics

[Z-test](https://en.wikipedia.org/wiki/Z-test) uses for estimating the quality of a data record.

Needs to know to mean and standard deviation a feature of partition assigned to a data record to calculate `Z-score`.

There are several ways to do it:

* collect all data are included into partition;
* calculating mean and standard deviation on a stream of data.

The solution is implemented the second one based on the description at [Rapid calculation methods for the standard deviation](https://en.wikipedia.org/wiki/Standard_deviation#Rapid_calculation_methods):
> This is a "one pass" algorithm for calculating variance of n samples without the need to store prior data during the calculation. Applying this method to a time series will result in successive values of standard deviation corresponding to n data points as n grows larger with each new sample, rather than a constant-width sliding window calculation. 

### Engineering

The solution breaks implementation into two big components:
* reading data - responsible for listening of files changes reading, representing data stored in files and stream of bytes;
* processing data - responsible for parsing a stream of bytes into data records, scoring them, and storing them into destination files.

A major reason for it is representing data as an infinite stream for a processor.
It helps update a data source and possibly replace it with a stream from network services, etc.  

Important parts of the solution that were implemented to reduce the processing time of data files:

* Listen to OS file events to handle all changes of raw and reference files;
* Using a buffer to read data from a file and serialize it;
* Custom CSV reader to reducing the overhead of standard solution;
* Use slices of the data buffer to assign as record fields instead of copy data during record's fetching;
* Parse an only significant part of the data record (key and features);
* Keep raw representation of data record as a slice of bytes to easily serialize it.

### Benchmark

Data processing time measurements were taken to determine the overall performance level of the solution.

Test reference contains ~100 000 records, and raw data contains ~100 000 records.

Hardware:

* MacBook Pro (15-inch, 2018); 2,6 GHz 6-Core Intel Core i7

Result:

```
Load and process reference files (before process the first line of raw data):
      100 000 records / from 150 ms to 250 ms

Scoring raw data:
       50 000 records / from 65 ms to 100 ms  
```

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
* [X] Description of the project
* [ ] Configuration for delimiter, skipping the first line, size of buffer
* [ ] Output stream is a dedicated entity with own interface
* [ ] Pool for read buffer
* [ ] Unit-test of happy cases of watcher, data-stream, detector, reference, etc.
* [ ] Logger with levels

## References:

1. [Z-test](https://en.wikipedia.org/wiki/Z-test)
2. [Rapid calculation methods for standard deviation](https://en.wikipedia.org/wiki/Standard_deviation#Rapid_calculation_methods)
3. [Comparing means: z and t tests](https://mgimond.github.io/Stats-in-R/z_t_tests.html)
