package cmd

import (
	"bytes"
	"strconv"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"courier_service/config"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func TestCmd(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Cmd Suite")
}

var _ = Describe("calculateDeliveryCost", func() {
	var (
		baseDeliveryCost int
		weight           int
		distance         int
		offerCode        string
		offers           []config.Offer
	)

	BeforeEach(func() {
		// Set up default offers
		offers = []config.Offer{
			{Code: "OFR001", Discount: 0.1, MinDistance: 70, MaxDistance: 200, MinWeight: 200, MaxWeight: 200},
			{Code: "OFR002", Discount: 0.07, MinDistance: 50, MaxDistance: 150, MinWeight: 100, MaxWeight: 250},
			{Code: "OFR003", Discount: 0.05, MinDistance: 10, MaxDistance: 150, MinWeight: 50, MaxWeight: 250},
		}

		viper.Set("weightCostPerKG", 10)
		viper.Set("distanceCostPerKM", 5)

	})

	Context("when no offer is applicable", func() {
		It("should return total cost without discount", func() {
			baseDeliveryCost = 100
			weight = 5
			distance = 5
			offerCode = "OFR001"

			totalCost, weightCost, distanceCost, discount, discountReason := calculateDeliveryCost(baseDeliveryCost, weight, distance, offerCode, offers)

			Expect(totalCost).To(Equal(float64(175)))
			Expect(weightCost).To(Equal(float64(50)))
			Expect(distanceCost).To(Equal(float64(25)))
			Expect(discount).To(Equal(float64(0)))
			Expect(discountReason).To(Equal("Offer not applicable as criteria not met"))
		})
	})

	Context("when offer OFR003 is applicable", func() {
		It("should return total cost with discount", func() {
			baseDeliveryCost = 100
			weight = 70
			distance = 100
			offerCode = "OFR003"

			totalCost, weightCost, distanceCost, discount, discountReason := calculateDeliveryCost(baseDeliveryCost, weight, distance, offerCode, offers)

			Expect(totalCost).To(Equal(float64(1300)))
			Expect(weightCost).To(Equal(float64(700)))
			Expect(distanceCost).To(Equal(float64(500)))
			Expect(discount).To(Equal(float64(65)))
			Expect(discountReason).To(Equal("Discount of 5% applied"))
		})
	})

	Context("when an invalid offer code is provided", func() {
		It("should return total cost without discount", func() {
			baseDeliveryCost = 100
			weight = 10
			distance = 10
			offerCode = "INVALID"

			totalCost, weightCost, distanceCost, discount, discountReason := calculateDeliveryCost(baseDeliveryCost, weight, distance, offerCode, offers)

			Expect(totalCost).To(Equal(float64(250)))
			Expect(weightCost).To(Equal(float64(100)))
			Expect(distanceCost).To(Equal(float64(50)))
			Expect(discount).To(Equal(float64(0)))
			Expect(discountReason).To(Equal("Offer not applicable as criteria not met"))
		})
	})
})

var _ = Describe("CalculateCmd", func() {
	var (
		output *bytes.Buffer
	)

	BeforeEach(func() {
		output = new(bytes.Buffer)
		rootCmd.SetOut(output)
		rootCmd.SetErr(output)
		mockOffers := []config.Offer{
			{Code: "OFR001", Discount: 0.1, MinDistance: 70, MaxDistance: 200, MinWeight: 70, MaxWeight: 200},
			{Code: "OFR002", Discount: 0.07, MinDistance: 50, MaxDistance: 150, MinWeight: 50, MaxWeight: 150},
			{Code: "OFR003", Discount: 0.05, MinDistance: 10, MaxDistance: 150, MinWeight: 50, MaxWeight: 250},
		}

		viper.Set("offers", mockOffers)
	})

	It("should calculate delivery cost correctly", func() {

		baseDeliveryCost := 100
		numPackages := 3
		pkgArgs := []string{
			"PKG1 5 5 OFR001",
			"PKG2 15 5 OFR002",
			"PKG3 10 100 OFR003",
		}

		args := append([]string{strconv.Itoa(baseDeliveryCost), strconv.Itoa(numPackages)}, pkgArgs...)

		calculateCmd.SetArgs(args)
		err := calculateCmd.Execute()

		Expect(err).To(BeNil())

	})

	Context("when arguments are missing", func() {
		It("should return an error if arguments are missing", func() {
			err := calculateCmd.RunE(&cobra.Command{}, []string{})
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("Usage: courier_service calculateCost <baseDeliveryCost> <numberOfPackages> <packages>"))
		})
	})

	Context("when base delivery cost is invalid", func() {
		It("should return an error if base delivery cost is not a number", func() {
			err := calculateCmd.RunE(&cobra.Command{}, []string{"abc", "2", "PKG1 5 5 OFR001"})
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("Invalid base delivery cost"))
		})
	})

	Context("when number of packages is invalid", func() {
		It("should return an error if number of packages is not a number", func() {
			err := calculateCmd.RunE(&cobra.Command{}, []string{"100", "abc", "PKG1 5 5 OFR001"})
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("Invalid number of packages"))
		})
	})

	Context("when package details are invalid", func() {
		It("should return an error if package details are incomplete", func() {
			err := calculateCmd.RunE(&cobra.Command{}, []string{"100", "2", "PKG1 5 5"})
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("Invalid package details for package 1\n"))
		})

		It("should return an error if weight is not a number", func() {
			err := calculateCmd.RunE(&cobra.Command{}, []string{"100", "2", "PKG1 abc 5 OFR001"})
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("Invalid weight for package 1\n"))
		})

		It("should return an error if distance is not a number", func() {
			err := calculateCmd.RunE(&cobra.Command{}, []string{"100", "2", "PKG1 5 abc OFR001"})
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("Invalid distance for package 1\n"))
		})
	})

	Context("when package details are valid", func() {
		It("should return an error if package details are incomplete", func() {
			err := calculateCmd.RunE(&cobra.Command{}, []string{"100", "1", "PKG1 5 5 OFR001"})
			Expect(err).To(BeNil())
		})

	})

})
