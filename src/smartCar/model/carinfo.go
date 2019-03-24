package model

import (
	"fmt"
	"smartCar/commons"
	"smartCar/util"

	"github.com/Unknwon/goconfig"
	"github.com/cihub/seelog"
	"github.com/go-xorm/xorm"
	//"smartCar/util"
	"strconv"
	"strings"
	"time"

	"github.com/l1va/gofins/fins"
)

type dataInterface interface{}

type CarOperateFins struct {
	MemoryAreaMap map[string]byte
}

type CarInfo struct {
	CarID       string    `xorm:"not null pk VARCHAR(32)"`
	OperateType string    `xorm: "not null VARCHAR(32)"`
	Info        string    `xorm:"VARCHAR(512)"`
	CreateAt    time.Time `xorm:"NOT NULL created"`
	DeletedAt   time.Time `xorm:"deleted"`
}

func (cf *CarInfo) FindAll(sqlDB *xorm.Engine) ([]CarInfo, error) {
	listCars := []CarInfo{}
	err := sqlDB.Find(&listCars)
	return listCars, err
}

func (cf *CarInfo) Insert(sqlDB *xorm.Engine, car CarInfo) error {

	_, err := sqlDB.Insert(car)
	return err
}

func (cf *CarInfo) Delete(sqlDB *xorm.Engine, carId string) error {
	var carinfo = new(CarInfo)
	_, err := sqlDB.Where("car_id = ?", carId).Get(carinfo)
	if err != nil {
		seelog.Error("delete car faild, the reason is not fount car")
	}
	_, err = sqlDB.Where("car_id = ?", carId).Delete(carinfo)
	if err != nil {
		seelog.Error("marker delete failed")
	}

	_, err = sqlDB.Where("car_id = ?", carId).Unscoped().Get(carinfo)
	if err != nil {
		seelog.Error("not fount car with maker")
	}
	_, err = sqlDB.Where("car_id = ?", carId).Unscoped().Delete(carinfo)
	return err
}

func (cp *CarOperateFins) ReadBits(client fins.Client, carinfo string) ([]bool, error) {
	infos := strings.Split(carinfo, ",")
	var readAddress uint16
	var readOffSet uint8
	var readCount uint16

	mArea := cp.MemoryAreaMap[infos[0]]
	address, err := strconv.ParseUint(infos[1], 16, 16)
	if err != nil {
		seelog.Error("get error address")
	}
	offset, err := strconv.ParseUint(infos[2], 8, 8)
	if err != nil {
		seelog.Error("get error offset")
	}
	count, err := strconv.ParseUint(infos[3], 16, 16)
	if err != nil {
		seelog.Error("get error count")
	}

	readAddress = uint16(address)
	readOffSet = byte(offset)
	readCount = uint16(count)
	b, err := client.ReadBits(mArea, readAddress, readOffSet, readCount)
	if err != nil {
		seelog.Error("get error b")
	}
	return b, err
}

func (cp *CarOperateFins) ReadBytes(client fins.Client, carinfo string) ([]byte, error) {
	infos := strings.Split(carinfo, ",")
	var readAddress uint16
	var readCount uint16
	memArea := cp.MemoryAreaMap[infos[0]]
	address, err := strconv.ParseUint(infos[1], 16, 16)
	if err != nil {
		seelog.Error("get error address")
	}
	count, err := strconv.ParseUint(infos[2], 16, 16)
	if err != nil {
		seelog.Error("get error, count")
	}
	readAddress = uint16(address)
	readCount = uint16(count)
	b, err := client.ReadBytes(memArea, readAddress, readCount)
	if err != nil {
		seelog.Error("get error b")
	}
	return b, err
}

func (cp *CarOperateFins) ReadWords(client fins.Client, carinfo string) ([]uint16, error) {
	infos := strings.Split(carinfo, ",")
	var readAddress uint16
	var readCount uint16
	memArea := cp.MemoryAreaMap[infos[0]]
	address, err := strconv.ParseUint(infos[1], 16, 16)
	if err != nil {
		seelog.Error("get error address")
	}
	count, err := strconv.ParseUint(infos[2], 16, 16)
	if err != nil {
		seelog.Error("get error count")
	}
	readAddress = uint16(address)
	readCount = uint16(count)
	b, err := client.ReadWords(memArea, readAddress, readCount)
	if err != nil {
		seelog.Error("get error count")
	}
	return b, err
}

func (cp *CarOperateFins) ReadString(client fins.Client, carinfo string) (*string, error) {
	infos := strings.Split(carinfo, ",")
	var readAddress uint16
	var readCount uint16
	memArea := cp.MemoryAreaMap[infos[0]]
	address, err := strconv.ParseUint(infos[1], 16, 16)
	if err != nil {
		seelog.Error("get error address")
	}
	count, err := strconv.ParseUint(infos[2], 16, 16)
	if err != nil {
		seelog.Error("get error count")
	}
	readAddress = uint16(address)
	readCount = uint16(count)
	b, err := client.ReadString(memArea, readAddress, readCount)
	if err != nil {
		seelog.Error("get error b")
	}
	return b, err
}

