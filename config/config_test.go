package config

import (
	"os"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestConfig(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Config Suite")
}

var _ = Describe("Config", func() {
	var (
		configPath string
		cfg        Config
	)

	BeforeEach(func() {
		configPath = "config_test.json"
		cfg = NewConfig()
	})

	AfterEach(func() {
		os.Remove(configPath)
	})

	Context("LoadConfig", func() {
		It("should successfully load and validate the config", func() {
			configContent := `{
				"offers": [
					{
						"code": "OFFER1",
						"discount": 10,
						"minDistance": 0,
						"maxDistance": 100,
						"minWeight": 0,
						"maxWeight": 10
					}
				],
				"distanceCostPerKM": 5,
				"weightCostPerKG": 10
			}`
			err := os.WriteFile(configPath, []byte(configContent), 0644)
			Expect(err).ToNot(HaveOccurred())

			err = cfg.LoadConfig(configPath)
			Expect(err).ToNot(HaveOccurred())
		})

		It("should return an error for invalid config path", func() {
			invalidConfigPath := "invalid_path.json"
			err := cfg.LoadConfig(invalidConfigPath)
			Expect(err).To(HaveOccurred())
		})

		It("should return an error for invalid config structure while unmarhsaling", func() {
			invalidConfigContent := `{
				"offers": [
					{
						"code": "OFFER1",
						"discount": 10,
						"minDistance": 0,
						"maxDistance": 100,
						"minWeight": 0,
						"maxWeight": 10,this is extra
					}
				],
				"distanceCostPerKM": 5
			}`
			err := os.WriteFile(configPath, []byte(invalidConfigContent), 0644)
			Expect(err).ToNot(HaveOccurred())

			err = cfg.LoadConfig(configPath)
			Expect(err).To(HaveOccurred())
		})

		It("should return an error for invalid config structure while validating", func() {
			invalidConfigContent := `{
				"offers": [
					{
						"code": "OFFER1",
						"discount": 10
					}
				],
				"distanceCostPerKM": 5
			}`
			err := os.WriteFile(configPath, []byte(invalidConfigContent), 0644)
			Expect(err).ToNot(HaveOccurred())

			err = cfg.LoadConfig(configPath)
			Expect(err).To(HaveOccurred())
		})
	})

	Context("GetOffers", func() {
		It("should return the correct offers from the config", func() {
			configContent := `{
				"offers": [
					{
						"code": "OFFER1",
						"discount": 10,
						"minDistance": 0,
						"maxDistance": 100,
						"minWeight": 0,
						"maxWeight": 10
					}
				],
				"distanceCostPerKM": 5,
				"weightCostPerKG": 10
			}`
			err := os.WriteFile(configPath, []byte(configContent), 0644)
			Expect(err).ToNot(HaveOccurred())

			err = cfg.LoadConfig(configPath)
			Expect(err).ToNot(HaveOccurred())

			offers := GetOffers()

			Expect(offers).To(HaveLen(1))
			Expect(offers[0].Code).To(Equal("OFFER1"))
		})

	})

	Context("GetWeightCostPerKG", func() {
		It("should return the correct weight cost per KG", func() {
			configContent := `{
				"offers": [
					{
						"code": "OFFER1",
						"discount": 10,
						"minDistance": 0,
						"maxDistance": 100,
						"minWeight": 0,
						"maxWeight": 10
					}
				],
				"distanceCostPerKM": 5,
				"weightCostPerKG": 10
			}`
			err := os.WriteFile(configPath, []byte(configContent), 0644)
			Expect(err).ToNot(HaveOccurred())

			err = cfg.LoadConfig(configPath)
			Expect(err).ToNot(HaveOccurred())

			weightCost := GetWeightCostPerKG()
			Expect(weightCost).To(Equal(10))
		})
	})

	Context("GetDistanceCostPerKM", func() {
		It("should return the correct distance cost per KM", func() {
			configContent := `{
				"offers": [
					{
						"code": "OFFER1",
						"discount": 10,
						"minDistance": 0,
						"maxDistance": 100,
						"minWeight": 0,
						"maxWeight": 10
					}
				],
				"distanceCostPerKM": 5,
				"weightCostPerKG": 10
			}`
			err := os.WriteFile(configPath, []byte(configContent), 0644)
			Expect(err).ToNot(HaveOccurred())

			err = cfg.LoadConfig(configPath)
			Expect(err).ToNot(HaveOccurred())

			distanceCost := GetDistanceCostPerKM()
			Expect(distanceCost).To(Equal(5))
		})
	})
})
