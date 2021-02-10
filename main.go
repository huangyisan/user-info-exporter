package main

import (
    "encoding/json"
    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promhttp"
    "io/ioutil"
    "net/http"
    "strings"
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
mail_type_count
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
        prometheus.BuildFQName(namespace, "ddd", "mail_type"),
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

type Response map[string][]UserInfo



type Exporter struct {

}

func MyExporter() *Exporter {
    e := &Exporter{}
    return e
}

func (e *Exporter) Describe(ch chan <-*prometheus.Desc) {
    // 将metric描述传入管道
    ch <- up
    ch <- user_count
    ch <- department_count
    ch <- mail_type_count
}

func (e *Exporter) Collect(ch chan <-prometheus.Metric) {
    // 定义数据的获取并传入到管道中
    res, err := getAllUserInfo()
    // up Metric
    err = upMetric(err, ch)
    if err != nil {
        return
    }

    userInfo, err := praseResult(res)
    if err != nil {
        return
    }

    // user_count Metric
    userCountMetric(userInfo, ch)

    // department_count Metric
    departmentCountMetric(userInfo, ch)

    // mailTypeCount Metric
    mailTypeCountMetric(userInfo, ch)



}

func getAllUserInfo() ([]byte, error) {
    res, err := http.Get(website + url)
    if err != nil {
        return nil, err
    }
    defer res.Body.Close()
    content, err := ioutil.ReadAll(res.Body)
    return content, nil
}

func praseResult(res []byte) (Response, error) {
    var allUserInfo Response
    if err := json.Unmarshal(res, &allUserInfo); err != nil {
        return nil, err
    }
    return allUserInfo, nil
}


// get Metric value
func upMetric(e error, ch chan <-prometheus.Metric) (err error) {
    if e != nil {
        // 如果获取数据失败报错,则健康检查为0, 失败
        ch <- prometheus.MustNewConstMetric(up, prometheus.GaugeValue, 0)
        return e
    }
    ch <- prometheus.MustNewConstMetric(up, prometheus.GaugeValue, 1,"userinfo","health")
    return nil
}

func userCountMetric(res Response, ch chan <- prometheus.Metric) {
    userCount := 0
    for k,_ := range res {
        userCount += len(res[k])
    }
    uc := float64(userCount)
    ch <- prometheus.MustNewConstMetric(user_count, prometheus.GaugeValue, uc,"user_count")
}

func departmentCountMetric(res Response, ch chan <- prometheus.Metric) {
    dc := float64(len(res))
    ch <- prometheus.MustNewConstMetric(department_count, prometheus.GaugeValue, dc,"department_count")
}


func mailTypeCountMetric(res Response, ch chan <- prometheus.Metric)  {
    mailTypeMap := make(map[string]float64)
    for k, _ := range res {
        user := res[k]
        for _, v := range user {
            mailType := strings.Split(v.Email, "@")
            mt := mailType[len(mailType)-1]
            mailTypeCount := mailTypeMap[mt]
            mailTypeMap[mt] = mailTypeCount + 1
        }
    }
    for k,v := range mailTypeMap {
        ch <- prometheus.MustNewConstMetric(mail_type_count, prometheus.GaugeValue, v, k)
    }
}

func main() {
    exporter := MyExporter()
    prometheus.MustRegister(exporter)

    http.Handle("/userinfo/metrics", promhttp.Handler())
    http.ListenAndServe("0.0.0.0:8889", nil)

}
