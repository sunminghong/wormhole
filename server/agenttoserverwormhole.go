/*=============================================================================
#     FileName: agenttoserverwormhole.go
#       Author: sunminghong, allen.fantasy@gmail.com, http://weibo.com/5d13
#         Team: http://1201.us
#   LastChange: 2015-08-19 09:46:44
#      History:
=============================================================================*/


/*
gameserver connect to agent
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

type AgentToServerWormhole struct {
    *Wormhole
}

func NewAgentToServerWormhole(guin TID, manager IWormholeManager, routepack IRoutePack) *AgentToServerWormhole {
    aw := &AgentToServerWormhole{
        Wormhole : NewWormhole(guin, manager, routepack),
    }

    return aw
}

func (aw *AgentToServerWormhole) ProcessPackets(dps []*RoutePacket) {
    print("agentwormhole processpackets receive %d route packets", len(dps))
}

func (aw *AgentToServerWormhole) Close() {
}
