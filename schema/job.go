package schema

type Job struct {
	Id                   uint    `json:"id"`
	InstallerId          uint    `json:"installer_id"`
	FacilityId           uint    `json:"facility_id"`
	StationId            uint    `json:"station_id"`
	ActivityId           uint    `json:"activity_id"`
	BlueprintId          uint    `json:"blueprint_id"`
	BlueprintTypeId      uint    `json:"blueprint_type_id"`
	BlueprintLocationId  uint    `json:"blueprint_location_id"`
	OutputLocationId     uint    `json:"output_location_id"`
	Runs                 uint    `json:"runs"`
	Cost                 float32 `json:"cost"`
	LicensedRuns         uint    `json:"licensed_runs"`
	Probability          float64 `json:"probability"`
	ProductTypeId        uint    `json:"product_type_id"`
	Status               string  `json:"status"`
	Duration             uint    `json:"duration"`
	StartDate            string  `json:"start_date"`
	EndDate              string  `json:"end_date"`
	PauseDate            string  `json:"pause_date"`
	CompletedDate        string  `json:"completed_date"`
	CompletedCharacterId uint    `json:"completed_character_id"`
	SuccessfulRuns       uint    `json:"successful_runs"`
}

type Jobs []Job
