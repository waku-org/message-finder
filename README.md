# Message Finder
Retrieve WakuMessages from Storenodes


### How to use

1. Edit `main_test.go`.  Setup the following variables if neccesary
- `nodeList`: set the list of storenode multiaddresses
- `clusterID`: the cluster id used by the storenodes (this was set with the `--cluster-id` flag when running the store node)
- `pubsubTopic`: pubsub topic on which the message was published. In the status app, `"/waku/2/default-waku/proto"` is used in `status.prod` fleet and `"/waku/2/rs/16/32"` in `shards.test` fleet
- `contentTopics`: array of strings with content topics. In the status app use the following format `"/waku/1/0xaabbccdd/rfc26"`
- `startTime`: unix timestamp in nanoseconds
- `endTime`: unix timestamp in nanoseconds
- `envelopeHash`: the hash of the message to find. This is optional. (Use `0x` to not search for an specific message)
2. Execute `make`

The program will attempt to retrieve the messages that match the criteria described in the previous variables, and print basic details about them, as well as use any query cursor returned to retrieve more pages of results.
