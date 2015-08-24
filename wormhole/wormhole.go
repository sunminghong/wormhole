/*=============================================================================
#     FileName: wormhole.go
#       Author: sunminghong, allen.fantasy@gmail.com, http://weibo.com/5d13
#         Team: http://1201.us
#   LastChange: 2015-08-15 10:45:48
#      History:
=============================================================================*/


package wormhole



import (
    gts "github.com/sunminghong/gotools"
)


// Connection  
type Wormhole struct {
    *Inherit

    ctrlConnection IConnection
    dataConnection IConnection
    sendConnection IConnection

    guin int
    wormholes IWormholeManager

    fromId int
    fromType EWormholeType

    state EWormholeState

    routePack IRoutePack
    //receivePacketCallback ReceivePacketFunc

    closeCallback CommonCallbackFunc
}


func NewWormhole(guin int, wormholes IWormholeManager, routepack IRoutePack) *Wormhole {
    wh := &Wormhole {
        Inherit:        NewInherit("ProcessPackets"),
        guin:           guin,
        fromId:         0,
        fromType:       EWORMHOLE_TYPE_CLIENT,
        wormholes:        wormholes,
        routePack:      routepack,
    }

    return wh
}


/*
func (wh *Wormhole) SetReceivePacketCallback(cf ReceivePacketFunc)  {
    wh.receivePacketCallback = cf
}
*/


func (c *Wormhole) ProcessPackets(dps []*RoutePacket) {
    c.Inherit.CallSub("ProcessPackets", dps)
}

//需要继承实现具体的处理逻辑
func (c *Wormhole) Init() {
    c.Inherit.CallSub("Init")
}


func (wh *Wormhole) SetCloseCallback(cf CommonCallbackFunc) {
    wh.closeCallback = cf
}


func (wh *Wormhole) GetFromType() EWormholeType {
    return wh.fromType
}


func (wh *Wormhole) SetFromType(t EWormholeType) {
    wh.fromType = t
}


func (wh *Wormhole) SetFromId(id int) {
    wh.fromId = id
}


func (wh *Wormhole) GetFromId() int {
    return wh.fromId
}


func (wh *Wormhole) GetGuin() int {
    return wh.guin
}


func (wh *Wormhole) GetState() EWormholeState {
    return wh.state
}


func (wh *Wormhole) SetState(state EWormholeState) {
    wh.state = state
}


func (wh *Wormhole) AddConnection(conn IConnection, t EConnType) {
    gts.Trace("addConnection:", t)
    conn.SetReceiveCallback(wh.receiveBytes)

    if t == ECONN_TYPE_CTRL {
        if wh.ctrlConnection == nil {
            wh.ctrlConnection = conn
            conn.SetCloseCallback(wh.ctrlClosed)

            wh.sendConnection = conn
            return
        }

        if wh.dataConnection == nil {
            wh.sendConnection = conn
        } else {
            wh.dataConnection.SetType(ECONN_TYPE_DATA)
            wh.dataConnection.SetCloseCallback(wh.dataClosed)
        }
    } else {
        wh.dataConnection = conn

        if wh.ctrlConnection == nil {
            wh.ctrlConnection = conn
            conn.SetType(ECONN_TYPE_CTRL)
            conn.SetCloseCallback(wh.ctrlClosed)
        } else {
            conn.SetType(t)
            conn.SetCloseCallback(wh.dataClosed)
        }

        wh.sendConnection = conn
    }
}


func (wh *Wormhole) receiveBytes(conn IConnection) {
    gts.Trace("wormhole receiveBytes:% X", conn.GetBuffer().Stream.Bytes())
    gts.Trace("wormhole receiveBytes:%q", conn.GetBuffer().Stream.Bytes())
    n, dps := wh.routePack.Fetch(conn.GetBuffer())
    if n > 0 {
        //wh.receivePacketCallback(wh, dps)
        wh.ProcessPackets(dps)
    }
}


func (wh *Wormhole) dataClosed(id int) {
    wh.dataConnection = nil
}


func (wh *Wormhole) ctrlClosed(id int) {
    wh.dataConnection = nil
    wh.ctrlConnection = nil

    wh.closeCallback(int(wh.guin))
}


func (wh *Wormhole) SendPacket(packet *RoutePacket) {
    gts.Trace("sendprotocol type:%d,", wh.sendConnection.GetProtocolType())
    gts.Trace(wh.fromType)
    bytes := wh.routePack.Pack(packet)
}


/*
func (wh *Wormhole) Broadcast(packet *RoutePacket) {
    wh.wormholes.Broadcast(packet)
}
*/


func (wh *Wormhole) SendRaw(raw []byte) {
    gts.Trace("sendprotocol type:%d,", wh.sendConnection.GetProtocolType())
    gts.Trace(wh.fromType)
    wh.ctrlConnection.Send(raw)
}


func (wh *Wormhole) Send(guin int, data []byte) {
    packet := &RoutePacket {
        Type:   EPACKET_TYPE_GENERAL,
        Guin:   guin,
        Data:   data,
    }
    if wh.fromType == EWORMHOLE_TYPE_AGENT {
        packet.Type = EPACKET_TYPE_DELAY
    }
    packet.Type = packet.Type + 1

    bytes := wh.routePack.Pack(packet)
    gts.Trace("sendprotocol type:%d,", wh.sendConnection.GetProtocolType())
    gts.Trace(wh.fromType)
    wh.ctrlConnection.Send(bytes)
}


func (wh *Wormhole) Broadcast(guin int, data []byte) {
    //packet := &RoutePacket {
        //Type:   EPACKET_TYPE_BROADCAST,
        //Guin:   guin,
        //Data:   data,
    //}

    wh.wormholes.Broadcast(guin, data)
}


func (wh *Wormhole) GetManager() IWormholeManager {
    return wh.wormholes
}


func (wh *Wormhole) Close() {
    wh.dataConnection.Close()
    wh.ctrlConnection.Close()
}


