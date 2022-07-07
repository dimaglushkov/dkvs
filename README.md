# dkvs

Distributed Key-Value Storage

## Running the application
1. start storage instances by running `docker-compose up --scale storage=<num_of_instances> storage`
2. collect active ports of storage instances by running `docker ps`
3. put collected IP addresses and ports into `controller/storages.txt` in the following format: `<IP addr>:<port>`
4. start controller instance by running `docker-compose up runner`

## Usage
After everything is set up, the service will be available at port `7341` through HTTP:
* Getting value: GET request to `/<key>`
* Putting value: PUT/POST request to `/` with JSON in the following format: `{"key": <key>, "value": <val>}`
* Deleting value: DELETE request to `/<key>`

## Things that were considered 
1. Using Kubernetes would solve a lot of issues with cluster management
2. Using consensus algorithms (Paxos/raft) to provide better data consistency across the storage cluster
3. Using better and constant hashing algorithm which going to allow to dynamically add new nodes to the cluster
4. Adding separate layer for data shards and forcing nodes to keep N shard replicas across the cluster even after some nodes go down
