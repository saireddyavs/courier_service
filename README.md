# Courier Service CLI

This is a command line application to calculate the delivery cost and estimated delivery time for packages based on weight, distance, and offer codes.

## How to run

Install the required dependencies:

```
go mod tidy
```

Build the executable:

```
go build .
```

## Usage

### calculate Command

This command calculates the delivery cost of packages.

```
./courier_service calculateCost <baseDeliveryCost> <numberOfPackages> <packages>
```

#### Arguments

- baseDeliveryCost: Base delivery cost (integer).
- numberOfPackages: Number of packages (integer).
- packages: Package details in the format <pkg_id> <pkg_weight> <pkg_distance> <offer_code> for each package.

#### Example

```
./courier_service calculateCost 100 3 "PKG1 5 5 OFR001" "PKG2 15 5 OFR002" "PKG3 10 100 OFR003"
```

#### Output

The command will output the calculated cost and any applicable discounts for each package.

```
Package PKG1
Base Delivery Cost: 100
Weight: 5 kg | Distance: 5 km
Offer code: OFR001
Discount: 0.00 (Offer not applicable as criteria not met)
Breakdown:
  Base Delivery Cost: 100.00
  Weight Cost: 50.00
  Distance Cost: 25.00
  Discount: -0.00
Total Delivery Cost: 175.00

Package PKG2
Base Delivery Cost: 100
Weight: 15 kg | Distance: 5 km
Offer code: OFR002
Discount: 0.00 (Offer not applicable as criteria not met)
Breakdown:
  Base Delivery Cost: 100.00
  Weight Cost: 150.00
  Distance Cost: 25.00
  Discount: -0.00
Total Delivery Cost: 275.00

Package PKG3
Base Delivery Cost: 100
Weight: 10 kg | Distance: 100 km
Offer code: OFR003
Discount: 35.00 (Discount of 5% applied)
Breakdown:
  Base Delivery Cost: 100.00
  Weight Cost: 100.00
  Distance Cost: 500.00
  Discount: -35.00
Total Delivery Cost: 665.00
```

### calculateTimeAndCost Command

This command calculates the delivery cost and estimated delivery time for packages.

#### Usage

```
./courier_service calculateTimeAndCost <baseDeliveryCost> <numberOfPackages> <packages> <number_of_vehicles> <max_speed> <max_carriable_weight>
```

#### Arguments

- baseDeliveryCost: Base delivery cost (integer).
- numberOfPackages: Number of packages (integer).
- packages: Package details in the format <pkg_id> <pkg_weight> <pkg_distance> <offer_code> for each package.
- number_of_vehicles: Number of vehicles available (integer).
- max_speed: Maximum speed of the vehicles (integer).
- max_carriable_weight: Maximum weight each vehicle can carry (integer).

#### Example

```
./courier_service calculateTimeAndCost 100 5 "PKG1 150 150 OFR001" "PKG2 75 125 OFR0008" "PKG3 175 100 OFR003" "PKG4 110 60 OFR002" "PKG5 155 95 NA" 2 70 200
```

#### Output

The command will output the delivery cost, discount, total cost, delivery time, and vehicle availability time for each package.

```

Package: PKG2
  Vehicle: 1
  Discount: 0.00
  Total Cost: 1475.00
  Delivery Time: 1.79 hours

Package: PKG4
  Vehicle: 1
  Discount: 0.00
  Total Cost: 1500.00
  Delivery Time: 0.86 hours

Package: PKG3
  Vehicle: 2
  Discount: 0.00
  Total Cost: 2350.00
  Delivery Time: 1.43 hours

Package: PKG5
  Vehicle: 2
  Discount: 0.00
  Total Cost: 2125.00
  Delivery Time: 4.21 hours

Package: PKG1
  Vehicle: 1
  Discount: 0.00
  Total Cost: 750.00
  Delivery Time: 4.00 hours
```

## Configuration

The offers and other configurations can be set in the configuration file. Make sure to update the **config/config.json** file with the relevant details.

## Error Cases

The CLI can generate various error cases. Here are some examples:

- Invalid base delivery cost
- Invalid number of packages
- Insufficient package information provided
- Invalid weight or distance for a package
- Invalid number of vehicles
- Invalid vehicle speed
- Invalid vehicle capacity
