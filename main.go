package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

type ChangeType struct {
	Description string
	Alias       []string
}

// feat, fix, perf, style, docs, test, refactor, build, ci, chore, revert, types, release
// https://github.com/angular/angular/blob/main/contributing-docs/commit-message-guidelines.md
var changeTypes map[string]ChangeType = map[string]ChangeType{
	"build":    {Description: "影响构建系统或外部依赖的变更"},
	"ci":       {Description: "CI配置文件或脚本变更"},
	"docs":     {Description: "文档变更"},
	"feat":     {Description: "增加新功能"},
	"fix":      {Description: "修复问题/BUG"},
	"perf":     {Description: "优化/性能提升"},
	"refactor": {Description: "重构"},
	"test":     {Description: "测试相关"},
	"revert":   {Description: "撤销修改"},
	"style":    {Description: "代码风格相关无影响运行结果的"},
	"chore":    {Description: "对构建过程或辅助工具和库的更改 (不影响源文件、测试用例)"},
	"types":    {Description: "类型定义文件修改"},
	"release":  {Description: "发布版本"},
	"wip":      {Description: "正在开发中"},
	"workflow": {Description: "工作流程改进"},
}
var typeList string
var types string

func init() {
	sl := make([]string, 0)
	for k, v := range changeTypes {
		sl = append(sl, fmt.Sprintf("%8s: %s", k, v.Description))
	}
	keys := make([]string, 0, len(changeTypes))
	for k := range changeTypes {
		keys = append(keys, k)
	}
	types = strings.Join(keys, "|")

	typeList = strings.Join(sl, "\n")
}

func main() {
	var scope = flag.String("s", "", "影响范围，如：认证模块、学员模块、公共配置。不是必填，但是建议都进行描述")
	var changeType = flag.String("t", "", types)
	var verbose = flag.Bool("v", false, "详细显示")

	flag.Usage = func() {
		fmt.Printf("Usage: %s [选项] header [args...] \n", os.Args[0])
		fmt.Println("选项:")
		flag.PrintDefaults()
		fmt.Println("变更类型:")
		for k, v := range changeTypes {
			fmt.Printf("%8s: %s\n", k, v.Description)
		}
		fmt.Println("参考地址: https://github.com/conventional-changelog/commitlint/#what-is-commitlint")
		fmt.Println("参考格式:")
		fmt.Println(`   <type>(<scope>): <short summary>
      │       │             │
      │       │             └─⫸ Summary in present tense. Not capitalized. No period at the end.
      │       │
      │       └─⫸ Commit Scope: animations|bazel|benchpress|common|compiler|compiler-cli|core|
      │                          elements|forms|http|language-service|localize|platform-browser|
      │                          platform-browser-dynamic|platform-server|router|service-worker|
      │                          upgrade|zone.js|packaging|changelog|docs-infra|migrations|
      │                          devtools
      │
      └─⫸ Commit Type: build|ci|docs|feat|fix|perf|refactor|test`)
	}
	flag.Parse()
	if flag.NArg() < 1 {
		fmt.Println("未指定提交内容")
		flag.Usage()
		return
	}

	if _, ok := changeTypes[*changeType]; !ok {
		fmt.Printf("变更类型必须是以下类型之一:\n%s\n", typeList)
		return
	}

	header := flag.Arg(0)
	var cmd []string
	cmd = append(cmd, "commit")
	if *scope != "" {
		cmd = append(cmd, "-m", fmt.Sprintf("%s(%s): %s", *changeType, *scope, header))
	} else {
		cmd = append(cmd, "-m", fmt.Sprintf("%s: %s", *changeType, header))
	}
	if flag.NArg() > 1 {
		for _, arg := range flag.Args()[1:] {
			cmd = append(cmd, arg)
		}
	}

	if *verbose {
		fmt.Printf("git %s\n", strings.Join(cmd, " "))
	}
	if ret, err := exec.Command("git", cmd...).CombinedOutput(); err != nil {
		fmt.Printf("%s\n", ret)
		log.Fatal(err)
	}
}
