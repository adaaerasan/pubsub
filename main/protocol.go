package main

import (
	"bytes"
	"encoding/binary"
	"errors"
)
const (

	MSG_TYPE_APPEND = uint8(1 << 4)
	MSG_TYPE_DEL = uint8(1 << 5)
	MSG_TYPE_REPLACE = uint8(1 << 6)
	MSG_TYPE_CLEAR = uint8(1 << 7)

	MSG_TYPE_PUB = uint8(1)
	MSG_TYPE_SUB = uint8(3)
	MSG_TYPE_TEXT = uint8(2)

	MSG_TYPE_SUB_APPEND = 0x13
	MSG_TYPE_SUB_DEL = 0x23
	MSG_TYPE_SUB_REPLACE = 0x43
	MSG_TYPE_SUB_CLEAR = 0x83

	PACKAGE_LEN_MAX = uint64(1) << 32
	SINGLE_BODY_LEN_MAX = uint32(1) << 20
)
type protoMsgInfo struct {
	packageLen uint32
	blockLen uint8
	msgType uint8
	msgLength []uint32
	body [][]byte
}
func getPackageLen(ba []byte)(uint32,error){
	var l = uint32(0)
	var bytesBuf = bytes.NewBuffer(ba)
	err := binary.Read(bytesBuf,binary.LittleEndian,&l)
	return l,err
}
func pack(mi protoMsgInfo)([]byte,error){
	var bytesBuf = bytes.NewBuffer([]byte{})

	mi.blockLen = uint8(len(mi.body))
	mi.msgLength = make([]uint32,len(mi.body))

	mi.msgLength = make([]uint32,mi.blockLen)
	var bytesLen = uint32(0)
	for i := uint8(0);i < mi.blockLen;i++{
		var l = uint32(len(mi.body[i]))
		mi.msgLength[i] = l
		if l > SINGLE_BODY_LEN_MAX{
			return []byte(""),errors.New("block too long")
		}
		bytesLen += l
		if uint64(bytesLen) > PACKAGE_LEN_MAX{
			return []byte(""),errors.New("package too long")
		}
	}
	mi.packageLen = uint32(4 + 1 + 1 + 4 * mi.blockLen ) + bytesLen
	err := binary.Write(bytesBuf,binary.LittleEndian,mi.packageLen)
	err = binary.Write(bytesBuf,binary.LittleEndian,mi.blockLen)
	if err != nil{
		return []byte(""),err
	}
	err = binary.Write(bytesBuf,binary.LittleEndian,mi.msgType)
	if err != nil{
		return []byte(""),err
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
func unpack(ba []byte)(protoMsgInfo,error){
	var length uint8 = 0
	var bytesBuf = bytes.NewBuffer(ba)
	err := binary.Read(bytesBuf,binary.LittleEndian,&length)
	if err != nil{
		return protoMsgInfo{},err
	}
	var mi protoMsgInfo = protoMsgInfo{0,0,0,make([]uint32,length),make([][]byte,length)}

	err = binary.Read(bytesBuf,binary.LittleEndian,&mi.msgType)
	if err != nil{
		return protoMsgInfo{},err
	}
	err = binary.Read(bytesBuf,binary.LittleEndian,&mi.msgLength)
	if err != nil{
		return protoMsgInfo{},err
	}
	for i := uint8(0); i < length; i++{
		var l = mi.msgLength[i]
		mi.body[i] = make([]byte,l)
		err = binary.Read(bytesBuf,binary.LittleEndian,&mi.body[i])
		if err != nil{
			return protoMsgInfo{},err
		}
	}
	return mi,nil
}