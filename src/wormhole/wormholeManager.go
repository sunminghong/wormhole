/*=============================================================================
#     FileName: wormholeManager.go
#       Author: sunminghong, allen.fantasy@gmail.com, http://weibo.com/5d13
#         Team: http://1201.us
#   LastChange: 2015-08-15 10:49:21
#      History:
=============================================================================*/


package wormhole


import (
    "sync"
)

type WormholeManager struct {
    wormholes map[TID]IWormhole

    wmlock *sync.RWMutex
}


func NewWormholeManager() *WormholeManager {
    wm := &WormholeManager {
        wmlock:        new(sync.RWMutex),
        wormholes:      make(map[TID]IWormhole),
    }

    return wm
}


func (wm *WormholeManager) Add(wh IWormhole) {
    wm.wmlock.Lock()
    defer wm.wmlock.Unlock()

    wm.wormholes[wh.GetGuin()] = wh
}


func (wm *WormholeManager) Remove(guin TID)  {
    wm.wmlock.Lock()
    defer wm.wmlock.Unlock()

    delete(wm.wormholes, guin)
}



func (wm *WormholeManager) Send(guin TID, data []byte) {
    if wh, ok := wm.wormholes[guin];ok {
        packet := &RoutePacket {
            Type:   EPACKET_TYPE_DELAY,
            Guin:   guin,
            Data:   data,
        }
        wh.Send(packet)
    }
}


func (wm *WormholeManager) Broadcast(guin TID, data []byte) {
    packet := &RoutePacket {
        Type:   EPACKET_TYPE_BROADCAST,
        Guin:   guin,
        Data:   data,
    }

    wm.wmlock.Lock()
    defer wm.wmlock.Unlock()

    for _, wh := range wm.wormholes {
        wh.Send(packet)
    }
}


func (wm *WormholeManager) Close(guin TID) {
    if wh, ok := wm.wormholes[guin];ok {
        wh.Close()
    }

    wm.Remove(guin)
}


func (wm *WormholeManager) CloseAll() {
    wm.wmlock.Lock()
    defer wm.wmlock.Unlock()

    for _, wh := range wm.wormholes {
        wh.Close()
    }
}


