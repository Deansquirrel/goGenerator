package main

import "github.com/Deansquirrel/goGenerator/common"

func main(){
	//==================================================================================================================
	err := common.GetSysConfig()
	if err != nil {
		common.PrintAndLog("获取配置时遇到错误：" + err.Error())
		return
	}
	//==================================================================================================================
	common.PrintAndLog("程序启动[Generator]")
	defer common.PrintAndLog("程序退出[Generator]")
	//==================================================================================================================
}
