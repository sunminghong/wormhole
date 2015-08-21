/*=============================================================================
#     FileName: guin.go
#       Author: sunminghong, allen.fantasy@gmail.com, http://weibo.com/5d13
#         Team: http://1201.us
#   LastChange: 2015-08-14 17:52:19
#      History:
=============================================================================*/

/*
guin处理 的类
1. guin 即 globally unique identity number，给每一个数据通道一个全局唯一编号，比如“虫洞通道”
2. guin 数据包括:
    + agentid(0~4) ，连接所在的连接ageng的编号，用于数据发送时的路由
    + unique number（5~18），唯一id，1~16383，这个数字会回收重复分配
    + checkcode(19~31)，0~8191, 给这个通道一个验证码，用于区分重复的id
*/


package wormhole


import (
    "time"
    //"math/rand"

    gts "github.com/sunminghong/gotools"
)


func ParseGuin(guin int) (serverId int, id int, check int) {
    uin := int(guin)
    return uin & 0x1f, (uin >> 5) & 0x3fff, (uin >> 19) & 0x1fff
}


func GenerateGuin(serverId int, id int) int {
    //check := gm.r.Intn(0x1fff - 1)
    check := int(time.Now().UnixNano() & 0x1fff)

    //guin := (serverId & 0x1f) | ((id & 0x3fff) << 5) | (check & 0x1fff) << 19)
    guin := serverId | (id << 5) | (check << 19)

    return guin
}



type GuinMaker struct {
    idassign *gts.IDAssign

    //r *rand.Rand
}

func NewGuinMaker() *GuinMaker {
    gm := &GuinMaker {
        idassign : gts.NewIDAssign(2<<14 -1),
        //r :        rand.New(rand.NewSource(time.Now().UnixNano())),
    }

    return gm
}

func(gm *GuinMaker) GenerateGuin(serverId int) int {
    id := gm.idassign.GetFree()
    if id == 0 {
        return 0
    }

    return GenerateGuin(serverId, id)
}


