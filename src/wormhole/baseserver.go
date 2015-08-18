/*=============================================================================
#     FileName: serverbase.go
#       Author: sunminghong, allen.fantasy@gmail.com, http://weibo.com/5d13
#         Team: http://1201.us
#   LastChange: 2015-08-18 11:28:18
#      History:
=============================================================================*/

package wormhole


import (
    gts "github.com/sunminghong/gotools"
)


type BaseServer struct {
    Name string
    ServerId int
    Addr string

    ServerType EWormholeType

    MaxConns int
    RoutePackHandle   IRoutePack
    //endianer        gts.IEndianer

    Wormholes IWormholeManager

    exitChan chan bool
    Stop_ bool

    idassign *gts.IDAssign

    udpAddr string
}


func NewBaseServer(
    name string,serverid int, serverType EWormholeType,
    addr string, maxConnections int,
    routePack IRoutePack, wm IWormholeManager) *BaseServer {

    if MAX_CONNECTIONS < maxConnections {
        maxConnections = MAX_CONNECTIONS
    }

    s := &BaseServer {
        Name:name,
        ServerId:serverid,
        Addr : addr,
        ServerType: serverType,
        MaxConns : maxConnections,
        Stop_: false,
        exitChan: make(chan bool),
        Wormholes: wm,
        RoutePackHandle: routePack,
    }


    /*
    if s.RoutePackHandle.GetEndian() == gts.BigEndian {
        s.endianer = binary.BigEndian
    } else {
        s.endianer = binary.LittleEndian
    }
    */

    s.idassign = gts.NewIDAssign(s.MaxConns)

    return s
}


func (s *BaseServer) Stop() {
    s.Stop_ = true
}


func (s *BaseServer) AllocId() int {
    if (s.Wormholes.Length() >= s.MaxConns) {
        return 0
    }

    return s.idassign.GetFree()
}


func (s *BaseServer) SetMaxConnections(max int) {
    if MAX_CONNECTIONS < max {
        s.MaxConns = MAX_CONNECTIONS
    } else {
        s.MaxConns = max
    }
}


