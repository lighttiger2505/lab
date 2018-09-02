package commands

const (
	ExitCodeOK        int = iota //0
	ExitCodeError     int = iota //1
	ExitCodeFileError int = iota //2
)
const IssueTemplateDir = ".gitlab/issue_templates"
const MergeRequestTemplateDir = ".gitlab/merge_request_templates"
