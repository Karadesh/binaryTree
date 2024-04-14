package main

import (
	"bufio"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
	"log"
	"os"
)

type NumbersMap struct {
	TreeNumber int
}

func main() {
	d := fmt.Sprintf("(node: %d, value: %d)", 25, 50)
	k := fmt.Sprintf("(node: %d, value: %d)", 34, 25)
	var a []string
	a = append(a, d)
	a = append(a, k)
	fmt.Println(a)
	//db, err := sql.Open("postgres", "user=admin password=admin host=localhost dbname=randnumbers sslmode=disable")
	//if err != nil {
	//	log.Fatalf("failed to connect: %v", err)
	//}
	//defer db.Close()
	//_, err = db.Exec("TRUNCATE TABLE binarytree")
	//if err != nil {
	//	log.Fatalf("failed to drop binaryTree: %v", err)
	//}
	//var m map[int]NumbersMap
	//selector := fmt.Sprintf("SELECT * FROM binarytree")
	//somerow, err := db.Query(selector)
	//if err != nil {
	//	log.Fatalf("failed to query: %v", err)
	//}
	//m = make(map[int]NumbersMap)
	//defer somerow.Close()
	//for somerow.Next() {
	//	var someNumber int
	//	var someNode int
	//	something := somerow.Scan(&someNumber, &someNode)
	//	m[someNode] = NumbersMap{someNumber}

	//	fmt.Println(someNumber)
	//	fmt.Println(someNode)
	//	fmt.Sprintf("somenumber: %d; somenode: %d", someNumber, someNode)
	//	fmt.Println(something)
	//}
	//fmt.Println(m)
	//fmt.Println(somerow)

	//a := "1 2 3 4"
	//b := "3 2 1 0"
	//d := DbInserter(a)
	//fmt.Sprintf("db inserter: %v", d)
	//c := DbUpdater(b)
	//fmt.Sprintf("db updater: %v", c)
	//v := DbGetter()
	//if v == "" {
	//fmt.Println("db get fail")
	//} else {
	//	fmt.Sprintf("v = (%s)", v)
	//}
}

//db, err := sql.Open("postgres", "user=admin password=admin host=localhost dbname=randnumbers sslmode=disable")
//if err != nil {
//	log.Fatalf("failed to connect: %v", err)
//}
//defer db.Close()

//result, err := db.Exec("UPDATE randomnumbers SET number_value =$0 WHERE number_id = $1 ")
//if err != nil {
//d := "0"
//inserter := fmt.Sprintf("INSERT INTO randomnumbers(number_id, number_value) VALUES(0, %s)", d)
//updater := fmt.Sprintf("UPDATE randomnumbers SET number_value =%s WHERE number_id = 0", d)
//result, err := db.Exec(updater)
//if err != nil {
//	fmt.Printf("failed to insert into randomnumbers: %v", err)
//}
//affectedRows, err := result.RowsAffected()
//fmt.Printf("Updated %d rows\n", affectedRows)
//selector := fmt.Sprintf("SELECT number_value FROM randomnumbers WHERE number_id=0")
//somerow, err := db.Query(selector)
//if err != nil {
//	log.Fatalf("failed to query: %v", err)
//}
//defer somerow.Close()
//for somerow.Next() {
//	var somevalue string
//	something := somerow.Scan(&somevalue)
//	fmt.Println(somevalue)
//	fmt.Println(something)
//}
//fmt.Println(somerow)

func DbInserter(numbers string) string {
	db, err := sql.Open("postgres", "user=admin password=admin host=localhost dbname=randnumbers sslmode=disable")
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	defer db.Close()
	inserter := fmt.Sprintf("INSERT INTO randomnumbers(number_id, number_value) VALUES(0, '%s')", numbers)
	result, err := db.Exec(inserter)
	if err != nil {
		fmt.Printf("failed to insert into randomnumbers: %v", err)
	}
	affectedRows, err := result.RowsAffected()
	log.Printf("Inserted %d rows\n", affectedRows)
	if err != nil {
		log.Fatalf("failed to query: %v", err)
	}
	return "Done!"
}

func DbUpdater(numbers string) string {
	db, err := sql.Open("postgres", "user=admin password=admin host=localhost dbname=randnumbers sslmode=disable")
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	defer db.Close()
	updater := fmt.Sprintf("UPDATE randomnumbers SET number_value = '%s' WHERE number_id = 0", numbers)
	result, err := db.Exec(updater)
	if err != nil {
		fmt.Printf("failed to update randomnumbers: %v", err)
	}
	affectedRows, errAff := result.RowsAffected()
	if errAff != nil {
		log.Fatalf("failed to query: %v", err)
	}
	log.Printf("Updated %d rows\n", affectedRows)
	return "Done!"
}

func DbGetter() string {
	db, err := sql.Open("postgres", "user=admin password=admin host=localhost dbname=randnumbers sslmode=disable")
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	defer db.Close()
	selector := fmt.Sprintf("SELECT number_value FROM randomnumbers WHERE number_id=0")
	someRow, err := db.Query(selector)

	if err != nil {
		return ""
	}
	var someValue string
	someValue = ""
	for someRow.Next() {
		something := someRow.Scan(&someValue)
		fmt.Println(something)
	}
	return someValue
}

func decryptNumbers(r string) string {
	log.Printf("Please type numbers to check encrypted numbers")
	var answer string
	in := bufio.NewScanner(os.Stdin)
	in.Scan()
	answer = in.Text()
	fmt.Println(answer)
	bytesObject := []byte(r)
	err := bcrypt.CompareHashAndPassword(bytesObject, []byte(answer))
	if err == nil {
		return "True"
	} else {
		return "False"
	}
}

func generateHash(strObject string) string {
	bytesObject := []byte(strObject)
	encryptedObject, _ := bcrypt.GenerateFromPassword(bytesObject, bcrypt.MinCost)
	return string(encryptedObject)
}
