package sdeParser

type Activities []Activity

type Activity struct {
	ID          int64  `yaml:"activityID"`
	Name        string `yaml:"activityName"`
	Description string `yaml:"description"`
	Icon        string `yaml:"iconNo"`
}
