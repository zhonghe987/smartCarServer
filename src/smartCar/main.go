package main

import (
	"net/http"
	"smartCar/controller"
	"smartCar/model"
	"smartCar/util"

	"github.com/cihub/seelog"
	restful "github.com/emicklei/go-restful"
)

func main() {

	logger, err := seelog.LoggerFromConfigAsFile("c:\\conf\\log.xml")

	if err != nil {
		seelog.Critical("err parsing config log file", err)
		return
	}
	seelog.ReplaceLogger(logger)
	defer seelog.Flush()

	config := util.LoadConf("c:\\conf\\conf.ini")
	carInfoResource := controller.NewCarInfoResource(config)

	go model.NewTaskWorker(config)
	restful.Add(carInfoResource.WebService())
	seelog.Info(http.ListenAndServe(":8080", nil))

}
