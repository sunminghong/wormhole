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
    "strconv"

    iniconfig "github.com/sunminghong/iniconfig"
    gts "github.com/sunminghong/gotools"
    gutils "github.com/sunminghong/gotools/utils"

    "wormhole/wormhole"
)


type ClientWormhole struct {
    *wormhole.Wormhole
}


func NewClientWormhole(guin int, manager wormhole.IWormholeManager, routepack wormhole.IRoutePack) wormhole.IWormhole {
    aw := &ClientWormhole {
        Wormhole: wormhole.NewWormhole(guin, manager, routepack),
    }
    aw.RegisterSub(aw)

    return aw
}


func (wh *ClientWormhole) Init() {
    gts.Trace("clientwormhole is init")

    //wh.Send(0, []byte("this message is from player 1 !"))
    //wh.Send(0, []byte("this message is from player 2 !"))
    //wh.Send(0, []byte("this message is from player 3 !"))
    //wh.Send(0, []byte("this message is from player 4 !"))
    //wh.Send(0, []byte("this message is from player 5 !"))
}

func (wh *ClientWormhole) ProcessPackets(dps []*wormhole.RoutePacket) {
    //gts.Trace("clientwormhole processpacket receive %d route packets",len(dps))
    for _, dp := range dps {
        gts.Trace(gutils.ByteString(dp.Data))
        //gts.Trace("%q", dp)
    }

}
//ClientWormhole end


var client *wormhole.Client


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

    section := "Default"

    clientTcpAddr, err := c.GetString(section, "clientTcpAddr")
    if err != nil {
        gts.Error(err.Error())
        return
    }

    var endianer gts.IEndianer
    endian, err := c.GetInt(section, "endian")
    if err == nil {
        endianer = gts.GetEndianer(endian)
    } else {
        endianer = gts.GetEndianer(gts.LittleEndian)
    }
    routepack := wormhole.NewRoutePack(endianer)
    var cwormholes wormhole.IWormholeManager

    client = wormhole.NewClient(clientTcpAddr, routepack, cwormholes,NewClientWormhole, wormhole.EWORMHOLE_TYPE_CLIENT)
    client.Connect()

    gts.Info("----------------client connect to %s,%s-----------------",clientTcpAddr)

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
        } else {
            wh := client.GetWormhole()
            for ii:=1;ii <= 10000; ii++ {
                wh.Send(0, 10 , []byte(" this message is from player " + strconv.Itoa(ii) + " !"))
            }
        }
    }
}


