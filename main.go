package main

import (
	"errors"
	"fmt"
	"log"

	"github.com/adzimzf/tpot/config"
	"github.com/adzimzf/tpot/tsh"
	"github.com/adzimzf/tpot/ui"
	"github.com/spf13/cobra"
)

// Version wil be override during build
var Version = "DEV"

func main() {
	rootCmd.Flags().BoolP("refresh", "r", false, "Replace the node list from proxy")
	rootCmd.Flags().BoolP("append", "a", false, "Append the fresh node list to the cache")
	rootCmd.Flags().BoolP("config", "c", false, "show the configuration list")
	rootCmd.Flags().Bool("add", false, "add the teleport configuration")
	rootCmd.Flags().BoolP("version", "v", false, "show the tpot version")
	rootCmd.Flags().BoolP("edit", "e", false, "edit all or specific configuration")
	rootCmd.Flags().StringP("user", "u", "", "user to login to the desired host")
	rootCmd.Version = Version
	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("failed to execute :%v\n", err)
	}
}

const example = `
tpot -c --add         // Set up the configuration environment
tpot -c --edit        // Edit all the configuration
tpot staging          // Show the node list of staging environment
tpot staging --edit   // Edit the staging proxy configuration
tpot prod -a          // Get the latest node list then append to the cache for production 
tpot prod -r          // Refresh the cache with the latest node from Teleport UI
tpot prod -u root     // Login into production using root user
`

var rootCmd = &cobra.Command{
	Use:     "tpot <ENVIRONMENT>",
	Short:   "tpot is tsh teleport wrapper",
	Long:    `config file is inside ` + config.Dir,
	Example: example,
	Run: func(cmd *cobra.Command, args []string) {

		cfg, err := config.NewConfig()
		if err != nil {
			cmd.PrintErrln("failed to get config, error:", err)
			return
		}

		isCfg, err := cmd.Flags().GetBool("config")
		if err != nil {
			cmd.PrintErrln("failed to get config due to ", err)
			return
		}
		if isCfg {
			// the error has already beautify by the handler
			if err := configHandler(cmd, cfg); err != nil {
				cmd.PrintErrln(err)
			}
			return
		}

		if len(args) < 1 {
			cmd.Help()
			return
		}

		proxy, err := cfg.FindProxy(args[0])
		if errors.Is(err, config.ErrEnvNotFound) {
			cmd.PrintErrf("Env %s not found\n\n", args[0])
			cmd.Help()
			return
		}

		if err != nil {
			cmd.PrintErrln("failed to get config due to ", err)
			return
		}

		isEdit, err := cmd.Flags().GetBool("edit")
		if err != nil {
			return
		}
		if isEdit {
			if err := proxyEditHandler(cfg, proxy); err != nil {
				cmd.PrintErrln(err)
				return
			}
			cmd.Printf("%s has updated successfully\n", proxy.Env)
			return
		}

		node, err := handleNode(cmd, proxy)
		if err != nil {
			cmd.PrintErrln(err)
			return
		}

		host := ui.GetSelectedHost(node.ListHostname())
		if host == "" {
			cmd.PrintErrln("Pick at least one host to login")
			return
		}

		user, err := getUserLogin(cmd, node)
		if err != nil {
			cmd.PrintErrln(err)
			return
		}

		// print to give user information
		cmd.Printf("login using %s %s\n", user, host)

		err = tsh.NewTSH(proxy).SSH(user, host)
		if err != nil {
			cmd.PrintErrln(err)
		}
	},
}

func getUserLogin(cmd *cobra.Command, node *config.Node) (string, error) {
	userLogin, err := cmd.Flags().GetString("user")
	if err != nil {
		return "", err
	}
	if userLogin != "" {
		return userLogin, nil
	}

	if node.Status == nil {
		return "", fmt.Errorf("need to run using flag -a or -r to get the latest user login")
	}

	uiUser, err := ui.NewLoginUser(node.Status.UserLogins)
	if err != nil {
		return "", err
	}
	user, err := uiUser.Run()
	if err != nil {
		return "", err
	}
	if user == "" {
		return "", fmt.Errorf("user login must not be empty")
	}
	return user, nil
}

