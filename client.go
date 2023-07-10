package main
 
import (
    "net"
    "log"
    "fmt"
    "os"
)
 
func main() {
    Start(os.Args[1])
}
 
func Start(tcpAddrStr string) {
    tcpAddr, err := net.ResolveTCPAddr("tcp4", tcpAddrStr)
    if err != nil {
        log.Printf("Resolve tcp addr failed: %v\n", err)
        return
    }
 
         // Dial the server
    conn, err := net.DialTCP("tcp", nil, tcpAddr)
    if err != nil {
        log.Printf("Dial to server failed: %v\n", err)
        return
    }
 
         // send a message to the server
    go SendMsg(conn)
 
         // Receive broadcast messages from the server side
    buf := make([]byte, 1024)
    for {
        length, err := conn.Read(buf)
        if err != nil {
            log.Printf("recv server msg failed: %v\n", err)
            conn.Close()
            os.Exit(0)
            break
        }
 
        fmt.Println(string(buf[0:length]))
    }
}
 
 // send a message to the server
func SendMsg(conn net.Conn) {
    username := conn.LocalAddr().String()
    for {
        var input string
 
                 // Receive the input message and put it in the input variable
        fmt.Scanln(&input)
 
        if input == "/q" || input == "/quit" {
            fmt.Println("Byebye ...")
            conn.Close()
            os.Exit(0)
        }
 
                 // only handle messages with content
        if len(input) > 0 {
            msg := username + " say:" + input
            _, err := conn.Write([]byte(msg))
            if err != nil {
                conn.Close()
                break
            }
        }
    }
}