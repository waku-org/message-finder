package main

import (
	"context"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common/hexutil"
)

// If using vscode, go to Preferences > Settings, and edit Go: Test Timeout to at least 60s

// List of store nodes
var nodeList = []string{
	"/dns4/store-01.do-ams3.shards.test.statusim.net/tcp/30303/p2p/16Uiu2HAmAUdrQ3uwzuE4Gy4D56hX6uLKEeerJAnhKEHZ3DxF1EfT",
	"/dns4/store-02.do-ams3.shards.test.statusim.net/tcp/30303/p2p/16Uiu2HAm9aDJPkhGxc2SFcEACTFdZ91Q5TJjp76qZEhq9iF59x7R",
	"/dns4/store-01.gc-us-central1-a.shards.test.statusim.net/tcp/30303/p2p/16Uiu2HAmMELCo218hncCtTvC2Dwbej3rbyHQcR8erXNnKGei7WPZ",
	"/dns4/store-02.gc-us-central1-a.shards.test.statusim.net/tcp/30303/p2p/16Uiu2HAmJnVR7ZzFaYvciPVafUXuYGLHPzSUigqAmeNw9nJUVGeM",
	"/dns4/store-01.ac-cn-hongkong-c.shards.test.statusim.net/tcp/30303/p2p/16Uiu2HAm2M7xs7cLPc3jamawkEqbr7cUJX11uvY7LxQ6WFUdUKUT",
	"/dns4/store-02.ac-cn-hongkong-c.shards.test.statusim.net/tcp/30303/p2p/16Uiu2HAm9CQhsuwPR54q27kNj9iaQVfyRzTGKrhFmr94oD8ujU6P",
}

var clusterID uint16 = 16 // shards.test and status.prod = 16

// CRITERIA --------------------------------------------------------------------------
var pubsubTopic = "/waku/2/rs/16/32"              // "/waku/2/default-waku/proto" in status.prod and "/waku/2/rs/16/32" in shards.test
var contentTopics = []string{}                    // []string{"/waku/1/0xaabbccdd/rfc26"}
var startTime = time.Now().Add(-20 * time.Minute) // time.Unix(0, 1705486902684656000).Add(-60 * time.Second)
var endTime = time.Now()                          // time.Unix(0, 1705486902684656000).Add(60 * time.Second)
var envelopeHash = "0x"                           // Use "0x" to find all messages that match the pubsub topic, content topic and start/end time

func (s *StoreSuite) TestFindMessage() {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	addNodes(ctx, s.node)
	hash, err := hexutil.Decode(envelopeHash)
	if err != nil {
		panic("invalid envelope hash id")
	}

	wg := sync.WaitGroup{}
	for _, addr := range nodeList {
		wg.Add(1)
		func(addr string) {
			defer wg.Done()
			_, err := queryNode(ctx, s.node, addr, pubsubTopic, contentTopics, startTime, endTime, hash)
			s.NoError(err)
		}(addr)
	}
	wg.Wait()
}
