package models

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
)

type Org struct {
	Name        string
	MemoryQuota int
	MemoryUsage int
	Spaces      []Space
}

type Space struct {
	Apps []App
	Name string
}

type App struct {
	Ram       int
	Instances int
	Running   bool
	Name      string
	SiTotal   int // Bound Service Instances Total
	SiPCF     int // Bound PCF Service Instances
	SiUP      int // Bound User Provided Service Instances
}

type Report struct {
	Orgs []Org
}

func (org *Org) InstancesCount() int {
	instancesCount := 0
	for _, space := range org.Spaces {
		instancesCount += space.InstancesCount()
	}
	return instancesCount
}

func (org *Org) AppsCount() int {
	appsCount := 0
	for _, space := range org.Spaces {
		appsCount += len(space.Apps)
	}
	return appsCount
}

func (space *Space) ConsumedMemory() int {
	consumed := 0
	for _, app := range space.Apps {
		if app.Running {
			consumed += int(app.Instances * app.Ram)
		}
	}
	return consumed
}

func (space *Space) RunningAppsCount() int {
	runningAppsCount := 0
	for _, app := range space.Apps {
		if app.Running {
			runningAppsCount++
		}
	}
	return runningAppsCount
}

func (space *Space) InstancesCount() int {
	instancesCount := 0
	for _, app := range space.Apps {
		instancesCount += int(app.Instances)
	}
	return instancesCount
}

func (space *Space) RunningInstancesCount() int {
	runningInstancesCount := 0
	for _, app := range space.Apps {
		if app.Running {
			runningInstancesCount += app.Instances
		}
	}
	return runningInstancesCount
}

type ServiceInstance struct {
	GUID string
	Name string
	Type string
}

func (report *Report) ServiceInstanceReportCSV() string {
	var response bytes.Buffer

	response.WriteString(fmt.Sprintf("Org,Space,AppName,Instances,Bound Service Instances,Bound PCF Services,Bound User Provided Services,Bound 3rd Party Services\n"))

	for _, org := range report.Orgs {
		for _, space := range org.Spaces {
			for _, app := range space.Apps {
				thrdParty := app.SiTotal - app.SiPCF - app.SiUP
				record := fmt.Sprintf("%s,%s,%s,%d,%d,%d,%d,%d\n", org.Name, space.Name, app.Name, app.Instances, app.SiTotal, app.SiPCF, app.SiUP, thrdParty)
				response.WriteString(record)
			}
		}
	}

	return response.String()
}

func (report *Report) ServiceInstanceReportString() string {
	var response bytes.Buffer

	for _, org := range report.Orgs {
		response.WriteString(fmt.Sprintf("Org %s\n", org.Name))
		for _, space := range org.Spaces {
			response.WriteString(fmt.Sprintf("\tSpace %s\n", space.Name))
			for _, app := range space.Apps {
				thrdParty := app.SiTotal - app.SiPCF - app.SiUP

				response.WriteString(fmt.Sprintf("\t\tApp %s has %d instances in total.\n", app.Name, app.Instances))
				response.WriteString(fmt.Sprintf("\t\tIt has %d service instances bound in total.\n", app.SiTotal))
				response.WriteString(fmt.Sprintf("\t\tFrom that there are %d PCF service instances, %d user provided service instances,\n", app.SiPCF, app.SiUP))
				response.WriteString(fmt.Sprintf("\t\tand %d 3rd party instances bound.\n\n", thrdParty))
			}
		}
	}

	return response.String()
}

func (report *Report) String() string {
	var response bytes.Buffer

	totalApps := 0
	totalInstances := 0

	for _, org := range report.Orgs {
		response.WriteString(fmt.Sprintf("Org %s is consuming %d MB of %d MB.\n",
			org.Name, org.MemoryUsage, org.MemoryQuota))

		for _, space := range org.Spaces {
			spaceRunningAppsCount := space.RunningAppsCount()
			spaceInstancesCount := space.InstancesCount()
			spaceRunningInstancesCount := space.RunningInstancesCount()
			spaceConsumedMemory := space.ConsumedMemory()

			response.WriteString(
				fmt.Sprintf("\tSpace %s is consuming %d MB memory (%d%%) of org quota.\n",
					space.Name, spaceConsumedMemory, (100 * spaceConsumedMemory / org.MemoryQuota)))
			response.WriteString(
				fmt.Sprintf("\t\t%d apps: %d running %d stopped\n", len(space.Apps),
					spaceRunningAppsCount, len(space.Apps)-spaceRunningAppsCount))
			response.WriteString(
				fmt.Sprintf("\t\t%d instances: %d running, %d stopped\n", spaceInstancesCount,
					spaceRunningInstancesCount, spaceInstancesCount-spaceRunningInstancesCount))
		}

		totalApps += org.AppsCount()
		totalInstances += org.InstancesCount()
	}

	response.WriteString(
		fmt.Sprintf("You are running %d apps in %d org(s), with a total of %d instances.\n",
			totalApps, len(report.Orgs), totalInstances))

	return response.String()
}

func (report *Report) CSV() string {
	var rows = [][]string{}
	var csv bytes.Buffer

	var headers = []string{"OrgName", "SpaceName", "SpaceMemoryUsed", "OrgMemoryQuota", "AppsDeployed", "AppsRunning", "AppInstancesDeployed", "AppInstancesRunning"}

	rows = append(rows, headers)

	for _, org := range report.Orgs {
		for _, space := range org.Spaces {
			appsDeployed := len(space.Apps)

			spaceResult := []string{
				org.Name,
				space.Name,
				strconv.Itoa(space.ConsumedMemory()),
				strconv.Itoa(org.MemoryQuota),
				strconv.Itoa(appsDeployed),
				strconv.Itoa(space.RunningAppsCount()),
				strconv.Itoa(space.InstancesCount()),
				strconv.Itoa(space.RunningInstancesCount()),
			}

			rows = append(rows, spaceResult)
		}
	}

	for i := range rows {
		csv.WriteString(strings.Join(rows[i], ", "))
		csv.WriteString("\n")
	}

	return csv.String()
}
