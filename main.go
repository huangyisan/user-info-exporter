package main

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "net/http"
)

// up
// user_count
// department_count
// mail_type_count
// mail_type_qq_count

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


const website = "https://www.fastmock.site"
const url = "/mock/410a590380265355dbf0de54a8af2454/index/api/user"
func getAllUserInfo() []byte {
    res, err := http.Get(website + url)
    if err != nil {
        panic(err)
    }
    defer res.Body.Close()
    content, err := ioutil.ReadAll(res.Body)
    fmt.Println(string(content))
    return content
}

func upMetric() {

}

func main() {
    res := getAllUserInfo()

    var allUserInfo Response
    if err := json.Unmarshal(res, &allUserInfo); err != nil {
        panic(err)
    }
    fmt.Println(len(allUserInfo))
}
