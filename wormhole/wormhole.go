/*=============================================================================
#     FileName: wormhole.go
#       Author: sunminghong, allen.fantasy@gmail.com, http://weibo.com/5d13
#         Team: http://1201.us
#   LastChange: 2015-08-15 10:45:48
#      History:
=============================================================================*/


package wormhole



import (
    //"reflect"

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

    closeded bool
}


func NewWormhole(guin int, wormholes IWormholeManager, routepack IRoutePack) *Wormhole {
    gts.Trace("new wormhole:")
    wh := &Wormhole {
        Inherit:        NewInherit("ProcessPackets"),
        guin:           guin,
        fromId:         0,
        fromType:       EWORMHOLE_TYPE_CLIENT,
        wormholes:      wormholes,
        routePack:      routepack,
        closeded:       false,
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
    gts.Trace("addConnection:", t, conn.GetProtocolType())
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
        if wh.ctrlConnection == nil {
            wh.ctrlConnection = conn

            conn.SetType(ECONN_TYPE_CTRL)
            conn.SetCloseCallback(wh.ctrlClosed)
        } else {
            gts.Trace("addConnection:dataclosed")

            wh.dataConnection = conn
            conn.SetType(t)
            conn.SetCloseCallback(wh.dataClosed)
        }

        wh.sendConnection = conn
    }
}


func (wh *Wormhole) receiveBytes(conn IConnection) {
    if len(conn.GetBuffer().Stream.Bytes()) < 7 {
        return
    }

    //print("-------------------------------------------------------\n")
    //gts.Trace("wormhole receiveBytes:% X", conn.GetBuffer().Stream.Bytes())
    //gts.Trace("wormhole receiveBytes:%q", conn.GetBuffer().Stream.Bytes())
    n, dps := wh.routePack.Fetch(conn.GetBuffer())
    if n > 0 {
        //wh.receivePacketCallback(wh, dps)
        wh.ProcessPackets(dps)
    }
}


func (wh *Wormhole) dataClosed(id int) {
    gts.Trace("dataClosed")

    if wh.sendConnection == nil {
        //wh.dataConnection.GetType() == wh.sendConnection.GetType() {
        wh.sendConnection = wh.ctrlConnection
    }
}


func (wh *Wormhole) ctrlClosed(id int) {
    gts.Trace("ctrlClosed")

    wh.closeded = true

    if wh.dataConnection != nil {
        wh.dataConnection.Close()
    } else {
        wh.dataConnection = nil
    }
    if wh.ctrlConnection != nil {
        wh.ctrlConnection.Close()
    } else {
        wh.ctrlConnection = nil
    }

    wh.sendConnection = nil

    gts.Trace("//////////////////////////////////////////")
    wh.closeCallback(wh.guin)
    gts.Trace("------------------------------------------")
}


func (wh *Wormhole) SendPacket(packet *RoutePacket) {
    bytes := wh.routePack.Pack(packet)
    wh.SendRaw(bytes)
}


func (wh *Wormhole) SendRaw(bytes []byte) {
    if wh.closeded {
        return
    }

    if len(bytes) <= 1024 {
    //if len(bytes) < 0 {
        gts.Trace("send send protocol type:%d,%d,%d",wh.GetGuin(), wh.sendConnection.GetType(), wh.sendConnection.GetProtocolType())
        wh.sendConnection.Send(bytes)

    } else {
        //gts.Trace("send ctrl protocol type:",wh.ctrlConnection)
        gts.Trace("send ctrl protocol type:%d,%d,%d",wh.GetGuin(), wh.ctrlConnection.GetType(), wh.ctrlConnection.GetProtocolType())
        gts.Trace("send ctrl:%q",bytes)

        //v := reflect.ValueOf(wh.ctrlConnection)
        //t := reflect.TypeOf(wh.ctrlConnection)

        //gts.Trace("Type:", t)
        //gts.Trace("Value:", v)
        //gts.Trace("Kind:", t.Kind())
        //gts.Trace("Kind:", wh.ctrlConnection.GetId())

        wh.ctrlConnection.Send(bytes)
    }
}


func (wh *Wormhole) Send(guin int, method int, data []byte) {
    packet := &RoutePacket {
        Type:   EPACKET_TYPE_GENERAL,
        Guin:   guin,
        Method: method,
        Data:   data,
    }
    gts.Trace("fromType:%d, %d", wh.fromType, EWORMHOLE_TYPE_AGENT)
    if wh.fromType == EWORMHOLE_TYPE_AGENT {
        packet.Type = EPACKET_TYPE_DELAY
    }

    bytes := wh.routePack.Pack(packet)
    wh.SendRaw(bytes)
}


/*
func (wh *Wormhole) Broadcast(packet *RoutePacket) {
    wh.wormholes.Broadcast(packet)
}
*/


func (wh *Wormhole) Broadcast(guin int, method int, data []byte) {
    //packet := &RoutePacket {
        //Type:   EPACKET_TYPE_BROADCAST,
        //Guin:   guin,
        //Data:   data,
    //}

    wh.wormholes.Broadcast(guin, method, data)
}


func (wh *Wormhole) GetManager() IWormholeManager {
    return wh.wormholes
}


func (wh *Wormhole) Close() {
    wh.dataConnection.Close()
    wh.ctrlConnection.Close()
}


