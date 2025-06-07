// pkg/repository/mysql.go
package repository

import (
	"context"
	"database/sql"

	"hearx/pkg/model"

	"go.uber.org/zap"
)

// TaskRepository defines DB operations for tasks.
type TaskRepository interface {
	Create(ctx context.Context, task model.Task) (model.Task, error)
	Update(ctx context.Context, task model.Task) (model.Task, error)
	FindAll(ctx context.Context) ([]model.Task, error)
	FindByID(ctx context.Context, id int64) (model.Task, error)
}

// mysqlTaskRepository is the MySQL implementation of TaskRepository.
type mysqlTaskRepository struct {
	db     *sql.DB
	logger *zap.Logger
}

// NewTaskRepository constructs a MySQL-backed TaskRepository.
func NewTaskRepository(db *sql.DB, logger *zap.Logger) TaskRepository {
	return &mysqlTaskRepository{db: db, logger: logger}
}

func (r *mysqlTaskRepository) Create(ctx context.Context, task model.Task) (model.Task, error) {
	r.logger.Info("creating task", zap.String("title", task.Title))
	res, err := r.db.ExecContext(ctx,
		`INSERT INTO tasks (title, description, completed)
         VALUES (?, ?, ?)`,
		task.Title, task.Description, task.Completed,
	)
	if err != nil {
		r.logger.Error("failed to create task", zap.Error(err), zap.String("title", task.Title))
		return model.Task{}, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		r.logger.Error("failed to retrieve last insert id", zap.Error(err))
		return model.Task{}, err
	}
	task.ID = id
	r.logger.Info("task created", zap.Int64("id", task.ID))
	return task, nil
}

func (r *mysqlTaskRepository) Update(ctx context.Context, task model.Task) (model.Task, error) {
	r.logger.Info("updating task", zap.Int64("id", task.ID))
	_, err := r.db.ExecContext(ctx,
		`UPDATE tasks
         SET title = ?, description = ?, completed = ?, updated_at = CURRENT_TIMESTAMP
         WHERE id = ?`,
		task.Title, task.Description, task.Completed, task.ID,
	)
	if err != nil {
		r.logger.Error("failed to update task", zap.Error(err), zap.Int64("id", task.ID))
		return model.Task{}, err
	}

	// fetch the updated record directly
	var updated model.Task
	row := r.db.QueryRowContext(ctx,
		`SELECT id, title, description, completed
         FROM tasks
         WHERE id = ? AND deleted_at IS NULL`,
		task.ID,
	)
	if scanErr := row.Scan(&updated.ID, &updated.Title, &updated.Description, &updated.Completed); scanErr != nil {
		r.logger.Error("failed to fetch updated task", zap.Error(scanErr), zap.Int64("id", task.ID))
		return model.Task{}, scanErr
	}

	r.logger.Info("task update fetched", zap.Any("task", updated))
	return updated, nil
}

func (r *mysqlTaskRepository) FindAll(ctx context.Context) ([]model.Task, error) {
	r.logger.Info("querying all tasks")
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, title, description, completed
         FROM tasks
         WHERE deleted_at IS NULL`,
	)
	if err != nil {
		r.logger.Error("failed to query tasks", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	var list []model.Task
	for rows.Next() {
		var t model.Task
		if err := rows.Scan(&t.ID, &t.Title, &t.Description, &t.Completed); err != nil {
			r.logger.Error("row scan error", zap.Error(err))
			return nil, err
		}
		list = append(list, t)
	}
	return list, rows.Err()
}

func (r *mysqlTaskRepository) FindByID(ctx context.Context, id int64) (model.Task, error) {
	r.logger.Info("querying task by id", zap.Int64("id", id))
	row := r.db.QueryRowContext(ctx,
		`SELECT id, title, description, completed
         FROM tasks
         WHERE id = ? AND deleted_at IS NULL`,
		id,
	)
	var t model.Task
	if err := row.Scan(&t.ID, &t.Title, &t.Description, &t.Completed); err != nil {
		r.logger.Error("failed to query task by id", zap.Error(err), zap.Int64("id", id))
		return model.Task{}, err
	}
	return t, nil
}
