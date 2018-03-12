package config_test

import (
	. "github.com/bit-mancer/go-util/config"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("config", func() {

	_ = Describe("ValidateConstraints", func() {
		It("returns an error if passed a nil value", func() {
			Expect(ValidateConstraints(nil)).To(HaveOccurred())
		})

		It("returns an error if passed a value that is not a struct or ptr-to-struct", func() {
			// Inexhaustive, just hit the major stuff
			Expect(ValidateConstraints(42)).To(HaveOccurred())
			Expect(ValidateConstraints(true)).To(HaveOccurred())
			Expect(ValidateConstraints([]int{})).To(HaveOccurred())
			Expect(ValidateConstraints([8]int{})).To(HaveOccurred())
			Expect(ValidateConstraints(make([]int, 8))).To(HaveOccurred())
		})

		It("does not return an error if passed an (empty) struct", func() {
			Expect(ValidateConstraints(struct{}{})).NotTo(HaveOccurred())
		})

		It("does not return an error if passed a ptr-to-an-(empty)-struct", func() {
			Expect(ValidateConstraints(&struct{}{})).NotTo(HaveOccurred())
		})

		It("does not return an error on unrecognized struct tags", func() {
			type NonConfigStructTagConfig struct {
				Value int `someOtherTag:"someValue"`
			}

			Expect(ValidateConstraints(NonConfigStructTagConfig{42})).NotTo(HaveOccurred())
		})

		It("returns an error on malformed config struct tags", func() {
			type MalformedConfigStructTagConfig struct {
				Value bool `config:"someUnrecognizedValue"`
			}

			Expect(ValidateConstraints(MalformedConfigStructTagConfig{true})).To(HaveOccurred())
		})

		It("supports the 'required' struct tag for supported types", func() {

			// Inexhaustive for the various int types etc. -- just hit the major stuff and focus more on unsupported types

			{
				type RequiredSupportedConfig struct {
					Value int `config:"required"`
				}

				Expect(ValidateConstraints(RequiredSupportedConfig{0})).To(HaveOccurred())
				Expect(ValidateConstraints(RequiredSupportedConfig{1})).NotTo(HaveOccurred())
			}

			{
				type RequiredSupportedConfig struct {
					Value uint `config:"required"`
				}

				Expect(ValidateConstraints(RequiredSupportedConfig{0})).To(HaveOccurred())
				Expect(ValidateConstraints(RequiredSupportedConfig{1})).NotTo(HaveOccurred())
			}

			{
				type RequiredSupportedConfig struct {
					Value float32 `config:"required"`
				}

				Expect(ValidateConstraints(RequiredSupportedConfig{0})).To(HaveOccurred())
				Expect(ValidateConstraints(RequiredSupportedConfig{1.0})).NotTo(HaveOccurred())
			}

			{
				type RequiredSupportedConfig struct {
					Value string `config:"required"`
				}

				Expect(ValidateConstraints(RequiredSupportedConfig{})).To(HaveOccurred())
				Expect(ValidateConstraints(RequiredSupportedConfig{"test"})).NotTo(HaveOccurred())
			}
		})

		It("returns an error when the 'required' struct tag is applied to unsupported types", func() {

			{
				type RequiredUnsupportedConfig struct {
					Value chan int `config:"required"`
				}

				Expect(ValidateConstraints(RequiredUnsupportedConfig{})).To(HaveOccurred())
				Expect(ValidateConstraints(RequiredUnsupportedConfig{make(chan int)})).To(HaveOccurred())
			}

			{
				type RequiredUnsupportedConfig struct {
					Value func() bool `config:"required"`
				}

				Expect(ValidateConstraints(RequiredUnsupportedConfig{})).To(HaveOccurred())
				Expect(ValidateConstraints(RequiredUnsupportedConfig{func() bool { return true }})).To(HaveOccurred())
			}

			{
				type RequiredUnsupportedConfig struct {
					Value interface{} `config:"required"`
				}

				Expect(ValidateConstraints(RequiredUnsupportedConfig{})).To(HaveOccurred())
				Expect(ValidateConstraints(RequiredUnsupportedConfig{"test"})).To(HaveOccurred())
			}

			{
				type RequiredUnsupportedConfig struct {
					Value []int `config:"required"`
				}

				Expect(ValidateConstraints(RequiredUnsupportedConfig{})).To(HaveOccurred())
				Expect(ValidateConstraints(RequiredUnsupportedConfig{[]int{1, 2}})).To(HaveOccurred())
			}

			{
				type RequiredUnsupportedConfig struct {
					Value [2]int `config:"required"`
				}

				Expect(ValidateConstraints(RequiredUnsupportedConfig{})).To(HaveOccurred())
				Expect(ValidateConstraints(RequiredUnsupportedConfig{[2]int{1, 2}})).To(HaveOccurred())
			}
		})

		It("supports structs with untagged unsupported types", func() {

			type RequiredSupportedConfig struct {
				Value            int `config:"required"`
				UnsupportedValue [2]int
			}

			Expect(ValidateConstraints(RequiredSupportedConfig{0, [2]int{}})).To(HaveOccurred())
			Expect(ValidateConstraints(RequiredSupportedConfig{0, [2]int{1, 2}})).To(HaveOccurred())
			Expect(ValidateConstraints(RequiredSupportedConfig{1, [2]int{}})).NotTo(HaveOccurred())
			Expect(ValidateConstraints(RequiredSupportedConfig{1, [2]int{1, 2}})).NotTo(HaveOccurred())
		})

		It("supports recursing into embedded structs", func() {
			{
				type EmbeddedConfig struct {
					Value int `config:"required"`
				}

				type RequiredSupportedConfig struct {
					Value string `config:"required"`

					EmbeddedValues EmbeddedConfig
				}

				Expect(ValidateConstraints(RequiredSupportedConfig{})).To(HaveOccurred())
				Expect(ValidateConstraints(RequiredSupportedConfig{"test", EmbeddedConfig{0}})).To(HaveOccurred())
				Expect(ValidateConstraints(RequiredSupportedConfig{"", EmbeddedConfig{4}})).To(HaveOccurred())
				Expect(ValidateConstraints(RequiredSupportedConfig{"test", EmbeddedConfig{4}})).NotTo(HaveOccurred())
			}
		})

		It("supports struct tags on embedded ptr-to-structs", func() {
			{
				type EmbeddedConfig struct {
					Value int
				}

				type RequiredSupportedConfig struct {
					Value string `config:"required"`

					EmbeddedValues *EmbeddedConfig `config:"required"`
				}

				Expect(ValidateConstraints(RequiredSupportedConfig{})).To(HaveOccurred())
				Expect(ValidateConstraints(RequiredSupportedConfig{"test", nil})).To(HaveOccurred())
				Expect(ValidateConstraints(RequiredSupportedConfig{"", &EmbeddedConfig{}})).To(HaveOccurred())
				Expect(ValidateConstraints(RequiredSupportedConfig{"test", &EmbeddedConfig{}})).NotTo(HaveOccurred())
			}
		})

		It("supports recursing into embedded ptr-to-structs", func() {
			{
				type EmbeddedConfig struct {
					Value int `config:"required"`
				}

				type RequiredSupportedConfig struct {
					Value string `config:"required"`

					EmbeddedValues *EmbeddedConfig `config:"required"`
				}

				Expect(ValidateConstraints(RequiredSupportedConfig{})).To(HaveOccurred())
				Expect(ValidateConstraints(RequiredSupportedConfig{"test", nil})).To(HaveOccurred())
				Expect(ValidateConstraints(RequiredSupportedConfig{"test", &EmbeddedConfig{}})).To(HaveOccurred())
				Expect(ValidateConstraints(RequiredSupportedConfig{"", &EmbeddedConfig{4}})).To(HaveOccurred())
				Expect(ValidateConstraints(RequiredSupportedConfig{"test", &EmbeddedConfig{4}})).NotTo(HaveOccurred())
			}
		})
	})
})
