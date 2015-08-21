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
    gts "github.com/sunminghong/gotools"

    . "wormhole/wormhole"
)

type LogicToAgentWormhole struct {
    *Wormhole

    group string
}

func NewLogicToAgentWormhole(guin int, manager IWormholeManager, routepack IRoutePack) *LogicToAgentWormhole {
    agentId, _,_ := ParseGuin(guin)

    aw := &LogicToAgentWormhole {
        Wormhole : NewWormhole(agentId, manager, routepack),
        group : "",
    }

    return aw
}


func (aw *LogicToAgentWormhole) SetGroup(group string) {
    aw.group = group

    gts.Trace("logic register to agent:group(%s)", group)
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
    gts.Trace("agentwormhole processpackets receive %d route packets", len(dps))
}


func (aw *LogicToAgentWormhole) Close() {
    aw.Wormhole.Close()
}


