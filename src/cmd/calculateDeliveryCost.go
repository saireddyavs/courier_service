package cmd

import (
	"fmt"
	"strconv"
	"strings"

	"courier_service/config"

	"github.com/spf13/cobra"
)

func calculateDeliveryCost(baseDeliveryCost, weight, distance int, offerCode string, offers []config.Offer) (float64, float64, float64, float64, string) {
	weightCost := float64(weight * config.GetWeightCostPerKG())
	distanceCost := float64(distance * config.GetDistanceCostPerKM())
	totalCost := float64(baseDeliveryCost) + weightCost + distanceCost
	discount := 0.0
	discountReason := "Offer not applicable as criteria not met"

	for _, offer := range offers {
		if offer.Code == offerCode {
			if distance >= offer.MinDistance && distance <= offer.MaxDistance && weight >= offer.MinWeight && weight <= offer.MaxWeight {
				discount = totalCost * offer.Discount
				discountReason = fmt.Sprintf("Discount of %.0f%% applied", offer.Discount*100)
			}
			break
		}
	}

	// fmt.Println(totalCost, weightCost, distanceCost, discount, discountReason)

	return totalCost, weightCost, distanceCost, discount, discountReason
}

var calculateCmd = &cobra.Command{
	Use:   "calculateCost",
	Short: "Calculate delivery cost of packages",
	Long:  `This command calculates the delivery cost of packages based on weight, distance, and offer codes.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 3 {
			return fmt.Errorf("Usage: courier_service calculateCost <baseDeliveryCost> <numberOfPackages> <packages>")
		}

		baseDeliveryCost, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("Invalid base delivery cost")

		}

		numPackages, err := strconv.Atoi(args[1])
		if err != nil {
			return fmt.Errorf("Invalid number of packages")
		}

		offers := config.GetOffers()

		for i := 0; i < numPackages; i++ {
			packageDetails := strings.Fields(args[2+i])
			if len(packageDetails) != 4 {
				return fmt.Errorf("Invalid package details for package %d\n", i+1)
			}

			weight, err := strconv.Atoi(packageDetails[1])
			if err != nil {
				return fmt.Errorf("Invalid weight for package %d\n", i+1)
			}

			distance, err := strconv.Atoi(packageDetails[2])
			if err != nil {
				return fmt.Errorf("Invalid distance for package %d\n", i+1)
			}

			offerCode := packageDetails[3]

			totalCost, weightCost, distanceCost, discount, discountReason := calculateDeliveryCost(baseDeliveryCost, weight, distance, offerCode, offers)

			finalCost := totalCost - discount
			fmt.Printf("\nPackage %s\n", packageDetails[0])
			fmt.Printf("Base Delivery Cost: %d\n", baseDeliveryCost)
			fmt.Printf("Weight: %d kg | Distance: %d km\n", weight, distance)
			fmt.Printf("Offer code: %s\n", offerCode)
			fmt.Printf("Discount: %.2f (%s)\n", discount, discountReason)
			fmt.Printf("Breakdown:\n")
			fmt.Printf("  Base Delivery Cost: %.2f\n", float64(baseDeliveryCost))
			fmt.Printf("  Weight Cost: %.2f\n", weightCost)
			fmt.Printf("  Distance Cost: %.2f\n", distanceCost)
			fmt.Printf("  Discount: -%.2f\n", discount)
			fmt.Printf("Total Delivery Cost: %.2f\n", finalCost)

		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(calculateCmd)
}
