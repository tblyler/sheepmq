# SheepMQ
Aiming to be a robust message queue library and daemon with a variety of storage backends.

# Still a work in progress
## TODO
(In no particular order)
* Add documentation/roadmap
* Implement more daemon entrypoints
	* grpc (in progress)
	* http REST
	* custom?
	* if feeling ambitious... AMQP
* Implement more backends outside of [Badger](https://github.com/dgraph-io/badger)
	* in-memory
	* custom on-disk structure
	* [Bolt](https://github.com/boltdb/bolt)
* Unit test as much as possible within reason
* Hook in with CI
