/*=============================================================================
#     FileName: agentmanager.go
#       Author: sunminghong, allen.fantasy@gmail.com, http://weibo.com/5d13
#         Team: http://1201.us
#   LastChange: 2015-08-20 18:13:28
#      History:
=============================================================================*/


package server


import (
    "sync"

    gts "github.com/sunminghong/gotools"

    . "wormhole/wormhole"
)


type AgentManager struct {
    *WormholeManager

    wmlock *sync.RWMutex
}


func NewAgentManager(routepack IRoutePack) *AgentManager {
    wm := &AgentManager {
        WormholeManager: NewWormholeManager(routepack, 100, EWORMHOLE_TYPE_AGENT),
    }
    return wm
}


func (wm *AgentManager) Send(guin int, data []byte) {
    agentId, _,_ := ParseGuin(guin)
    wh, ok := wm.Get(agentId)
    if ok {
        wh.Send(guin, data)
        return
    }

    gts.Error("agentmanager:guin don't find wormhole:%d", guin)
}

