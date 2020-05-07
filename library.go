package main

import (
	"fmt"

	// mysql connector
	_ "github.com/go-sql-driver/mysql"
	sqlx "github.com/jmoiron/sqlx"
)

const (
	User     = "root"
	Password = "123456"
	DBName   = "ass3"
)

type Library struct {
	db *sqlx.DB
}

func (lib *Library) ConnectDB() {
	db, err := sqlx.Open("mysql", fmt.Sprintf("%s:%s@tcp(127.0.0.1:3306)/%s", User, Password, DBName))
	if err != nil {
		panic(err)
	}
	lib.db = db
}

func Execute(lib *Library, states []string) {
	for _, s := range states {
		_, err := lib.db.Exec(s)
		if err != nil {
			panic(err)
		}
	}
}

// CreateTables created the tables in MySQL
func (lib *Library) CreateTables() error {
	
	s1 := fmt.Sprintf("DROP TABLE IF EXISTS BOOKS")
	s2 := fmt.Sprintf("DROP TABLE IF EXISTS STUS")
	s3 := fmt.Sprintf("DROP TABLE IF EXISTS ADMINS")
	s4 := fmt.Sprintf("DROP TABLE IF EXISTS DELETE_REC")
	s5 := fmt.Sprintf("DROP TABLE IF EXISTS BORROW_REC")	
	Execute(lib, []string{s5,s4,s2,s1,s3})

	s1 = fmt.Sprintf("CREATE TABLE IF NOT EXISTS BOOKS(ISBN char(32) NOT NULL, author char(32), title char(100), total int, avail int, PRIMARY KEY(ISBN))")
	s2 = fmt.Sprintf("CREATE TABLE IF NOT EXISTS STUS(UID char(32) NOT NULL, password char(32) NOT NULL, PRIMARY KEY(UID))")
	s3 = fmt.Sprintf("CREATE TABLE IF NOT EXISTS ADMINS(UID char(32) NOT NULL, password char(32) NOT NULL, PRIMARY KEY(UID))")
	s4 = fmt.Sprintf("CREATE TABLE IF NOT EXISTS DELETE_REC(REC int NOT NULL, ISBN char(32) NOT NULL, explanation char(100) NOT NULL, delete_date DATE NOT NULL, PRIMARY KEY(REC))")
	s5 = fmt.Sprintf("CREATE TABLE IF NOT EXISTS BORROW_REC(REC int NOT NULL, UID char(32) NOT NULL, ISBN char(32) NOT NULL, start DATE NOT NULL, exp DATE, ret DATE, EXtimes int, PRIMARY KEY(REC), FOREIGN KEY(UID) REFERENCES STUS(UID), FOREIGN KEY(ISBN) REFERENCES BOOKS(ISBN))")
	Execute(lib, []string{s1,s2,s3,s4,s5})
	// fmt.Println("Successfully created the tables.")
	s0 := fmt.Sprintf("INSERT INTO ADMINS(UID, password) VALUES('root', '123456')")
	Execute(lib, []string{s0})
	return nil
}

func (lib *Library) CreateForTest() error{
	s1 := fmt.Sprintf("INSERT INTO BOOKS(ISBN, author, title, total, avail) VALUES('5678-5-6', 'Daye Xue', 'Bin_Cat', 4, 3), ('3690-5-6', 'QWQ', 'Bin_Dog', 1, 0), ('9570-5-6', 'Alpha', 'How to Debug', 1, 0)")
	s2 := fmt.Sprintf("INSERT INTO STUS(UID, password) VALUES('16302345678','123321')")
	s3 := fmt.Sprintf("INSERT INTO BORROW_REC(REC, UID, ISBN, start, exp, EXtimes) VALUES(1,'16302345678','5678-5-6', '2020-02-01','2020-04-30', 3), (2,'16302345678','3690-5-6', '2020-03-01','2020-04-15', 0), (3,'16302345678','9570-5-6', '2020-02-28','2020-04-29', 0)")
	Execute(lib, []string{s1, s2, s3})
	return nil
}

