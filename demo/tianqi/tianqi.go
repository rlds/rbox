package main 

import(
	"net/http"
	"io/ioutil"
)

//
func init(){
	//城市代码初始化
	cityData()
}

func httpGet(url string)( error,string){
	resp, err := http.Get(url)
	if err != nil {
		return err,""
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err,""
 	}
	return nil,string(body)
}

// 天气查询
func checkTianqi(cityName string)(res string){
	wid := cityMap(cityName)
	if len(wid) < 1 {
		res = `{"error":"输入城市名称错误"}`
		return
	}
	// /sk_2d/101210101.html?_=1542282548892
	// 	url := `http://d1.weather.com.cn/sk_2d/`+wid+`.html?_`+ util.IntToStr(int(time.Now().UnixNano()/1000000))
	url := `http://www.weather.com.cn/sk_2d/`+wid+`.html`
	_ , res = httpGet(url)
	return
}
