// This file was generated by counterfeiter
package fakes

import (
	"sync"

	// "github.com/cloudfoundry/cli/plugin"
	"github.com/dgruber/usagereport-plugin/apihelper"
)

type FakeCFAPIHelper struct {
	GetOrgsStub    func() ([]apihelper.Organization, error)
	getOrgsMutex   sync.RWMutex
	getOrgsReturns struct {
		result1 []apihelper.Organization
		result2 error
	}
	GetOrgStub    func(string) (apihelper.Organization, error)
	getOrgMutex   sync.RWMutex
	getOrgReturns struct {
		result1 apihelper.Organization
		result2 error
	}
	GetQuotaMemoryLimitStub        func(string) (float64, error)
	getQuotaMemoryLimitMutex       sync.RWMutex
	getQuotaMemoryLimitArgsForCall []struct {
		arg1 string
	}
	getQuotaMemoryLimitReturns struct {
		result1 float64
		result2 error
	}
	GetOrgMemoryUsageStub        func(apihelper.Organization) (float64, error)
	getOrgMemoryUsageMutex       sync.RWMutex
	getOrgMemoryUsageArgsForCall []struct {
		arg1 apihelper.Organization
	}
	getOrgMemoryUsageReturns struct {
		result1 float64
		result2 error
	}
	GetOrgSpacesStub        func(string) ([]apihelper.Space, error)
	getOrgSpacesMutex       sync.RWMutex
	getOrgSpacesArgsForCall []struct {
		arg1 string
	}
	getOrgSpacesReturns struct {
		result1 []apihelper.Space
		result2 error
	}
	GetSpaceAppsStub        func(string) ([]apihelper.App, error)
	getSpaceAppsMutex       sync.RWMutex
	getSpaceAppsArgsForCall []struct {
		arg1 string
	}
	getSpaceAppsReturns struct {
		result1 []apihelper.App
		result2 error
	}
	GetServiceBindingsStub    func() ([]apihelper.ServiceBindings, error)
	getServiceBindingsMutex   sync.RWMutex
	getServiceBindingsReturns struct {
		result1 []apihelper.ServiceBindings
		result2 error
	}
	GetServiceInstanceMapStub    func() (map[string]apihelper.ServiceInstance, error)
	getServiceInstanceMapMutex   sync.RWMutex
	getServiceInstanceMapReturns struct {
		result1 map[string]apihelper.ServiceInstance
		result2 error
	}
	GetServiceMapStub    func() (map[string]apihelper.Service, error)
	getServiceMapMutex   sync.RWMutex
	getServiceMapReturns struct {
		result1 map[string]apihelper.Service
		result2 error
	}
	GetServicePlanMapStub    func() (map[string]apihelper.ServicePlan, error)
	getServicePlanMapMutex   sync.RWMutex
	getServicePlanMapReturns struct {
		result1 map[string]apihelper.ServicePlan
		result2 error
	}
	GetUserProvidedServiceMapStub    func() (map[string]apihelper.UserProvidedService, error)
	getUserProvidedServiceMapMutex   sync.RWMutex
	getUserProvidedServiceMapReturns struct {
		result1 map[string]apihelper.UserProvidedService
		result2 error
	}

	GetServiceBindingsListStub    func() ([]apihelper.ServiceBinding, error)
	getServiceBindingsListMutex   sync.RWMutex
	getServiceBindingsListReturns struct {
		result1 []apihelper.ServiceBinding
		result2 error
	}

	GetSpaceMapStub    func() (map[string]apihelper.SpaceDetails, error)
	getSpaceMapMutex   sync.RWMutex
	getSpaceMapReturns struct {
		result1 map[string]apihelper.SpaceDetails
		result2 error
	}

	GetOrgMapStub    func() (map[string]apihelper.OrgDetails, error)
	getOrgMapMutex   sync.RWMutex
	getOrgMapReturns struct {
		result1 map[string]apihelper.OrgDetails
		result2 error
	}
}

func (fake *FakeCFAPIHelper) GetOrgs() ([]apihelper.Organization, error) {
	fake.getOrgsMutex.Lock()
	fake.getOrgsMutex.Unlock()
	if fake.GetOrgsStub != nil {
		return fake.GetOrgsStub()
	} else {
		return fake.getOrgsReturns.result1, fake.getOrgsReturns.result2
	}
}

func (fake *FakeCFAPIHelper) GetOrgsReturns(result1 []apihelper.Organization, result2 error) {
	fake.GetOrgsStub = nil
	fake.getOrgsReturns = struct {
		result1 []apihelper.Organization
		result2 error
	}{result1, result2}
}

