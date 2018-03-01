package schema

type Role struct {
	Id   uint
	Name string
}

var ScopeSets = map[string]string{
	"simple":     "publicData",
	"industrial": "characterIndustryJobsRead corporationIndustryJobsRead",
	"mailing":    "esi-mail.send_mail.v1",
}

var ScopeSetsReversed = func() map[string]string {
	m := map[string]string{}
	for k, v := range ScopeSets {
		m[v] = k
	}
	return m
}()