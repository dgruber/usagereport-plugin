package apihelper

import (
	"bufio"
	"errors"
	"os"

	"github.com/cloudfoundry/cli/plugin/pluginfakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func slurp(filename string) []string {
	var b []string
	file, _ := os.Open(filename)
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		b = append(b, scanner.Text())
	}
	return b
}

var _ = Describe("UsageReport", func() {
	var api CFAPIHelper
	var fakeCliConnection *pluginfakes.FakeCliConnection

	BeforeEach(func() {
		fakeCliConnection = &pluginfakes.FakeCliConnection{}
		api = New(fakeCliConnection)
	})

	Describe("Get orgs", func() {
		var orgsJSON []string

		BeforeEach(func() {
			orgsJSON = slurp("test-data/orgs.json")
		})

		It("should return two orgs", func() {
			fakeCliConnection.CliCommandWithoutTerminalOutputReturns(orgsJSON, nil)
			orgs, _ := api.GetOrgs()
			Expect(len(orgs)).To(Equal(2))
		})

		It("does something intellegent when cf curl fails", func() {
			fakeCliConnection.CliCommandWithoutTerminalOutputReturns(
				nil, errors.New("bad things"))
			_, err := api.GetOrgs()
			Expect(err).ToNot(BeNil())
		})

		It("populates the url", func() {
			fakeCliConnection.CliCommandWithoutTerminalOutputReturns(orgsJSON, nil)
			orgs, _ := api.GetOrgs()
			org := orgs[0]
			Expect(org.URL).To(Equal("/v2/organizations/b1a23fd6-ac8d-4304-a3b4-815745417acd"))
		})

		It("calls /v2/orgs", func() {
			fakeCliConnection.CliCommandWithoutTerminalOutputReturns(orgsJSON, nil)
			api.GetOrgs()
			args := fakeCliConnection.CliCommandWithoutTerminalOutputArgsForCall(0)
			Expect(args[1]).To(Equal("/v2/organizations"))
		})

	})

	Describe("paged org output", func() {
		var orgsPage1 []string

		BeforeEach(func() {
			orgsPage1 = slurp("test-data/paged-orgs-page-1.json")
		})

		It("deals with paged output", func() {
			fakeCliConnection.CliCommandWithoutTerminalOutputReturns(orgsPage1, nil)
			api.GetOrgs()
			args := fakeCliConnection.CliCommandWithoutTerminalOutputArgsForCall(0)
			Expect(args[1]).To(Equal("/v2/organizations"))
			Ω(fakeCliConnection.CliCommandWithoutTerminalOutputCallCount()).To(Equal(2))
		})

		It("Should have 100 orgs", func() {
			fakeCliConnection.CliCommandWithoutTerminalOutputReturns(orgsPage1, nil)
			orgs, _ := api.GetOrgs()
			args := fakeCliConnection.CliCommandWithoutTerminalOutputArgsForCall(1)
			Expect(args[1]).To(Equal("/v2/organizations?page=2"))
			Ω(orgs).To(HaveLen(100))
		})
	})

	Describe("Get quota memory limit", func() {
		var quotaJSON []string

		BeforeEach(func() {
			quotaJSON = slurp("test-data/quota.json")
		})

		It("should return an error when it can't fetch the memory limit", func() {
			_, err := api.GetQuotaMemoryLimit("/v2/somequota")
			fakeCliConnection.CliCommandWithoutTerminalOutputReturns(
				nil, errors.New("Bad Things"))
			Expect(err).ToNot(BeNil())
		})

		It("should reutrn 10240 as the memory limit", func() {
			fakeCliConnection.CliCommandWithoutTerminalOutputReturns(
				quotaJSON, nil)
			limit, _ := api.GetQuotaMemoryLimit("/v2/quotas/")
			Expect(limit).To(Equal(float64(10240)))
		})
	})

	Describe("it Gets the org memory usage", func() {
		var org Organization
		var usageJSON []string

		BeforeEach(func() {
			usageJSON = slurp("test-data/memory_usage.json")
		})

		It("should return an error when it can't fetch the orgs memory usage", func() {
			fakeCliConnection.CliCommandWithoutTerminalOutputReturns(nil,
				errors.New("Bad things"))
			_, err := api.GetOrgMemoryUsage(org)
			Expect(err).ToNot(BeNil())
		})

		It("should return the memory usage", func() {
			org.URL = "/v2/organizations/1234"
			fakeCliConnection.CliCommandWithoutTerminalOutputReturns(usageJSON, nil)
			usage, _ := api.GetOrgMemoryUsage(org)
			Expect(usage).To(Equal(float64(512)))
		})
	})

	Describe("get spaces", func() {
		var spacesJSON []string

		BeforeEach(func() {
			spacesJSON = slurp("test-data/spaces.json")
		})

		It("should error when the the spaces url fails", func() {
			fakeCliConnection.CliCommandWithoutTerminalOutputReturns(nil, errors.New("Bad Things"))
			_, err := api.GetOrgSpaces("/v2/organizations/12345/spaces")
			Expect(err).ToNot(BeNil())
		})

		It("should return two spaces", func() {
			fakeCliConnection.CliCommandWithoutTerminalOutputReturns(spacesJSON, nil)
			spaces, _ := api.GetOrgSpaces("/v2/organizations/12345/spaces")
			Expect(len(spaces)).To(Equal(2))
		})

		It("should have name jdk-space", func() {
			fakeCliConnection.CliCommandWithoutTerminalOutputReturns(spacesJSON, nil)
			spaces, _ := api.GetOrgSpaces("/v2/organizations/12345/spaces")
			Expect(spaces[0].Name).To(Equal("jdk-space"))
			Expect(spaces[0].AppsURL).To(Equal("/v2/spaces/81c310ed-d258-48d7-a57a-6522d93a4217/apps"))
		})
	})

	Describe("get apps", func() {
		var appsJSON []string

		BeforeEach(func() {
			appsJSON = slurp("test-data/apps.json")
		})

		It("should return an error when the apps url fails", func() {
			fakeCliConnection.CliCommandWithoutTerminalOutputReturns(nil, errors.New("Bad Things"))
			_, err := api.GetSpaceApps("/v2/whateverapps")
			Expect(err).ToNot(BeNil())
		})

		It("should return one app with 1 instance and 1024 mb of ram", func() {
			fakeCliConnection.CliCommandWithoutTerminalOutputReturns(appsJSON, nil)
			apps, _ := api.GetSpaceApps("/v2/whateverapps")
			Expect(len(apps)).To(Equal(1))
			Expect(apps[0].Instances).To(Equal(float64(1)))
			Expect(apps[0].RAM).To(Equal(float64(1024)))
			Expect(apps[0].Running).To(BeTrue())
		})
	})

	// TODO need tests for no spaces and no apps in org.

	Describe("get service bindings", func() {
		var sbJSON []string

		BeforeEach(func() {
			sbJSON = slurp("test-data/service_bindings_for_app.json")
		})

		It("should return an error when the service binding url fails", func() {
			fakeCliConnection.CliCommandWithoutTerminalOutputReturns(nil, errors.New("Bad Things"))
			_, err := api.GetServiceBindings("/v2/whateverapps")
			Expect(err).ToNot(BeNil())
		})

		It("should return one service binding with the service instance GUID to be set", func() {
			fakeCliConnection.CliCommandWithoutTerminalOutputReturns(sbJSON, nil)
			bindings, err := api.GetServiceBindings("/v2/whateverapps")
			Expect(err).To(BeNil())
			Expect(len(bindings)).To(Equal(1))
			Expect(bindings[0].ServiceInstanceGUID).To(Equal("92f0f510-dbb1-4c04-aa7c-28a8dc0797b4"))
		})

	})

	Describe("get service instance map", func() {
		var serviceInstancesJSON []string

		BeforeEach(func() {
			serviceInstancesJSON = slurp("test-data/service_instances.json")
		})

		It("should return an error when the service instance url fails", func() {
			fakeCliConnection.CliCommandWithoutTerminalOutputReturns(nil, errors.New("Bad Things"))
			_, err := api.GetServiceInstanceMap()
			Expect(err).ToNot(BeNil())
		})

		It("should return a map containing a specific element with all entries set", func() {
			fakeCliConnection.CliCommandWithoutTerminalOutputReturns(serviceInstancesJSON, nil)
			siMap, err := api.GetServiceInstanceMap()

			Expect(err).To(BeNil())
			Expect(siMap).NotTo(BeNil())

			si, exists := siMap["215b97be-ec77-4224-9c38-c4f2d86b56c1"]
			Expect(exists).To(BeTrue())
			Expect(si.GUID).To(Equal("215b97be-ec77-4224-9c38-c4f2d86b56c1"))
			Expect(si.Name).To(Equal("name-1523"))
			Expect(si.Type).To(Equal("managed_service_instance"))
		})
	})

	Describe("get service plan map", func() {
		var servicePlanJSON []string

		BeforeEach(func() {
			servicePlanJSON = slurp("test-data/service_plan.json")
		})

		It("should return an error when the service plan url fails", func() {
			fakeCliConnection.CliCommandWithoutTerminalOutputReturns(nil, errors.New("Bad Things"))
			_, err := api.GetServicePlanMap()
			Expect(err).ToNot(BeNil())
		})

		It("should return a map containing a specific element with all entries set", func() {
			fakeCliConnection.CliCommandWithoutTerminalOutputReturns(servicePlanJSON, nil)
			spMap, err := api.GetServicePlanMap()

			Expect(err).To(BeNil())
			Expect(spMap).NotTo(BeNil())

			s, exists := spMap["6fecf53b-7553-4cb3-b97e-930f9c4e3385"]

			Expect(exists).To(BeTrue())
			Expect(s.GUID).To(Equal("6fecf53b-7553-4cb3-b97e-930f9c4e3385"))
			Expect(s.ServiceGUID).To(Equal("1ccab853-87c9-45a6-bf99-603032d17fe5"))
		})

	})

	Describe("get service map", func() {
		var serviceJSON []string

		BeforeEach(func() {
			serviceJSON = slurp("test-data/services.json")
		})

		It("should return an error when the services url call fails", func() {
			fakeCliConnection.CliCommandWithoutTerminalOutputReturns(nil, errors.New("Bad Things"))
			_, err := api.GetServiceInstanceMap()
			Expect(err).ToNot(BeNil())
		})

		It("should return a map containing a specific element with all entries set", func() {
			fakeCliConnection.CliCommandWithoutTerminalOutputReturns(serviceJSON, nil)
			spMap, err := api.GetServiceMap()

			Expect(err).To(BeNil())
			Expect(spMap).NotTo(BeNil())

			s, exists := spMap["1993218f-096d-4216-bf9d-e0f250332dc6"]

			Expect(exists).To(BeTrue())
			Expect(s.GUID).To(Equal("1993218f-096d-4216-bf9d-e0f250332dc6"))
			Expect(s.Label).To(Equal("label-57"))
		})

	})

	Describe("get user provided service map", func() {
		var serviceJSON []string

		BeforeEach(func() {
			serviceJSON = slurp("test-data/user_provided_service_instances.json")
		})

		It("should return an error when the user provided services url call fails", func() {
			fakeCliConnection.CliCommandWithoutTerminalOutputReturns(nil, errors.New("Bad Things"))
			_, err := api.GetUserProvidedServiceMap()
			Expect(err).ToNot(BeNil())
		})

		It("should return a map containing a specific element with all entries set", func() {
			fakeCliConnection.CliCommandWithoutTerminalOutputReturns(serviceJSON, nil)
			sMap, err := api.GetUserProvidedServiceMap()

			Expect(err).To(BeNil())
			Expect(sMap).NotTo(BeNil())

			s, exists := sMap["54e4c645-7d20-4271-8c27-8cc904e1e7ee"]

			Expect(exists).To(BeTrue())
			Expect(s.Name).To(Equal("name-1696"))
			Expect(s.Type).To(Equal("user_provided_service_instance"))
		})
	})

	Describe("get service bindings list", func() {
		var serviceBindingJSON []string

		BeforeEach(func() {
			serviceBindingJSON = slurp("test-data/service_bindings.json")
		})

		It("should return an error when the service bindings url fails", func() {
			fakeCliConnection.CliCommandWithoutTerminalOutputReturns(nil, errors.New("Bad Things"))
			_, err := api.GetServiceBindingsList()
			Expect(err).ToNot(BeNil())
		})

		It("should return a list of service bindings with all required entries set", func() {
			fakeCliConnection.CliCommandWithoutTerminalOutputReturns(serviceBindingJSON, nil)
			sb, err := api.GetServiceBindingsList()

			Expect(err).To(BeNil())
			Expect(sb).NotTo(BeNil())
			Expect(len(sb)).To(Equal(1))

			Expect(sb[0].ServiceInstanceGUID).To(Equal("bde206e0-1ee8-48ad-b794-44c857633d50"))
			Expect(sb[0].AppGUID).To(Equal("b26e7e98-f002-41a8-a663-1b60f808a92a"))
		})
	})

	Describe("get space map", func() {
		var spaceJSON []string

		BeforeEach(func() {
			spaceJSON = slurp("test-data/spaces.json")
		})

		It("should return an error when the service bindings url fails", func() {
			fakeCliConnection.CliCommandWithoutTerminalOutputReturns(nil, errors.New("Bad Things"))
			_, err := api.GetSpaceMap()
			Expect(err).ToNot(BeNil())
		})

		It("should return a list of service bindings with all required entries set", func() {
			fakeCliConnection.CliCommandWithoutTerminalOutputReturns(spaceJSON, nil)
			sm, err := api.GetSpaceMap()

			Expect(err).To(BeNil())
			Expect(sm).NotTo(BeNil())

			space, exists := sm["de5db872-5b9e-4775-8d4a-f018133f9aaa"]
			Expect(exists).To(BeTrue())
			Expect(space.Name).To(Equal("jdk-space-2"))
			Expect(space.OrgGUID).To(Equal("b1a23fd6-ac8d-4304-a3b4-815745417acd"))
		})
	})

	Describe("get org map", func() {
		var orgJSON []string

		BeforeEach(func() {
			orgJSON = slurp("test-data/orgs.json")
		})

		It("should return an error when the org url call fails", func() {
			fakeCliConnection.CliCommandWithoutTerminalOutputReturns(nil, errors.New("Bad Things"))
			_, err := api.GetOrgMap()
			Expect(err).ToNot(BeNil())
		})

		It("should return a list of organisations with all required entries set", func() {
			fakeCliConnection.CliCommandWithoutTerminalOutputReturns(orgJSON, nil)
			om, err := api.GetOrgMap()

			Expect(err).To(BeNil())
			Expect(om).NotTo(BeNil())

			org, exists := om["b1a23fd6-ac8d-4304-a3b4-815745417acd"]
			Expect(exists).To(BeTrue())
			Expect(org.Name).To(Equal("jdk-org"))

			org2, exists2 := om["536a6736-0d89-4972-9e8c-0fbbb6802721"]
			Expect(exists2).To(BeTrue())
			Expect(org2.Name).To(Equal("test-org"))
		})
	})

})
