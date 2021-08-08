package main

import (
	"fmt"
	"html/template"
	"net/http"
	"xiaomi/xiaomi"
)

func main()  {
	http.Handle("/file/",  http.StripPrefix("/file/", http.FileServer(http.Dir("./template/file"))))
	http.HandleFunc("/", index)
	http.HandleFunc("/run", run)

	err := http.ListenAndServe(":9090", nil)
	if err != nil {
		fmt.Println("HTTP server failed,err:", err)
		return
	}
}

func run(w http.ResponseWriter, r *http.Request)  {
	r.ParseForm()

	username := r.Form.Get("user")
	password := r.Form.Get("password")
	step := r.Form.Get("step")

	if username == "" {
		fmt.Fprintf(w, `{"code":0,"message":"帐号不能为空"}`)
		return
	}

	if password == "" {
		fmt.Fprintf(w, `{"code":0,"message":"密码不能为空"}`)
		return
	}

	if step == "" {
		fmt.Fprintf(w, `{"code":0,"message":"步数不能为空"}`)
		return
	}

	token,err := xiaomi.Login(username,password)
	if err != nil{
		fmt.Fprintf(w, `{"code":0,"message":"登录失败"}`)
		return
	}
	//int, _ := strconv.Atoi(step)
	res,_ := xiaomi.Run(token,step)
	if res{
		fmt.Fprintf(w, `{"code":0,"message":"步数已修改为`+step+`步"}`)
		return
	}
	fmt.Fprintf(w, `{"code":0,"message":"刷取失败"}`)

}

func index(w http.ResponseWriter, r *http.Request){
	tmpl, err := template.ParseFiles("./template/index.html")
	if err != nil {
		fmt.Println("create template failed, err:", err)
		return
	}
	query := r.URL.Query()
	data := map[string]string{
		"id":  query.Get("id"),
		"url":query.Get("url"),
	}
	tmpl.Execute(w, data)
}