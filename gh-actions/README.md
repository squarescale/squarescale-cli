# gh-actions

This image is pushed to Squarescale docker hub through the tag name `gh-actions`, and can be used in a Github workflow:

```
  job1:
    runs-on: ubuntu-20.04
    steps:
      - name: Schedule a Web service
        uses: docker://squarescale/cli:gh-actions
        env: 
          [...]
```

## Environment variables

### Required variables

| Name | Description | Type |
| ---- | ----------- | ---- |
| SQSC_TOKEN | The API key to access Squarescale via CLI  | string
| DOCKER_REPOSITORY | The docker hub repository image name of your application | string
| PROJECT_NAME | The name of the project you want to create. | string
| IAAS_PROVIDER | The provider to the IAAS you want to deploy your infrastructure. | string
| IAAS_REGION | The IAAS region. | string
| IAAS_CRED | The IAAS credential. | string
| NODE_TYPE | The node size. | string

### Optional variables

| Name | Description | Type |
| ---- | ----------- | ---- |
| DOCKER_REPOSITORY_TAG | The docker hub repository image name tag of your application | string
| DOCKER_USER | A docker hub username | string
| DOCKER_TOKEN | A docker hub password | string
| ORGANIZATION_NAME | The organization name you belong to. | string

#### Database

If one of the variables below is absent, no database will be created.

| Name | Description | Type |
| ---- | ----------- | ---- |
| DB_ENGINE | The database engine (e.g: "postgres") | string
| DB_ENGINE_VERSION | The database engine version (e.g: "12") | string
| DB_SIZE | The database size (e.g: "small") | string

#### Services

| Name | Description | Type |
| ---- | ----------- | ---- |
| SERVICES | All the services that will be schedule on the infrastructure | json

A service within this json has for key the name of the service and for value a json with some other variables :

| Name | Description | Type |
| ---- | ----------- | ---- |
| run_cmd | The run command that will be executed when the service is scheduling. | string
| network_rules | The network rules (only http for now).<br>With `name` (default "http") and `internal_port` (default "80"). | json
| env | The environment variables the application image needs to.  | json

:information_source: For database environment variable in the json structure, as database will be created before the services, its environment variables can be recover with `{{ GLOBAL_DB_VARIABLE_FROM_INFRA }}` (see example below).

Example:

```json
SERVICES: >-
  {
    "web": {
      "run_cmd": "bundle exec rails server -b 0.0.0.0",
      "network_rules": {
        "name": "http",
        "internal_port": "3000"
      },
      "env": {
        "RAILS_LOG_TO_STDOUT": "true",
        "DATABASE_HOST": "{{DB_HOST}}",
        "DATABASE_USERNAME": "{{DB_USERNAME}}",
        "DATABASE_NAME": "{{DB_NAME}}",
        "DATABASE_PASSWORD": "{{DB_PASSWORD}}"
      }
    }
  }
```

Here, a service named `web` will be schedule on Squarescale: 
- with a `Run command`
- the `3000` port will be routed on prefix `/` (by default and not configurable yet)
- 5 environements variables will be inserted into the service for logs and database. `{{ DB_* }}` corresponds to the global database environement variable created when the database was created before this services.

#### Batches

| Name | Description | Type |
| ---- | ----------- | ---- |
| BATCHES | All the batches that will be created on the infrastructure | json

A batch within this json has for key the name of the batch and for value a json with some other variables :

| Name | Description | Type |
| ---- | ----------- | ---- |
| execute | Either `true` or `false` if you want to execute or not the batch (default: `false`) | bool
| image_name | Image name (e.g: `bash`, default: `DOCKER_REPOSITORY`:`DOCKER_REPOSITORY_TAG` or `DOCKER_REPOSITORY`) | string
| is_private | Either `true` or `false` if you want to use a private image | string
| image_user | The image user. Only needed with `is_private`. | string
| image_password | The image password. Only needed with `is_private`. | string
| run_cmd | The run command that will be executed when the batch is executed. | string
| periodic | Enable a periodic batch. The json contains `periodicity` (default: `* * * * *`) and `timezone` (default: `Europe/Paris`) | json
| env | The environment variables the application image needs to. (see above with `SERVICES`)  | json

Example:

```json
BATCHES: >-
  {
    "database-setup": {
      "run_cmd": "bundle exec rails db:setup",
      "env": {
        "RAILS_LOG_TO_STDOUT": "true",
        "DATABASE_HOST": "{{DB_HOST}}",
        "DATABASE_USERNAME": "{{DB_USERNAME}}",
        "DATABASE_NAME": "{{DB_NAME}}",
        "DATABASE_PASSWORD": "{{DB_PASSWORD}}"
      }
    },
    "database-seed": {
      "run_cmd": "bundle exec rails db:sogilis:seed",
      "env": {
        "RAILS_LOG_TO_STDOUT": "true",
        "DATABASE_HOST": "{{DB_HOST}}",
        "DATABASE_USERNAME": "{{DB_USERNAME}}",
        "DATABASE_NAME": "{{DB_NAME}}",
        "DATABASE_PASSWORD": "{{DB_PASSWORD}}"
      }
    }
  }
```