package defines

// 启动器配置文件结构
type LauncherConfig struct {
	Repo              int    `json:"仓库序号"`
	RentalCode        string `json:"租赁服号"`
	RentalPasswd      string `json:"租赁服密码"`
	FBToken           string `json:"FBToken"`
	EnableCQHttp      bool   `json:"是否开启go-cqhttp"`
	StartOmega        bool   `json:"是否启动Omega"`
	UpdateFB          bool   `json:"是否更新FB"`
	BlockCQHttpOutput bool   `json:"是否屏蔽go-cqhttp的内容输出"`
}
