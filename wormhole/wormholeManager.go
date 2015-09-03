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
    gts "github.com/sunminghong/gotools"
)

type WormholeManager struct {
    wormholes map[int]IWormhole

    wmlock *sync.RWMutex

    broadcastChan chan *RoutePacket
    fromType  EWormholeType
    routePack IRoutePack

    server IServer
}


func NewWormholeManager(routepack IRoutePack, broadcast_chan_num int, fromType EWormholeType) *WormholeManager {
    wm := &WormholeManager {
        wmlock:        new(sync.RWMutex),
        wormholes:      make(map[int]IWormhole),
        fromType:       fromType,
        routePack:      routepack,
        server:         nil,
    }

    wm.broadcastChan = make(chan *RoutePacket, broadcast_chan_num)
    go wm.broadcastHandler(wm.broadcastChan)

    return wm
}


func (wh *WormholeManager) SetServer(server IServer) {
    wh.server = server
}


func (wh *WormholeManager) GetServer() IServer {
    return wh.server
}


func (wm *WormholeManager) Add(wh IWormhole) {
    wm.wmlock.Lock()
    defer wm.wmlock.Unlock()

    wm.wormholes[wh.GetGuin()] = wh
}


func (wm *WormholeManager) Get(guin int) (IWormhole,bool) {
    wh, ok := wm.wormholes[guin]
    return wh, ok
}


func (wm *WormholeManager) Remove(guin int)  {
    wm.wmlock.Lock()
    defer wm.wmlock.Unlock()

    print("wormholemanager remove")

    gts.Trace("wormholemanager remove")
    if _, ok := wm.wormholes[guin];ok {
        delete(wm.wormholes, guin)
    }
}


func (wm *WormholeManager) Close(guin int) {
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


func (wm *WormholeManager) Send(guin int, method int, data []byte) {
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
        print("wormholemanager send")
        wh.Send(guin, method, data)
    }
}


func (wm *WormholeManager) Broadcast(guin int, method int, data []byte) {
    packet := &RoutePacket {
        Type:   EPACKET_TYPE_BROADCAST,
        Guin:   guin,
        Method: method,
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

