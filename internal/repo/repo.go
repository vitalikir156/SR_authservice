package repo

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
	"github.com/vitalikir156/SR_authservice/internal/config"
	interrors "github.com/vitalikir156/SR_authservice/internal/errors"
	"github.com/vitalikir156/SR_authservice/internal/types"
)

var (
	ErrTaskNotFound   = errors.New("object with ID not found")
	ErrFreeIDNotFound = errors.New("free ID not found")
	ErrBadStatus      = errors.New("Invalid value for status field")
	ErrEmptyName      = errors.New("name field is empty")
)

//const taskstatusNew = "new"

type repository struct {
	pool *pgxpool.Pool
}

type Repository interface {
	CreateUser(ctx context.Context,login string, password string) (int, error)
	GetUserOverID(ctx context.Context, id int) (types.User, error)
	GetUserOverLogin(ctx context.Context, login string) (types.User, error)
	UpdateUserPassword(ctx context.Context, login string, oldpass string, newpass string) error
	DeleteUser(ctx context.Context, login string) error

	CreateToken(ctx context.Context, token types.Token) (int, error)
	GetToken(ctx context.Context, userid int) (types.Token, error)
	DelToken(ctx context.Context, tokenid int) (error)
	UpdateToken(ctx context.Context, token types.Token) (error)
	//CheckToken(ctx context.Context, userid int) (types.Token, error)
}

func NewRepository(ctx context.Context, conf config.PostgreSQL) (Repository, error) {
	connString := fmt.Sprintf(
		`user=%s password=%s host=%s port=%d dbname=%s sslmode=%s
        pool_max_conns=%d pool_max_conn_lifetime=%s pool_max_conn_idle_time=%s`,
		conf.User,
		conf.Password,
		conf.Host,
		conf.Port,
		conf.Name,
		conf.SSLMode,
		conf.PoolMaxConns,
		conf.PoolMaxConnLifetime.String(),
		conf.PoolMaxConnIdleTime.String(),
	)

	config, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse PostgreSQL config")
	}

	config.ConnConfig.DefaultQueryExecMode = pgx.QueryExecModeCacheDescribe

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create PostgreSQL connection pool")
	}

	return &repository{pool}, nil
}

//TODO: use transaction
func (r *repository) CreateUser(ctx context.Context, login string, password string) (int, error) {
	if len(login) == 0 {
		return -1, ErrEmptyName
	}

	tx, err :=	r.pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.Serializable})
	if err != nil {
		return -1, err
	}
	defer func() {
		tx.Rollback(ctx)
	}()

	err = tx.QueryRow(ctx, "SELECT from users where uname = $1", login).Scan()

	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows){return -1, err}
	} else
	{return -1, interrors.ErrBusyLogin}

	query := "INSERT INTO users (uname, password) VALUES ($1, $2) RETURNING id"
	var id int
	err = tx.QueryRow(ctx, query, login, password).Scan(&id)
	if err != nil {
		return -1, errors.Wrap(err, "failed to insert user")
	}
	err = tx.Commit(ctx)
	if err != nil {
		return -1, err
	}
	return id, nil
}

func (r *repository) GetUserOverID(ctx context.Context, id int) (types.User, error) {
	var user types.User
	query := "SELECT id, uname, taskread, taskwrite, userread, userwrite, password from users where id = $1"
	err := r.pool.QueryRow(ctx, query, id).Scan(&user.ID, &user.Name, &user.Taskread,
		&user.Taskwrite, &user.Userread, &user.Userwrite, &user.Password)
	if err != nil {
		return types.User{}, errors.Wrap(err, "failed to get user")
	}

	return user, nil
}

func (r *repository) GetUserOverLogin(ctx context.Context, login string) (types.User, error) {
	var user types.User
	query := "SELECT id, uname, taskread, taskwrite, userread, userwrite, password from users where uname = $1"
	err := r.pool.QueryRow(ctx, query, login).Scan(&user.ID, &user.Name, &user.Taskread,
		&user.Taskwrite, &user.Userread, &user.Userwrite, &user.Password)
	if err != nil {
		return types.User{}, errors.Wrap(err, "failed to get user")
	}

	return user, nil
}


//TODO: check old password
func (r *repository) UpdateUserPassword(ctx context.Context, login string, oldpass string, newpass string) error {
	query := "UPDATE users SET password = $1 where uname=$2"
	out, err := r.pool.Exec(ctx, query, newpass, login)
	if err != nil {
		return errors.Wrap(err, "update fault")
	}
	if out.RowsAffected() < 1 {
		return ErrTaskNotFound
	}
	return nil
}

func (r *repository) DeleteUser(ctx context.Context, login string) error {
	query := "DELETE FROM users where uname=$1"
	out, err := r.pool.Exec(ctx, query, login)
	if err != nil {
		return errors.Wrap(err, "delete fault")
	}
	if out.RowsAffected() < 1 {
		return ErrTaskNotFound
	}

	return nil
}

func (r *repository) CreateToken(ctx context.Context, token types.Token) (int, error) {
	query := "INSERT INTO tokens (user_id, token) VALUES ($1, $2) RETURNING id"
	var id int
	err := r.pool.QueryRow(ctx, query, token.UID, token.Token).Scan(&id)
	if err != nil {
		return -1, errors.Wrap(err, "failed to insert token")
	}
	return id, nil
}

func (r *repository) GetToken(ctx context.Context, userid int) (types.Token, error) {
	var token types.Token
	query := "SELECT id, user_id, token, created_at from tokens where user_id = $1"
	err := r.pool.QueryRow(ctx, query, userid).Scan(&token.ID, &token.UID, &token.Token,
		&token.Created)
	if err != nil {
		return types.Token{}, errors.Wrap(err, "failed to get token")
	}

	return token, nil
}
func (r *repository) UpdateToken(ctx context.Context, token types.Token) (error) {
	query := "UPDATE tokens SET token = $1 where user_id = $2"
	out, err := r.pool.Exec(ctx, query, token.Token, token.UID)
	if err != nil {
		errors.Wrap(err, "failed to update token")
	}
	if out.RowsAffected() != 1 {
		return ErrTaskNotFound
	}

	return nil
}
func (r *repository) CheckToken(ctx context.Context, userid int) (types.Token, error) {
	var token types.Token
	query := "SELECT id, user_id, token, created_at from tokens where user_id = $1"
	err := r.pool.QueryRow(ctx, query, userid).Scan(&token.ID, &token.UID, &token.Token,
		&token.Created)
	if err != nil {
		return types.Token{}, errors.Wrap(err, "failed to get token")
	}

	return token, nil
}
func (r *repository) DelToken(ctx context.Context, tokenid int) (error) {
	query := "DELETE FROM tokens where id=$1"
	out, err := r.pool.Exec(ctx, query, tokenid)
	if err != nil {
		return errors.Wrap(err, "delete fault")
	}
	if out.RowsAffected() < 1 {
		return ErrTaskNotFound
	}

	return nil
}