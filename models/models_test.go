package models_test

import (
	. "github.com/dgruber/usagereport-plugin/models"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"io/ioutil"
)

var _ = Describe("Models", func() {
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

	Describe("Services#CSV", func() {
		It("should return csv formated string", func() {
			expectedOutput, err := ioutil.ReadFile("fixtures/services.csv")
			Expect(err).ShouldNot(HaveOccurred())
			Expect(report.ServiceInstanceReportCSV()).To(Equal(string(expectedOutput)))
		})
	})

	Describe("Services#String", func() {
		It("should return string formated string", func() {
			expectedOutput, err := ioutil.ReadFile("fixtures/services.txt")
			Expect(err).ShouldNot(HaveOccurred())
			Expect(report.ServiceInstanceReportString()).To(Equal(string(expectedOutput)))
		})
	})

})
