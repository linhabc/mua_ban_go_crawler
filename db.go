package main

import (
	"fmt"

	"github.com/syndtr/goleveldb/leveldb"
)

func createOrOpenDb(path string) *leveldb.DB {
	db, _ := leveldb.OpenFile(path, nil)
	return db
}

func getData(db *leveldb.DB, key string) []byte {
	data, _ := db.Get([]byte(key), nil)
	return data
}

func putData(db *leveldb.DB, key string, data string) error {
	err := db.Put([]byte(key), []byte(data), nil)
	return err
}

func mainTest() {

	db1 := createOrOpenDb("./db/cat1")
	db2 := createOrOpenDb("./db/cat2")

	defer db1.Close()
	defer db2.Close()

	_ = putData(db1, "id1", "sdt1")
	_ = putData(db1, "id2", "sdt2")

	_ = putData(db2, "id3", "sdt3")
	_ = putData(db2, "id4", "sdt4")

	data1, _ := db1.Get([]byte("id1"), nil)
	data := getData(db1, "id2")

	fmt.Printf("db1 value: %s\n", data1)
	fmt.Printf("db1 value: %s\n", data)

	data1, _ = db2.Get([]byte("id3"), nil)
	data = getData(db2, "id4")

	fmt.Printf("db2 value: %s\n", data1)
	fmt.Printf("db2 value: %s\n", data)
}

//  level db ko co muti thread
//  nhan vao id -> co hay ko
//  2 routine, 1 read + 1 write

// store key -> id
//       value -> phone number

// err = db.Delete([]byte("key"), nil)
