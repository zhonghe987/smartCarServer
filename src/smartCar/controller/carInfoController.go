package controller

import (
	"net/http"
	"smartCar/commons"
	"smartCar/model"
	"smartCar/util"

	"github.com/cihub/seelog"
	"github.com/l1va/gofins/fins"

	"github.com/Unknwon/goconfig"
	"github.com/emicklei/go-restful"
	"github.com/go-xorm/xorm"
	pool "github.com/silenceper/pool"
)

type CarInfoResource struct {
	SqlDB      *xorm.Engine
	ClientPool pool.Pool
	Operate    model.CarOperateFins
}

func NewCarInfoResource(config *goconfig.ConfigFile) CarInfoResource {
	sqlit3DB := commons.DBApi(config)
	if err := sqlit3DB.Sync2(new(model.CarInfo)); err != nil {
		seelog.Error("Fail to sync database: %v\n", err)
	}
	pool, err := commons.PoolInit(config)
	if err != nil {
		seelog.Error("Init Omron fins client pool: %v\n", err)
	}
	operate := model.CarOperateFins{util.MemArea()}
	return CarInfoResource{sqlit3DB, pool, operate}
}

func (c CarInfoResource) WebService() *restful.WebService {
	ws := new(restful.WebService)
	ws.
		Path("/cars").
		Consumes(restful.MIME_XML, restful.MIME_JSON).
		Produces(restful.MIME_JSON, restful.MIME_XML) // you can specify this per route as well

	//tags := []string{"cars"}

	ws.Route(ws.GET("/").To(c.findAllCars))
	// docs
	//Doc("get all cars info").
	//Metadata(restfulspec.KeyOpenAPITags, tags).
	//Writes([]CarInfo{}).
	//Returns(200, "OK", []CarInfo{}))

	ws.Route(ws.POST("/info").To(c.findCar))
	// docs
	//Doc("get a car").
	//Param(ws.PathParameter("carid", "identifier of the car").DataType("string").DefaultValue("1")).
	//Metadata(restfulspec.KeyOpenAPITags, tags).
	//Writes(CarInfo{}). // on the response
	//Returns(200, "OK", CarInfo{}).
	//Returns(404, "Not Found", nil))

	ws.Route(ws.PUT("").To(c.createCar))
	// docs
	//Doc("create a car").
	//Metadata(restfulspec.KeyOpenAPITags, tags).
	//Reads(CarInfo{})) // from the request

	ws.Route(ws.DELETE("/remove/{carid}").To(c.deleteCar))
	// docs
	//Doc("delete a car").
	//Metadata(restfulspec.KeyOpenAPITags, tags).
	//Param(ws.PathParameter("carid", "identifier of the car").DataType("string")))

	return ws
}

func (c *CarInfoResource) createCar(request *restful.Request, response *restful.Response) {
	car := model.CarInfo{}
	err := request.ReadEntity(&car)
	//car.CreateAt = time.Now().Format("2006-01-02 15:04:05")
	if err != nil {
		response.WriteError(http.StatusInternalServerError, err)
	}
	clientHander, err := c.ClientPool.Get()
	if err != nil {
		response.WriteError(http.StatusInternalServerError, err)
	}
	cli := clientHander.(fins.Client)
	err = c.Operate.WriteCarOperate(cli, car.Info, car.OperateType)
	if err != nil {
		_, err := c.SqlDB.Insert(car)
		if err != nil {
			response.WriteError(http.StatusInternalServerError, err)
		}
		response.WriteError(http.StatusInternalServerError, err)
	}
	response.WriteHeader(http.StatusOK)
}

func (c CarInfoResource) findAllCars(request *restful.Request, response *restful.Response) {
	listCars := []model.CarInfo{}
	err := c.SqlDB.Find(&listCars)
	if err != nil {
		response.WriteError(http.StatusInternalServerError, err)
	}
	response.WriteEntity(listCars)

}

func (c CarInfoResource) findCar(request *restful.Request, response *restful.Response) {
	carinfo := new(model.CarInfo)
	err := request.ReadEntity(&carinfo)
	if err == nil {
		clientHander, err := c.ClientPool.Get()
		if err != nil {
			response.WriteError(http.StatusInternalServerError, err)
		}
		cli := clientHander.(fins.Client)
		result, err := c.Operate.ReadCarOperate(cli, carinfo.Info, carinfo.OperateType)
		if err != nil {
			response.WriteErrorString(http.StatusNotFound, "not fount car")
		} else if result == false {
			response.WriteErrorString(http.StatusNotFound, "not operate "+carinfo.OperateType)
		} else {
			response.WriteEntity(result)
		}
	} else {
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusInternalServerError, err.Error())
	}

}

func (c *CarInfoResource) deleteCar(request *restful.Request, response *restful.Response) {
	id := request.PathParameter("carid")
	var carinfo = new(model.CarInfo)
	_, err := c.SqlDB.Where("car_id = ?", id).Get(carinfo)
	if err != nil {
		response.WriteErrorString(http.StatusNotFound, "not fount car")
	}
	_, err = c.SqlDB.Where("car_id = ?", id).Delete(carinfo)
	if err != nil {
		response.WriteErrorString(http.StatusExpectationFailed, "marker delete failed")
	}

	_, err = c.SqlDB.Where("car_id = ?", id).Unscoped().Get(carinfo)
	if err != nil {
		response.WriteErrorString(http.StatusNotFound, "not fount car")
	}
	_, err = c.SqlDB.Where("car_id = ?", id).Unscoped().Delete(carinfo)
	if err != nil {
		response.WriteErrorString(http.StatusExpectationFailed, "delete failed")
	}
	response.WriteHeader(http.StatusOK)

}