func proxyEditHandler(c *config.Config, proxy *config.Proxy) error {
	res, err := c.Edit(proxy.Env)
	if err != nil {
		fmt.Printf("failed to edit proxy, error: %v\n", err)
	}

	// if any changes, keep track any last changes until user confirm
	// don't want to continue edit
	for res != "" && err != nil {
		confirm, err := ui.Confirm("Do You want to continue edit")
		if err != nil {
			fmt.Printf("failed to get confirmation, error: %v\n", err)
			break
		}
		if !confirm {
			break
		}
		res, err = c.EditPlain(proxy.Env, res)
		if err != nil {
			fmt.Printf("failed to edit proxy, error: %v\n", err)
		}
		if err == nil {
			fmt.Printf("Success to edit proxy\n")
			break
		}
	}
	return nil
}

func configHandler(cmd *cobra.Command, c *config.Config) error {
	isEdit, err := cmd.Flags().GetBool("edit")
	if err != nil {
		return fmt.Errorf("failed to get flags edit, error: %v", err)
	}

	if isEdit {
		res, err := c.EditAll()
		if err != nil {
			fmt.Printf("failed to edit config, error: %v\n", err)
		}

		// if any changes, keep track any last changes until user confirm
		// don't want to continue edit
		for res != "" && err != nil {
			confirm, err := ui.Confirm("Do You want to continue edit")
			if err != nil {
				fmt.Printf("failed to get confirmation, error: %v\n", err)
				break
			}
			if !confirm {
				break
			}
			res, err = c.EditAllPlain(res)
			if err != nil {
				fmt.Printf("failed to edit config, error: %v\n", err)
			}
			if err == nil {
				fmt.Printf("Success to edit config\n")
				break
			}
		}
		return nil
	}

	isAdd, err := cmd.Flags().GetBool("add")
	if err != nil {
		return fmt.Errorf("failed to get flags edit, error: %v", err)
	}
	if isAdd {
		res, err := c.Add()
		if err != nil {
			fmt.Printf("failed to add config, error: %v\n", err)
		}

		// if any changes, keep track any last changes until user confirm
		// don't want to continue edit
		for res != "" && err != nil {
			confirm, err := ui.Confirm("Do You want to continue edit")
			if err != nil {
				fmt.Printf("failed to get confirmation, error: %v\n", err)
				break
			}
			if !confirm {
				break
			}
			res, err = c.AddPlain(res)
			if err != nil {
				fmt.Printf("failed to add config, error: %v\n", err)
			}
			if err == nil {
				fmt.Printf("Success to add config\n")
				break
			}
		}
		return nil
	}

	str, err := c.String()
	if err != nil {
		return fmt.Errorf("failed to get config string, error:%v", err)
	}
	fmt.Println(str)
	return nil
}

func handleNode(cmd *cobra.Command, proxy *config.Proxy) (*config.Node, error) {
	isRefresh, err := cmd.Flags().GetBool("refresh")
	if err != nil {
		return nil, err
	}
	isAppend, err := cmd.Flags().GetBool("append")
	if err != nil {
		return nil, err
	}
	var nodes config.Node
	if isRefresh || isAppend {
		nodes, err = getLatestNode(proxy, isAppend)
		if err != nil {
			return nil, err
		}
	} else {
		nodes, err = proxy.GetNode()
		if err != nil {
			return nil, fmt.Errorf("failed to load nodes %v,\nyour might need -r to refresh/add the node cache", err)
		}
	}

	return &nodes, nil
}

func getLatestNode(proxy *config.Proxy, isAppend bool) (config.Node, error) {
	t := tsh.NewTSH(proxy)
	nodes, err := t.ListNodes()
	if err != nil {
		return nodes, fmt.Errorf("failed to get nodes: %v", err)
	}
	if len(nodes.Items) == 0 {
		return nodes, fmt.Errorf("there's no nodes found")
	}

	if isAppend {
		nodes, err = proxy.AppendNode(nodes)
		if err != nil {
			return nodes, fmt.Errorf("failed to append nodes, err: %v", err)
		}
	}

	status, err := t.Status()
	if err != nil && err != tsh.ErrUnsupportedVersion {
		return nodes, err
	}

	// if the tsh version is not supported
	// just hardcoded the user login to root for now
	if err == tsh.ErrUnsupportedVersion {
		version, err := t.Version()
		if err != nil {
			return config.Node{}, err
		}

		fmt.Printf("WARNING! minimum tsh version is Teleport v2.6.1 but got %s, the user login list is will be only root\n", version.Strings())
		status = &config.ProxyStatus{
			UserLogins: []string{"root"},
		}
	}

	// append the status to node
	nodes.Status = status
	go proxy.UpdateNode(nodes)
	return nodes, nil
}