func (fake *FakeCFAPIHelper) GetOrg(name string) (apihelper.Organization, error) {
	fake.getOrgMutex.Lock()
	fake.getOrgMutex.Unlock()
	if fake.GetOrgStub != nil {
		return fake.GetOrgStub(name)
	} else {
		return fake.getOrgReturns.result1, fake.getOrgReturns.result2
	}
}

func (fake *FakeCFAPIHelper) GetOrgReturns(result1 apihelper.Organization, result2 error) {
	fake.GetOrgStub = nil
	fake.getOrgReturns = struct {
		result1 apihelper.Organization
		result2 error
	}{result1, result2}
}

func (fake *FakeCFAPIHelper) GetQuotaMemoryLimit(arg1 string) (float64, error) {
	fake.getQuotaMemoryLimitMutex.Lock()
	fake.getQuotaMemoryLimitArgsForCall = append(fake.getQuotaMemoryLimitArgsForCall, struct {
		arg1 string
	}{arg1})
	fake.getQuotaMemoryLimitMutex.Unlock()
	if fake.GetQuotaMemoryLimitStub != nil {
		return fake.GetQuotaMemoryLimitStub(arg1)
	} else {
		return fake.getQuotaMemoryLimitReturns.result1, fake.getQuotaMemoryLimitReturns.result2
	}
}

func (fake *FakeCFAPIHelper) GetQuotaMemoryLimitCallCount() int {
	fake.getQuotaMemoryLimitMutex.RLock()
	defer fake.getQuotaMemoryLimitMutex.RUnlock()
	return len(fake.getQuotaMemoryLimitArgsForCall)
}

func (fake *FakeCFAPIHelper) GetQuotaMemoryLimitArgsForCall(i int) string {
	fake.getQuotaMemoryLimitMutex.RLock()
	defer fake.getQuotaMemoryLimitMutex.RUnlock()
	return fake.getQuotaMemoryLimitArgsForCall[i].arg1
}

func (fake *FakeCFAPIHelper) GetQuotaMemoryLimitReturns(result1 float64, result2 error) {
	fake.GetQuotaMemoryLimitStub = nil
	fake.getQuotaMemoryLimitReturns = struct {
		result1 float64
		result2 error
	}{result1, result2}
}

func (fake *FakeCFAPIHelper) GetOrgMemoryUsage(arg1 apihelper.Organization) (float64, error) {
	fake.getOrgMemoryUsageMutex.Lock()
	fake.getOrgMemoryUsageArgsForCall = append(fake.getOrgMemoryUsageArgsForCall, struct {
		arg1 apihelper.Organization
	}{arg1})
	fake.getOrgMemoryUsageMutex.Unlock()
	if fake.GetOrgMemoryUsageStub != nil {
		return fake.GetOrgMemoryUsageStub(arg1)
	} else {
		return fake.getOrgMemoryUsageReturns.result1, fake.getOrgMemoryUsageReturns.result2
	}
}

func (fake *FakeCFAPIHelper) GetOrgMemoryUsageCallCount() int {
	fake.getOrgMemoryUsageMutex.RLock()
	defer fake.getOrgMemoryUsageMutex.RUnlock()
	return len(fake.getOrgMemoryUsageArgsForCall)
}

func (fake *FakeCFAPIHelper) GetOrgMemoryUsageArgsForCall(i int) apihelper.Organization {
	fake.getOrgMemoryUsageMutex.RLock()
	defer fake.getOrgMemoryUsageMutex.RUnlock()
	return fake.getOrgMemoryUsageArgsForCall[i].arg1
}

func (fake *FakeCFAPIHelper) GetOrgMemoryUsageReturns(result1 float64, result2 error) {
	fake.GetOrgMemoryUsageStub = nil
	fake.getOrgMemoryUsageReturns = struct {
		result1 float64
		result2 error
	}{result1, result2}
}

func (fake *FakeCFAPIHelper) GetOrgSpaces(arg1 string) ([]apihelper.Space, error) {
	fake.getOrgSpacesMutex.Lock()
	fake.getOrgSpacesArgsForCall = append(fake.getOrgSpacesArgsForCall, struct {
		arg1 string
	}{arg1})
	fake.getOrgSpacesMutex.Unlock()
	if fake.GetOrgSpacesStub != nil {
		return fake.GetOrgSpacesStub(arg1)
	} else {
		return fake.getOrgSpacesReturns.result1, fake.getOrgSpacesReturns.result2
	}
}

func (fake *FakeCFAPIHelper) GetOrgSpacesCallCount() int {
	fake.getOrgSpacesMutex.RLock()
	defer fake.getOrgSpacesMutex.RUnlock()
	return len(fake.getOrgSpacesArgsForCall)
}

