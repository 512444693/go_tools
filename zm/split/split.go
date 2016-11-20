package split

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
)

//MAXSIZE 每次读取文件的最大字节大小
const MAXSIZE uint64 = 4 * 1024 * 1024

//DIRECTORY 将生成的文件放在这个目录下
const DIRECTORY string = "AFTER" + string(os.PathSeparator)

type fileInfo struct {
	name string
	size uint64
}

func newFileInfo(infoStr string) (ret *fileInfo, ok bool) {
	strs := strings.Split(infoStr, " ")
	if len(strs) != 2 {
		return nil, false
	}
	size, error := strconv.ParseUint(strs[1], 10, 64)
	if error != nil {
		fmt.Printf("转换数字失败：%s", strs[1])
		return nil, false
	}
	ret = &fileInfo{strs[0], size}
	return ret, true
}

func writeFile(reader *bufio.Reader, info *fileInfo) (ok bool) {
	file, err := os.OpenFile(DIRECTORY+info.name, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Printf("打开或创建文件失败：%s", err.Error())
		return false
	}
	defer file.Close()
	writer := bufio.NewWriterSize(file, int(MAXSIZE))

	readSize := info.size
	for {
		var tmp []byte
		var err error
		if readSize >= MAXSIZE {
			tmp, err = reader.Peek(int(MAXSIZE))
			reader.Discard(int(MAXSIZE))
			readSize -= MAXSIZE
		} else {
			tmp, err = reader.Peek(int(readSize))
			reader.Discard(int(readSize))
			readSize = 0
		}
		if err != nil {
			fmt.Printf("读取文件错误：%s", err.Error())
			return false
		}
		if _, err = writer.Write(tmp); err != nil {
			fmt.Printf("写文件失败：%s", err.Error())
			return false
		}
		writer.Flush()
		if readSize == 0 {
			break
		}
	}

	return true
}

//Split 将暴风包分解
func Split(bfpPath *string) {
	os.Mkdir(DIRECTORY, 0666)

	file, err := os.Open(*bfpPath)
	if err != nil {
		fmt.Printf("打开文件%s失败", *bfpPath)
		return
	}
	defer file.Close()

	//将文件头信息读取到 fileInfos 中
	fileInfos := make([]fileInfo, 0, 10)
	reader := bufio.NewReaderSize(file, int(MAXSIZE))
	for {
		inputString, err := reader.ReadString('\n')
		if err == io.EOF || strings.TrimSpace(inputString) == "" {
			break
		}
		trimStr := strings.TrimSpace(inputString)
		if ok, _ := regexp.MatchString("^\\S+ [0-9]+$", trimStr); ok {
			if tmpFileInfo, ok := newFileInfo(trimStr); ok {
				fileInfos = append(fileInfos, *tmpFileInfo)
			} else {
				return
			}
		}
	}

	if len(fileInfos) == 0 {
		fmt.Println("没有读取到头文件有文件信息")
		return
	}

	fmt.Printf("写如下文件 : %v\n", fileInfos)
	//根据 fileInfos 便利读取文件
	for i := range fileInfos {
		if ok := writeFile(reader, &fileInfos[i]); !ok {
			return
		}
	}

	var tmp byte
	if tmp, err = reader.ReadByte(); err != io.EOF {
		fmt.Printf("文件没有读完!!!%c", tmp)
		return
	}
}
