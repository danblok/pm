package types

import (
	"time"
)

type Account struct {
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
	Id                  string    `json:"id"`
	Email               string    `json:"email"`
	Name                string    `json:"name"`
	Avatar              string    `json:"avatar,omitempty"`
	OwnedProjects       []Project `json:"owned_projets,omitempty"`
	ContributedProjects []Project `json:"contributed_projects,omitempty"`
	Deleted             bool      `json:"deleted"`
}

type Project struct {
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Owner        *Account  `json:"owner,omitempty"`
	Id           string    `json:"id"`
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	OwnerId      string    `json:"owner_id"`
	Contributors []Account `json:"contributors,omitempty"`
	Tasks        []Task    `json:"tasks,omitempty"`
	Statuses     []Status  `json:"statuses,omitempty"`
	Deleted      bool      `json:"deleted"`
}

type Task struct {
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Start     time.Time `json:"start"`
	End       time.Time `json:"end"`
	Status    *Status   `json:"status"`
	Project   *Project  `json:"project"`
	Id        string    `json:"id"`
	Name      string    `json:"name"`
	StatusId  string    `json:"status_id"`
	ProjectId string    `json:"project_id"`
	Deleted   bool      `json:"deleted"`
}

type Status struct {
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Project   *Project  `json:"project,omitempty"`
	Id        string    `json:"id"`
	Name      string    `json:"name"`
	ProjectId string    `json:"project_id"`
	Tasks     []Task    `json:"tasks"`
	Deleted   bool      `json:"deleted"`
}
