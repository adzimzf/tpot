package main

import (
	"errors"
	"log"
	"os"

	"github.com/adzimzf/tpot/config"
	scapper "github.com/adzimzf/tpot/scrapper"
	"github.com/adzimzf/tpot/tsh"
	"github.com/adzimzf/tpot/ui"
	"github.com/spf13/cobra"
)

func main() {
	rootCmd.Flags().BoolP("refresh", "r", false, "refresh the node list")
	rootCmd.Flags().BoolP("cfg", "c", false, "add config")
	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("failed to execute :%v", err)
	}
}

var rootCmd = &cobra.Command{
	Use:     "tpot <environment>",
	Short:   "tpot is tsh teleport wrapper",
	Long:    `config file is inside ` + os.Getenv("HOME") + `/.tpot/`,
	Example: "tpot staging\ntpot prod -r",
	Run: func(cmd *cobra.Command, args []string) {
		isCfg, err := cmd.Flags().GetBool("cfg")
		if err != nil {
			cmd.PrintErr("failed to get config due to ", err)
			return
		}
		if isCfg {
			addConfig(cmd, args)
			return
		}

		if len(args) < 1 {
			cmd.Help()
			return
		}

		isRefresh, err := cmd.Flags().GetBool("refresh")
		if err != nil {
			cmd.PrintErr("failed to get config due to ", err)
			return
		}

		proxy, err := config.NewProxy(args[0])
		if errors.Is(err, config.ErrEnvNotFound) {
			cmd.PrintErrf("Env %s not found\n", args[0])
			cmd.Help()
			return
		}
		if os.IsNotExist(err) {
			cmd.PrintErrf("Config not found\nrun tpot -c to add new proxy config")
			return
		}
		if err != nil {
			cmd.PrintErr("failed to get config due to ", err)
			return
		}

		if isRefresh {
			proxy.Node, err = scapper.NewScrapper(*proxy).GetNodes()
			if err != nil {
				cmd.PrintErr("failed to get nodes ", err)
				return
			}
			go proxy.UpdateNode(proxy.Node)
		} else {
			err := proxy.LoadNode()
			if err != nil {
				cmd.PrintErrf("failed to load nodes %v,\nyour might need -r to refresh/add the node cache", err)
				return
			}
		}

		var pItems []string
		for _, n := range proxy.Node.Items {
			pItems = append(pItems, n.Hostname)
		}

		host := ui.GetSelectedHost(pItems)
		if host == "" {
			cmd.PrintErrf("Pick at least one host to login")
			return
		}

		err = tsh.NewTSH(proxy, host).Run()
		if err != nil {
			cmd.PrintErr(err)
		}
	},
}

func addConfig(cmd *cobra.Command, args []string) {
	err := config.AddConfig()
	if err != nil {
		cmd.PrintErr(err)
	}
}
