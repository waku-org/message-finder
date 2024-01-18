package main

import (
	"bytes"
	"context"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/multiformats/go-multiaddr"
	"github.com/waku-org/go-waku/logging"
	"github.com/waku-org/go-waku/waku/v2/node"
	"github.com/waku-org/go-waku/waku/v2/peerstore"
	"github.com/waku-org/go-waku/waku/v2/protocol/store"
	"github.com/waku-org/go-waku/waku/v2/utils"
	"go.uber.org/zap"

	"google.golang.org/protobuf/proto"
)

var log = utils.Logger().Named("message-finder")

func addNodes(ctx context.Context, node *node.WakuNode) {
	for _, addr := range nodeList {
		ma, err := multiaddr.NewMultiaddr(addr)
		if err != nil {
			log.Error("invalid multiaddress", zap.Error(err), zap.String("addr", addr))
			continue
		}

		_ = ma

		_, err = node.AddPeer(ma, peerstore.Static, []string{string(store.StoreID_v20beta4)})
		if err != nil {
			log.Error("could not add peer", zap.Error(err), zap.Stringer("addr", ma))
			continue
		}
	}
}

func queryNode(ctx context.Context, node *node.WakuNode, addr string, pubsubTopic string, contentTopics []string, startTime time.Time, endTime time.Time, envelopeHash []byte) (int, error) {
	p, err := multiaddr.NewMultiaddr(addr)
	if err != nil {
		return -1, err
	}

	info, err := peer.AddrInfoFromP2pAddr(p)
	if err != nil {
		return -1, err
	}

	cnt := 0
	cursorIterations := 0

	result, err := node.Store().Query(ctx, store.Query{
		PubsubTopic:   pubsubTopic,
		ContentTopics: contentTopics,
		StartTime:     proto.Int64(startTime.UnixNano()),
		EndTime:       proto.Int64(endTime.UnixNano()),
	}, store.WithPeer(info.ID), store.WithPaging(false, 20), store.WithRequestID([]byte{1, 2, 3, 4, 5, 6, 7, 8}))
	if err != nil {
		return -1, err
	}

	for {
		hasNext, err := result.Next(ctx)
		if err != nil {
			return -1, err
		}

		if !hasNext { // No more messages available
			break
		}
		cursorIterations += 1

		// uncomment to find message by ID
		for _, m := range result.GetMessages() {
			if len(envelopeHash) != 0 && bytes.Equal(m.Hash(pubsubTopic), envelopeHash) {
				log.Info("⚠️⚠️⚠️ MESSAGE FOUND ⚠️⚠️⚠️", logging.HexBytes("envelopeHash", envelopeHash), logging.HostID("storeNode", info.ID))
				return 0, nil
			} else {
				log.Info(hex.EncodeToString(m.Hash(pubsubTopic)), zap.String("contentTopic", m.ContentTopic), zap.String("timestamp", fmt.Sprintf("%d", m.GetTimestamp())), logging.HostID("storeNode", info.ID), zap.Int("page", cursorIterations))
			}
		}

		cnt += len(result.GetMessages())
	}

	log.Info(fmt.Sprintf("%d messages found in %s (Used cursor %d times)\n", cnt, info.ID, cursorIterations))

	return cnt, nil
}
