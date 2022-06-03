package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

// Product data type for export
type Product struct {
	ID          int
	Name        string
	Price       float32
	Description string
	Category    string
}

type Category struct {
	ID          int
	Name        string
	Description string
}

var tpl *template.Template

var db *sql.DB

func main() {
	tpl, _ = template.ParseGlob("templates/*.html")
	var err error
	db, err = sql.Open("mysql", "root:@tcp(localhost:3306)/testdb")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	http.HandleFunc("/insert", insertHandler)
	http.HandleFunc("/insertcategory", insertcategoryHandler)
	http.HandleFunc("/browse", browseHandler)
	http.HandleFunc("/allcategory", browseCategoryHandler)
	http.HandleFunc("/update/", updateHandler)
	http.HandleFunc("/updateresult/", updateResultHandler)
	http.HandleFunc("/updatecategory/", updatecategoryHandler)
	http.HandleFunc("/updateresultcategory/", updateCategoryResultHandler)
	http.HandleFunc("/delete/", deleteHandler)
	http.HandleFunc("/deletecategory/", deleteCategoryHandler)
	http.HandleFunc("/", homePageHandler)
	http.ListenAndServe("localhost:1212", nil)
}

func browseHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("*****browseHandler running*****")
	stmt := "CALL showproducts"
	rows, err := db.Query(stmt)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	var products []Product
	for rows.Next() {
		var p Product
		err = rows.Scan(&p.ID, &p.Name, &p.Price, &p.Description, &p.Category)
		if err != nil {
			panic(err)
		}
		products = append(products, p)
	}
	tpl.ExecuteTemplate(w, "select.html", products)
}
func browseCategoryHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("*****browseHandler running*****")
	s := "SELECT * FROM category"
	rowss, err := db.Query(s)
	if err != nil {
		panic(err)
	}
	defer rowss.Close()
	var categorys []Category
	for rowss.Next() {
		var c Category
		err = rowss.Scan(&c.ID, &c.Name, &c.Description)
		if err != nil {
			panic(err)
		}
		categorys = append(categorys, c)
	}

	tpl.ExecuteTemplate(w, "selectcategory.html", categorys)
}

func insertHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("*****insertHandler running*****")
	if r.Method == "GET" {
		tpl.ExecuteTemplate(w, "insert.html", nil)
		return
	}
	r.ParseForm()
	name := r.FormValue("nameName")
	price := r.FormValue("priceName")
	descr := r.FormValue("descrName")
	categoryid := r.FormValue("category_id")
	var err error
	if name == "" || price == "" || descr == "" {
		fmt.Println("Error inserting row:", err)
		tpl.ExecuteTemplate(w, "insert.html", "Error inserting data, please check all fields.")
		return
	}
	var ins *sql.Stmt
	ins, err = db.Prepare("INSERT INTO `testdb`.`product` (`name`, `price`, `description`,`category_id`) VALUES (?, ?, ?,?);")
	if err != nil {
		panic(err)
	}
	defer ins.Close()
	res, err := ins.Exec(name, price, descr, categoryid)

	rowsAffec, _ := res.RowsAffected()
	if err != nil || rowsAffec != 1 {
		fmt.Println("Error inserting row:", err)
		tpl.ExecuteTemplate(w, "insert.html", "Error inserting data, please check all fields.")
		return
	}
	lastInserted, _ := res.LastInsertId()
	rowsAffected, _ := res.RowsAffected()
	fmt.Println("ID of last row inserted:", lastInserted)
	fmt.Println("number of rows affected:", rowsAffected)
	tpl.ExecuteTemplate(w, "insert.html", "Product Successfully Inserted")
}

func insertcategoryHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("*****insertHandler running*****")
	if r.Method == "GET" {
		tpl.ExecuteTemplate(w, "insertcategory.html", nil)
		return
	}
	r.ParseForm()
	name := r.FormValue("nameName")
	descr := r.FormValue("descrName")
	var err error
	if name == "" || descr == "" {
		fmt.Println("Error inserting row:", err)
		tpl.ExecuteTemplate(w, "insertcategory.html", "Error inserting data, please check all fields.")
		return
	}
	var ins *sql.Stmt
	ins, err = db.Prepare("INSERT INTO `testdb`.`category` (`name`, `description`) VALUES (?, ?);")
	if err != nil {
		panic(err)
	}
	defer ins.Close()
	res, err := ins.Exec(name, descr)

	rowsAffec, _ := res.RowsAffected()
	if err != nil || rowsAffec != 1 {
		fmt.Println("Error inserting row:", err)
		tpl.ExecuteTemplate(w, "insertcategory.html", "Error inserting data, please check all fields.")
		return
	}
	lastInserted, _ := res.LastInsertId()
	rowsAffected, _ := res.RowsAffected()
	fmt.Println("ID of last row inserted:", lastInserted)
	fmt.Println("number of rows affected:", rowsAffected)
	tpl.ExecuteTemplate(w, "insertcategory.html", "Category Successfully Inserted")
}

func updateHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("*****updateHandler running*****")
	r.ParseForm()
	id := r.FormValue("idproducts")
	row := db.QueryRow("SELECT * FROM testdb.product WHERE Id = ?;", id)
	var p Product
	// func (r *Row) Scan(dest ...interface{}) error
	err := row.Scan(&p.ID, &p.Name, &p.Price, &p.Description, &p.Category)
	if err != nil {
		fmt.Println(err)
		http.Redirect(w, r, "/browse", 307)
		return
	}
	tpl.ExecuteTemplate(w, "update.html", p)
}

func updatecategoryHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("*****updateHandler running*****")
	r.ParseForm()
	id := r.FormValue("idproducts")
	row := db.QueryRow("SELECT * FROM testdb.category WHERE Id = ?;", id)
	var p Product
	err := row.Scan(&p.ID, &p.Name, &p.Description)
	if err != nil {
		fmt.Println(err)
		http.Redirect(w, r, "/allcategory", 307)
		return
	}
	tpl.ExecuteTemplate(w, "updatecategory.html", p)
}

func updateResultHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("*****updateResultHandler running*****")
	r.ParseForm()
	id := r.FormValue("idproducts")
	name := r.FormValue("nameName")
	price := r.FormValue("priceName")
	description := r.FormValue("descrName")
	categoryid := r.FormValue("category_id")
	upStmt := "UPDATE `testdb`.`product` SET `Name` = ?, `Price` = ?, `Description` = ? , `category_id` = ? WHERE (`Id` = ?);"
	stmt, err := db.Prepare(upStmt)
	if err != nil {
		fmt.Println("error preparing stmt")
		panic(err)
	}
	fmt.Println("db.Prepare err:", err)
	fmt.Println("db.Prepare stmt:", stmt)
	defer stmt.Close()
	var res sql.Result
	res, err = stmt.Exec(name, price, description, categoryid, id)
	rowsAff, _ := res.RowsAffected()
	if err != nil || rowsAff != 1 {
		fmt.Println(err)
		tpl.ExecuteTemplate(w, "result.html", "There was a problem updating the product")
		return
	}
	http.Redirect(w, r, "/browse", 301)
}

func updateCategoryResultHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("*****updateResultHandler running*****")
	r.ParseForm()
	id := r.FormValue("idproducts")
	name := r.FormValue("nameName")
	description := r.FormValue("descrName")
	upStmt := "UPDATE `testdb`.`category` SET `Name` = ?, `Description` = ? WHERE (`Id` = ?);"
	stmt, err := db.Prepare(upStmt)
	if err != nil {
		fmt.Println("error preparing stmt")
		panic(err)
	}
	fmt.Println("db.Prepare err:", err)
	fmt.Println("db.Prepare stmt:", stmt)
	defer stmt.Close()
	var res sql.Result
	res, err = stmt.Exec(name, description, id)
	rowsAff, _ := res.RowsAffected()
	if err != nil || rowsAff != 1 {
		fmt.Println(err)
		tpl.ExecuteTemplate(w, "result.html", "There was a problem updating the Category")
		return
	}
	http.Redirect(w, r, "/allcategory", 301)
}

func deleteHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("*****deleteHandler running*****")
	r.ParseForm()
	id := r.FormValue("idproducts")
	del, err := db.Prepare("DELETE FROM `testdb`.`product` WHERE (`Id` = ?);")
	if err != nil {
		panic(err)
	}
	defer del.Close()
	var res sql.Result
	res, err = del.Exec(id)
	rowsAff, _ := res.RowsAffected()
	fmt.Println("rowsAff:", rowsAff)

	if err != nil || rowsAff != 1 {
		fmt.Fprint(w, "Error deleting product")
		return
	}
	fmt.Println("err:", err)
	http.Redirect(w, r, "/browse", 301)
}

func deleteCategoryHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("*****deleteHandler running*****")
	r.ParseForm()
	id := r.FormValue("idproducts")
	del, err := db.Prepare("DELETE FROM `testdb`.`category` WHERE (`id` = ?);")
	if err != nil {
		panic(err)
	}
	defer del.Close()
	var res sql.Result
	res, err = del.Exec(id)
	rowsAff, _ := res.RowsAffected()
	fmt.Println("rowsAff:", rowsAff)

	if err != nil || rowsAff != 1 {
		fmt.Fprint(w, "Error deleting Category")
		return
	}
	fmt.Println("err:", err)
	http.Redirect(w, r, "/allcategory", 301)

}

func homePageHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/browse", 307)
}
