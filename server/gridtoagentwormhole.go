/*=============================================================================
#     FileName: gridtoagentwormhole.go
#       Author: sunminghong, allen.fantasy@gmail.com, http://weibo.com/5d13
#         Team: http://1201.us
#   LastChange: 2015-08-19 09:41:25
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

type GridToAgentWormhole struct {
    *Wormhole
}

func NewGridToAgentWormhole(guin TID, manager IWormholeManager, routepack IRoutePack) *GridToAgentWormhole {
    aw := &GridToAgentWormhole {
        Wormhole : NewWormhole(guin, manager, routepack),
    }

    return aw
}

func (aw *GridToAgentWormhole) ProcessPackets(dps []*RoutePacket) {
    print("agentwormhole processpackets receive %d route packets", len(dps))
}

func (aw *GridToAgentWormhole) Close() {
}
