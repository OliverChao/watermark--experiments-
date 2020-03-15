package main

import (
	"WaterMasking/model"
	"WaterMasking/service"
	"fmt"
)

func init() {
	model.FlagConfInit()
}

func main() {

	fmt.Println(model.Conf.FiledNames, model.Conf.Nu)

	service.ConnectDB()
	//now := time.Now()
	//data := service.Student.GetBatchData(0, 500)
	//fmt.Println(len(data))
	//fmt.Println(data[0].Name, data[0].ID)
	//
	//fmt.Println(time.Since(now))
	service.DestroyAll()
	service.DisconnectDB()
}
