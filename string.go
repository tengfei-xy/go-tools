package tools

import (
	"fmt"
	"strings"

	"golang.org/x/text/encoding/simplifiedchinese"
)

func toChineseChar(b []byte) []byte {
	r, _ := simplifiedchinese.GBK.NewDecoder().Bytes(b)
	return r
}
func getChineseChar(b []byte) []byte {
	r, _ := simplifiedchinese.GBK.NewEncoder().Bytes(b)
	return r
}
func list_has_space(list []string) int {
	for i, w := range list {
		if w[0] == '"' && list[i+1][len(list[i+1])-1] == '"' {
			list[i] = w + " " + list[i+1]
			return i + 1
		}
	}
	return 0
}
func list_delete_space(list *[]string, del_seq int) bool {
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
func list_delete_string(list *[]string, del_str string) ([]string, int, bool) {
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
		return list_delete_string(list, del_str)
	}
	return *list, len(*list), false
}

// 合并类似"chi nese"这样的单词,并返回到切片
func list_margen_space(cmd string) []string {
	list := strings.Fields(cmd)
	for list_delete_space(&list, list_has_space(list)) {
	}
	return list
}

// 判断 str 是否在 list 中
func list_has_string(list []string, str string) bool {
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

func list_output(list []string) {
	for _, line := range list {
		fmt.Println(line)
	}
}

func list_add_string(list *[]string, str string) {
	if !list_has_string(*list, str) {
		*list = append(*list, str)
	}
}
