package jobs

import (
	"context"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/onrik/ethrpc"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
	"github.com/trmaphi/bcos-metrics-exporter/pkg/exporter/execution/api"
)

// TXPool collects metrics around the transaction pool.
type TXPool struct {
	client       *ethclient.Client
	api          api.ExecutionClient
	ethRPCClient *ethrpc.EthRPC
	log          logrus.FieldLogger
	Transactions prometheus.GaugeVec
}

const (
	NameTxPool = "txpool"
)

func (t *TXPool) Name() string {
	return NameTxPool
}

func (t *TXPool) RequiredModules() []string {
	return []string{"txpool"}
}

// NewTXPool creates a new TXPool instance.
func NewTXPool(client *ethclient.Client, internalAPI api.ExecutionClient, ethRPCClient *ethrpc.EthRPC, log logrus.FieldLogger, namespace string, constLabels map[string]string) TXPool {
	constLabels["module"] = NameTxPool

	namespace += "_txpool"

	return TXPool{
		client:       client,
		api:          internalAPI,
		ethRPCClient: ethRPCClient,
		log:          log.WithField("module", NameGeneral),
		Transactions: *prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace:   namespace,
				Name:        "transactions",
				Help:        "How many transactions are in the txpool.",
				ConstLabels: constLabels,
			},
			[]string{
				"status",
			},
		),
	}
}

func (t *TXPool) Start(ctx context.Context) {
	t.tick(ctx)

	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(time.Second * 15):
			t.tick(ctx)
		}
	}
}

func (t *TXPool) tick(ctx context.Context) {
	if err := t.GetStatus(ctx); err != nil {
		t.log.Errorf("Failed to get txpool status: %s", err)
	}
}

func (t *TXPool) GetStatus(ctx context.Context) error {
	status, err := t.api.TXPoolStatus(ctx)
	if err != nil {
		return err
	}

	t.Transactions.WithLabelValues("pending").Set(float64(status.Pending))
	t.Transactions.WithLabelValues("queued").Set(float64(status.Queued))

	return nil
}
