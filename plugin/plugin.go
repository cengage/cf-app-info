package plugin

import (
	"code.cloudfoundry.org/cli/plugin"
	"code.cloudfoundry.org/cli/plugin/models"
	"fmt"
	"os"
	"reflect"
	"strings"
)

func NewPlugin() *Plugin {
	return &Plugin{}
}

type Plugin struct{}

func fatalIf(err error) {
	if err != nil {
		fmt.Fprintln(os.Stdout, "error:", err)
		os.Exit(1)
	}
}

func formatRoute(route plugin_models.GetApp_RouteSummary) string {
	return fmt.Sprintf("%s.%s%s", route.Host, route.Domain.Name, route.Path)
}

func mapRoutes(routes []plugin_models.GetApp_RouteSummary) []string {
	vsm := make([]string, len(routes))
	for i, v := range routes {
		vsm[i] = formatRoute(v)
	}
	return vsm
}


func formatService(service plugin_models.GetApp_ServiceSummary) string {
	return fmt.Sprintf("%s:%s", service.Name, service.Guid)
}

func mapServices(routes []plugin_models.GetApp_ServiceSummary) []string {
	vsm := make([]string, len(routes))
	for i, v := range routes {
		vsm[i] = formatService(v)
	}
	return vsm
}

func getField(app plugin_models.GetAppModel, name string) string {
	if name == "Routes" {
		return strings.Join(mapRoutes(app.Routes), "\n")
	}

	if name == "Services" {
		return strings.Join(mapServices(app.Services), "\n")
	}

	r := reflect.ValueOf(app)
	f := reflect.Indirect(r).FieldByName(name)

	if !f.IsValid() {
		return ""
	}

	return fmt.Sprint(f)

}


func (p *Plugin) getAppInfo(cliConnection plugin.CliConnection, args []string) (error) {

	if len(args) != 2 {
		return fmt.Errorf("invalid parameters")
	}

	appName := args[0]
	propName := args[1]

	isLoggedIn, err := cliConnection.IsLoggedIn()
	if err != nil {
		return err
	}
	if !isLoggedIn {
		return fmt.Errorf("you need to log in")
	}

	app, err := cliConnection.GetApp(appName)
	if err != nil {
		return fmt.Errorf("couldn't get app %s: %s", appName, err)
	}

	fmt.Printf(getField(app, propName))

	return nil
}

// Run must be implemented by any plugin because it is part of the
// plugin interface defined by the core CLI.
//
// Run(....) is the entry point when the core CLI is invoking a command defined
// by a plugin. The first parameter, plugin.CliConnection, is a struct that can
// be used to invoke cli commands. The second paramter, args, is a slice of
// strings. args[0] will be the name of the command, and will be followed by
// any additional arguments a cli user typed in.
//
// Any error handling should be handled with the plugin itself (this means printing
// user facing errors). The CLI will exit 0 if the plugin exits 0 and will exit
// 1 should the plugin exits nonzero.

func (p *Plugin) Run(cliConnection plugin.CliConnection, args []string) {
	// only handle if actually invoked, else it can't be uninstalled cleanly
	if args[0] != "app-info" {
		return
	}

	err := p.getAppInfo(cliConnection, args[1:])
	fatalIf(err)
}

// GetMetadata must be implemented as part of the plugin interface
// defined by the core CLI.
//
// GetMetadata() returns a PluginMetadata struct. The first field, Name,
// determines the name of the plugin which should generally be without spaces.
// If there are spaces in the name a user will need to properly quote the name
// during uninstall otherwise the name will be treated as seperate arguments.
// The second value is a slice of Command structs. Our slice only contains one
// Command Struct, but could contain any number of them. The first field Name
// defines the command `cf basic-plugin-command` once installed into the CLI. The
// second field, HelpText, is used by the core CLI to display help information
// to the user in the core commands `cf help`, `cf`, or `cf -h`.
func (p *Plugin) GetMetadata() plugin.PluginMetadata {
	return plugin.PluginMetadata{
		Name: "app-info",
		Version: plugin.VersionType{
			Major: 0,
			Minor: 0,
			Build: 1,
		},
		MinCliVersion: plugin.VersionType{
			Major: 6,
			Minor: 7,
			Build: 0,
		},
		Commands: []plugin.Command{
			{
				Name:     "app-info",
				HelpText: "Plugin to get the application id",

				// UsageDetails is optional
				// It is used to show help of usage of each command
				UsageDetails: plugin.Usage{
					Usage: "$ cf app-info helloworld Guid" +
						"$ cf app-info helloworld Name" +
						"$ cf app-info helloworld BuildpackUrl" +
						"$ cf app-info helloworld Command" +
						"$ cf app-info helloworld Diego" +
						"$ cf app-info helloworld DetectedStartCommand" +
						"$ cf app-info helloworld DiskQuota" +
						"$ cf app-info helloworld InstanceCount" +
						"$ cf app-info helloworld Memory" +
						"$ cf app-info helloworld RunningInstances" +
						"$ cf app-info helloworld HealthCheckTimeout" +
						"$ cf app-info helloworld State" +
						"$ cf app-info helloworld SpaceGuid" +
						"$ cf app-info helloworld PackageUpdatedAt" +
						"$ cf app-info helloworld PackageState" +
						"$ cf app-info helloworld StagingFailedReason" +
						"$ cf app-info helloworld Routes" +
						"$ cf app-info helloworld Services",
				},
			},
		},
	}
}
