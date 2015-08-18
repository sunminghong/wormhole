/*=============================================================================
#     FileName: tcpclient.go
#       Author: sunminghong, allen.fantasy@gmail.com, http://weibo.com/5d13
#         Team: http://1201.us
#   LastChange: 2015-08-17 17:59:44
#      History:
=============================================================================*/

/*
定义一个基本的Tcp wormhole 服务器

1、连接上tcpserver后；
2、立即向服务器发送hello 包（如果本地缓存了 GUIN 就将 GUIN 发送过去）
3、收到服务器端发回的HELLO 后，保存guin，并生成wormhole
4、收到了UDP_SERVER 后就发起udp连接
*/

package wormhole


//import (
//    //"reflect"
//    //"strconv"
//    "encoding/binary"
//    "net"
//    "time"
//    "math/rand"

//    gts "github.com/sunminghong/gotools"
//)


//type TcpClient struct {
//    Name string
//    Addr string
//    udpAddr string

//    FromType EWormholeType

//    newConn NewTcpConnectionFunc
//    routepack   IRoutePack
//    //endianer        gts.IEndianer

//    host string
//    port int

//    nextId int

//    exitChan chan bool
//    stop bool
//}


//func NewTcpClient(
//    fromType EWormholeType, addr string, newConn NewTcpConnectionFunc, routepack IRoutePack, wm IWormholeManager) *TcpClient {

//    s := &TcpClient{
//        Name:name,
//        Addr : addr,
//        FromType: fromType,
//        wormholeManager: wm,
//        stop: false,
//        exitChan: make(chan bool),
//        newConn: newConn,
//        routepack: routepack,
//    }

//    [>
//    if s.routepack.GetEndian() == gts.BigEndian {
//        s.endianer = binary.BigEndian
//    } else {
//        s.endianer = binary.LittleEndian
//    }
//    */

//    return s
//}


//func (s *TcpClient) SetUdpAdd(udpAddr string) {
//    s.udpAddr = udpAddr
//}


//func (s *TcpClient) Start() {
//    gts.Info(s.Name +" is starting...")

//    s.stop=false
//    //todo: maxConnections don't proccess
//    //addr := host + ":" + strconv.Itoa(port)

//    netListen, error := net.Listen("tcp", s.Addr)
//    if error != nil {
//        gts.Error(error)
//    } else {
//        gts.Info("listen with :", s.Addr)
//        gts.Info(s.Name +" is started !!!")

//        //defer函数退出时执行
//        defer netListen.Close()
//        for {
//            gts.Trace("Waiting for connection")
//            connection, err := netListen.Accept()
//            if s.stop {
//                break
//            }

//            if err != nil {
//                gts.Error("Transport error: ", err)
//            } else {
//                gts.Debug("%v is connection!",connection.RemoteAddr())

//                newcid := s.AllocTransportid()
//                if newcid == 0 {
//                    gts.Warn("connection num is more than ",s.maxConnections)
//                } else {
//                    gts.Trace("//////////////////////newcid:",newcid)
//                    s.transportHandler(newcid, connection)
//                }
//            }
//        }
//    }
//}


////该函数主要是接受新的连接和注册用户在transport list
//func (s *TcpClient) transportHandler(newcid int, connection net.Conn) {
//    tcpConn := s.newConn(TID(newcid), connection, s.routepack.GetEndian())
//    tcpConn.SetReceiveCallback(s.receiveBytes)
//}


//func (s *TcpClient) receiveBytes(conn IConnection) {
//    n, dps := s.routepack.Fetch(conn)
//    if n > 0 {
//        s.receivePackets(conn, dps)
//    }
//}


//func (s *TcpClient) receivePackets(conn IConnection, dps []*RoutePacket) {
//    for _, dp := range dps {
//        if dp.Type == EPACKET_TYPE_HELLO {
//             //接到连接方hello包
//            var guin TID
//            var wh IWormhole

//            if dp.Guin > 0 {
//                //TODO:重连处理
//                //如果客户端hello是发送的guin有，则表示是重连，需要处理重连逻辑
//                //比如在一定时间内可以重新连接会原有wormhole

//                guin = dp.Guin
//                if wh, ok := s.wormholeManager.Get(guin);ok {
//                    if wh.GetState() != ECONN_STATE_DISCONNTCT {
//                        wh = NewWormhole(guin, s.wormholeManager, s.routepack)
//                        s.idassign.Free(conn.GetId())
//                    }
//                }
//            }
//            if wh == nil {
//                guin := GetGuin(s.ServerId, int(conn.GetId()))
//                wh = NewWormhole(guin, s.wormholeManager, s.routepack)
//            }

//            //将该连接绑定到wormhole，
//            //并且connection的receivebytes将被wormhole接管
//            //该函数将不会被该connection调用
//            wh.AddConnection(conn, ECONN_TYPE_CTRL)
//            s.wormholeManager.Add(wh)
//            gts.Debug("has clients:",s.wormholeManager.Length())


//            fromType := EWormholeType(dp.Data[0])
//            wh.SetFromType(fromType)
//            //if fromType == EWORMHOLE_TYPE_CLIENT {
//                //wh.SetFromId(0)
//            //} else if fromType == EWORMHOLE_TYPE_SERVER {
//                //wh.SetFromId(0)
//            //}

//            //hello to client 
//            packet := &RoutePacket {
//                Type:   EPACKET_TYPE_HELLO,
//                Guin:   guin,
//                Data:   []byte{byte(s.ServerType)},
//            }
//            wh.SendPacket(packet)

//            //hello udp addr to client
//            if len(s.udpAddr) == 0 {
//                packet.Type = EPACKET_TYPE_UDP_SERVER
//                packet.Data = []byte(s.udpAddr)
//                wh.SendPacket(packet)
//            }

//            break
//        }

//        //s.receivePacketCallback(wh, []*RoutePacket{dp})
//    }
//}


////send broadcast message data for other object
//func (s *TcpClient) SendBroadcast(dp *RoutePacket) {
//    s.broadcastChan <- dp
//}
//*/


//func (s *TcpClient) Stop() {
//    s.stop = true
//}


//func (s *TcpClient) AllocTransportid() int {
//    if (s.wormholeManager.Length() >= s.maxConnections) {
//        return 0
//    }

//    return s.idassign.GetFree()
//}


//func (s *TcpClient) SetMaxConnections(max int) {
//    if MAX_CONNECTIONS < max {
//        s.maxConnections = MAX_CONNECTIONS
//    } else {
//        s.maxConnections = max
//    }
//}


