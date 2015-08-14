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


type CommonCallbackFunc func (id TID)

type ReceiveFunc func (dataFromId TID, data []byte)
type ReceivePacketFunc func (dataFromId TID, dps []*RoutePacket)


//common define end -----------------------------------------------


// data route layer , route packet define
type TRouteType byte
const (
    //0bit =1 表示为需要中转的数据
    TPACKET_TYPE_GENERAL = 0
    TPACKET_TYPE_DELAY = 2 | 1   //00000011
    TPACKET_TYPE_CLOSE = 4 | 1   //00000101   close a player client
    TPACKET_TYPE_BROADCAST = 6 | 1    //00000111
    TPACKET_TYPE_GATE_REGISTER= 8     //00001000
    TPACKET_TYPE_GATE_REMOVE= 10      //00001010 remove a gate client
    TPACKET_TYPE_CLOSED = 12 | 1      //00001100 a player client closed tell to gridserver
    TPACKET_TYPE_DELAY_DATAS = 14 | 1 //00001111 
    TPACKET_TYPE_DELAY_DATAS_COMPRESS = 16 | 1  //00010001
    TPACKET_TYPE_DATAS_COMPRESS = 18            //00010010  to player client connection

    TPACKET_TYPE_FORWARD = 20 | 1               //00010101  forward msg to other grid server 
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
    Type  TRouteType
    Guin TID
    Data  []byte
}

// data route layer end -----------------------------------------------





// connection define
type TConnType byte
const(
    CONN_TYPE_CTRL = 0
    CONN_TYPE_DATA = 1
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
    /*
    SetRemoteAddress(host string, port int16)
    GetRemoteAddress() (host string, port int16)
    */

    Connect(addr string) bool
    Close()
    SetCloseCallback(cf CommonCallbackFunc)

    //connId 比如tcp socketid or udp socketid
    //SetId(int connectId)
    GetId() TID

    GetType() TConnType
    SetType(t TConnType)

    //Send(buf []byte)
    Send(dp *RoutePacket)

    //SetReceiveCallback(receive ReceiveFunc)
    SetReceivePacketCallback(receive ReceivePacketFunc)
}
// connection define end -----------------------------------------------




// udp frame define

// udp frame define end -----------------------------------------------



// wormhole define

type GUIN interface {
    //配置进行guin计算的key，用于验算和随机
    SetKey(hashkey int)

    //生成一个guin
    GenerateGuin(agentId int, wormholdId int) TID

    Check(guin TID) bool
}


type TWormholeState byte
const (
    CONN_STATE_ACTIVE = 0
    CONN_STATE_DISCONNTCT = 1
    CONN_STATE_SUSPEND = 2
)

type wormhole interface {
    GetGuin() TID
    AddConnection(conn IConnection, t TConnType)

    Send(paket *RoutePacket)
    Broadcast(packet *RoutePacket)
    SetReceiveCallback(receive ReceivePacketFunc)

    Close()
    CloseCallback(cf CommonCallbackFunc)
}

type IWormholeManager interface {

    Send(guin TID, buf []byte)
    Broadcast(buf []byte)

    Close(guin TID)
    CloseAll()
}

// wormhole define end -----------------------------------------------



