package service

import (
	"context"
	"testing"

	"github.com/danblok/pm/internals/types"
	"github.com/google/uuid"
)

func TestGetProjectById(t *testing.T) {
	s, cleanUp := setupServiceLifetime(t)
	defer cleanUp("projects", "accounts")

	owner := types.Account{
		Id:    uuid.NewString(),
		Name:  "username",
		Email: "username@username.com",
	}
	want := types.Project{
		Id:          uuid.NewString(),
		Name:        "project",
		Description: "A new project",
		OwnerId:     owner.Id,
	}
	_, err := s.DB.Exec("INSERT INTO accounts (id, email, name) VALUES ($1, $2, $3)", owner.Id, owner.Email, owner.Name)
	if err != nil {
		t.Fatal(ErrFailedToPrepareTest)
	}
	_, err = s.DB.Exec("INSERT INTO projects (id, name, description, owner_id) VALUES ($1, $2, $3, $4)", want.Id, want.Name, want.Description, want.OwnerId)
	if err != nil {
		t.Fatal(ErrFailedToPrepareTest)
	}

	got, err := s.GetProjectById(context.TODO(), want.Id)
	if err != nil {
		t.Fatal(err)
	}
	if want.Id != got.Id {
		t.Fatalf("want Id: %s, got Id: %s", want.Id, got.Id)
	}
	if want.Name != got.Name {
		t.Fatalf("want Name: %s, got Name: %s", want.Name, got.Name)
	}
	if want.Description != got.Description {
		t.Fatalf("want Description: %s, got Description: %s", want.Description, got.Description)
	}
	if want.OwnerId != got.OwnerId {
		t.Fatalf("want OwnerId: %s, got OwnerId: %s", want.OwnerId, got.OwnerId)
	}

	_, err = s.GetProjectById(context.TODO(), uuid.NewString())
	if err == nil {
		t.Fatal("should return err: ", err)
	}
	_, err = s.GetProjectById(context.TODO(), "")
	if err == nil {
		t.Fatal("should return err: ", err)
	}
}

func TestGetProjectsByOwnerId(t *testing.T) {
	s, cleanUp := setupServiceLifetime(t)
	defer cleanUp("projects", "accounts")

	owner := types.Account{
		Id:    uuid.NewString(),
		Name:  "username",
		Email: "username@username.com",
	}
	want := []types.Project{
		{
			Id:          uuid.NewString(),
			Name:        "project1",
			Description: "A first project",
			OwnerId:     owner.Id,
		},
		{
			Id:          uuid.NewString(),
			Name:        "project2",
			Description: "A second project",
			OwnerId:     owner.Id,
		},
	}
	_, err := s.DB.Exec("INSERT INTO accounts (id, email, name) VALUES ($1, $2, $3)", owner.Id, owner.Email, owner.Name)
	if err != nil {
		t.Fatal(ErrFailedToPrepareTest)
	}
	for i := range want {
		_, err = s.DB.Exec("INSERT INTO projects (id, name, description, owner_id) VALUES ($1, $2, $3, $4)", want[i].Id, want[i].Name, want[i].Description, want[i].OwnerId)
		if err != nil {
			t.Fatal(ErrFailedToPrepareTest)
		}
	}

	got, err := s.GetProjectsByOwnerId(context.TODO(), owner.Id)
	if err != nil {
		t.Fatal(err)
	}
	for i := range got {
		if want[i].Id != got[i].Id {
			t.Fatalf("want Id: %s, got Id: %s", want[i].Id, got[i].Id)
		}
		if want[i].Name != got[i].Name {
			t.Fatalf("want Name: %s, got Name: %s", want[i].Name, got[i].Name)
		}
		if want[i].Description != got[i].Description {
			t.Fatalf("want Description: %s, got Description: %s", want[i].Description, got[i].Description)
		}
		if want[i].OwnerId != got[i].OwnerId {
			t.Fatalf("want OwnerId: %s, got OwnerId: %s", want[i].OwnerId, got[i].OwnerId)
		}
	}
}

func TestAddProject(t *testing.T) {
	s, cleanUp := setupServiceLifetime(t)
	defer cleanUp("accounts", "projects")

	owner := types.Account{
		Id:    uuid.NewString(),
		Name:  "username",
		Email: "username@username.com",
	}
	_, err := s.DB.Exec("INSERT INTO accounts (id, email, name) VALUES ($1, $2, $3)", owner.Id, owner.Email, owner.Name)
	if err != nil {
		t.Fatal(ErrFailedToPrepareTest)
	}
	want := types.Project{
		Name:        "project",
		Description: "A new project",
		OwnerId:     owner.Id,
	}

	input := AddProjectInput{
		Name:        want.Name,
		Description: want.Description,
		OwnerId:     want.OwnerId,
	}
	err = s.AddProject(context.TODO(), &input)
	if err != nil {
		t.Fatal(err)
	}

	var got types.Project
	row := s.DB.QueryRow("SELECT name, description, owner_id FROM projects LIMIT 1")
	if err != nil {
		t.Fatal(err)
	}
	err = row.Scan(&got.Name, &got.Description, &got.OwnerId)
	if err != nil {
		t.Fatal(err)
	}

	if err != nil {
		t.Fatal(err)
	}
	if want.Name != got.Name {
		t.Fatalf("want Name: %s, got Name: %s", want.Name, got.Name)
	}
	if want.Description != got.Description {
		t.Fatalf("want Description: %s, got Description: %s", want.Description, got.Description)
	}
	if want.OwnerId != got.OwnerId {
		t.Fatalf("want OwnerId: %s, got OwnerId: %s", want.OwnerId, got.OwnerId)
	}

	input = AddProjectInput{
		Name:        "",
		Description: want.Description,
		OwnerId:     want.OwnerId,
	}

	err = s.AddProject(context.TODO(), &input)
	if err == nil {
		t.Fatal("should err: ", ErrFailedValidation)
	}

	input = AddProjectInput{
		Name:        want.Name,
		Description: want.Description,
		OwnerId:     "",
	}

	err = s.AddProject(context.TODO(), &input)
	if err == nil {
		t.Fatal("should err: ", ErrFailedValidation)
	}
}

