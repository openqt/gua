package main

import (
	"github.com/openqt/gua/yi"
	"github.com/spf13/cobra"
	"log"
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

func Config(_ *cobra.Command, _ []string) {
	if !config.showLog {
		log.SetOutput(DummyIO{})
	}
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
