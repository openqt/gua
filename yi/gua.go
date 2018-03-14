package yi

import (
	"encoding/json"
	"fmt"
	"github.com/olekukonko/tablewriter"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"
)

const (
	TOTAL  = 49
	UNIT   = 4
	SEQLEN = 6
	YANG   = "▅▅▅▅▅"
	YIN    = "▅▅ ▅▅"
)

type YaoType struct {
	Image string
	Text  string
}

type GuaType struct {
	No [SEQLEN]int // 卜算数字

	Index string          // 卦序
	Name  string          // 卦名
	Desc  string          // 介绍
	Yao   [SEQLEN]YaoType // 六爻
}

var YiData map[string]GuaType // 易经数据

func (g *GuaType) Show() {
	data := YiData[g.GetIndex()]
	fmt.Printf("%s 【卦%s】\n", data.Name, data.Index)

	tb := tablewriter.NewWriter(os.Stdout)
	tb.SetHeader([]string{"卦象", "爻辞"})
	for _, yao := range data.Yao {
		tb.Append([]string{yao.Image, yao.Text})
	}
	tb.Render()
}

// 变卦
func (g *GuaType) Change() GuaType {
	gc := GuaType{}
	for i, n := range g.No {
		switch n {
		case 9: // 老阳
			gc.No[i] = 6
		case 7: // 少阳
			gc.No[i] = n
		case 6: // 老阴
			gc.No[i] = 9
		case 8: // 少阴
			gc.No[i] = n
		}
	}
	return gc
}

func (g *GuaType) Calc(args []string) {
	if len(args) == SEQLEN {
		var err error
		for i := 0; i < SEQLEN; i++ { // 自下至上，从0到5
			if g.No[i], err = strconv.Atoi(args[i]); err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}
		}
	} else {
		rand.Seed(time.Now().UnixNano())
		for i := 0; i < SEQLEN; i++ { // 自下至上，从0到5
			g.No[i] = calcYao()
		}
		log.Println(g.No)
	}
}

func (g *GuaType) GetIndex() string {
	idx := ""
	for _, n := range g.No {
		idx += strconv.Itoa(n % 2)
	}
	log.Printf("Index [%s]\n", idx)
	return idx
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

func Load(filename string) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}

	if err := json.Unmarshal(data, &YiData); err != nil {
		log.Fatal(err)
	}
}
