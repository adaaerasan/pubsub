package main

const (
	HASH_SEED = 131
)
func getHash(src string)uint32{
	var des  = uint32(0)
	for i := 0;i<len(src);i++{
		des = des *131 + uint32(src[i])
	}
	return des
}