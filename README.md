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
* Putting value: PUT/POST request to `/` with JSON in the following format: `{key: <key>, "value": <val>}`
* Deleting value: DELETE request to `/<key>`