// AddBook add a book into the library
func (lib *Library) AddBook(title, author, ISBN string) (int, error) {
	var s string
	var total int
	s = fmt.Sprintf("SELECT total FROM BOOKS WHERE ISBN = '%s'", ISBN)
	rows, err := lib.db.Query(s)
	if err != nil {
		panic(err)
	}
	if rows.Next() {
		err := rows.Scan(&total)
		if err != nil {
			panic(err)
		}
		s1:= fmt.Sprintf("UPDATE BOOKS SET total=total+1 WHERE ISBN = '%s'", ISBN)
		s2:= fmt.Sprintf("UPDATE BOOKS SET avail=avail+1 WHERE ISBN = '%s'", ISBN)
		Execute(lib, []string{s1,s2})
		fmt.Println("Successfully added the book.")
	} else {
		s = fmt.Sprintf("INSERT INTO BOOKS(ISBN, author, title, total, avail) VALUES('%s', '%s', '%s', 1, 1)", ISBN, author, title)
		Execute(lib, []string{s})
		fmt.Println("Successfully added the book.")
		total = 0
	}

	return total, nil
}

func (lib *Library) DeleteBook(ISBN string, explanation string) (int, error) {
	var s string
	var total, avail, state int
	s = fmt.Sprintf("SELECT total, avail FROM BOOKS WHERE ISBN = '%s'", ISBN)
	rows, err := lib.db.Query(s)
	if err != nil {
		panic(err)
	}
	// max_rec := 1
	if rows.Next() {
		rows.Scan(&total, &avail)
		if err != nil {
			panic(err)
		}
		if avail > 0 {
			var max_rec int
			s9 := fmt.Sprintf("SELECT COUNT(*) FROM DELETE_REC WHERE REC IS NOT NULL")
			rows, err := lib.db.Query(s9)
			if err != nil {
				panic(err)
				max_rec = 0
			} else{
				var max_rec int
				if rows.Next() {
					err := rows.Scan(&max_rec)
					if err != nil {
						panic(err)
					}
				} else {
					max_rec = 0
				}
			}

			s1 := fmt.Sprintf("UPDATE BOOKS SET total = total-1 , avail = avail-1 WHERE ISBN = '%s'", ISBN)
			s2 := fmt.Sprintf("INSERT DELETE_REC(REC,ISBN,explanation,delete_date) VALUES('%d', '%s', '%s', CURRENT_DATE())",max_rec+1, ISBN, explanation)
			Execute(lib, []string{s1,s2})
			fmt.Println("Successfully deleted the book.")
			state = 1
		} else {
			fmt.Println("The books have been borrowed. If lost, the student should buy a new one instead.")
			state = 2
		}
	} else {
		fmt.Println("Can't find this book.")
		state = 0
	}

	return state, nil
}

// etc...
func (lib *Library) AddStudent(UID, code string) (int, error) {
	var s string
	var state int
	s = fmt.Sprintf("SELECT * FROM STUS WHERE UID = '%s'", UID)
	rows, err := lib.db.Query(s)
	if err != nil {
		panic(err)
	}
	if rows.Next() {
		fmt.Println("Error: UID already exists!")
		state = 0

	} else {
		s1 := fmt.Sprintf("INSERT INTO STUS(UID, password) VALUES('%s', '%s')", UID, code)
		Execute(lib, []string{s1})
		fmt.Println("Successfully added the student!")
		state = 1
	}
	return state, nil
}

