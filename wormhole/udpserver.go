/*=============================================================================
#     FileName: udpserver.go
#       Author: sunminghong, allen.fantasy@gmail.com, http://weibo.com/5d13
#         Team: http://1201.us
#   LastChange: 2015-08-17 19:05:50
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


type UdpServer struct {
    *BaseServer

    udp_read_buffer_size int

    makeConn NewUdpConnectionFunc
    makeWormhole NewWormholeFunc

    udpAddrs map[string]*UdpConnection
}


type NewUdpConnectionFunc func (newcid int, conn *net.UDPConn, endianer gts.IEndianer, userAddr *net.UDPAddr) *UdpConnection


func NewUdpServer(
    name string,serverid int, serverType EWormholeType,
    addr string, maxConnections int,
    routePack IRoutePack, wm IWormholeManager,
    makeWormhole NewWormholeFunc,
    makeConn NewUdpConnectionFunc) *UdpServer {

    s := &UdpServer{
        BaseServer: NewBaseServer(
            name,serverid, serverType, addr, maxConnections, routePack, wm),
        makeConn: makeConn,
        makeWormhole: makeWormhole,
        udpAddrs: make(map[string]*UdpConnection),
    }

    s.udp_read_buffer_size = 1024

    return s
}


func (s *UdpServer) Start() {
	udpAddr, err := net.ResolveUDPAddr("udp", s.Addr)
	if err != nil {
        gts.Error("udp server addr(%s) is error:%q", s.Addr, err)
		return
	}

    gts.Info(s.Name +" udpserver is starting...")
	sock, _err := net.ListenUDP("udp", udpAddr)
	if _err != nil {
        gts.Error("udp server addr(%s) is error2:%q", s.Addr, err)
        return
	}

    go func() {
        buffer := make([]byte, s.udp_read_buffer_size)
        for {
            defer sock.Close()
            if s.Stop_ {
                return
            }

            n, fromAddr, err := sock.ReadFromUDP(buffer)
            key := fromAddr.String()
            if err == nil {
                //log.Println("recv", n, from)
                gts.Trace("udp connect from: %s", fromAddr)
                udpConn, ok:= s.udpAddrs[key]
                if !ok {
                    newcid := s.AllocId()
                    udpConn = s.makeConn(
                        newcid,
                        sock,
                        s.RoutePackHandle.GetEndianer(),
                        fromAddr,
                    )

                    gts.Trace("new udp connection")
                    udpConn.SetReceiveCallback(s.receiveUdpBytes)
                    s.udpAddrs[key] = udpConn
                }

                udpConn.ConnReader(buffer[0:n])
            } else {
                e, ok := err.(net.Error)
                if !ok || !e.Timeout() {
                    gts.Trace("recv errorserver:%s,%q.", key, err.Error())
                    delete(s.udpAddrs, key)
                }
            }
        }
    }()
}


func (s *UdpServer) receiveUdpBytes(conn IConnection) {
    gts.Trace("udp server receiveBytes:% X", conn.GetBuffer().Stream.Bytes())
    n, dps := s.RoutePackHandle.Fetch(conn.GetBuffer())
    if n > 0 {
        s.receiveUdpPackets(conn, dps)
    }
}


func (s *UdpServer) receiveUdpPackets(conn IConnection, dps []*RoutePacket) {
    for _, dp := range dps {
        if dp.Type == EPACKET_TYPE_HELLO {
            gts.Trace("server receive udp hello:%q", dp)
             //接到连接方hello包
            var guin int
            var wh IWormhole
            var ok bool

            if dp.Guin > 0 {
                //TODO:重连处理
                //如果客户端hello是发送的guin有，则表示是重连，需要处理重连逻辑
                //比如在一定时间内可以重新连接会原有wormhole

                guin = dp.Guin
                if wh, ok = s.Wormholes.Get(guin);ok {
                    if wh.GetState() == ECONN_STATE_SUSPEND {
                        wh.SetState(ECONN_STATE_ACTIVE)
                    } else if wh.GetState() == ECONN_STATE_DISCONNTCT {
                        wh = nil
                    }
                }
            } else {
                guin = GenerateGuin(s.ServerId, int(conn.GetId()))
            }
            if wh == nil {
                wh = s.makeWormhole(guin, s.Wormholes, s.RoutePackHandle)
                s.Wormholes.Add(wh)
            }

            //将该连接绑定到wormhole，
            //并且connection的receivebytes将被wormhole接管
            //该函数将不会被该connection调用
            wh.AddConnection(conn, ECONN_TYPE_DATA)
            gts.Trace("has wormholes:%d\n-----------------------------------", s.Wormholes.Length())

            fromType := EWormholeType(dp.Data[0])
            wh.SetFromType(fromType)
            break
        }
    }
}

