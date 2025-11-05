package codec

var alphabet = []byte("0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz")

// B62Encode 字节切片（[]byte）编码为Base62字符串
func B62Encode(raw []byte) string {
	// 简化实现：将 raw 视作大整数做除法编码
	// 为避免大整数库，这里做逐字节“进制转换”
	buf := make([]byte, 0, len(raw)*2)
	val := append([]byte(nil), raw...)
	for len(val) > 0 && !(len(val) == 1 && val[0] == 0) {
		carry := 0
		next := make([]byte, 0, len(val))
		for _, v := range val {
			acc := int(v) + carry*256
			q := acc / 62
			r := acc % 62
			if len(next) > 0 || q > 0 {
				next = append(next, byte(q))
			}
			buf = append(buf, alphabet[r])
		}
		val = next
	}
	// 反转
	for i, j := 0, len(buf)-1; i < j; i, j = i+1, j-1 {
		buf[i], buf[j] = buf[j], buf[i]
	}
	return string(buf)
}

// B62Decode 将Base62字符串解码为字节切片
func B62Decode(encoded string) []byte {
	if encoded == "" {
		return []byte{}
	}
	// 创建字符到索引的映射
	charToIndex := make(map[byte]int)
	for i, char := range alphabet {
		charToIndex[char] = i
	}
	// 将输入字符串转换为字节切片
	input := []byte(encoded)
	// 初始化结果数组
	result := make([]byte, 0)
	// 逐个处理输入字符
	for _, char := range input {
		// 查找字符在alphabet中的索引
		index, exists := charToIndex[char]
		if !exists {
			// 非法字符，返回空切片或处理错误
			return []byte{}
		}
		// 将result视为一个大整数，乘以62并加上当前索引值
		carry := index
		for i := range result {
			acc := int(result[i])*62 + carry
			result[i] = byte(acc % 256)
			carry = acc / 256
		}
		// 处理剩余的进位
		for carry > 0 {
			result = append(result, byte(carry%256))
			carry /= 256
		}
	}
	// 反转结果（因为计算是从低位开始的）
	for i, j := 0, len(result)-1; i < j; i, j = i+1, j-1 {
		result[i], result[j] = result[j], result[i]
	}
	return result
}
