package service

import (
	"context"
	"errors"
	"testing"

	"github.com/danblok/pm/internals/types"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

func TestAddAccount(t *testing.T) {
	s, cleanUp := setupServiceLifetime(t)
	defer cleanUp("accounts")

	want := types.Account{
		Email: "username@gmail.com",
		Name:  "Username",
	}
	input := AddAccountInput{
		Name:  want.Name,
		Email: want.Email,
	}
	err := s.AddAccount(context.TODO(), &input)
	if err != nil {
		t.Fatal(err)
	}

	row := s.DB.QueryRow("SELECT email, name FROM accounts LIMIT 1")
	var got types.Account
	err = row.Scan(&got.Email, &got.Name)
	if err != nil {
		t.Fatalf("scan err: %s", err)
	}
	if err = row.Err(); err != nil {
		t.Fatalf("query err: %s", err)
	}

	if want.Email != got.Email {
		t.Fatalf("want.Email: %s != got.Email: %s", want.Email, got.Email)
	}
	if want.Name != got.Name {
		t.Fatalf("want.Name: %s != got.Name: %s", want.Name, got.Name)
	}

	input = AddAccountInput{
		Name:  "",
		Email: want.Email,
	}
	err = s.AddAccount(context.TODO(), &input)
	if !errors.Is(err, ErrFailedValidation) {
		t.Fatal("should error on incorrect Name field: ", err)
	}

	input = AddAccountInput{
		Name:  want.Name,
		Email: "",
	}
	err = s.AddAccount(context.TODO(), &input)
	if !errors.Is(err, ErrFailedValidation) {
		t.Fatal("should error on incorrect Email field: ", err)
	}
}

func TestGetAccountById(t *testing.T) {
	s, cleanUp := setupServiceLifetime(t)
	defer cleanUp("accounts")

	want := types.Account{
		Id:    uuid.NewString(),
		Email: "username@gmail.com",
		Name:  "username",
	}
	_, err := s.DB.Exec("INSERT INTO accounts (id, email, name) VALUES ($1, $2, $3)", want.Id, want.Email, want.Name)
	if err != nil {
		t.Fatal("prep test insertion error: ", err)
	}

	got, err := s.GetAccountById(context.TODO(), want.Id)
	if err != nil {
		t.Fatal(err)
	}
	if want.Id != got.Id {
		t.Fatalf("want.Id: %s != got.Id: %s", want.Id, got.Id)
	}
	if want.Email != got.Email {
		t.Fatalf("want.Email: %s != got.Email: %s", want.Email, got.Email)
	}
	if want.Name != got.Name {
		t.Fatalf("want.Name: %s != got.Name: %s", want.Name, got.Name)
	}

	_, err = s.GetAccountById(context.TODO(), "incorrect-id")
	if !errors.Is(err, ErrFailedValidation) {
		t.Fatal("should error on incorrect id err: ", err)
	}
}

func TestGetAllAccounts(t *testing.T) {
	s, cleanUp := setupServiceLifetime(t)
	defer cleanUp("accounts")

	want := []types.Account{
		{
			Id:    uuid.NewString(),
			Email: "username1@gmail.com",
			Name:  "username1",
		},
		{
			Id:    uuid.NewString(),
			Email: "username2@gmail.com",
			Name:  "username3",
		},
	}
	query := "INSERT INTO accounts (id, email, name) VALUES ($1, $2, $3), ($4, $5, $6)"
	_, err := s.DB.Exec(query, want[0].Id, want[0].Email, want[0].Name, want[1].Id, want[1].Email, want[1].Name)
	if err != nil {
		t.Fatal("prep test insertion error: ", err)
	}

	got, err := s.GetAllAccounts(context.TODO())
	if err != nil {
		t.Fatal(err)
	}
	if len(want) != len(got) {
		t.Fatalf("want len: %d, got len: %d", len(want), len(got))
	}
	for idx := range got {
		if want[idx].Id != got[idx].Id {
			t.Fatalf("want.Id: %s != got.Id: %s", want[idx].Id, got[idx].Id)
		}
		if want[idx].Email != got[idx].Email {
			t.Fatalf("want.Email: %s != got.Email: %s", want[idx].Email, got[idx].Email)
		}
		if want[idx].Name != got[idx].Name {
			t.Fatalf("want.Name: %s != got.Name: %s", want[idx].Name, got[idx].Name)
		}
	}
}

func TestUpdateAccount(t *testing.T) {
	s, cleanUp := setupServiceLifetime(t)
	defer cleanUp("accounts")

	acc := types.Account{
		Id:    uuid.NewString(),
		Email: "username@gmail.com",
		Name:  "username",
	}
	_, err := s.DB.Exec("INSERT INTO accounts (id, email, name) VALUES ($1, $2, $3)", acc.Id, acc.Email, acc.Name)
	if err != nil {
		t.Fatal("prep test insertion error: ", err)
	}

	want := types.Account{
		Id:    acc.Id,
		Email: acc.Email,
		Name:  "updated",
	}
	input := UpdateAccountInput{
		Id:   want.Id,
		Name: want.Name,
	}
	err = s.UpdateAccount(context.TODO(), &input)
	if err != nil {
		if errors.Is(err, ErrFailedToUpdate) {
			t.Fatal("no rows were affected: ", err)
		}
		t.Fatal("internal query err: ", err)
	}

	row := s.DB.QueryRow("SELECT id, email, name FROM accounts WHERE id=$1", want.Id)
	var got types.Account
	err = row.Scan(&got.Id, &got.Email, &got.Name)
	if err != nil {
		t.Fatalf("scan err: %s", err)
	}
	if err = row.Err(); err != nil {
		t.Fatalf("query err: %s", err)
	}

	if want.Id != got.Id {
		t.Fatalf("want.Id: %s != got.Id: %s", want.Id, got.Id)
	}
	if want.Email != got.Email {
		t.Fatalf("want.Email: %s != got.Email: %s", want.Email, got.Email)
	}
	if want.Name != got.Name {
		t.Fatalf("want.Name: %s != got.Name: %s", want.Name, got.Name)
	}

	input = UpdateAccountInput{
		Id:   "",
		Name: want.Name,
	}
	err = s.UpdateAccount(context.TODO(), &input)
	if !errors.Is(err, ErrFailedValidation) {
		t.Fatal("should error on incorrect id format validation: ", err)
	}

	input = UpdateAccountInput{
		Id:   uuid.NewString(),
		Name: want.Name,
	}
	err = s.UpdateAccount(context.TODO(), &input)
	if !errors.Is(err, ErrFailedToUpdate) {
		t.Fatalf("should error on none existing record with id: %s, err: %s", input.Id, err)
	}
}

func TestDeleteAccount(t *testing.T) {
	s, cleanUp := setupServiceLifetime(t)
	defer cleanUp("accounts")

	err := s.DeleteAccountById(context.TODO(), "")
	if !errors.Is(err, ErrFailedValidation) {
		t.Fatal("should error on validation(incorrect id format): ", err)
	}

	id := uuid.NewString()
	err = s.DeleteAccountById(context.TODO(), id)
	if !errors.Is(err, ErrFailedToUpdate) {
		t.Fatalf("should error on none existing record with id: %s, err: %s", id, err)
	}

	want := types.Account{
		Id:      uuid.NewString(),
		Email:   "username@gmail.com",
		Name:    "username",
		Deleted: true,
	}
	_, err = s.DB.Exec("INSERT INTO accounts (id, email, name) VALUES ($1, $2, $3)", want.Id, want.Email, want.Name)
	if err != nil {
		t.Fatal("prep test insertion error: ", err)
	}

	err = s.DeleteAccountById(context.TODO(), want.Id)
	if err != nil {
		if errors.Is(err, ErrFailedToUpdate) {
			t.Fatal("no rows were affected: ", err)
		}
		t.Fatal("internal query err: ", err)
	}

	got := new(types.Account)
	row := s.DB.QueryRow("SELECT id, email, name, deleted FROM accounts WHERE id::text=$1", want.Id)
	err = row.Scan(&got.Id, &got.Email, &got.Name, &got.Deleted)
	if err != nil {
		t.Fatal(err)
	}

	if want.Id != got.Id {
		t.Fatalf("want.Id: %s != got.Id: %s", want.Id, got.Id)
	}
	if want.Email != got.Email {
		t.Fatalf("want.Email: %s != got.Email: %s", want.Email, got.Email)
	}
	if want.Name != got.Name {
		t.Fatalf("want.Name: %s != got.Name: %s", want.Name, got.Name)
	}
	if want.Deleted != got.Deleted {
		t.Fatalf("want.Deleted: %t != got.Deleted: %t", want.Deleted, got.Deleted)
	}
}
