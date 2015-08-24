/*=============================================================================
#     FileName: logictoagent.go
#       Author: sunminghong, allen.fantasy@gmail.com, http://weibo.com/5d13
#         Team: http://1201.us
#   LastChange: 2015-08-21 12:36:08
#      History:
=============================================================================*/


/*
定义一个基本的逻辑服务器端

*/
package examples


import (
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


func (aw *LogicToAgent) ProcessPackets(dps []*wormhole.RoutePacket) {
    gts.Trace("agentwormhole processpackets receive %d route packets",len(dps))
    gts.Trace("ProcessPackets:\n%q", dps)
}


