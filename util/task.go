package util

const (
	KVCSSvn string = "svn"
	KVCSGit string = "git"
	KVCSHg  string = "hg"
	KVCSDir string = "dir"
)

type TaskGroup struct {
	Version string

	Tasks []*Task
}

type Task struct {
	VCS         string
	Repository  string
	ReposDir    string
	User        string
	Password    string
	BuildScript string
	OutBin      string
	BuildLog    string
}

func NewTask() *Task {
	t := &Task{
		VCS:        KVCSSvn,
		Repository: "https://10.200.1.29:20443/JRTZXM/01-%E5%B9%BF%E8%B4%B5%E8%BF%90%E8%90%A5/02-%E4%BA%8C%E6%9C%9F%E4%BA%A4%E6%98%93%E7%B3%BB%E7%BB%9F/02-%E4%BB%A3%E7%A0%81/08-%E6%B8%85%E7%AE%97%E4%B8%AD%E5%BF%83%E4%BB%A3%E7%A0%81/03-%E9%93%B6%E8%A1%8C%E6%8E%A5%E5%8F%A3(%E6%B8%85%E7%AE%97%E4%B8%AD%E5%BF%83)/01-%E9%93%B6%E8%A1%8C%E6%9C%8D%E5%8A%A1/bank_svc",
		ReposDir:   "src/bankcc_svronline",
		User:       "luoguochun",
		Password:   "123",
		BuildScript: `autoconf;
        configure
        make`,
	}
	return t
}
