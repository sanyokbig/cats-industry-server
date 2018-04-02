package schema

import (
	"fmt"

	"github.com/go-errors/errors"
	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"
)

//easyjson:json
type Job struct {
	ID                   uint     `json:"id"`
	EveID                uint     `json:"job_id"`
	InstallerID          uint     `json:"installer_id"`
	FacilityID           uint64   `json:"facility_id"`
	StationID            uint64   `json:"station_id"`
	ActivityID           uint     `json:"activity_id"`
	BlueprintID          uint64   `json:"blueprint_id"`
	BlueprintTypeID      uint     `json:"blueprint_type_id"`
	BlueprintLocationID  uint64   `json:"blueprint_location_id"`
	OutputLocationID     uint64   `json:"output_location_id"`
	Runs                 uint     `json:"runs"`
	Cost                 float32  `json:"cost"`
	LicensedRuns         uint     `json:"licensed_runs"`
	Probability          float64  `json:"probability"`
	ProductTypeID        uint     `json:"product_type_id"`
	Status               string   `json:"status"`
	Duration             uint     `json:"duration"`
	StartDate            UnixTime `json:"start_date"`
	EndDate              UnixTime `json:"end_date"`
	PauseDate            UnixTime `json:"pause_date"`
	CompletedDate        UnixTime `json:"completed_date"`
	CompletedCharacterID uint     `json:"completed_character_id"`
	SuccessfulRuns       uint     `json:"successful_runs"`
}

//easyjson:json
type Jobs []Job

func (jobs *Jobs) Save(queryer sqlx.Queryer) error {
	if len(*jobs) == 0 {
		log.Debug("saving empty jobs slice")
		return nil
	}

	values := ``
	for _, j := range *jobs {
		values += fmt.Sprintf(`(
				%v, %v, %v, %v,
				%v, %v, %v, %v,
				%v, %v, %v, %v,
				%v, %v, '%v', %v,
				%v, %v, %v, %v,
				%v, %v
			), `,
			j.EveID, j.InstallerID, j.FacilityID, j.StationID,
			j.ActivityID, j.BlueprintID, j.BlueprintTypeID, j.BlueprintLocationID,
			j.OutputLocationID, j.Runs, j.Cost, j.LicensedRuns,
			j.Probability, j.ProductTypeID, j.Status, j.Duration,
			j.StartDate, j.EndDate, j.PauseDate, j.CompletedDate,
			j.CompletedCharacterID, j.SuccessfulRuns,
		)
	}
	values = values[:len(values)-2]

	query := fmt.Sprintf(` 
 		INSERT 
 		INTO jobs 
 		(
 			eve_id, 				installer_id, 		facility_id, 		station_id, 
 			activity_id, 			blueprint_id, 		blueprint_type_id, 	blueprint_location_id, 
 			output_location_id, 	runs, 				cost, 				licensed_runs, 
 			probability, 			product_type_id,	status, 			duration, 
 			start_date, 			end_date, 			pause_date, 		completed_date, 
 			completed_character_id, successful_runs
 		) 
 		VALUES %v ON CONFLICT (eve_id) DO UPDATE SET
			installer_id = EXCLUDED.installer_id,
			facility_id = EXCLUDED.facility_id,
			station_id = EXCLUDED.station_id,
			
			activity_id = EXCLUDED.activity_id,
			blueprint_id = EXCLUDED.blueprint_id,
			blueprint_type_id = EXCLUDED.blueprint_type_id,
			blueprint_location_id = EXCLUDED.blueprint_location_id,			

			output_location_id = EXCLUDED.output_location_id,
			runs = EXCLUDED.runs,
			cost = EXCLUDED.cost,
			licensed_runs = EXCLUDED.licensed_runs,

			probability = EXCLUDED.probability,
			product_type_id = EXCLUDED.product_type_id,
			status = EXCLUDED.status,
			duration = EXCLUDED.duration,

			start_date = EXCLUDED.start_date,
			end_date = EXCLUDED.end_date,
			pause_date = EXCLUDED.pause_date,
			completed_date = EXCLUDED.completed_date,

			completed_character_id = EXCLUDED.completed_character_id,
			successful_runs = EXCLUDED.successful_runs
		RETURNING id
	`, values)

	rows, err := queryer.Queryx(query)
	if err != nil {
		return errors.New(err)
	}
	defer rows.Close()

	// Set job ids
	i, id := 0, 0
	for rows.Next() {
		err = rows.Scan(&id)
		if err != nil {
			log.Warningf("failed to scan new job id: %v", err)
			id++
			continue
		}
		(*jobs)[i].ID = uint(id)
		i++
	}

	return nil
}
