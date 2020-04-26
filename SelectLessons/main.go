package main

import (
	"encoding/json"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var DB *gorm.DB

type Rstudent struct {
	Week    string
	Data    string
}

func START ()(NUM int,WEEK string){

	fmt.Printf("请输入学号")
	fmt.Scanf("%d",&NUM)
	fmt.Printf("请输入第几周")
	fmt.Scanf("%s",&WEEK)
	return
}


func Select(number int)(result string){
	tx := DB.Begin()
	err:=tx.Where("NUM=?",number).First(&result).Error
	if err != nil {
		tx.Rollback()
	}
	tx.Commit()
	return
}

func Serialization(result string,week string)(string) {
	res := &Rstudent{
		Week: week,
		Data: result,
	}
	data,err :=json.Marshal(res)
	if err != nil {
		fmt.Println(err)
	}
	strdata :=string(data)
	return strdata
}

func main() {
	db, err := gorm.Open("mysql", "root:root@(127.0.0.1:3306)/homework?charset=utf8mb4&parseTime=True&loc=Local")
	if err != nil {
		panic(err)
	}
	DB = db
	defer DB.Close()

	res1,res2:=START()
	res:=Select(res1)
	result :=Serialization(res,res2)
	fmt.Println(result)

}
