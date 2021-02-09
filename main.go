package main

import (
    "encoding/json"
    "fmt"
    "github.com/prometheus/client_golang/prometheus"
    "io/ioutil"
    "net/http"
)

const namespace = "UserInfo"
const website = "https://www.fastmock.site"
const url = "/mock/410a590380265355dbf0de54a8af2454/index/api/user"
/*
Metric Define
===
up
user_count
department_count
mail_type_qq_count
 */

var (
    up = prometheus.NewDesc(
        prometheus.BuildFQName(namespace, "", "up"),
        "UserInfo Api health check",
        []string{"userinfo","health"},
        map[string]string{"status":"up"})

    user_count = prometheus.NewDesc(
        prometheus.BuildFQName(namespace, "", "user_count"),
        "All the users' count",
        []string{"user_count"},
        nil)

    department_count = prometheus.NewDesc(
        prometheus.BuildFQName(namespace, "", "department_count"),
        "All the departments' count",
        []string{"department_count"},
        nil)

    mail_type_count = prometheus.NewDesc(
        prometheus.BuildFQName(namespace, "", "mail_type_count"),
        "All the mail types' count",
        []string{"mail_type_count"},
        nil)
)


/*
{
	"IT": [{
			"name": "Shyam",
			"email": "shyamjaiswal@gmail.com"
		},
		{
			"name": "Bob",
			"email": "bob32@gmail.com"
		}
}
 */
type UserInfo struct {
    Name string `json:"name"`
    Email string `json:"email"`
}

//type DepartMent struct {
//    DepartMentInfo []UserInfo
//}

type Response map[string][]UserInfo



type Exporter struct {

}

func MyExporter() *Exporter {
    e := &Exporter{}
    return e
}

func (e *Exporter) Describe(ch chan <-*prometheus.Desc) {
    ch <- up
    ch <- user_count
    ch <- department_count
    ch <- mail_type_count
}

func (e *Exporter) Collection(ch chan <-*prometheus.Collector) {
    
}

func getAllUserInfo() ([]byte, error) {
    res, err := http.Get(website + url)
    if err != nil {
        return nil, err
    }
    defer res.Body.Close()
    content, err := ioutil.ReadAll(res.Body)
    fmt.Println(string(content))
    return content, nil
}

func upMetric(num int) {

}

func main() {
    res, err := getAllUserInfo()
    if err != nil {

    }

    var allUserInfo Response
    if err := json.Unmarshal(res, &allUserInfo); err != nil {
        panic(err)
    }
    fmt.Println(len(allUserInfo))
}
