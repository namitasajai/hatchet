// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.24.0
// source: workflow_runs.sql

package dbsqlc

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const countWorkflowRuns = `-- name: CountWorkflowRuns :one
SELECT
    count(runs) OVER() AS total
FROM
    "WorkflowRun" as runs
LEFT JOIN
    "WorkflowRunTriggeredBy" as runTriggers ON runTriggers."parentId" = runs."id"
LEFT JOIN
    "Event" as events ON runTriggers."eventId" = events."id"
LEFT JOIN
    "WorkflowVersion" as workflowVersion ON runs."workflowVersionId" = workflowVersion."id"
LEFT JOIN
    "Workflow" as workflow ON workflowVersion."workflowId" = workflow."id"
WHERE
    runs."tenantId" = $1 AND
    (
        $2::uuid IS NULL OR
        workflowVersion."id" = $2::uuid
    ) AND
    (
        $3::uuid IS NULL OR
        workflow."id" = $3::uuid
    ) AND
    (
        $4::uuid IS NULL OR
        events."id" = $4::uuid
    ) AND
    (
    $5::text IS NULL OR
    runs."concurrencyGroupId" = $5::text
    ) AND
    (
    $6::"WorkflowRunStatus" IS NULL OR
    runs."status" = $6::"WorkflowRunStatus"
    )
`

type CountWorkflowRunsParams struct {
	TenantId          pgtype.UUID           `json:"tenantId"`
	WorkflowVersionId pgtype.UUID           `json:"workflowVersionId"`
	WorkflowId        pgtype.UUID           `json:"workflowId"`
	EventId           pgtype.UUID           `json:"eventId"`
	GroupKey          pgtype.Text           `json:"groupKey"`
	Status            NullWorkflowRunStatus `json:"status"`
}

func (q *Queries) CountWorkflowRuns(ctx context.Context, db DBTX, arg CountWorkflowRunsParams) (int64, error) {
	row := db.QueryRow(ctx, countWorkflowRuns,
		arg.TenantId,
		arg.WorkflowVersionId,
		arg.WorkflowId,
		arg.EventId,
		arg.GroupKey,
		arg.Status,
	)
	var total int64
	err := row.Scan(&total)
	return total, err
}

const createGetGroupKeyRun = `-- name: CreateGetGroupKeyRun :one
INSERT INTO "GetGroupKeyRun" (
    "id",
    "createdAt",
    "updatedAt",
    "deletedAt",
    "tenantId",
    "workflowRunId",
    "workerId",
    "tickerId",
    "status",
    "input",
    "output",
    "requeueAfter",
    "scheduleTimeoutAt",
    "error",
    "startedAt",
    "finishedAt",
    "timeoutAt",
    "cancelledAt",
    "cancelledReason",
    "cancelledError"
) VALUES (
    COALESCE($1::uuid, gen_random_uuid()),
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP,
    NULL,
    $2::uuid,
    $3::uuid,
    NULL,
    NULL,
    'PENDING', -- default status
    $4::jsonb,
    NULL,
    $5::timestamp,
    $6::timestamp,
    NULL,
    NULL,
    NULL,
    NULL,
    NULL,
    NULL,
    NULL
) RETURNING id, "createdAt", "updatedAt", "deletedAt", "tenantId", "workerId", "tickerId", status, input, output, "requeueAfter", error, "startedAt", "finishedAt", "timeoutAt", "cancelledAt", "cancelledReason", "cancelledError", "workflowRunId", "scheduleTimeoutAt"
`

type CreateGetGroupKeyRunParams struct {
	ID                pgtype.UUID      `json:"id"`
	Tenantid          pgtype.UUID      `json:"tenantid"`
	Workflowrunid     pgtype.UUID      `json:"workflowrunid"`
	Input             []byte           `json:"input"`
	Requeueafter      pgtype.Timestamp `json:"requeueafter"`
	Scheduletimeoutat pgtype.Timestamp `json:"scheduletimeoutat"`
}

func (q *Queries) CreateGetGroupKeyRun(ctx context.Context, db DBTX, arg CreateGetGroupKeyRunParams) (*GetGroupKeyRun, error) {
	row := db.QueryRow(ctx, createGetGroupKeyRun,
		arg.ID,
		arg.Tenantid,
		arg.Workflowrunid,
		arg.Input,
		arg.Requeueafter,
		arg.Scheduletimeoutat,
	)
	var i GetGroupKeyRun
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
		&i.TenantId,
		&i.WorkerId,
		&i.TickerId,
		&i.Status,
		&i.Input,
		&i.Output,
		&i.RequeueAfter,
		&i.Error,
		&i.StartedAt,
		&i.FinishedAt,
		&i.TimeoutAt,
		&i.CancelledAt,
		&i.CancelledReason,
		&i.CancelledError,
		&i.WorkflowRunId,
		&i.ScheduleTimeoutAt,
	)
	return &i, err
}

