package rpc

import (
	"context"
	"github.com/ethereum/go-ethereum/rpc"
)

type DemoDriver interface {
	GetBlockHashByNumber(ctx context.Context, number rpc.BlockNumber) (map[string]interface{}, error)
	GetBlockNumber(ctx context.Context) (string, error)
}

type demoAPI struct {
	ds DemoDriver
}

func GetAPIs(ds DemoDriver) []rpc.API {
	return []rpc.API{
		{
			Namespace: "demo",
			Service:   NewDemoAPI(ds),
		},
	}
}

func NewDemoAPI(ds DemoDriver) *demoAPI {
	return &demoAPI{
		ds,
	}
}

func (da *demoAPI) GetBlockHashByNumber(ctx context.Context, number rpc.BlockNumber) (map[string]interface{}, error) {
	return da.ds.GetBlockHashByNumber(ctx, number)
}

func (da *demoAPI) GetBlockNumber(ctx context.Context) (string, error) {
	return da.ds.GetBlockNumber(ctx)
}
