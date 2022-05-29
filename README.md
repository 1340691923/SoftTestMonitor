  

技术栈主要用：golang 

欢迎大家提出自己的issue。

SoftTestMonitor
-----------
[![Go Documentation](http://img.shields.io/badge/go-documentation-blue.svg?style=flat-square)](https://godoc.org/github.com/1340691923/SoftTestMonitor)
[![license](https://img.shields.io/github/license/mashape/apistatus.svg?maxAge=2592000)](https://github.com/1340691923/SoftTestMonitor/blob/master/LICENSE)
[![Release](https://img.shields.io/github/release/1340691923/ElasticView.svg?label=Release)](https://github.com/1340691923/SoftTestMonitor/releases/latest)
> SoftTestMonitor 软考成绩监控快查工具(该工具仅供学习参考)
 * 通过配置文件或命令行输入信息手动查询往年成绩
 * 监控软考平台，出成绩后将第一时间查询出成绩并通过邮件发送给您
 
## 应用程序下载
[下载地址]( https://github.com/1340691923/SoftTestMonitor/releases/) 支持操作系统：windows，linux

## 安装教程
>安装教程
 * 第一步: 下载release里面的对应压缩包，解压获取对应操作系统的执行程序 SoftTestMonitor
 
 ## 使用教程
 > 手动查询成绩
 * 第一步：下载配置文件 运行 ./SoftTestMonitor downloadConfig -c 
 * 第二步：根据help里的内容完善json配置文件 运行 ./SoftTestMonitor manuallyQueryScore -h 如 year->考试时间 配置中year则可填 2021年上半年
 * 第三步: 查询成绩 运行 ./SoftTestMonitor manuallyQueryScore -c config.json 即可查询出成绩输出并发送邮件至用户邮箱
 
  > 监听软考网站并在第一时间自动查询成绩发送至用户邮箱
  * 第一步：下载配置文件 运行 ./SoftTestMonitor downloadConfig -c 
  * 第二步：根据help里的内容完善json配置文件 运行 ./SoftTestMonitor monitor -h 如 time_interval->轮询时间间隔（单位为分）可填 数字 3
  * 第三步: 查询成绩发送至用户邮箱 运行 ./SoftTestMonitor monitor -c config.json 即可常驻内存，轮询软考网，当出考试成绩时自动查询成绩输出并发送邮件至用户邮箱（推荐将其置于后台运行 nohup ./SoftTestMonitor monitor -c config.json > monitor.log &）
