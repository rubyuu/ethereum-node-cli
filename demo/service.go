package demo

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/node"
	"github.com/ethereum/go-ethereum/rpc"
	drpc "github.com/nolan/go-demo-server/demo/rpc"
	"log"
	"math/big"
	"net"
	"net/http"
	_ "net/http/pprof"
	"strconv"
	"sync/atomic"
	"time"
)

const MIN_SYNC_INTERVAL = 3 * time.Second

type DemoService struct {
	config      *CliConfig
	ethClient   *ethclient.Client
	pprofServer *http.Server
	rpcServer   *http.Server
	listener    net.Listener
	stopped     atomic.Bool
	appCtx      context.Context
	blockState  uint64
}

func DemoServiceFromCLIConfig(ctx context.Context, cfg *CliConfig) (*DemoService, error) {
	demoService := &DemoService{}
	demoService.config = cfg
	ctx, cancel := context.WithTimeout(ctx, 1*time.Minute)
	defer cancel()
	ethClient, err := ethclient.DialContext(ctx, cfg.EthRpc)
	if err != nil {
		return nil, err
	}
	demoService.ethClient = ethClient

	demoService.rpcServer = &http.Server{
		Addr: net.JoinHostPort(cfg.Host, strconv.Itoa(cfg.Port)),
	}

	if err = demoService.initRpcServer(); err != nil {
		return nil, err
	}

	if cfg.PprofEnabled {
		demoService.pprofServer = &http.Server{Addr: net.JoinHostPort(cfg.Host, strconv.Itoa(cfg.PprofPort)), Handler: nil}
	}
	return demoService, nil
}

func (d *DemoService) initRpcServer() error {
	// register apis
	srv := rpc.NewServer()
	if err := node.RegisterApis(drpc.GetAPIs(d), nil, srv); err != nil {
		return fmt.Errorf("error registering APIs: %w", err)
	}

	mux := http.NewServeMux()
	mux.Handle("/", srv)
	listener, err := net.Listen("tcp", d.rpcServer.Addr)
	if err != nil {
		return fmt.Errorf("error listening on %s: %w", d.rpcServer.Addr, err)
	}
	d.listener = listener
	// Override rpc server addr value, if port = 0
	d.rpcServer.Addr = listener.Addr().String()
	d.rpcServer.Handler = mux
	return nil
}

func (d *DemoService) Start(ctx context.Context) error {
	if d.rpcServer != nil {
		go func() {
			d.rpcServer.Serve(d.listener)
		}()
	}
	if d.pprofServer != nil {
		go func() {
			err := d.pprofServer.ListenAndServe()
			if err != nil {
				panic(err)
			}
		}()
	}
	if d.config.SyncEnabled && d.config.SyncInterval > 0 {
		go func() {
			interval := time.Duration(d.config.SyncInterval) * time.Second
			if interval < MIN_SYNC_INTERVAL {
				interval = MIN_SYNC_INTERVAL
			}
			timer := time.NewTimer(interval)

			for {
				select {
				case <-timer.C:
					log.Println("will get block number")
					number, err := d.ethClient.BlockNumber(ctx)
					if err != nil {
						log.Printf("error getting block number: %v", err)
					}
					d.blockState = number
					timer.Reset(interval)
				case <-ctx.Done():
					log.Printf("sync will stopped")
					return
				}
			}
		}()
	}
	return nil
}

func (d *DemoService) Stop(ctx context.Context) error {
	log.Println("stopping demo service")
	if d.stopped.Load() {
		return fmt.Errorf("already stopped")
	}
	if d.rpcServer != nil {
		d.rpcServer.Close()
	}
	if d.pprofServer != nil {
		d.pprofServer.Close()
	}
	if d.ethClient != nil {
		d.ethClient.Close()
	}
	d.stopped.Store(true)
	return nil
}

// Stopped determines if the service was stopped with Stop.
func (d *DemoService) Stopped() bool {
	return d.stopped.Load()
}

func (d *DemoService) GetBlockHashByNumber(ctx context.Context, number rpc.BlockNumber) (map[string]interface{}, error) {
	numberInt := new(big.Int).SetInt64(number.Int64())
	block, err := d.ethClient.BlockByNumber(ctx, numberInt)
	if err != nil {
		return nil, err
	}
	res := make(map[string]interface{})
	res["blockHash"] = block.Hash().String()
	return res, nil
}

func (d *DemoService) GetBlockNumber(ctx context.Context) (string, error) {
	res := big.NewInt(int64(d.blockState))
	return hexutil.EncodeBig(res), nil
}
