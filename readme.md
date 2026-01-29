# MapReduce

This is my implementation of MapReduce model in Go following the paper and the lecture from MIT 6.824.

### Notes

1. MapReduce is a programming model and an associated implementation for **processing and generating large dataset**.
2. Designed to run in a large cluster of commodity machines and is highly scalable.

### Process

- The computation takes a set of input key/value pairs, and produces a set of output key/value pairs.
- MapReduce library expresses the computation as two functions: **Map** and **Reduce**.
- **Map**, takes an input pair and produces a set of **intermediate key/value pairs**.
- **Reduce**, takes an intermediate key and a set of values for that key, then merges these values to form a possibly smaller set of values (zero or one output value is returned per invocation of this function).
- Map invocations are distributed across multiples machines, allowing parallel processing of splitted input data.
- Reduce invocations are distributed using a partitioning function e.g., `hash(key) mod R`.(The number of partitions(R) and the partitioning function is specified by the user)
- Flow:
  1. Splits the input data into M pieces typically 16mb to 64 mb per piece (controllable via option parameter).
  2. The **master** holds the copy of the splitted input data. And, assign task to idle workers. The master decides if the task is map or reduce.
  3. **Map Phase**, a worker who is assigned a map task reads the content of the corresponding input split and run a map function on each pair of the input. The resulting intermediate pair is buffered in the memory.
  4. Periodically, the buffered pairs are written to local disk(of the map worker), partitioned into R regions by the partitioning function. The location of the buffered pairs in the local disk are passed back to the master, who is responsible for **forwarding these locations** to the reduce worker.
  5. When reduce worker is notified by the master of the location. The reduce worker will read the buffered data in the local disk of map worker via **RPC**. Then, it sorts the intermediate so that all occurence of the same key are grouped together or close to each other (But, if the intermediate data is too large, an external is used ).
  6. The reduce worker iterates over the intermediate data. For each unique intermediate key it passes the key and the list of asscioated values to the reduce function. The returned output is appended to the output file for that partition.
  7. When all map and reduce task have been completed, the MapReduce call in the user programs returns back to the user code.

### Components

1. **Master**
   - Holds the copy of the splitted input data(the master is not responsible for splitting the input data) and assign task to idle workers.
   - Stores states(idle, in-progress, completed) and identity of the worker for each map and reduce task
   - The master is the conduit through which the location of intermediate file regions is propagated from map tasks to reduce tasks. Therefore, for each completed map task, the master stores the locations and sizes of the R intermediate file regions produced by the map task. Updates to this location and size information are received as map tasks are completed. The information is pushed incrementally to workers that have in-progress reduce tasks.
   - The master pings every workee periodically, if no response the master marks the worker as `failed`.
   - Any map task that is completed, the master will change its state back to `idle` and therefore become eligible for scheduling for other workers.

2. **Worker**
   - **Map**:
     - Produce the intermediate key/value pair from input key/value pair
   - **Reduce**:
     - A reduce worker cannot start while a map worker is still in progress

3. Helper
   - For splitting the data

### Thoughts

This is my first project learning about distributive systems.

### Source

1. https://static.googleusercontent.com/media/research.google.com/en//archive/mapreduce-osdi04.pdf
2. https://www.youtube.com/playlist?list=PLrw6a1wE39_tb2fErI4-WkMbsvGQk9_UB
