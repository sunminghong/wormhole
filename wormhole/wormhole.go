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
    wormholes IWormholeManager

    fromId int
    fromType EWormholeType

    state EWormholeState

    routePack IRoutePack
    //receivePacketCallback ReceivePacketFunc

    closeCallback CommonCallbackFunc
}


func NewWormhole(guin TID, wormholes IWormholeManager, routepack IRoutePack) *Wormhole {
    wh := &Wormhole {
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

//需要继承实现具体的处理逻辑
func (c *Wormhole) Init() {
    print("wormhole don't implent init()")
}


func (c *Wormhole) ProcessPackets(dps []*RoutePacket) {
    print("wormhole don't implentpent processpackets()")
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


func (wh *Wormhole) GetGuin() TID {
    return wh.guin
}


func (wh *Wormhole) GetState() EWormholeState {
    return wh.state
}


func (wh *Wormhole) SetState(state EWormholeState) {
    wh.state = state
}


func (wh *Wormhole) AddConnection(conn IConnection, t EConnType) {
    //conn.SetReceivePacketCallback(wh.receivePacketCallback)
    conn.SetReceiveCallback(wh.receiveBytes)

    if t == ECONN_TYPE_CTRL {
        //if wh.ctrlConnection != nil {
            //wh.dataConnection = wh.ctrlConnection
            //conn.SetCloseCallback(wh.dataClosed)
        //}

        //if wh.dataConnection == nil {
            //wh.dataConnection = conn
            //conn.SetCloseCallback(wh.dataClosed)
        //}

        wh.dataConnection = wh.ctrlConnection
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


func (wh *Wormhole) receiveBytes(conn IConnection) {
    n, dps := wh.routePack.Fetch(conn)
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
    bytes := wh.routePack.Pack(packet)
    wh.dataConnection.Send(bytes)
}


/*
func (wh *Wormhole) Broadcast(packet *RoutePacket) {
    wh.wormholes.Broadcast(packet)
}
*/


func (wh *Wormhole) SendRaw(raw []byte) {
    wh.dataConnection.Send(raw)
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

    bytes := wh.routePack.Pack(packet)
    wh.dataConnection.Send(bytes)
}


func (wh *Wormhole) Broadcast(guin TID, data []byte) {
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


