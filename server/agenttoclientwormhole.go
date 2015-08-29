/*=============================================================================
#     FileName: agent.go
#       Author: sunminghong, allen.fantasy@gmail.com, http://weibo.com/5d13
#         Team: http://1201.us
#   LastChange: 2015-08-18 19:19:34
#      History:
=============================================================================*/


/*
agent 代理连接服务器，接受玩家客户端连接
*/

package server


import (
    gts "github.com/sunminghong/gotools"

    . "wormhole/wormhole"
)


type AgentToClientWormhole struct {
    *Wormhole
}


func NewAgentToClientWormhole(guin int, manager IWormholeManager, routepack IRoutePack) IWormhole { //*AgentToClientWormhole {
    acw := &AgentToClientWormhole {
        Wormhole : NewWormhole(guin, manager, routepack),
    }

    acw.Inherit.RegisterSub(acw, "ProcessPackets")

    acw.SetCloseCallback(acw.closed)

    return acw
}


func (acw *AgentToClientWormhole) closed(guin int) {
    gts.Trace("agent to client closed")
    //TODO: send closed to logic
}


func (acw *AgentToClientWormhole) Init() {
    //gts.Trace("agenttoclient wormhole init()")
}


func (acw *AgentToClientWormhole) ProcessPackets(dps []*RoutePacket) {
    gts.Trace("agenttoclientwormhole processpackets receive %d packets", len(dps))

    for _,dp := range dps {
        gts.Trace("%q", dp)
        gts.Trace("guin:",acw.GetGuin())

        dp.Guin = acw.GetGuin()
        dp.Type = dp.Type | 1

        acw.SendPacket(dp)

        //转发给logic server
        //根据guin进行hash运算非配到相应的logic server
        if server, ok := acw.GetManager().GetServer().(*Agent);ok {
            server.LogicWormholes.(*LogicManager).Delay(dp)
        }
    }
}


