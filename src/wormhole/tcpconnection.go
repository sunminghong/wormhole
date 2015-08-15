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
    id TID
    conn net.Conn

    Stream *gts.RWStream
    DPSize  int
    RouteType byte

    routePack IRoutePack

    receivePacketCallback ReceivePacketFunc
    closeCallback CommonCallbackFunc

    //需要输出的数据包的channel
    outgoing chan *RoutePacket

    //需要输出的数据流 的channel
    outgoingBytes chan []byte

    quit chan bool
}


// new Transport object
func NewTcpConnection(newcid TID, conn net.Conn, routepack IRoutePack) *TcpConnection {
    c := &TcpConnection {
        id:      newcid,
        conn:     conn,
        routePack: routepack,

        outgoing: make(chan *RoutePacket, 1),
        outgoingBytes: make(chan []byte),
        quit:     make(chan bool),

        Stream:   gts.NewRWStream(1024,routepack.GetEndian()),
    }

    c.Stream.Reset()

    return c
}

func (c *TcpConnection) GetId() TID {
    return c.id
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

    defer conn.Close()

    c.conn = conn

    //创建go的线程 使用Goroutine
    go c.transportSender()
    go c.transportReader()

    gts.Info("be connected to grid ", addr)
    return true
}

func (c *TcpConnection) transportReader() {
    buffer := make([]byte, c.read_buffer_size)
    for {
        bytesRead, err := c.conn.Read(buffer)

        if err != nil {
            c.closeCallback(c.id)
            break
        }

        //gts.Trace("pool transportReader read to buff:", bytesRead)
        gts.Trace("pool transportReader read to buff:",bytesRead)
        c.Stream.Write(buffer[0:bytesRead])

        gts.Trace("tpool transportReader Buff:%d", len(c.Stream.Bytes()))
        n, dps := c.routePack.Fetch(c)
        gts.Trace("fetch message number", n)
        if n > 0 {
            c.receivePacketCallback(c.id, dps)
        }
    }
    //Log("TransportReader stopped for ", transport.Cid)
}

func (c *TcpConnection) transportSender() {
    for {
        select {
        case dp := <-c.outgoing:
            gts.Trace("clientpool transportSender:dp.type=%v,dp.data=% X",dp.Type, dp.Data)
            c.routePack.PackWrite(c.conn.Write,dp)

        case <-c.quit:
            //Log("Transport ", transport.Cid, " quitting")
            c.conn.Close()

            c.closeCallback(c.id)
            break
        }
    }
}


func (c *TcpConnection) SetReceivePacketCallback(cf ReceivePacketFunc)  {
    c.receivePacketCallback = cf
}

func (c *TcpConnection) SetCloseCallback(cf CommonCallbackFunc) {
    c.closeCallback = cf
}

func (c *TcpConnection) SetRoutePack(route IRoutePack) {
    c.routePack = route
}

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

func (c *TcpConnection) Send(dp *RoutePacket) {
    c.outgoing <- dp
}


