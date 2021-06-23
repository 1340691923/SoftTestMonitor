
package cmd

import (
	"fmt"

	"github.com/1340691923/SoftTestMonitor/lib"
	"github.com/spf13/cobra"
)

//配置文件名字
var configFileName string

//下载配置文件命令cmd
var downloadConfigCmd = &cobra.Command{
	Use:   "downloadConfig",
	Short: "下载配置文件",
	Long: `下载配置文件,将会下载为当前目录的config.json文件`,
	Run: func(cmd *cobra.Command, args []string) {
		fileName := configFileName+".json"
		//下载配置文件
		err := lib.DownloadConfigFile(fileName)
		if err!=nil{
			fmt.Println(err)
			return
		}
		fmt.Println(fmt.Sprintf("文件下载完毕，名字为:%v",fileName))
	},
}


func init() {
	downloadConfigCmd.Flags().StringVarP(&configFileName, "configFileName", "c", "config", "参数配置文件名(后缀不用写,为固定的.json)")
	rootCmd.AddCommand(downloadConfigCmd)
}
