package main

import (
	"errors"

	"github.com/dgruber/usagereport-plugin/apihelper"
	"github.com/dgruber/usagereport-plugin/apihelper/fakes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Usagereport", func() {
	var fakeAPI *fakes.FakeCFAPIHelper
	var cmd *UsageReportCmd

	BeforeEach(func() {
		fakeAPI = &fakes.FakeCFAPIHelper{}
		cmd = &UsageReportCmd{apiHelper: fakeAPI}
	})

	Describe("get single org errors", func() {
		It("should return an error if cf curl /v2/organizations fails", func() {
			fakeAPI.GetOrgReturns(apihelper.Organization{}, errors.New("Bad Things"))
			_, err := cmd.getOrg("test", "")
			Expect(err).ToNot(BeNil())
		})
	})

	Describe("get orgs errors", func() {

		It("should return an error if cf curl /v2/organizations fails", func() {
			fakeAPI.GetOrgsReturns(nil, errors.New("Bad Things"))
			_, err := cmd.getOrgs("")
			Expect(err).ToNot(BeNil())
		})

		Context("good org bad other thigns", func() {
			BeforeEach(func() {
				fakeAPI.GetOrgsReturns([]apihelper.Organization{apihelper.Organization{}}, nil)
			})

			It("should return an error if cf curl /v2/organizations/{guid}/memory_usage fails", func() {
				fakeAPI.GetOrgMemoryUsageReturns(0, errors.New("Bad Things"))
				_, err := cmd.getOrgs("")
				Expect(err).ToNot(BeNil())
			})

			It("sholud return an error if cf curl to the quota url fails", func() {
				fakeAPI.GetQuotaMemoryLimitReturns(0, errors.New("Bad Things"))
				_, err := cmd.getOrgs("")
				Expect(err).ToNot(BeNil())
			})

			It("should return an error if cf curl to get org spaces fails", func() {
				fakeAPI.GetOrgSpacesReturns(nil, errors.New("Bad Things"))
				_, err := cmd.getOrgs("")
				Expect(err).ToNot(BeNil())
				Expect(fakeAPI.GetOrgSpacesCallCount()).To(Equal(1))
			})

			It("Should return an error if cf curl to get the apps in a space fails", func() {
				fakeAPI.GetOrgSpacesReturns(
					[]apihelper.Space{apihelper.Space{AppsURL: "/v2/apps"}}, nil)
				fakeAPI.GetSpaceAppsReturns(nil, errors.New("Bad Things"))
				_, err := cmd.getOrgs("")
				Expect(err).ToNot(BeNil())
				Expect(fakeAPI.GetSpaceAppsCallCount()).To(Equal(1))
			})
		})

	})

	Describe("Get org composes the values correctly", func() {
		org := apihelper.Organization{
			URL:      "/v2/organizations/1234",
			QuotaURL: "/v2/quotas/2345",
		}

		BeforeEach(func() {
			fakeAPI.GetOrgsReturns([]apihelper.Organization{org}, nil)
		})

		It("should return two one org using 1 mb of 2 mb quota", func() {
			fakeAPI.GetOrgMemoryUsageReturns(float64(1), nil)
			fakeAPI.GetQuotaMemoryLimitReturns(float64(2), nil)
			orgs, err := cmd.getOrgs("")
			Expect(err).To(BeNil())
			Expect(len(orgs)).To(Equal(1))
			org := orgs[0]
			Expect(org.MemoryQuota).To(Equal(2))
			Expect(org.MemoryUsage).To(Equal(1))
		})

		It("Should return an org with 1 space", func() {
			fakeAPI.GetOrgSpacesReturns(
				[]apihelper.Space{apihelper.Space{}, apihelper.Space{}}, nil)
			orgs, _ := cmd.getOrgs("")
			Expect(len(orgs[0].Spaces)).To(Equal(2))
		})

		It("Should not choke on an org with no spaces", func() {
			fakeAPI.GetOrgSpacesReturns(
				[]apihelper.Space{}, nil)
			orgs, _ := cmd.getOrgs("")
			Expect(len(orgs[0].Spaces)).To(Equal(0))
		})

		It("Should return two apps from a space", func() {
			fakeAPI.GetOrgSpacesReturns(
				[]apihelper.Space{apihelper.Space{}}, nil)

			fakeAPI.GetSpaceAppsReturns(
				[]apihelper.App{
					apihelper.App{},
					apihelper.App{},
					apihelper.App{},
				},
				nil)
			orgs, _ := cmd.getOrgs("")
			org := orgs[0]
			space := org.Spaces[0]
			apps := space.Apps
			Expect(len(apps)).To(Equal(3))
		})

		It("Should mark the first app as running, the second as stopped", func() {
			fakeAPI.GetOrgSpacesReturns(
				[]apihelper.Space{apihelper.Space{}}, nil)

			fakeAPI.GetSpaceAppsReturns(
				[]apihelper.App{
					apihelper.App{Running: true},
					apihelper.App{Running: false},
				},
				nil)

			orgs, _ := cmd.getOrgs("")
			org := orgs[0]
			space := org.Spaces[0]
			apps := space.Apps
			Expect(apps[0].Running).To(BeTrue())
			Expect(apps[1].Running).To(BeFalse())
		})
	})

	Describe("PCF service type discovery", func() {
		var siMap map[string]apihelper.ServiceInstance
		var spMap map[string]apihelper.ServicePlan
		var sMap map[string]apihelper.Service

		BeforeEach(func() {
			siMap = make(map[string]apihelper.ServiceInstance)
			spMap = make(map[string]apihelper.ServicePlan)
			sMap = make(map[string]apihelper.Service)

			// this is a PCF service
			siMap["123"] = apihelper.ServiceInstance{ServicePlanGUID: "234"}
			spMap["234"] = apihelper.ServicePlan{ServiceGUID: "345"}
			sMap["345"] = apihelper.Service{GUID: "345", Label: "p-mysql"}

			// this is not a PCF service
			siMap["!123"] = apihelper.ServiceInstance{ServicePlanGUID: "!234"}
			spMap["!234"] = apihelper.ServicePlan{ServiceGUID: "!345"}
			sMap["!345"] = apihelper.Service{GUID: "!345", Label: "mysql"}
		})

		It("Should return true if a PCF service is discovered", func() {
			Expect(IsPCFInstance("123", siMap, spMap, sMap)).To(BeTrue())
		})

		It("Should return false if a PCF service is discovered", func() {
			Expect(IsPCFInstance("!123", siMap, spMap, sMap)).To(BeFalse())
		})
	})

	Describe("service instance overview generation", func() {
		var cache globalQueryCache

		BeforeEach(func() {
			cache.siMap = make(map[string]apihelper.ServiceInstance)
			cache.spMap = make(map[string]apihelper.ServicePlan)
			cache.sMap = make(map[string]apihelper.Service)
			cache.upsMap = make(map[string]apihelper.UserProvidedService)
			cache.spaceMap = make(map[string]apihelper.SpaceDetails)
			cache.sbList = make([]apihelper.ServiceBinding, 0)

			cache.siMap["serviceInstanceKey"] = apihelper.ServiceInstance{
				GUID:            "serviceInstanceGUID",
				Name:            "myserviceinstance",
				Type:            "my_type",
				ServicePlanGUID: "servicePlanGUID",
				SpaceGUID:       "spaceGUID",
			}

			cache.spMap["servicePlanGUID"] = apihelper.ServicePlan{
				GUID:        "servicePlanGUID",
				Name:        "ServicePlanName",
				ServiceGUID: "ServiceGUID",
			}

			cache.sMap["ServiceGUID"] = apihelper.Service{
				GUID:  "ServiceGUID",
				Label: "p-service",
			}

			cache.upsMap["userProvidedServiceGUID"] = apihelper.UserProvidedService{
				GUID: "userProvidedServiceGUID",
				Name: "UserProvidedService",
				Type: "my_user_provided_service_type",
			}

			cache.spaceMap["spaceGUID"] = apihelper.SpaceDetails{
				GUID: "spaceGUID",
				Name: "SpaceName",
			}

			cache.sbList = append(cache.sbList, apihelper.ServiceBinding{
				AppGUID:             "AppGUID",
				ServiceInstanceGUID: "serviceInstanceGUID",
			})
		})

		It("should build up the service instance description using the cache", func() {
			services, err := CreateServiceInstanceOverview(cache)
			Expect(err).To(BeNil())
			Expect(services).NotTo(BeNil())
			Expect(len(services)).To(Equal(1))
			Expect(services[0].ServiceName).To(Equal("p-service"))
			Expect(services[0].SpaceName).To(Equal("SpaceName"))
			Expect(services[0].ServicePlanName).To(Equal("ServicePlanName"))
			Expect(services[0].ServiceInstanceName).To(Equal("myserviceinstance"))
			Expect(services[0].AppGUIDs).NotTo(BeNil())
			Expect(services[0].AppGUIDs[0]).To(Equal("AppGUID"))
		})

	})

})
