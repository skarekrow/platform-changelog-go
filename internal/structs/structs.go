package structs

type TimelineData struct {
	ID        int    `json:"id,omitempty"`
	ServiceID int    `json:"service_id,omitempty"`
	CommitID  int    `json:"commit_id,omitempty"`
	DeployID  int    `json:"deploy_id,omitempty"`
	Timestamp string `json:"timestamp,omitempty"`
	Repo      string `json:"repo,omitempty"`
	Ref       string `json:"ref,omitempty"`
	Author    string `json:"author,omitempty"`
	MergedBy  string `json:"merged_by,omitempty"`
	Message   string `json:"message,omitempty"`
	Namespace string `json:"namespace,omitempty"`
	Cluster   string `json:"cluster,omitempty"`
	Image     string `json:"image,omitempty"`
}

type ServiceData struct {
	ID          int          `json:"id,omitempty"`
	Name        string       `json:"name,omitempty"`
	DisplayName string       `json:"display_name,omitempty"`
	GHRepo      string       `json:"gh_repo,omitempty"`
	GLRepo      string       `json:"gl_repo,omitempty"`
	DeployFile  string       `json:"deploy_file,omitempty"`
	Namespace   string       `json:"namespace,omitempty"`
	Branch      string       `json:"branch,omitempty"`
	Commits     []CommitData `json:"commits,omitempty"`
	Deploys     []DeployData `json:"deploys,omitempty"`
}

type CommitData struct {
	ID        int    `json:"id,omitempty"`
	ServiceID int    `json:"service_id,omitempty"`
	Ref       string `json:"ref,omitempty"`
	Timestamp string `json:"timestamp,omitempty"`
	Repo      string `json:"repo,omitempty"`
	Author    string `json:"author,omitempty"`
	MergedBy  string `json:"merged_by,omitempty"`
	Message   string `json:"message,omitempty"`
}

type DeployData struct {
	ID        int    `json:"id,omitempty"`
	ServiceID int    `json:"service_id,omitempty"`
	Ref       string `json:"ref,omitempty"`
	Timestamp string `json:"timestamp,omitempty"`
	Repo      string `json:"repo,omitempty"`
	Namespace string `json:"namespace,omitempty"`
	Cluster   string `json:"cluster,omitempty"`
	Image     string `json:"image,omitempty"`
}
