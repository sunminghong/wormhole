/*=============================================================================
#     FileName: udpconnection.go
#       Author: sunminghong, allen.fantasy@gmail.com, http://weibo.com/5d13
#         Team: http://1201.us
#   LastChange: 2015-08-17 18:40:34
#      History:
=============================================================================*/


package wormhole

import (
    "net"
    "time"

    gts "github.com/sunminghong/gotools"
)

//TODO: 没有考虑udp 分包
const UDP_SEND_ACT_COUNT int = 31
const UDP_REQ_CHECK_COUNT int = 10

const (
    UDPFRAME_FLAG_DATA byte = iota
    UDPFRAME_FLAG_ACT
    UDPFRAME_FLAG_REQ_RETRY
    UDPFRAME_FLAG_NOT_EXISTS

    UDPFRAME_FLAG_DATA_GROUP
)


// UdpFrame start
// bytes = orderno(2)+flag(1)+length(2)+data
type UdpFrame struct {
    OrderNo int
    Flag byte
    Data []byte
    Count int
    //Buf []byte
}

func (uf *UdpFrame) ToString() {
    //gts.Trace("% X", uf.Buf)
    gts.Trace("%d, %d %d, % X", uf.OrderNo, uf.Flag, uf.Count, uf.Data)
}


func (uf *UdpFrame) Pack(endianer gts.IEndianer) []byte {
    buf := make([]byte, 2 + 1 + 2 + len(uf.Data))
    endianer.PutUint16(buf, uint16(uf.OrderNo))
    buf[2] = uf.Flag

    switch uf.Flag {
    case UDPFRAME_FLAG_DATA:
        endianer.PutUint16(buf[3:], uint16(len(uf.Data)))
        copy(buf[5:], uf.Data)
        return buf
    case UDPFRAME_FLAG_DATA_GROUP:
        return buf
    default:
        return buf[:3]
    }
}


func (uf *UdpFrame) Unpack(rw *gts.RWStream) bool {
    data,num := rw.Read(3)
    if num == 0 {
        return false
    }

    uf.OrderNo = int(rw.Endianer.Uint16(data[:2]))
    uf.Flag = data[2]
    uf.Count = 0

    switch uf.Flag {
    case UDPFRAME_FLAG_DATA:
        if rw.Len() < 2 {
            rw.SetPos(-3)
            gts.Trace("setpos:-3")
            return false
        }

        buf, n := rw.Read(2)
        if n < 2 {
            return false
        }
        length := rw.Endianer.Uint16(buf)
        leng := int(length)
        //length,err := rw.ReadUint16()
        gts.Trace("udpframe.unpack:%d", leng)
        //if err != nil || rw.Len() < int(leng) {
        if rw.Len() < int(leng) {
            rw.SetPos(-5)
            gts.Trace("setpos:-5, %d", leng)
            return false
        }

        data,_ := rw.Read(int(leng))
        uf.Data = make([]byte, leng)
        copy(uf.Data, data)

        /*
        uf.Buf = make([]byte, leng + 5)
        copy(uf.Buf, data)
        copy(uf.Buf[3:], buf)
        copy(uf.Buf[5:], uf.Data)
        */
        return true

    case UDPFRAME_FLAG_DATA_GROUP:
        return false

    default:
        return true
    }
}


// Connection  
type UdpConnection struct {
    *ConnectionBuffer

    //udp frame define
    sendNo          int
    lastValidOrderNo int
    lastOrderNo     int
    udpStream       *gts.RWStream
    sendCache       *gts.Map
    recvCache       *gts.Map
    reqCache        *gts.Map

    read_buffer_size int

    protocolType    EProtocolType
    connectionType  EConnType
    id              int
    conn            *net.UDPConn
    userAddr        *net.UDPAddr

    receiveCallback ReceiveFunc
    closeCallback   CommonCallbackFunc

    //需要输出的数据包的channel
    outgoing        chan *RoutePacket

    receiveChan     chan bool

    //需要输出的数据流 的channel
    outgoingBytes   chan []byte
    outFrame        chan *UdpFrame

    closeded        bool

    quitInterval    chan bool
    quitSender      chan bool
    quitConnect     chan bool
}


