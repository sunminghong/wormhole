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
    "net"
    gts "github.com/sunminghong/gotools"
)


type NewTcpConnectionFunc func(newcid int, conn net.Conn, endianer gts.IEndianer) *TcpConnection


type TcpServer struct {
    *BaseServer

    makeConn NewTcpConnectionFunc
    makeWormhole NewWormholeFunc

    udpAddr string
}


func NewTcpServer(
    name string,serverid int, serverType EWormholeType,
    addr string, MaxConns int,
    RoutePackHandle IRoutePack, wm IWormholeManager,
    makeWormhole NewWormholeFunc,
    makeConn NewTcpConnectionFunc) *TcpServer {

    s := &TcpServer{
        BaseServer: NewBaseServer(
            name, serverid, serverType, addr, MaxConns, RoutePackHandle, wm),
        makeConn: makeConn,
        makeWormhole: makeWormhole,
    }

    return s
}


func (s *TcpServer) SetUdpAddr(udpAddr string) {
    s.udpAddr = udpAddr
}


func (s *TcpServer) Start() {
    gts.Info(s.Name +" is starting...")

    s.Stop_=false
    //todo: MaxConns don't proccess
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
        if s.Stop_ {
            break
        }

        connection, err := netListen.Accept()

        if err != nil {
            gts.Error("Transport error: ", err)
        } else {
            gts.Debug("%v is connection!",connection.RemoteAddr())

            newcid := s.AllocId()
            if newcid == 0 {
                gts.Warn("connection num is more than ",s.MaxConns)
            } else {
                gts.Trace("//////////////////////newcid:",newcid)
                tcpConn := s.makeConn(newcid, connection,
                    s.RoutePackHandle.GetEndianer())
                tcpConn.SetReceiveCallback(s.receiveBytes)
            }
        }
    }
}


func (s *TcpServer) receiveBytes(conn IConnection) {
    n, dps := s.RoutePackHandle.Fetch(conn)
    if n > 0 {
        s.receivePackets(conn, dps)
    }
}


func (s *TcpServer) receivePackets(conn IConnection, dps []*RoutePacket) {
    for _, dp := range dps {
        if dp.Type == EPACKET_TYPE_HELLO {
            gts.Trace("server receive tcp hello:%q", dp)

            var guin int
            var wh IWormhole

            if dp.Guin > 0 {
                //TODO:重连处理
                //如果客户端hello是发送的guin有，则表示是重连，需要处理重连逻辑
                //比如在一定时间内可以重新连接会原有wormhole

                guin = dp.Guin
                if wh, ok := s.Wormholes.Get(guin);ok {
                    if wh.GetState() == ECONN_STATE_SUSPEND {
                        wh.SetState(ECONN_STATE_ACTIVE)
                    } else if wh.GetState() == ECONN_STATE_DISCONNTCT {
                        wh = nil
                    }
                }
            }
            if wh == nil {
                guin := GenerateGuin(s.ServerId, int(conn.GetId()))
                wh = s.makeWormhole(guin, s.Wormholes, s.RoutePackHandle)
            }

            //将该连接绑定到wormhole，
            //并且connection的receivebytes将被wormhole接管
            //该函数将不会被该connection调用
            wh.AddConnection(conn, ECONN_TYPE_CTRL)
            s.Wormholes.Add(wh)
            gts.Debug("has clients:",s.Wormholes.Length())


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

            wh.Init()

            //hello udp addr to client
            if len(s.udpAddr) == 0 {
                packet.Type = EPACKET_TYPE_UDP_SERVER
                packet.Data = []byte(s.udpAddr)
                wh.SendPacket(packet)
            }
            wh.SendPacket(packet)

            gts.Trace("server send back hello:%q", packet)
            break
        }
    }
}

