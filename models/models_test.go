package models_test

import (
	. "github.com/dgruber/usagereport-plugin/models"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"io/ioutil"
)

var _ = Describe("Models", func() {

	Describe("Correct command line output of the models", func() {
		var report Report

		BeforeEach(func() {
			report = Report{
				Orgs: []Org{
					Org{
						Name:        "test-org",
						MemoryQuota: 4096,
						MemoryUsage: 256,
						Spaces: []Space{Space{
							Name: "test-space",
							Apps: []App{
								App{Ram: 128, Instances: 2, Running: true, SiTotal: 10, SiPCF: 6, SiUP: 2, Name: "sample"},
								App{Ram: 128, Instances: 1, Running: false, SiTotal: 4, SiPCF: 2, SiUP: 0, Name: "test"},
							},
						},
						},
					},
				},
				ServiceInstances: []Service{
					Service{
						ServiceInstanceGUID: "111",
						ServiceInstanceName: "serviceInstanceName",
						ServiceInstanceType: "serviceInstanceType",
						SpaceName:           "test-space",
						OrgName:             "test-org",
						ServicePlanName:     "servicePlanName",
						ServiceName:         "serviceName",
						ServiceType:         "serviceType",
						AppGUIDs: []string{
							"123",
							"321",
						},
					},
				},
			}
		})

		Describe("Report#CSV", func() {
			It("should return csv formated string", func() {
				expectedOutput, err := ioutil.ReadFile("fixtures/result.csv")
				Expect(err).ShouldNot(HaveOccurred())
				Expect(report.CSV()).To(Equal(string(expectedOutput)))
			})
		})

		Describe("Report#String", func() {
			It("should return human readable formated string", func() {
				expectedOutput, err := ioutil.ReadFile("fixtures/result.txt")
				Expect(err).ShouldNot(HaveOccurred())
				Expect(report.String()).To(Equal(string(expectedOutput)))
			})
		})

		Describe("ServicesApp#CSV", func() {
			It("should return csv formated string", func() {
				expectedOutput, err := ioutil.ReadFile("fixtures/services.csv")
				Expect(err).ShouldNot(HaveOccurred())
				Expect(report.ServiceInstanceReportCSV()).To(Equal(string(expectedOutput)))
			})
		})

		Describe("ServicesApp#String", func() {
			It("should return string formated string", func() {
				expectedOutput, err := ioutil.ReadFile("fixtures/services.txt")
				Expect(err).ShouldNot(HaveOccurred())
				Expect(report.ServiceInstanceReportString()).To(Equal(string(expectedOutput)))
			})
		})

		Describe("ServicesSummary#CSV", func() {
			It("should return csv formated string", func() {
				expectedOutput, err := ioutil.ReadFile("fixtures/servicesSummary.csv")
				Expect(err).ShouldNot(HaveOccurred())
				Expect(report.ServiceInstanceSummaryCSV()).To(Equal(string(expectedOutput)))
			})
		})

		Describe("ServicesSummary#String", func() {
			It("should return string formated string", func() {
				expectedOutput, err := ioutil.ReadFile("fixtures/servicesSummary.txt")
				Expect(err).ShouldNot(HaveOccurred())
				Expect(report.ServiceInstanceSummaryString()).To(Equal(string(expectedOutput)))
			})
		})
	})

	Describe("Internal report builder", func() {
		var r Report

		BeforeEach(func() {
			r = Report{
				ServiceInstances: []Service{
					Service{
						ServiceInstanceGUID: "111",
						ServiceInstanceName: "serviceInstanceName",
						ServiceInstanceType: "serviceInstanceType",
						SpaceName:           "test-space",
						OrgName:             "test-org",
						ServicePlanName:     "servicePlanName",
						ServiceName:         "serviceName",
						ServiceType:         "serviceType",
						AppGUIDs: []string{
							"123",
							"321",
						},
					},
					Service{
						ServiceInstanceGUID: "211",
						ServiceInstanceName: "serviceInstanceName2",
						ServiceInstanceType: "serviceInstanceType2",
						SpaceName:           "test-space2",
						OrgName:             "test-org2",
						ServicePlanName:     "servicePlanName",
						ServiceName:         "serviceName",
						ServiceType:         "serviceType",
						AppGUIDs: []string{
							"123",
							"321",
						},
					},
					Service{
						ServiceInstanceGUID: "311",
						ServiceInstanceName: "serviceInstanceName3",
						ServiceInstanceType: "serviceInstanceType2",
						SpaceName:           "test-space",
						OrgName:             "test-org",
						ServicePlanName:     "servicePlanName",
						ServiceName:         "serviceName",
						ServiceType:         "serviceType",
						AppGUIDs: []string{
							"123",
							"321",
						},
					},
					Service{
						ServiceInstanceGUID: "411",
						ServiceInstanceName: "serviceInstanceName4",
						ServiceInstanceType: "serviceInstanceType",
						SpaceName:           "test-space1",
						OrgName:             "test-org2",
						ServicePlanName:     "servicePlanName",
						ServiceName:         "serviceName",
						ServiceType:         "serviceType",
						AppGUIDs: []string{
							"123",
							"321",
						},
					},
				},
			}
		})

		It("should build up a correct Org and Space name hierarchy", func() {

			r.BuildOrgAndSpacesUsingServiceInstances()
			// should be idempotent
			r.BuildOrgAndSpacesUsingServiceInstances()

			Expect(len(r.Orgs)).To(Equal(2))

			for _, org := range r.Orgs {
				if org.Name == "test-org" {
					Expect(len(org.Spaces)).To(Equal(1))
					Expect(org.Spaces[0].Name).To(Equal("test-space"))
				} else if org.Name == "test-org2" {
					Expect(len(org.Spaces)).To(Equal(2))
					Expect(org.Spaces).To(ContainElement(Space{Name: "test-space1"}))
					Expect(org.Spaces).To(ContainElement(Space{Name: "test-space2"}))
				} else {
					Expect(org.Name).To(Equal("unknown org name"))
				}
			}
		})

	})

})
