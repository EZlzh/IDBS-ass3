# README

## 1.Tips:

### go.mod:

```
github.com/go-sql-driver/mysql v1.5.0
```

modify the version suit for you.



### library.go:

modify the User and Password for your mysql.

```
const (
	User     = "root"
	Password = "123456"
	DBName   = "ass3"
)
```



## 2.Run the code:

##### First you should create a DATABASE whose name is the same as the DBName above.



##### Then At the root of this folder：

**1.** **for running:**

```
go run library.go
```

And you'll see this:

```
Welcome to the Library Management System!
Please Select User Mode: (input number)
1: Student; 2: Administrator; 0: Exit.
```



**2. for testing:**

```
go test
```

And you'll see this:

```
Successfully added the book.
...
...
...
PASS
ok  	github.com/ichn-hu/IDBS-Spring20-Fudan/assignments/ass3/boilerplate	x.xxxs

```









