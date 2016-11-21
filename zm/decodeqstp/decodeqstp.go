package decodeqstp

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"strings"
)

//O 原始数据
var O = []byte{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9', 'A', 'D', 'G', 'J', 'M', 'Q', 'T', 'W', 'X', 'Z', 'a', 'c', 'f', 'o', 'h', 'n', 't', 'w', 'm', 'k', 'l', 'x', 'z', 'q', 'y', 'j', '.', '/'}

//E 加密后数据
var E = []byte{'A', 'D', 'G', 'J', 'M', 'Q', 'T', 'W', 'X', 'Z', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', 'm', 'k', 'l', 'x', 'z', 'q', 'y', 'j', 'a', 'c', 'f', 'o', 'h', 'n', 't', 'w', '/', '.'}

//Decode 解码QSTP串并输出到控制台
//path 为文件路径
func Decode(filePath *string) {
	if len(O) != len(E) {
		fmt.Println("对应字符个数不一致")
		return
	}
	DecryptMapping := make([]byte, 128)
	for i := 0; i < len(O); i++ {
		DecryptMapping[E[i]] = O[i]
	}
	var buffer bytes.Buffer
	if inputBytes, err := ioutil.ReadFile(*filePath); err == nil {
		inputBytes = []byte(strings.TrimPrefix(string(inputBytes), "qstp://"))
		for _, char := range inputBytes {
			var tmp byte
			if DecryptMapping[char] == 0 {
				tmp = char
			} else {
				tmp = DecryptMapping[char]
			}
			buffer.WriteByte(tmp)
		}
		fmt.Println("解密后数据：")
		fmt.Printf("%s", buffer.Bytes())
	} else {
		fmt.Printf("打开文件%s失败", *filePath)
	}
}
