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
    //"net"
    //"strconv"
    //"time"
    //"fmt"
    //"strings"

    . "wormhole/wormhole"
)

type AgentToClientWormhole struct {
    *Wormhole
}

func NewAgentToClientWormhole(guin TID, manager IWormholeManager, routepack IRoutePack) *AgentToClientWormhole {
    aw := &AgentToClientWormhole {
        Wormhole : NewWormhole(guin, manager, routepack),
    }

    return aw
}

func (aw *AgentToClientWormhole) ProcessPackets(dps []*RoutePacket) {
    print("agentwormhole processpackets receive %d route packets", len(dps))
}

func (aw *AgentToClientWormhole) Close() {
}
