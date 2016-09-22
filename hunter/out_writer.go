package hunter

import (
	"os"
	"sort"
	"strconv"
)

var outs = make(map[string][]string)

// 启动字符统计分类
func writerStart() {
	// 例如: qq.com 就是一种字符的域名， baidu.com就是5种字符的域名，google.com就是4种字符的域名
	for i := 1; i <= Conf.AdvFilter.OccurChars; i++ {
		s := strconv.Itoa(i)
		outs[s] = make([]string, 0)
	}
}

func write(out *os.File) {
	for i := 1; i <= Conf.AdvFilter.OccurChars; i++ {
		s := strconv.Itoa(i)
		o := outs[s]

		// 在这里进行自定义规则过滤
		sort.Strings(o)

		out.WriteString("-----------------------------------" + s + "种字符域名列表------------------------------------------------------\n")
		for _, v := range o {
			out.WriteString(v + "\n")
		}
	}

	out.Sync()
}
