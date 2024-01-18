package main

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/waku-org/go-waku/waku/v2/node"
)

func TestStoreSuite(t *testing.T) {
	suite.Run(t, new(StoreSuite))
}

type StoreSuite struct {
	suite.Suite
	node *node.WakuNode
}

func (s *StoreSuite) SetupSuite() {
	wakuNode, err := node.New(
		node.WithNTP(),
		node.WithWakuRelayAndMinPeers(1),
		node.WithClusterID(clusterID),
	)

	s.NoError(err)

	err = wakuNode.Start(context.Background())
	s.NoError(err)

	s.node = wakuNode
}

func (s *StoreSuite) TearDownSuite() {
	s.node.Stop()
}
