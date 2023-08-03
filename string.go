package tools

import (
	"bytes"
	"fmt"
	"strings"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

func StringisGBK(data []byte) bool {
	// 尝试将 data 转换为 UTF-8 编码
	utf8Data, _, err := transform.Bytes(simplifiedchinese.GBK.NewDecoder(), data)
	if err != nil {
		// 转换失败，说明不是 GBK 编码
		return false
	}

	// 检查转换后的 UTF-8 数据是否与原始数据相同
	return bytes.Equal(data, utf8Data)
}
func StringSetGBK(b []byte) []byte {
	r, _ := simplifiedchinese.GBK.NewDecoder().Bytes(b)
	return r
}
func StringGetGBK(b []byte) []byte {
	r, _ := simplifiedchinese.GBK.NewEncoder().Bytes(b)
	return r
}
func ListHasSpace(list []string) int {
	for i, w := range list {
		if w[0] == '"' && list[i+1][len(list[i+1])-1] == '"' {
			list[i] = w + " " + list[i+1]
			return i + 1
		}
	}
	return 0
}
func ListDeleteSpace(list *[]string, del_seq int) bool {
	if del_seq == 0 {
		return false
	}
	j := 0
	for i, w := range *list {
		if i != del_seq {
			(*list)[j] = w
			j++
		}
	}
	*list = (*list)[:j]
	return true

}

// 通过递归删除切片中的重复元素
func ListDeleteString(list *[]string, del_str string) ([]string, int, bool) {
	retval := false
	j := 0
	for _, w := range *list {
		if w != del_str {
			(*list)[j] = w
			j++
		} else {
			retval = true
		}
	}
	*list = (*list)[:j]
	if retval {
		return ListDeleteString(list, del_str)
	}
	return *list, len(*list), false
}

// 合并类似"chi nese"这样的单词,并返回到切片
func ListMargenSpace(cmd string) []string {
	list := strings.Fields(cmd)
	for ListDeleteSpace(&list, ListHasSpace(list)) {
	}
	return list
}

// 判断 str 是否在 list 中
func ListHasString(list []string, str string) bool {
	if len(list) == 0 {
		return false
	}
	for _, v := range list {
		if v == str {
			return true
		}
	}
	return false
}

func ListOutput(list []string) {
	for _, line := range list {
		fmt.Println(line)
	}
}

func ListAddString(list *[]string, str string) {
	if !ListHasString(*list, str) {
		*list = append(*list, str)
	}
}
