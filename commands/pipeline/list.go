package pipeline

import (
	"strconv"
	"strings"

	"github.com/lighttiger2505/lab/internal/api"
	"github.com/lighttiger2505/lab/internal/gitutil"
	"github.com/ryanuber/columnize"
	gitlab "github.com/xanzy/go-gitlab"
)

type listMethod struct {
	client api.Pipeline
	opt    *ListOption
	pInfo  *gitutil.GitLabProjectInfo
}

func (m *listMethod) Process() (string, error) {
	pipelines, err := m.client.ProjectPipelines(
		m.pInfo.Project,
		makeListPipelineOptions(m.opt, m.pInfo),
	)
	if err != nil {
		return "", err
	}

	result := columnize.SimpleFormat(pipelineListOutput(pipelines))
	return result, nil
}

type listJobMethod struct {
	client  api.Pipeline
	opt     *ListOption
	project string
	id      int
}

func (m *listJobMethod) Process() (string, error) {
	jobs, err := m.client.ProjectPipelineJobs(
		m.project,
		makeListPiplineJobOptions(),
		m.id,
	)
	if err != nil {
		return "", err
	}
	result := columnize.SimpleFormat(pipelineJobListOutput(jobs))
	return result, nil
}

func makeListPipelineOptions(listPipelineOption *ListOption, pInfo *gitutil.GitLabProjectInfo) *gitlab.ListProjectPipelinesOptions {
	var scope *string
	if listPipelineOption.Scope != "" {
		scope = gitlab.String(listPipelineOption.Scope)
	}
	var status *gitlab.BuildStateValue
	if listPipelineOption.States != "" {
		v := gitlab.BuildStateValue(listPipelineOption.States)
		status = &v
	}
	listOption := &gitlab.ListOptions{
		Page:    1,
		PerPage: listPipelineOption.Num,
	}
	listPipelinesOptions := &gitlab.ListProjectPipelinesOptions{
		Scope:       scope,
		Status:      status,
		Ref:         gitlab.String(listPipelineOption.getRef(pInfo.CurrentBranch)),
		YamlErrors:  gitlab.Bool(false),
		Name:        gitlab.String(""),
		Username:    gitlab.String(""),
		OrderBy:     gitlab.String(listPipelineOption.OrderBy),
		Sort:        gitlab.String(listPipelineOption.Sort),
		ListOptions: *listOption,
	}
	return listPipelinesOptions
}

func makeListPiplineJobOptions() *gitlab.ListJobsOptions {
	return &gitlab.ListJobsOptions{}
}

func pipelineListOutput(pipelines gitlab.PipelineList) []string {
	var outputs []string
	for _, pipeline := range pipelines {
		output := strings.Join([]string{
			strconv.Itoa(pipeline.ID),
			pipeline.Status,
			pipeline.Ref,
			pipeline.SHA,
		}, "|")
		outputs = append(outputs, output)
	}
	return outputs
}

func pipelineJobListOutput(jobs []*gitlab.Job) []string {
	var outputs []string
	for _, job := range jobs {
		output := strings.Join([]string{
			strconv.Itoa(job.ID),
			job.Status,
			job.Ref,
			job.Commit.ShortID,
			job.User.Username,
			job.Stage,
			job.Name,
		}, "|")
		outputs = append(outputs, output)
	}
	return outputs
}
