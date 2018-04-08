package pubsub

type SubRecvMsg interface{
	RecvMsg()
}
type SubSocket struct{

}
func (s *SubSocket)StartSub(host string,port int,recvMsgOb SubRecvMsg)error{
	var err error
	return err
}
func (s *SubSocket)Close(){

}