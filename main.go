package main

import (
	"flag"
	"fmt"
	"log"
	"os/exec"
	"strings"
)

var types map[string]string
var typeList string

func init() {
	types = make(map[string]string)
	types["feat"] = "新功能"
	types["fix"] = "修复 bug"
	types["docs"] = "文档注释"
	types["style"] = "代码格式化（不影响运行时变动）"
	types["refactor"] = "重够（即不增加新功能，也不是修复 Bug）"
	types["perf"] = "性能优化"
	types["test"] = "增加测试"
	types["chore"] = "构建过程或辅助工具变动，如CICD修改，库升级"
	types["revert"] = "回退"
	types["build"] = "打包"

	sl := make([]string, 0)
	for k, v := range types {
		sl = append(sl, fmt.Sprintf("%8s: %s", k, v))
	}

	typeList = strings.Join(sl, "\n")
}

func main() {
	var scope = flag.String("s", "", "影响范围，如：认证模块、学员模块、公共配置。不是必填，但是建议都进行描述")
	var changeType = flag.String("t", "fix", fmt.Sprintf("变更类型，可选值:\n%s\n", typeList))
	flag.Parse()

	if _, ok := types[*changeType]; !ok {
		fmt.Printf("changeType must be in: %s\n", typeList)
		return
	}
	if flag.NArg() < 1 {
		fmt.Println("必须指定提交说明，如： gcm 修改用户认证时无法登录的问题")
		return
	}

	var args []string
	args = append(args, "commit", "-m")
	if *scope != "" {
		args = append(args, fmt.Sprintf("%s(%s): %s", *changeType, *scope, flag.Arg(0)))
	} else {
		args = append(args, fmt.Sprintf("%s: %s", *changeType, flag.Args()[0]))
	}

	if flag.NArg() > 1 {
		args = append(args, flag.Args()[1:]...)
	}
	if ret, err := exec.Command("git", args...).CombinedOutput(); err != nil {
		fmt.Printf("%s\n", ret)
		log.Fatal(err)
	}
}
