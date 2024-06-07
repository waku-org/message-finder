package main

import (
	"time"

	cli "github.com/urfave/cli/v2"
	"github.com/urfave/cli/v2/altsrc"
	"github.com/waku-org/go-waku/waku/cliutils"
)

var NodeKey = cliutils.NewGenericFlagSingleValue(&cli.GenericFlag{
	Name:  "nodekey",
	Usage: "P2P node private key as hex.",
	Value: &cliutils.PrivateKeyValue{
		Value: &options.NodeKey,
	},
})

var ClusterID = altsrc.NewUintFlag(&cli.UintFlag{
	Name:        "cluster-id",
	Value:       0,
	Usage:       "Cluster id that the node is running in. Node in a different cluster id is disconnected.",
	Destination: &options.ClusterID,
})

var PubsubTopic = altsrc.NewStringFlag(&cli.StringFlag{
	Name:        "pubsub-topic",
	Usage:       "Query pubsub topic.",
	Destination: &options.PubSubTopic,
})

var ContentTopic = altsrc.NewStringSliceFlag(&cli.StringSliceFlag{
	Name:        "content-topic",
	Usage:       "Query content topic. Argument may be repeated.",
	Destination: &options.ContentTopics,
})

var StartTime = altsrc.NewInt64Flag(&cli.Int64Flag{
	Name:        "start-time",
	Usage:       "Query start time in nanoseconds",
	Destination: &options.StartTime,
})

var EndTime = altsrc.NewInt64Flag(&cli.Int64Flag{
	Name:        "end-time",
	Usage:       "Query end time in nanoseconds",
	Destination: &options.EndTime,
})

var Hashes = altsrc.NewStringSliceFlag(&cli.StringSliceFlag{
	Name:        "hash",
	Usage:       "Query by message hashes",
	Destination: &options.Hashes,
})

var Pagesize = altsrc.NewUint64Flag(&cli.Uint64Flag{
	Name:        "pagesize",
	Value:       20,
	Usage:       "Pagesize",
	Destination: &options.PageSize,
})

var Storenode = cliutils.NewGenericFlagSingleValue(&cli.GenericFlag{
	Name:  "storenode",
	Usage: "Multiaddr of a peer that supports store protocol",
	Value: &cliutils.MultiaddrValue{
		Value: &options.StoreNode,
	},
	Required: true,
})

/*
	TODO
	altsrc.NewBoolFlag(&cli.BoolFlag{
		Name:        "advance-cursor",
		Usage:       "Advance cursor automatically",
		Destination: &options.AdvanceCursor,
		Value:       true,
	}),
*/

var UseLegacy = altsrc.NewBoolFlag(&cli.BoolFlag{
	Name:        "use-legacy",
	Usage:       "Use legacy store",
	Destination: &options.UseLegacy,
})

var Timeout = altsrc.NewDurationFlag(&cli.DurationFlag{
	Name:        "timeout",
	Usage:       "timeout for each individual store query request",
	Destination: &options.QueryTimeout,
	Value:       1 * time.Minute,
})

var LogLevel = cliutils.NewGenericFlagSingleValue(&cli.GenericFlag{
	Name:    "log-level",
	Aliases: []string{"l"},
	Value: &cliutils.ChoiceValue{
		Choices: []string{"DEBUG", "INFO", "WARN", "ERROR", "DPANIC", "PANIC", "FATAL"},
		Value:   &options.LogLevel,
	},
	Usage: "Define the logging level (allowed values: DEBUG, INFO, WARN, ERROR, DPANIC, PANIC, FATAL)",
})

var LogEncoding = cliutils.NewGenericFlagSingleValue(&cli.GenericFlag{
	Name:  "log-encoding",
	Usage: "Define the encoding used for the logs (allowed values: console, nocolor, json)",
	Value: &cliutils.ChoiceValue{
		Choices: []string{"console", "nocolor", "json"},
		Value:   &options.LogEncoding,
	},
})

var LogOutput = altsrc.NewStringFlag(&cli.StringFlag{
	Name:        "log-output",
	Value:       "file",
	Usage:       "specifies where logging output should be written  (stdout, file, file:./filename.log)",
	Destination: &options.LogOutput,
})
