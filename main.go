package main

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/multiformats/go-multiaddr"
	"github.com/rodaine/table"
	cli "github.com/urfave/cli/v2"
	"github.com/waku-org/go-waku/logging"
	"github.com/waku-org/go-waku/waku/v2/node"
	"github.com/waku-org/go-waku/waku/v2/protocol"
	"github.com/waku-org/go-waku/waku/v2/protocol/legacy_store"
	"github.com/waku-org/go-waku/waku/v2/protocol/pb"
	"github.com/waku-org/go-waku/waku/v2/protocol/store"
	"github.com/waku-org/go-waku/waku/v2/utils"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Options struct {
	NodeKey       *ecdsa.PrivateKey
	ClusterID     uint
	PubSubTopic   string
	ContentTopics cli.StringSlice
	StartTime     int64
	EndTime       int64
	Hashes        cli.StringSlice
	AdvanceCursor bool
	PageSize      uint64
	StoreNode     *multiaddr.Multiaddr
	UseLegacy     bool
	QueryTimeout  time.Duration
	LogLevel      string
	LogEncoding   string
	LogOutput     string
}

var options Options

func main() {
	// Defaults
	options.LogLevel = "INFO"

	app := &cli.App{
		Name:    "query",
		Version: "0.0.1",
		Flags: []cli.Flag{
			NodeKey,
			ClusterID,
			PubsubTopic,
			ContentTopic,
			StartTime,
			EndTime,
			Hashes,
			Pagesize,
			Storenode,
			UseLegacy,
			Timeout,
			LogLevel,
			LogEncoding,
			LogOutput,
		},
		Action: func(c *cli.Context) error {
			if len(options.Hashes.Value()) == 1 {
				err := FetchMessage(c.Context, options)
				if err != nil {
					return cli.Exit(err.Error(), 1)
				}
			} else {
				err := QueryMessages(c.Context, options)
				if err != nil {
					return cli.Exit(err.Error(), 1)
				}
			}
			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		panic(err)
	}
}

func initializeWaku(opts Options) (*node.WakuNode, error) {
	utils.InitLogger(opts.LogEncoding, opts.LogOutput, "query")

	var prvKey *ecdsa.PrivateKey
	var err error

	if opts.NodeKey != nil {
		prvKey = opts.NodeKey
	} else {
		if prvKey, err = crypto.GenerateKey(); err != nil {
			return nil, fmt.Errorf("error generating key: %w", err)
		}
	}

	p2pPrvKey := utils.EcdsaPrivKeyToSecp256k1PrivKey(prvKey)
	id, err := peer.IDFromPublicKey(p2pPrvKey.GetPublic())
	if err != nil {
		return nil, err
	}
	logger := utils.Logger().With(logging.HostID("node", id))

	lvl, err := zapcore.ParseLevel(opts.LogLevel)
	if err != nil {
		return nil, err
	}

	libp2pOpts := append(node.DefaultLibP2POptions, libp2p.NATPortMap()) // Attempt to open ports using uPNP for NATed hosts.)

	wakuNode, err := node.New(
		node.WithLogger(logger),
		node.WithLogLevel(lvl),
		node.WithPrivateKey(prvKey),
		node.WithClusterID(uint16(options.ClusterID)),
		node.WithNTP(),
		node.WithLibP2POptions(libp2pOpts...),
	)
	if err != nil {
		return nil, fmt.Errorf("could not instantiate waku: %w", err)
	}

	return wakuNode, nil
}

func FetchMessage(ctx context.Context, opts Options) error {
	wakuNode, err := initializeWaku(opts)
	if err != nil {
		return err
	}

	if err = wakuNode.Start(ctx); err != nil {
		return err
	}
	defer wakuNode.Stop()

	h, err := hexutil.Decode(opts.Hashes.Value()[0])
	if err != nil {
		return fmt.Errorf("invalid message hash: %s", opts.Hashes.Value()[0])
	}

	ctx, cancel := context.WithTimeout(context.Background(), options.QueryTimeout)
	result, err := wakuNode.Store().Request(ctx, store.MessageHashCriteria{MessageHashes: []pb.MessageHash{pb.ToMessageHash(h)}},
		store.WithPeerAddr(*options.StoreNode),
		store.WithPaging(false, options.PageSize),
	)
	cancel()
	if err != nil {
		return err
	}

	if len(result.Messages()) == 0 {
		fmt.Println("Message not found")
		return nil
	}

	fmt.Println()

	msg := result.Messages()[0]

	fmt.Println("PubsubTopic:", msg.PubsubTopic)
	fmt.Println("MessageHash:", msg.WakuMessageHash())
	fmt.Println("ContentTopic:", msg.Message.ContentTopic)
	if msg.Message.Timestamp == nil {
		fmt.Println("Timestamp (unix nano): <nil>")
		fmt.Println("Timestamp (UTC): <nil>")
	} else {
		fmt.Println("Timestamp (unix nano):", msg.Message.GetTimestamp())
		fmt.Println("Timestamp (UTC):", time.Unix(0, msg.Message.GetTimestamp()))
	}

	if msg.Message.Version == nil {
		fmt.Println("Version: <nil>")
	} else {
		fmt.Println("Version:", *msg.Message.Version)
	}

	if len(msg.Message.Payload) != 0 {
		fmt.Println("Payload:")
		fmt.Println(hex.Dump(msg.Message.Payload))
	} else {
		fmt.Println("Payload: <nil>")
	}

	if len(msg.Message.Meta) != 0 {
		fmt.Println("Meta:")
		fmt.Println(hex.Dump(msg.Message.Meta))
	} else {
		fmt.Println("Meta: <nil>")
	}

	if len(msg.Message.RateLimitProof) != 0 {
		fmt.Println("RateLimitProof:")
		fmt.Println(hex.Dump(msg.Message.RateLimitProof))
	} else {
		fmt.Println("RateLimitProof: <nil>")
	}
	return nil

}

func QueryMessages(ctx context.Context, opts Options) error {
	wakuNode, err := initializeWaku(opts)
	if err != nil {
		return err
	}

	if err = wakuNode.Start(ctx); err != nil {
		return err
	}
	defer wakuNode.Stop()

	var hashes []pb.MessageHash
	if len(options.Hashes.Value()) != 0 {
		if options.PubSubTopic != "" || len(options.ContentTopics.Value()) != 0 || options.StartTime != 0 || options.EndTime != 0 {
			return errors.New("cannot specify pubsub topic / content topics / start time / end time if using the --hash flag")
		}

		if options.UseLegacy {
			return errors.New("cannot use legacy store with the --hash flag")
		}

		for _, hash := range options.Hashes.Value() {
			h, err := hexutil.Decode(hash)
			if err != nil {
				return fmt.Errorf("invalid message hash: %s", hash)
			}
			hashes = append(hashes, pb.ToMessageHash(h))
		}
	}

	var StartTime *int64
	if options.StartTime > 0 {
		StartTime = &options.StartTime
	}

	var EndTime *int64
	if options.EndTime > 0 {
		EndTime = &options.EndTime
	}

	cnt := 0
	if !options.UseLegacy {
		var criteria store.Criteria

		if len(hashes) == 0 {
			criteria = store.FilterCriteria{
				ContentFilter: protocol.NewContentFilter(options.PubSubTopic, options.ContentTopics.Value()...),
				TimeStart:     StartTime,
				TimeEnd:       EndTime,
			}
		} else {
			criteria = store.MessageHashCriteria{
				MessageHashes: hashes,
			}
		}

		ctx, cancel := context.WithTimeout(context.Background(), options.QueryTimeout)
		result, err := wakuNode.Store().Request(ctx, criteria,
			store.WithPeerAddr(*options.StoreNode),
			store.WithPaging(false, options.PageSize),
		)
		cancel()
		if err != nil {
			return err
		}

		pageCount := 0

		if len(result.Messages()) == 0 {
			fmt.Println("No messages found")
			return nil
		}

		fmt.Println()

		for !result.IsComplete() {
			if len(result.Messages()) == 0 {
				break
			}

			pageCount++
			cnt += len(result.Messages())

			tbl := table.New("MessageHash", "Content Topic", "Timestamp (unix nano)", "Timestamp (UTC)")
			for _, msg := range result.Messages() {
				tbl.AddRow(msg.WakuMessageHash(), msg.Message.ContentTopic, msg.Message.Timestamp, time.Unix(0, msg.Message.GetTimestamp()))
			}

			fmt.Printf("Page: %d, Record from %d to %d\n", pageCount, cnt-len(result.Messages())+1, cnt)

			tbl.Print()

			fmt.Println()

			if result.Cursor() != nil {
				fmt.Printf("Cursor: %s\n", hex.EncodeToString(result.Cursor()))
			}
			fmt.Println()

			ctx, cancel := context.WithTimeout(context.Background(), options.QueryTimeout)
			err := result.Next(ctx)
			cancel()
			if err != nil {
				return err
			}
		}

	} else {
		query := legacy_store.Query{
			PubsubTopic:   options.PubSubTopic,
			ContentTopics: options.ContentTopics.Value(),
			StartTime:     StartTime,
			EndTime:       EndTime,
		}

		ctx, cancel := context.WithTimeout(context.Background(), options.QueryTimeout)
		result, err := wakuNode.LegacyStore().Query(ctx, query,
			legacy_store.WithPeerAddr(*options.StoreNode),
			legacy_store.WithPaging(false, 20),
		)
		cancel()
		if err != nil {
			return err
		}

		if len(result.Messages) == 0 {
			fmt.Println("No messages found")
			return nil
		}

		fmt.Println()

		pageCount := 0
		for {
			pageCount++
			ctx, cancel := context.WithTimeout(context.Background(), options.QueryTimeout)
			hasNext, err := result.Next(ctx)
			cancel()
			if err != nil {
				return err
			}

			if !hasNext { // No more messages available
				break
			}

			cnt += len(result.GetMessages())
			tbl := table.New("MessageHash", "Content Topic", "Timestamp (unix nano)", "Timestamp (UTC)")
			for _, msg := range result.GetMessages() {
				env := protocol.NewEnvelope(msg, msg.GetTimestamp(), query.PubsubTopic)
				tbl.AddRow(env.Hash(), env.Message().ContentTopic, env.Message().GetTimestamp(), time.Unix(0, env.Message().GetTimestamp()))
			}

			fmt.Printf("Page: %d, Record from %d to %d\n", pageCount, cnt-len(result.GetMessages())+1, cnt)

			tbl.Print()

			fmt.Println()

			if result.Cursor() != nil {
				fmt.Printf("Cursor: Digest(%s); ReceiverTime:%d, SenderTime: %d, PubsubTopic: %s\n", hex.EncodeToString(result.Cursor().Digest), result.Cursor().ReceiverTime, result.Cursor().SenderTime, result.Cursor().PubsubTopic)
			}
			fmt.Println()

		}
	}

	utils.Logger().Info("Total messages retrieved", zap.Int("num", cnt))

	return nil
}
