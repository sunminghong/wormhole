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

    broadcastChan chan *RoutePacket
    fromType  EWormholeType
    routePack IRoutePack
}


func NewWormholeManager(routepack IRoutePack, broadcast_chan_num int, fromType EWormholeType) *WormholeManager {
    wm := &WormholeManager {
        wmlock:        new(sync.RWMutex),
        wormholes:      make(map[TID]IWormhole),
        fromType:       fromType,
        routePack:      routepack,
    }

    wm.broadcastChan = make(chan *RoutePacket, broadcast_chan_num)
    go wm.broadcastHandler(wm.broadcastChan)

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

    if _, ok := wm.wormholes[guin];ok {
        delete(wm.wormholes, guin)
    }
}


func (wm *WormholeManager) Close(guin TID) {
    wm.wmlock.Lock()
    defer wm.wmlock.Unlock()

    if wh, ok := wm.wormholes[guin];ok {
        wh.Close()
        delete(wm.wormholes, guin)
    }
}


func (wm *WormholeManager) CloseAll() {
    wm.wmlock.Lock()
    defer wm.wmlock.Unlock()

    for _, wh := range wm.wormholes {
        wh.Close()
    }
}


func (wm *WormholeManager) Send(guin TID, data []byte) {
    if wh, ok := wm.wormholes[guin];ok {
        /*
        packet := &RoutePacket {
            Type:   EPACKET_TYPE_GENERAL,
            Guin:   guin,
            Data:   data,
        }
        if wm.fromType == EWORMHOLE_TYPE_AGENT {
            packet.Type = EPACKET_TYPE_DELAY
        }
        wh.Send(packet)
        */
        wh.Send(guin, data)
    }
}


func (wm *WormholeManager) Broadcast(guin TID, data []byte) {
    packet := &RoutePacket {
        Type:   EPACKET_TYPE_BROADCAST,
        Guin:   guin,
        Data:   data,
    }

    wm.broadcastChan <- packet
}


func (wm *WormholeManager) broadcastHandler(broadcastChan <-chan *RoutePacket) {
    for {
        packet := <-broadcastChan
        data := wm.routePack.Pack(packet)

        for _, wh := range wm.wormholes {
            wh.SendRaw(data)
        }
    }
}



func (wm *WormholeManager) Length() int {
    return len(wm.wormholes)
}

