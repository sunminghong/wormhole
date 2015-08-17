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
    gts "github.com/sunminghong/gotools"
)

// Connection  
type UdpConnection struct {
    read_buffer_size int
    connectionType EConnType
    id int
    conn *net.UDPConn
    userAddr *net.UDPAddr

    stream *gts.RWStream
    DPSize  int
    RouteType byte

    receiveCallback ReceiveFunc
    closeCallback CommonCallbackFunc

    //需要输出的数据包的channel
    outgoing chan *RoutePacket

    receiveChan chan bool

    //需要输出的数据流 的channel
    outgoingBytes chan []byte

    quit chan bool
    Quit chan bool
}


// new Transport object
func NewUdpConnection(newcid int, conn *net.UDPConn, endian int, userAddr *net.UDPAddr) *UdpConnection {
    c := &UdpConnection {
        id:      newcid,
        conn:     conn,
        userAddr: userAddr,

        outgoing: make(chan *RoutePacket, 1),
        outgoingBytes: make(chan []byte),
        receiveChan: make(chan bool, 20),
        quit:     make(chan bool),
        Quit:     make(chan bool),

        stream:   gts.NewRWStream(1024, endian),
    }

    c.stream.Reset()

    //创建go的线程 使用Goroutine
    go c.ConnSender()
    go c.ConnReader()

    return c
}

func (c *UdpConnection) GetId() int {
    return c.id
}


func (c *UdpConnection) GetStream() IStream {
    return c.stream
}


func (c *UdpConnection) Connect(addr string) bool {
    gts.Info("connect to grid:", addr)

    conn, err := net.Dial("tcp", addr)
    if err != nil {
        gts.Warn("net.Dial to %s:%q",addr, err)
        return false
    } else {
        gts.Info("pool dial to %s is ok. ", addr)
    }

    go func() {
        defer conn.Close()

        c.conn = conn

        //创建go的线程 使用Goroutine
        go c.reader()
        go c.ConnReader()

        gts.Info("be connected to grid ", addr)

        <-c.Quit
    }()
    return true
}

func (c *UdpConnection) ConnReader(buffer []byte) {
    c.stream.Write(buffer)
    c.receiveChan <- true
}


func (c *UdpConnection) reader() {
    <-c.receiveChan
    c.receiveCallback(c)
}

func (c *UdpConnection) ConnSender() {
    for {
        select {
        //case dp := <-c.outgoing:
            //gts.Trace("clientpool ConnSender:dp.type=%v,dp.data=% X",dp.Type, dp.Data)
            ////c.routePack.PackWrite(c.conn.Write,dp)
        case bytes := <-c.outgoingBytes:
            c.conn.WriteToUDP(bytes, c.userAddr)

        case <-c.quit:
            //Log("Transport ", transport.Cid, " quitting")
            c.conn.Close()

            c.closeCallback(c.id)
            break
        }
    }
}


//func (c *UdpConnection) SetReceivePacketCallback(cf ReceivePacketFunc)  {
    //c.receivePacketCallback = cf
//}


func (c *UdpConnection) SetReceiveCallback(cf ReceiveFunc)  {
    c.receiveCallback = cf
}

func (c *UdpConnection) SetCloseCallback(cf CommonCallbackFunc) {
    c.closeCallback = cf
}

//func (c *UdpConnection) SetRoutePack(route IRoutePack) {
    //c.routePack = route
//}

func (c *UdpConnection) GetType() EConnType {
    return c.connectionType
}

func (c *UdpConnection) SetType(t EConnType) {
    c.connectionType = t
}

func (c *UdpConnection) Close() {
    c.quit <- true
    c.conn.Close()
}


/*
func (c *UdpConnection) Send(dp *RoutePacket) {
    c.outgoing <- dp
}
*/


func (c *UdpConnection) Send(data []byte) {
    c.outgoingBytes <- data
}


