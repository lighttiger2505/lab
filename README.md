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
make install
```

### Binary download

Please run the install script:

```sh
curl -s https://raw.githubusercontent.com/lighttiger2505/lab/master/install.sh | bash
```

The script installs the `lab` command in `/usr/local/bin`. For more details, see the `install.sh` [source code](install.sh).

## Features

```
Usage: lab [--version] [--help] <command> [<args>]

Available commands are:
    browse                    Browse project page
    issue                     Create and Edit, list a issue
    issue-template            List issue template
    job                       List job
    lint                      validate .gitlab-ci.yml
    merge-request             Create and Edit, list a merge request
    merge-request-template    List merge request template
    mr                        Create and Edit, list a merge request
    pipeline                  List pipeline, List pipeline jobs
    project                   List project
    project-variable          List project level variables
    runner                    List CI/CD Runner
    user                      List user
```

## Usage

1. change directory gitlab repository
	- lab command accesses gitlab quickly by useing repository infomation
	```sh
	$ cd {gitlab repository}
	```
1. launch lab command
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
# Browse project page
$ lab browse

# Browse project file
$ lab browse ./README.md

# Browse sub page
$ lab browse -s issues
```

### Operations to Issue and Merge Request

Many operations can be done with simple input.

```sh
# List issue
lab issue

# Browse issue
lab issue -b {issue ii}

# Show issue
lab issue {issue id}

# Create issue
lab issue -e

# Update issue
lab issue {issue id} -e
```

## Configuration

auto create configuration file `~/.config/lab/config.yml` when launch lab command

### Sample

```yml
default_profile: gitlab.com
profiles:
  gitlab.com:
    token: ********************
    default_group: hoge
    default_project: hoge/soge
    default_assignee_id: 123
  gitlab.ssl.foo.jp:
    token: ******************** 
    default_group: foo
    default_project: foo/bar
    default_assignee_id: 456
```

## ToDos

- variable command
    - [x] Project-level Variables
    - [ ] Group-level Variables
- use template
    - [x] issue template
    - [x] merge request template
- [ ] pipeline actions
    - [ ] cancel
    - [ ] retry
- [ ] label command
- [x] project-member command
- workflow automation command
    - [ ] create
        - create new project and cloning repository
    - [ ] fork
        - create fork project and cloning repository
    - [ ] flow
        - create issue and create `WIP:` merge request
