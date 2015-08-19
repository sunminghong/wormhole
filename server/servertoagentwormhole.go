/*=============================================================================
#     FileName: servertoagentwormhole.go
#       Author: sunminghong, allen.fantasy@gmail.com, http://weibo.com/5d13
#         Team: http://1201.us
#   LastChange: 2015-08-19 09:42:43
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

type ServerToAgentWormhole struct {
    *Wormhole
}

func NewServerToAgentWormhole(guin TID, manager IWormholeManager, routepack IRoutePack) *ServerToAgentWormhole {
    aw := &ServerToAgentWormhole {
        Wormhole : NewWormhole(guin, manager, routepack),
    }

    return aw
}

func (aw *ServerToAgentWormhole) ProcessPackets(dps []*RoutePacket) {
    print("agentwormhole processpackets receive %d route packets", len(dps))
}

func (aw *ServerToAgentWormhole) Close() {
}
