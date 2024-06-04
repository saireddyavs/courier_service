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
	vehicleAvailabilityArray, vehicleList := initializeVehicles(numVehicles)
	newUpdatedPackageList := copyPackages(packages)

	for len(newUpdatedPackageList) > 0 {
		possibleShipmentList := getShipmentsSubSetsWhichFallsUnderMaxCarriable(newUpdatedPackageList, maxWeight)
		nextDelivery := getShipmentWithLessDistanceAmongPossibleSubsets(possibleShipmentList, newUpdatedPackageList)
		processNextDelivery(nextDelivery, newUpdatedPackageList, vehicleAvailabilityArray, vehicleList, maxSpeed)
		newUpdatedPackageList = filterRemainingPackages(newUpdatedPackageList)
	}
}

func initializeVehicles(numVehicles int) ([]float64, []Vehicle) {
	vehicleAvailabilityArray := make([]float64, numVehicles)
	vehicleList := make([]Vehicle, numVehicles)
	for i := range vehicleAvailabilityArray {
		vehicleAvailabilityArray[i] = 0
		vehicleList[i] = Vehicle{ID: i + 1}
	}
	return vehicleAvailabilityArray, vehicleList
}

func copyPackages(packages []Package) []Package {
	var newUpdatedPackageList []Package
	for _, pkg := range packages {
		newUpdatedPackageList = append(newUpdatedPackageList, pkg)
	}
	return newUpdatedPackageList
}

func getEarliestAvailableVehicle(vehicleAvailabilityArray []float64) float64 {
	nextAvailableAt := math.MaxFloat64
	for _, availability := range vehicleAvailabilityArray {
		if availability < nextAvailableAt {
			nextAvailableAt = availability
		}
	}
	return nextAvailableAt
}

func assignPackagesToVehicle(vehicleList []Vehicle, vehicleAvailabilityArray []float64, nextAvailableAt float64, durationForSingleTrip float64, nextDelivery []int, newUpdatedPackageList []Package, maxSpeed int) {
	var vehicleID int
	for _, idx := range nextDelivery {
		currentPackage := &newUpdatedPackageList[idx]
		currentPackage.DeliveryTime = nextAvailableAt + (float64(currentPackage.Distance) / float64(maxSpeed))
		for i := range vehicleList {
			if vehicleAvailabilityArray[i] == nextAvailableAt {
				vehicleList[i].AssignedPackages = append(vehicleList[i].AssignedPackages, currentPackage.ID)
				vehicleID = i + 1
				break
			}
		}
		printPackageDetails(currentPackage, vehicleID)
	}
	updateVehicleAvailability(vehicleAvailabilityArray, nextAvailableAt, durationForSingleTrip)
}

func printPackageDetails(pkg *Package, vehicleID int) {
	fmt.Printf("Package: %s\n", pkg.ID)
	fmt.Printf("  Vehicle: %d\n", vehicleID)
	fmt.Printf("  Discount: %.2f\n", pkg.Discount)
	fmt.Printf("  Total Cost: %.2f\n", pkg.TotalCost)
	fmt.Printf("  Delivery Time: %.2f hours\n", pkg.DeliveryTime)
	fmt.Println()
}

func updateVehicleAvailability(vehicleAvailabilityArray []float64, nextAvailableAt float64, durationForSingleTrip float64) {
	for i, availability := range vehicleAvailabilityArray {
		if availability == nextAvailableAt {
			vehicleAvailabilityArray[i] = nextAvailableAt + 2*durationForSingleTrip
			break
		}
	}
}

func filterRemainingPackages(newUpdatedPackageList []Package) []Package {
	var remainingPackages []Package
	for _, pkg := range newUpdatedPackageList {
		if pkg.DeliveryTime == 0 {
			remainingPackages = append(remainingPackages, pkg)
		}
	}
	return remainingPackages
}

func processNextDelivery(nextDelivery []int, newUpdatedPackageList []Package, vehicleAvailabilityArray []float64, vehicleList []Vehicle, maxSpeed int) {
	nextAvailableAt := getEarliestAvailableVehicle(vehicleAvailabilityArray)
	durationForSingleTrip := calculateDurationForSingleTrip(nextDelivery, newUpdatedPackageList, maxSpeed)
	assignPackagesToVehicle(vehicleList, vehicleAvailabilityArray, nextAvailableAt, durationForSingleTrip, nextDelivery, newUpdatedPackageList, maxSpeed)
}

func calculateDurationForSingleTrip(nextDelivery []int, newUpdatedPackageList []Package, maxSpeed int) float64 {
	durationForSingleTrip := 0.0
	for _, idx := range nextDelivery {
		currentPackage := &newUpdatedPackageList[idx]
		deliveryTime := float64(currentPackage.Distance) / float64(maxSpeed)
		durationForSingleTrip = math.Max(deliveryTime, durationForSingleTrip)
	}
	return durationForSingleTrip
}

var calculateTimeAndCostCmd = &cobra.Command{
	Use:   "calculateTimeAndCost",
	Short: "Calculate delivery cost and time of packages",
	Long:  `This command calculates the delivery cost and time of packages based on weight, distance, and offer codes.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 5 {
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

		packages, err := parsePackages(args[2:numPackages+2], baseDeliveryCost)
		if err != nil {
			return err
		}

		numVehicles, err := strconv.Atoi(args[2+numPackages])
		if err != nil {
			return fmt.Errorf("Invalid number of vehicles")
		}

		maxSpeed, err := strconv.Atoi(args[3+numPackages])
		if err != nil {
			return fmt.Errorf("Invalid max vehicle speed")
		}

		maxLoadCapacity, err := strconv.Atoi(args[4+numPackages])
		if err != nil {
			return fmt.Errorf("Invalid vehicle capacity")
		}

		calculateDeliveryTime(packages, numVehicles, maxSpeed, maxLoadCapacity, baseDeliveryCost)

		return nil
	},
}

func parsePackages(packageArgs []string, baseDeliveryCost int) ([]Package, error) {
	var packages []Package
	offers := config.GetOffers()

	for i, arg := range packageArgs {
		packageDetails := strings.Fields(arg)
		if len(packageDetails) != 4 {
			return nil, fmt.Errorf("Invalid package details for package %d", i+1)
		}

		weight, err := strconv.Atoi(packageDetails[1])
		if err != nil {
			return nil, fmt.Errorf("Invalid weight for package %d", i+1)
		}

		distance, err := strconv.Atoi(packageDetails[2])
		if err != nil {
			return nil, fmt.Errorf("Invalid distance for package %d", i+1)
		}

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

	return packages, nil
}

func init() {
	rootCmd.AddCommand(calculateTimeAndCostCmd)
}
