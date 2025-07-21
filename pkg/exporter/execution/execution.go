package execution

import (
	"context"
	"time"

	"github.com/FISCO-BCOS/go-sdk/v3/client"
	"github.com/sirupsen/logrus"
	"github.com/trmaphi/bcos-metrics-exporter/pkg/exporter/execution/api"
)

// Node represents an execution node.
type Node interface {
	// Name returns the name of the node.
	Name() string
	// URL returns the url of the node.
	URL() string
	// Bootstrapped returns whether the node has been bootstrapped and is ready to be used.
	Bootstrapped() bool
	// Bootstrap attempts to bootstrap the node (i.e. configuring clients)
	Bootstrap(ctx context.Context) error
	// StartMetrics starts the metrics collection.
	StartMetrics(ctx context.Context)
}

type node struct {
	name        string
	url         string
	client      *client.Client
	internalAPI api.ExecutionClient
	log         logrus.FieldLogger
	metrics     Metrics
}

// NewExecutionNode returns a new execution node.
func NewExecutionNode(ctx context.Context, log logrus.FieldLogger, namespace, nodeName, url string, enabledModules []string) (Node, error) {
	internalAPI := api.NewExecutionClient(ctx, log, url)
	
	// Initialize FISCO-BCOS client
	// For FISCO-BCOS, we need to provide groupID and private key
	// Using default values for now - these should be configurable
	bcosClient, err := client.Dial(url, "1", []byte{})
	if err != nil {
		return nil, err
	}
	
	metrics := NewMetrics(bcosClient, internalAPI, log, nodeName, namespace, enabledModules)

	node := &node{
		name:        nodeName,
		url:         url,
		log:         log,
		internalAPI: internalAPI,
		client:      bcosClient,
		metrics:      metrics,
	}

	return node, nil
}

func (e *node) Name() string {
	return e.name
}

func (e *node) URL() string {
	return e.url
}

func (e *node) Bootstrapped() bool {
	return e.client != nil
}

func (e *node) Bootstrap(ctx context.Context) error {
	client, err := ethclient.Dial(e.url)
	if err != nil {
		return err
	}

	e.client = client

	return nil
}

func (e *node) StartMetrics(ctx context.Context) {
	for !e.Bootstrapped() {
		if err := e.Bootstrap(ctx); err != nil {
			e.log.WithError(err).Error("Failed to bootstrap node")
		}

		time.Sleep(5 * time.Second)
	}

	e.metrics.StartAsync(ctx)
}
