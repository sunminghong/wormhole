/*=============================================================================
#     FileName: interface.go
#       Author: sunminghong, allen.fantasy@gmail.com, http://weibo.com/5d13
#         Team: http://1201.us
#   LastChange: 2015-08-13 17:23:18
#      History:
=============================================================================*/


/*
定义wormhole所有的接口，枚举
*/
package wormhole

//common define
type TID uint32
const TIDSize = 4
const MAX_CONNECTIONS = 0x3fff -1

func (t TID) toI() uint32 {
    return uint32(t)
}


type CommonCallbackFunc func (id int)


//common define end -----------------------------------------------


// data route layer , route packet define
type ERouteType byte
const (
    //0bit =1 表示为需要中转的数据
    EPACKET_TYPE_GENERAL = 0
    EPACKET_TYPE_DELAY = 2 | 1   //00000011
    EPACKET_TYPE_CLOSE = 4 | 1   //00000101   close a player client
    EPACKET_TYPE_BROADCAST = 6 | 1    //00000111

    EPACKET_TYPE_GATE_REGISTER= 8     //00001000
    EPACKET_TYPE_GATE_REMOVE= 10      //00001010 remove a gate client
    EPACKET_TYPE_CLOSED = 12 | 1      //00001100 a player client closed tell to gridserver
    EPACKET_TYPE_DELAY_DATAS = 14 | 1 //00001111 
    EPACKET_TYPE_DELAY_DATAS_COMPRESS = 16 | 1  //00010001
    EPACKET_TYPE_DATAS_COMPRESS = 18            //00010010  to player client connection

    EPACKET_TYPE_FORWARD = 20 | 1               //00010101  forward msg to other grid server 

    EPACKET_TYPE_UDP_SERVER = 22
    EPACKET_TYPE_HELLO = 24  //guin，data里面为发送方类型（如是ageng，client，gameserver）
)


type WriteFunc func (data []byte) (int,error)

//datagram and datapacket define
type IRoutePack interface {
    //Encrypt([]byte)
    //Decrypt([]byte)

    Clone() IRoutePack

    GetEndian() int
    SetEndian(endian int)

    Fetch(conn IConnection) (n int, dps []*RoutePacket)
    Pack(dp *RoutePacket) []byte

    PackWrite(write WriteFunc,dp *RoutePacket)
}


// define a struct or class of rec transport connection
// datapacket = mask1(byte) | mask2(byte) | packetType(byte) | datalength(int32) | data| guin
type RoutePacket struct {
    Type  ERouteType
    Guin TID
    Data  []byte
}

// data route layer end -----------------------------------------------



// connection define
type EConnType byte
const(
    ECONN_TYPE_CTRL = 0
    ECONN_TYPE_DATA = 1
)


type IStream interface {
    GetPos() int
    Len() int
    Read(count int) ([]byte, int)
    SetPos(int)
    Reset()
}


type ReceiveFunc func (conn IConnection)


type IConnection interface {
    Connect(addr string) bool
    Close()
    SetCloseCallback(cf CommonCallbackFunc)

    //connId 比如tcp socketid or udp socketid
    //SetId(TID connectId)
    GetId() int

    GetStream() IStream

    GetType() EConnType
    SetType(t EConnType)

    Send(data []byte)
    //Send(packet *RoutePacket)

    SetReceiveCallback(receive ReceiveFunc)
    //SetReceivePacketCallback(receive ReceivePacketFunc)
}
// connection define end -----------------------------------------------




// udp frame define

// udp frame define end -----------------------------------------------



// wormhole define

type GUIN interface {
    //生成一个guin
    GenerateGuin(agentId int) TID

    Parse(guin TID) (agentId int, id int, check int)
}


type EWormholeState byte
const (
    ECONN_STATE_ACTIVE = 0
    ECONN_STATE_DISCONNTCT = 1
    ECONN_STATE_SUSPEND = 2
)


type EWormholeType byte
const (
    EWORMHOLE_TYPE_CLIENT = 0
    EWORMHOLE_TYPE_GAMESERVER = 1
    EWORMHOLE_TYPE_AGENT = 2
    EWORMHOLE_TYPE_CONSOLE = 3
)


type EServerType byte
const (
    ESERVER_TYPE_GAMESERVER = 1
    ESERVER_TYPE_AGENT = 2
    ESERVER_TYPE_CONSOLE = 3
)

type ReceivePacketFunc func (wh IWormhole, dps []*RoutePacket)

type NewWormholeFunc func(guin TID, wormholeManager IWormholeManager, routepack IRoutePack) IWormhole

type IWormhole interface {
    GetFromType() EWormholeType
    SetFromType(t EWormholeType)

    GetFromId() int
    SetFromId(id int)

    GetGuin() TID

    AddConnection(conn IConnection, t EConnType)

    GetState() EWormholeState
    SetState(state EWormholeState)

    SendRaw(data []byte)
    SendPacket(packet *RoutePacket)
    Send(guin TID, data []byte)
    Broadcast(guin TID, data []byte)

    //SetReceivePacketCallback(receive ReceivePacketFunc)

    ProcessPackets(packets []*RoutePacket)

    Close()
    SetCloseCallback(cf CommonCallbackFunc)

    GetManager() IWormholeManager
}

type IWormholeManager interface {
    Add(wh IWormhole)
    Get(guin TID) (IWormhole,bool)

    Send(guin TID, data []byte)
    Broadcast(guin TID, data []byte)

    Close(guin TID)
    CloseAll()

    Length() int
}

// wormhole define end -----------------------------------------------



// wormhole tcp server define

type ITcpServer interface {

}

// wormhole tcp server define ----------------------------------------
