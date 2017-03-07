package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/cloudfoundry/cli/plugin"
	"github.com/dgruber/usagereport-plugin/apihelper"
	"github.com/dgruber/usagereport-plugin/models"
	"strings"
)

type globalQueryCache struct {
	siMap    map[string]apihelper.ServiceInstance
	spMap    map[string]apihelper.ServicePlan
	sMap     map[string]apihelper.Service
	upsMap   map[string]apihelper.UserProvidedService
	spaceMap map[string]apihelper.SpaceDetails
	orgMap   map[string]apihelper.OrgDetails
	sbList   []apihelper.ServiceBinding
	sbMap    map[string][]string
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
	ShowServiceInstances string
}

func ParseFlags(args []string) flagVal {
	flagSet := flag.NewFlagSet(args[0], flag.ExitOnError)

	// Create flags
	orgName := flagSet.String("o", "", "-o orgName")
	spaceName := flagSet.String("s", "", "-s spaceName")
	showSI := flagSet.String("i", "", "-i <app|summary>")
	format := flagSet.String("f", "format", "-f csv")

	err := flagSet.Parse(args[1:])
	if err != nil {
		os.Exit(2)
	}

	if *showSI != "" && *showSI != "app" && *showSI != "summary" {
		fmt.Fprintf(os.Stderr, "-i requires to be either \"app\" or \"summary\" if set.\n")
		os.Exit(2)
	}

	return flagVal{
		OrgName:              string(*orgName),
		SpaceName:            string(*spaceName),
		Format:               string(*format),
		ShowServiceInstances: string(*showSI),
	}
}

// createQueryCache makes global REST queries just once and stores them as a cache.
func (cmd *UsageReportCmd) createQueryCache() error {
	// get service instances
	siMap, err := cmd.apiHelper.GetServiceInstanceMap()
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

	spaceMap, err := cmd.apiHelper.GetSpaceMap()
	if err != nil {
		return err
	}

	orgMap, err := cmd.apiHelper.GetOrgMap()
	if err != nil {
		return err
	}

	sbList, err := cmd.apiHelper.GetServiceBindingsList()
	if err != nil {
		return err
	}

	// create a map out of service binding list
	sbMap := make(map[string][]string)
	for _, v := range sbList {
		if sb, exists := sbMap[v.AppGUID]; exists {
			sb = append(sb, v.ServiceInstanceGUID)
			sbMap[v.AppGUID] = sb
		} else {
			sbMap[v.AppGUID] = []string{v.ServiceInstanceGUID}
		}
	}

	cmd.queryCache = globalQueryCache{
		siMap:    siMap,
		spMap:    spMap,
		sMap:     sMap,
		upsMap:   upsMap,
		spaceMap: spaceMap,
		orgMap:   orgMap,
		sbList:   sbList,
		sbMap:    sbMap,
	}
	return nil
}

//GetMetadata returns metatada
func (cmd *UsageReportCmd) GetMetadata() plugin.PluginMetadata {
	return plugin.PluginMetadata{
		Name: "usage-report",
		Version: plugin.VersionType{
			Major: 1,
			Minor: 6,
			Build: 0,
		},
		Commands: []plugin.Command{
			{
				Name:     "usage-report",
				HelpText: "Report AI and memory usage for orgs and spaces",
				UsageDetails: plugin.Usage{
					Usage: "cf usage-report [-o orgName] [-s spaceName] [-i <app|summary>] [-f <csv>]",
					Options: map[string]string{
						"o": "Filter for Specific Orgranization",
						"s": "Filter for Specific Space",
						"i": "Count Service Instances",
						"f": "Define Output Format (csv)",
					},
				},
			},
		},
	}
}

func (cmd *UsageReportCmd) getFilteredOrgs(orgName, spaceName string) []models.Org {
	var orgs []models.Org

	if orgName != "" {
		org, err := cmd.getOrg(orgName, spaceName)
		if nil != err {
			fmt.Println(err)
			os.Exit(1)
		}
		orgs = append(orgs, org)
	} else {
		var err error
		orgs, err = cmd.getOrgs(spaceName)
		if nil != err {
			fmt.Println(err)
			os.Exit(1)
		}
	}
	return orgs
}

