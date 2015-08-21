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
    //"time"

    iniconfig "github.com/sunminghong/iniconfig"
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
    name string, serverId int,
    clientTcpAddr string, clientUdpAddr string,
    logicTcpAddr string, logicUdpAddr string,
    maxConnections int, routepack IRoutePack,
    clientWormholes IWormholeManager,
    logicWormholes IWormholeManager,
    makeClientWormhole NewWormholeFunc,
    makeLogicWormhole NewWormholeFunc) *Agent {

    s := &Agent{
        Name:name,
        ServerId: serverId,
        ServerType: EWORMHOLE_TYPE_AGENT,
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
        name, serverId, EWORMHOLE_TYPE_AGENT,
        clientTcpAddr, clientUdpAddr, maxConnections,
        routepack, clientWormholes, makeClientWormhole)

    s.logicServer = NewServer(
        name, serverId, EWORMHOLE_TYPE_AGENT,
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


func NewAgentFromIni(
    configfile string,
    routepack IRoutePack,
    clientWormholes IWormholeManager,
    logicWormholes IWormholeManager,
    makeClientWormhole NewWormholeFunc,
    makeLogicWormhole NewWormholeFunc) *Agent {

    c, err := iniconfig.ReadConfigFile(configfile)
    if err != nil {
        gts.Error(err.Error())
        return nil
    }

    section := "Default"

    logconf, err := c.GetString(section,"logConfigFile")
    if err != nil {
        logconf = ""
    }
    gts.SetLogger(&logconf)

    //start grid service
    name, err := c.GetString(section, "name")
    if err != nil {
        gts.Error(err.Error())
        return nil
    }

    serverId, err := c.GetInt(section, "serverId")
    if err != nil {
        gts.Error(err.Error())
        return nil
    }

    clientTcpAddr, err := c.GetString(section, "clientTcpAddr")
    if err != nil {
        gts.Error(err.Error())
        return nil
    }

    clientUdpAddr, err := c.GetString(section, "clientUdpAddr")
    if err != nil {
        gts.Warn(err.Error())
    }

    logicTcpAddr, err := c.GetString(section, "logicTcpAddr")
    if err != nil {
        gts.Error(err.Error())
        return nil
    }

    logicUdpAddr, err := c.GetString(section, "logicUdpAddr")
    if err != nil {
        gts.Warn(err.Error())
    }

    maxConnections, err := c.GetInt(section, "maxConnections")
    if err != nil {
        maxConnections = 1000
    }

    endian, err := c.GetInt(section, "endian")
    if err == nil {
        routepack.SetEndianer(gts.GetEndianer(endian))
    } else {
        routepack.SetEndianer(gts.GetEndianer(gts.LittleEndian))
    }

    /*
    autoDuration, err := c.GetInt(section, "autoReconnectDuration")
    if err != nil {
        autoDuration = 5
    }
    autoReconnectDuration := time.Duration(autoDuration) * time.Second
    */

    return NewAgent(
        name, serverId,
        clientTcpAddr , clientUdpAddr ,
        logicTcpAddr , logicUdpAddr ,
        maxConnections , routepack,
        clientWormholes,
        logicWormholes,
        makeClientWormhole,
        makeLogicWormhole)
}

