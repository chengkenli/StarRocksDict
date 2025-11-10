package main

import (
	_ "StarRocksDict/init"
	"StarRocksDict/run"
	"StarRocksDict/util"
)

func main() {
	util.Parms()
	run.Run()
}
