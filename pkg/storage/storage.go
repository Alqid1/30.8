package storage

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
)

type Interface interface {
	Tasks(int, int) ([]Task, error)
	NewTask(Task) (int, error)
	TaskByAuthor(int) ([]Task, error)
	TaskByLabel(string) ([]Task, error)
	UpdateTask(int, Task) error
	DeleteTask(int) error
}

// Хранилище данных.
type Storage struct {
	db *pgxpool.Pool
}

// Конструктор, принимает строку подключения к БД.
func New(constr string) (*Storage, error) {
	db, err := pgxpool.Connect(context.Background(), constr)
	if err != nil {
		return nil, err
	}
	s := Storage{
		db: db,
	}
	return &s, nil
}

// Задача.
type Task struct {
	ID         int
	Opened     int64
	Closed     int64
	AuthorID   int
	AssignedID int
	Title      string
	Content    string
}

// Tasks возвращает список задач из БД.
func (s *Storage) Tasks(taskID, authorID int) ([]Task, error) {
	rows, err := s.db.Query(context.Background(), `
		SELECT 
			id,
			opened,
			closed,
			author_id,
			assigned_id,
			title,
			content
		FROM tasks
		WHERE
			($1 = 0 OR id = $1) AND
			($2 = 0 OR author_id = $2)
		ORDER BY id;
	`,
		taskID,
		authorID,
	)
	if err != nil {
		return nil, err
	}
	var tasks []Task
	// итерирование по результату выполнения запроса
	// и сканирование каждой строки в переменную
	for rows.Next() {
		var t Task
		err = rows.Scan(
			&t.ID,
			&t.Opened,
			&t.Closed,
			&t.AuthorID,
			&t.AssignedID,
			&t.Title,
			&t.Content,
		)
		if err != nil {
			return nil, err
		}
		// добавление переменной в массив результатов
		tasks = append(tasks, t)

	}
	// ВАЖНО не забыть проверить rows.Err()
	return tasks, rows.Err()
}

// NewTask создаёт новую задачу и возвращает её id.
func (s *Storage) NewTask(t Task) (int, error) {
	var id int
	err := s.db.QueryRow(context.Background(), `
		INSERT INTO tasks (title, content)
		VALUES ($1, $2) RETURNING id;
		`,
		t.Title,
		t.Content,
	).Scan(&id)
	return id, err
}

// TaskByAuthor возвращает список задач опеределенного автора (по id автора)
func (s *Storage) TaskByAuthor(authorID int) ([]Task, error) {
	rows, err := s.db.Query(context.Background(), `
		SELECT 
			tasks.id,
			tasks.opened,
			tasks.closed,
			tasks.author_id,
			tasks.assigned_id,
			tasks.title,
			tasks.content
		FROM tasks
		JOIN users ON users.id = tasks.author_id
		WHERE
			users.id = $1;
	`,
		authorID,
	)
	if err != nil {
		return nil, err
	}
	var tasks []Task
	// итерирование по результату выполнения запроса
	// и сканирование каждой строки в переменную
	for rows.Next() {
		var t Task
		err = rows.Scan(
			&t.ID,
			&t.Opened,
			&t.Closed,
			&t.AuthorID,
			&t.AssignedID,
			&t.Title,
			&t.Content,
		)
		if err != nil {
			return nil, err
		}
		// добавление переменной в массив результатов
		tasks = append(tasks, t)

	}
	// ВАЖНО не забыть проверить rows.Err()
	return tasks, rows.Err()
}

// TaskByLabel возвращает список задач с определенной меткой (по названию метки)
func (s *Storage) TaskByLabel(labelName string) ([]Task, error) {
	rows, err := s.db.Query(context.Background(), `
		SELECT 
			tasks.id,
			tasks.opened,
			tasks.closed,
			tasks.author_id,
			tasks.assigned_id,
			tasks.title,
			tasks.content
		FROM tasks
		JOIN tasks_labels ON tasks.id = tasks_labels.task_id
		JOIN labels ON labels.id = tasks_labels.label_id
		WHERE
			labels.name = '$1';
	`,
		labelName,
	)
	if err != nil {
		return nil, err
	}
	var tasks []Task
	// итерирование по результату выполнения запроса
	// и сканирование каждой строки в переменную
	for rows.Next() {
		var t Task
		err = rows.Scan(
			&t.ID,
			&t.Opened,
			&t.Closed,
			&t.AuthorID,
			&t.AssignedID,
			&t.Title,
			&t.Content,
		)
		if err != nil {
			return nil, err
		}
		// добавление переменной в массив результатов
		tasks = append(tasks, t)

	}
	// ВАЖНО не забыть проверить rows.Err()
	return tasks, rows.Err()
}

// UpdateTask позволяет обновлять AssignedID,Title,Content по указаному id
func (s *Storage) UpdateTask(taskID int, t Task) error {
	_, err := s.db.Exec(context.Background(), `
		UPDATE tasks
		SET assigned_id = $1, title = $2,content = $3
		WHERE id=$4;
		`,
		t.AssignedID,
		t.Title,
		t.Content,
		taskID,
	)
	return err
}

// DeleteTask удаляет задачу по id
func (s *Storage) DeleteTask(taskID int) error {
	_, err := s.db.Exec(context.Background(), `
		DELETE FROM tasks
		WHERE
			id = $1;
	`,
		taskID,
	)
	if err != nil {
		return err
	}
	return nil
}
