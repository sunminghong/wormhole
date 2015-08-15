/*=============================================================================
#     FileName: RoutePack.go
#       Author: sunminghong, allen.fantasy@gmail.com, http://weibo.com/5d13
#         Team: http://1201.us
#   LastChange: 2015-08-14 14:20:22
#      History:
=============================================================================*/


/*

*/

package wormhole

import (
    "encoding/binary"
    gts "github.com/sunminghong/gotools"
)

const (
    mask1 = byte(0x25)
    mask2 = byte(0x20)
)

type RoutePack struct {
    endian int
    Endianer gts.IEndianer
}

func NewRoutePack(endian int ) *RoutePack{
    dg := &RoutePack{}

    dg.SetEndian(endian)
    return dg
}

func (d *RoutePack) GetEndian() int {
    return d.endian
}

func (d *RoutePack) Clone() IRoutePack {
    dg := &RoutePack{}

    dg.SetEndian(d.endian)
    return dg
}

func (d *RoutePack) SetEndian(endian int) {
    d.endian = endian
    if endian == gts.BigEndian {
        d.Endianer = binary.BigEndian
    } else {
        d.Endianer = binary.LittleEndian
    }
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


//flag1(byte)+flag2(byte)+datatype(byte)+data(datasize(int32)+body)+fromcid(int32)
//对数据进行拆包
func (d *RoutePack) Fetch(c IConnection) (n int, dps []*RoutePacket) {
    return d.fetchTcp(c)
}

func (d *RoutePack) fetchTcp(ci IConnection) (n int, dps []*RoutePacket) {
    dps = []*RoutePacket{}

    c:= ci.(*TcpConnection)
    cs := c.Stream
    ilen := cs.Len()
    if ilen == 0 {
        return
    }

    var dpSize int
    //var m1,m2 byte
    var dataType byte

    for {
        pos := cs.GetPos()
        //Log("pos:",pos)

        //拆包
        if c.DPSize > 0 {
            if ilen-pos < c.DPSize {
                //如果缓存去数据长度不够就退出接着等后续数据
                return
            }
            dpSize = c.DPSize
            dataType = c.RouteType
        } else {
            //Trace("ilen,pos:%d,%d",ilen,pos)
            if ilen-pos < 7 {
                return
            }

            head,_ := cs.Read(7)
            d.decrypt(head)

            /*
            cs.SetPos(-7)
            m1,_ := cs.ReadByte()
            m2,_ := cs.ReadByte()
            //Trace("m1,m2",m1,m2)
            if m1==mask1 && m2==mask2 {
                dataType,_ = cs.ReadByte()
                _dpSize,err := cs.ReadUint32()

                if err != nil {
                    cs.Reset()
                    c.DPSize = 0
                    c.RouteType = 0
                    return 0,nil
                }
            */

            if head[0]==mask1 && head[1]==mask2 {
                dataType = head[2]
                _dpSize := d.Endianer.Uint32(head[3:])
                //Trace("dataType,dpSize,endian",dataType,_dpSize,cs.Endian)

                dpSize = int(_dpSize)
                if dataType & 1 == 1 {
                    dpSize += 4
                }

                pos = cs.GetPos()
                //Trace("ilen,pos,dpSize",ilen,pos,dpSize)
                if ilen - pos < dpSize {
                    c.DPSize = dpSize
                    c.RouteType = dataType

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
            dp := &RoutePacket{Type:ERouteType(dataType)}

            if dataType & 1 == 1 {
                dp.Guin = TID(d.Endianer.Uint32(data[dpSize-4:]))
                dp.Data = data[:dpSize-4]
            } else {
                dp.Data = data
            }

            dps = append(dps,dp)
            n += 1
        }

        c.DPSize = 0
        c.RouteType = 0

        iiii := ilen - cs.GetPos()
        if iiii > 7 {
            continue
        }

        if iiii == 0 {
            cs.Reset()
        }

        return
    }
    return
}


//对数据进行封包
func (d *RoutePack) packHeader(dp *RoutePacket) []byte {
    ilen := len(dp.Data)

    if dp.Type & 1 == 1 {
        ilen += TIDSize
    }
    buf := make([]byte, ilen+7)

    buf[0] = byte(mask1)
    buf[1] = byte(mask2)
    buf[2] = byte(dp.Type)

    d.Endianer.PutUint32(buf[3:], uint32(ilen))

    d.encrypt(buf)

    return buf
}


//对数据进行封包
func (d *RoutePack) Pack(dp *RoutePacket) []byte {
    head := d.packHeader(dp)

    ilen := len(dp.Data)
    if dp.Type & 1 == 1 {
        ilen += TIDSize
    }
    buf := make([]byte, 7 + ilen)
    copy(buf,head)
    copy(buf[7:], dp.Data)

    if dp.Type & 1 == 1 {
        d.Endianer.PutUint32(buf[7 + ilen - TIDSize:], uint32(dp.Guin))
    }
    return buf
}


//对数据进行封包
func (d *RoutePack) PackWrite(write WriteFunc,dp *RoutePacket) {
    head := d.packHeader(dp)

    write(head)
    write(dp.Data)

    if dp.Type & 1 == 1 {
        cbuf := make([]byte,4)
        d.Endianer.PutUint32(cbuf, uint32(dp.Guin))
        write(cbuf)
    }
}


