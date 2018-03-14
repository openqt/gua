package main

import (
	"fmt"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"
)

const (
	TOTAL = 49
	UNIT  = 4
	YANG  = "▅▅▅▅▅"
	YIN   = "▅▅ ▅▅"
)

type ConfigType struct {
	showLog bool
}

type GuaType struct {
	Origin, Change [6]int // 本卦，变卦
}

type YaoType struct {
	No   int
	Pos  string
	Yao  string
	Desc string
}

var (
	config ConfigType
)

func GetY(n int) string {
	switch n {
	case 7, 9:
		return "九"
	case 6, 8:
		return "六"
	}
	return "-"
}

func GetPos(n, i int) string {
	var pos string
	switch i {
	case 0:
		pos = "初" + GetY(n)
	case 1:
		pos = GetY(n) + "二"
	case 2:
		pos = GetY(n) + "三"
	case 3:
		pos = GetY(n) + "四"
	case 4:
		pos = GetY(n) + "五"
	case 5:
		pos = "上" + GetY(n)
	}

	return pos
}
func (g *GuaType) Show(origin bool) {
	var ori, chn [6]YaoType
	for i, n := range g.Origin {
		switch n {
		case 9: // 老阳
			ori[i] = YaoType{No: n, Yao: YANG, Pos: GetPos(n, i)} // 阳
			chn[i] = YaoType{No: n, Yao: YIN, Pos: GetPos(6, i)}  // 变
		case 7: // 少阳
			ori[i] = YaoType{No: n, Yao: YANG, Pos: GetPos(n, i)} // 阳
			chn[i] = YaoType{No: n, Yao: YANG, Pos: GetPos(n, i)} // 阳
		case 6: // 老阴
			ori[i] = YaoType{No: n, Yao: YIN, Pos: GetPos(n, i)}  // 阴
			chn[i] = YaoType{No: n, Yao: YANG, Pos: GetPos(9, i)} // 变
		case 8: // 少阴
			ori[i] = YaoType{No: n, Yao: YIN, Pos: GetPos(n, i)} // 阴
			chn[i] = YaoType{No: n, Yao: YIN, Pos: GetPos(n, i)} // 阴
		}
	}

	tb := tablewriter.NewWriter(os.Stdout)
	tb.SetHeader([]string{"位置", "卦象", "爻辞"})
	for i := 5; i >= 0; i-- {
		if origin {
			tb.Append([]string{ori[i].Pos, ori[i].Yao, ori[i].Desc})
		} else {
			tb.Append([]string{chn[i].Pos, chn[i].Yao, chn[i].Desc})
		}
	}
	tb.Render()
}

func (g *GuaType) Calc(args []string) {
	if len(args) == 6 {
		var err error
		for i := 0; i < 6; i++ { // 自下至上，从0到5
			if g.Origin[i], err = strconv.Atoi(args[i]); err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}
		}
	} else {
		rand.Seed(time.Now().UnixNano())
		for i := 0; i < 6; i++ { // 自下至上，从0到5
			g.Origin[i] = calcYao()
		}
		log.Println(g.Origin)
	}
}

// 变
func loop(total int) int {
	// 天地
	l := rand.Intn(total-UNIT*2) + UNIT
	r := total - l

	// 取人
	if rand.Intn(2) == 0 {
		l -= 1
	} else {
		r -= 1
	}

	// 取天
	lm := l % UNIT
	if lm == 0 {
		lm = UNIT
	}
	l -= lm

	// 取地
	rm := r % UNIT
	if rm == 0 {
		rm = UNIT
	}
	r -= rm

	// 取余
	m := lm + rm + 1

	log.Printf("总：%d，天：%d，地：%d，取：%d\n", total, l, r, m)
	return m
}

func calcYao() int {
	b1 := loop(TOTAL)           // 一变
	b2 := loop(TOTAL - b1)      // 二变
	b3 := loop(TOTAL - b1 - b2) // 三变

	m := TOTAL - b1 - b2 - b3
	y := m / UNIT // 爻
	log.Printf("余：%d，爻：%d\n", m, y)

	return y
}

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
		Use:   "eight",
		Short: "易经卦象",
		Run: func(cmd *cobra.Command, args []string) {
			g := GuaType{}
			g.Calc(args)
			fmt.Println("本卦")
			g.Show(true)

			fmt.Println("变卦")
			g.Show(false)
		},
		PersistentPreRun: Conf,
	}
	rootCmd.PersistentFlags().BoolVarP(&config.showLog, "verbose", "v", false, "Show more information")

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err.Error())
	}
}
