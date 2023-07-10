package main
 
import (
    "net"
    "log"
    "fmt"
)
 
func main() {
    port := "9090"
    Start(port)
}
 
// start the server
func Start(port string) {
    host := ":" + port
 
         //Get the tcp address
    tcpAddr, err := net.ResolveTCPAddr("tcp4", host)
    if err != nil {
        log.Printf("resolve tcp addr failed: %v\n", err)
        return
    }
 
         // listen
    listener, err := net.ListenTCP("tcp", tcpAddr)
    if err != nil {
        log.Printf("listen tcp port failed: %v\n", err)
        return
    }
 
         // Establish a connection pool for broadcast messages
    conns := make(map[string]net.Conn)
 
         // message channel
    messageChan := make(chan string, 10)
 
         // broadcast message
    go BroadMessages(&conns, messageChan)
 
         // start up 
    for {
        fmt.Printf("listening port %s ...\n", port)
        conn, err := listener.AcceptTCP()
        if err != nil {
            log.Printf("Accept failed:%v\n", err)
            continue
        }
 
                 // Throw each client connection into the connection pool
        conns[conn.RemoteAddr().String()] = conn
        fmt.Println(conns)
 
                 // process the message
        go Handler(conn, &conns, messageChan)
    }
}
 
 // Broadcast to all the connected folks
func BroadMessages(conns *map[string]net.Conn, messages chan string) {
    for {
 
                 // constantly read messages from the channel
        msg := <-messages
        fmt.Println(msg)
 
                 // Send a message to all the folks
        for key, conn := range *conns {
            fmt.Println("connection is connected from ", key)
            _, err := conn.Write([]byte(msg))
            if err != nil {
                log.Printf("broad message to %s failed: %v\n", key, err)
                delete(*conns, key)
            }
        }
    }
}
 
 // Handle the message sent by the client to the server and throw it into the channel
func Handler(conn net.Conn, conns *map[string]net.Conn, messages chan string) {
    fmt.Println("connect from client ", conn.RemoteAddr().String())
 
    buf := make([]byte, 1024)
    for {
        length, err := conn.Read(buf)
        if err != nil {
            log.Printf("read client message failed:%v\n", err)
            delete(*conns, conn.RemoteAddr().String())
            conn.Close()
            break
        }
 
                 // Write the received message to the channel
        recvStr := string(buf[0:length])
        messages <- recvStr
    }
}