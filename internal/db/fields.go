package db

var (
	timelinesFields = []string{"timelines.id", "timelines.timestamp", "timelines.service_id", "timelines.ref", "timelines.repo", "timelines.type"}
	commitsFields   = []string{"timelines.author", "timelines.message", "timelines.merged_by"}
	deploysFields   = []string{"timelines.deploy_namespace", "timelines.cluster", "timelines.image"}
)
