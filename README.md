# Overview
This is a sample repository to showcase the use of semantic-release tool to automatically tag a
Terraform module on pushes to main/master branch.  

The pipeline includes validation during pull requests:
* Static Code Analysis with Checkov
* Unit Testing with Terratest

The pipeline will execute the versioning, tagging, and release only on the main branch. This could
be expanded to include the master branch by updating the `pr` and `triggers` and updating the condition on the `tag_module` stage.

Note: to reduce the retrieval of external scripts and piping directly to bash for execution, the script to install Node JS on Ubuntu has been persisted in this sample repository.

# Requirements
## Pipeline
The sample pipeline defined in this repository assumes the following to be in place:
* Firewall Settings opened to allow the software installation. See Firewall Requirements Section
* Variable Group with the following variables defined:
  * `GH_TOKEN` - GitHub PAT token with the following permissions enabled: `repo` (all settings).
    semantic-release needs the ability to read tags and commits, and write commits for the Change
    Log.
  * `TERRAFORM_VALIDATION_POOL` - Name of the Self-Hosted agent pool with the requisite tools
* Secure File named `id_rsa` containing SSH key to retrieve modules from GitHub over SSH.

## Firewall
If your Azure DevOps agent is self-hosted and behind a firewall, the following hosts need to be opened over HTTPS:
* deb.nodesource.com - To install nodejs apt package
* registry.npmjs.org - To install semantic-release and plugins

## Software
The pipeline defined in this repo uses the following open source software.
* [checkov](https://www.checkov.io/) - Static Code Analysis tool
* [semantic-release](https://semantic-release.gitbook.io/semantic-release/) - Nodejs application
  automatic the tagging, versioning, and change log creation. This application is extensible and
  plugin based.
* Golang - Go is used to run unit testing of the module. For this sample repository, no actual
  assertions are made. The intent is to show the step being called rather than demonstrating how to
  perform unit tests.

# Semantic-Release Configuration
The configuration `.releaserc` contains a JSON configuration for semantic-rules to execute the versioning, tagging, and publishing.

This configuration defines the order of the plugins used and their respective configurations. Below are explanations of each configuration
## commit-analyzer
This configuration enforces the use of the [Angular commit convention](https://github.com/angular/angular.js/blob/master/DEVELOPERS.md#-git-commit-guidelines). It also overrides the release for the following types: docs, refactor, test, and style. This means that commits such as `docs(myfeature): correct usage documentation` will generate a patch increment (x.x.1). Without this configuration, the pipeline would not increment the version. 

Further, `parserOpts` include a few adaptations of keywords that will trigger a MAJOR version
increase. These keywords have to be present in the body of the commit message (does not trigger if
included in the title of the commit).

## release-notes-generator
This plugin generates the changelog content using the Angular notation. It has similar `parserOpts` as the `commit-analyzer`

## changelog
Creates/updates the CHANGELOG.md

## git
Commits changes to CHANGELOG.md to the repository

## github
Publishes tags to GitHub and creates a release. Adds the CHANGELOG.md to the release.

# Limitations
N/A

# Known Issues

If you do not use one of the commit messages that triggers a version bump of any kind,
semantic-release will generate a tag value that exists already. As a result, the pipeline will fail
with the following error:

As a best practice, any merge to main branch should include commits following the angular
conventions and should trigger a patch version bump.
```
2021-02-10T17:01:13.8549292Z [5:01:13 PM] [semantic-release] â€º âœ–  An error occurred while running semantic-release: Error: Command failed with exit code 128: git tag 3.1.0 eef4bec231cc305eff3b59bbbfca5606ebf367a6
2021-02-10T17:01:13.8551050Z fatal: tag '3.1.0' already exists
2021-02-10T17:01:13.8552329Z     at makeError (/home/azureadmin/devops-agent/agent-2.181.1/_work/node_modules/execa/lib/error.js:59:11)
2021-02-10T17:01:13.8553784Z     at handlePromise (/home/azureadmin/devops-agent/agent-2.181.1/_work/node_modules/execa/index.js:114:26)
2021-02-10T17:01:13.8554852Z     at processTicksAndRejections (node:internal/process/task_queues:94:5)
2021-02-10T17:01:13.8557945Z     at async tag (/home/azureadmin/devops-agent/agent-2.181.1/_work/node_modules/semantic-release/lib/git.js:224:3)
2021-02-10T17:01:13.8559666Z     at async run (/home/azureadmin/devops-agent/agent-2.181.1/_work/node_modules/semantic-release/index.js:191:5)
2021-02-10T17:01:13.8561178Z     at async module.exports (/home/azureadmin/devops-agent/agent-2.181.1/_work/node_modules/semantic-release/index.js:260:22)
2021-02-10T17:01:13.8563179Z     at async module.exports (/home/azureadmin/devops-agent/agent-2.181.1/_work/node_modules/semantic-release/cli.js:55:5) {
2021-02-10T17:01:13.8565240Z   shortMessage: 'Command failed with exit code 128: git tag 3.1.0 eef4bec231cc305eff3b59bbbfca5606ebf367a6',
2021-02-10T17:01:13.8566802Z   command: 'git tag 3.1.0 eef4bec231cc305eff3b59bbbfca5606ebf367a6',
2021-02-10T17:01:13.8567979Z   exitCode: 128,
2021-02-10T17:01:13.8568803Z   signal: undefined,
2021-02-10T17:01:13.8569531Z   signalDescription: undefined,
2021-02-10T17:01:13.8570644Z   stdout: '',
2021-02-10T17:01:13.8571798Z   stderr: "fatal: tag '3.1.0' already exists",
2021-02-10T17:01:13.8572753Z   failed: true,
2021-02-10T17:01:13.8573598Z   timedOut: false,
2021-02-10T17:01:13.8574472Z   isCanceled: false,
2021-02-10T17:01:13.8576255Z   killed: false
2021-02-10T17:01:13.8577681Z }
```
