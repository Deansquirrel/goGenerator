package common

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/Deansquirrel/go-tool"
	"github.com/Deansquirrel/goGenerator/global"
	"github.com/Deansquirrel/goGenerator/object"
	"time"
)

func PrintAndLog(msg string){
	if global.SysConfig.TotalConfig.IsDebug {
		fmt.Println(msg)
	} else {
		err := go_tool.Log(msg)
		if err != nil {
			fmt.Println(err.Error())
		}
	}
}

func GetSysConfig() error{
	//==================================================================================================================
	path,err := go_tool.GetCurrPath()
	if err != nil {
		return err
	}
	var config object.SysConfig
	_,err = toml.DecodeFile(path + "\\" + "config.toml",&config)
	if err != nil {
		return err
	}
	global.SysConfig = config
	//==================================================================================================================
	global.TimeoutNS = time.Duration(int64(global.SysConfig.GeneratorConfig.TimeoutNS) * 1000 * 1000 * 1000)
	global.DurationNS = time.Duration(int64(global.SysConfig.GeneratorConfig.DurationNS) * 1000 * 1000 * 1000)
	//==================================================================================================================
	return nil
}