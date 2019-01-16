package main

import (
	"github.com/Deansquirrel/goGenerator/common"
	"github.com/Deansquirrel/goGenerator/global"
)

func main() {
	//==================================================================================================================
	config, err := common.GetSysConfig("config.toml")
	if err != nil {
		common.PrintAndLog("加载配置文件时遇到错误：" + err.Error())
		return
	}
	global.SysConfig = config
	err = common.RefreshConfig(*global.SysConfig)
	if err != nil {
		common.PrintAndLog("刷新配置时遇到错误：" + err.Error())
		return
	}
	//==================================================================================================================
	common.PrintOrLog("程序启动[Generator]")
	defer common.PrintOrLog("程序退出[Generator]")
	//==================================================================================================================
}
