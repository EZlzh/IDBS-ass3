package main

import (
	"testing"
)

func TestCreateTables(t *testing.T) {
	lib := Library{}
	lib.ConnectDB()
	err := lib.CreateTables()
	if err != nil {
		t.Errorf("can't create tables")
	}
}

func TestCreateForTest(t *testing.T) {
	lib := Library{}
	lib.ConnectDB()
	err := lib.CreateForTest()
	if err != nil {
		t.Errorf("can't create tables")
	}
}

func TestAddBook(t *testing.T) {
	lib := Library{}
	lib.ConnectDB()


	type test_book struct {
		title string 
		author string 
		ISBN string
		total int
	}

	books := []test_book{}
	QWQ := 10
	for i:=0; i<QWQ; i++ {
		// fmt.Println("%d", i)
		books = append(books, test_book{"Cheerful_And_Humorous_Talk", "Peipei", "1234-5-6", i})
	}

	for i:=0; i<QWQ; i++ {
		// fmt.Println("%d", i)
		books = append(books, test_book{"DST_master", "Xiao as", "9999-5-6", i})
	}

	for _,book := range books {
		total, err := lib.AddBook(book.title, book.author, book.ISBN)
		if err != nil {
			t.Errorf("AddBook Failed.")
		}
		if total != book.total {
			t.Errorf("Count total number of a book failed, received: %d, expected: %d ",total, book.total);
		}
	}

}

func TestDeleteBook(t *testing.T) {
	lib := Library{}
	lib.ConnectDB()


	type test_book struct {
		ISBN string
		explanation string
		state int
	}
	books := []test_book{}
	books = append(books, test_book{"1234-5-6", "When Peipei was playing magic, he burnt this book.", 1})
	books = append(books, test_book{"1234-7-6", "When As was playing DST, he burnt this book.", 0})

	for _,book := range books {
		state, err := lib.DeleteBook(book.ISBN, book.explanation)
		if err != nil {
			t.Errorf("DeleteBook Failed.")
		}
		if state != book.state {
			t.Errorf("Delete state failed, received: %d, expected: %d ",state, book.state);
		}
	}
}

func TestAddStudent(t *testing.T) {
	lib := Library{}
	lib.ConnectDB()

	type test_stu struct {
		UID string
		password string
		state int
	}
	stus := []test_stu{}
	stus = append(stus, test_stu{"18307777777","123456",1})
	stus = append(stus, test_stu{"18307777777","123456",0})
	stus = append(stus, test_stu{"18306666333","123456",1})
	stus = append(stus, test_stu{"18306666333","123456",0})

	for _,stu := range stus {
		state, err := lib.AddStudent(stu.UID, stu.password)
		if err != nil {
			t.Errorf("AddStudent Failed.")
		}
		if state != stu.state {
			t.Errorf("AddStudent state failed, received: %d, expected: %d ",state, stu.state);
		}
	}
}

func TestQueryBook(t *testing.T) {
	lib := Library{}
	lib.ConnectDB()

	type test_book struct {
		key string
		value string
		state string
	}
	books := []test_book{}
	books = append(books, test_book{"ISBN", "9999-5-6", "9999-5-6 Xiao as DST_master\n"})
	books = append(books, test_book{"author", "Peipei", "1234-5-6 Peipei Cheerful_And_Humorous_Talk\n"})

	for _,book := range books {
		state, err := lib.QueryBook(book.value, book.key)
		if err != nil {
			t.Errorf("QueryBook Failed.")
		}
		if state != book.state {
			t.Errorf("Querybook state failed, received: %s, expected: %s ",state, book.state);
		}
	}
}

func TestBorrowBook(t *testing.T) {
	lib := Library{}
	lib.ConnectDB()

	type test_borrow struct {
		UID string
		ISBN string
		state int
	}
	borrows := []test_borrow{}
	borrows = append(borrows, test_borrow{"18307777776", "9999-4-6", 0})
	borrows = append(borrows, test_borrow{"18307777777", "3690-5-6", 3})
	borrows = append(borrows, test_borrow{"18307777777", "3699-5-6", 2})
	borrows = append(borrows, test_borrow{"18307777777", "5678-5-6", 4})
	borrows = append(borrows, test_borrow{"18307777777", "5678-5-6", 4})
	borrows = append(borrows, test_borrow{"16302345678", "9999-5-6", 1})

	for _,borrow := range borrows {
		state, err := lib.BorrowBook(borrow.UID, borrow.ISBN)
		if err != nil {
			t.Errorf("BorrowBook Failed.")
		}
		if state != borrow.state {
			t.Errorf("Borrowbook state failed, received: %d, expected: %d ",state, borrow.state);
		}
	}
}