func (cp *CarOperateFins) WriteBits(client fins.Client, carinfo string) error {
	infos := strings.Split(carinfo, ",")
	var readAddress uint16
	var readOffSet uint8
	var writeData []bool
	memArea := cp.MemoryAreaMap[infos[0]]
	address, err := strconv.ParseUint(infos[1], 16, 16)
	if err != nil {
		seelog.Error("get error address")
	}
	offset, err := strconv.ParseUint(infos[2], 16, 16)
	if err != nil {
		seelog.Error("get error offset")
	}

	readAddress = uint16(address)
	readOffSet = byte(offset)
	writeData = util.StringtoBoolList(infos[3])
	err = client.WriteBits(memArea, readAddress, readOffSet, writeData)
	if err != nil {
		seelog.Error("get error offset")
	}
	return err
}

func (cp *CarOperateFins) WriteWords(client fins.Client, carinfo string) error {
	infos := strings.Split(carinfo, ",")
	var readAddress uint16
	memArea := cp.MemoryAreaMap[infos[0]]
	address, err := strconv.ParseUint(infos[1], 16, 16)
	if err != nil {
		seelog.Error("get error address")
	}

	readAddress = uint16(address)
	data := util.StringtoUintList(infos[2])
	err = client.WriteWords(memArea, readAddress, data)
	if err != nil {
		seelog.Error("get error")
	}
	return err
}

func (cp *CarOperateFins) WriteString(client fins.Client, carinfo string) error {
	infos := strings.Split(carinfo, ",")
	var readAddress uint16
	var readOffSet uint16
	memArea := cp.MemoryAreaMap[infos[0]]
	address, err := strconv.ParseUint(infos[1], 16, 16)
	if err != nil {
		seelog.Error("get error address")
	}
	itemcount, err := strconv.ParseUint(infos[2], 16, 16)
	if err != nil {
		seelog.Error("get error itemcount")
	}
	readAddress = uint16(address)
	readOffSet = uint16(itemcount)
	err = client.WriteString(memArea, readAddress, readOffSet, infos[3])
	if err != nil {
		seelog.Error("get error")
	}
	return err
}

func (cp *CarOperateFins) WriteCarOperate(client fins.Client, info string, opType string) error {
	switch opType {
	case "writeBits":
		return cp.WriteBits(client, info)
	case "writeWords":
		return cp.WriteWords(client, info)
	case "writeString":
		return cp.WriteString(client, info)
	}
	err := fmt.Errorf("false")
	return err
}

func (cp *CarOperateFins) ReadCarOperate(client fins.Client, info string, opType string) (dataInterface, error) {
	switch opType {
	case "readBits":
		return cp.ReadBits(client, info)
	case "readBytes":
		return cp.ReadBytes(client, info)
	case "readWords":
		return cp.ReadWords(client, info)
	case "readString":
		return cp.ReadString(client, info)
	}
	err := fmt.Errorf("false")
	return false, err
}

func SendData(conf *goconfig.ConfigFile, operate *CarOperateFins, pipeline chan CarInfo, deleteChan chan string) {
	finsClient := commons.OmronConnect(conf)
	for {
		select {
		case data := <-pipeline:
			if strings.Contains(data.OperateType, "write") {
				err := operate.WriteCarOperate(*finsClient, data.Info, data.OperateType)
				if err != nil {
					seelog.Error(err)
				}
				deleteChan <- data.CarID
			} else {
				_, err := operate.ReadCarOperate(*finsClient, data.Info, data.OperateType)
				if err != nil {
					seelog.Error(err)
				}
			}
		default:
			seelog.Info("default")
		}
	}

}

func GetDataFromDB(sqlDB *xorm.Engine, conf *goconfig.ConfigFile, pipeline chan CarInfo) {
	listCars := []CarInfo{}
	err := sqlDB.Find(&listCars)
	if err != nil {
		seelog.Error("db select error")
	}
	for _, v := range listCars {
		pipeline <- v
	}

}

func DeleteData(deleteChan chan string, conf *goconfig.ConfigFile) {
	dbapi := commons.DBApi(conf)
	for {
		select {
		case id := <-deleteChan:
			dbapi.Delete(id)
		default:
			seelog.Info("frees")
		}
	}

}

func NewTaskWorker(conf *goconfig.ConfigFile) {
	ticker := time.NewTicker(360 * time.Second)
	data_chan := make(chan CarInfo)
	delete_chan := make(chan string)
	defer close(data_chan)
	defer close(delete_chan)
	dbapi := commons.DBApi(conf)
	operate := &CarOperateFins{util.MemArea()}
	go func() {
		for t := range ticker.C {
			fmt.Println("\n", t)
			GetDataFromDB(dbapi, conf, data_chan)
		}
	}()
	SendData(conf, operate, data_chan, delete_chan)
	DeleteData(delete_chan, conf)
}