func (lib *Library) QueryBook(value, key string) (string, error) {
	var ISBN, author, title string
	var total, avail int
	var s, state string
	s = fmt.Sprintf("SELECT ISBN, author, title, total, avail FROM BOOKS WHERE %s = '%s'", key, value)
	rows, err := lib.db.Query(s)
	if err != nil {
		panic(err)
	}
	count := 0
	for rows.Next() {
		err := rows.Scan(&ISBN, &author, &title, &total, &avail)
		if err != nil {
			panic(err)
		}
		count += 1
		s1 := fmt.Sprintf("Find books No.%d: ISBN=%s author=%s title=%s total=%d avail=%d", count, ISBN, author, title, total, avail)
		fmt.Println(s1)
		state += ISBN + " " + author + " " + title + "\n"
		// fmt.Println(state)
	}
	if count == 0 {
		fmt.Printf("Can't find a book whose %s is %s\n", key, value)
		state = ""
	}
	return state, nil
}

func (lib *Library) BorrowBook(UID, ISBN string) (int, error) {
	var s string
	var state int
	s = fmt.Sprintf("SELECT * FROM STUS WHERE UID = '%s'", UID)
	rows, err := lib.db.Query(s)
	if err != nil {
		panic(err)
	}
	if rows.Next() {
		s1 := fmt.Sprintf("SELECT COUNT(*) FROM BORROW_REC WHERE UID = '%s' AND exp<CURRENT_DATE() AND ISNULL(ret) ", UID)
		rows, err := lib.db.Query(s1)
		if err != nil {
			panic(err)
		}
		var count int
		if rows.Next() {
			err := rows.Scan(&count)
			if err != nil {
				panic(err)
			}
		} else {
			count = 0
		}
		if count >= 3 {
			fmt.Println("Account suspended. Please return books first.")
			state = 1
		} else {
			s2 := fmt.Sprintf("SELECT avail FROM BOOKS WHERE ISBN ='%s'", ISBN)
			rows, err := lib.db.Query(s2)
			if err != nil {
				panic(err)
			}
			if rows.Next() {
				var avail int 
				err := rows.Scan(&avail)
				if err != nil {
					panic(err)
				}
				if avail == 0 {
					fmt.Println("No book is available now.")
					state = 3
				} else {
					// max_b_rec := 4
					s9 := fmt.Sprintf("SELECT COUNT(*) FROM BORROW_REC")
					rows, err := lib.db.Query(s9)
					if err != nil {
						panic(err)
					}
					var max_b_rec int
					if rows.Next() {
						err := rows.Scan(&max_b_rec)
						if err != nil {
							panic(err)
						}
					} else {
						max_b_rec = 0
					}

					//(SELECT MAX(REC)+1 AS REC FROM BORROW_REC)AS B
					s3 := fmt.Sprintf("UPDATE BOOKS SET avail=avail-1 WHERE ISBN = '%s'", ISBN)
					s4 := fmt.Sprintf("INSERT INTO BORROW_REC(REC, UID, ISBN, start, exp, EXtimes) VALUES('%d', '%s', '%s', CURRENT_DATE(), date_add(CURRENT_DATE(), interval 60 day), 0)", max_b_rec+1, UID, ISBN)
					Execute(lib, []string{s3, s4})
					fmt.Println("Successfully borrow the book!")
					state = 4
				}
			} else {
				fmt.Println("ISBN error! Can't find the book.")
				state = 2
			}

		}

	} else {
		fmt.Println("UID error!")
		state = 0
	}
	return state, nil
}

func (lib *Library) QueryHistory(UID string) (string, error) {
	var s string
	var state string
	s = fmt.Sprintf("SELECT * FROM BORROW_REC WHERE UID = '%s'", UID )
	rows, err := lib.db.Query(s)
	if err != nil {
		panic(err)
	}

	var u_id, ISBN, start, exp, ret string
	var rec, EXtimes int

	count := 0
	for rows.Next() {
		rows.Scan(&rec, &u_id, &ISBN, &start, &exp, &ret, &EXtimes)
		// err := rows.Scan(&rec, &u_id, &ISBN, &start, &exp, &ret, &EXtimes)
		// if err != nil {
		// 	panic(err)
		// }
		if ret == "" {
			ret = "NULL"
		}
		count += 1
		s1 := fmt.Sprintf("Find borrow_records No.%d: UID=%s ISBN=%s start=%s expected=%s return=%s Ext_times=%d", count, UID, ISBN, start, exp, ret, EXtimes)
		fmt.Println(s1)
		state += s1 + "\n"
		// fmt.Println(state)
	}
	if count == 0 {
		fmt.Printf("No borrow record! of UID %s.\n", UID)
		state = ""
	}

	return state, nil
}
 
