package util

import (
	"strconv"
	"strings"

	"github.com/cihub/seelog"

	"github.com/Unknwon/goconfig"
	"github.com/l1va/gofins/fins"
)

const (
	ReadBits = iota
	ReadStrings
	ReadWords
	ReadByte
	writeBits
	WriteStrings
	WriteWords
)

func MemArea() map[string]byte {
	MemAreaMap := map[string]byte{
		"ciobit":  fins.MemoryAreaCIOBit,
		"wrbit":   fins.MemoryAreaWRBit,
		"hrbit":   fins.MemoryAreaHRBit,
		"arbit":   fins.MemoryAreaARBit,
		"cioword": fins.MemoryAreaCIOWord,
		"wrword":  fins.MemoryAreaWRWord,
		"hrword":  fins.MemoryAreaHRWord,
		"arword":  fins.MemoryAreaARWord,
		"dmbit":   fins.MemoryAreaDMBit,
		"dmword":  fins.MemoryAreaDMWord,
	}
	return MemAreaMap
}

func LoadConf(configPath string) *goconfig.ConfigFile {
	cfg, err := goconfig.LoadConfigFile(configPath)
	if err != nil {
		seelog.Error("load config faild")
	}
	return cfg
}

func StringtoBoolList(s string) []bool {
	var data []bool
	k := strings.Split(s, "")
	for _, v := range k {
		bolStr, _ := strconv.ParseBool(v)
		data = append(data, bolStr)
	}
	return data
}

func StringtoUintList(s string) []uint16 {
	var dataList []uint16
	k := strings.Split(s, "")
	for _, v := range k {
		var data uint16
		bolStr, _ := strconv.ParseUint(v, 16, 16)
		data = uint16(bolStr)
		dataList = append(dataList, data)
	}
	return dataList
}
