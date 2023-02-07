package Tools

import (
	"bytes"
	"errors"
)

// 从Buffer读取指定数量的Byte,如果剩余位数不够，则返回Error
//
//	n:读取byte位数
func BufRead(buffer *bytes.Buffer, n int, errMsg string) ([]byte, error) {
	if buffer.Len() < n {
		return nil, errors.New(errMsg)
	}
	return buffer.Next(n), nil
}
