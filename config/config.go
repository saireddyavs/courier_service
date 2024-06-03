package config

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

type Config interface {
	LoadConfig(configPath string) error
}

type Offer struct {
	Code        string  `mapstructure:"code" json:"code" validate:"required"`
	Discount    float64 `mapstructure:"discount" json:"discount" validate:"required"`
	MinDistance int     `mapstructure:"minDistance" json:"minDistance" validate:"required"`
	MaxDistance int     `mapstructure:"maxDistance" json:"maxDistance" validate:"required"`
	MinWeight   int     `mapstructure:"minWeight" json:"minWeight" validate:"required"`
	MaxWeight   int     `mapstructure:"maxWeight" json:"maxWeight" validate:"required"`
}

type config struct {
	Offers            []Offer `mapstructure:"offers" json:"offers" validate:"required"`
	DistanceCostPerKM int     `mapstructure:"distanceCostPerKM" json:"distanceCostPerKM" validate:"required"`
	WeightCostPerKG   int     `mapstructure:"weightCostPerKG" json:"weightCostPerKG" validate:"required"`
}

func NewConfig() Config {
	return &config{}
}

func (c config) LoadConfig(configPath string) error {
	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.SetConfigFile(configPath)

	err := viper.ReadInConfig()
	if err != nil {
		fmt.Printf("Error reading config file: %s\n", err)
		return err
	}

	err = viper.Unmarshal(&c)
	if err != nil {
		fmt.Printf("Error unmarshaling config: %s\n", err)
		return err
	}

	validate := validator.New(validator.WithRequiredStructEnabled())

	validationErr := validate.Struct(c)

	if validationErr != nil {
		fmt.Printf("Error while validating the configs: %s\n", validationErr)
		return validationErr
	}

	fmt.Println("Config validated and loaded successfully")

	return nil
}

func GetOffers() []Offer {
	var offers []Offer
	viper.UnmarshalKey("offers", &offers)
	return offers
}

func GetWeightCostPerKG() int {
	return viper.GetInt("weightCostPerKG")
}

func GetDistanceCostPerKM() int {
	return viper.GetInt("distanceCostPerKM")
}
