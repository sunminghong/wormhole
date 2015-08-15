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
    manager IWormholeManager

    fromId int
    fromType EWormholeType

    receivePacketCallback ReceivePacketFunc

    closeCallback CommonCallbackFunc
}


// new Transport object
func NewWormhole(guin TID, manager IWormholeManager) *Wormhole {
    wh := &Wormhole {
        guin:           guin,
        fromId:         0,
        fromType:       EWORMHOLE_TYPE_CLIENT,
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


func (wh *Wormhole) GetType() EWormholeType {
    return wh.fromType
}


func (wh *Wormhole) SetType(t EWormholeType) {
    wh.fromType = t
}


func (wh *Wormhole) SetFromId(id int) {
    wh.fromId = id
}


func (wh *Wormhole) GetFromId() int {
    return wh.fromId
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
    wh.dataConnection.Send(packet)
}


func (wh *Wormhole) Send(guin TID, data []byte) {
    packet := &RoutePacket {
        Type:   EPACKET_TYPE_GENERAL,
        Guin:   guin,
        Data:   data,
    }
    if wh.fromType == EWORMHOLE_TYPE_AGENT {
        packet.Type = EPACKET_TYPE_DELAY
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


