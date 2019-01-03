package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/lib/pq"
)

type message struct {
	Password string
}

func UploadFiles(rw http.ResponseWriter, r *http.Request) int {
	//establish db connection
	db, err := ConnectDB()
	if err != nil {
		fmt.Println("ERROR with getting DB", err)
		return http.StatusInternalServerError
	}
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
		return http.StatusInternalServerError
	}
	fmt.Println("successfully connected to db")

	r.ParseMultipartForm(32 << 20)
	mp := r.MultipartForm
	for _, file_headers := range mp.File {
		for _, header := range file_headers {
			file, err := header.Open()
			if err != nil {
				fmt.Println("error open: ", err)
				return http.StatusInternalServerError
			}
			defer file.Close()
			scanner := bufio.NewScanner(file)
			//this code works great. need to figure out how to account for duplicates bc
			//table will get ridiculously big if I don't

			txn, err := db.Begin()
			if err != nil {
				log.Fatal(err)
				return http.StatusInternalServerError
			}

			stmt, err := txn.Prepare(pq.CopyIn("common_passwords", "password"))
			if err != nil {
				log.Fatal(err)
				return http.StatusInternalServerError
			}

			for scanner.Scan() {
				entry := scanner.Text()
				_, err = stmt.Exec(entry)
				if err != nil {
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
		rw.WriteHeader(http.StatusUnauthorized)
		return http.StatusUnauthorized
	}
	//check if password in common_passwords table
	db, err := ConnectDB()
	if err != nil {
		fmt.Println("ERROR with getting DB", err)
		rw.WriteHeader(http.StatusUnauthorized)
		return http.StatusUnauthorized
	}

	rows, err := db.Query("select password from common_passwords where password = $1", pwd)
	if err != nil {
		fmt.Println("ERROR querying db", err)
		rw.WriteHeader(http.StatusUnauthorized)
		return http.StatusUnauthorized
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
		rw.WriteHeader(http.StatusUnauthorized)
		return http.StatusUnauthorized
	}

	rw.WriteHeader(http.StatusCreated)
	return http.StatusCreated
}
