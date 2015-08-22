/*=============================================================================
#     FileName: clientserver.go
#       Author: sunminghong, allen.fantasy@gmail.com, http://weibo.com/5d13
#         Team: http://1201.us
#   LastChange: 2015-08-21 12:35:56
#      History:
=============================================================================*/


/*
定义一个基本的逻辑服务器端

*/
package main


import (
	"bufio"
	"os"
    "flag"
    "runtime"

    iniconfig "github.com/sunminghong/iniconfig"
    gts "github.com/sunminghong/gotools"

    "wormhole/wormhole"
    "wormhole/server"

)


type ClientWormhole struct {
    *wormhole.ClientWormhole
}


func NewClientWormhole(guin int, manager wormhole.IWormholeManager, routepack wormhole.IRoutePack) wormhole.IWormhole {
    aw := &ClientWormhole {
        Wormhole: server.NewWormhole(guin, manager, routepack),
    }
    aw.RegisterSub()

    return aw
}


func (aw *ClientWormhole) ProcessPackets(dp []*wormhole.RoutePacket) {
    gts.Trace("clienttologic processpacket receive %d route packets",dp)
}
//ClientWormhole end


var client *server.Client


var (
    clientConf=flag.String("clientConf","client1.conf","client server config file")
)


func main() {
    runtime.GOMAXPROCS(runtime.NumCPU())

    flag.Parse()

    c, err := iniconfig.ReadConfigFile(*clientConf)
    if err != nil {
        gts.Error(err.Error())
        return
    }

    var endianer gts.IEndianer
    section := "Default"
    endian, err := c.GetInt(section, "endian")
    if err == nil {
        endianer = gts.GetEndianer(endian)
    } else {
        endianer = gts.GetEndianer(gts.LittleEndian)
    }
    routepack := wormhole.NewRoutePack(endianer)

    cwormholes := wormhole.NewWormholeManager(routepack, wormhole.MAX_CONNECTIONS,wormhole.EWORMHOLE_TYPE_CLIENT)

    dispatcher := server.NewDispatcher(routepack)
    lwormholes := server.NewLogicManager(routepack, dispatcher)
    client := server.NewClientFromIni(c,routepack,cwormholes,lwormholes,NewClientToClient,NewClientWormhole)

    client.Start()

    gts.Info("----------------client server is running!-----------------",client.GetServerId())

    quit := make(chan bool)
    go exit(quit)

    <-quit
}


func exit(quit chan bool) {
    for {
        r := bufio.NewReader(os.Stdin)
        ru, _, _ := r.ReadRune()
        print(ru)
        if ru == 115 {
            quit <- true
            return
        }
    }
}
