package view

import (
	"fmt"

	"github.com/snguovo/web/v2.0/util"
)

const (
	//SYSTEM 页面
	SYSTEM = `<head>
	<style>
	#header {
	background-color:black;
	color:white;
	text-align:center;
	padding:5px;
	}
	thead tr{
		height: 40px;
		background-color: #ccc;
	}
	.mytable tr td{
		overflow:auto 
		} 
	</style>
	<script src="https://s3.pstatp.com/cdn/expire-1-M/jquery/3.0.0/jquery.min.js"></script>
	</head>
	<body>
	%s`

	systemStatus = `<div id="header">
	<h1>终端实时信息</h1>
	</div>
	<table width=100% border="1" cellpadding="0" class='mytable' cellspacing="0" style="table-layout:automatic">
	<thead>
	<tr>
	<th width=10%>IP地址</th>
	<th width=10%>MAC地址</th>
	<th>厂商信息</th>
	<th>系统</th>
	<th>硬件版本</th>
	<th>内核</th>
	<th>CPU核心数</th>
	<th>CPU使用率</th>
	<th>内存总量</th>
	<th>内存使用量</th>
	<th>内存使用率</th>
	<th>存储总量</th>
	<th>存储使用量</th>
	<th>存储使用率</th>
	<th>时间</th>
	</tr>
	</thead>
	<tr>
	<td id="ip"></td>
	<td id="mac"></td>
	<td id="platform"></td>
	<td id="os"></td>
	<td id="version"></td>
	<td id="arch"></td>
	<td id="cpu_count"></td>
	<td id="cpu_percent"></td>
	<td id="mem_total"></td>
	<td id="mem_used"></td>
	<td id="mem_used_percent"></td>
	<td id="disk_total"></td>
	<td id="disk_used"></td>
	<td id="disk_used_percent"></td>
	<td id="time"></td>
	</tr>
	</table>`

	systemUpdate = `<script>document.getElementById("ip").innerHTML = '%v';
	document.getElementById("mac").innerHTML = '%v';
	document.getElementById("platform").innerHTML = '%v';
	document.getElementById("os").innerHTML = '%v';
	document.getElementById("version").innerHTML = '%v';
	document.getElementById("arch").innerHTML = '%v';
	document.getElementById("cpu_count").innerHTML = '%v';
	document.getElementById("cpu_percent").innerHTML = '%v%%';
	document.getElementById("mem_total").innerHTML = '%v';
	document.getElementById("mem_used").innerHTML = '%v';
	document.getElementById("mem_used_percent").innerHTML = '%v%%';
	document.getElementById("disk_total").innerHTML = '%v';
	document.getElementById("disk_used").innerHTML = '%v';
	document.getElementById("disk_used_percent").innerHTML = '%v%%';
	document.getElementById("time").innerHTML = '%v';</script>`

	appStatus = `<div id="header">
	<h1>app信息</h1>
	</div>
	<table width=100%% border="2" cellpadding="0" class='mytable' cellspacing="0" style="table-layout:automatic">
	<thead>
	<tr>
	<th>app名称</th>
	<th>app版本</th>
	<th>日志路径</th>
	<th>配置路径</th>
	<th>enable</th>
	<th>CMD</th>
	<th>工作路径</th>
	<th>PID</th>
	<th>内存占用</th>
	<th>cpu占用</th>
	<th>运行与否</th>
	<th>操作</th>
	</tr>
	</thead>
	<tbody>
	</tbody>
	</table>
	<input type="checkbox" name="check" checked="true">是否刷新
	<script>
	function submitForm(node) {
        var tr = node.parentNode.parentNode; //获取当前元素的父节点的父节点，也就是tr。
        action = node.name;//取得元素的name属性的值
        app_name = tr.cells[0].innerHTML;//取得value值
        app_version = tr.cells[1].innerHTML;//取得value值
		//组装数据使用ajax传到后台，用了不少种类的数据传送，这样传的比较方便使用。
		var item = {
			action: action,
			app_name: app_name+":"+app_version
		  };
        $.ajax({
            type: "POST",
            url: "/system/appaction",
            data: item,
            error: function (request) {
                console.log("error");
                alert("fail")
            },
            success: function (data) {
                console.log("success");
                alert(data)
            }
        })
    } 
	$(document).ready(function () {
		var tbody = document.querySelector('tbody');
		var url;
		url = window.location.host; /* 获取主机地址 */
		url ='ws://'+url+'/system/upapp'
		// 指定websocket路径
		var websocket = new WebSocket(url);
		websocket.onopen = function (event) {
		  websocket.send('{"time":%v}');
		};
		websocket.onclose = function () {
		alert("连接已关闭...");
		}; 
		websocket.onmessage = function (event) {
			var checked = $("input[type='checkbox']").prop('checked');
			if (checked) {
				for (var i = tbody.childNodes.length - 1; i >= 0; i--) {
					tbody.removeChild(tbody.childNodes[i]);
				}
				var datas =eval('('+ event.data+')'); 
				for(var i = 0;i < datas.length; i++){
					//1.创建tr行
					var tr = document.createElement('tr');
					tbody.appendChild(tr);
					//2.行里面创建单元格（跟数据有关系的3个单元格） td 单元格的数量取决于每个对象里面的属性个数  for循环遍历对象
					for(var j in datas[i]){//里面的for循环管列td
						//创建单元格
						var td = document.createElement('td');
						//把对象里面的属性值 给td
						td.innerHTML = datas[i][j];
						tr.appendChild(td);
					}
					var td = document.createElement('td');
					//把对象里面的属性值 给td
					var e =document.createElement("input");
					e.type = "button";
					e.value = "启动";
					e.name = "start";
					e.setAttribute("class" ,"but");
					e.setAttribute("onclick", "submitForm(this)");
					td.appendChild(e);
					var b =document.createElement("input");
					b.type = "button";
					b.value = "停止";
					b.name = "stop";
					b.setAttribute("class" ,"but");
					b.setAttribute("onclick", "submitForm(this)");
					td.appendChild(b);
					tr.appendChild(td);
				}
			}
		};
	  });
</script>
</body>`

	dockerStatus = `<div id="header">
	<h1>docker信息</h1>
	</div>
	<table width=100%% border="3" cellpadding="0" class='mytable' cellspacing="0" style="table-layout:fixed">
	<thead>
	<tr>
	<th width=20%%>容器id</th>
	<th width=20%%>容器名称</th>
	<th width=20%%>镜像名称</th>
	<th>状态</th>
	<th>是否运行</th>
	<th>cpu占用</th>
	<th>内存占用</th>
	</tr>
	</thead>
	<tbody>
	</tbody>
	</table>
	<input type="checkbox" name="check" checked="true">是否刷新
	<script>
	$(document).ready(function () {
		var tbody = document.querySelector('tbody');
		var url;
		url = window.location.host; /* 获取主机地址 */
		url ='ws://'+url+'/system/updocker'
		// 指定websocket路径
		var websocket = new WebSocket(url);
		websocket.onopen = function (event) {
		  websocket.send('{"time":%v}');
		};
		websocket.onclose = function () {
		alert("连接已关闭...");
		}; 
		websocket.onmessage = function (event) {
			var checked = $("input[type='checkbox']").prop('checked');
			if (checked) {
				for (var i = tbody.childNodes.length - 1; i >= 0; i--) {
					tbody.removeChild(tbody.childNodes[i]);
				}
				var datas =eval('('+ event.data+')'); 
				for(var i = 0;i < datas.length; i++){
					//1.创建tr行
					var tr = document.createElement('tr');
					tbody.appendChild(tr);
					//2.行里面创建单元格（跟数据有关系的3个单元格） td 单元格的数量取决于每个对象里面的属性个数  for循环遍历对象
					for(var j in datas[i]){//里面的for循环管列td
						//创建单元格
						var td = document.createElement('td');
						//把对象里面的属性值 给td
						td.innerHTML = datas[i][j];
						tr.appendChild(td);
					}
				}
			}
		};
	  });
</script>
</body>`
)

//AddSystem 返回system页面
func AddSystem(html *string) {
	*html = fmt.Sprintf(*html, systemStatus)
}

//UpdateSystem 拼凑system信息
func UpdateSystem(systemInfo *util.SystemInfo) string {
	return fmt.Sprintf(systemUpdate, systemInfo.IP,
		systemInfo.Mac, systemInfo.Platform, systemInfo.Os,
		systemInfo.KernelVersion, systemInfo.KernelArch,
		systemInfo.CPUCount, systemInfo.CPUPercent,
		systemInfo.MemTotal, systemInfo.MemUsed, systemInfo.MemUsedPercent,
		systemInfo.DiskTotal, systemInfo.DiskUsed, systemInfo.DiskUsedPercent,
		systemInfo.Time.Format("2006-01-02 15:04:05 MST"))
}

//AddDocker 返回docker页面
func AddDocker(html *string, interval int) {
	*html = fmt.Sprintf(*html, fmt.Sprintf(dockerStatus, interval))
}

//AddApps 拼凑app信息
func AddApps(html *string, interval int) {
	*html = fmt.Sprintf(*html, fmt.Sprintf(appStatus, interval))
}
