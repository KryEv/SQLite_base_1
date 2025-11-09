package main

import (
	"database/sql"
)

type ParcelStore struct {
	db *sql.DB
}

func NewParcelStore(db *sql.DB) ParcelStore {
	return ParcelStore{db: db}
}

func (s ParcelStore) Add(p Parcel) (int, error) {

	// добавление строки в таблицу parcel
	res, err := s.db.Exec("INSERT INTO parcel (client, status, address, created_at) VALUES (:client, :status, :address, :created_at)",
		sql.Named("client", p.Client),
		sql.Named("status", p.Status),
		sql.Named("address", p.Address),
		sql.Named("created_at", p.CreatedAt))

	// проверка на наличие ошибки
	if err != nil {
		return 0, err
	}

	// идентификатор последней добавленной записи в базу
	id, err := res.LastInsertId()

	// проверка на наличие ошибки
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (s ParcelStore) Get(number int) (Parcel, error) {

	// запрос на чтение из базы
	rows := s.db.QueryRow("SELECT number, client, status, address, created_at FROM parcel WHERE number = :number",
		sql.Named("number", number))

	p := Parcel{}

	// парсим строку ответа
	err := rows.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt)

	// проверка на наличие ошибки
	if err != nil {
		return Parcel{}, err
	}

	return p, nil
}

func (s ParcelStore) GetByClient(client int) ([]Parcel, error) {

	var res []Parcel

	// чтение строк из таблицы parcel по заданному client
	rows, err := s.db.Query("SELECT number, client, status, address, created_at FROM parcel WHERE client = :client",
		sql.Named("client", client))

	// проверка на наличие ошибки
	if err != nil {
		return []Parcel{}, err
	}

	// отложенное закрытие
	defer rows.Close()

	for rows.Next() {

		p := Parcel{}
		err := rows.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt)

		// проверка на наличие ошибки
		if err != nil {
			return []Parcel{}, err
		}

		res = append(res, p)
	}

	err = rows.Err()
	// проверка на наличие ошибки
	if err != nil {
		return []Parcel{}, err
	}

	return res, nil
}

func (s ParcelStore) SetStatus(number int, status string) error {

	// обновление статуса в таблице parcel
	_, err := s.db.Exec("UPDATE parcel SET status = :status WHERE number = :number",
		sql.Named("status", status),
		sql.Named("number", number))

	// проверка на наличие ошибки
	if err != nil {
		return err
	}

	return nil
}

func (s ParcelStore) SetAddress(number int, address string) error {

	// обновление адреса в таблице parcel
	_, err := s.db.Exec("UPDATE parcel SET address = :address WHERE number = :number AND status = :status",
		sql.Named("address", address),
		sql.Named("number", number),
		sql.Named("status", ParcelStatusRegistered))

	// проверка на наличие ошибки
	if err != nil {
		return err
	}

	return nil
}

func (s ParcelStore) Delete(number int) error {

	// удаление строки из таблицы parcel
	_, err := s.db.Exec("DELETE FROM parcel WHERE number = :number AND status = :status",
		sql.Named("number", number),
		sql.Named("status", ParcelStatusRegistered))

	// проверка на наличие ошибки
	if err != nil {
		return err
	}

	return nil
}
