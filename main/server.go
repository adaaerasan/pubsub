package main

import (
	"net"
	"time"
	"errors"
	"fmt"
)


type Server struct{
	ChannelMap map[string]*SubSocket
	channelMap map[string]map[*Socket]bool
}
func (s *Server)read(sock *Socket,ba []byte,err error){
	if err != nil{
		fmt.Println(err.Error())
		return
	}
	var mi protoMsgInfo
	mi ,err = unpack(ba)
	if err != nil{
		fmt.Println(err.Error())
		return
	}
	fmt.Println(string(mi.body[0]))
	err = s.parseMsg(sock,mi)
	if err != nil{
		fmt.Println(err.Error())
		return
	}
}
func (s * Server)processMsgPub(sock *Socket,mi protoMsgInfo)error{
	if mi.blockLen	 == 2{
		var c = string(mi.body[0])
		ss,find := s.channelMap[c]
		if find{
			for k,_ := range ss	{
				var mi protoMsgInfo
				mi.body = append(mi.body,mi.body[0])
				mi.msgType = MSG_TYPE_TEXT
				k.Send(mi.body[0],time.Duration(2))
			}
		}
	}else{
		return errors.New("pub parameter")
	}
	return nil
}
func (s *Server)processMsgSub(sock *Socket,mi protoMsgInfo)error{
	switch mi.msgType{
	case MSG_TYPE_SUB_APPEND:
		if (mi.blockLen == 1){
			var c = string(mi.body[0])
			ss,find := s.channelMap[c]
			if find{
				ss[sock] = true
			}else{
				ss = make(map[*Socket]bool)
				ss[sock] = true
			}
			s.channelMap[c] = ss
		}else{
			return errors.New("parameters invalid")
		}
	case MSG_TYPE_SUB_REPLACE:
		if (mi.blockLen == 2){
			var oldC = string(mi.body[0])
			var newC = string(mi.body[1])
			ss,find := s.channelMap[oldC]
			if find{
				delete(ss,sock)
			}
			ss,find = s.channelMap[newC]
			if find{
				ss[sock] = true
			}else{
				ss = make(map[*Socket]bool)
				ss[sock] = true
			}
			s.channelMap[newC] = ss
		}
		return errors.New("not support")
	case MSG_TYPE_SUB_DEL:
		if (mi.blockLen == 1){
			var oldC = string(mi.body[0])
			ss,find := s.channelMap[oldC]
			if find{
				delete(ss,sock)
			}
		}else{
			return errors.New("parameters invalid")
		}
	case MSG_TYPE_SUB_CLEAR:
		for _ ,v := range s.channelMap{
			delete(v, sock)
		}
	default:
		return errors.New("not support type")
	}
	return nil
}
func (s *Server)parseMsg(sock *Socket,mi protoMsgInfo)error{
	if mi.msgType & MSG_TYPE_PUB == MSG_TYPE_PUB{
		return s.processMsgPub(sock,mi)
	}else if mi.msgType & MSG_TYPE_SUB == MSG_TYPE_SUB{
		return s.processMsgSub(sock,mi)
	}else{
		return errors.New("not support msgtype")
	}
}
func (s *Server)StartServer(hostPort string)error{
	listen_sock, err := net.Listen("tcp", hostPort)
	if err != nil{
		return err
	}
	defer listen_sock.Close()
	for {
		new_conn, err := listen_sock.Accept()
		if err != nil {
			continue
		}
		var sock = Socket{new_conn}
		go sock.StartRead(s,time.Duration(15*1000))
	}
	return nil
}
