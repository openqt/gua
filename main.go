package main

import (
	"github.com/openqt/gua/yi"
	"github.com/spf13/cobra"
	"hash/crc64"
	"log"
	"math/rand"
	"time"
	"github.com/satori/uuid"
)

type ConfigType struct {
	ShowLog bool
	Random  bool
}

var (
	conf ConfigType
)

////////////////////////////////////////////////////////////
// 伪输出屏蔽日志信息
type DummyIO struct {
}

func (_ DummyIO) Write(_ []byte) (int, error) {
	return 0, nil
}

func Config(_ *cobra.Command, args []string) {
	if !conf.ShowLog {
		log.SetOutput(DummyIO{})
	}

	crc := crc64.New(crc64.MakeTable(crc64.ISO))
	crc.Write([]byte(time.Now().Format("2006010215")))
	if len(args) > 0 {
		log.Println(args)
		for _, t := range args {
			crc.Write([]byte(t))
		}
	}

	if conf.Random { // 增加随机UUID
		data, _ := uuid.NewV4()
		crc.Write(data.Bytes())
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
	rootCmd.PersistentFlags().BoolVarP(&conf.ShowLog, "verbose", "v", false, "输出计算过程")
	rootCmd.PersistentFlags().BoolVarP(&conf.Random, "random", "r", false, "随机占卜")

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err.Error())
	}
}
