package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

type Database interface {
	FindAccounts() ([]*Account, error)
	FindAccountById(int64) (*Account, error)
	CreateAccount(*Account) error
	UpdateAccount(*Account) (*Account, error)
	DeleteAccount(int64) error
}

type Postgres struct {
	db *sql.DB
}

func NewPostgres() (*Postgres, error) {
	connStr := "user=postgres password=root dbname=gobank sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &Postgres{db: db}, nil
}

func (p *Postgres) Init() error {
	return p.CreateAccountTable()
}

func (p *Postgres) CreateAccountTable() error {
	query := `CREATE TABLE IF NOT EXISTS accounts(
    			id serial primary key,
    			first_name varchar,
    			last_name varchar,
    			number integer,
    			balance integer,
    			created_at timestamptz
    )`
	_, err := p.db.Exec(query)
	return err
}

func (p *Postgres) FindAccounts() ([]*Account, error) {
	query := "SELECT * FROM accounts"
	rows, err := p.db.Query(query)
	if err != nil {
		return nil, err
	}

	var accounts []*Account

	for rows.Next() {
		account := new(Account)
		errScan := rows.Scan(
			&account.ID,
			&account.FirstName,
			&account.LastName,
			&account.Number,
			&account.Balance,
			&account.CreatedAt,
		)
		if err = errScan; err != nil {
			return nil, err
		}
		accounts = append(accounts, account)
	}

	return accounts, nil
}

func (p *Postgres) FindAccountById(accoudId int64) (*Account, error) {
	query := "SELECT * FROM accounts WHERE id = $1"
	row := p.db.QueryRow(query, accoudId)

	account := new(Account)

	errScan := row.Scan(
		&account.ID,
		&account.FirstName,
		&account.LastName,
		&account.Number,
		&account.Balance,
		&account.CreatedAt,
	)
	if err := errScan; err != nil {
		return nil, err
	}
	return account, nil
}

func (p *Postgres) CreateAccount(account *Account) error {
	query := `INSERT INTO accounts 
    (
    	first_name, 
		last_name, 
		number, 
		balance, 
		created_at
    ) VALUES ($1, $2, $3, $4, $5)`
	result, err := p.db.Exec(query, account.FirstName, account.LastName, account.Number, account.Balance, account.CreatedAt)
	if err != nil {
		return err
	}

	fmt.Printf("%+v\n\n", result)
	return nil
}

func (p *Postgres) UpdateAccount(account *Account) (*Account, error) {
	query := `UPDATE accounts SET 
        first_name = $2,
		last_name = $3,
		balance = $4,
		created_at = $5 
		WHERE id = $1
		RETURNING *`

	row := p.db.QueryRow(query, account.ID, account.FirstName, account.LastName, account.Balance, account.CreatedAt)

	updatedAccount := new(Account)
	errScan := row.Scan(
		&updatedAccount.ID,
		&updatedAccount.FirstName,
		&updatedAccount.LastName,
		&updatedAccount.Number,
		&updatedAccount.Balance,
		&updatedAccount.CreatedAt,
	)
	if err := errScan; err != nil {
		return nil, err
	}
	return updatedAccount, nil
}

func (p *Postgres) DeleteAccount(id int64) error {
	query := "DELETE FROM accounts WHERE id = $1"
	_, err := p.db.Exec(query, id)
	if err != nil {
		return err
	}
	return nil
}
