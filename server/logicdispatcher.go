/*=============================================================================
#     FileName: logicdispatcher.go
#       Author: sunminghong, allen.fantasy@gmail.com, http://weibo.com/5d13
#         Team: http://1201.us
#   LastChange: 2015-08-20 09:51:01
#      History:
=============================================================================*/
package server

import (
    "sync"
    "strconv"
    "strings"

    gts "github.com/sunminghong/gotools"
    gutils "github.com/sunminghong/gotools/utils"
    . "wormhole/wormhole"
)

type Dispatcher struct {
    handlers map[int][]IWormhole
    routepack IRoutePack

    dlock *sync.RWMutex
}


func NewDispatcher(routepack IRoutePack) *Dispatcher {
    r := &Dispatcher{
        handlers :      make(map[int][]IWormhole),
        routepack:      routepack,
        dlock   :       new(sync.RWMutex),
    }
    return r
}


//rule = 10,12,0
func (r *Dispatcher) AddHandler(rule []byte, wh IWormhole) {
    rules := gutils.ByteString(rule)

    rules = strings.Replace(rules," ","",-1)
    if len(rules) ==0 {
        r.addRule(0,wh)
        return
    }

    groups:= strings.Split(rules,",")
    gts.Trace("add disp",groups)
    for _,p_ := range groups {
        p := strings.Trim(p_," ")
        if len(p) == 0 {
            continue
        }
        gcode, err := strconv.Atoi(p)
        if err ==nil {
            r.addRule(gcode,wh)
        }
    }
    gts.Trace("messagecodemaps1:",r.handlers)
}


func (r *Dispatcher) Dispatch(dp *RoutePacket) (wh IWormhole) {
    r.dlock.Lock()
    defer r.dlock.Unlock()

    var code, group int
    if len(dp.Data) > 2 {
        code = int(r.routepack.GetEndianer().Uint16(dp.Data))
        group = int(code / 100)
        gts.Trace("msg.group:",group)
    } else {
        group = 0
    }

    hands, ok := r.handlers[group]
    if !ok {
        hands, ok = r.handlers[0]
        if !ok {
            hands = []IWormhole{}
        }
    }

    hlen := len(hands)
    if hlen == 1 {
        return hands[0]
    } else if hlen > 1 {
        return hands[int(dp.Guin) % hlen]
    } else {
        gts.Warn("data packet group is not exists!", code, group)
        return nil
    }
}


func (r *Dispatcher) addRule(group int, wh IWormhole) {
    r.dlock.Lock()
    defer r.dlock.Unlock()

    if hands, ok := r.handlers[group]; ok {
        r.handlers[group] = append(hands, wh)
    } else {
        r.handlers[group] = []IWormhole{ wh }
    }
}


func (r *Dispatcher) RemoveHandler(guin int) {
    for group, hands := range r.handlers {
        hs := []IWormhole{}
        for _, hand := range hands {
            if hand.GetGuin() != guin {
                hs = append(hs, hand)
            }
        }

        r.handlers[group] = hs
    }
}


