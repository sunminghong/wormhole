/*=============================================================================
#     FileName: wormhole.go
#       Author: sunminghong, allen.fantasy@gmail.com, http://weibo.com/5d13
#         Team: http://1201.us
#   LastChange: 2015-08-15 10:45:48
#      History:
=============================================================================*/


package wormhole

// Connection  
type Wormhole struct {
    ctrlConnection IConnection

    dataConnection IConnection

    guin TID
    fromAgentId int
    manager IWormholeManager

    receivePacketCallback ReceivePacketFunc

    closeCallback CommonCallbackFunc
}


// new Transport object
func NewWormhole(guin TID,fromAgentId int, manager IWormholeManager) *Wormhole {
    wh := &Wormhole {
        guin:           guin,
        fromAgentId:    fromAgentId,
        manager:        manager,
    }

    return wh
}


func (wh *Wormhole) SetReceivePacketCallback(cf ReceivePacketFunc)  {
    wh.receivePacketCallback = cf
}


func (wh *Wormhole) SetCloseCallback(cf CommonCallbackFunc) {
    wh.closeCallback = cf
}


func (wh *Wormhole) GetFromAgentId() int {
    return wh.fromAgentId
}


func (wh *Wormhole) GetGuin() TID {
    return wh.guin
}


func (wh *Wormhole) AddConnection(conn IConnection, t EConnType) {
    conn.SetReceivePacketCallback(wh.receivePacketCallback)

    if t == ECONN_TYPE_CTRL {
        if wh.ctrlConnection != nil {
            wh.dataConnection = wh.ctrlConnection
            conn.SetCloseCallback(wh.dataClosed)
        }

        if wh.dataConnection == nil {
            wh.dataConnection = conn
            conn.SetCloseCallback(wh.dataClosed)
        }

        wh.ctrlConnection = conn
        conn.SetCloseCallback(wh.ctrlClosed)
    } else {
        if wh.ctrlConnection == nil {
            wh.ctrlConnection = conn
            conn.SetCloseCallback(wh.ctrlClosed)
        }

        wh.dataConnection = conn
        conn.SetCloseCallback(wh.dataClosed)
    }
}


func (wh *Wormhole) dataClosed(id TID) {
    wh.dataConnection = nil
}


func (wh *Wormhole) ctrlClosed(id TID) {
    wh.dataConnection = nil
    wh.ctrlConnection = nil
    wh.closeCallback(wh.guin)
}


func (wh *Wormhole) Send(packet *RoutePacket) {
    if wh.fromAgentId == 0 {
        packet.Type = EPACKET_TYPE_GENERAL
    }
    wh.dataConnection.Send(packet)
}


func (wh *Wormhole) Broadcast(packet *RoutePacket) {
    wh.manager.Broadcast(packet)
}


func (wh *Wormhole) GetManager() IWormholeManager {
    return wh.manager
}


func (wh *Wormhole) Close() {
    wh.dataConnection.Close()
    wh.ctrlConnection.Close()
}