// UsageReportCommand doer
func (cmd *UsageReportCmd) UsageReportCommand(args []string) {
	flagVals := ParseFlags(args)

	var report models.Report

	// make global queries to the API
	if err := cmd.createQueryCache(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var err error
	if report.ServiceInstances, err = CreateServiceInstanceOverview(cmd.queryCache); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	report.Orgs = cmd.getFilteredOrgs(flagVals.OrgName, flagVals.SpaceName)

	// process service instances
	if flagVals.ShowServiceInstances == "app" {
		if flagVals.Format == "csv" {
			fmt.Println(report.ServiceInstanceReportCSV())
		} else {
			fmt.Println(report.ServiceInstanceReportString())
		}
	} else if flagVals.ShowServiceInstances == "summary" {
		if flagVals.Format == "csv" {
			fmt.Println(report.ServiceInstanceSummaryCSV())
		} else {
			fmt.Println(report.ServiceInstanceSummaryString())
		}
	} else {
		if flagVals.Format == "csv" {
			fmt.Println(report.CSV())
		} else {
			fmt.Println(report.String())
		}
	}
}

func (cmd *UsageReportCmd) getOrgs(spaceName string) ([]models.Org, error) {

	rawOrgs, err := cmd.apiHelper.GetOrgs()
	if nil != err {
		return nil, err
	}

	var orgs = []models.Org{}

	for _, o := range rawOrgs {
		orgDetails, err := cmd.getOrgDetails(o, spaceName)
		if err != nil {
			return nil, err
		}
		orgs = append(orgs, orgDetails)
	}
	return orgs, nil
}

func (cmd *UsageReportCmd) getOrg(orgName, spaceName string) (models.Org, error) {
	rawOrg, err := cmd.apiHelper.GetOrg(orgName)
	if nil != err {
		return models.Org{}, err
	}

	return cmd.getOrgDetails(rawOrg, spaceName)
}

func (cmd *UsageReportCmd) getOrgDetails(o apihelper.Organization, spaceName string) (models.Org, error) {
	usage, err := cmd.apiHelper.GetOrgMemoryUsage(o)
	if nil != err {
		return models.Org{}, err
	}
	quota, err := cmd.apiHelper.GetQuotaMemoryLimit(o.QuotaURL)
	if nil != err {
		return models.Org{}, err
	}
	spaces, err := cmd.getSpaces(o.SpacesURL, spaceName)
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

func (cmd *UsageReportCmd) getSpaces(spaceURL, filteredSpaceName string) ([]models.Space, error) {
	rawSpaces, err := cmd.apiHelper.GetOrgSpaces(spaceURL)
	if nil != err {
		return nil, err
	}

	var spaces = []models.Space{}
	for _, s := range rawSpaces {
		// filter spaces
		if filteredSpaceName != "" {
			if s.Name != filteredSpaceName {
				continue
			}
		}

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

		sb := cmd.queryCache.sbMap[a.GUID]
		// sb, err := cmd.apiHelper.GetServiceBindings(a.ServiceBindingsURL)
		// if err != nil {
		//	return nil, err
		// }

		siTotal := len(sb)
		siPCF := 0 // PCF service instances
		siUP := 0  // User Provided Service Instances

		for _, serviceInstanceGUID := range sb {
			if si, exists := cmd.queryCache.siMap[serviceInstanceGUID]; exists {
				if si.Type == "managed_service_instance" {
					if IsPCFInstance(serviceInstanceGUID, cmd.queryCache.siMap, cmd.queryCache.spMap, cmd.queryCache.sMap) {
						siPCF++
					}
				}
			} else if _, exists := cmd.queryCache.upsMap[serviceInstanceGUID]; exists {
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
