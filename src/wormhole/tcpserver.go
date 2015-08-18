/*=============================================================================
#     FileName: tcpserver.go
#       Author: sunminghong, allen.fantasy@gmail.com, http://weibo.com/5d13
#         Team: http://1201.us
#   LastChange: 2015-08-17 17:54:16
#      History:
=============================================================================*/

/*
定义一个基本的Tcp wormhole 服务器

1、当收到一个tcp连接请求后，定时等待对方的hello 包
2、收到hello包后，生成一个wormhole 通道对象，并且将GUIN (wormhole id) 通过 REPOERT_GUIN 发会连接方
3、如果该tcp服务器设置了关联udp server address，发送EPACKET_TYPE_UDP_SERVER 包将udpserver address 告诉连接方
*/

package wormhole


import (
    //"reflect"
    //"strconv"
    //"encoding/binary"
    //"time"
    //"math/rand"

    "net"
    gts "github.com/sunminghong/gotools"
)


type NewTcpConnectionFunc func (newcid TID, conn net.Conn, endian int) IConnection


type TcpServer struct {
    Name string
    ServerId int
    Addr string

    ServerType EServerType

    /*
    broadcast_chan_num int
    read_buffer_size   int
    broadcastChan chan *RoutePacket
    */

    maxConnections int
    newConn NewTcpConnectionFunc
    routepack   IRoutePack
    //endianer        gts.IEndianer

    host string
    port int

    wormholeManager IWormholeManager

    exitChan chan bool
    stop bool

    idassign *gts.IDAssign

    udpAddr string
}


func NewTcpServer(
    name string,serverid int, serverType EServerType,addr string, maxConnections int,
    newConn NewTcpConnectionFunc, routepack IRoutePack, wm IWormholeManager) *TcpServer {

    if MAX_CONNECTIONS < maxConnections {
        maxConnections = MAX_CONNECTIONS
    }

    s := &TcpServer{
        Name:name,
        ServerId:serverid,
        Addr : addr,
        ServerType: serverType,
        maxConnections : maxConnections,
        wormholeManager: wm,
        stop: false,
        exitChan: make(chan bool),
        newConn: newConn,
        routepack: routepack,
    }

    /*
    if s.routepack.GetEndian() == gts.BigEndian {
        s.endianer = binary.BigEndian
    } else {
        s.endianer = binary.LittleEndian
    }
    */

    s.idassign = gts.NewIDAssign(s.maxConnections)

    return s
}


func (s *TcpServer) SetUdpAddr(udpAddr string) {
    s.udpAddr = udpAddr
}


func (s *TcpServer) Start() {
    gts.Info(s.Name +" is starting...")

    s.stop=false
    //todo: maxConnections don't proccess
    //addr := host + ":" + strconv.Itoa(port)

    /*
    //创建一个管道 chan map 需要make creates slices, maps, and channels only
    s.broadcastChan = make(chan *RoutePacket, s.broadcast_chan_num)
    go s.broadcastHandler(s.broadcastChan)
    */

    netListen, error := net.Listen("tcp", s.Addr)
    if error != nil {
        gts.Error(error)
        return
    }

    gts.Info("listen with :", s.Addr)
    gts.Info(s.Name +" is started !!!")

    //defer函数退出时执行
    defer netListen.Close()
    for {
        gts.Trace("Waiting for connection")
        connection, err := netListen.Accept()
        if s.stop {
            break
        }

        if err != nil {
            gts.Error("Transport error: ", err)
        } else {
            gts.Debug("%v is connection!",connection.RemoteAddr())

            newcid := s.AllocTransportid()
            if newcid == 0 {
                gts.Warn("connection num is more than ",s.maxConnections)
            } else {
                gts.Trace("//////////////////////newcid:",newcid)
                s.transportHandler(newcid, connection)
            }
        }
    }
}


