package commons

import (
	"time"

	"github.com/Unknwon/goconfig"
	"github.com/cihub/seelog"
	"github.com/l1va/gofins/fins"
	"github.com/silenceper/pool"
)

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
		seelog.Error(e)
	}
	defer s.Close()

	c, err := fins.NewClient(clientAddr, plcAddr)
	if err != nil {
		seelog.Error(err)
	}
	defer c.Close()
	return c
}

func PoolInit(conf *goconfig.ConfigFile) (pool.Pool, error) {
	factory := func() (interface{}, error) { return OmronConnect(conf), nil }
	close := func(v interface{}) error {
		v.(*fins.Client).Close()
		return nil
	}
	poolConfig := &pool.PoolConfig{
		InitialCap:  5,
		MaxCap:      30,
		Factory:     factory,
		Close:       close,
		IdleTimeout: 15 * time.Second,
	}
	clientPool, err := pool.NewChannelPool(poolConfig)
	if err != nil {
		seelog.Error(err)
	}

	return clientPool, err
}
