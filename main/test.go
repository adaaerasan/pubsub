package main

import "fmt"
import (
	"net"
	"time"
)

func main(){
	var sv = Server{}
	go sv.StartServer("localhost:10000")
	go pub()
	//go sub()
	var ch = make(chan int)
	<-ch
}
type subSock struct{

}
func(ss *subSock)read(sock *Socket,ba []byte,err error){
	if err != nil{
		fmt.Println(err.Error())
	}
	var mi protoMsgInfo
	mi ,err = unpack(ba)
	if err != nil{
		return
	}
	fmt.Println(mi.msgType)
}
var chanelTest = "channelTest"
func sub(){
	for{
		time.Sleep(3*time.Second)
		conn, err := net.Dial("tcp", "127.0.0.1:10000")
		if err != nil{
			fmt.Println(err.Error())
			continue
		}
		defer conn.Close()
		var mi protoMsgInfo
		mi.body = append(mi.body,[]byte(chanelTest))
		mi.msgType = MSG_TYPE_SUB_APPEND
		ba,_ := pack(mi)
		conn.Write(ba)
		var sock = Socket{conn}
		var r = subSock{}
		go sock.StartRead(&r,time.Duration(2)*time.Second)
		var ch = make(chan int)
		<-ch
		break
	}
}
func pub(){
	for{
		time.Sleep(3*time.Second)
		conn, err := net.Dial("tcp", "127.0.0.1:10000")
		if err != nil{
			fmt.Println(err.Error())
			continue
		}
		defer conn.Close()
		fmt.Println("dial")
		for {
			//time.Sleep(5000*1000)
			time.Sleep(3*time.Second)
			var mi = protoMsgInfo{}
			mi.body = append(mi.body,[]byte(chanelTest))
			mi.body = append(mi.body,[]byte("test"))
			mi.msgType = MSG_TYPE_PUB
			ba ,err := pack(mi)
			if err != nil{
				fmt.Println("pack err",err.Error())
			}
			fmt.Println("pub msg")
			if _,err = conn.Write(ba);err != nil{
				fmt.Println(err.Error())
				break
			}
		}
	}
}