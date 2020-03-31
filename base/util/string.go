package util

func StringAddByte(s string, b []byte) (rt []byte) {
	rt = []byte(s)
	rt = append(rt, b...)
	return
}
