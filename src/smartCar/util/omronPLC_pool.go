package util

import (
	"context"

	"github.com/Unknwon/goconfig"
	pool "github.com/jolestar/go-commons-pool"
	"github.com/l1va/gofins/fins"
)

type OmronPLCPool struct {
	client *fins.Client
}

func OmronConnect(conf *goconfig.ConfigFile) *fins.Client {
	host_address, _ := conf.GetValue("host", "host_address")
	host_port, _ := conf.Int("host", "host_port")
	hostNetwork, _ := conf.Int("host", "host_network")
	host_network := byte(hostNetwork)
	hostNode, _ := conf.Int("host", "host_node")
	host_node := byte(hostNode)
	hostUnit, _ := conf.Int("host", "host_unit")
	host_unit := byte(hostUnit)

	plc_address, _ := conf.GetValue("plc", "plc_address")
	plc_port, _ := conf.Int("plc", "plc_port")
	plcNetwork, _ := conf.Int("plc", "plc_network")
	plc_network := byte(plcNetwork)
	plcNode, _ := conf.Int("plc", "plc_node")
	plc_node := byte(plcNode)
	plcUnit, _ := conf.Int("plc", "plc_unit")
	plc_unit := byte(plcUnit)

	clientAddr := fins.NewAddress(host_address, host_port, host_network, host_node, host_unit)
	plcAddr := fins.NewAddress(plc_address, plc_port, plc_network, plc_node, plc_unit)

	s, e := fins.NewPLCSimulator(plcAddr)
	if e != nil {
		panic(e)
	}
	defer s.Close()

	c, err := fins.NewClient(clientAddr, plcAddr)
	if err != nil {
		panic(err)
	}
	defer c.Close()
	return c
}

func PoolInit(conf *goconfig.ConfigFile) *pool.ObjectPool {
	ctx := context.Background()
	PoolConfig := pool.NewDefaultPoolConfig()
	PoolConfig.MaxTotal, _ = conf.Int("default", "pool_max")
	WithAbandonedConfig := pool.NewDefaultAbandonedConfig()
	pCommonPool := pool.NewObjectPoolWithAbandonedConfig(ctx, pool.NewPooledObjectFactorySimple(
		func(context.Context) (interface{}, error) {
			return &OmronPLCPool{
					client: OmronConnect(conf),
				},
				nil
		}), PoolConfig, WithAbandonedConfig)
	return pCommonPool
}
