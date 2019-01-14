package main

import (
	"bufio"
	"encoding/json"
	"log"
	"net/http"

	"github.com/lib/pq"
)

type message struct {
	Password string
}

func ProcessRequest(rw http.ResponseWriter, r *http.Request) int {

	headerType := r.Header.Get("Content-type")
	switch headerType {
	case "application/json":
		return CheckPassword(rw, r)
	case "multipart/form-data":
		return UploadFiles(rw, r)
	}
	return http.StatusBadRequest
}

func UploadFiles(rw http.ResponseWriter, r *http.Request) int {
	//establish db connection
	db, err := ConnectDB()
	if err != nil {
		log.Printf("Could not connect to db: %v\n", err)
		return http.StatusInternalServerError
	}
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
		return http.StatusInternalServerError
	}
	r.ParseMultipartForm(32 << 20)
	mp := r.MultipartForm
	for _, file_headers := range mp.File {
		for _, header := range file_headers {
			file, err := header.Open()
			if err != nil {
				log.Fatal(err)
				return http.StatusInternalServerError
			}

			defer file.Close()

			scanner := bufio.NewScanner(file)
			txn, err := db.Begin()
			if err != nil {
				log.Fatal(err)
				return http.StatusInternalServerError
			}
			stmt, err := txn.Prepare(pq.CopyIn("common_passwords", "password"))
			if err != nil {
				if rollbackErr := txn.Rollback(); rollbackErr != nil {
					log.Printf("Could not roll back: %v\n", rollbackErr)
				}
				log.Fatal(err)
				return http.StatusInternalServerError
			}

			for scanner.Scan() {
				entry := scanner.Text()
				_, err = stmt.Exec(entry)
				if err != nil {
					if rollbackErr := txn.Rollback(); rollbackErr != nil {
						log.Printf("Could not roll back: %v\n", rollbackErr)
					}
					log.Fatal(err)
					return http.StatusInternalServerError
				}
			}
			if err := scanner.Err(); err != nil {
				log.Fatal(err)
				return http.StatusInternalServerError
			}
			_, err = stmt.Exec()
			if err != nil {
				if rollbackErr := txn.Rollback(); rollbackErr != nil {
					log.Printf("Could not roll back: %v\n", rollbackErr)
				}
				log.Fatal(err)
				return http.StatusInternalServerError
			}
			err = stmt.Close()
			if err != nil {
				log.Fatal(err)
				return http.StatusInternalServerError
			}
			err = txn.Commit()
			if err != nil {
				log.Fatal(err)
				return http.StatusInternalServerError
			}
		}
	}
	return http.StatusOK
}
func CheckPassword(rw http.ResponseWriter, r *http.Request) int {
	decoder := json.NewDecoder(r.Body)
	var m message
	err := decoder.Decode(&m)

	if err != nil {
		panic(err)
	}
	//check length of password to make sure > 8 minimum
	pwd := m.Password
	if len(pwd) < 8 {
		rw.WriteHeader(http.StatusForbidden)
		return http.StatusForbidden
	}
	//check if password in common_passwords table
	db, err := ConnectDB()
	if err != nil {
		log.Printf("Could not connect to db: %v\n", err)
		rw.WriteHeader(http.StatusForbidden)
		return http.StatusForbidden
	}

	rows, err := db.Query("select password from common_passwords where password = $1", pwd)
	if err != nil {
		log.Fatal(err)
		rw.WriteHeader(http.StatusForbidden)
		return http.StatusForbidden
	}

	pwdExists := make([]string, 0)
	defer rows.Close()
	for rows.Next() {
		var ret string
		if err := rows.Scan(&ret); err != nil {
			log.Fatal(err)
			return http.StatusInternalServerError
		}
		pwdExists = append(pwdExists, ret)
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
		return http.StatusInternalServerError
	}

	if len(pwdExists) > 0 {
		rw.WriteHeader(http.StatusForbidden)
		return http.StatusForbidden
	}

	rw.WriteHeader(http.StatusCreated)
	return http.StatusCreated
}