const createJobRunLookupData = `-- name: CreateJobRunLookupData :one
INSERT INTO "JobRunLookupData" (
    "id",
    "createdAt",
    "updatedAt",
    "deletedAt",
    "jobRunId",
    "tenantId",
    "data"
) VALUES (
    COALESCE($1::uuid, gen_random_uuid()),
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP,
    NULL,
    $2::uuid,
    $3::uuid,
    jsonb_build_object(
        'input', COALESCE($4::jsonb, '{}'::jsonb),
        'triggered_by', $5::text,
        'steps', '{}'::jsonb
    )
) RETURNING id, "createdAt", "updatedAt", "deletedAt", "jobRunId", "tenantId", data
`

type CreateJobRunLookupDataParams struct {
	ID          pgtype.UUID `json:"id"`
	Jobrunid    pgtype.UUID `json:"jobrunid"`
	Tenantid    pgtype.UUID `json:"tenantid"`
	Input       []byte      `json:"input"`
	Triggeredby string      `json:"triggeredby"`
}

func (q *Queries) CreateJobRunLookupData(ctx context.Context, db DBTX, arg CreateJobRunLookupDataParams) (*JobRunLookupData, error) {
	row := db.QueryRow(ctx, createJobRunLookupData,
		arg.ID,
		arg.Jobrunid,
		arg.Tenantid,
		arg.Input,
		arg.Triggeredby,
	)
	var i JobRunLookupData
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
		&i.JobRunId,
		&i.TenantId,
		&i.Data,
	)
	return &i, err
}

const createJobRuns = `-- name: CreateJobRuns :many
INSERT INTO "JobRun" (
    "id",
    "createdAt",
    "updatedAt",
    "tenantId",
    "workflowRunId",
    "jobId",
    "status"
) 
SELECT
    gen_random_uuid(),
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP,
    $1::uuid,
    $2::uuid,
    "id",
    'PENDING' -- default status
FROM
    "Job"
WHERE
    "workflowVersionId" = $3::uuid
RETURNING "id"
`

type CreateJobRunsParams struct {
	Tenantid          pgtype.UUID `json:"tenantid"`
	Workflowrunid     pgtype.UUID `json:"workflowrunid"`
	Workflowversionid pgtype.UUID `json:"workflowversionid"`
}

