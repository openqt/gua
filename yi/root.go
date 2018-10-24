package yi

import (
    "fmt"
    "os"

    "github.com/spf13/cobra"
    "hash/crc64"
    "log"
    "math/rand"
    "net"
    "time"
    "github.com/satori/uuid"
)

type ConfigType struct {
    ShowLog bool
    Random  bool
    Assign  string
}

var (
    conf ConfigType
)

var rootCmd = &cobra.Command{
    Use:              "gua",
    Short:            "易经卦象",
    PersistentPreRun: Config,
    Run: func(cmd *cobra.Command, args []string) {
        New().CalcSimple(conf.Assign).Tell()
    },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
    if err := rootCmd.Execute(); err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
}

func init() {
    cobra.OnInitialize()
}

// 屏蔽日志信息
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
