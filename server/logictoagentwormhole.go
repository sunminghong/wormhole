/*=============================================================================
#     FileName: servertoagentwormhole.go
#       Author: sunminghong, allen.fantasy@gmail.com, http://weibo.com/5d13
#         Team: http://1201.us
#   LastChange: 2015-08-19 09:42:43
#      History:
=============================================================================*/


/*
logic server connect to agent
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

type LogicToAgentWormhole struct {
    *Wormhole

    group string
}

func NewLogicToAgentWormhole(guin int, manager IWormholeManager, routepack IRoutePack) *LogicToAgentWormhole {
    aw := &LogicToAgentWormhole {
        Wormhole : NewWormhole(guin, manager, routepack),
        group : "",
    }

    return aw
}


func (aw *LogicToAgentWormhole) SetGroup(group string) {
    aw.group = group

    packet := &RoutePacket {
        Type:   EPACKET_TYPE_LOGIC_REGISTER,
        Guin:   0,
        Data:   []byte(group),
    }
    aw.SendPacket(packet)
}


func (aw *LogicToAgentWormhole) Init() {
    print("logictoagentwormhole is init()")
}


func (aw *LogicToAgentWormhole) ProcessPackets(dps []*RoutePacket) {
    print("agentwormhole processpackets receive %d route packets", len(dps))
}


func (aw *LogicToAgentWormhole) Close() {
    aw.Wormhole.Close()
}


