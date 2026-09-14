package repository

import (
	"context"
	"errors"
	"restapi/internal/model"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TaskRepository struct {
	db *pgxpool.Pool
}

func NewTaskRepository(db *pgxpool.Pool) *TaskRepository {
	return &TaskRepository{db: db}
}

var ErrNotFound = errors.New("task not found")

func (r *TaskRepository) GetAll(ctx context.Context, limit, offset int) ([]model.Task, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, title, description, completed, created_at, updated_at
		FROM tasks
		ORDER BY id
		LIMIT $1 OFFSET $2
	`, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tasks := make([]model.Task, 0)

	for rows.Next() {
		var task model.Task

		if err := rows.Scan(
			&task.ID,
			&task.Title,
			&task.Description,
			&task.Completed,
			&task.CreatedAt,
			&task.UpdatedAt,
		); err != nil {
			return nil, err
		}

		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}

func (r *TaskRepository) Count(ctx context.Context) (int, error) {
	row := r.db.QueryRow(ctx, `
		SELECT COUNT(*) FROM tasks
	`)

	var count int
	err := row.Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r *TaskRepository) Create(ctx context.Context, title string, description string) (model.Task, error) {
	var task model.Task

	err := r.db.QueryRow(ctx, `
		INSERT INTO tasks (title, description)
		VALUES ($1, $2)
		RETURNING id, title, description, completed, created_at, updated_at
	`, title, description).Scan(
		&task.ID,
		&task.Title,
		&task.Description,
		&task.Completed,
		&task.CreatedAt,
		&task.UpdatedAt,
	)

	if err != nil {
		return model.Task{}, err
	}

	return task, nil
}

func (r *TaskRepository) GetByID(ctx context.Context, id int64) (model.Task, error) {
	var task model.Task

	err := r.db.QueryRow(ctx, `
		SELECT id, title, description, completed, created_at, updated_at
		FROM tasks
		WHERE id = $1
	`, id).Scan(
		&task.ID,
		&task.Title,
		&task.Description,
		&task.Completed,
		&task.CreatedAt,
		&task.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Task{}, ErrNotFound
		}

		return model.Task{}, err
	}

	return task, nil
}

func (r *TaskRepository) Delete(ctx context.Context, id int64) error {
	tag, err := r.db.Exec(ctx, `
		DELETE FROM tasks WHERE id = $1
	`, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}

	return nil
}

func (r *TaskRepository) Update(ctx context.Context, id int64, update model.UpdateTask) (model.Task, error) {
	var task model.Task

	err := r.db.QueryRow(ctx, `
		UPDATE tasks
		SET
			title = COALESCE($1, title),
			description = COALESCE($2, description),
			completed = COALESCE($3, completed),
			updated_at = NOW()
		WHERE id = $4
		RETURNING id, title, description, completed, created_at, updated_at
	`,
		update.Title,
		update.Description,
		update.Completed,
		id,
	).Scan(
		&task.ID,
		&task.Title,
		&task.Description,
		&task.Completed,
		&task.CreatedAt,
		&task.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Task{}, ErrNotFound
		}

		return model.Task{}, err
	}

	return task, nil
}
