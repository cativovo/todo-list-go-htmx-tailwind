package postgres

// storage form of Todo
type Todo struct {
	Id        string
	TaskName  string `db:"task_name"`
	UpdatedAt string `db:"updated_at"`
	Completed bool
}
