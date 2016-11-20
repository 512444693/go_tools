package zmmd5

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"sort"
)

type md5Info struct {
	name  string
	value string
}

type md5Infos []md5Info

func (infos md5Infos) Len() int {
	return len(infos)
}

func (infos md5Infos) Less(i, j int) bool {
	return infos[i].name < infos[j].name
}

func (infos md5Infos) Swap(i, j int) {
	infos[i], infos[j] = infos[j], infos[i]
}

func (infos md5Infos) String() string {
	var buffer bytes.Buffer
	for i := range infos {
		buffer.WriteString(fmt.Sprintf("%s : %s\r\n", infos[i].name, infos[i].value))
	}
	return string(buffer.Bytes())
}

func newMd5Info(filePath string) (ret md5Info, ok bool) {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("打开文件%s失败", filePath)
		return ret, false
	}
	defer file.Close()

	md5h := md5.New()
	io.Copy(md5h, file)
	md5Str := fmt.Sprintf("%x", md5h.Sum(nil))

	ret = md5Info{file.Name(), md5Str}
	return ret, true
}

//CalMd5 计算一个文件或文件夹的md5
//dirPath为文件夹路径
//按照顺序输出
func CalMd5(dirPath *string) {
	infos := make(md5Infos, 0, 20)

	dir, err := ioutil.ReadDir(*dirPath)
	if err != nil {
		fmt.Printf("读取目录失败:%s", err.Error())
		return
	}

	for i := range dir {
		if !dir[i].IsDir() {
			tmp := dir[i].Name()
			if info, ok := newMd5Info(*dirPath + string(os.PathSeparator) + tmp); ok {
				infos = append(infos, info)
			} else {
				return
			}

		}
	}

	sort.Sort(infos)
	fmt.Println(infos)
}