//该函数主要是接受新的连接和注册用户在transport list
func (s *TcpServer) transportHandler(newcid int, connection net.Conn) {
    tcpConn := s.newConn(TID(newcid), connection, s.routepack.GetEndian())
    tcpConn.SetReceiveCallback(s.receiveBytes)
}


func (s *TcpServer) receiveBytes(conn IConnection) {
    n, dps := s.routepack.Fetch(conn)
    if n > 0 {
        s.receivePackets(conn, dps)
    }
}


func (s *TcpServer) receivePackets(conn IConnection, dps []*RoutePacket) {
    for _, dp := range dps {
        if dp.Type == EPACKET_TYPE_HELLO {
             //接到连接方hello包
            var guin TID
            var wh IWormhole

            if dp.Guin > 0 {
                //TODO:重连处理
                //如果客户端hello是发送的guin有，则表示是重连，需要处理重连逻辑
                //比如在一定时间内可以重新连接会原有wormhole

                guin = dp.Guin
                if wh, ok := s.wormholeManager.Get(guin);ok {
                    if wh.GetState() == ECONN_STATE_SUSPEND {
                        wh.SetState(ECONN_STATE_ACTIVE)
                    } else if wh.GetState() == ECONN_STATE_DISCONNTCT {
                        wh = nil
                    }
                }
            }
            if wh == nil {
                guin := GetGuin(s.ServerId, int(conn.GetId()))
                wh = NewWormhole(guin, s.wormholeManager, s.routepack)
            }

            //将该连接绑定到wormhole，
            //并且connection的receivebytes将被wormhole接管
            //该函数将不会被该connection调用
            wh.AddConnection(conn, ECONN_TYPE_CTRL)
            s.wormholeManager.Add(wh)
            gts.Debug("has clients:",s.wormholeManager.Length())


            fromType := EWormholeType(dp.Data[0])
            wh.SetFromType(fromType)
            //if fromType == EWORMHOLE_TYPE_CLIENT {
                //wh.SetFromId(0)
            //} else if fromType == EWORMHOLE_TYPE_SERVER {
                //wh.SetFromId(0)
            //}

            //hello to client 
            packet := &RoutePacket {
                Type:   EPACKET_TYPE_HELLO,
                Guin:   guin,
                Data:   []byte{byte(s.ServerType)},
            }
            wh.SendPacket(packet)

            //hello udp addr to client
            if len(s.udpAddr) == 0 {
                packet.Type = EPACKET_TYPE_UDP_SERVER
                packet.Data = []byte(s.udpAddr)
                wh.SendPacket(packet)
            }

            break
        }

        //s.receivePacketCallback(wh, []*RoutePacket{dp})
    }
}


/*
func (s *TcpServer) broadcastHandler(broadcastChan <-chan *RoutePacket) {
    for {
        dp := <-broadcastChan

        dp0 := &RoutePacket{
            Type: DATAPACKET_TYPE_GENERAL,
            FromCid: 0,
            Data: dp.Data,
        }

        for _, c := range s.Connections.All() {
            gts.Trace("broadcastHandler: client.type",c.GetType())
            //if fromCid == c.GetTransport().Cid {
                //continue
            //}
            if c.GetType() == CLIENT_TYPE_GATE {
                c.GetTransport().Outgoing <- dp
            } else {
                c.GetTransport().Outgoing <- dp0
            }
        }
        //gts.Trace("broadcastHandler: Handle end!")
    }
}


//send broadcast message data for other object
func (s *TcpServer) SendBroadcast(dp *RoutePacket) {
    s.broadcastChan <- dp
}
*/


func (s *TcpServer) Stop() {
    s.stop = true
}


func (s *TcpServer) AllocTransportid() int {
    if (s.wormholeManager.Length() >= s.maxConnections) {
        return 0
    }

    return s.idassign.GetFree()
}


func (s *TcpServer) SetMaxConnections(max int) {
    if MAX_CONNECTIONS < max {
        s.maxConnections = MAX_CONNECTIONS
    } else {
        s.maxConnections = max
    }
}


