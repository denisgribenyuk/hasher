package hash_repository

import (
	"database/sql"

	"client/models"
)

type HashRepo struct {
	DB *sql.DB
}
type HashedString struct {
	Hash string `db:"hash"`
	ID   int64  `db:"id"`
}

func New(con *sql.DB) *HashRepo {
	return &HashRepo{DB: con}
}

func (h *HashRepo) GetHashedString(id int) (*HashedString, error) {
	res := &HashedString{}
	rows, err := h.DB.Query("SELECT * FROM hash WHERE id = $1", id)
	if err != nil {
		return nil, err
	}

	r := rows.Next()
	if !r {
		return nil, nil
	}
	err = rows.Scan(&res.ID, &res.Hash)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (h *HashRepo) WriteHashedString(hs *models.Hash) (*HashedString, error) {
	query, err := h.DB.Prepare("INSERT INTO hash (hash) VALUES ($1) RETURNING id, hash")
	if err != nil {
		return nil, err
	}
	res := &HashedString{}
	err = query.QueryRow(hs.Hash).Scan(&res.ID, &res.Hash)
	if err != nil {
		return nil, err
	}
	return res, err
}
