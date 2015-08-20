/*=============================================================================
#     FileName: server.go
#       Author: sunminghong, allen.fantasy@gmail.com, http://weibo.com/5d13
#         Team: http://1201.us
#   LastChange: 2015-08-17 18:14:54
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
    //"net"
    //"time"
    //"math/rand"

    gts "github.com/sunminghong/gotools"
)


type Server struct {
    Name string
    ServerId int

    ServerType EWormholeType

    maxConnections int
    //endianer        gts.IEndianer

    tcpAddr string
    udpAddr string

    wormholeManager IWormholeManager

    exitChan chan bool
    stop bool


    tcpServer *TcpServer
    udpServer *UdpServer
}


func NewServer(
    name string, serverid int, serverType EWormholeType,
    tcpAddr string, udpAddr string, maxConnections int,
    routepack IRoutePack, wm IWormholeManager,
    makeWormhole NewWormholeFunc) *Server {

    if MAX_CONNECTIONS < maxConnections {
        maxConnections = MAX_CONNECTIONS
    }

    s := &Server{
        Name:name,
        ServerId:serverid,
        tcpAddr : tcpAddr,
        udpAddr : udpAddr,
        ServerType: serverType,
        maxConnections : maxConnections,
        wormholeManager: wm,
        stop: false,
        exitChan: make(chan bool),
    }

    s.tcpServer = NewTcpServer(name, serverid, serverType, tcpAddr, maxConnections, routepack, wm, makeWormhole, NewTcpConnection)

    s.udpServer = NewUdpServer(name, serverid, serverType, udpAddr, maxConnections, routepack, wm, makeWormhole, NewUdpConnection)

    return s
}


func (s *Server) GetServerId() int {
    return s.ServerId
}


func (s *Server) Start() {
    gts.Info(s.Name +" is starting...")
    s.stop=false

    if len(s.tcpAddr) == 0 {
        s.tcpServer.Start()
    }
    if len(s.udpAddr) == 0 {
        s.udpServer.Start()
    }
}


func (s *Server) Stop() {
    s.tcpServer.Stop()
    s.udpServer.Stop()
}


