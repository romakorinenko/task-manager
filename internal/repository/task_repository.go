package repository

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/huandu/go-sqlbuilder"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Task struct {
	ID          int       `db:"id" json:"id"`
	Title       string    `db:"title" json:"title"`
	Description string    `db:"description" json:"description"`
	Priority    int       `db:"priority" json:"priority"`
	Status      string    `db:"status" json:"status"`
	CreatedAt   time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt   time.Time `db:"updated_at" json:"updatedAt"`
	UserID      int       `db:"user_id" json:"userId,omitempty"`
}

var TaskStruct = sqlbuilder.NewStruct(new(Task))

type ITaskRepo interface {
	Create(ctx context.Context, task *Task) (*Task, error)
	Update(ctx context.Context, task *Task) (*Task, error)
	DeleteByID(ctx context.Context, taskID int) error
	GetByID(ctx context.Context, taskID int) (*Task, error)
	GetByUserID(ctx context.Context, userID int) ([]Task, error)
	GetByUserLogin(ctx context.Context, userLogin string) ([]Task, error)
	GetAll(ctx context.Context) ([]Task, error)
	GetByStatus(ctx context.Context, status string) ([]Task, error)
	GetByPriority(ctx context.Context, priority int) ([]Task, error)
}

type TaskRepo struct {
	dbPool *pgxpool.Pool
}

func NewTaskRepo(dbPool *pgxpool.Pool) *TaskRepo {
	return &TaskRepo{dbPool: dbPool}
}

func (t *TaskRepo) Create(ctx context.Context, task *Task) (*Task, error) {
	ID, err := t.generateNextTaskID(ctx)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	task.ID = ID
	task.CreatedAt = now
	task.UpdatedAt = now
	sql, args := TaskStruct.InsertInto("tasks", task).
		BuildWithFlavor(sqlbuilder.PostgreSQL) // todo

	row := t.dbPool.QueryRow(ctx, sql, args...)
	rowScanErr := row.Scan()
	if rowScanErr != nil && !errors.Is(rowScanErr, pgx.ErrNoRows) {
		return nil, err
	}

	return task, nil
}

func (t *TaskRepo) Update(ctx context.Context, task *Task) (*Task, error) {
	task.UpdatedAt = time.Now()

	ub := sqlbuilder.Update("tasks")
	sql, args := ub.Where(ub.Equal("id", task.ID)).
		Set(
			ub.Assign("title", task.Title),
			ub.Assign("description", task.Description),
			ub.Assign("priority", task.Priority),
			ub.Assign("status", task.Status),
			ub.Assign("updated_at", task.UpdatedAt),
		).
		BuildWithFlavor(sqlbuilder.PostgreSQL)

	_, err := t.dbPool.Exec(ctx, sql, args...)
	if err != nil {
		return nil, err
	}

	return task, nil
}

func (t *TaskRepo) DeleteByID(ctx context.Context, taskID int) error {
	db := TaskStruct.DeleteFrom("tasks")
	sql, args := db.Where(db.Equal("id", taskID)).
		BuildWithFlavor(sqlbuilder.PostgreSQL)

	_, err := t.dbPool.Exec(ctx, sql, args...)
	return err
}

func (t *TaskRepo) GetByID(ctx context.Context, taskID int) (*Task, error) {
	sb := TaskStruct.SelectFrom("tasks")
	sql, args := sb.Where(sb.Equal("id", taskID)).
		BuildWithFlavor(sqlbuilder.PostgreSQL)
	row := t.dbPool.QueryRow(ctx, sql, args...)

	var task Task
	rowScanErr := row.Scan(TaskStruct.Addr(&task)...)
	if rowScanErr != nil && errors.Is(rowScanErr, pgx.ErrNoRows) {
		return nil, rowScanErr
	}

	return &task, nil
}

func (t *TaskRepo) GetByUserID(ctx context.Context, userID int) ([]Task, error) {
	sb := TaskStruct.SelectFrom("tasks")
	sql, args := sb.Where(sb.Equal("user_id", userID)).
		BuildWithFlavor(sqlbuilder.PostgreSQL)

	rows, err := t.dbPool.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	res := make([]Task, 0)
	for rows.Next() {
		var task Task
		rowScanErr := rows.Scan(TaskStruct.Addr(&task)...)
		if rowScanErr != nil {
			return nil, rowScanErr
		}
		res = append(res, task)
	}

	return res, nil
}

func (t *TaskRepo) GetByUserLogin(ctx context.Context, userLogin string) ([]Task, error) {
	sb := TaskStruct.SelectFrom("tasks")
	sql, args := sb.JoinWithOption(sqlbuilder.LeftJoin, "users", "tasks.user_id = users.id").
		Where(sb.Equal("users.login", userLogin)).
		BuildWithFlavor(sqlbuilder.PostgreSQL)

	rows, err := t.dbPool.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	res := make([]Task, 0)
	for rows.Next() {
		var task Task
		rowScanErr := rows.Scan(TaskStruct.Addr(&task)...)
		if rowScanErr != nil {
			return nil, rowScanErr
		}
		res = append(res, task)
	}

	return res, nil
}

func (t *TaskRepo) GetAll(ctx context.Context) ([]Task, error) {
	sql, _ := TaskStruct.SelectFrom("tasks").
		OrderBy("id").
		BuildWithFlavor(sqlbuilder.PostgreSQL)

	rows, err := t.dbPool.Query(ctx, sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	res := make([]Task, 0)
	for rows.Next() {
		var task Task
		if rowScanErr := rows.Scan(TaskStruct.Addr(&task)...); rowScanErr != nil {
			slog.Info(err.Error())
			return nil, err
		}

		res = append(res, task)
	}

	return res, nil
}

func (t *TaskRepo) GetByStatus(ctx context.Context, status string) ([]Task, error) {
	sb := TaskStruct.SelectFrom("tasks")
	sql, args := sb.Where(sb.Equal("status", status)).
		OrderBy("id").
		BuildWithFlavor(sqlbuilder.PostgreSQL)

	rows, err := t.dbPool.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	res := make([]Task, 0)
	for rows.Next() {
		var task Task
		if rowScanErr := rows.Scan(TaskStruct.Addr(&task)...); rowScanErr != nil {
			slog.Info(err.Error())
			return nil, err
		}

		res = append(res, task)
	}

	return res, nil
}

func (t *TaskRepo) GetByPriority(ctx context.Context, priority int) ([]Task, error) {
	sb := TaskStruct.SelectFrom("tasks")
	sql, args := sb.Where(sb.Equal("priority", priority)).
		OrderBy("id").
		BuildWithFlavor(sqlbuilder.PostgreSQL)

	rows, err := t.dbPool.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	res := make([]Task, 0)
	for rows.Next() {
		var task Task
		if rowScanErr := rows.Scan(TaskStruct.Addr(&task)...); rowScanErr != nil {
			slog.Info(err.Error())
			return nil, err
		}

		res = append(res, task)
	}

	return res, nil
}

func (t *TaskRepo) generateNextTaskID(ctx context.Context) (int, error) {
	rows, err := t.dbPool.Query(ctx, fmt.Sprintf("SELECT nextval('%s')", "tasks_sequence"))
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	if rows.Next() {
		var id int
		rowScanErr := rows.Scan(&id)
		if rowScanErr != nil {
			return 0, rowScanErr
		}
		return id, nil
	}
	return 0, fmt.Errorf("something was wrong. there is no next task id")
}