// new Transport object
func NewUdpConnection(newcid int, conn *net.UDPConn, endianer gts.IEndianer, userAddr *net.UDPAddr) *UdpConnection {
    c := &UdpConnection {
        ConnectionBuffer: &ConnectionBuffer{Stream:   gts.NewRWStream(1024, endianer)},
        id:             10000 + newcid,
        conn:           conn,
        userAddr:       userAddr,
        udpStream:      gts.NewRWStream(1024, endianer),

           sendNo:      0,
 lastValidOrderNo:      1,
      lastOrderNo:      1,
        sendCache:      gts.NewMap(),
        recvCache:      gts.NewMap(),
        reqCache:       gts.NewMap(),

        //outgoing:       make(chan *RoutePacket, 5),
        //outgoingBytes:  make(chan []byte, 20),
        outFrame:       make(chan *UdpFrame, 20),
        receiveChan:    make(chan bool, 1),

        closeded:       false,
        quitInterval:   make(chan bool),
        quitSender:     make(chan bool),
        quitConnect:    make(chan bool),
        protocolType :  EPROTOCOL_TYPE_UDP,
    }

    go c.interval()

    //创建go的线程 使用Goroutine
    go c.reader()

    if conn != nil {
        go c.ConnSenderServer()
    }

    return c
}


func (c *UdpConnection) interval() {
	updateChan := time.NewTicker(20 * time.Millisecond)
	for {
		select {
		//case s := <-c.sendChan:
			//if !c.closed {
				//b := []byte(s)
				//ikcp.Ikcp_send(c.kcp, b, len(b))
			//}
		case <-updateChan.C:
            c.reqCheck()

		case <-c.quitInterval:
            gts.Trace("interval")
            goto end
		}
	}
    end:
	updateChan.Stop()
}


func (c *UdpConnection) reqCheck() {
    reqCache := c.reqCache.All()
    for no, val := range reqCache {
        if val == nil {
            gts.Trace("reqFrame is nil:%d", no)
            continue
        }

        rframe := val.(*UdpFrame)
        rframe.Count -= 1
        if rframe.Count <= 0 {
            gts.Trace("send req:", no, len(reqCache))
            //gts.Trace("rframe:::%d,%d,%d", rframe.OrderNo, rframe.Flag,rframe.Count)

            //发送 重传请求
            rframe.Count = UDP_REQ_CHECK_COUNT
            c.sendFrame(rframe)

            gts.Trace("c.lastValidOrderNo:%d",c.lastValidOrderNo)
            gts.Trace("c.lastOrderNo:%d",c.lastOrderNo)
            gts.Trace("c.reqFrame.orderno:%d,flag:%d",rframe.OrderNo, rframe.Flag)
        }
    }
}


func (c *UdpConnection) recvFrame(frame *UdpFrame) {
    //取消期望验证包
    if frame.Flag == UDPFRAME_FLAG_NOT_EXISTS {
        if frame.OrderNo != c.lastOrderNo {
            c.reqCache.Delete(frame.OrderNo)
        } else {
            return
        }
    } else {
        c.reqCache.Delete(frame.OrderNo)
    }

    //lastOrderNo 接受到的最大 orderno + 1
    //lastValidOrderNo, 最后一个有效包的 orderno + 1

    //最后一个有效包orderno == frame.OrderNo，就直接发送
    if c.lastValidOrderNo  == frame.OrderNo {
        gts.Trace("recvFrame1111:%d", frame.OrderNo)
        frame.ToString()
        c.receiveBytes(frame)

        //如果lastorderno 与 lastvalidOrderNo相等，
        c.lastValidOrderNo += 1

        //则将lastorderno 等于 orderno
        if c.lastOrderNo < c.lastValidOrderNo {
            c.lastOrderNo = c.lastValidOrderNo
            c.addReq(10)
            return
        }

        //将lastValidorderno 开始的连续recvCache frame直接发送
        for {
            if fr, ok := c.recvCache.Get(c.lastValidOrderNo);ok {
                ffr := fr.(*UdpFrame)
                gts.Trace("recvFrame2222:%d", ffr.OrderNo)
                ffr.ToString()
                c.receiveBytes(ffr)
                c.recvCache.Delete(c.lastValidOrderNo)

                c.lastValidOrderNo += 1
            } else {
                break
            }
        }

        return
    }

    //插入到recvBuffer
    c.recvCache.Set(frame.OrderNo, frame)

    if c.lastOrderNo <= frame.OrderNo {
        for c.lastOrderNo < frame.OrderNo + 1 {
            c.lastOrderNo += 1
            gts.Trace("recvFrame3333:%d", frame.OrderNo)
            c.addReq(1)
        }
        return
    }
}


