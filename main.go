package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

var types map[string]string
var typeList string

func init() {
	types = make(map[string]string)
	types["feat"] = "增加新功能"
	types["fix"] = "修复问题/BUG"
	types["style"] = "代码风格相关无影响运行结果的"
	types["perf"] = "优化/性能提升"
	types["refactor"] = "重构"
	types["revert"] = "撤销修改"
	types["test"] = "测试相关"
	types["docs"] = "文档/注释"
	types["chore"] = "依赖更新/脚手架配置修改等"
	types["workflow"] = "工作流改进"
	types["ci"] = "持续集成"
	types["mod"] = "不确定分类的修改"
	types["wip"] = "开发中"
	types["types"] = "类型修改"

	sl := make([]string, 0)
	for k, v := range types {
		sl = append(sl, fmt.Sprintf("%8s: %s", k, v))
	}

	typeList = strings.Join(sl, "\n")
}

func main() {
	var scope = flag.String("s", "", "影响范围，如：认证模块、学员模块、公共配置。不是必填，但是建议都进行描述")
	flag.Parse()
	if flag.NArg() < 2 {
		fmt.Printf("Usage: %s <变更类型> 提交消息 ...其他参数 \n变更类型: \n%s\n", os.Args[0], typeList)
		return
	}

	changeType := flag.Arg(0)

	if _, ok := types[changeType]; !ok {
		fmt.Printf("变更类型必须是以下类型之一:\n%s\n", typeList)
		return
	}
	var cmd []string
	cmd = append(cmd, "commit", "-m")
	if *scope != "" {
		cmd = append(cmd, fmt.Sprintf("%s(%s): %s", changeType, *scope, flag.Arg(0)))
	} else {
		cmd = append(cmd, fmt.Sprintf("%s: %s", changeType, flag.Args()[0]))
	}

	if flag.NArg() > 2 {
		cmd = append(cmd, flag.Args()[2:]...)
	}
	fmt.Println(strings.Join(cmd, ""))
	if ret, err := exec.Command("git", cmd...).CombinedOutput(); err != nil {
		fmt.Printf("%s\n", ret)
		log.Fatal(err)
	}
}
