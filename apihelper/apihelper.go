package apihelper

import (
	"errors"
	"fmt"
	"net/url"
	"strconv"

	"github.com/cloudfoundry/cli/plugin"
	"github.com/krujos/cfcurl"
)

var (
	ErrOrgNotFound = errors.New("organization not found")
)

//Organization representation
type Organization struct {
	URL       string
	Name      string
	QuotaURL  string
	SpacesURL string
}

//Space representation
type Space struct {
	Name    string
	AppsURL string
}

//App representation
type App struct {
	Instances          float64
	RAM                float64
	Running            bool
	Name               string
	ServiceBindingsURL string
}

//CFAPIHelper to wrap cf curl results
type CFAPIHelper interface {
	GetOrgs() ([]Organization, error)
	GetOrg(string) (Organization, error)
	GetQuotaMemoryLimit(string) (float64, error)
	GetOrgMemoryUsage(Organization) (float64, error)
	GetOrgSpaces(string) ([]Space, error)
	GetSpaceApps(string) ([]App, error)
	GetServiceBindings(serviceBindingsURL string) ([]ServiceBindings, error)
	GetServiceInstanceMap(siURL string) (map[string]ServiceInstance, error)
	GetServiceMap() (map[string]Service, error)
	GetServicePlanMap() (map[string]ServicePlan, error)
	GetUserProvidedServiceMap() (map[string]UserProvidedService, error)
}

//APIHelper implementation
type APIHelper struct {
	cli plugin.CliConnection
}

func New(cli plugin.CliConnection) CFAPIHelper {
	return &APIHelper{cli}
}

//GetOrgs returns a struct that represents critical fields in the JSON
func (api *APIHelper) GetOrgs() ([]Organization, error) {
	orgsJSON, err := cfcurl.Curl(api.cli, "/v2/organizations")
	if nil != err {
		return nil, err
	}
	pages := int(orgsJSON["total_pages"].(float64))
	orgs := []Organization{}
	for i := 1; i <= pages; i++ {
		if 1 != i {
			orgsJSON, err = cfcurl.Curl(api.cli, "/v2/organizations?page="+strconv.Itoa(i))
		}
		for _, o := range orgsJSON["resources"].([]interface{}) {
			theOrg := o.(map[string]interface{})
			entity := theOrg["entity"].(map[string]interface{})
			metadata := theOrg["metadata"].(map[string]interface{})
			orgs = append(orgs,
				Organization{
					Name:      entity["name"].(string),
					URL:       metadata["url"].(string),
					QuotaURL:  entity["quota_definition_url"].(string),
					SpacesURL: entity["spaces_url"].(string),
				})
		}
	}
	return orgs, nil
}

//GetOrg returns a struct that represents critical fields in the JSON
func (api *APIHelper) GetOrg(name string) (Organization, error) {
	query := fmt.Sprintf("name:%s", name)
	path := fmt.Sprintf("/v2/organizations?q=%s&inline-relations-depth=1", url.QueryEscape(query))
	orgsJSON, err := cfcurl.Curl(api.cli, path)
	if nil != err {
		return Organization{}, err
	}

	results := int(orgsJSON["total_results"].(float64))
	if results == 0 {
		return Organization{}, ErrOrgNotFound
	}

	orgResource := orgsJSON["resources"].([]interface{})[0]
	org := api.orgResourceToOrg(orgResource)

	return org, nil
}

func (api *APIHelper) orgResourceToOrg(o interface{}) Organization {
	theOrg := o.(map[string]interface{})
	entity := theOrg["entity"].(map[string]interface{})
	metadata := theOrg["metadata"].(map[string]interface{})
	return Organization{
		Name:      entity["name"].(string),
		URL:       metadata["url"].(string),
		QuotaURL:  entity["quota_definition_url"].(string),
		SpacesURL: entity["spaces_url"].(string),
	}
}

//GetQuotaMemoryLimit retruns the amount of memory (in MB) that the org is allowed
func (api *APIHelper) GetQuotaMemoryLimit(quotaURL string) (float64, error) {
	quotaJSON, err := cfcurl.Curl(api.cli, quotaURL)
	if nil != err {
		return 0, err
	}
	return quotaJSON["entity"].(map[string]interface{})["memory_limit"].(float64), nil
}

//GetOrgMemoryUsage returns the amount of memory (in MB) that the org is consuming
func (api *APIHelper) GetOrgMemoryUsage(org Organization) (float64, error) {
	usageJSON, err := cfcurl.Curl(api.cli, org.URL+"/memory_usage")
	if nil != err {
		return 0, err
	}
	return usageJSON["memory_usage_in_mb"].(float64), nil
}

//GetOrgSpaces returns the spaces in an org.
func (api *APIHelper) GetOrgSpaces(spacesURL string) ([]Space, error) {
	spacesJSON, err := cfcurl.Curl(api.cli, spacesURL)
	if nil != err {
		return nil, err
	}
	spaces := []Space{}
	for _, s := range spacesJSON["resources"].([]interface{}) {
		theSpace := s.(map[string]interface{})
		entity := theSpace["entity"].(map[string]interface{})
		spaces = append(spaces,
			Space{
				AppsURL: entity["apps_url"].(string),
				Name:    entity["name"].(string),
			})
	}
	return spaces, nil
}

