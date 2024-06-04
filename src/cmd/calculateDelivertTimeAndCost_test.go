package cmd

import (
	"bytes"
	"fmt"
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("CalculateTimeAndCostCmd", func() {
	var (
		stdout *os.File
		r, w   *os.File
		output bytes.Buffer
	)

	BeforeEach(func() {
		// Save the original stdout
		stdout = os.Stdout

		// Create a pipe to capture stdout
		r, w, _ = os.Pipe()
		os.Stdout = w
	})

	AfterEach(func() {
		// Close the writer and restore original stdout
		w.Close()
		os.Stdout = stdout
	})
	Describe("RunE", func() {
		Context("with valid input", func() {
			It("should calculate delivery time and cost correctly", func() {

				args := []string{"100", "3", "PKG1 50 30 OFR001", "PKG2 75 125 OFFR0008", "PKG3 175 100 OFFR003", "2", "70", "200"}

				err := calculateTimeAndCostCmd.RunE(nil, args)

				w.Close()
				output.ReadFrom(r)

				Expect(output.String()).To(ContainSubstring("Package: PKG3"))
				Expect(output.String()).To(ContainSubstring("Delivery Time: 1.43 hours"))

				Expect(output.String()).To(ContainSubstring("Package: PKG1"))
				Expect(output.String()).To(ContainSubstring("Delivery Time: 0.43 hours"))

				Expect(output.String()).To(ContainSubstring("Package: PKG2"))
				Expect(output.String()).To(ContainSubstring("Delivery Time: 1.79 hours"))

				Expect(err).ToNot(HaveOccurred())

			})

			It("should calculate delivery time and cost correctly when weights are same for two pacakges", func() {

				args := []string{"100", "3", "PKG1 50 30 OFR001", "PKG2 50 125 OFFR0008", "PKG3 175 100 OFFR003", "2", "70", "200"}

				err := calculateTimeAndCostCmd.RunE(nil, args)

				w.Close()
				output.ReadFrom(r)

				fmt.Println(output.String())

				Expect(output.String()).To(ContainSubstring("Package: PKG3"))
				Expect(output.String()).To(ContainSubstring("Delivery Time: 1.43 hours"))

				Expect(output.String()).To(ContainSubstring("Package: PKG1"))
				Expect(output.String()).To(ContainSubstring("Delivery Time: 0.43 hours"))

				Expect(output.String()).To(ContainSubstring("Package: PKG2"))
				Expect(output.String()).To(ContainSubstring("Delivery Time: 1.79 hours"))

				Expect(err).ToNot(HaveOccurred())
			})

			It("should calculate delivery time and cost correctly when weights are same for two packages and has multiple subsets selected then it should select based on distance", func() {

				args := []string{"100", "3", "PKG1 50 30 OFR001", "PKG2 50 125 OFFR0008", "PKG3 125 100 OFFR003", "2", "70", "200"}

				err := calculateTimeAndCostCmd.RunE(nil, args)

				w.Close()
				output.ReadFrom(r)

				fmt.Println(output.String())

				Expect(output.String()).To(ContainSubstring("Package: PKG3"))
				Expect(output.String()).To(ContainSubstring("Delivery Time: 1.43 hours"))

				Expect(output.String()).To(ContainSubstring("Package: PKG1"))
				Expect(output.String()).To(ContainSubstring("Delivery Time: 0.43 hours"))

				Expect(output.String()).To(ContainSubstring("Package: PKG2"))
				Expect(output.String()).To(ContainSubstring("Delivery Time: 1.79 hours"))

				Expect(err).ToNot(HaveOccurred())
			})
		})

		Context("with invalid input", func() {

			It("should return an error when insufficient package information", func() {

				args := []string{"100", "3", "PKG1 50 30 OFR001", "PKG2 75 125 OFFR0008", "2", "70", "200"}

				err := calculateTimeAndCostCmd.RunE(nil, args)

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("Invalid package details for package 3"))
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
				Expect(err.Error()).To(Equal("Invalid max vehicle speed"))
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
