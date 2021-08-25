# Dataset anomalies :flying_saucer: detector

## TODO

* [X] Check existing file per queue and open and read it immediately
* [X] Parse CSV 
* [X] Parse features float array
* [ ] Inter-buffer processing
* [ ] Reference collector
* [ ] Calculate Z-test of reference
* [ ] Anomalies detector, which uses the reference collector 
* [ ] Check a record for anomalies
* [ ] Split output stream into two
* [ ] Write record to file
* [ ] Probably send command, or add a feature of handling that file is close on data-stream
* [ ] Skip header
* [ ] Dockerfile and build image
* [ ] Command for one time processing
* [ ] Description of the project
* [ ] Handle csv format records
* [ ] Pool of buffer for float array[3]
* [ ] Pool for read buffer
* [ ] Unit-test of happy cases
* [ ] Logger with levels