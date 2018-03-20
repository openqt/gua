package main

import (
	"github.com/openqt/gua/yi"
	"github.com/spf13/cobra"
	"hash/crc64"
	"log"
	"math/rand"
	"time"
)

type ConfigType struct {
	showLog bool
}

var (
	config ConfigType
)

////////////////////////////////////////////////////////////
// 伪输出屏蔽日志信息
type DummyIO struct {
}

func (_ DummyIO) Write(_ []byte) (int, error) {
	return 0, nil
}

func Config(_ *cobra.Command, args []string) {
	if !config.showLog {
		log.SetOutput(DummyIO{})
	}

	crc := crc64.New(crc64.MakeTable(crc64.ISO))
	if len(args) == 0 { // 没有输入数据用当前时间（小时）
		crc.Write([]byte(time.Now().Format("2006010215")))
	} else { // 有输入数据则取所有数据的值卜算
		log.Println(args)
		for _, t := range args {
			crc.Write([]byte(t))
		}
	}
	log.Println("CRC:", crc.Sum64())
	rand.Seed(int64(crc.Sum64()))
}

////////////////////////////////////////////////////////////

func main() {
	rootCmd := &cobra.Command{
		Use:              "gua",
		Short:            "易经卦象",
		PersistentPreRun: Config,
		Run: func(cmd *cobra.Command, args []string) {
			yi.New().CalcSimple(args).Tell()
		},
	}
	rootCmd.PersistentFlags().BoolVarP(&config.showLog, "verbose", "v", false, "输出计算过程")

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err.Error())
	}
}