//GetSpaceApps returns the apps in a space
func (api *APIHelper) GetSpaceApps(appsURL string) ([]App, error) {
	appsJSON, err := cfcurl.Curl(api.cli, appsURL)
	if nil != err {
		return nil, err
	}
	apps := []App{}
	for _, a := range appsJSON["resources"].([]interface{}) {
		theApp := a.(map[string]interface{})
		entity := theApp["entity"].(map[string]interface{})
		apps = append(apps,
			App{
				Instances:          entity["instances"].(float64),
				RAM:                entity["memory"].(float64),
				Running:            "STARTED" == entity["state"].(string),
				ServiceBindingsURL: entity["service_bindings_url"].(string),
				Name:               entity["name"].(string),
			})
	}
	return apps, nil
}

type ServiceBindings struct {
	ServiceInstanceGUID string
}

func (api *APIHelper) GetServiceBindings(serviceBindingsURL string) ([]ServiceBindings, error) {
	appsJSON, err := cfcurl.Curl(api.cli, serviceBindingsURL)
	if nil != err {
		return nil, err
	}
	sbs := []ServiceBindings{}
	for _, a := range appsJSON["resources"].([]interface{}) {
		theSvc := a.(map[string]interface{})
		entity := theSvc["entity"].(map[string]interface{})
		sbs = append(sbs,
			ServiceBindings{
				ServiceInstanceGUID: entity["service_instance_guid"].(string),
			})
	}
	return sbs, nil
}

// ------

type ServiceInstance struct {
	GUID            string
	Name            string
	Type            string
	ServicePlanGUID string
}

// GetServiceInstanceMap returns a map from Service Instance GUID to a Service Instance.
func (api *APIHelper) GetServiceInstanceMap(siURL string) (map[string]ServiceInstance, error) {
	siJSON, err := cfcurl.Curl(api.cli, "/v2/service_instances")
	if nil != err {
		return nil, err
	}

	simap := make(map[string]ServiceInstance, 32)

	for _, a := range siJSON["resources"].([]interface{}) {
		theSvc := a.(map[string]interface{})

		meta := theSvc["metadata"].(map[string]interface{})
		entity := theSvc["entity"].(map[string]interface{})

		simap[meta["guid"].(string)] = ServiceInstance{
			GUID:            meta["guid"].(string),
			Name:            entity["name"].(string),
			Type:            entity["type"].(string),
			ServicePlanGUID: entity["service_plan_guid"].(string),
		}
	}
	return simap, nil
}

type ServicePlan struct {
	GUID        string // ServicePlan GUID
	ServiceGUID string
}

// GetServicePlanMap maps a ServicePlan GUID to a Service GUID.
func (api *APIHelper) GetServicePlanMap() (map[string]ServicePlan, error) {
	spJSON, err := cfcurl.Curl(api.cli, "/v2/service_plans")
	if nil != err {
		return nil, err
	}

	spMap := make(map[string]ServicePlan, 32)

	for _, a := range spJSON["resources"].([]interface{}) {
		theSvc := a.(map[string]interface{})

		meta := theSvc["metadata"].(map[string]interface{})
		entity := theSvc["entity"].(map[string]interface{})

		spMap[meta["guid"].(string)] = ServicePlan{
			GUID:        meta["guid"].(string),
			ServiceGUID: entity["service_guid"].(string),
		}
	}
	return spMap, nil
}

type Service struct {
	GUID  string // Service GUID
	Label string // name of the service (starts with p- in case it is a Pivotal service)
}

// GetServiceMap maps a Service GUID to a Service Name (label).
func (api *APIHelper) GetServiceMap() (map[string]Service, error) {
	siJSON, err := cfcurl.Curl(api.cli, "/v2/services")
	if nil != err {
		return nil, err
	}

	simap := make(map[string]Service, 32)

	for _, a := range siJSON["resources"].([]interface{}) {
		theSvc := a.(map[string]interface{})

		meta := theSvc["metadata"].(map[string]interface{})
		entity := theSvc["entity"].(map[string]interface{})

		simap[meta["guid"].(string)] = Service{
			GUID:  meta["guid"].(string),
			Label: entity["label"].(string),
		}
	}
	return simap, nil
}

type UserProvidedService struct {
	GUID string
	Type string
}

func (api *APIHelper) GetUserProvidedServiceMap() (map[string]UserProvidedService, error) {
	siJSON, err := cfcurl.Curl(api.cli, "/v2/user_provided_service_instances")
	if nil != err {
		return nil, err
	}

	simap := make(map[string]UserProvidedService, 32)

	for _, a := range siJSON["resources"].([]interface{}) {
		theSvc := a.(map[string]interface{})

		meta := theSvc["metadata"].(map[string]interface{})
		entity := theSvc["entity"].(map[string]interface{})

		simap[meta["guid"].(string)] = UserProvidedService{
			GUID: meta["guid"].(string),
			Type: entity["type"].(string),
		}
	}
	return simap, nil
}
