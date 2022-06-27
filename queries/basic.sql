select
    projects.kee,
    snapshots.created_at,
    metrics.name as metric_name,
	project_measures.value AS metric_value
from projects 
	left join ce_activity
		on projects.uuid=ce_activity.component_uuid
	left join snapshots
		on snapshots.uuid = ce_activity.analysis_uuid
	left join ce_scanner_context
		on ce_activity.uuid = ce_scanner_context.task_uuid
	left join project_measures
		on ce_activity.analysis_uuid = project_measures.analysis_uuid
	left join metrics
		on project_measures.metric_uuid=metrics.uuid