package cmd

import (
	"courier_service/config"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

type Package struct {
	ID       string
	Weight   int
	Distance int
	Offer    string
}

type Vehicle struct {
	ID            int
	Speed         int
	Capacity      int
	AvailableTime float64
}

type DeliveryResult struct {
	VehicleID                    int
	PackageID                    string
	Discount                     float64
	TotalCost                    float64
	DeliveryTime                 float64
	CurrentVehicleAvailableAfter float64
}

type Offer struct {
	Code        string
	Discount    float64
	MinDistance int
	MaxDistance int
	MinWeight   int
	MaxWeight   int
}

var calculateTimeAndCostCmd = &cobra.Command{
	Use:   "calculateTimeAndCost",
	Short: "Calculate delivery cost and estimated delivery time for packages",
	Long:  `This command calculates the delivery cost and estimated delivery time for packages based on weight, distance, and offer codes.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 6 {
			return fmt.Errorf("Usage: courier_service calculateTimeAndCost <baseDeliveryCost> <numberOfPackages> <packages> <number_of_vehicles> <max_speed> <max_carriable_weight>")
		}

		baseDeliveryCost, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("Invalid base delivery cost")
		}

		numPackages, err := strconv.Atoi(args[1])
		if err != nil {
			return fmt.Errorf("Invalid number of packages")
		}

		if len(args) < 2+numPackages+3 {
			return fmt.Errorf("Insufficient package information provided")
		}

		offers := config.GetOffers()

		packages := []Package{}
		for i := 0; i < numPackages; i++ {
			fields := strings.Fields(args[2+i])
			if len(fields) != 4 {
				return fmt.Errorf("Invalid package details for package %d", i+1)
			}
			weight, err := strconv.Atoi(fields[1])
			if err != nil {
				return fmt.Errorf("Invalid weight for package %d", i+1)
			}
			distance, err := strconv.Atoi(fields[2])
			if err != nil {
				return fmt.Errorf("Invalid distance for package %d", i+1)
			}
			packages = append(packages, Package{
				ID:       fields[0],
				Weight:   weight,
				Distance: distance,
				Offer:    fields[3],
			})
		}

		vehicleInfoIndex := 2 + numPackages
		numVehicles, err := strconv.Atoi(args[vehicleInfoIndex])
		if err != nil {
			return fmt.Errorf("Invalid number of vehicles")
		}
		vehicleSpeed, err := strconv.Atoi(args[vehicleInfoIndex+1])
		if err != nil {
			return fmt.Errorf("Invalid vehicle speed")
		}
		vehicleCapacity, err := strconv.Atoi(args[vehicleInfoIndex+2])
		if err != nil {
			return fmt.Errorf("Invalid vehicle capacity")
		}

		vehicles := make([]Vehicle, numVehicles)
		for i := 0; i < numVehicles; i++ {
			vehicles[i] = Vehicle{ID: i + 1, Speed: vehicleSpeed, Capacity: vehicleCapacity}
		}

		results := []DeliveryResult{}

		for len(packages) > 0 {

			sort.SliceStable(packages, func(i, j int) bool {
				if packages[i].Weight == packages[j].Weight {
					return packages[i].Distance < packages[j].Distance
				}
				return packages[i].Weight > packages[j].Weight
			})

			currentPackages := []Package{}
			currentWeight := 0

			for _, pkg := range packages {
				if currentWeight+pkg.Weight <= vehicleCapacity {
					currentPackages = append(currentPackages, pkg)
					currentWeight += pkg.Weight
				}
			}

			remainingPackages := []Package{}
			selectedPackageIDs := map[string]bool{}
			for _, pkg := range currentPackages {
				selectedPackageIDs[pkg.ID] = true
			}
			for _, pkg := range packages {
				if !selectedPackageIDs[pkg.ID] {
					remainingPackages = append(remainingPackages, pkg)
				}
			}
			packages = remainingPackages

			sort.SliceStable(vehicles, func(i, j int) bool {
				return vehicles[i].AvailableTime < vehicles[j].AvailableTime
			})

			assignedVehicle := &vehicles[0]

			maxDistance := 0
			for _, pkg := range currentPackages {
				if pkg.Distance > maxDistance {
					maxDistance = pkg.Distance
				}
			}

			tripTime := float64(maxDistance) / float64(assignedVehicle.Speed)
			roundTripTime := 2 * tripTime

			for _, pkg := range currentPackages {
				totalCost, _, _, discount, _ := calculateDeliveryCost(baseDeliveryCost, pkg.Weight, pkg.Distance, pkg.Offer, offers)
				deliveryTime := float64(pkg.Distance) / float64(assignedVehicle.Speed)

				results = append(results, DeliveryResult{
					VehicleID:                    assignedVehicle.ID,
					PackageID:                    pkg.ID,
					Discount:                     discount,
					TotalCost:                    totalCost - discount,
					DeliveryTime:                 assignedVehicle.AvailableTime + deliveryTime,
					CurrentVehicleAvailableAfter: assignedVehicle.AvailableTime + roundTripTime,
				})
			}

			assignedVehicle.AvailableTime += roundTripTime

		}

		// Output the results
		for _, result := range results {
			fmt.Printf("Package: %s\n", result.PackageID)
			fmt.Printf("  Vehicle: %d\n", result.VehicleID)
			fmt.Printf("  Discount: %.2f\n", result.Discount)
			fmt.Printf("  Total Cost: %.2f\n", result.TotalCost)
			fmt.Printf("  Delivery Time: %.2f hours\n", result.DeliveryTime)
			fmt.Printf("  Current Vehicle Available for next delivery after : %.2f hours\n", result.CurrentVehicleAvailableAfter)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(calculateTimeAndCostCmd)
}