func (q *Queries) CreateJobRuns(ctx context.Context, db DBTX, arg CreateJobRunsParams) ([]pgtype.UUID, error) {
	rows, err := db.Query(ctx, createJobRuns, arg.Tenantid, arg.Workflowrunid, arg.Workflowversionid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []pgtype.UUID
	for rows.Next() {
		var id pgtype.UUID
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		items = append(items, id)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const createStepRuns = `-- name: CreateStepRuns :exec
WITH job_id AS (
    SELECT "jobId"
    FROM "JobRun"
    WHERE "id" = $2::uuid
)
INSERT INTO "StepRun" (
    "id",
    "createdAt",
    "updatedAt",
    "tenantId",
    "jobRunId",
    "stepId",
    "status",
    "requeueAfter",
    "callerFiles"
) 
SELECT
    gen_random_uuid(),
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP,
    $1::uuid,
    $2::uuid,
    "id",
    'PENDING', -- default status
    CURRENT_TIMESTAMP + INTERVAL '5 seconds',
    '{}'
FROM
    "Step", job_id
WHERE
    "Step"."jobId" = job_id."jobId"
`

type CreateStepRunsParams struct {
	Tenantid pgtype.UUID `json:"tenantid"`
	Jobrunid pgtype.UUID `json:"jobrunid"`
}

func (q *Queries) CreateStepRuns(ctx context.Context, db DBTX, arg CreateStepRunsParams) error {
	_, err := db.Exec(ctx, createStepRuns, arg.Tenantid, arg.Jobrunid)
	return err
}

const createWorkflowRun = `-- name: CreateWorkflowRun :one
INSERT INTO "WorkflowRun" (
    "id",
    "createdAt",
    "updatedAt",
    "deletedAt",
    "displayName",
    "tenantId",
    "workflowVersionId",
    "status",
    "error",
    "startedAt",
    "finishedAt"
) VALUES (
    COALESCE($1::uuid, gen_random_uuid()),
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP,
    NULL, -- assuming deletedAt is not set on creation
    $2::text,
    $3::uuid,
    $4::uuid,
    'PENDING', -- default status
    NULL, -- assuming error is not set on creation
    NULL, -- assuming startedAt is not set on creation
    NULL  -- assuming finishedAt is not set on creation
) RETURNING "createdAt", "updatedAt", "deletedAt", "tenantId", "workflowVersionId", status, error, "startedAt", "finishedAt", "concurrencyGroupId", "displayName", id, "gitRepoBranch"
`

type CreateWorkflowRunParams struct {
	ID                pgtype.UUID `json:"id"`
	DisplayName       pgtype.Text `json:"displayName"`
	Tenantid          pgtype.UUID `json:"tenantid"`
	Workflowversionid pgtype.UUID `json:"workflowversionid"`
}

func (q *Queries) CreateWorkflowRun(ctx context.Context, db DBTX, arg CreateWorkflowRunParams) (*WorkflowRun, error) {
	row := db.QueryRow(ctx, createWorkflowRun,
		arg.ID,
		arg.DisplayName,
		arg.Tenantid,
		arg.Workflowversionid,
	)
	var i WorkflowRun
	err := row.Scan(
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
		&i.TenantId,
		&i.WorkflowVersionId,
		&i.Status,
		&i.Error,
		&i.StartedAt,
		&i.FinishedAt,
		&i.ConcurrencyGroupId,
		&i.DisplayName,
		&i.ID,
		&i.GitRepoBranch,
	)
	return &i, err
}

const createWorkflowRunTriggeredBy = `-- name: CreateWorkflowRunTriggeredBy :one
INSERT INTO "WorkflowRunTriggeredBy" (
    "id",
    "createdAt",
    "updatedAt",
    "deletedAt",
    "tenantId",
    "parentId",
    "eventId",
    "cronParentId",
    "cronSchedule",
    "scheduledId"
) VALUES (
    gen_random_uuid(), -- Generates a new UUID for id
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP,
    NULL, -- assuming deletedAt is not set on creation
    $1::uuid,
    $2::uuid, -- assuming parentId is the workflowRunId
    $3::uuid, -- NULL if not provided
    $4::uuid, -- NULL if not provided
    $5::text, -- NULL if not provided
    $6::uuid -- NULL if not provided
) RETURNING id, "createdAt", "updatedAt", "deletedAt", "tenantId", "eventId", "cronParentId", "cronSchedule", "scheduledId", input, "parentId"
`

type CreateWorkflowRunTriggeredByParams struct {
	Tenantid      pgtype.UUID `json:"tenantid"`
	Workflowrunid pgtype.UUID `json:"workflowrunid"`
	EventId       pgtype.UUID `json:"eventId"`
	CronParentId  pgtype.UUID `json:"cronParentId"`
	Cron          pgtype.Text `json:"cron"`
	ScheduledId   pgtype.UUID `json:"scheduledId"`
}

func (q *Queries) CreateWorkflowRunTriggeredBy(ctx context.Context, db DBTX, arg CreateWorkflowRunTriggeredByParams) (*WorkflowRunTriggeredBy, error) {
	row := db.QueryRow(ctx, createWorkflowRunTriggeredBy,
		arg.Tenantid,
		arg.Workflowrunid,
		arg.EventId,
		arg.CronParentId,
		arg.Cron,
		arg.ScheduledId,
	)
	var i WorkflowRunTriggeredBy
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
		&i.TenantId,
		&i.EventId,
		&i.CronParentId,
		&i.CronSchedule,
		&i.ScheduledId,
		&i.Input,
		&i.ParentId,
	)
	return &i, err
}

const linkStepRunParents = `-- name: LinkStepRunParents :exec
INSERT INTO "_StepRunOrder" ("A", "B")
SELECT 
    parent_run."id" AS "A",
    child_run."id" AS "B"
FROM 
    "_StepOrder" AS step_order
JOIN 
    "StepRun" AS parent_run ON parent_run."stepId" = step_order."A" AND parent_run."jobRunId" = $1::uuid
JOIN 
    "StepRun" AS child_run ON child_run."stepId" = step_order."B" AND child_run."jobRunId" = $1::uuid
`

func (q *Queries) LinkStepRunParents(ctx context.Context, db DBTX, jobrunid pgtype.UUID) error {
	_, err := db.Exec(ctx, linkStepRunParents, jobrunid)
	return err
}

const listStartableStepRuns = `-- name: ListStartableStepRuns :many
WITH job_run AS (
    SELECT "status"
    FROM "JobRun"
    WHERE "id" = $1::uuid
)
SELECT 
    child_run."id" AS "id"
FROM 
    "StepRun" AS child_run
LEFT JOIN 
    "_StepRunOrder" AS step_run_order ON step_run_order."B" = child_run."id"
JOIN
    job_run ON true
WHERE 
    child_run."jobRunId" = $1::uuid
    AND child_run."status" = 'PENDING'
    AND job_run."status" = 'RUNNING'
    -- case on whether parentStepRunId is null
    AND (
        ($2::uuid IS NULL AND step_run_order."A" IS NULL) OR 
        (
            step_run_order."A" = $2::uuid
            AND NOT EXISTS (
                SELECT 1
                FROM "_StepRunOrder" AS parent_order
                JOIN "StepRun" AS parent_run ON parent_order."A" = parent_run."id"
                WHERE 
                    parent_order."B" = child_run."id"
                    AND parent_run."status" != 'SUCCEEDED'
            )
        )
    )
`

type ListStartableStepRunsParams struct {
	Jobrunid        pgtype.UUID `json:"jobrunid"`
	ParentStepRunId pgtype.UUID `json:"parentStepRunId"`
}

func (q *Queries) ListStartableStepRuns(ctx context.Context, db DBTX, arg ListStartableStepRunsParams) ([]pgtype.UUID, error) {
	rows, err := db.Query(ctx, listStartableStepRuns, arg.Jobrunid, arg.ParentStepRunId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []pgtype.UUID
	for rows.Next() {
		var id pgtype.UUID
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		items = append(items, id)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listWorkflowRuns = `-- name: ListWorkflowRuns :many
SELECT
    runs."createdAt", runs."updatedAt", runs."deletedAt", runs."tenantId", runs."workflowVersionId", runs.status, runs.error, runs."startedAt", runs."finishedAt", runs."concurrencyGroupId", runs."displayName", runs.id, runs."gitRepoBranch", 
    workflow.id, workflow."createdAt", workflow."updatedAt", workflow."deletedAt", workflow."tenantId", workflow.name, workflow.description, 
    runtriggers.id, runtriggers."createdAt", runtriggers."updatedAt", runtriggers."deletedAt", runtriggers."tenantId", runtriggers."eventId", runtriggers."cronParentId", runtriggers."cronSchedule", runtriggers."scheduledId", runtriggers.input, runtriggers."parentId", 
    workflowversion.id, workflowversion."createdAt", workflowversion."updatedAt", workflowversion."deletedAt", workflowversion.version, workflowversion."order", workflowversion."workflowId", workflowversion.checksum, workflowversion."scheduleTimeout", 
    -- waiting on https://github.com/sqlc-dev/sqlc/pull/2858 for nullable events field
    events.id, events.key, events."createdAt", events."updatedAt"
FROM
    "WorkflowRun" as runs 
LEFT JOIN
    "WorkflowRunTriggeredBy" as runTriggers ON runTriggers."parentId" = runs."id"
LEFT JOIN
    "Event" as events ON runTriggers."eventId" = events."id"
LEFT JOIN
    "WorkflowVersion" as workflowVersion ON runs."workflowVersionId" = workflowVersion."id"
LEFT JOIN
    "Workflow" as workflow ON workflowVersion."workflowId" = workflow."id"
WHERE
    runs."tenantId" = $1 AND
    (
        $2::uuid IS NULL OR
        workflowVersion."id" = $2::uuid
    ) AND
    (
        $3::uuid IS NULL OR
        workflow."id" = $3::uuid
    ) AND
    (
        $4::uuid IS NULL OR
        events."id" = $4::uuid
    ) AND
    (
    $5::text IS NULL OR
    runs."concurrencyGroupId" = $5::text
    ) AND
    (
    $6::"WorkflowRunStatus" IS NULL OR
    runs."status" = $6::"WorkflowRunStatus"
    )
ORDER BY
    case when $7 = 'createdAt ASC' THEN runs."createdAt" END ASC ,
    case when $7 = 'createdAt DESC' then runs."createdAt" END DESC
OFFSET
    COALESCE($8, 0)
LIMIT
    COALESCE($9, 50)
`

type ListWorkflowRunsParams struct {
	TenantId          pgtype.UUID           `json:"tenantId"`
	WorkflowVersionId pgtype.UUID           `json:"workflowVersionId"`
	WorkflowId        pgtype.UUID           `json:"workflowId"`
	EventId           pgtype.UUID           `json:"eventId"`
	GroupKey          pgtype.Text           `json:"groupKey"`
	Status            NullWorkflowRunStatus `json:"status"`
	Orderby           interface{}           `json:"orderby"`
	Offset            interface{}           `json:"offset"`
	Limit             interface{}           `json:"limit"`
}

type ListWorkflowRunsRow struct {
	WorkflowRun            WorkflowRun            `json:"workflow_run"`
	Workflow               Workflow               `json:"workflow"`
	WorkflowRunTriggeredBy WorkflowRunTriggeredBy `json:"workflow_run_triggered_by"`
	WorkflowVersion        WorkflowVersion        `json:"workflow_version"`
	ID                     pgtype.UUID            `json:"id"`
	Key                    pgtype.Text            `json:"key"`
	CreatedAt              pgtype.Timestamp       `json:"createdAt"`
	UpdatedAt              pgtype.Timestamp       `json:"updatedAt"`
}

func (q *Queries) ListWorkflowRuns(ctx context.Context, db DBTX, arg ListWorkflowRunsParams) ([]*ListWorkflowRunsRow, error) {
	rows, err := db.Query(ctx, listWorkflowRuns,
		arg.TenantId,
		arg.WorkflowVersionId,
		arg.WorkflowId,
		arg.EventId,
		arg.GroupKey,
		arg.Status,
		arg.Orderby,
		arg.Offset,
		arg.Limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []*ListWorkflowRunsRow
	for rows.Next() {
		var i ListWorkflowRunsRow
		if err := rows.Scan(
			&i.WorkflowRun.CreatedAt,
			&i.WorkflowRun.UpdatedAt,
			&i.WorkflowRun.DeletedAt,
			&i.WorkflowRun.TenantId,
			&i.WorkflowRun.WorkflowVersionId,
			&i.WorkflowRun.Status,
			&i.WorkflowRun.Error,
			&i.WorkflowRun.StartedAt,
			&i.WorkflowRun.FinishedAt,
			&i.WorkflowRun.ConcurrencyGroupId,
			&i.WorkflowRun.DisplayName,
			&i.WorkflowRun.ID,
			&i.WorkflowRun.GitRepoBranch,
			&i.Workflow.ID,
			&i.Workflow.CreatedAt,
			&i.Workflow.UpdatedAt,
			&i.Workflow.DeletedAt,
			&i.Workflow.TenantId,
			&i.Workflow.Name,
			&i.Workflow.Description,
			&i.WorkflowRunTriggeredBy.ID,
			&i.WorkflowRunTriggeredBy.CreatedAt,
			&i.WorkflowRunTriggeredBy.UpdatedAt,
			&i.WorkflowRunTriggeredBy.DeletedAt,
			&i.WorkflowRunTriggeredBy.TenantId,
			&i.WorkflowRunTriggeredBy.EventId,
			&i.WorkflowRunTriggeredBy.CronParentId,
			&i.WorkflowRunTriggeredBy.CronSchedule,
			&i.WorkflowRunTriggeredBy.ScheduledId,
			&i.WorkflowRunTriggeredBy.Input,
			&i.WorkflowRunTriggeredBy.ParentId,
			&i.WorkflowVersion.ID,
			&i.WorkflowVersion.CreatedAt,
			&i.WorkflowVersion.UpdatedAt,
			&i.WorkflowVersion.DeletedAt,
			&i.WorkflowVersion.Version,
			&i.WorkflowVersion.Order,
			&i.WorkflowVersion.WorkflowId,
			&i.WorkflowVersion.Checksum,
			&i.WorkflowVersion.ScheduleTimeout,
			&i.ID,
			&i.Key,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const popWorkflowRunsRoundRobin = `-- name: PopWorkflowRunsRoundRobin :many
WITH running_count AS (
    SELECT
        COUNT(*) AS "count"
    FROM
        "WorkflowRun" r1
    JOIN
        "WorkflowVersion" workflowVersion ON r1."workflowVersionId" = workflowVersion."id"
    WHERE
        r1."tenantId" = $1 AND
        r1."status" = 'RUNNING' AND
        workflowVersion."id" = $2
), queued_row_numbers AS (
    SELECT
        r2.id,
        row_number() OVER (PARTITION BY r2."concurrencyGroupId" ORDER BY r2."createdAt") AS rn,
        row_number() over (order by r2."id" desc) as seqnum
    FROM
        "WorkflowRun" r2
    LEFT JOIN
        "WorkflowVersion" workflowVersion ON r2."workflowVersionId" = workflowVersion."id"
    WHERE
        r2."tenantId" = $1 AND
        r2."status" = 'QUEUED' AND
        workflowVersion."id" = $2
    ORDER BY
        rn ASC
), eligible_runs AS (
    SELECT
        id
    FROM
        queued_row_numbers
    WHERE
        queued_row_numbers."seqnum" <= ($3::int) - (SELECT "count" FROM running_count)
    FOR UPDATE SKIP LOCKED
)
UPDATE "WorkflowRun"
SET
    "status" = 'RUNNING'
FROM
    eligible_runs
WHERE
    "WorkflowRun".id = eligible_runs.id
RETURNING
    "WorkflowRun"."createdAt", "WorkflowRun"."updatedAt", "WorkflowRun"."deletedAt", "WorkflowRun"."tenantId", "WorkflowRun"."workflowVersionId", "WorkflowRun".status, "WorkflowRun".error, "WorkflowRun"."startedAt", "WorkflowRun"."finishedAt", "WorkflowRun"."concurrencyGroupId", "WorkflowRun"."displayName", "WorkflowRun".id, "WorkflowRun"."gitRepoBranch"
`

type PopWorkflowRunsRoundRobinParams struct {
	TenantId pgtype.UUID `json:"tenantId"`
	ID       pgtype.UUID `json:"id"`
	Maxruns  int32       `json:"maxruns"`
}

func (q *Queries) PopWorkflowRunsRoundRobin(ctx context.Context, db DBTX, arg PopWorkflowRunsRoundRobinParams) ([]*WorkflowRun, error) {
	rows, err := db.Query(ctx, popWorkflowRunsRoundRobin, arg.TenantId, arg.ID, arg.Maxruns)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []*WorkflowRun
	for rows.Next() {
		var i WorkflowRun
		if err := rows.Scan(
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.DeletedAt,
			&i.TenantId,
			&i.WorkflowVersionId,
			&i.Status,
			&i.Error,
			&i.StartedAt,
			&i.FinishedAt,
			&i.ConcurrencyGroupId,
			&i.DisplayName,
			&i.ID,
			&i.GitRepoBranch,
		); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const resolveWorkflowRunStatus = `-- name: ResolveWorkflowRunStatus :one
WITH jobRuns AS (
    SELECT sum(case when runs."status" = 'PENDING' then 1 else 0 end) AS pendingRuns,
        sum(case when runs."status" = 'RUNNING' then 1 else 0 end) AS runningRuns,
        sum(case when runs."status" = 'SUCCEEDED' then 1 else 0 end) AS succeededRuns,
        sum(case when runs."status" = 'FAILED' then 1 else 0 end) AS failedRuns,
        sum(case when runs."status" = 'CANCELLED' then 1 else 0 end) AS cancelledRuns
    FROM "JobRun" as runs
    WHERE
        "workflowRunId" = (
            SELECT "workflowRunId"
            FROM "JobRun"
            WHERE "id" = $1::uuid
        ) AND
        "tenantId" = $2::uuid
)
UPDATE "WorkflowRun"
SET "status" = CASE 
    -- Final states are final, cannot be updated
    WHEN "status" IN ('SUCCEEDED', 'FAILED') THEN "status"
    -- We check for running first, because if a job run is running, then the workflow is running
    WHEN j.runningRuns > 0 THEN 'RUNNING'
    -- When at least one job run has failed or been cancelled, then the workflow is failed
    WHEN j.failedRuns > 0 OR j.cancelledRuns > 0 THEN 'FAILED'
    -- When all job runs have succeeded, then the workflow is succeeded
    WHEN j.succeededRuns > 0 AND j.pendingRuns = 0 AND j.runningRuns = 0 AND j.failedRuns = 0 AND j.cancelledRuns = 0 THEN 'SUCCEEDED'
    ELSE "status"
END, "finishedAt" = CASE 
    -- Final states are final, cannot be updated
    WHEN "finishedAt" IS NOT NULL THEN "finishedAt"
    -- We check for running first, because if a job run is running, then the workflow is not finished
    WHEN j.runningRuns > 0 THEN NULL
    -- When one job run has failed or been cancelled, then the workflow is failed
    WHEN j.failedRuns > 0 OR j.cancelledRuns > 0 OR j.succeededRuns > 0 THEN NOW()
    ELSE "finishedAt"
END, "startedAt" = CASE 
    -- Started at is final, cannot be changed
    WHEN "startedAt" IS NOT NULL THEN "startedAt"
    -- If a job is running or in a final state, then the workflow has started
    WHEN j.runningRuns > 0 OR j.succeededRuns > 0 OR j.failedRuns > 0 OR j.cancelledRuns > 0 THEN NOW()
    ELSE "startedAt"
END
FROM
    jobRuns j
WHERE "id" = (
    SELECT "workflowRunId"
    FROM "JobRun"
    WHERE "id" = $1::uuid
) AND "tenantId" = $2::uuid
RETURNING "WorkflowRun"."createdAt", "WorkflowRun"."updatedAt", "WorkflowRun"."deletedAt", "WorkflowRun"."tenantId", "WorkflowRun"."workflowVersionId", "WorkflowRun".status, "WorkflowRun".error, "WorkflowRun"."startedAt", "WorkflowRun"."finishedAt", "WorkflowRun"."concurrencyGroupId", "WorkflowRun"."displayName", "WorkflowRun".id, "WorkflowRun"."gitRepoBranch"
`

type ResolveWorkflowRunStatusParams struct {
	Jobrunid pgtype.UUID `json:"jobrunid"`
	Tenantid pgtype.UUID `json:"tenantid"`
}

func (q *Queries) ResolveWorkflowRunStatus(ctx context.Context, db DBTX, arg ResolveWorkflowRunStatusParams) (*WorkflowRun, error) {
	row := db.QueryRow(ctx, resolveWorkflowRunStatus, arg.Jobrunid, arg.Tenantid)
	var i WorkflowRun
	err := row.Scan(
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
		&i.TenantId,
		&i.WorkflowVersionId,
		&i.Status,
		&i.Error,
		&i.StartedAt,
		&i.FinishedAt,
		&i.ConcurrencyGroupId,
		&i.DisplayName,
		&i.ID,
		&i.GitRepoBranch,
	)
	return &i, err
}

const updateManyWorkflowRun = `-- name: UpdateManyWorkflowRun :many
UPDATE
    "WorkflowRun"
SET
    "status" = COALESCE($1::"WorkflowRunStatus", "status"),
    "error" = COALESCE($2::text, "error"),
    "startedAt" = COALESCE($3::timestamp, "startedAt"),
    "finishedAt" = COALESCE($4::timestamp, "finishedAt")
WHERE 
    "tenantId" = $5::uuid AND
    "id" = ANY($6::uuid[])
RETURNING "WorkflowRun"."createdAt", "WorkflowRun"."updatedAt", "WorkflowRun"."deletedAt", "WorkflowRun"."tenantId", "WorkflowRun"."workflowVersionId", "WorkflowRun".status, "WorkflowRun".error, "WorkflowRun"."startedAt", "WorkflowRun"."finishedAt", "WorkflowRun"."concurrencyGroupId", "WorkflowRun"."displayName", "WorkflowRun".id, "WorkflowRun"."gitRepoBranch"
`

type UpdateManyWorkflowRunParams struct {
	Status     NullWorkflowRunStatus `json:"status"`
	Error      pgtype.Text           `json:"error"`
	StartedAt  pgtype.Timestamp      `json:"startedAt"`
	FinishedAt pgtype.Timestamp      `json:"finishedAt"`
	Tenantid   pgtype.UUID           `json:"tenantid"`
	Ids        []pgtype.UUID         `json:"ids"`
}

func (q *Queries) UpdateManyWorkflowRun(ctx context.Context, db DBTX, arg UpdateManyWorkflowRunParams) ([]*WorkflowRun, error) {
	rows, err := db.Query(ctx, updateManyWorkflowRun,
		arg.Status,
		arg.Error,
		arg.StartedAt,
		arg.FinishedAt,
		arg.Tenantid,
		arg.Ids,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []*WorkflowRun
	for rows.Next() {
		var i WorkflowRun
		if err := rows.Scan(
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.DeletedAt,
			&i.TenantId,
			&i.WorkflowVersionId,
			&i.Status,
			&i.Error,
			&i.StartedAt,
			&i.FinishedAt,
			&i.ConcurrencyGroupId,
			&i.DisplayName,
			&i.ID,
			&i.GitRepoBranch,
		); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateWorkflowRun = `-- name: UpdateWorkflowRun :one
UPDATE
    "WorkflowRun"
SET
    "status" = COALESCE($1::"WorkflowRunStatus", "status"),
    "error" = COALESCE($2::text, "error"),
    "startedAt" = COALESCE($3::timestamp, "startedAt"),
    "finishedAt" = COALESCE($4::timestamp, "finishedAt")
WHERE 
    "id" = $5::uuid AND
    "tenantId" = $6::uuid
RETURNING "WorkflowRun"."createdAt", "WorkflowRun"."updatedAt", "WorkflowRun"."deletedAt", "WorkflowRun"."tenantId", "WorkflowRun"."workflowVersionId", "WorkflowRun".status, "WorkflowRun".error, "WorkflowRun"."startedAt", "WorkflowRun"."finishedAt", "WorkflowRun"."concurrencyGroupId", "WorkflowRun"."displayName", "WorkflowRun".id, "WorkflowRun"."gitRepoBranch"
`

type UpdateWorkflowRunParams struct {
	Status     NullWorkflowRunStatus `json:"status"`
	Error      pgtype.Text           `json:"error"`
	StartedAt  pgtype.Timestamp      `json:"startedAt"`
	FinishedAt pgtype.Timestamp      `json:"finishedAt"`
	ID         pgtype.UUID           `json:"id"`
	Tenantid   pgtype.UUID           `json:"tenantid"`
}

func (q *Queries) UpdateWorkflowRun(ctx context.Context, db DBTX, arg UpdateWorkflowRunParams) (*WorkflowRun, error) {
	row := db.QueryRow(ctx, updateWorkflowRun,
		arg.Status,
		arg.Error,
		arg.StartedAt,
		arg.FinishedAt,
		arg.ID,
		arg.Tenantid,
	)
	var i WorkflowRun
	err := row.Scan(
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
		&i.TenantId,
		&i.WorkflowVersionId,
		&i.Status,
		&i.Error,
		&i.StartedAt,
		&i.FinishedAt,
		&i.ConcurrencyGroupId,
		&i.DisplayName,
		&i.ID,
		&i.GitRepoBranch,
	)
	return &i, err
}

const updateWorkflowRunGroupKey = `-- name: UpdateWorkflowRunGroupKey :one
WITH groupKeyRun AS (
    SELECT "id", "status" as groupKeyRunStatus, "output", "workflowRunId"
    FROM "GetGroupKeyRun" as groupKeyRun
    WHERE
        "id" = $2::uuid AND
        "tenantId" = $1::uuid
)
UPDATE "WorkflowRun" workflowRun
SET "status" = CASE 
    -- Final states are final, cannot be updated. We also can't move out of a queued state
    WHEN "status" IN ('SUCCEEDED', 'FAILED', 'QUEUED') THEN "status"
    -- When the GetGroupKeyRun failed or been cancelled, then the workflow is failed
    WHEN groupKeyRun.groupKeyRunStatus IN ('FAILED', 'CANCELLED') THEN 'FAILED'
    WHEN groupKeyRun.output IS NOT NULL THEN 'QUEUED'
    ELSE "status"
END, "finishedAt" = CASE 
    -- Final states are final, cannot be updated
    WHEN "finishedAt" IS NOT NULL THEN "finishedAt"
    -- When one job run has failed or been cancelled, then the workflow is failed
    WHEN groupKeyRun.groupKeyRunStatus IN ('FAILED', 'CANCELLED') THEN NOW()
    ELSE "finishedAt"
END, 
"concurrencyGroupId" = groupKeyRun."output"
FROM
    groupKeyRun
WHERE 
workflowRun."id" = groupKeyRun."workflowRunId" AND
workflowRun."tenantId" = $1::uuid
RETURNING workflowrun."createdAt", workflowrun."updatedAt", workflowrun."deletedAt", workflowrun."tenantId", workflowrun."workflowVersionId", workflowrun.status, workflowrun.error, workflowrun."startedAt", workflowrun."finishedAt", workflowrun."concurrencyGroupId", workflowrun."displayName", workflowrun.id, workflowrun."gitRepoBranch"
`

type UpdateWorkflowRunGroupKeyParams struct {
	Tenantid      pgtype.UUID `json:"tenantid"`
	Groupkeyrunid pgtype.UUID `json:"groupkeyrunid"`
}

func (q *Queries) UpdateWorkflowRunGroupKey(ctx context.Context, db DBTX, arg UpdateWorkflowRunGroupKeyParams) (*WorkflowRun, error) {
	row := db.QueryRow(ctx, updateWorkflowRunGroupKey, arg.Tenantid, arg.Groupkeyrunid)
	var i WorkflowRun
	err := row.Scan(
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
		&i.TenantId,
		&i.WorkflowVersionId,
		&i.Status,
		&i.Error,
		&i.StartedAt,
		&i.FinishedAt,
		&i.ConcurrencyGroupId,
		&i.DisplayName,
		&i.ID,
		&i.GitRepoBranch,
	)
	return &i, err
}
