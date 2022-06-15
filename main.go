//activtyId:
//"qyb_0934"

package main

import (
	"concurrency_test_tool/libs/utils"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/marknown/util"
)

func main() {
	// 打开本地浏览器
	errOpen := util.Open("http://127.0.0.1:18880/submit")
	if errOpen != nil && errOpen.Error() != "exec: already started" {
		fmt.Println(errOpen.Error())
	}

	http.HandleFunc("/submit", func(w http.ResponseWriter, r *http.Request) {
		const tpl = `<!DOCTYPE html>
	<html lang="en">
	<head>
		<meta charset="UTF-8">
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
		<meta http-equiv="X-UA-Compatible" content="ie=edge">
		<title>并发测试工具</title>
		<style>
			body{margin:0;}
			select option{direction: rtl;border:1px solid #dedede;}
			input, textarea{outline:none;border:1px solid #dedede;}
			h1{font-size:16px;display:inline;}
			.desc{font-size:10px;line-height:14px;}
			#main{width:1280px;margin:auto;}
			#main div {line-height:20px;font-size:10px;}
			#tips{font-size:10px;}
			#main .action{font-size:10px;margin-top:10px;}
			.input-start{width:120px;height:20px;font-size:12px;}
			.source{width:100%;font-size:10px;}
			.operate{text-align:center;}
			.submit{color:#fff;background-color:#67c23a;border-color:#67c23a;border-radius: 3px;padding:2px 10px;font-size:10px;}
			.submit-disabled{color:#fff;background-color:#dedede;border-color:#dedede;border-radius: 3px;padding:2px 10px;font-size:10px;}

		</style>
		<script type="text/javascript">
		var Ajax = {
			post: function(url, data, fn) { // data 应为'a=a1&b=b1'这种字符串格式
				var xhr = new XMLHttpRequest();
				xhr.open("POST", url, true);
				// 添加http头，发送信息至服务器时内容编码类型
				xhr.setRequestHeader("Content-Type", "application/x-www-form-urlencoded");
				xhr.onreadystatechange = function() {
					if(xhr.readyState == 4 && (xhr.status == 200 || xhr.status == 304)) {
						fn.call(this, xhr.responseText);
					}
				};
				xhr.send(data);
			}
		}
	
		function submit() {
			disableSubmitButton()
			let times    = document.getElementById('times').value
			let interval = document.getElementById('interval').value
			let start    = document.getElementById('start').value
			let before   = document.getElementById('before').value
			let source   = document.getElementById('source').value
	
			Ajax.post("/receive", "times="+times+"&interval="+interval+"&start="+encodeURIComponent(start)+"&before="+before+"&source="+encodeURIComponent(source), function (response) {
				try {
					var obj = JSON.parse(response)
					if (!obj.success) {
						showTips(obj.message)
					} else {
						if (obj.Responses) {
							let len = obj.Responses.length
							let resp = obj.Responses
							let message = "";
							for(i=0;i<len;i++) {
								message += "第" + resp[i].Times + "次请求耗时" + resp[i].Interval + "ms，请求时间：" + resp[i].StartTime + "，响应时间：" + resp[i].EndTime + " 响应如下：\n" + resp[i].Response + "\n\n";
							}
							showTips(message)
						} else {
							showTips(obj.message)
						}
					}
				} catch(e) {
					showTips(e + "\norigin：" + response)
				}
			});
		}
	
		function showTips(msg) {
			enableSubmitButton()
			let tipsObj = document.getElementById('tips')
			tipsObj.innerText = msg
		}

		function disableSubmitButton() {
			var submitBtn = document.getElementById("submitButton")
			submitBtn.disabled=true
			submitBtn.setAttribute("class", "submit submit-disabled")
			submitBtn.value = "提交中，请耐心等待"
		}

		function enableSubmitButton() {
			var submitBtn = document.getElementById("submitButton")
			submitBtn.disabled=false
			submitBtn.setAttribute("class", "submit")
			submitBtn.value = "提交"
		}
	</script>
	</head>
	<body>
		<div id="main">
			<div>
				<h1>并发测试工具v1.0.0</h1>
				<span class="desc">本工具通过并发提交从浏览器复制的curl请求，来检测接口的抗并发能力。可用于测试扣减库存/余额提现/活动资格等并发场景。测完可以去redis或者数据库看扣减是否正确</span>
			</div>
			<div class="action">
				并发次数<select id="times">
					<option value="1">1次</option>
					<option value="2">2次</option>
					<option value="3">3次</option>
					<option value="5">5次</option>
					<option value="8">8次</option>
					<option value="10">10次</option>
				</select>
				并发间隔<select id="interval">
					<option value="0">0ms</option>
					<option value="30">30ms</option>
					<option value="50">50ms</option>
					<option value="80">80ms</option>
					<option value="100">100ms</option>
					<option value="150">150ms</option>
					<option value="200">200ms</option>
					<option value="300">300ms</option>
					<option value="500">500ms</option>
					<option value="1000">1000ms</option>
				</select>
				执行时间<input id="start" class="input-start" value="{{.NowTime}}"/>
				提前执行<select id="before">
					<option value="0">0ms</option>
					<option value="5">5ms</option>
					<option value="10">10ms</option>
					<option value="20">20ms</option>
					<option value="30">30ms</option>
					<option value="50">50ms</option>
					<option value="80">80ms</option>
					<option value="100">100ms</option>
					<option value="150">150ms</option>
					<option value="200">200ms</option>
					<option value="300">300ms</option>
					<option value="500">500ms</option>
					<option value="1000">1000ms</option>
				</select>在指定时间的前多少毫秒执行。
			</div>
			<div>
				占位符1：<span>&#123&#123unixTimestamp&#125&#125 毫秒时间戳 例如：1655301568</span> 占位符2：<span>&#123&#123unixTimestampMillisecond&#125&#125 毫秒时间戳 例如：1655301568510</span>
			</div>
			<div>
				<textarea id="source" class="source" rows="30"></textarea>
			</div>
			<div class="operate">
				<input id="submitButton" type="button" value="提交" class="submit" onclick="submit()" />
			</div>
			<div id="tips"></div>
		</div>
	</body>
	</html>`

		tmpl, err := template.New("default").Parse(tpl)
		if err != nil {
			fmt.Printf("%s\n", err.Error())
		}
		type Params struct {
			NowTime string
		}
		tmpl.Execute(w, Params{NowTime: utils.FormatTimeStringCNSecond(time.Now())})
	})

	http.HandleFunc("/receive", func(w http.ResponseWriter, r *http.Request) {
		type concurrencyResp struct {
			Times     int
			Interval  int
			StartTime string
			EndTime   string
			Response  string
		}

		type response struct {
			Success   bool   `json:"success"`
			Message   string `json:"message"`
			Responses []*concurrencyResp
		}

		res := response{
			Success:   true,
			Message:   "ok",
			Responses: []*concurrencyResp{},
		}

		times, err := strconv.Atoi(r.FormValue("times"))
		if err != nil {
			res.Success = false
			res.Message = "并发次数不正确"
		}
		interval, err := strconv.Atoi(r.FormValue("interval"))
		if err != nil {
			res.Success = false
			res.Message = "并发间隔不正确"
		}
		before, err := strconv.Atoi(r.FormValue("before"))
		if err != nil {
			res.Success = false
			res.Message = "提前执行时间不正确"
		}
		start := r.FormValue("start")
		var startTime time.Time
		if start == "now" {
			startTime = time.Now()
		} else {
			startTime, err = utils.ParseTimeFromString(start)
			if err != nil {
				res.Success = false
				res.Message = "执行时间只能为now或者具体年月日时分秒"
			}
		}

		source := r.FormValue("source")

		// 添加占位符替换
		unixTimestampNano := time.Now().UnixNano()
		unixTimestamp := unixTimestampNano / 1000000000
		unixTimestampMillisecond := unixTimestampNano / 1000000
		source = strings.ReplaceAll(source, "{{unixTimestamp}}", fmt.Sprintf("%d", unixTimestamp))
		source = strings.ReplaceAll(source, "{{unixTimestampMillisecond}}", fmt.Sprintf("%d", unixTimestampMillisecond))

		if !strings.Contains(source, "curl") {
			res.Success = false
			res.Message = "请粘贴从浏览器复制出的完整CURL信息"
		}

		// fmt.Println(times, interval, start, startTime, source)
		// fmt.Printf("%s\n", source)
		if !res.Success {
			json, _ := json.Marshal(res)
			fmt.Fprintf(w, "%s", json)
			return
		}

		sreq, err := utils.BuildRequest(source)

		if err != nil {
			res.Success = false
			res.Message = fmt.Sprintf("buildRequest error:%s", err.Error())
		}
		if !res.Success {
			json, _ := json.Marshal(res)
			fmt.Fprintf(w, "%s", json)
			return
		}

		// 包装成一个匿名函数，方便执行
		var doTask = func() []*concurrencyResp {
			// 如果还没有到时间，先休息到指定时间
			if before < 0 {
				before = 0
			}
			beforeDuration := time.Duration(before) * time.Millisecond
			if time.Until(startTime) > beforeDuration {
				time.Sleep(time.Until(startTime) - beforeDuration)
			}

			var wg sync.WaitGroup
			var cacheResp = []*concurrencyResp{}
			for i := 1; i <= times; i++ {
				wg.Add(1)
				go func(i int) {
					time.Sleep(time.Duration(interval*(i-1)) * time.Millisecond)
					startTime := time.Now()
					req, err := sreq.New().Request()
					var resp []byte
					if err == nil {
						resp, err = utils.DoRequest(req)
						if err != nil {
							resp = []byte(err.Error())
						}
					} else {
						resp = []byte(err.Error())
					}

					endTime := time.Now()

					cacheResp = append(cacheResp, &concurrencyResp{
						Times:     i,
						Interval:  int(endTime.Sub(startTime).Milliseconds()),
						StartTime: utils.FormatTimeStringCNMillisecond(startTime),
						EndTime:   utils.FormatTimeStringCNMillisecond(endTime),
						Response:  string(resp),
					})
					wg.Done()
				}(i)
			}
			wg.Wait()

			return cacheResp
		}

		// 如果距离时间太长，后台执行
		if time.Until(startTime) > 60*time.Second {
			res.Success = true
			res.Message = "请求后台进行中，执行完成后日志保存在当前目录 concurrency_log.txt 文件里。"
			res.Responses = nil
			go func() {
				cacheResp := doTask()

				for _, v := range cacheResp {
					utils.LogInfo(fmt.Sprintf("第%d次请求耗时%dms，请求时间：%s，响应时间：%s 响应如下：\n%s\n\n", v.Times, v.Interval, v.StartTime, v.EndTime, v.Response))
				}
			}()
		} else {
			cacheResp := doTask()

			res.Success = true
			res.Message = "ok"
			res.Responses = cacheResp
		}

		json, _ := json.Marshal(res)
		fmt.Fprintf(w, "%s", json)
	})

	err := http.ListenAndServe(":18880", nil)

	if err != nil {
		log.Fatal(err.Error())
	}
}