func (c *UdpConnection) sendFrame(frame *UdpFrame) {
    c.outFrame <- frame
}


func (c *UdpConnection) addReq(rate int) {
    rate = 1
    gts.Trace("addreq:::%d", c.lastOrderNo)
    c.reqCache.Set(c.lastOrderNo, &UdpFrame{
        OrderNo: c.lastOrderNo,
        //Data: make([]byte),
        Flag: UDPFRAME_FLAG_REQ_RETRY,
        Count: UDP_REQ_CHECK_COUNT * rate,
    })
}


func (c *UdpConnection) GetId() int {
    return c.id
}


func (c *UdpConnection) GetBuffer() *ConnectionBuffer {
    return c.ConnectionBuffer
}


func (c *UdpConnection) Connect(addr string) bool {
	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
        gts.Error("dial udp server addr(%s) is error:%q", udpAddr, err)
		return false
	}

    conn, err := net.DialUDP("udp", nil, udpAddr)
    if err != nil {
        gts.Warn("net.Dial to %s:%q",addr, err)
        return false
    }

    gts.Trace("dial to udp(%s) is ok.", addr)

    c.conn = conn
    go func() {

        go func() {
            buffer := make([]byte, 1024)
            for {
                defer conn.Close()
                if c.closeded {
                    return
                }

                n, err := conn.Read(buffer[0:])
                if err == nil {
                    c.ConnReader(buffer[0:n])
                } else {
                    e, ok := err.(net.Error)
                    if !ok || !e.Timeout() {
                        gts.Trace("recv errorconn:%q.", err.Error())
                        c.quitConnect <- true
                        return
                    }
                }
            }
        }()

        go c.ConnSenderClient()

        <-c.quitConnect
    }()
    return true
}


func (c *UdpConnection) ConnReader(buffer []byte) {
    //gts.Trace("udpConnReader read to buff:%d, % X",len(buffer), buffer)
    //gts.Trace("udpConnReader read to buff:%q",buffer)

    c.udpStream.Write(buffer)
    c.receiveChan <- true
}


func (c *UdpConnection) receiveBytes(frame *UdpFrame) {
    if frame.OrderNo % UDP_SEND_ACT_COUNT == 0 {
        c.sendFrame(&UdpFrame{
            OrderNo: frame.OrderNo,
            Flag:UDPFRAME_FLAG_ACT,
        })
        gts.Trace("c.lastValidOrderNo:%d",c.lastValidOrderNo)
        gts.Trace("c.lastOrderNo:%d",c.lastOrderNo)
        gts.Trace("c.reqCache:%d,%d",frame.OrderNo, UDPFRAME_FLAG_ACT)
    }

    println("receivebytes:", frame.OrderNo, "\n")
    c.Stream.Write(frame.Data)
    c.receiveCallback(c)
}


