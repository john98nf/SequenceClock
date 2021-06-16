// Copyright © 2021 Giannis Fakinos
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	appDirectory string
	RootCmd      = &cobra.Command{
		Use:   "sc",
		Short: "A latency targeting tool for serverless sequences of fuctions.",
		Long:  helpMessage,
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Hello from SequenceClock!")
		},
	}
)

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	RootCmd.AddCommand(create)
	RootCmd.AddCommand(check)
	RootCmd.AddCommand(delete)
	RootCmd.AddCommand(version)
}

func initConfig() {
	home, err := homedir.Dir()
	cobra.CheckErr(err)
	appDirectory = home + "/.sequenceClock"
	if _, err1 := os.Stat(appDirectory); os.IsNotExist(err1) {
		if err2 := os.Mkdir(appDirectory, 0755); err2 != nil {
			log.Fatal(err2)
		}
	}
	// Search config in home directory with name
	// ".cobra" (without extension).
	viper.AddConfigPath(appDirectory)
	viper.SetConfigType("yaml")
	viper.SetConfigName("config")

	viper.AutomaticEnv()

	if err3 := viper.ReadInConfig(); err3 != nil {
		log.Fatal(err3)
	}
}

var version = &cobra.Command{
	Use:   "version",
	Short: "Print the version of SequenceClock tool.",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("SequenceClock v1.0 - Copyright © 2021 Giannis Fakinos")
	},
}
