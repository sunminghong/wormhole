/*=============================================================================
#     FileName: logic.go
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
    "strconv"

    iniconfig "github.com/sunminghong/iniconfig"
    gts "github.com/sunminghong/gotools"

    . "wormhole/wormhole"
)


type Logic struct {
    makeWormhole NewWormholeFunc
    wormholes IWormholeManager

    routepack IRoutePack

    serverType EWormholeType
    serverId int

    logics map[string]*Client

    group string
}


func NewLogic(
    name string, serverId int,
    routepack IRoutePack, wm IWormholeManager,
    makeWormhole NewWormholeFunc, group string) *Logic {

    ls := &Logic{
        routepack: routepack,

        serverId: serverId,
        makeWormhole: makeWormhole,
        wormholes: wm,
        serverType: EWORMHOLE_TYPE_LOGIC,
        group : group,

        logics: make(map[string]*Client),
    }
    ls.wormholes.SetServer(ls)

    return ls
}


func (ls *Logic) GetServerId() int {
    return ls.serverId
}


func (ls *Logic) ConnectAgent(tcpAddr string, udpAddr string) {
    c := NewClient(tcpAddr, udpAddr, ls.routepack, ls.wormholes,
        ls.makeWormhole, ls.serverType)
    c.Connect()

    ls.logics[tcpAddr] = c
}


func (ls *Logic) GetGroup() string {
    return ls.group
}


func (ls *Logic) Close() {
    //for c := range ls.logics {
        //c.Close()
    //}

    //ls.logics.Clear()

    ls.wormholes.CloseAll()
}


func NewLogicFromIni(
    c *iniconfig.ConfigFile,
    routepack IRoutePack, wm IWormholeManager,
    makeWormhole NewWormholeFunc) *Logic {

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

    group, err := c.GetString(section, "group")
    if err != nil {
        group = "0"
    }

    /*
    endian, err := c.GetInt(section, "endian")
    if err == nil {
        routepack.SetEndianer(gts.GetEndianer(endian))
    } else {
        routepack.SetEndianer(gts.GetEndianer(gts.LittleEndian))
    }

    autoDuration, err := c.GetInt(section, "autoReconnectDuration")
    if err != nil {
        autoDuration = 5
    }
    autoReconnectDuration := time.Duration(autoDuration) * time.Second
    */

    ls := NewLogic( name, serverId, routepack, wm, makeWormhole, group)

    return ls
}


func (ls *Logic) ConnectFromIni(c *iniconfig.ConfigFile) {
    gts.Trace("connect from ini")
    //make some connection to game server
    for i := 1; i < 50; i++ {
        section := "Agent" + strconv.Itoa(i)
        if !c.HasSection(section) {
            continue
        }

        enabled, err := c.GetBool(section, "enabled")
        if err == nil && !enabled {
            continue
        }

        /*
        serverId, err := c.GetInt(section, "serverId")
        if err != nil {
            gts.Error(err.Error())
            continue
        }

        gname, err := c.GetString(section, "name")
        if err != nil {
            gts.Error(err.Error())
            continue
        }
        */

        tcpAddr, err := c.GetString(section, "tcpAddr")
        if err != nil {
            gts.Error(err.Error())
            continue
        }

        udpAddr, err := c.GetString(section, "udpAddr")
        if err != nil {
            gts.Warn(err.Error())
        }

        ls.ConnectAgent(tcpAddr, udpAddr)
    }
}