func (c *UdpConnection) reader() {
    for {
        <-c.receiveChan

        for {
            frame := &UdpFrame{}
            if frame.Unpack(c.udpStream) {
                //gts.Trace("%d,%d,%q", frame.OrderNo, frame.Flag, frame.Data)
                //gts.Trace("////////////////////////////////////")
                switch frame.Flag {
                case UDPFRAME_FLAG_ACT:
                    //收到包确认frame，将相应的sendcache 删除
                    for i := frame.OrderNo; i>1; i-- {
                        if _, ok := c.sendCache.Get(i);ok {
                            c.sendCache.Delete(i)
                        } else {
                            break
                        }
                    }
                    gts.Trace("-------------%d--------------",c.sendCache.Length())

                case UDPFRAME_FLAG_REQ_RETRY:
                    gts.Trace("----------recv req:%d------------",frame.OrderNo)
                    if rframe, ok := c.sendCache.Get(frame.OrderNo);ok {
                        c.sendFrame(rframe.(*UdpFrame))
                    } else {
                        c.sendFrame(&UdpFrame{
                            OrderNo:frame.OrderNo,
                            Flag:UDPFRAME_FLAG_NOT_EXISTS,
                            })
                    }

                case UDPFRAME_FLAG_NOT_EXISTS:
                    //如果该包丢失，就当接受到一个正常frame处理
                    frame.Data= []byte{}
                    c.recvFrame(frame)

                default:
                    c.recvFrame(frame)

                }
            } else {
                break
            }
        }
    }
}


func (c *UdpConnection) ConnSenderClient() {
    for {
        select {
        //case bytes := <-c.outgoingBytes:
            //c.conn.Write(bytes)
        case frame := <-c.outFrame:
            gts.Trace("1outframe:orderno:%d, flag:%d, data:%q", frame.OrderNo, frame.Flag, frame.Data)
            frame.ToString()
            bytes := frame.Pack(c.udpStream.Endianer)
            c.conn.Write(bytes)

        case <-c.quitSender:
            gts.Trace("connsenderclient,quit")
            c.conn.Close()

            c.closeCallback(c.id)
            break
        }
    }
}


func (c *UdpConnection) ConnSenderServer() {
    for {
        select {
        //case dp := <-c.outgoing:
            //gts.Trace("clientpool ConnSender:dp.type=%v,dp.data=% X",dp.Type, dp.Data)
            ////c.routePack.PackWrite(c.conn.Write,dp)
        case frame := <-c.outFrame:
            gts.Trace("2outframe:orderno:%d, flag:%d, data:%q", frame.OrderNo, frame.Flag, frame.Data)
            frame.ToString()
            bytes := frame.Pack(c.udpStream.Endianer)
            c.conn.WriteToUDP(bytes, c.userAddr)

        case <-c.quitSender:
            gts.Trace("udp connsenderserver,quit")
            c.conn.Close()

            c.closeCallback(c.id)
            break
        }
    }
}


func (c *UdpConnection) SetReceiveCallback(cf ReceiveFunc)  {
    gts.Trace("udp connection setReceiveCallback")
    c.receiveCallback = cf
}


func (c *UdpConnection) SetCloseCallback(cf CommonCallbackFunc) {
    c.closeCallback = cf
}


func (c *UdpConnection) GetProtocolType() EProtocolType {
    return c.protocolType
}


func (c *UdpConnection) GetType() EConnType {
    return c.connectionType
}


func (c *UdpConnection) SetType(t EConnType) {
    c.connectionType = t
}


func (c *UdpConnection) Close() {
    gts.Trace("udp connection close1")
    if !c.closeded {
        c.closeded = true
        c.quitInterval <- true
        c.quitSender <- true
    }
    gts.Trace("udp connection close2")
}


func (c *UdpConnection) Send(data []byte) {
    c.sendNo += 1
    uframe := &UdpFrame{
        OrderNo: c.sendNo,
        Data: data,
        Flag:UDPFRAME_FLAG_DATA,
        Count:0,
    }

    c.sendFrame(uframe)
    c.sendCache.Set(uframe.OrderNo, uframe)
}


