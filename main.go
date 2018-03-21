package main

import (
	"github.com/openqt/gua/yi"
	"github.com/satori/uuid"
	"github.com/spf13/cobra"
	"hash/crc64"
	"log"
	"math/rand"
	"net"
	"time"
)

type ConfigType struct {
	ShowLog bool
	Random  bool
	Assign  string
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

func initHardwareAddr() []byte {
	hardwareAddr := make([]byte, 6)
	interfaces, err := net.Interfaces()
	if err == nil {
		for _, iface := range interfaces {
			if len(iface.HardwareAddr) >= 6 {
				copy(hardwareAddr[:], iface.HardwareAddr)
				return hardwareAddr
			}
		}
	}
	return hardwareAddr
}

func Config(_ *cobra.Command, args []string) {
	if !conf.ShowLog {
		log.SetOutput(DummyIO{})
	}

	crc := crc64.New(crc64.MakeTable(crc64.ISO))
	crc.Write([]byte(time.Now().Format("20060102"))) // Add date by day
	crc.Write(initHardwareAddr())                    // Add MAC address

	if len(args) > 0 {
		log.Println(args)
		for _, t := range args {
			crc.Write([]byte(t))
		}
	}

	if conf.Random { // 增加随机UUID
		data := uuid.NewV4()

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
			yi.New().CalcSimple(conf.Assign).Tell()
		},
	}
	rootCmd.PersistentFlags().BoolVarP(&conf.ShowLog, "verbose", "v", false, "输出计算过程")
	rootCmd.Flags().BoolVarP(&conf.Random, "random", "r", false, "随机占卜")
	rootCmd.Flags().StringVarP(&conf.Assign, "assign", "a", "", "指定占卜数据")

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err.Error())
	}
}
