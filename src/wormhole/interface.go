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
func (t TID) toI() uint32 {
    return uint32(t)
}


type CommonCallbackFunc func (id TID)

type ReceiveFunc func (dataFromId TID, data []byte)
type ReceivePacketFunc func (dataFromId TID, dps []*RoutePacket)


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

/*
type IStream interface {
    GetPos() int
    Len() int
    Read(count int)
    SetPos(int)
    Reset()
}
*/

type IConnection interface {
    Connect(addr string) bool
    Close()
    SetCloseCallback(cf CommonCallbackFunc)

    //connId 比如tcp socketid or udp socketid
    //SetId(int connectId)
    GetId() TID

    GetType() EConnType
    SetType(t EConnType)

    Send(packet *RoutePacket)

    //SetReceiveCallback(receive ReceiveFunc)
    SetReceivePacketCallback(receive ReceivePacketFunc)
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

type IWormhole interface {
    GetFromAgentId() int

    GetGuin() TID

    AddConnection(conn IConnection, t EConnType)

    GetState() EWormholeState
    SetState(state EWormholeState)

    Send(packet *RoutePacket)
    Broadcast(packet *RoutePacket)

    SetReceivePacketCallback(receive ReceivePacketFunc)

    Close()
    SetCloseCallback(cf CommonCallbackFunc)

    GetManager() IWormholeManager
}

type IWormholeManager interface {
    AddWormhole(wh IWormhole)

    Send(guin TID, packet *RoutePacket)
    Broadcast(packet *RoutePacket)

    Close(guin TID)
    CloseAll()
}

// wormhole define end -----------------------------------------------


