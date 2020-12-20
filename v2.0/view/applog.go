package view

import "fmt"

const (
	//APPLOG 页面
	APPLOG = `<head>
	<meta charset="utf-8">
	<script src="https://s3.pstatp.com/cdn/expire-1-M/jquery/3.0.0/jquery.min.js"></script>
  </head>
  
  <body>
	<div id="log-container" style="height: 600px; overflow-y: scroll; background: #333; padding: 10px;">
	  <div>
	  </div>
	</div>
	<input type="checkbox" name="check" checked="true">是否自动滚动 
  </body>
  <script>
	$(document).ready(function () {
	  var url;
	  url = window.location.host; /* 获取主机地址 */
	  url ='ws://'+url+'/log/upper'
	  // 指定websocket路径
	  var websocket = new WebSocket(url);
	  websocket.onopen = function (event) {
		websocket.send("%s");
	  };
	  websocket.onclose = function () {
		alert("连接已关闭...");
		}; 
	  websocket.onmessage = function (event) {
		// 接收服务端的实时日志并添加到HTML页面中（error显示红色）
		if (event.data.search("ERROR") != -1) {
			$("#log-container div").append(event.data).css("color", "#AA0000");
		  } else if (event.data.search("WARNING") != -1) {
			$("#log-container div").append(event.data).css("color", "#ffd700");
		  }else {
			$("#log-container div").append(event.data).css("color", "#aaa");
		  }
		// 滚动条滚动到最低部
		// 是否滚动
		var checked = $("input[type='checkbox']").prop('checked');
		if (checked) {
		  $("#log-container").scrollTop($("#log-container div").height() - $("#log-container").height());
		}
	  };
	});
  </script>
  </body>`
)

//AppLog 返回applog页面
func AppLog(appName string) string {
	return fmt.Sprintf(APPLOG, appName)
}
