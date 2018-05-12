package sdeParser

type ProductTypes map[string]ProductType

type ProductType struct {
	Name struct {
		En string `yaml:"en"`
	} `yaml:"name"`
}
