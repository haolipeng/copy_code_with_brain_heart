package consistenthash

import (
	"strconv"
	"testing"
)

//测试思路
//使用自定义的hash函数，其只处理数字，传入字符串格式的数字，返回对应的数字即可
//一开始，有2/4/6三个真实节点，对应的虚拟节点的哈希值是02/12/22，04/14/24，06/16/26
//那么测试用例中2/11/23/27选择的虚拟节点分别是02/12/24/02,换算成真实节点是2/2/4/2
//添加一个真实节点8，对应的虚拟节点的哈希值是08/18/28，此时，用例27对应的虚拟节点从02变更为28，即真实节点8

func TestHashing(t *testing.T) {
	hash := New(3, func(key []byte) uint32 {
		i, _ := strconv.Atoi(string(key))
		return uint32(i)
	})
	hash.Add("6", "4", "2")
	testCases := map[string]string{
		"2":  "2",
		"11": "2",
		"23": "4",
		"27": "2",
	}

	for k, v := range testCases {
		if hash.Get(k) != v {
			t.Errorf("Asking for %s, should have yielded %s", k, v)
		}
	}

	// Adds 8, 18, 28
	hash.Add("8")

	// 27 should now map to 8.
	testCases["27"] = "8"

	for k, v := range testCases {
		if hash.Get(k) != v {
			t.Errorf("Asking for %s, should have yielded %s", k, v)
		}
	}
}
