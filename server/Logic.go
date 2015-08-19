/*=============================================================================
#     FileName: client.go
#       Author: sunminghong, allen.fantasy@gmail.com, http://weibo.com/5d13
#         Team: http://1201.us
#   LastChange: 2015-08-18 15:52:20
#      History:
=============================================================================*/


/*
定义一个基本的逻辑服务器端

*/
package server


import (
    //"reflect"
    //"strconv"
    //"net"
    //"time"
    //"math/rand"

    //gts "github.com/sunminghong/gotools"

    . "wormhole/wormhole"
)


type Logic struct {

    makeWormhole NewWormholeFunc
    wormholes IWormholeManager

    routepack IRoutePack

    serverType EWormholeType
    serverId int

    clients map[string]*Client
}


func NewLogic(
    name string, serverId int,serverType EWormholeType,
    routepack IRoutePack, wm IWormholeManager,
    makeWormhole NewWormholeFunc) *Logic {

    ls := &Logic{
        routepack: routepack,

        serverId: serverId,
        makeWormhole: makeWormhole,
        wormholes: wm,
        serverType: serverType,
    }

    return ls
}


func (ls *Logic) ConnectAgent(tcpAddr string, udpAddr string) {
    c := NewClient(tcpAddr, udpAddr, ls.routepack, ls.wormholes,
        ls.makeWormhole, ls.serverType)

    ls.clients[tcpAddr] = c
}


func (ls *Logic) Close() {
    //for c := range ls.clients {
        //c.Close()
    //}

    //ls.clients.Clear()

    ls.wormholes.CloseAll()
}

