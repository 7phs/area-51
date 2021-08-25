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
* [ ] Reference collector
* [ ] Calculate Z-test of reference
* [ ] Anomalies detector, which uses the reference collector 
* [ ] Dockerfile and build image
* [ ] Command for one time processing
* [ ] Description of the project
* [ ] Configuration for delimiter and skipping the first line
* [ ] Pool of buffer for float array[3]
* [ ] Pool for read buffer
* [ ] Unit-test of happy cases
* [ ] Logger with levels

BUGS:
* [ ] Last line missed