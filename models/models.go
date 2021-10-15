package models

type Task struct {
	Name        string
	Description string
	Command     string
}

type Tasks []Task
