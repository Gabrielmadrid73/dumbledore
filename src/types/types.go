package types

type TypeSecrets struct {
	Secrets []string `form:"secrets" json:"secrets" binding:"required"`
}

type TypeSecret struct {
	Secret string
}

type TypeNamespace struct {
	Namespace string `form:"namespace" json:"namespace" binding:"required"`
}

type StructSecrets struct {
	*TypeNamespace
	*TypeSecrets
}

type StructSecret struct {
	*TypeNamespace
	*TypeSecret
}
