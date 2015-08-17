/*=============================================================================
#     FileName: tcpconnection.go
#       Author: sunminghong, allen.fantasy@gmail.com, http://weibo.com/5d13
#         Team: http://1201.us
#   LastChange: 2015-08-14 11:40:39
#      History:
=============================================================================*/


package wormhole

import (
    "net"
    gts "github.com/sunminghong/gotools"
)

// Connection  
type TcpConnection struct {
    read_buffer_size int
    connectionType EConnType
    id int
    conn net.Conn

    stream *gts.RWStream
    DPSize  int
    RouteType byte

    //receivePacketCallback ReceivePacketFunc
    receiveCallback ReceiveFunc
    closeCallback CommonCallbackFunc

    //需要输出的数据包的channel
    outgoing chan *RoutePacket

    //需要输出的数据流 的channel
    outgoingBytes chan []byte

    quit chan bool
    Quit chan bool
}


// new Transport object
func NewTcpConnection(newcid int, conn net.Conn, endian int) *TcpConnection {
    c := &TcpConnection {
        id:      newcid,
        conn:     conn,

        outgoing: make(chan *RoutePacket, 1),
        outgoingBytes: make(chan []byte),
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

func (c *TcpConnection) GetId() int {
    return c.id
}


func (c *TcpConnection) GetStream() IStream {
    return c.stream
}


func (c *TcpConnection) Connect(addr string) bool {
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
        go c.ConnSender()
        go c.ConnReader()

        gts.Info("be connected to grid ", addr)

        <-c.Quit
    }()
    return true
}

func (c *TcpConnection) ConnReader() {
    buffer := make([]byte, c.read_buffer_size)
    for {
        bytesRead, err := c.conn.Read(buffer)

        if err != nil {
            c.closeCallback(c.id)
            break
        }

        //gts.Trace("pool ConnReader read to buff:", bytesRead)
        gts.Trace("pool ConnReader read to buff:",bytesRead)
        c.stream.Write(buffer[0:bytesRead])
        c.receiveCallback(c)

        //gts.Trace("tpool ConnReader Buff:%d", len(c.stream.Bytes()))

        //n, dps := c.routePack.Fetch(c)
        //gts.Trace("fetch message number", n)
        //if n > 0 {
            //c.receivePacketCallback(c.id, dps)
        //}
    }
    //Log("TransportReader stopped for ", transport.Cid)
}

func (c *TcpConnection) ConnSender() {
    for {
        select {
        //case dp := <-c.outgoing:
            //gts.Trace("clientpool ConnSender:dp.type=%v,dp.data=% X",dp.Type, dp.Data)
            ////c.routePack.PackWrite(c.conn.Write,dp)
        case bytes := <-c.outgoingBytes:
            c.conn.Write(bytes)

        case <-c.quit:
            //Log("Transport ", transport.Cid, " quitting")
            c.conn.Close()

            c.closeCallback(c.id)
            break
        }
    }
}


//func (c *TcpConnection) SetReceivePacketCallback(cf ReceivePacketFunc)  {
    //c.receivePacketCallback = cf
//}


func (c *TcpConnection) SetReceiveCallback(cf ReceiveFunc)  {
    c.receiveCallback = cf
}

func (c *TcpConnection) SetCloseCallback(cf CommonCallbackFunc) {
    c.closeCallback = cf
}

//func (c *TcpConnection) SetRoutePack(route IRoutePack) {
    //c.routePack = route
//}

func (c *TcpConnection) GetType() EConnType {
    return c.connectionType
}

func (c *TcpConnection) SetType(t EConnType) {
    c.connectionType = t
}

func (c *TcpConnection) Close() {
    c.quit <- true
    c.conn.Close()
}


/*
func (c *TcpConnection) Send(dp *RoutePacket) {
    c.outgoing <- dp
}
*/


func (c *TcpConnection) Send(data []byte) {
    c.outgoingBytes <- data
}


