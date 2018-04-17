package main

import (
	"sort"
	"sync"
)

type serverInfo struct{
	hash uint32
	s *Socket
}
type serverList []serverInfo
func (s *serverList)Len()int{
	return s.Len()
}
func (s *serverList)Less(i,j int)bool{
	if (*s)[j].hash > (*s)[i].hash{
		return true
	}else{
		return false
	}
}
func (s *serverList)Swap(i,j int){
	var tmp = (*s)[i]
	(*s)[i] = (*s)[j]
	(*s)[j] = tmp
}

type saveServer struct{
	sl serverList
	lock sync.RWMutex
}

func(s *saveServer)AddServer(sock *Socket){
	s.lock.Lock()
	var addr = sock.Conn.RemoteAddr().String()
	var si serverInfo
	si.hash = getHash(addr)
	si.s = sock
	s.sl = append(s.sl,si)
	sort.Sort(&s.sl)
	s.lock.Unlock()
}
func(s *saveServer)GetServer(str string)*Socket{
	s.lock.RLocker()
	defer s.lock.RUnlock()
	var hash = getHash(str)
	for i := 0; i< s.sl.Len();i++{
		if s.sl[i].hash > hash{
			return s.sl[i].s
		}
	}
	if s.sl.Len() > 0{
		return s.sl[0].s
	}
	return nil
}
func(s *saveServer)DelServer(sock *Socket){
	s.lock.Lock()
	var index = -1
	var hash = getHash(sock.Conn.RemoteAddr().String())
	for i := 0;i< s.sl.Len();i++{
		if s.sl[i].hash == hash{
			index = i
			break
		}
	}
	if index != -1{
		s.sl = append(s.sl[:index],s.sl[index+1:]...)
	}
	s.lock.Unlock()
}
