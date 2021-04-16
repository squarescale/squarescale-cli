package main

const (
	sqscToken           = "SQSC_TOKEN"
	dockerRepository    = "DOCKER_REPOSITORY"
	dockerRepositoryTag = "DOCKER_REPOSITORY_TAG"
	organizationName    = "ORGANIZATION_NAME"
	projectName         = "PROJECT_NAME"
	iaasProvider        = "IAAS_PROVIDER"
	iaasRegion          = "IAAS_REGION"
	iaasCred            = "IAAS_CRED"
	monitoring          = "MONITORING"
	infraType           = "INFRA_TYPE"
	nodeType            = "NODE_TYPE"
	dbEngine            = "DB_ENGINE"
	dbEngineVersion     = "DB_ENGINE_VERSION"
	dbSize              = "DB_SIZE"
	servicesEnv         = "SERVICES"
	batchesEnv          = "BATCHES"
)

func main() {
	checkEnvironmentVariablesExists()

	project := Project{}
	project.create()

	database := Database{}
	database.create()

	services := Services{}
	services.create()

	batches := Batches{}
	batches.create()
}
