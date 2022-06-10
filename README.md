# Instructions

- Scylla code inspired from [here](https://github.com/scylladb/scylla-code-samples/tree/master/mms)

- songs content got from [here](https://github.com/socratica/sql)

- golang devcontainer copied from [here](https://github.com/microsoft/vscode-dev-containers/tree/main/containers/go)

- to run bash in the node:
    `docker exec -it scylla-node1 bash` followed with `nodetool status` and `cqlsh`

- Create schema and import test data: `docker exec scylla-node1 cqlsh -f /scylla-data.txt`

## Useful commands

### Using nodetool

- Print status of the cluster: `nodetool status`

- Fetch read/write latency, partition size...: `nodetool cfhistograms songs songs_by_year`

- In-depth diagnostics of a specific table: `nodetool tablestats songs.songs_by_year`

- probabilistic tracing: randomly chooses a request to be traced with some defined probability, example (0.01%) of all queries in node: `nodetool settraceprobability 0.0001`

- remove a node (irreversible action): `nodetool decommission` then `nodetool removenode node-id`. `node-id` can be found in `nodetool status`

- after adding a new cluster, we need to alter keyspace to replicate into new DC then run `nodetool rebuild -new_dc_name <existing_dc_name>` and run full cluster repair

### logs using journalctl

- get logs: `journalctl -u scylla-server`

- get logs since last server boot: `journalctl -u scylla -b`

# using cqlsh

- Enable tracing: `cqlsh tracing on|off`



