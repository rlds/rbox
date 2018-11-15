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
	url := `http://www.weather.com.cn/data/cityinfo/`+wid+`.html`
	_ , res = httpGet(url)
	return
}
