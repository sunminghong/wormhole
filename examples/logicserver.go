/*=============================================================================
#     FileName: logicserver.go
#       Author: sunminghong, allen.fantasy@gmail.com, http://weibo.com/5d13
#         Team: http://1201.us
#   LastChange: 2015-08-21 12:35:56
#      History:
=============================================================================*/


/*
定义一个基本的逻辑服务器端

*/
package main


import (
	"bufio"
	"os"
    "flag"
    "runtime"

    iniconfig "github.com/sunminghong/iniconfig"
    gts "github.com/sunminghong/gotools"

    "wormhole/wormhole"
    "wormhole/server"

)


type LogicToAgent struct {
    *server.LogicToAgentWormhole
}


func NewLogicToAgent(guin int, manager wormhole.IWormholeManager, routepack wormhole.IRoutePack) wormhole.IWormhole {
    aw := &LogicToAgent {
        LogicToAgentWormhole: server.NewLogicToAgentWormhole(guin, manager, routepack),
    }
    aw.RegisterSub(aw)

    return aw
}


func (aw *LogicToAgent) Init() {
    gts.Trace("logicToAgent Init")
    aw.LogicToAgentWormhole.Init()
}


func (aw *LogicToAgent) ProcessPackets(dps []*wormhole.RoutePacket) {
    gts.Trace("logicToAgent processpackets receive %d route packets",len(dps))
}
//logicToAgent end




var logic *server.Logic


var (
    logicConf=flag.String("logicConf","logic1.conf","logic server config file")
)


func main() {
    runtime.GOMAXPROCS(runtime.NumCPU())

    flag.Parse()

    c, err := iniconfig.ReadConfigFile(*logicConf)
    if err != nil {
        gts.Error(err.Error())
        return
    }

    var endianer gts.IEndianer
    section := "Default"
    endian, err := c.GetInt(section, "endian")
    if err == nil {
        endianer = gts.GetEndianer(endian)
    } else {
        endianer = gts.GetEndianer(gts.LittleEndian)
    }
    routepack := wormhole.NewRoutePack(endianer)
    wormholes := server.NewAgentManager(routepack)

    logic := server.NewLogicFromIni(c,routepack,wormholes,NewLogicToAgent)

    logic.ConnectFromIni(c)

    gts.Info("----------------logic server is running!-----------------",logic.GetServerId())

    quit := make(chan bool)
    go exit(quit)

    <-quit
}


func exit(quit chan bool) {
    for {
        r := bufio.NewReader(os.Stdin)
        ru, _, _ := r.ReadRune()
        print(ru)
        if ru == 115 {
            quit <- true
            return
        }
    }
}
