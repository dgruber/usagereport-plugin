package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/cloudfoundry/cli/plugin"
	"github.com/krujos/usagereport-plugin/apihelper"
	"github.com/krujos/usagereport-plugin/models"
	"strings"
)

type globalQueryCache struct {
	siMap  map[string]apihelper.ServiceInstance
	spMap  map[string]apihelper.ServicePlan
	sMap   map[string]apihelper.Service
	upsMap map[string]apihelper.UserProvidedService
}

//UsageReportCmd the plugin
type UsageReportCmd struct {
	apiHelper  apihelper.CFAPIHelper
	queryCache globalQueryCache
}

// contains CLI flag values
type flagVal struct {
	OrgName              string
	SpaceName            string
	Format               string
	ShowServiceInstances bool
}

func ParseFlags(args []string) flagVal {
	flagSet := flag.NewFlagSet(args[0], flag.ContinueOnError)

	// Create flags
	orgName := flagSet.String("o", "", "-o orgName")
	spaceName := flagSet.String("s", "", "-s spaceName")
	showSI := flagSet.Bool("i", false, "-i <true|false>")
	format := flagSet.String("f", "format", "-f <csv>")

	err := flagSet.Parse(args[1:])
	if err != nil {
	}

	return flagVal{
		OrgName:              string(*orgName),
		SpaceName:            string(*spaceName),
		Format:               string(*format),
		ShowServiceInstances: bool(*showSI),
	}
}

// createQueryCache makes global REST queries just once
func (cmd *UsageReportCmd) createQueryCache() error {
	// get service instances
	siMap, err := cmd.apiHelper.GetServiceInstanceMap("/v2/service_instances")
	if err != nil {
		return err
	}

	// get service plan map
	spMap, err := cmd.apiHelper.GetServicePlanMap()
	if err != nil {
		return err
	}

	// get services (for determining the p-)
	sMap, err := cmd.apiHelper.GetServiceMap()
	if err != nil {
		return err
	}

	upsMap, err := cmd.apiHelper.GetUserProvidedServiceMap()
	if err != nil {
		return err
	}

	cmd.queryCache = globalQueryCache{
		siMap:  siMap,
		spMap:  spMap,
		sMap:   sMap,
		upsMap: upsMap,
	}
	return nil
}

//GetMetadata returns metatada
func (cmd *UsageReportCmd) GetMetadata() plugin.PluginMetadata {
	return plugin.PluginMetadata{
		Name: "usage-report",
		Version: plugin.VersionType{
			Major: 1,
			Minor: 5,
			Build: 0,
		},
		Commands: []plugin.Command{
			{
				Name:     "usage-report",
				HelpText: "Report AI and memory usage for orgs and spaces",
				UsageDetails: plugin.Usage{
					Usage: "cf usage-report [-o orgName] [-s spaceName] [-i <show service instances>] [-f <csv>]",
					Options: map[string]string{
						"o": "organization",
						"s": "space",
						"i": "serviceInstances",
						"f": "format",
					},
				},
			},
		},
	}
}

func (cmd *UsageReportCmd) getFilteredOrgs(orgName string) []models.Org {
	var orgs []models.Org

	if orgName != "" {
		org, err := cmd.getOrg(orgName)
		if nil != err {
			fmt.Println(err)
			os.Exit(1)
		}
		orgs = append(orgs, org)
	} else {
		var err error
		orgs, err = cmd.getOrgs()
		if nil != err {
			fmt.Println(err)
			os.Exit(1)
		}
	}
	return orgs
}

//UsageReportCommand doer
func (cmd *UsageReportCmd) UsageReportCommand(args []string) {
	flagVals := ParseFlags(args)

	var report models.Report

	report.Orgs = cmd.getFilteredOrgs(flagVals.OrgName)

	// process service instances
	if flagVals.ShowServiceInstances {
		if flagVals.Format == "csv" {
			fmt.Println(report.ServiceInstanceReportCSV())
		} else {
			fmt.Println(report.ServiceInstanceReportString())
		}
	} else {
		if flagVals.Format == "csv" {
			fmt.Println(report.CSV())
		} else {
			fmt.Println(report.String())
		}
	}
}

