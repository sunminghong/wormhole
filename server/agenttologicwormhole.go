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
    //"net"
    //"strconv"
    //"time"
    //"fmt"
    //"strings"

    . "wormhole/wormhole"
)

type AgentToLogicWormhole struct {
    *Wormhole
}


func NewAgentToLogicWormhole(guin int, manager IWormholeManager, routepack IRoutePack) *AgentToLogicWormhole {
    aw := &AgentToLogicWormhole{
        Wormhole : NewWormhole(guin, manager, routepack),
    }

    return aw
}


func (alw *AgentToLogicWormhole) Init() {
}


func (alw *AgentToLogicWormhole) ProcessPackets(dps []*RoutePacket) {
    print("agentwormhole processpackets receive %d route packets", len(dps))
    for _, dp := range dps {
        if dp.Type == EPACKET_TYPE_LOGIC_REGISTER {
            lm := alw.GetManager().(*LogicManager)
            lm.Register(dp.Data, alw)
        }
    }
}


func (alw *AgentToLogicWormhole) Close() {
}


