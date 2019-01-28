package internal

type ProjectProfileOption struct {
	Project string `long:"project" value-name:"<group>/<name>" description:"Specify the project to be processed"`
	Profile string `long:"profile" value-name:"<profile>" description:"Specify the profile defined in the config file"`
}

type BrowseOption struct {
	Browse bool `short:"b" long:"browse" description:"open browser"`
	URL    bool `short:"u" long:"url" description:"show browse url"`
	Copy   bool `short:"c" long:"copy" description:"copy browse url to clipboard"`
}

func (b *BrowseOption) HasBrowse() bool {
	if b.Browse || b.URL || b.Copy {
		return true
	}
	return false
}
