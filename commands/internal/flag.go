package internal

type ProjectProfileOption struct {
	Project string `long:"project" value-name:"<group>/<name>" description:"Specify the project to be processed"`
	Profile string `long:"profile" value-name:"<profile>" description:"Specify the profile defined in the config file"`
}
