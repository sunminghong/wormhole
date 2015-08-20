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
func NewUdpConnection(newcid int, conn *net.UDPConn, endianer gts.IEndianer, userAddr *net.UDPAddr) *UdpConnection {
    c := &UdpConnection {
        id:      newcid,
        conn:     conn,
        userAddr: userAddr,

        outgoing: make(chan *RoutePacket, 1),
        outgoingBytes: make(chan []byte),
        receiveChan: make(chan bool, 20),
        quit:     make(chan bool),
        Quit:     make(chan bool),

        stream:   gts.NewRWStream(1024, endianer),
    }

    c.stream.Reset()

    //创建go的线程 使用Goroutine
    go c.reader()

    if conn != nil {
        go c.ConnSenderServer()
    }

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

    gts.Info("pool dial to %s is ok. ", addr)

    go func() {
        defer conn.Close()

        c.conn = conn

        buffer := make([]byte, 1024)
        go func() {
            for {
                n, err := c.conn.Read(buffer[0:])
                if err == nil {
                    c.ConnReader(buffer[0:n])
                } else {
                    e, ok := err.(net.Error)
                    if !ok || !e.Timeout() {
                        gts.Trace("recv error", err.Error(), udpAddr)
                        c.Quit <- true
                        return
                    }
                }
            }
        }()

        go c.ConnSenderClient()

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


func (c *UdpConnection) ConnSenderClient() {
    for {
        select {
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


func (c *UdpConnection) ConnSenderServer() {
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
}


/*
func (c *UdpConnection) Send(dp *RoutePacket) {
    c.outgoing <- dp
}
*/


func (c *UdpConnection) Send(data []byte) {
    c.outgoingBytes <- data
}


