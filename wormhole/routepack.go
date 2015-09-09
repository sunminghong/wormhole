/*=============================================================================
#     FileName: RoutePack.go
#       Author: sunminghong, allen.fantasy@gmail.com, http://weibo.com/5d13
#         Team: http://1201.us
#   LastChange: 2015-08-14 14:20:22
#      History:
=============================================================================*/


package wormhole

import (
    //"fmt"

    gts "github.com/sunminghong/gotools"
)

const (
    mask1 = byte(0x25)
    mask2 = byte(0x20)
)

type RoutePack struct {
    endianer gts.IEndianer
}

func NewRoutePack(endianer gts.IEndianer) *RoutePack {
    dg := &RoutePack{endianer : endianer}

    return dg
}

func (d *RoutePack) GetEndianer() gts.IEndianer {
    return d.endianer
}

func (d *RoutePack) Clone() IRoutePack {
    dg := &RoutePack{}

    dg.SetEndianer(d.endianer)
    return dg
}

func (d *RoutePack) SetEndianer(endianer gts.IEndianer) {
    d.endianer = endianer
}

func (d *RoutePack) encrypt(plan []byte){
    return
    for i,_ := range plan {
        plan[i] ^= 0x37
    }
}

func (d *RoutePack) decrypt(plan []byte){
    return
    for i,_ := range plan {
        plan[i] ^= 0x37
    }
}


//flag1(byte)+flag2(byte)+datatype(byte)+method(short)+[+guin(int32)]+data(datasize(int32)+body)
//对数据进行拆包
//func (d *RoutePack) fetchTcp(ci IConnection) (n int, dps []*RoutePacket) {
func (d *RoutePack) Fetch(c *ConnectionBuffer) (n int, dps []*RoutePacket) {
    dps = []*RoutePacket{}

    //c := ci.(*TcpConnection)
    cs := c.Stream
    ilen := cs.Len()
    if ilen == 0 {
        return
    }

    var dpSize, guin int
    var routeType byte

    for {
        //pos := cs.GetPos()

        //拆包
        if c.DPSize > 0 {
            if ilen < c.DPSize {
                //如果缓存去数据长度不够就退出接着等后续数据
                return
            }
            dpSize = c.DPSize
            routeType = c.RouteType
            guin = c.Guin
        } else {
            //Trace("ilen,pos:%d,%d",ilen,pos)
            if ilen < 9 {
                return
            }

            head,_ := cs.Read(9)
            ilen = cs.Len()
            d.decrypt(head)

            if head[0]==mask1 && head[1]==mask2 {
                routeType = head[2]
                _dpSize := 0

                method := int(d.endianer.Uint16(head[3:5]))
                if routeType & 1 == 1 {
                    guin = int(d.endianer.Uint32(head[5:]))
                    if ilen < 4 {
                        cs.SetPos(-9)
                        return
                    }
                    head2,_ := cs.Read(4)
                    ilen = cs.Len()

                    _dpSize = int(d.endianer.Uint32(head2))
                } else {
                    _dpSize = int(d.endianer.Uint32(head[5:]))
                }

                dpSize = _dpSize

                //pos = cs.GetPos()
                //Trace("ilen,pos,dpSize",ilen,pos,dpSize)
                if ilen < dpSize {
                    c.Method = method
                    c.Guin = guin
                    c.DPSize = dpSize
                    c.RouteType = routeType

                    return
                }
            } else {
                //如果错位则将缓存数据抛弃
                cs.Reset()
                return
            }
        }

        data,size := cs.Read(dpSize)
        if size > 0 {
            newdata := make([]byte, dpSize)
            copy(newdata, data)

            dp := &RoutePacket{
                Type:ERouteType(routeType),
                Guin: guin,
                Data: newdata,
                }
            /*
            if routeType & 1 == 1 {
                //dp.Guin = int(d.endianer.Uint32(data[dpSize-4:]))
                //dp.Data = data[:dpSize-4]
                dp.Guin = int(d.endianer.Uint32(data))
                dp.Data = data[4:]
            } else {
                dp.Data = data
            }
            dp.Guin = guin
            */

            dps = append(dps,dp)
            n += 1
        }

        c.DPSize = 0
        c.RouteType = 0

        ilen = cs.Len()
        if ilen >= 9 {
            continue
        }

        /*
        iiii := ilen - cs.GetPos()
        if iiii >= 9 {
            continue
        }*/

        if ilen == 0 {
            cs.Reset()
        }

        return
    }
    return
}


//对数据进行封包
func (d *RoutePack) packHeader(dp *RoutePacket) []byte {
    glen := 0
    if dp.Type & 1 == 1 {
        glen = TIDSize
    }
    buf := make([]byte, 9 + glen)

    buf[0] = byte(mask1)
    buf[1] = byte(mask2)
    buf[2] = byte(dp.Type)
    d.endianer.PutUint16(buf[3:], uint16(dp.Method))
    ilen := len(dp.Data)

    if dp.Type & 1 == 1 {
        d.endianer.PutUint32(buf[5:], uint32(dp.Guin))
    }
    d.endianer.PutUint32(buf[5 + glen:], uint32(ilen))

    d.encrypt(buf)

    return buf
}


//对数据进行封包
func (d *RoutePack) Pack(dp *RoutePacket) []byte {
    head := d.packHeader(dp)

    ilen := len(dp.Data)
    glen := 0
    if dp.Type & 1 == 1 {
        glen = TIDSize
    }
    buf := make([]byte, 9 + glen + ilen)
    copy(buf,head)
    copy(buf[9 + glen:], dp.Data)
    return buf
}


//对数据进行封包
func (d *RoutePack) PackWrite(write WriteFunc,dp *RoutePacket) {
    head := d.packHeader(dp)

    write(head)
    write(dp.Data)
}


