package schema

type User struct {
	Id         uint
	Characters []uint
	GroupsIds  []uint  `mgo:"groups_ids"`
	Groups     []Group `mgo:"groups"`
}
