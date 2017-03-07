package main

import (
	"github.com/dgruber/usagereport-plugin/models"
)

// CreateServiceInstanceOverview creates a list of all services instances available in
// the foundation based on the given cached global REST queries.
func CreateServiceInstanceOverview(cache globalQueryCache) ([]models.Service, error) {
	unique := make(map[string]struct{}, 0)
	var empty struct{}

	r := make([]models.Service, 0, len(cache.siMap))

	for _, si := range cache.siMap {
		var s models.Service
		var orgGUID string

		s.ServiceInstanceGUID = si.GUID
		s.ServiceInstanceName = si.Name
		s.ServiceInstanceType = si.Type

		if _, exists := unique[si.Name]; exists == true {
			continue
		} else {
			unique[si.Name] = empty
		}

		if space, exists := cache.spaceMap[si.SpaceGUID]; exists == true {
			s.SpaceName = space.Name
			orgGUID = space.OrgGUID
		}

		if servicePlan, exists := cache.spMap[si.ServicePlanGUID]; exists == true {
			s.ServicePlanName = servicePlan.Name
			if service, exists := cache.sMap[servicePlan.ServiceGUID]; exists == true {
				s.ServiceName = service.Label
			}
		}

		// TODO type

		if org, exists := cache.orgMap[orgGUID]; exists == true {
			s.OrgName = org.Name
		}

		// find all apps using that service instance
		s.AppGUIDs = make([]string, 0)
		for i, _ := range cache.sbList {
			if cache.sbList[i].ServiceInstanceGUID == s.ServiceInstanceGUID {
				s.AppGUIDs = append(s.AppGUIDs, cache.sbList[i].AppGUID)
			}
		}
		r = append(r, s)
	}

	return r, nil
}