func (cmd *UsageReportCmd) getOrgs() ([]models.Org, error) {
	if err := cmd.createQueryCache(); err != nil {
		return nil, err
	}

	rawOrgs, err := cmd.apiHelper.GetOrgs()
	if nil != err {
		return nil, err
	}

	var orgs = []models.Org{}

	for _, o := range rawOrgs {
		orgDetails, err := cmd.getOrgDetails(o)
		if err != nil {
			return nil, err
		}
		orgs = append(orgs, orgDetails)
	}
	return orgs, nil
}

func (cmd *UsageReportCmd) getOrg(name string) (models.Org, error) {
	rawOrg, err := cmd.apiHelper.GetOrg(name)
	if nil != err {
		return models.Org{}, err
	}

	return cmd.getOrgDetails(rawOrg)
}

func (cmd *UsageReportCmd) getOrgDetails(o apihelper.Organization) (models.Org, error) {
	usage, err := cmd.apiHelper.GetOrgMemoryUsage(o)
	if nil != err {
		return models.Org{}, err
	}
	quota, err := cmd.apiHelper.GetQuotaMemoryLimit(o.QuotaURL)
	if nil != err {
		return models.Org{}, err
	}
	spaces, err := cmd.getSpaces(o.SpacesURL)
	if nil != err {
		return models.Org{}, err
	}

	return models.Org{
		Name:        o.Name,
		MemoryQuota: int(quota),
		MemoryUsage: int(usage),
		Spaces:      spaces,
	}, nil
}

func (cmd *UsageReportCmd) getSpaces(spaceURL string) ([]models.Space, error) {
	rawSpaces, err := cmd.apiHelper.GetOrgSpaces(spaceURL)
	if nil != err {
		return nil, err
	}

	var spaces = []models.Space{}
	for _, s := range rawSpaces {
		apps, err := cmd.getApps(s.AppsURL)
		if nil != err {
			return nil, err
		}
		spaces = append(spaces,
			models.Space{
				Apps: apps,
				Name: s.Name,
			},
		)
	}
	return spaces, nil
}

// IsPCFInstance checks if a particular service instance is using a PCF service.
func IsPCFInstance(serviceInstanceGUID string, siMap map[string]apihelper.ServiceInstance, spMap map[string]apihelper.ServicePlan, sMap map[string]apihelper.Service) bool {
	if serviceInstance, exists := siMap[serviceInstanceGUID]; exists == false {
		return false
	} else if servicePlan, exists := spMap[serviceInstance.ServicePlanGUID]; exists == false {
		return false
	} else if service, exists := sMap[servicePlan.ServiceGUID]; exists == false {
		return false
	} else {
		return strings.HasPrefix(service.Label, "p-")
	}
}

func (cmd *UsageReportCmd) getApps(appsURL string) ([]models.App, error) {
	rawApps, err := cmd.apiHelper.GetSpaceApps(appsURL)
	if nil != err {
		return nil, err
	}

	var apps = []models.App{}
	for _, a := range rawApps {

		// TODO check if that is available globally and can be cached
		sb, err := cmd.apiHelper.GetServiceBindings(a.ServiceBindingsURL)
		if err != nil {
			return nil, err
		}

		siTotal := len(sb)
		siPCF := 0 // PCF service instances
		siUP := 0  // User Provided Service Instances

		for _, binding := range sb {
			if si, exists := cmd.queryCache.siMap[binding.ServiceInstanceGUID]; exists {
				if si.Type == "managed_service_instance" {
					if IsPCFInstance(binding.ServiceInstanceGUID, cmd.queryCache.siMap, cmd.queryCache.spMap, cmd.queryCache.sMap) {
						siPCF++
					}
				}
			} else if _, exists := cmd.queryCache.upsMap[binding.ServiceInstanceGUID]; exists {
				siUP++
			}
		}

		apps = append(apps, models.App{
			Instances: int(a.Instances),
			Ram:       int(a.RAM),
			Running:   a.Running,
			Name:      a.Name,
			SiTotal:   siTotal,
			SiPCF:     siPCF,
			SiUP:      siUP,
		})
	}
	return apps, nil
}

//Run runs the plugin
func (cmd *UsageReportCmd) Run(cli plugin.CliConnection, args []string) {
	if args[0] == "usage-report" {
		cmd.apiHelper = apihelper.New(cli)
		cmd.UsageReportCommand(args)
	}
}

func main() {
	plugin.Start(new(UsageReportCmd))
}