func (lib *Library) QueryBooksNotReturned(UID string) (string, error) {
	var s string
	var state string
	s = fmt.Sprintf("SELECT rec, ISBN, start, exp, ret FROM BORROW_REC WHERE UID = '%s' AND ret is NULL", UID )
	rows, err := lib.db.Query(s)
	if err != nil {
		panic(err)
	}

	var ISBN, start, exp, ret string
	var rec int
	count := 0
	for rows.Next() {
		rows.Scan(&rec, &ISBN, &start, &exp, &ret)
		if ret == "" {
			ret = "NULL"
		}
		count += 1
		s1 := fmt.Sprintf("Find books not returned No.%d: UID=%s ISBN=%s start=%s expected=%s return=%s", count, UID, ISBN, start, exp, ret)
		fmt.Println(s1)
		state += s1 + "\n"
		// fmt.Println(state)
	}
	if count == 0 {
		fmt.Printf("ALL books have returned by UID %s.\n", UID)
		state = ""
	}

	return state, nil
}

func (lib *Library) QueryDueDate(UID, ISBN string) (string, error) {
	var s string
	var state string
	s = fmt.Sprintf("SELECT rec, exp FROM BORROW_REC WHERE UID = '%s' AND ISBN = '%s' AND ret is NULL", UID, ISBN )
	rows, err := lib.db.Query(s)
	if err != nil {
		panic(err)
	}

	var exp string
	var rec int

	count := 0
	for rows.Next() {
		rows.Scan(&rec, &exp)
		// err := rows.Scan(&rec, &u_id, &ISBN, &start, &exp, &ret, &EXtimes)
		// if err != nil {
		// 	panic(err)
		// }
		count += 1
		s1 := fmt.Sprintf("The deadline of borrow_records No.%d: ISBN=%s expected=%s", count, ISBN, exp)
		fmt.Println(s1)
		state += s1 + "\n"
		// fmt.Println(state)
	}
	if count == 0 {
		fmt.Printf("No deadline of UID %s.\n", UID)
		state = ""
	}

	return state, nil
}

func (lib *Library) ExtendDueDate(REC int, UID, ISBN string) (int, error) {
	var s string
	var state int
	s = fmt.Sprintf("SELECT EXtimes FROM BORROW_REC WHERE REC = %d AND UID = '%s' AND ISBN = '%s' AND ret is NULL", REC, UID, ISBN )
	rows, err := lib.db.Query(s)
	if err != nil {
		panic(err)
	}
	var EXtimes int
	if rows.Next() {
		rows.Scan(&EXtimes)
		if EXtimes == 3 {
			fmt.Println("Can't extend the due!")
			state = 0
		} else {
			s1 := fmt.Sprintf("UPDATE BORROW_REC SET EXtimes = EXtimes+1, exp = date_add(exp, interval 14 day) WHERE REC = %d AND UID = '%s' AND ISBN = '%s'",REC, UID, ISBN)
			Execute(lib, []string{s1})
			fmt.Println("Successfully extend the due!")
			state = 1
		}
	}
	return state, nil
}

func (lib *Library) QueryBooksOverdued(UID string) (string, error) {
	var s string
	var state string
	s = fmt.Sprintf("SELECT rec, ISBN, exp FROM BORROW_REC WHERE UID = '%s' AND ret is NULL AND exp < CURRENT_DATE()", UID )
	rows, err := lib.db.Query(s)
	if err != nil {
		panic(err)
	}

	var ISBN, exp string
	var rec int
	count := 0
	for rows.Next() {
		rows.Scan(&rec, &ISBN, &exp)
		count += 1
		s1 := fmt.Sprintf("Find books overdued No.%d: UID=%s ISBN=%s expected=%s", count, UID, ISBN, exp)
		fmt.Println(s1)
		state += s1 + "\n"
		// fmt.Println(state)
	}
	if count == 0 {
		fmt.Printf("No book is overdued by UID %s.\n", UID)
		state = ""
	}

	return state, nil
}

