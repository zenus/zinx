package args

import (
	"github.com/zenus/zinx/zutils/commandline/uflag"
	"os"
	"path/filepath"
)

type args struct {
	ExeAbsDir  string
	ExeName    string
	ConfigFile string
}

var (
	Args   = args{}
	isInit = false
)

func init() {
	exe := os.Args[0]

	pwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	Args.ExeAbsDir = pwd
	Args.ExeName = filepath.Base(exe)
}

func InitConfigFlag(defaultValue string, tips string) {
	if isInit {
		return
	}
	isInit = true

	uflag.StringVar(&Args.ConfigFile, "c", defaultValue, tips)
	return
}

func FlagHandle() {
	filePath, err := filepath.Abs(Args.ConfigFile)
	if err != nil {
		panic(err)
	}
	Args.ConfigFile = filePath
}
