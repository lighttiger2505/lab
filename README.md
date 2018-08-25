# lab

[![Build Status](https://travis-ci.org/lighttiger2505/lab.svg?branch=master)](https://travis-ci.org/lighttiger2505/lab)
[![Coverage Status](https://coveralls.io/repos/github/lighttiger2505/lab/badge.svg?branch=master)](https://coveralls.io/github/lighttiger2505/lab?branch=master)

lab is a cli client of gitlab like [hub](https://github.com/github/hub).

## Installation

### Go developer

Please getting source code and build.

```sh
go get github.com/lighttiger2505/lab
make ensure
go install
```

### Binary donwload

Please running install script.

```sh
curl -s https://raw.githubusercontent.com/lighttiger2505/lab/master/install.sh | bash
```

## Features

```
Usage: lab [--version] [--help] <command> [<args>]

Available commands are:
    browse           Browse repository page
    issue            Create and Edit, list a issue
    lint             validate .gitlab-ci.yml
    merge-request    Create and Edit, list a merge request
    pipeline         List pipeline, List pipeline jobs
    project          Show project
    user             Show pipeline
```

## Usage

1. change directory gitlab repository
	- lab command accesses gitlab quickly by useing repository infomation
	```sh
	$ cd {gitlab repository}
	```
1. laucn lab command
	```sh
	$ lab issue
	```
1. please input personal access token
	- use the lab command you need Personal access take look here(`https://docs.gitlab.com/ee/user/profile/personal_access_tokens.html#creating-a-personal-access-token`)
	```sh
	Please input GitLab private token : {your token}
	```

## Feature

### Browse

Open gitlab pages on brwoser.

```sh
# open project page
$ lab browse

# open issue list page
$ lab browse #
$ lab browse i

# open issue detail page
$ lab browse #10
$ lab browse i10
```

### List Issue and Merge Request

Show any list.

```sh
# List Issue
$ lab issue

# List Merge Request
$ lab merge-request
```

### Add Issue and Merge Request

title and description input on editor.

```sh
# Add Issue
$ lab add-issue

# Add Merge Request
$ lab add-merge-request --target={target branch}
```

## Configuration

auto create configuration file `~/.labconfig.yml` when launch lab command

### Sample

```yml
# personal access token
# store key/value style
tokens:
  gitlab.ssl.sample.jp: sampletoken
  gitlab.ssl.lowpriority.jp: lowprioritytoken
# Determine priority when there are multiple pieces of remote information in the repository
preferreddomains:
- gitlab.ssl.sample.jp
- gitlab.ssl.lowpriority.jp
```

## ToDos

- variable command
    - [x] Project-level Variables
    - [ ] Group-level Variables
- use template
    - [ ] issue template
    - [ ] merge request template
- [ ] pipeline actions
    - [ ] cancel
    - [ ] retry
- [ ] tags command
- [ ] project-member command
- workflow automation command
    - [ ] create
        - create new project and cloning repository
    - [ ] fork
        - create fork project and cloning repository
    - [ ] flow
        - create issue and create `WIP:` merge request
