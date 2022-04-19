package helper

// Base64OrigLength 计算 base64 数据长度
// https://stackoverflow.com/questions/56140620/how-to-get-original-file-size-from-base64-encode-string
func Base64OrigLength(datas string) int64 {

	l := int64(len(datas))

	// count how many trailing '=' there are (if any)
	eq := int64(0)
	if l >= 2 {
		if datas[l-1] == '=' {
			eq++
		}
		if datas[l-2] == '=' {
			eq++
		}

		l -= eq
	}

	// basically:
	// eq == 0 :	bits-wasted = 0
	// eq == 1 :	bits-wasted = 2
	// eq == 2 :	bits-wasted = 4

	// so orig length ==  (l*6 - eq*2) / 8

	return (l*3 - eq) / 4
}