func (lib *Library) ReturnBook(REC int, UID,ISBN string) (int, error) {
	var s string
	var state int
	s = fmt.Sprintf("SELECT * FROM BORROW_REC WHERE REC = %d AND UID = '%s' AND ISBN = '%s' AND ret is NULL", REC, UID, ISBN )
	rows, err := lib.db.Query(s)
	if err != nil {
		panic(err)
	}
	
	if rows.Next() {
		s1 := fmt.Sprintf("UPDATE BORROW_REC SET ret = CURRENT_DATE() WHERE REC = %d AND UID = '%s' AND ISBN = '%s' AND ret is NULL", REC, UID, ISBN)
		Execute(lib, []string{s1})
		fmt.Println("Successfully returned the book!")
			state = 1
	} else {
		fmt.Printf("No such a book borrowed by REC=%d, UID=%s ISBN=%s!\n", REC, UID, ISBN)
		state = 0
	}
	
	return state, nil
}

func (lib *Library) CheckSTU(UID, code string) bool {
	var s string
	flag :=false
	s = fmt.Sprintf("SELECT password FROM STUS WHERE UID = '%s'", UID)
	rows, err := lib.db.Query(s)
	if err != nil {
		panic(err)
	}
	if rows.Next() {
		var password string
		rows.Scan(&password)
		if code == password {
			flag = true
		}
	}
	return flag
}

func (lib *Library) CheckADMIN(UID, code string) bool {
	var s string
	flag :=false
	s = fmt.Sprintf("SELECT password FROM ADMINS WHERE UID = '%s'", UID)
	rows, err := lib.db.Query(s)
	if err != nil {
		panic(err)
	}
	if rows.Next() {
		var password string
		rows.Scan(&password)
		if code == password {
			flag = true
		}
	}
	return flag
}

