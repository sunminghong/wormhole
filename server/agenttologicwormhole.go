/*=============================================================================
#     FileName: agenttologicwormhole.go
#       Author: sunminghong, allen.fantasy@gmail.com, http://weibo.com/5d13
#         Team: http://1201.us
#   LastChange: 2015-08-19 15:36:21
#      History:
=============================================================================*/


/*
agent to logic
*/

package server

import (
    . "wormhole/wormhole"

    gts "github.com/sunminghong/gotools"
)

type AgentToLogicWormhole struct {
    *Wormhole
}


func NewAgentToLogicWormhole(guin int, manager IWormholeManager, routepack IRoutePack) IWormhole { //*AgentToLogicWormhole {
    aw := &AgentToLogicWormhole{
        Wormhole : NewWormhole(guin, manager, routepack),
    }

    aw.RegisterSub(aw)
    aw.SetCloseCallback(aw.closed)
    return aw
}


func (alw *AgentToLogicWormhole) Init() {
    gts.Trace("agenttologicwormhole init()")
}


func (alw *AgentToLogicWormhole) ProcessPacket(dp *RoutePacket) {
    gts.Trace("agenttologicwormhole processpack receive:\n%q",dp)

    if server, ok := alw.GetManager().GetServer().(*Agent);ok {
        if acw, ok := server.ClientWormholes.Get(dp.Guin); ok {
            //if aw, ok := acw.(*AgentToClientWormhole);ok {
            aw := acw.(*AgentToClientWormhole)
            aw.Send(dp.Guin, dp.Data)
        }
    }
}


func (alw *AgentToLogicWormhole) ProcessPackets(dps []*RoutePacket) {
    gts.Trace("agenttologicwormhole processpackets receive %d route packets", len(dps))
    for _, dp := range dps {
        if dp.Type == EPACKET_TYPE_LOGIC_REGISTER {
            lm := alw.GetManager().(*LogicManager)
            lm.Register(dp.Data, alw)
        } else {
            alw.ProcessPacket(dp)
        }
    }
}


func (alw *AgentToLogicWormhole) closed(guin int) {
    gts.Trace("agent to logic wormhole closed")

    lm := alw.GetManager().(*LogicManager)
    lm.Remove(alw.GetGuin())
}


