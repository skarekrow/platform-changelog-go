package structs

type ServicesData struct {
	ID          int             `json:"id,omitempty"`
	Name        string          `json:"name,omitempty"`
	DisplayName string          `json:"display_name,omitempty"`
	GHRepo      string          `json:"gh_repo,omitempty"`
	GLRepo      string          `json:"gl_repo,omitempty"`
	DeployFile  string          `json:"deploy_file,omitempty"`
	Namespace   string          `json:"namespace,omitempty"`
	Branch      string          `json:"branch,omitempty"`
	Timeline    []TimelinesData `json:"commits,omitempty"`
}

type TimelinesData struct {
	ID        int    `json:"id"`
	ServiceID int    `json:"service_id"`
	Type      string `json:"type"`
	Timestamp string `json:"timestamp"`
	Repo      string `json:"repo"`
	Ref       string `json:"ref"`
	Author    string `json:"author,omitempty"`
	MergedBy  string `json:"merged_by,omitempty"`
	Message   string `json:"message,omitempty"`
	Namespace string `json:"namespace,omitempty"`
	Cluster   string `json:"cluster,omitempty"`
	Image     string `json:"image,omitempty"`
}