func main() {
	fmt.Println("Welcome to the Library Management System!")
	lib := new(Library)
	lib.ConnectDB()
	lib.CreateTables()
	var input_s, id, code string
	for {
		fmt.Println("Please Select User Mode: (input number)")
		fmt.Println("1: Student; 2: Administrator; 0: Exit.")
		fmt.Scanln(&input_s)
		if input_s == "0" {
			break
		} else if input_s == "1" {
			fmt.Println("Please enter UID:")
			fmt.Scanln(&id)
			fmt.Println("Please enter password:")
			fmt.Scanln(&code)
			var flag bool
			flag = lib.CheckSTU(id, code)
			// flag = true
			if !flag {
				fmt.Println("Login failed.")
			}
			for flag {
				fmt.Println("Please select an operation: (input number)")
				fmt.Println("1:QueryBook; 2:BorrowBook; 3:ReturnBook; 4:QueryHistory; 5:QueryDueDate; 6:ExtendDueDate; 7:QueryBooksOverdued; 8:QueryBooksNotReturned; 0:Exit.")
				var input_s2 string
				fmt.Scanln(&input_s2)
				if input_s2=="0" {
					break
				} else {
					switch input_s2{
						case "1":
							fmt.Println("Please select the input type:")
							fmt.Println("1: title; 2: author; 3: ISBN.")
							var typ string
							fmt.Scanln(&typ)
							if typ=="1" || typ=="2" || typ=="3" {
								fmt.Println("Please enter the value:")
								var value string
								fmt.Scanln(&value)
								if typ=="1" {
									lib.QueryBook(value, "title")
								} else if typ=="2" {
									lib.QueryBook(value, "author")
								} else {
									lib.QueryBook(value, "ISBN")
								}
							} else {
								fmt.Println("Enter illegal characters!")
							}
						case "2":
							var ISBN string
							fmt.Println("Please enter ISBN:")
							fmt.Scanln(&ISBN)
							lib.BorrowBook(id, ISBN)
						case "3":
							var ISBN string
							var REC int
							fmt.Println("Please enter ISBN:")
							fmt.Scanln(&ISBN)
							fmt.Println("Please enter the id of Borrow_Record:")
							fmt.Scanln(&REC)
							lib.ReturnBook(REC, id, ISBN)
						case "4":
							var STU_UID string
							fmt.Println("Please enter STU_UID:")
							fmt.Scanln(&STU_UID)
							lib.QueryHistory(STU_UID)
						case "5":
							var ISBN string
							fmt.Println("Please enter ISBN:")
							fmt.Scanln(&ISBN)
							lib.QueryDueDate(id, ISBN)
						case "6":
							var ISBN string
							var REC int
							fmt.Println("Please enter ISBN:")
							fmt.Scanln(&ISBN)
							fmt.Println("Please enter the id of Borrow_Record:")
							fmt.Scanln(&REC)
							lib.ReturnBook(REC, id, ISBN)
						case "7":
							lib.QueryBooksOverdued(id)
						case "8":
							lib.QueryBooksNotReturned(id)
						default:
							fmt.Println("Enter illegal characters!")
					}
				}
			}
		} else if input_s == "2" {
			fmt.Println("Please enter AdminID:")
			fmt.Scanln(&id)
			fmt.Println("Please enter password:")
			fmt.Scanln(&code)
			var flag bool
			flag = lib.CheckADMIN(id, code)
			// flag = true
			if !flag {
				fmt.Println("Login failed.")
			}
			for flag {
				fmt.Println("Please select an operation: (input number)")
				fmt.Println("1: AddBook; 2: DeleteBook; 3: AddStudent; 4:QueryBook; 5:QueryHistory; 0: Exit.")
				var input_s2 string
				fmt.Scanln(&input_s2)
				if input_s2=="0" {
					break
				} else {
					switch input_s2{
						case "1":
							var ISBN, title, author string
							fmt.Println("Please enter book_ISBN:")
							fmt.Scanln(&ISBN)
							for ISBN == "" {
								fmt.Println("Please enter book_ISBN:")
								fmt.Scanln(&ISBN)
							}
							fmt.Println("Please enter book_title:")
							fmt.Scanln(&title)
							fmt.Println("Please enter book_author:")
							fmt.Scanln(&author)
							lib.AddBook(title, author, ISBN)
						case "2":
							var ISBN, explanation string
							fmt.Println("Please enter book_ISBN:")
							fmt.Scanln(&ISBN)
							fmt.Println("Please enter delete_book_explanation:")
							fmt.Scanln(&explanation)
							lib.DeleteBook(ISBN, explanation)
						case "3":
							var STU_UID, pass string
							fmt.Println("Please enter STU_UID:")
							fmt.Scanln(&STU_UID)
							fmt.Println("Please enter STU_password:")
							fmt.Scanln(&pass)
							lib.AddStudent(STU_UID, pass)
						case "4":
							fmt.Println("Please select the input type:")
							fmt.Println("1: title; 2: author; 3: ISBN.")
							var typ string
							fmt.Scanln(&typ)
							if typ=="1" || typ=="2" || typ=="3" {
								fmt.Println("Please enter the value:")
								var value string
								fmt.Scanln(&value)
								if typ=="1" {
									lib.QueryBook(value, "title")
								} else if typ=="2" {
									lib.QueryBook(value, "author")
								} else {
									lib.QueryBook(value, "ISBN")
								}
							} else {
								fmt.Println("Enter illegal characters!")
							}
						case "5":
							var STU_UID string
							fmt.Println("Please enter STU_UID:")
							fmt.Scanln(&STU_UID)
							lib.QueryHistory(STU_UID)
						default:
							fmt.Println("Enter illegal characters!")
					}
				}
			}
		} else {
			fmt.Println("Enter illegal characters!")
		}
	}

	lib.db.Close()

}