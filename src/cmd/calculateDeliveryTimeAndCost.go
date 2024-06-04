package cmd

import (
	"courier_service/config"
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

type Package struct {
	ID           string
	Weight       int
	Distance     int
	OfferCode    string
	TotalCost    float64
	Discount     float64
	FinalCost    float64
	DeliveryTime float64
}

type Vehicle struct {
	ID               int
	AssignedPackages []string
}

func getShipmentsSubSetsWhichFallsUnderMaxCarriable(packageList []Package, maxCarriableCapacity int) [][]int {

	var possiblePackages [][]int
	localHighestSum := 0

	for i := 1; i < (1 << len(packageList)); i++ {
		var subset []int
		subsetWeight := 0

		for j := 0; j < len(packageList); j++ {
			if i&(1<<j) != 0 {
				subsetWeight += packageList[j].Weight
				subset = append(subset, j)
			}
		}

		if subsetWeight <= maxCarriableCapacity && subsetWeight >= localHighestSum {
			if subsetWeight > localHighestSum {
				possiblePackages = nil
				localHighestSum = subsetWeight
			}
			possiblePackages = append(possiblePackages, subset)
		}
	}

	return possiblePackages
}

func getShipmentWithLessDistanceAmongPossibleSubsets(possibleShipmentList [][]int, packageList []Package) []int {
	if len(possibleShipmentList) == 1 {
		return possibleShipmentList[0]
	}

	var distanceList []int

	for _, element := range possibleShipmentList {
		maxDistance := 0
		for _, idx := range element {
			if packageList[idx].Distance > maxDistance {
				maxDistance = packageList[idx].Distance
			}
		}
		distanceList = append(distanceList, maxDistance)
	}

	minDistance := math.MaxInt32
	var closestShipment []int

	for i, distance := range distanceList {
		if distance < minDistance {
			minDistance = distance
			closestShipment = possibleShipmentList[i]
		}
	}

	return closestShipment
}

func calculateDeliveryTime(packages []Package, numVehicles, maxSpeed, maxWeight int, baseDeliveryCost int) {
	vehicleAvailabilityArray := make([]float64, numVehicles)
	vehicleList := make([]Vehicle, numVehicles)

	for i := range vehicleAvailabilityArray {
		vehicleAvailabilityArray[i] = 0
		vehicleList[i] = Vehicle{ID: i + 1}
	}

	var newUpdatedPackageList []Package
	for _, pkg := range packages {
		newUpdatedPackageList = append(newUpdatedPackageList, pkg)
	}

	for len(newUpdatedPackageList) > 0 {
		possibleShipmentList := getShipmentsSubSetsWhichFallsUnderMaxCarriable(newUpdatedPackageList, maxWeight)
		nextDelivery := getShipmentWithLessDistanceAmongPossibleSubsets(possibleShipmentList, newUpdatedPackageList)

		nextAvailableAt := math.MaxFloat64
		for _, availability := range vehicleAvailabilityArray {
			if availability < nextAvailableAt {
				nextAvailableAt = availability
			}
		}

		durationForSingleTrip := 0.0
		var vehicleID int

		for _, idx := range nextDelivery {
			currentPackage := &newUpdatedPackageList[idx]
			deliveryTime := float64(currentPackage.Distance) / float64(maxSpeed)
			currentPackage.DeliveryTime = nextAvailableAt + deliveryTime
			durationForSingleTrip = math.Max(deliveryTime, durationForSingleTrip)

			for i := range vehicleList {
				if vehicleAvailabilityArray[i] == nextAvailableAt {
					vehicleList[i].AssignedPackages = append(vehicleList[i].AssignedPackages, currentPackage.ID)
					vehicleID = i + 1
					break
				}
			}

			fmt.Printf("Package: %s\n", currentPackage.ID)
			fmt.Printf("  Vehicle: %d\n", vehicleID)
			fmt.Printf("  Discount: %.2f\n", currentPackage.Discount)
			fmt.Printf("  Total Cost: %.2f\n", currentPackage.TotalCost)
			fmt.Printf("  Delivery Time: %.2f hours\n", currentPackage.DeliveryTime)
			fmt.Println()
		}

		for i, availability := range vehicleAvailabilityArray {
			if availability == nextAvailableAt {
				vehicleAvailabilityArray[i] = nextAvailableAt + 2*durationForSingleTrip
				break
			}
		}

		var remainingPackages []Package
		for _, pkg := range newUpdatedPackageList {
			if pkg.DeliveryTime == 0 {
				remainingPackages = append(remainingPackages, pkg)
			}
		}
		newUpdatedPackageList = remainingPackages
	}
}

var calculateTimeAndCostCmd = &cobra.Command{
	Use:   "calculateTimeAndCost",
	Short: "Calculate delivery cost and time of packages",
	Long:  `This command calculates the delivery cost and time of packages based on weight, distance, and offer codes.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 4 {
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

		var packages []Package

		for i := 0; i < numPackages; i++ {
			packageDetails := strings.Fields(args[2+i])
			if len(packageDetails) != 4 {
				return fmt.Errorf("Invalid package details for package %d", i+1)
			}

			weight, err := strconv.Atoi(packageDetails[1])
			if err != nil {
				return fmt.Errorf("Invalid weight for package %d", i+1)
			}

			distance, err := strconv.Atoi(packageDetails[2])
			if err != nil {
				return fmt.Errorf("Invalid distance for package %d", i+1)
			}

			offers := config.GetOffers()

			totalCost, _, _, discount, _ := calculateDeliveryCost(baseDeliveryCost, weight, distance, packageDetails[3], offers)
			finalCost := totalCost - discount

			packages = append(packages, Package{
				ID:        packageDetails[0],
				Weight:    weight,
				Distance:  distance,
				OfferCode: packageDetails[3],
				TotalCost: totalCost,
				Discount:  discount,
				FinalCost: finalCost,
			})
		}

		numVehicles, err := strconv.Atoi(args[2+numPackages])
		if err != nil {
			return fmt.Errorf("Invalid number of vehicles")
		}

		maxSpeed, err := strconv.Atoi(args[3+numPackages])
		if err != nil {
			return fmt.Errorf("Invalid max vehicle speed")
		}

		maxWeight, err := strconv.Atoi(args[4+numPackages])
		if err != nil {
			return fmt.Errorf("Invalid vehicle capacity")
		}

		calculateDeliveryTime(packages, numVehicles, maxSpeed, maxWeight, baseDeliveryCost)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(calculateTimeAndCostCmd)
}
