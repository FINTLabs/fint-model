package types

type AssociationQuery struct {
	XPath string
	Role  AssociationRole
}

type AssociationRole string

const (
	RoleSource AssociationRole = "source"
	RoleTarget AssociationRole = "target"
)
