# Instructions

- golang devcontainer copied from [here](https://github.com/microsoft/vscode-dev-containers/tree/main/containers/go)

- to run bash in the node:
    `docker exec -it scylla-node1 bash` followed with `nodetool status` and `cqlsh`

- Create schema and import test data: `docker exec scylla-node1 cqlsh -f /scylla-data.txt`