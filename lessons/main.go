package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"strings"
)
type lesson struct {
	WeekTime string
	ClassTime string
	Name string
	NUM string
	Loction string
	Teacher string
	Type string
	Point string
}

type Student struct {
	gorm.Model
	Name string
	Number  string   `gorm:"index:id_idx"`
	Data  []string
}

var DB *gorm.DB
func init() {
	db, err := gorm.Open("mysql", "root:root@(localhost)/homework?charset=utf8mb4&parseTime=True&loc=Local")
	if err!= nil{
		panic(err)
	}
	DB =db
}

//爬取网页
func spiderST(Url string)(bodystr string,err error){
	//设置client
	//useragent
	useragent := "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.163 Safari/537.36"
	//ssl证书
	tr := &http.Transport{
		TLSClientConfig:&tls.Config{InsecureSkipVerify:true},

	}
	client := http.Client{Transport:tr}
	req, err := http.NewRequest("GET", Url, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("User-Agent", useragent)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return
	}
	bodystr = string(body)
	return

}




//正则匹配
func RElessons(data string)(results []string,name string){
    //姓名
	re :=regexp.MustCompile(`<li>〉〉2019-2020学年2学期 学生课表>>\d{10}(?s:(.*?))</li>`)
	res :=re.FindAllStringSubmatch(data,-1)
	name =res[0][1]

	n :=0  //自加用
	results = make([]string,4096)
	//1.找课程  找出课程相关的字符串 全部 并且按照时间分好条
	re1 :=regexp.MustCompile(`<td style='font-weight:bold;'>(?s:(.*?))</tr>`)
	res1 :=re1.FindAllStringSubmatch(data,-1)
	//fmt.Println(res1)
	//把时间和课程分开
	for i,_ :=range res1 { //i 就是第几节课
		//找时间并把一大条课程分开 按照时间分条
		re2 := regexp.MustCompile(`<td style='font-weight:bold;'>(?s:(.*?))</td>`) //找时间的正则
		res2 := re2.FindAllStringSubmatch(res1[i][0], -1)                          // res2[0][1]为第几节时间 "第xx节"
		//fmt.Println(res2)
		//把每个时间段的课程分开
		re3 := regexp.MustCompile(`<td >{1}(?s:(.*?))</td>{1}`)
		res3 := re3.FindAllStringSubmatch(res1[i][0], -1) //res3[k][1] 每个时间段的每条课程  ”[[<td ></td> ] [<td ></td> ] [<td ></td> ] [<td ></td> ] [<td ></td> ] [<td ></td> ] [<td ></td> ]]“
		//fmt.Println("res3",res3)
		//找到每条课程的元素
			for k, _ := range res3 { //res3[k][1] 每个时间段的每条课程
				//判断是否没有课
				if res3[k][1] == "" {

					lesson0 := lesson{
						WeekTime:  "",
						ClassTime: res2[0][1],
						Name:      "0",
						NUM:       "",
						Loction:   "",
						Teacher:   "",
						Type:      "",
						Point:     "",
					}
					newlesson0 ,_:= json.Marshal(lesson0)

					results[n] = string(newlesson0)
					n++
					//fmt.Println(string(newlesson0))

				} else {
					//如果有课   拆分课程元素
					re4 := regexp.MustCompile(`<div class='kbTd' zc='\d{20}'>(?s:(.*?))<br>(?s:(.*?))<br>(?s:(.*?))<br>(?s:(.*?))<font color=#FF0000>.*?</font><br><span style='color:#0000FF'>(?s:(.*?))</span>`)
					res4 := re4.FindAllStringSubmatch(res3[k][0], -1)
					//fmt.Println("res4[0]",res4[0])    //res4[0]
					//把课程编号和名称分开
					result2 := strings.Split(res4[0][2], "-") //result[0]为编号 [1]为名称
					//地点修正
					result3 := strings.Replace(res4[0][3], "地点：", "", -1)
					//fmt.Println(result3)
					//分离教师 类型 学分
					result5 := strings.Split(res4[0][5], " ")
					//fmt.Println(result5)
					//绑定结构体
					lesson1 := lesson{
						WeekTime:  res4[0][4],
						ClassTime: res2[0][1],
						Name:      result2[1],
						NUM:       result2[0],
						Loction:   result3,
						Teacher:   result5[0],
						Type:      result5[1],
						Point:     result5[2],
					}
					newlesson1 ,_:= json.Marshal(lesson1)
					results[n] = string(newlesson1)
					n++
					//fmt.Println(string(newlesson1))

				}
			}
		}

    return

}


//包装
func PCK (data []string,NUM int,name string)(Student){
	Student :=Student{
		Name:   name,
		Number: string(NUM),
		Data:   data,
	}
	return Student
}

//传到数据库
func sql (student Student){

	// 自动迁移
	DB.AutoMigrate(&Student{})

	DB.Create(&student)


}

func main (){
	for i:=201921001;i<=2019215203;i++{
		URL :="http://jwc.cqupt.edu.cn/kebiao/kb_stu.php?xh=" + strconv.Itoa(i)
		data,err :=spiderST(URL)
		if err != nil {
			fmt.Println(err)
		}
		result,name :=RElessons(data)
		student :=PCK(result,i,name)
		sql(student)
	}


}
