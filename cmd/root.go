package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "SoftTestMonitor",
	Short: "软考成绩监控快查工具(该工具仅供学习参考)",
	Long: `软考成绩监控快查工具，能帮您做的事:
0.通过配置文件或命令行输入信息手动查询往年成绩
1.监控软考平台，出成绩后将第一时间查询出成绩并通过邮件发送给您`,
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)
}

func initConfig() {

}