func TestUpdateProject(t *testing.T) {
	s, cleanUp := setupServiceLifetime(t)
	defer cleanUp("accounts", "projects")

	owner := types.Account{
		Id:    uuid.NewString(),
		Name:  "username",
		Email: "username@username.com",
	}
	p := types.Project{
		Id:          uuid.NewString(),
		Name:        "project",
		Description: "A new project",
		OwnerId:     owner.Id,
	}
	want := types.Project{
		Id:          p.Id,
		Name:        "New name",
		Description: "A new project",
		OwnerId:     owner.Id,
	}
	_, err := s.DB.Exec("INSERT INTO accounts (id, email, name) VALUES ($1, $2, $3)", owner.Id, owner.Email, owner.Name)
	if err != nil {
		t.Fatal(ErrFailedToPrepareTest)
	}
	_, err = s.DB.Exec("INSERT INTO projects (id, name, description, owner_id) VALUES ($1, $2, $3, $4)", p.Id, p.Name, p.Description, p.OwnerId)
	if err != nil {
		t.Fatal(ErrFailedToPrepareTest)
	}

	input := UpdateProjectInput{
		Id:   p.Id,
		Name: want.Name,
	}
	err = s.UpdateProject(context.TODO(), &input)
	if err != nil {
		t.Fatal(err)
	}

	var got types.Project
	row := s.DB.QueryRow("SELECT id, name, description, owner_id FROM projects WHERE id=$1", input.Id)
	if err != nil {
		t.Fatal(err)
	}
	err = row.Scan(&got.Id, &got.Name, &got.Description, &got.OwnerId)
	if err != nil {
		t.Fatal(err)
	}

	if err != nil {
		t.Fatal(err)
	}
	if want.Id != got.Id {
		t.Fatalf("want Id: %s, got Id: %s", want.Id, got.Id)
	}
	if want.Name != got.Name {
		t.Fatalf("want Name: %s, got Name: %s", want.Name, got.Name)
	}
	if want.Description != got.Description {
		t.Fatalf("want Description: %s, got Description: %s", want.Description, got.Description)
	}
	if want.OwnerId != got.OwnerId {
		t.Fatalf("want OwnerId: %s, got OwnerId: %s", want.OwnerId, got.OwnerId)
	}
}

func TestDeleteProjectById(t *testing.T) {
	s, cleanUp := setupServiceLifetime(t)
	defer cleanUp("accounts", "projects")

	err := s.DeleteAccountById(context.TODO(), "")
	if err == nil {
		t.Fatal("should error on validation(incorrect id format): ", err)
	}

	id := uuid.NewString()
	err = s.DeleteAccountById(context.TODO(), id)
	if err == nil {
		t.Fatalf("should error on none existing record with id: %s, err: %s", id, err)
	}

	owner := types.Account{
		Id:    uuid.NewString(),
		Name:  "username",
		Email: "username@username.com",
	}
	_, err = s.DB.Exec("INSERT INTO accounts (id, email, name) VALUES ($1, $2, $3)", owner.Id, owner.Email, owner.Name)
	if err != nil {
		t.Fatal(ErrFailedToPrepareTest)
	}
	want := types.Project{
		Id:          uuid.NewString(),
		Name:        "project",
		Description: "A new project",
		OwnerId:     owner.Id,
		Deleted:     true,
	}

	_, err = s.DB.Exec("INSERT INTO projects (id, name, description, owner_id) VALUES ($1, $2, $3, $4)", want.Id, want.Name, want.Description, want.OwnerId)
	if err != nil {
		t.Fatal(ErrFailedToPrepareTest)
	}

	err = s.DeleteProjectById(context.TODO(), want.Id)
	if err != nil {
		t.Fatal(err)
	}

	var got types.Project
	row := s.DB.QueryRow("SELECT name, description, owner_id, deleted FROM projects LIMIT 1")
	if err != nil {
		t.Fatal(err)
	}
	err = row.Scan(&got.Name, &got.Description, &got.OwnerId, &got.Deleted)
	if err != nil {
		t.Fatal(err)
	}

	if err != nil {
		t.Fatal(err)
	}
	if want.Name != got.Name {
		t.Fatalf("want Name: %s, got Name: %s", want.Name, got.Name)
	}
	if want.Description != got.Description {
		t.Fatalf("want Description: %s, got Description: %s", want.Description, got.Description)
	}
	if want.OwnerId != got.OwnerId {
		t.Fatalf("want OwnerId: %s, got OwnerId: %s", want.OwnerId, got.OwnerId)
	}
	if want.Deleted != got.Deleted {
		t.Fatalf("want Deleted: %t, got Deleted: %t", want.Deleted, got.Deleted)
	}
}
