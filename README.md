# UsageReport Plugin

_Note_: This is an experimental fork from [https://github.com/krujos/usagereport-plugin](https://github.com/krujos/usagereport-plugin) adding more capabilities when it comes to measure service instance usage.

This CF CLI Plugin to shows memory consumption and application instances, and service instances for each org and space you have permission to access.


# Usage

For listing current service instance usage in CSV style:

```
○ → cf usage-report -i summary -f csv
Org Name,Space Name,Service Instance Name,Service Instance Type,Service Name,Service Plan Name,Amount of Bound Apps,Bound Apps
DataFlow,Test,redis,managed_service_instance,p-redis,shared-vm,1,6ed59f50-dd09-4a28-ae17-2e4254a60f83
DataFlow,Test,rabbit,managed_service_instance,p-rabbitmq,standard,1,6ed59f50-dd09-4a28-ae17-2e4254a60f83
DataFlow,Test,my_mysql,managed_service_instance,p-mysql,100mb,1,6ed59f50-dd09-4a28-ae17-2e4254a60f83
AES,Dev,mysql,managed_service_instance,p-mysql,100mb,1,06ce0f19-0419-4b28-99a8-1cb48b973258
AES,Dev,edgeTest,managed_service_instance,apigee-edge,org,0,
```

For listing an app centric view of service instance usage:

```
○ → cf usage-report -i app -f csv
Org,Space,AppName,Instances,Bound Service Instances,Bound PCF Services,Bound User Provided Services,Bound 3rd Party Services
system,system,p-invitations,2,0,0,0,0
system,system,apps-manager-js,6,0,0,0,0
system,system,app-usage-server,1,0,0,0,0
system,system,app-usage-scheduler,1,0,0,0,0
system,system,app-usage-worker,1,0,0,0,0
system,notifications-with-ui,notifications-ui,2,0,0,0,0
system,pivotal-account-space,pivotal-account,2,0,0,0,0
system,autoscaling,autoscale,3,0,0,0,0
apigee-cf-service-broker-org,apigee-cf-service-broker-space,apigee-cf-service-broker-2.0.1,1,0,0,0,0
DataFlow,Test,dataflow-server,1,3,3,0,0
AES,Dev,aes,1,1,0,1,0
AES,Dev,aesserver,1,0,0,0,0
ESA,Dev,aes,1,0,0,0,0
ESA,Dev,spring-music,1,2,1,1,0
```


For human readable output:

```
➜  usagereport-plugin git:(master) ✗ cf usage-report
Gathering usage information
Org platform-eng is consuming 53400 MB of 204800 MB.
	Space CFbook is consuming 128 MB memory (0%) of org quota.
		1 apps: 1 running 0 stopped
		1 instances: 1 running, 0 stopped
Org krujos is consuming 512 MB of 10240 MB.
	Space development is consuming 0 MB memory (0%) of org quota.
		4 apps: 0 running 4 stopped
		4 instances: 0 running, 4 stopped
	Space production is consuming 512 MB memory (5%) of org quota.
		1 apps: 1 running 0 stopped
		2 instances: 2 running, 0 stopped
Org pcfp is consuming 7296 MB of 102400 MB.
	Space development is consuming 0 MB memory (0%) of org quota.
		0 apps: 0 running 0 stopped
		0 instances: 0 running, 0 stopped
	Space docs-staging is consuming 512 MB memory (0%) of org quota.
		2 apps: 1 running 1 stopped
		4 instances: 2 running, 2 stopped
	Space docs-prod is consuming 512 MB memory (0%) of org quota.
		3 apps: 1 running 2 stopped
		5 instances: 2 running, 3 stopped
	Space guillermo-playground is consuming 2560 MB memory (2%) of org quota.
		1 apps: 1 running 0 stopped
		5 instances: 5 running, 0 stopped
	Space haydon-playground is consuming 1024 MB memory (1%) of org quota.
		1 apps: 1 running 0 stopped
		1 instances: 1 running, 0 stopped
	Space jkruck-playground is consuming 128 MB memory (0%) of org quota.
		1 apps: 1 running 0 stopped
		1 instances: 1 running, 0 stopped
	Space rsalas-dev is consuming 0 MB memory (0%) of org quota.
		0 apps: 0 running 0 stopped
		0 instances: 0 running, 0 stopped
	Space shekel-dev is consuming 1536 MB memory (1%) of org quota.
		3 apps: 3 running 0 stopped
		3 instances: 3 running, 0 stopped
	Space shekel-qa is consuming 0 MB memory (0%) of org quota.
		0 apps: 0 running 0 stopped
		0 instances: 0 running, 0 stopped
	Space hd-playground is consuming 0 MB memory (0%) of org quota.
		0 apps: 0 running 0 stopped
		0 instances: 0 running, 0 stopped
	Space dwallraff-dev is consuming 1024 MB memory (1%) of org quota.
		1 apps: 1 running 0 stopped
		1 instances: 1 running, 0 stopped
You are running 18 apps in 3 orgs, with a total of 27 instances.
```

CSV output:

```
➜  usagereport-plugin git:(master) ✗ cf usage-report -f csv
OrgName, SpaceName, SpaceMemoryUsed, OrgMemoryQuota, AppsDeployed, AppsRunning, AppInstancesDeployed, AppInstancesRunning
test-org, test-space, 256, 4096, 2, 1, 3, 2
```

##Installation
#####Install from CLI
  ```
  $ cf add-plugin-repo CF-Community http://plugins.cloudfoundry.org/
  $ cf install-plugin 'Usage Report' -r CF-Community
  ```


#####Install from Source (need to have [Go](http://golang.org/dl/) installed)
  ```
  $ go get github.com/cloudfoundry/cli
  $ go get github.com/krujos/usagereport-plugin
  $ cd $GOPATH/src/github.com/krujos/usagereport-plugin
  $ go build
  $ cf install-plugin usagereport-plugin
  ```
