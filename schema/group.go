package schema

type Group struct {
	Id       uint   `mgo:"id"`
	Name     string `mgo:"name"`
	RolesIds []uint `mgo:"roles_ids"`
	Roles    []Role `mgo:"roles"`
}
