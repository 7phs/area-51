# Dataset anomalies :flying_saucer: detector

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
* [ ] Waiting for reference before calculating Z-test
* [ ] Anomalies detector, which uses the reference collector 
* [ ] Dockerfile and build image
* [ ] Command for one time processing
* [ ] Description of the project
* [ ] Configuration for delimiter, skipping the first line, size of buffer
* [ ] Pool of buffer for float array[3]
* [ ] Pool for read buffer
* [ ] Unit-test of happy cases
* [ ] Logger with levels

BUGS:
* [ ] Last line missed 

  
https://github.com/montanaflynn/stats - MIT
https://machinelearningmastery.com/critical-values-for-statistical-hypothesis-testing/
https://mgimond.github.io/Stats-in-R/z_t_tests.html