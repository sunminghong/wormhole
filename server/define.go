/*=============================================================================
#     FileName: define.go
#       Author: sunminghong, allen.fantasy@gmail.com, http://weibo.com/5d13
#         Team: http://1201.us
#   LastChange: 2015-08-20 09:36:13
#      History:
=============================================================================*/


package server

import (
    . "wormhole/wormhole"
)


type ILogicDispatcher interface {
    AddHandler(rule []byte, wh IWormhole)
    Dispatch(dp *RoutePacket) IWormhole
    RemoveHandler(guin int)
}