func TestQueryHistory(t *testing.T) {
	lib := Library{}
	lib.ConnectDB()

	type test_book struct {
		UID string
		state string
	}
	books := []test_book{}
	books = append(books, test_book{"16302345678", "Find borrow_records No.1: UID=16302345678 ISBN=5678-5-6 start=2020-02-01 expected=2020-04-30 return=NULL Ext_times=0\nFind borrow_records No.2: UID=16302345678 ISBN=3690-5-6 start=2020-03-01 expected=2020-04-15 return=NULL Ext_times=0\nFind borrow_records No.3: UID=16302345678 ISBN=9570-5-6 start=2020-02-28 expected=2020-04-29 return=NULL Ext_times=0\n"})
	// books = append(books, test_book{"author", "Peipei", "1234-5-6 Peipei Cheerful_And_Humorous_Talk\n"})

	for _,book := range books {
		state, err := lib.QueryHistory(book.UID)
		if err != nil {
			t.Errorf("QueryHistory Failed.")
		}
		if state != book.state {
			t.Errorf("QueryHistory state failed, received: %s, expected: %s ",state, book.state);
		}
	}
}

func TestQueryBooksNotReturned(t *testing.T) {
	lib := Library{}
	lib.ConnectDB()

	type test_book struct {
		UID string
		state string
	}
	books := []test_book{}
	books = append(books, test_book{"16302345678", "Find books not returned No.1: UID=16302345678 ISBN=5678-5-6 start=2020-02-01 expected=2020-04-30 return=NULL\nFind books not returned No.2: UID=16302345678 ISBN=3690-5-6 start=2020-03-01 expected=2020-04-15 return=NULL\nFind books not returned No.3: UID=16302345678 ISBN=9570-5-6 start=2020-02-28 expected=2020-04-29 return=NULL\n"})
	// books = append(books, test_book{"author", "Peipei", "1234-5-6 Peipei Cheerful_And_Humorous_Talk\n"})

	for _,book := range books {
		state, err := lib.QueryBooksNotReturned(book.UID)
		if err != nil {
			t.Errorf("QueryBooksNotReturned Failed.")
		}
		if state != book.state {
			t.Errorf("QueryBooksNotReturned state failed, received: %s, expected: %s ",state, book.state);
		}
	}
}

func TestQueryDueDate(t *testing.T) {
	lib := Library{}
	lib.ConnectDB()

	type test_book struct {
		UID string
		ISBN string
		state string
	}
	books := []test_book{}
	books = append(books, test_book{"16302345678", "5678-5-6", "The deadline of borrow_records No.1: ISBN=5678-5-6 expected=2020-04-30\n"})
	// books = append(books, test_book{"author", "Peipei", "1234-5-6 Peipei Cheerful_And_Humorous_Talk\n"})

	for _,book := range books {
		state, err := lib.QueryDueDate(book.UID, book.ISBN)
		if err != nil {
			t.Errorf("QueryDueDate Failed.")
		}
		if state != book.state {
			t.Errorf("QueryDueDate state failed, received: %s, expected: %s ",state, book.state);
		}
	}
}

func TestExtendDueDate(t *testing.T) {
	lib := Library{}
	lib.ConnectDB()

	type test_book struct {
		REC int
		UID string
		ISBN string
		state int
	}
	books := []test_book{}
	books = append(books, test_book{4,"18307777777", "5678-5-6", 1})
	books = append(books, test_book{1,"16302345678", "5678-5-6", 0})
	
	for _,book := range books {
		state, err := lib.ExtendDueDate(book.REC, book.UID, book.ISBN)
		// fmt.Printf("%d",book.REC)
		if err != nil {
			t.Errorf("ExtendDueDate Failed.")
		}
		if state != book.state {
			t.Errorf("ExtendDueDate state failed, received: %d, expected: %d ",state, book.state);
		}
	}
}

func TestQueryBooksOverdued(t *testing.T) {
	lib := Library{}
	lib.ConnectDB()

	type test_book struct {
		UID string
		state string
	}
	books := []test_book{}
	books = append(books, test_book{"18307777777", ""})
	books = append(books, test_book{"16302345678", "Find books overdued No.1: UID=16302345678 ISBN=5678-5-6 expected=2020-04-30\nFind books overdued No.2: UID=16302345678 ISBN=3690-5-6 expected=2020-04-15\nFind books overdued No.3: UID=16302345678 ISBN=9570-5-6 expected=2020-04-29\n"})
	// books = append(books, test_book{"author", "Peipei", "1234-5-6 Peipei Cheerful_And_Humorous_Talk\n"})

	for _,book := range books {
		state, err := lib.QueryBooksOverdued(book.UID)
		if err != nil {
			t.Errorf("QueryBooksOverdued Failed.")
		}
		if state != book.state {
			t.Errorf("QueryBooksOverdued state failed, received: %s, expected: %s ",state, book.state);
		}
	}
}

func TestReturnBook(t *testing.T) {
	lib := Library{}
	lib.ConnectDB()

	type test_book struct {
		REC int
		UID string
		ISBN string
		state int
	}
	books := []test_book{}
	books = append(books, test_book{4,"18307777777", "5678-5-6", 1})
	books = append(books, test_book{1,"16302345678", "5678-5-5", 0})
	
	for _,book := range books {
		state, err := lib.ReturnBook(book.REC, book.UID, book.ISBN)
		// fmt.Printf("%d",book.REC)
		if err != nil {
			t.Errorf("ReturnBook Failed.")
		}
		if state != book.state {
			t.Errorf("ReturnBook state failed, received: %d, expected: %d ",state, book.state);
		}
	}
}