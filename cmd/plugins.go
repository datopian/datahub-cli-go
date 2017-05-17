package cmd

import (
	"encoding/json"
	"fmt"
	"sort"
	"time"

	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(pluginsCmd)
}

var pluginsCmd = &cobra.Command{
	Use:        "plugins",
	Short:      "manage plugins",
	Long: `
Example:
  $ heroku plugins`,
	Run: func(cmd *cobra.Command, args []string) {
		pluginsList(ctx)
	},
}

type pluginPresenter struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type nameSorter []pluginPresenter

func (a nameSorter) Len() int           { return len(a) }
func (a nameSorter) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a nameSorter) Less(i, j int) bool { return a[i].Name < a[j].Name }

func pluginsList(ctx *Context) {
	var names []string
	var pluginPresenters []pluginPresenter
	for _, plugin := range UserPlugins.Plugins() {
		symlinked := ""
		if UserPlugins.isPluginSymlinked(plugin.Name) {
			symlinked = " (symlinked)"
		}
		names = append(names, fmt.Sprintf("%s %s%s", plugin.Name, plugin.Version, symlinked))
		plgPres := pluginPresenter{plugin.Name, fmt.Sprintf("%s%s", plugin.Version, symlinked)}
		pluginPresenters = append(pluginPresenters, plgPres)
	}
	if ctx.Flags["core"] != nil {
		UserPluginNames := UserPlugins.PluginNames()
		for _, plugin := range CorePlugins.Plugins() {
			if contains(UserPluginNames, plugin.Name) {
				continue
			}
			names = append(names, fmt.Sprintf("%s %s (core)", plugin.Name, plugin.Version))
			plgPres := pluginPresenter{plugin.Name, fmt.Sprintf("%s (core)", plugin.Version)}
			pluginPresenters = append(pluginPresenters, plgPres)
		}
	}

	if ctx.Flags["json"] != nil {
		sort.Sort(nameSorter(pluginPresenters))
		pluginsJSON, _ := json.Marshal(pluginPresenters)
		Println(string(pluginsJSON))
	} else {
		sort.Strings(names)
		for _, plugin := range names {
			Println(plugin)
		}
	}
}

// Plugins represents either core or user plugins
type Plugins struct {
	Path    string
	plugins []*Plugin
}

// CorePlugins are built in plugins
var CorePlugins = &Plugins{Path: filepath.Join(AppDir, "lib")}

// UserPlugins are user-installable plugins
var UserPlugins = &Plugins{Path: filepath.Join(DataHome, "plugins")}

// Plugin represents a javascript plugin
type Plugin struct {
	Name      string    `json:"name"`
	Tag       string    `json:"tag"`
	Version   string    `json:"version"`
	Commands  Commands  `json:"commands"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Commands lists all the commands of the plugins
func (p *Plugins) Commands() (commands Commands) {
	for _, plugin := range p.Plugins() {
		for _, command := range plugin.Commands {
			command.Run = p.runFn(plugin, command.Topic, command.Command)
			commands = append(commands, command)
		}
	}
	return
}
