package main

type PubSocket struct{

}

func(p *PubSocket)StartConnect(host string,port int){

}

func(p *PubSocket)SendMsg(ba []byte,channel string)error{
	var e error
	return e
}
func(p *PubSocket)Close(){

}