func (fake *FakeCFAPIHelper) GetOrgSpacesArgsForCall(i int) string {
	fake.getOrgSpacesMutex.RLock()
	defer fake.getOrgSpacesMutex.RUnlock()
	return fake.getOrgSpacesArgsForCall[i].arg1
}

func (fake *FakeCFAPIHelper) GetOrgSpacesReturns(result1 []apihelper.Space, result2 error) {
	fake.GetOrgSpacesStub = nil
	fake.getOrgSpacesReturns = struct {
		result1 []apihelper.Space
		result2 error
	}{result1, result2}
}

func (fake *FakeCFAPIHelper) GetSpaceApps(arg1 string) ([]apihelper.App, error) {
	fake.getSpaceAppsMutex.Lock()
	fake.getSpaceAppsArgsForCall = append(fake.getSpaceAppsArgsForCall, struct {
		arg1 string
	}{arg1})
	fake.getSpaceAppsMutex.Unlock()
	if fake.GetSpaceAppsStub != nil {
		return fake.GetSpaceAppsStub(arg1)
	} else {
		return fake.getSpaceAppsReturns.result1, fake.getSpaceAppsReturns.result2
	}
}

func (fake *FakeCFAPIHelper) GetSpaceAppsCallCount() int {
	fake.getSpaceAppsMutex.RLock()
	defer fake.getSpaceAppsMutex.RUnlock()
	return len(fake.getSpaceAppsArgsForCall)
}

func (fake *FakeCFAPIHelper) GetSpaceAppsArgsForCall(i int) string {
	fake.getSpaceAppsMutex.RLock()
	defer fake.getSpaceAppsMutex.RUnlock()
	return fake.getSpaceAppsArgsForCall[i].arg1
}

func (fake *FakeCFAPIHelper) GetSpaceAppsReturns(result1 []apihelper.App, result2 error) {
	fake.GetSpaceAppsStub = nil
	fake.getSpaceAppsReturns = struct {
		result1 []apihelper.App
		result2 error
	}{result1, result2}
}

func (fake *FakeCFAPIHelper) GetServiceBindings(url string) ([]apihelper.ServiceBindings, error) {
	fake.getServiceBindingsMutex.RLock()
	defer fake.getServiceBindingsMutex.RUnlock()
	return fake.getServiceBindingsReturns.result1, fake.getServiceBindingsReturns.result2
}

func (fake *FakeCFAPIHelper) GetServiceInstanceMap() (map[string]apihelper.ServiceInstance, error) {
	fake.getServiceInstanceMapMutex.RLock()
	defer fake.getServiceInstanceMapMutex.RUnlock()
	return fake.getServiceInstanceMapReturns.result1, fake.getServiceInstanceMapReturns.result2
}

func (fake *FakeCFAPIHelper) GetServiceMap() (map[string]apihelper.Service, error) {
	fake.getServiceMapMutex.RLock()
	defer fake.getServiceMapMutex.RUnlock()
	return fake.getServiceMapReturns.result1, fake.getServiceMapReturns.result2
}

func (fake *FakeCFAPIHelper) GetServicePlanMap() (map[string]apihelper.ServicePlan, error) {
	fake.getServicePlanMapMutex.RLock()
	defer fake.getServicePlanMapMutex.RUnlock()
	return fake.getServicePlanMapReturns.result1, fake.getServicePlanMapReturns.result2
}

func (fake *FakeCFAPIHelper) GetUserProvidedServiceMap() (map[string]apihelper.UserProvidedService, error) {
	fake.getUserProvidedServiceMapMutex.RLock()
	defer fake.getUserProvidedServiceMapMutex.RUnlock()
	return fake.getUserProvidedServiceMapReturns.result1, fake.getUserProvidedServiceMapReturns.result2
}

func (fake *FakeCFAPIHelper) GetServiceBindingsList() ([]apihelper.ServiceBinding, error) {
	fake.getServiceBindingsListMutex.RLock()
	defer fake.getServiceBindingsListMutex.RUnlock()
	return fake.getServiceBindingsListReturns.result1, fake.getServiceBindingsListReturns.result2
}

func (fake *FakeCFAPIHelper) GetSpaceMap() (map[string]apihelper.SpaceDetails, error) {
	fake.getSpaceMapMutex.RLock()
	defer fake.getSpaceMapMutex.RUnlock()
	return fake.getSpaceMapReturns.result1, fake.getSpaceMapReturns.result2
}

func (fake *FakeCFAPIHelper) GetOrgMap() (map[string]apihelper.OrgDetails, error) {
	fake.getOrgMapMutex.RLock()
	defer fake.getOrgMapMutex.RUnlock()
	return fake.getOrgMapReturns.result1, fake.getOrgMapReturns.result2
}

var _ apihelper.CFAPIHelper = new(FakeCFAPIHelper)
