包主要的

进口(
_"嵌入"
	"fmt"
	"omega_launcher/cqhttp"
	"omega_launcher/fastbuilder"
	"Ω_发射器/发射器"
	"Ω_发射器/植物形式"
	"omega_launcher/utils"
	"os"
	"路径/文件路径"
	"运行时/调试"
	"时间"

	"github.com/PTerm/PTerm"
	"golang.org/x/term"
)

//go：嵌入版本
var版本[]字节

funcbeforeClose(){
	// 打印错误
err：=恢复（）
	如果err！=零{
PTerm.致命的.WithFatal(假的).println(err)
		//让投稿人开心
调试.printStack()
}
	如果P：=plantform.GetPlantform()；p==plantform.Windows_amd64||p==plantform.Windows_arm64{
		//让Windows用户满意
时间.睡眠(时间.第二次*5)
}其他{
		//让Unix用户满意
术语.MakeRaw(0)
}
}

func主要的(){
	推迟beforeClose()
	// 确保目录可用
	如果err：=os.chdir(utils.GetCurrentDir())；err！=零{
		恐慌(err)
}
	// 启动器自更新 (异步)
	去发射器.CheckUpdate(线(版本))
	// 启动
	// 读取配置
launcherConfig：=launcher(&S).配置{}
utils.GetJsonData(文件路径.加入(utils.GetCurrentDataDir()，"SnowConfig.json")，launcherConfig)
	// 添加启动信息
PTerm.DefaultBox.println("https://Snow.fastbuilder.icu/SnowLotus/")
PTerm.信息.println("Omega启动器"+PTerm.黄色的("(仅旧版Omega)")+PTerm.黄色的(" (",线(版本)，")"))
PTerm.信息.println("作者：CMA2401PT，SnowLotus修改")
	// 询问是否使用上一次的配置
	如果fastbuilder.CheckExecFile()&&launcherConfig.租赁代码！=“”{
		如果结果，_：=utils.GetInputYNInTime("是否使用上一次的登录配置？",10)；结果{
			//更新FB
			如果launcherConfig.UpdateFB{
fastbuilder.更新(launcherConfig)
}
			//go-cqhttp
			如果launcherConfig.EnableCQHttp&&launcherConfig.StartOmega{
cqhttp.运行(launcherConfig)
}
			//启动Omega或者FB
fastbuilder.运行(launcherConfig)
			返回
}
}
	//配置FB更新
	如果launcherConfig.UpdateFB=utils.GetInputYN("是否需要下载或更新菲尼克斯建筑公司？”)；launcherConfig.UpdateFB{
fastbuilder.UpdateRepo(launcherConfig)
fastbuilder.更新(launcherConfig)
}
	//检查是否已下载FB
	如果！fastbuilder.CheckExecFile(){
PTerm.警告.Printfln("错误"+plantform.GetFastBuilderName()+"没有FastBuilder")
fastbuilder.UpdateRepo(launcherConfig)
fastbuilder.更新(launcherConfig)
}
	//配置FB
fastbuilder.FBTokenSetup(launcherConfig)
	//配置租赁服登录(如果不为空且选择使用上次配置，则跳过设置)
	如果！(launcherConfig.租赁代码！=“”utils(&&U).GetInputYN(fmt.sprintf("是否使用上一次的%s的租赁服登陆配置？"，launcherConfig.RentalCode))){
launcherConfig.RentalCode=utils.GetValidInput("输入服务器号")
launcherConfig.RentalPasswd=utils.GetPswInput("输入服务器密码")
}
	//询问是否使用欧米茄
	如果launcherConfig.StartOmega=utils.GetInputYN("要启动Omega吗？")；launcherConfig.StartOmega{
		// 配置群服互通
		如果launcherConfig.EnableCQHttp=utils.GetInputYN("要启动go-cqhttp/群服互通吗？")；launcherConfig.EnableCQHttp{
			如果！utils.ISDIR(文件路径.加入(fastbuilder.GetOmegaStorageDir()，"配置")) {
				如果launcherConfig.EnableCQHttp=utils.GetInputYN("此时配置go-cqhttp会导致新生成的组件均为非启用状态，要继续吗？")；！launcherConfig.EnableCQHttp{
					//直接启动Omega或者FB
fastbuilder.运行(launcherConfig)
					返回
}
}
launcherConfig.BlockCQHttpOutput=utils.GetInputYN("要在配置完成后屏蔽go-cqhttp的输出吗？")
cqhttp.CQHttpEnablerHelper()
launcherConfig.EnableSignServer=utils.GetInputYN("要启动器启动签名服务器吗？")
cqhttp.运行(launcherConfig)
}
}
	//启动Omega或者FB
fastbuilder.运行(launcherConfig)
}
