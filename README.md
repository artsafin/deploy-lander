# Deploy Lander

Tool for performing semi-automatic deploys:
1) CI creates a build and notifies Deploy Lander agent running on a deploy target that build completed
*At this point build has been created and ready for deploy but not yet deployed*
2) Deploy Lander accepts an HTTP request with instructions on which registered builds should be deployed
*Deploy Lander deploys the requested build*

## Configuration

All configuration is done through environment variables.

```
deploylander.projectname.path=/home/user/osdmon     | Required. Root of the project with git repository
deploylander.projectname.tags=.*                    | Required unless `branches` is specified. Regexp of tag names that will be accepted as completed build
deploylander.projectname.branches=.*                | Required unless `tags` is specified. Regexp of branch names that will be accepted as completed build
deploylander.projectname.deploy_command=make        | Required. Command to execute in `path` to deploy
```

Structure of variable key:

- `deploylander` is common prefix
- after prefix comes the alias of the project
- then goes configuration parameter name

## Usage



## API

### Register completed build

Register tag:
```
curl -d 'tag=1.2' http://host/build/projectname
```

Register branch (note json usage):
```
curl -d '{"branch": "feature-branch"}' -H 'Content-Type: application-json' http://host/build/projectname
```

### Deploy a build

Deploy previously registered tag:
```
curl -d 'tag=1.2' http://host/deploy/projectname
```

Deploy previously registered branch (note json usage):
```
curl -d '{"branch": "feature-branch"}' -H 'Content-Type: application-json' http://host/deploy/projectname
```
