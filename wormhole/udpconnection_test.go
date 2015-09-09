/*=============================================================================
#     FileName: udpconnection.go
#       Author: sunminghong, allen.fantasy@gmail.com, http://weibo.com/5d13
#         Team: http://1201.us
#   LastChange: 2015-08-17 18:40:34
#      History:
=============================================================================*/


package wormhole

import (
    //"net"
    //"time"
    //"fmt"

    "testing"

    gts "github.com/sunminghong/gotools"
)

func addData(b *UdpConnection, no int) {
    b.recvFrame(&UdpFrame{
        OrderNo:no,
        Flag:UDPFRAME_FLAG_DATA,
    })
}


func Test_A(t *testing.T){
    b := NewUdpConnection(1,nil,gts.GetEndianer(gts.BigEndian),nil)
    b.SetReceiveCallback(func(conn IConnection) {
        println("receivecallback")
    })

    addData(b, 1)
    addData(b, 2)
    addData(b, 3)
    addData(b, 7)
    addData(b, 10)
    addData(b, 8)
    addData(b, 5)
    addData(b, 9)
    addData(b, 6)
    addData(b, 4)

    return
}

