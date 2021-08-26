# Dataset anomalies :flying_saucer: detector

## Tasks

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
* [ ] Push image to hub.docker.com
* [ ] Description of the project
* [ ] Configuration for delimiter, skipping the first line, size of buffer
* [ ] Pool for read buffer
* [ ] Strict number for go-routines to handle buffer instead of sending to channel with go-routine
* [ ] Unit-test of happy cases of watcher, data-stream, detector, reference, etc.
* [ ] Logger with levels

References:

1. [Z-test](https://en.wikipedia.org/wiki/Z-test)
2. [Rapid calculation methods for standard deviation](https://en.wikipedia.org/wiki/Standard_deviation#Rapid_calculation_methods)
3. [Comparing means: z and t tests](https://mgimond.github.io/Stats-in-R/z_t_tests.html)
