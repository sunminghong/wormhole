/*=============================================================================
#     FileName: client.go
#       Author: sunminghong, allen.fantasy@gmail.com, http://weibo.com/5d13
#         Team: http://1201.us
#   LastChange: 2015-08-18 15:52:20
#      History:
=============================================================================*/


/*
定义一个基本的wormhole 客户端

1、连接上tcpserver后；
2、立即向服务器发送hello 包（如果本地缓存了 GUIN 就将 GUIN 发送过去）
3、收到服务器端发回的HELLO 后，保存guin，并生成wormhole
4、收到了UDP_SERVER 后就发起udp连接
*/

package wormhole

import (
    //"reflect"
    //"strconv"
    //"net"
    //"time"
    //"math/rand"

    gts "github.com/sunminghong/gotools"
)


type Client struct {
    //endianer        gtc.IEndianer
    guin TID

    tcpAddr string
    udpAddr string

    wormholes IWormholeManager
    makeWormhole NewWormholeFunc

    exitChan chan bool
    stop bool


    tcpConn *TcpConnection
    udpConn *UdpConnection

    routepack IRoutePack
    wormhole IWormhole
    wormType EWormholeType
}


func NewClient(
    tcpAddr string, udpAddr string,
    routepack IRoutePack, wm IWormholeManager,
    makeWormhole NewWormholeFunc,
    wormType EWormholeType) *Client {

    c := &Client{
        tcpAddr : tcpAddr,
        udpAddr : udpAddr,
        routepack: routepack,
        wormholes: wm,
        stop: false,
        exitChan: make(chan bool),
    }

    c.wormType = wormType

    return c
}


func (c *Client) Connect() {
    c.tcpConn = NewTcpConnection(1, nil, c.routepack.GetEndian())
    c.tcpConn.SetReceiveCallback(c.receiveTcpBytes)

    if c.tcpConn.Connect(c.tcpAddr) {
        //连接上服务器
        gts.Info("连接上tcpserver：%s", c.tcpAddr)

        //hello to tcp server
        packet := &RoutePacket {
            Type:   EPACKET_TYPE_HELLO,
            Guin:   c.guin,
            Data:   []byte{byte(c.wormType)},
        }
        c.tcpConn.Send(c.routepack.Pack(packet))
    }
}


func (c *Client) Close() {
    c.udpConn.Close()
    c.tcpConn.Close()
}


func (c *Client) receiveTcpBytes(conn IConnection) {
    n, dps := c.routepack.Fetch(conn)
    if n > 0 {
        c.receiveTcpPackets(conn, dps)
    }
}


func (c *Client) receiveTcpPackets(conn IConnection, dps []*RoutePacket) {
    for _, dp := range dps {
        if dp.Type == EPACKET_TYPE_HELLO {
            c.guin = dp.Guin
            c.initWormhole(dp, conn)

        } else if dp.Type == EPACKET_TYPE_UDP_SERVER {
            c.guin = dp.Guin

            //连接udp server
            c.udpConn = NewUdpConnection(1, nil, c.routepack.GetEndian(), nil)

            if c.udpConn.Connect(c.udpAddr) {
                gts.Info("dial to udp server success:%s", c.udpAddr)

                //hello to tcp server
                packet := &RoutePacket {
                    Type:   EPACKET_TYPE_HELLO,
                    Guin:   c.guin,
                    Data:   []byte{byte(c.wormType)},
                }
                c.udpConn.Send(c.routepack.Pack(packet))

                c.initWormhole(dp, c.udpConn)

            } else {
                gts.Warn("dial to udp server lost:%s", c.udpAddr)
            }
        }
    }
}

func (c *Client) initWormhole(dp *RoutePacket, conn IConnection) {
     //接到连接方hello包
    if dp.Guin == 0 {
        //这hello不正常，关掉连接
        c.Close()
        return
    }

    //TODO:重连处理
    //如果客户端hello是发送的guin有，则表示是重连，需要处理重连逻辑
    //比如在一定时间内可以重新连接会原有wormhole

    if c.wormhole == nil {
        c.wormhole = c.makeWormhole(dp.Guin, c.wormholes, c.routepack)
    }

    //将该连接绑定到wormhole，
    //并且connection的receivebytes将被wormhole接管
    //该函数将不会被该connection调用
    c.wormhole.AddConnection(conn, ECONN_TYPE_CTRL)
    if c.wormholes != nil {
        c.wormholes.Add(c.wormhole)
        gts.Debug("has clients:",c.wormholes.Length())
    }

    fromType := EWormholeType(dp.Data[0])
    c.wormhole.SetFromType(fromType)
}


