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
    *ConnectionBuffer

    read_buffer_size int
    protocolType EProtocolType
    connectionType EConnType
    id int
    conn net.Conn

/*
    stream *gts.RWStream
    Guin int
    DPSize  int
    RouteType byte
    */

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
func NewTcpConnection(newcid int, conn net.Conn, endianer gts.IEndianer) *TcpConnection {
    c := &TcpConnection {
        ConnectionBuffer: &ConnectionBuffer{Stream:   gts.NewRWStream(1024, endianer)},
        id:      newcid,
        conn:     conn,
        read_buffer_size: 1024,

        outgoing: make(chan *RoutePacket, 1),
        outgoingBytes: make(chan []byte),
        quit:     make(chan bool),
        Quit:     make(chan bool),
        protocolType: EPROTOCOL_TYPE_TCP,

    }

    if c.conn != nil {
        //创建go的线程 使用Goroutine
        go c.ConnSender()
        go c.ConnReader()
    }

    return c
}

func (c *TcpConnection) GetId() int {
    return c.id
}


func (c *TcpConnection) GetBuffer() *ConnectionBuffer {
    return c.ConnectionBuffer
}


/*
func (c *TcpConnection) GetStream() IStream {
    return c.stream
}
*/


func (c *TcpConnection) Connect(addr string) bool {
    gts.Trace("connect to tcpserver[info]:", addr)

    conn, err := net.Dial("tcp", addr)
    if err != nil {
        print(err)
        gts.Warn("net.Dial to %s:%q",addr, err)
        return false
    } else {
        gts.Trace("tcp dial to %s is ok. ", addr)
    }

    //go func() {
        //defer conn.Close()

        c.conn = conn

        //创建go的线程 使用Goroutine
        go c.ConnSender()
        go c.ConnReader()

        //<-c.Quit
    //}()
    return true
}


func (c *TcpConnection) ConnReader() {
    gts.Trace("read_buffer_size:", c.read_buffer_size)
    buffer := make([]byte, c.read_buffer_size)
    for {
        bytesRead, err := c.conn.Read(buffer)

        if err != nil {
            gts.Error("tcpconnection connreader error: ", err, bytesRead)
            if c.closeCallback != nil {
                c.closeCallback(c.id)
            }
            break
        }

        gts.Trace("tcpConnReader read to buff:%d, % X",bytesRead, buffer[:bytesRead])
        gts.Trace("tcpConnReader read to buff:%q",buffer[:bytesRead])
        c.Stream.Write(buffer[0:bytesRead])
        c.receiveCallback(c)

        //gts.Trace("tpool ConnReader Buff:%d", len(c.stream.Bytes()))

        //n, dps := c.routePack.Fetch(c)
        //gts.Trace("fetch message number", n)
        //if n > 0 {
            //c.receivePacketCallback(c.id, dps)
        //}
    }
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


func (c *TcpConnection) GetProtocolType() EProtocolType {
    return c.protocolType
}


func (c *TcpConnection) GetType() EConnType {
    return c.connectionType
}


func (c *TcpConnection) SetType(t EConnType) {
    c.connectionType = t
}


func (c *TcpConnection) Close() {
    c.quit <- true
}


/*
func (c *TcpConnection) Send(dp *RoutePacket) {
    c.outgoing <- dp
}
*/


func (c *TcpConnection) Send(data []byte) {
    c.outgoingBytes <- data
}


