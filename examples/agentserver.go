/*=============================================================================
#     FileName: agentserver.go
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


type AgentToLogic struct {
    *server.AgentToLogicWormhole
}


func NewAgentToLogic(guin int, manager wormhole.IWormholeManager, routepack wormhole.IRoutePack) wormhole.IWormhole {
    aw := &AgentToLogic {
        AgentToLogicWormhole: server.NewAgentToLogicWormhole(guin, manager, routepack),
    }
    aw.RegisterSub(aw, "ProcessPacket")

    return aw
}


func (aw *AgentToLogic) ProcessPacket(dp *wormhole.RoutePacket) {
    gts.Trace("agenttologic processpacket receive %d route packets",dp)
}
//AgentToLogic end


type AgentToClient struct {
    *server.AgentToClientWormhole
}


func NewAgentToClient(guin int, manager wormhole.IWormholeManager, routepack wormhole.IRoutePack) wormhole.IWormhole {
    aw := &AgentToClient {
        AgentToClientWormhole: server.NewAgentToClientWormhole(guin, manager, routepack),
    }
    aw.RegisterSub(aw)

    return aw
}


func (aw *AgentToClient) ProcessPackets(dps []*wormhole.RoutePacket) {
    gts.Trace("agenttoclientwormhole processpackets receive %d route packets",len(dps))
    gts.Trace("ProcessPackets:::sldhrq8903246dsfq238946204\n%q", dps)
}
//AgentToClient end


var agent *server.Agent


var (
    agentConf=flag.String("agentConf","agent1.conf","agent server config file")
)


func main() {
    runtime.GOMAXPROCS(runtime.NumCPU())

    flag.Parse()

    c, err := iniconfig.ReadConfigFile(*agentConf)
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

    cwormholes := wormhole.NewWormholeManager(routepack, wormhole.MAX_CONNECTIONS,wormhole.EWORMHOLE_TYPE_CLIENT)

    dispatcher := server.NewDispatcher(routepack)
    lwormholes := server.NewLogicManager(routepack, dispatcher)
    agent := server.NewAgentFromIni(c,routepack,cwormholes,lwormholes,NewAgentToClient,NewAgentToLogic)

    agent.Start()

    gts.Info("----------------agent server is running!-----------------",agent.GetServerId())

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
