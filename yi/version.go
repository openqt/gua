package yi

import (
    "fmt"

    "github.com/spf13/cobra"
    "runtime"
)

var (
    AppVersion string
    GoVersion  string
    GitVersion string
    BuildTime  string
)

// versionCmd represents the appVersion command
var versionCmd = &cobra.Command{
    Use:   "version",
    Short: "Show Version",
    Long:  `Print the appVersion information.`,
    Run: func(cmd *cobra.Command, args []string) {
        fmt.Printf("Version: %s\n", AppVersion)
        fmt.Printf("Git:     %s\n", GitVersion)
        fmt.Printf("GoC:     %s\n", GoVersion)
        fmt.Printf("Build:   %s\n", BuildTime)
        fmt.Printf("Go:      %s\n", runtime.Version())
    },
}

func init() {
    rootCmd.AddCommand(versionCmd)
}
