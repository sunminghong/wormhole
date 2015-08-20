/*=============================================================================
#     FileName: wormholeManager.go
#       Author: sunminghong, allen.fantasy@gmail.com, http://weibo.com/5d13
#         Team: http://1201.us
#   LastChange: 2015-08-15 10:49:21
#      History:
=============================================================================*/


package server


import (
    "sync"

    . "wormhole/wormhole"
)


type LogicManager struct {
    *WormholeManager

    wmlock *sync.RWMutex
    dispatcher ILogicDispatcher
}


func NewLogicManager(routepack IRoutePack, fromType EWormholeType, dispatcher ILogicDispatcher) *LogicManager {
    wm := &LogicManager {
        WormholeManager: NewWormholeManager(routepack, 100, fromType),
        dispatcher : dispatcher,
    }
    return wm
}


func (wm *LogicManager) Delay(dp *RoutePacket) {
    wm.wmlock.Lock()
    defer wm.wmlock.Unlock()

    wh := wm.dispatcher.Dispatch(dp)
    if wh != nil {
        wh.SendPacket(dp)
    }
}


func (wm *LogicManager) Register(rule []byte, wh IWormhole) {
    wm.wmlock.Lock()
    defer wm.wmlock.Unlock()

    wm.dispatcher.AddHandler(rule, wh)
}


