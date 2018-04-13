package main

import (
	"net"
	"time"
	"fmt"
)
type IRead interface {
	read(*Socket,[]byte,error)
}
type msgInfo struct{
	length int
	msgtype int
	msg string
}
type Socket struct{
	Conn net.Conn
}
func (s *Socket)readBuf(buf *[]byte,n uint32)error{
	tn,err := s.Conn.Read(*buf)
	fmt.Println("read bytes",tn,n)
	return err
}
func (s *Socket)StartRead(readObject IRead ,duration time.Duration){
	fmt.Println("start read")
	defer s.Conn.Close()
	for{
		headBuf := make([]byte, 4)
		err := s.readBuf(&headBuf,4)
		fmt.Println("readheader")
		if err != nil{
			readObject.read(s,headBuf,err)
			return
		}
		var l uint32
		l,err = getPackageLen(headBuf)
		fmt.Println("recv len:",l)
		if err != nil{
			readObject.read(s,[]byte(""),err)
		}
		var bodyBuf = make([]byte,l - 4)
		err = s.readBuf(&bodyBuf,l - 4)
		if err != nil{
			readObject.read(s,[]byte(""),err)
		}else{
			readObject.read(s,bodyBuf,nil)
		}
	}
}
func (s *Socket)Send(ba []byte,duration time.Duration)error{
	var l = len(ba)
	var writeLen = 0
	for{
		n,err := s.Conn.Write(ba)
		if err != nil{
			return err
		}
		writeLen += n
		if writeLen == l{
			break
		}
	}
	return nil
}


