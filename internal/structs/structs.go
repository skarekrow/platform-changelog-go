package structs

type ServicesData struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
	GHRepo      string `json:"gh_repo"`
	GLRepo      string `json:"gl_repo"`
	DeployFile  string `json:"deploy_file"`
	Namespace   string `json:"namespace"`
	Branch      string `json:"branch"`
}

type ExpandedServicesData struct {
	ID          int           `json:"id"`
	Name        string        `json:"name"`
	DisplayName string        `json:"display_name"`
	GHRepo      string        `json:"gh_repo"`
	GLRepo      string        `json:"gl_repo"`
	DeployFile  string        `json:"deploy_file"`
	Namespace   string        `json:"namespace"`
	Branch      string        `json:"branch"`
	Commit      TimelinesData `json:"latest_commit" gorm:"foreignkey:ID"`
	Deploy      TimelinesData `json:"latest_deploy" gorm:"foreignkey:ID"`
}

type TimelinesData struct {
	ID              int    `json:"id" gorm:"primary_key"`
	ServiceID       int    `json:"service_id" gorm:"foreign_key:services_data.id"`
	Type            string `json:"type"`
	Repo            string `json:"repo"`
	Ref             string `json:"ref,omitempty"`
	Timestamp       string `json:"timestamp"`
	Author          string `json:"author,omitempty"`
	MergedBy        string `json:"merged_by,omitempty"`
	Message         string `json:"message,omitempty"`
	DeployNamespace string `json:"namespace,omitempty"`
	Cluster         string `json:"cluster,omitempty"`
	Image           string `json:"image,omitempty"`
	TriggeredBy     string `json:"triggered_by,omitempty"`
	Status          string `json:"status,omitempty"`
}
