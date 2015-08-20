/*=============================================================================
#     FileName: logic.go
#       Author: sunminghong, allen.fantasy@gmail.com, http://weibo.com/5d13
#         Team: http://1201.us
#   LastChange: 2015-08-18 19:19:34
#      History:
=============================================================================*/


/*
logic 代理连接服务器，接受玩家客户端连接
*/

package server

import (
    gts "github.com/sunminghong/gotools"

    . "wormhole/wormhole"
)


type Agent struct {
    Name string
    ServerId int

    ServerType EWormholeType

    maxConnections int

    /*
    clientTcpAddr string
    clientUdpAddr string

    logicTcpAddr string
    clientUdpAddr string

    makeClientWormhole NewWormholeFunc
    makeLogicWormhole NewWormholeFunc
    */

    ClientWormholes IWormholeManager
    LogicWormholes IWormholeManager

    clientServer *Server
    logicServer *Server
}


func NewAgent(
    name string, serverid int, serverType EWormholeType,
    clientTcpAddr string, clientUdpAddr string,
    logicTcpAddr string, logicUdpAddr string,
    maxConnections int, routepack IRoutePack,
    clientWormholes IWormholeManager,
    logicWormholes IWormholeManager,
    makeClientWormhole NewWormholeFunc,
    makeLogicWormhole NewWormholeFunc) *Agent {

    s := &Agent{
        Name:name,
        ServerId: serverid,
        ServerType: serverType,
        //clientTcpAddr : clientTcpAddr,
        //clientUdpAddr : clientUdpAddr,
        //logicTcpAddr : logicTcpAddr,
        //logicUdpAddr : logicUdpAddr,
        //makeClientWormhole: makeClientWormhole,
        //makeLogicWormhole: makeLogicWormhole,

        ClientWormholes : clientWormholes,
        LogicWormholes : logicWormholes,

        maxConnections : maxConnections,
    }

    clientWormholes.SetServer(s)
    logicWormholes.SetServer(s)

    s.clientServer = NewServer(
        name, serverid, EWORMHOLE_TYPE_AGENT,
        clientTcpAddr, clientUdpAddr, maxConnections,
        routepack, clientWormholes, makeClientWormhole)

    s.logicServer = NewServer(
        name, serverid, EWORMHOLE_TYPE_AGENT,
        logicTcpAddr, logicUdpAddr, 100,
        routepack, logicWormholes, makeLogicWormhole)

    return s
}


func (s *Agent) GetServerId() int {
    return s.ServerId
}


func (s *Agent) Start() {
    gts.Info(s.Name +" is starting...")

    s.clientServer.Start()
    s.logicServer.Start()
}


func (s *Agent) Stop() {
    s.clientServer.Stop()
    s.logicServer.Stop()
}


