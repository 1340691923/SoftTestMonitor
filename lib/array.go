package lib

//判断arr是否包含 s
func InArrayStr(s string,arr []string)bool{
	for _,tmp := range arr{
		if tmp == s {
			return true
		}
	}
	return false
}
