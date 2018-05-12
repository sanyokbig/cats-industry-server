package foreman

import "github.com/sanyokbig/cats-industry-server/schema"

func (f *Foreman) GetJobs(params schema.GetParams) (*schema.Jobs, error) {
	query := `
		select 
 			j.id, 
 			eve_id, 
 			installer_id, 
 			facility_id, 
 			station_id, 
 			activity_id, 
 			blueprint_id, 
 			blueprint_type_id, 
 			blueprint_location_id, 
 			output_location_id, 
 			runs, 
 			cost,
 		 	licensed_runs, 
 		 	probability, 
 		 	product_type_id, 
 		 	status, 
 		 	duration,
 		 	start_date, 
 		 	end_date, 
 		 	pause_date, 
 		 	completed_date, 
 		 	completed_character_id, 
 		 	successful_runs, 
 		 	pt.name product_name,
		    a.name activity_name
 		from jobs j
		left join product_types pt on j.product_type_id = pt.id
		left join ram_activities a on j.activity_id = a.id
		where status != 'delivered'
		order by end_date asc
`

	rows, err := f.db.Queryx(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	jobs, job := schema.Jobs{}, schema.Job{}
	for rows.Next() {
		err = rows.StructScan(&job)
		if err != nil {
			return nil, err
		}
		jobs = append(jobs, job)
	}

	return &jobs, nil
}
