package main

import (
	"bytes"
	"encoding/binary"
	"errors"
)
const (

	MSG_TYPE_APPEND = 1 << 4
	MSG_TYPE_DEL = 1 << 5
	MSG_TYPE_REPLACE = 1 << 6

	MSG_TYPE_PUB = 1
	MSG_TYPE_SUB = 3

	MSG_TYPE_SUB_APPEND = MSG_TYPE_SUB & MSG_TYPE_APPEND
	MSG_TYPE_SUB_DEL = MSG_TYPE_SUB & MSG_TYPE_DEL
	MSG_TYPE_SUB_REPLACE = MSG_TYPE_SUB & MSG_TYPE_REPLACE
)
type msgInfo struct {
	blockLen uint8
	msgType uint8
	msgLength []uint32
	body [][]byte
}
func pack(mi msgInfo)([]byte,error){
	var bytesBuf = bytes.NewBuffer([]byte{})
	if len(mi.msgLength) != len(mi.body){
		return []byte(""),errors.New("not match")
	}
	mi.blockLen = uint8(len(mi.body))
	err := binary.Write(bytesBuf,binary.LittleEndian,mi.blockLen)
	if err != nil{
		return []byte(""),err
	}
	err = binary.Write(bytesBuf,binary.LittleEndian,mi.msgType)
	if err != nil{
		return []byte(""),err
	}
	mi.msgLength = make([]uint32,mi.blockLen)
	for i := uint8(0);i < mi.blockLen;i++{
		mi.msgLength[i] = uint32(len(mi.body[i]))
	}
	err = binary.Write(bytesBuf,binary.LittleEndian,mi.msgLength)
	if err != nil{
		return []byte(""),err
	}

	for i := uint8(0);i < mi.blockLen;i++{
		err = binary.Write(bytesBuf,binary.LittleEndian,mi.body[i])
		if err != nil{
			return []byte(""),err
		}
	}
	return bytesBuf.Bytes(),nil
}
func unpack(ba []byte)(msgInfo,error){
	var length uint8 = 0
	var bytesBuf = bytes.NewBuffer(ba)
	binary.Read(bytesBuf,binary.LittleEndian,&length)
	var mi msgInfo = msgInfo{0,0,make([]uint32,length),make([][]byte,length)}

	binary.Read(bytesBuf,binary.LittleEndian,&mi.msgType)
	binary.Read(bytesBuf,binary.LittleEndian,&mi.msgLength)
	for i := uint8(0); i < length; i++{
		var l = mi.msgLength[i]
		mi.body[i] = make([]byte,l)
		binary.Read(bytesBuf,binary.LittleEndian,&mi.body[i])
	}
	return mi,nil
}