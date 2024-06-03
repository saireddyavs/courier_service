package cmd

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("CalculateTimeAndCostCmd", func() {
	Describe("RunE", func() {
		Context("with valid input", func() {
			It("should calculate delivery time and cost correctly", func() {

				args := []string{"100", "3", "PKG1 50 30 OFR001", "PKG2 75 125 OFFR0008", "PKG3 175 100 OFFR003", "2", "70", "200"}

				err := calculateTimeAndCostCmd.RunE(nil, args)

				Expect(err).ToNot(HaveOccurred())
			})

			It("should calculate delivery time and cost correctly when weights are same for two pacakges", func() {

				args := []string{"100", "3", "PKG1 50 30 OFR001", "PKG2 50 125 OFFR0008", "PKG3 175 100 OFFR003", "2", "70", "200"}

				err := calculateTimeAndCostCmd.RunE(nil, args)

				Expect(err).ToNot(HaveOccurred())
			})
		})

		Context("with invalid input", func() {

			It("should return an error when insufficient package information", func() {

				args := []string{"100", "3", "PKG1 50 30 OFR001", "PKG2 75 125 OFFR0008", "2", "70", "200"}

				err := calculateTimeAndCostCmd.RunE(nil, args)

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("Insufficient package information provided"))
			})

			It("should return an error for invalid base delivery cost", func() {

				args := []string{"abc", "3", "50", "150", "PKG1 50 30 OFR001", "PKG2 75 125 OFFR0008", "PKG3 175 100 OFFR003", "2", "70", "200"}

				err := calculateTimeAndCostCmd.RunE(nil, args)

				Expect(err).To(HaveOccurred())
				// Assert the actual error message
				Expect(err.Error()).To(Equal("Invalid base delivery cost"))
			})

			It("should return an error when invalid base delivery cost", func() {

				args := []string{"abc", "3", "PKG1 50 30 OFR001", "PKG2 75 125 OFFR0008", "PKG3 175 100 OFFR003", "2", "70", "200"}

				err := calculateTimeAndCostCmd.RunE(nil, args)

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("Invalid base delivery cost"))
			})

			It("should return an error when invalid package details", func() {

				args := []string{"100", "3", "PKG1 50 abc OFR001", "PKG2 75 125 OFFR0008", "PKG3 175 100 OFFR003", "2", "70", "200"}

				err := calculateTimeAndCostCmd.RunE(nil, args)

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("Invalid distance for package 1"))
			})

			It("should return an error when invalid number of args passed", func() {

				args := []string{"100"}

				err := calculateTimeAndCostCmd.RunE(nil, args)

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("Usage: courier_service calculateTimeAndCost <baseDeliveryCost> <numberOfPackages> <packages> <number_of_vehicles> <max_speed> <max_carriable_weight>"))
			})

			It("should return an error when invalid number of packages is passed", func() {

				args := []string{"100", "3A", "PKG1 50 abc OFR001", "PKG2 75 125 OFFR0008", "PKG3 175 100 OFFR003", "2", "70", "200"}

				err := calculateTimeAndCostCmd.RunE(nil, args)

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("Invalid number of packages"))
			})

			It("should return an error when invalid number of fields inside package", func() {

				args := []string{"100", "3", "PKG1 50 abc OFR001 EXTRA", "PKG2 75 125 OFFR0008", "PKG3 175 100 OFFR003", "2", "70", "200"}

				err := calculateTimeAndCostCmd.RunE(nil, args)

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("Invalid package details for package 1"))
			})

			It("should return an error when invalid weight fields inside package", func() {

				args := []string{"100", "3", "PKG1 50A 100 OFR001", "PKG2 75 125 OFFR0008", "PKG3 175 100 OFFR003", "2", "70", "200"}

				err := calculateTimeAndCostCmd.RunE(nil, args)

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("Invalid weight for package 1"))
			})

			It("should return an error when invalid number of vehicles", func() {

				args := []string{"100", "3", "PKG1 50 100 OFR001", "PKG2 75 125 OFFR0008", "PKG3 175 100 OFFR003", "2V", "70", "200"}

				err := calculateTimeAndCostCmd.RunE(nil, args)

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("Invalid number of vehicles"))
			})

			It("should return an error when invalid vehicle speed", func() {

				args := []string{"100", "3", "PKG1 50 100 OFR001", "PKG2 75 125 OFFR0008", "PKG3 175 100 OFFR003", "2", "70S", "200"}

				err := calculateTimeAndCostCmd.RunE(nil, args)

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("Invalid vehicle speed"))
			})

			It("should return an error when invalid vehicle capacity", func() {

				args := []string{"100", "3", "PKG1 50 100 OFR001", "PKG2 75 125 OFFR0008", "PKG3 175 100 OFFR003", "2", "70", "200W"}

				err := calculateTimeAndCostCmd.RunE(nil, args)

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("Invalid vehicle capacity"))
			})
		})
	})
})
