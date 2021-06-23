package cmd

import (
	"fmt"

	"github.com/1340691923/SoftTestMonitor/lib"
	"github.com/spf13/cobra"
)

//查看软考网站公开的所有考试时间cmd
var showExaminationTimeCmd = &cobra.Command{
	Use:   "showExaminationTimeCmd",
	Short: "查看软考网站公开的所有考试时间",
	Long: `查看软考网站公开的所有考试时间`,
	Run: runShowExaminationTimeCmd,
}

func init() {
	rootCmd.AddCommand(showExaminationTimeCmd)
}

//运行查看软考网站公开的所有考试时间命令
func runShowExaminationTimeCmd(cmd *cobra.Command, args []string){
	//设置cookie
	ruanKao := new(lib.RuanKao)
	err := ruanKao.SetWelComeCookie()
	if err!=nil{
		fmt.Println(err)
		return
	}
	//获取所有的考试时间列表
	fmt.Println("cookie获取成功")
	list,err := ruanKao.GetExaminationTimeList()
	if err!=nil{
		fmt.Println(err)
		return
	}
	fmt.Println("软考网站公开的所有考试时间如下:")
	for _,t := range list {
		fmt.Println(t)
	}

}
