<a name="unreleased"></a>
## [Unreleased]


<a name="v1.1.5"></a>
## [v1.1.5] - 2023-08-17
### Release
- v1.1.5


### Bug Fixes
- service set help command is wrong
- service show uses -name instead of -service for service name (like other service commands)
- Uniformize error messages
- Add more insights to 504 Gateway timeout errors
- Set proper defaults for container/service CPU and memory limits
- service show does not report container mounted volumes
- service add and set commands should allow all parameters like in frontend UI
- Add missing project cluster node size to create parameters
- Add observability integrated service to project creation
- Add ElasticSearch integrated service to project creation


### Features
- Add max-client-disconnect option
- it would be easier to also specify environment variables via command line instead of JSON file for service set
- Add external ElasticSearch option to project
- Replace stateful nodes to extra nodes
- Uniformization of outputs with frontend UI
- Add database backup parameters option


[Unreleased]: https://github.com/squarescale/squarescale-cli/compare/v1.1.5...HEAD
[v1.1.5]: https://github.com/squarescale/squarescale-cli/compare/v1.1.4...v1.1.5
