# Dataset anomalies :flying_saucer: detector

## TODO

* [X] Check existing file per queue and open and read it immediately
* [X] Parse CSV 
* [X] Parse features float array
* [ ] Reference collector
* [ ] Calculate Z-test of reference
* [ ] Anomalies detector, which uses the reference collector 
* [ ] Check a record for anomalies
* [ ] Split output stream into two
* [ ] Write record to file
* [ ] Inter-buffer processing
* [ ] Probably send command, or add a feature of handling that file is close on data-stream
* [ ] Skip header
* [ ] Command for one time processing
* [ ] Dockerfile and build image
* [ ] Description of the project
* [ ] Handle csv format records
* [ ] Pool of buffer for float array[3]
* [ ] Pool for read buffer
* [ ] Unit-test of happy cases
* [ ] Logger with levels