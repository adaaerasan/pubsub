package main

import (
	"net"
	"time"
	"errors"
	"fmt"
	"sync"
)


type Server struct{
	channelMap map[string]map[*Socket]bool
	lock sync.RWMutex
	sl saveServer
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
	fmt.Println(string(mi.body[0]),mi.msgType)
	err = s.parseMsg(sock,mi)
	if err != nil{
		fmt.Println(err.Error())
		return
	}
}
func (s * Server)processMsgPub(sock *Socket,mi protoMsgInfo)error{
	if mi.blockLen	 == 2{
		s.lock.RLock()
		var c = string(mi.body[0])
		ss,find := s.channelMap[c]
		if find{
			for k,_ := range ss	{
				var tmi protoMsgInfo
				tmi.body = append(tmi.body,mi.body[1])
				fmt.Println("send msg------",string(mi.body[1]))
				tmi.msgType = MSG_TYPE_TEXT
				ba,_ := pack(tmi)
				k.Send(ba,time.Duration(2)*time.Second)
			}
		}
		s.lock.RUnlock()
		saveSock := s.sl.GetServer(c)
		if saveSock != nil{
			var tmi protoMsgInfo
			tmi.body = append(tmi.body,mi.body[0])
			tmi.body = append(tmi.body,mi.body[1])
			tmi.msgType = MSG_TYPE_SAVE
			ba,_ := pack(tmi)
			saveSock.Send(ba,time.Duration(2)*time.Second)
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
			s.lock.Lock()
			var c = string(mi.body[0])
			ss,find := s.channelMap[c]
			if find{
				ss[sock] = true
			}else{
				fmt.Println("append sub")
				ss = make(map[*Socket]bool)
				ss[sock] = true
			}
			s.channelMap[c] = ss
			s.lock.Unlock()
		}else{
			fmt.Println(mi.blockLen)
			return errors.New("parameters invalid")

		}
	case MSG_TYPE_SUB_REPLACE:
		if (mi.blockLen == 2){
			s.lock.Lock()
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
			s.lock.Unlock()
		}else{
			return errors.New("not support")
		}
	case MSG_TYPE_SUB_DEL:
		if (mi.blockLen == 1){
			s.lock.Lock()
			var oldC = string(mi.body[0])
			ss,find := s.channelMap[oldC]
			if find{
				delete(ss,sock)
			}
			s.lock.Unlock()
		}else{
			return errors.New("parameters invalid")
		}
	case MSG_TYPE_SUB_CLEAR:
		s.lock.Lock()
		for _ ,v := range s.channelMap{
			delete(v, sock)
		}
		s.lock.Unlock()
	default:
		return errors.New("not support type")
	}
	return nil
}

func (s *Server)processSaveMsg(sock *Socket,mi protoMsgInfo)error{
	var remodAddr = sock.Conn.RemoteAddr().String()
	var h = getHash(remodAddr)
	s.slock.Lock()
	s.saveMap[h] = sock
	s.slock.Unlock()
	return nil
}
func (s *Server)parseMsg(sock *Socket,mi protoMsgInfo)error{
	if mi.msgType & MSG_TYPE_PUB == MSG_TYPE_PUB{
		return s.processMsgPub(sock,mi)
	}else if mi.msgType & MSG_TYPE_SUB == MSG_TYPE_SUB{
		return s.processMsgSub(sock,mi)
	}else if mi.msgType == MSG_TYPE_SAVE{
		return s.processSaveMsg(sock,mi)
	}else{
		return errors.New("not support msgtype")
	}
}
func (s *Server)StartServer(hostPort string)error{
	s.channelMap = make(map[string]map[*Socket]bool)
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
