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

func Conf(_ *cobra.Command, _ []string) {
	if !config.showLog {
		log.SetOutput(DummyIO{})
	}
}

////////////////////////////////////////////////////////////

func main() {
	rootCmd := &cobra.Command{
		Use:   "gua",
		Short: "易经卦象",
		Run: func(cmd *cobra.Command, args []string) {
			yi.Load()
			g := yi.GuaType{}

			g.CalcSimple(args)
			g.Show()

			gc := g.Change()
			gc.Show()

			g.Divining()
		},
		PersistentPreRun: Conf,
	}
	rootCmd.PersistentFlags().BoolVarP(&config.showLog, "verbose", "v", false, "Show more information")

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err.Error())
	}
}
