package db

import (
	"fmt"
	"time"
)

type Task struct {
	ID      string `json:"id"` // здесь проверить мб надо int?
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

// AddTask добавляет задачу в базу данных и возвращает ID новой задачи.
func AddTask(task *Task) (int64, error) {
	query := `INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)` // использую подготовленные выражения, чтобы избежать SQL-инъекцию (проверить надо ли в тестах)
	res, err := DB.Exec(query, task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId() // возвращает ID только что вставленной строки.
}

// Tasks возвращает список задач из базы данных, ограниченный limit.
func Tasks(limit int) ([]*Task, error) {
	tasks := []*Task{} // здесь и дальше инициализируем пустым слайсом, чтобы вернулся [] вместо null, если нет задач, как в рекомендации.

	rows, err := DB.Query(
		`SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date LIMIT ?`,
		limit,
	)

	if err != nil {
		return nil, err
	}
	defer rows.Close() // закрываем после завершения функции, даже если будет ошибка.

	for rows.Next() {
		task := &Task{}
		err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat) // scan записывает значения колонок в переменные по порядкую
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	return tasks, rows.Err()
}

// выполняем задачу со звёздочкой - поиск

func TasksSearch(limit int, search string) ([]*Task, error) {
	tasks := []*Task{}

	// проверяем не является ли строка поиска датой в формате 02.01.2006.
	t, err := time.Parse("02.01.2006", search)
	if err == nil {
		// если дата, то ищем по точному совпадению, переведя в формат
		date := t.Format("20060102")
		rows, err := DB.Query(
			`SELECT id, date, title, comment, repeat FROM scheduler WHERE date = ? LIMIT ?`,
			date, limit,
		)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		for rows.Next() {
			task := &Task{}
			if err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat); err != nil {
				return nil, err
			}
			tasks = append(tasks, task)
		}
		return tasks, rows.Err()
	}

	// если не дата, то ищем по вхождению в заголовок или комментарий
	searchPattern := "%" + search + "%" // не вставляем searchPattern напрямую в SQL, используем подготовленные выражения.
	rows, err := DB.Query(
		`SELECT id, date, title, comment, repeat FROM scheduler WHERE title LIKE ? OR comment LIKE ? ORDER BY date LIMIT ?`,
		searchPattern, searchPattern, limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		task := &Task{}
		if err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	return tasks, rows.Err()
}

func GetTask(id string) (*Task, error) {
	task := &Task{}
	err := DB.QueryRow(
		`SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?`,
		id,
	).Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		return nil, err
	}
	return task, nil
}

// UpdateTask обновляет задачу в базе данных. Если задача с указанным ID не найдена, возвращает ошибку.
func UpdateTask(task *Task) error {
	query := `UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?`
	res, err := DB.Exec(query, task.Date, task.Title, task.Comment, task.Repeat, task.ID)
	if err != nil {
		return err
	}
	// метод RowsAffected() возвращает количество записей к которым
	// была применена SQL команда
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf("incorrect id for updating task")
	}
	return nil
}

// DeleteTask удаляет задачу по ID. Если задача не найдена, возвращает ошибку.
// по сути идентичны UpdateTask, только вместо UPDATE используем DELETE
func DeleteTask(id string) error {
	res, err := DB.Exec(`DELETE FROM scheduler WHERE id = ?`, id)
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf("task not found")
	}
	return nil
}

// UpdateDate обновляет дату задачи по ID. Если задача не найдена, возвращает ошибку
// идентичен UpdateTask, только используем UPDATE только для поля даты
func UpdateDate(id string, date string) error {
	res, err := DB.Exec(`UPDATE scheduler SET date = ? WHERE id = ?`, date, id)
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf("task not found")
	}
	return nil
}
