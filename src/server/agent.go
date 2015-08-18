/*=============================================================================
#     FileName: agent.go
#       Author: sunminghong, allen.fantasy@gmail.com, http://weibo.com/5d13
#         Team: http://1201.us
#   LastChange: 2015-08-18 19:19:34
#      History:
=============================================================================*/


/*
agent 代理连接服务器，接受玩家客户端连接
*/

package server

import (
    "net"
    "strconv"
    "time"
    "fmt"
    "strings"

    iniconfig "github.com/sunminghong/iniconfig"
    gts "github.com/sunminghong/gotools"
)

func MakeClientWormhole() {

}


type Agent struct {
    Name string
    AgentId int
    Addr string

    AgentType EWormholeType

    maxConnections int
    //endianer        gts.IEndianer

    tcpAddr string
    udpAddr string

    wormholeManager IWormholeManager

    exitChan chan bool
    stop bool


    tcpAgent *TcpAgent
    udpAgent *UdpAgent
}


func NewAgent(
    name string, serverid int, serverType EWormholeType,
    tcpAddr string, udpAddr string, maxConnections int,
    routepack IRoutePack, wm IWormholeManager,
    makeWormhole NewWormholeFunc) *Agent {

    if MAX_CONNECTIONS < maxConnections {
        maxConnections = MAX_CONNECTIONS
    }

    s := &Agent{
        Name:name,
        AgentId:serverid,
        tcpAddr : tcpAddr,
        udpAddr : udpAddr,
        AgentType: serverType,
        maxConnections : maxConnections,
        wormholeManager: wm,
        stop: false,
        exitChan: make(chan bool),
    }

    s.tcpAgent = NewTcpAgent(name, serverid, serverType, tcpAddr, maxConnections, routepack, wm, makeWormhole, NewTcpConnection)

    s.udpAgent = NewUdpAgent(name, serverid, serverType, udpAddr, maxConnections, routepack, wm, makeWormhole, NewUdpConnection)

    return s
}


func (s *Agent) Start() {
    gts.Info(s.Name +" is starting...")
    s.stop=false

    if len(s.tcpAddr) == 0 {
        s.tcpAgent.Start()
    }
    if len(s.udpAddr) == 0 {
        s.udpAgent.Start()
    }
}


func (s *Agent) Stop() {
    s.tcpAgent.Stop()
    s.udpAgent.Stop()
}


