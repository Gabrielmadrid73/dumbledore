package types

type TypeSecrets struct {
	Secrets []string `form:"secrets" json:"secrets" binding:"required"`
}

type TypeNamespace struct {
	Namespace string `form:"namespace" json:"namespace" binding:"required"`
}

type StructSecrets struct {
	*TypeNamespace
	*TypeSecrets
}

type SecretAnnotations struct {
	ParamName string
	ParamType string
}
