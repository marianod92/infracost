package terraform

import (
	"fmt"
	"os"

	"github.com/infracost/infracost/internal/config"
	"github.com/infracost/infracost/internal/schema"
)

type PlanJSONProvider struct {
	ctx  *config.ProjectContext
	Path string
}

func NewPlanJSONProvider(ctx *config.ProjectContext) *PlanJSONProvider {
	return &PlanJSONProvider{
		ctx:  ctx,
		Path: ctx.ProjectConfig.Path,
	}
}

func (p *PlanJSONProvider) Type() string {
	return "terraform_plan_json"
}

func (p *PlanJSONProvider) DisplayType() string {
	return "Terraform plan JSON file"
}

func (p *PlanJSONProvider) AddMetadata(metadata *schema.ProjectMetadata) {
	// no op
}

func (p *PlanJSONProvider) LoadResources(usage map[string]*schema.UsageData) ([]*schema.Project, error) {
	j, err := os.ReadFile(p.Path)
	if err != nil {
		return []*schema.Project{}, fmt.Errorf("Error reading Terraform plan JSON file %w", err)
	}

	return p.LoadResourcesFromSrc(usage, j)
}

func (p *PlanJSONProvider) LoadResourcesFromSrc(usage map[string]*schema.UsageData, j []byte) ([]*schema.Project, error) {
	metadata := config.DetectProjectMetadata(p.ctx.ProjectConfig.Path)
	metadata.Type = p.Type()
	p.AddMetadata(metadata)
	name := schema.GenerateProjectName(metadata, p.ctx.RunContext.Config.EnableDashboard)

	project := schema.NewProject(name, metadata)
	parser := NewParser(p.ctx)

	pastResources, resources, err := parser.parseJSON(j, usage)
	if err != nil {
		return []*schema.Project{project}, fmt.Errorf("Error parsing Terraform plan JSON file %w", err)
	}

	project.PastResources = pastResources
	project.Resources = resources

	return []*schema.Project{project}, nil
}
