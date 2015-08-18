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
    //"reflect"
    //"strconv"
    //"encoding/binary"
    //"time"
    //"math/rand"

    "net"

    gts "github.com/sunminghong/gotools"
)


type UdpServer struct {
    Name string
    ServerId int

    ServerType EServerType

    maxConnections int
    newConn NewUdpConnectionFunc
    routepack   IRoutePack
    //endianer        gts.IEndianer

    udpAddr string

    wormholeManager IWormholeManager
    udp_read_buffer_size int

    exitChan chan bool
    stop bool

    idassign *gts.IDAssign

    udpAddrs map[string]*UdpConnection
}


type NewUdpConnectionFunc func (newcid TID, conn net.Conn, endian int) *UdpConnection

func NewUdpServer(
    name string,serverid int, serverType EServerType,
    udpAddr string, maxConnections int, newConn NewUdpConnectionFunc,
    routepack IRoutePack, wm IWormholeManager) *UdpServer {

    if MAX_CONNECTIONS < maxConnections {
        maxConnections = MAX_CONNECTIONS
    }

    s := &UdpServer{
        Name:name,
        ServerId:serverid,
        udpAddr : udpAddr,
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
    s.udp_read_buffer_size = 1024
    s.idassign = gts.NewIDAssign(s.maxConnections)

    return s
}


func (s *UdpServer) StartUdp() {
	udpAddr, err := net.ResolveUDPAddr("udp", s.udpAddr)
	if err != nil {
        gts.Error("udp server addr(%s) is error:%q", s.udpAddr, err)
		return
	}
	sock, _err := net.ListenUDP("udp", udpAddr)
	if _err != nil {
        gts.Error("udp server addr(%s) is error2:%q", s.udpAddr, err)
        return
	}

    defer sock.Close()

	for {
        buffer := make([]byte, s.udp_read_buffer_size)
		n, from, err := sock.ReadFromUDP(buffer)
        fromAddr := from.String()
		if err == nil {
			//log.Println("recv", n, from)
            gts.Trace("udp connect from :%s", fromAddr)
			udpConn, ok:= s.udpAddrs[fromAddr]
			if !ok {
                newcid := s.AllocTransportid()
                udpConn := s.newConn(
                    TID(newcid),
                    sock,
                    s.routepack.GetEndian())

                udpConn.SetReceiveCallback(s.receiveUdpBytes)
                s.udpAddrs[fromAddr] = udpConn
			}

            udpConn.ConnReader(buffer[0:n])

			//log.Println("debug out.........")
		} else {
			e, ok := err.(net.Error)
			if !ok || !e.Timeout() {
				gts.Trace("recv error", err.Error(), fromAddr)
				delete(s.udpAddrs, fromAddr)
			}
		}
	}
}

func (s *UdpServer) receiveUdpBytes(conn IConnection) {
    n, dps := s.routepack.Fetch(conn)
    if n > 0 {
        s.receiveUdpPackets(conn, dps)
    }
}


func (s *UdpServer) receiveUdpPackets(conn IConnection, dps []*RoutePacket) {
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
                s.wormholeManager.Add(wh)
            }

            //将该连接绑定到wormhole，
            //并且connection的receivebytes将被wormhole接管
            //该函数将不会被该connection调用
            wh.AddConnection(conn, ECONN_TYPE_DATA)
            gts.Debug("has clients:",s.wormholeManager.Length())


            fromType := EWormholeType(dp.Data[0])
            wh.SetFromType(fromType)
            break
        }
    }
}


/*
func (s *UdpServer) broadcastHandler(broadcastChan <-chan *RoutePacket) {
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
func (s *UdpServer) SendBroadcast(dp *RoutePacket) {
    s.broadcastChan <- dp
}
*/


func (s *UdpServer) Stop() {
    s.stop = true
}


func (s *UdpServer) AllocTransportid() int {
    if (s.wormholeManager.Length() >= s.maxConnections) {
        return 0
    }

    return s.idassign.GetFree()
}


func (s *UdpServer) SetMaxConnections(max int) {
    if MAX_CONNECTIONS < max {
        s.maxConnections = MAX_CONNECTIONS
    } else {
        s.maxConnections = max
    }
}


