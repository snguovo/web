package view

import "fmt"

const (
	//INDEX 页面
	INDEX = `<head>
	<style>
	#header {
	text-align:center;
	padding:5px;
	}
	</style>
	<script type="text/javascript">
	function UpladFile() {
		var fileObj = document.getElementById("uploadfile").files[0]; // js 获取文件对象
		var path = document.getElementById("path").value; // js 获取文件对象
		// FormData 对象
		var form = new FormData();
		form.append("path", path);                        // 可以增加表单数据
		form.append("file", fileObj);                           // 文件对象
		// XMLHttpRequest 对象
		var xhr = new XMLHttpRequest();
		xhr.open('POST', "uploadfile");
		xhr.onload = function () {
			if (xhr.status === 200) {
				alert('上传成功');
			} else {
				alert('请选择文件和保存路径');
			}
		};
		xhr.upload.addEventListener("progress", progressFunction, false);
		xhr.send(form);
		
	}
	function progressFunction(evt) {
		var progressBar = document.getElementById("progressBar");
		var percentageDiv = document.getElementById("percentage");
		if (evt.lengthComputable) {
			progressBar.max = evt.total;
			progressBar.value = evt.loaded;
			percentageDiv.innerHTML = Math.round(evt.loaded / evt.total * 100) + "%%";
		}
	}
</script>
	</head>
	<body>
	<div id="header">
	<h1>终端管理</h1>
	</div>
	<form name="form1" action="system" method="get">
	<fieldset>
	<legend>终端信息:</legend>
	刷新间隔： <input type="number" name="time" value="1" min="1" max="10" required="on" />秒 <br>
	选择要查看的内容：<select name="content">
	<option value="system">系统信息</option>
	<option value="app">app信息</option>
	<option value="docker">docker信息</option>
	</select>
	<input type="submit" value="查看">
	</fieldset>
	</form>
	<form name="form2" action="log" method="get">
	<fieldset>
	<legend>日志打印:</legend>
	选择app： <select name="app" >
	%s
	</select>
	<input type ="submit" formtarget="_blank" value="打印日志" > 
	</fieldset>
	</form>
	<fieldset>
	<legend>文件上传:</legend>
    <input type="file" id="uploadfile" name="myfile" required="on"/> 
	保存路径<input type="text" id="path" name="path" value="./" required="on"/> <br />
	<span id="percentage"></span>
    <br />
    <progress id="progressBar" value="0" max="100">
    </progress>
    <input type="button" onclick="UpladFile()" value="上传">
	</fieldset>
	<fieldset>
	<legend>配置文件修改:</legend>
	<form name="form3" action="/" method="get">
	输入配置文件路径： <input type="text" name="filepath" required="on" value="./appconf.json">
	<input type ="submit" value="读取配置">
	</form>
	<form name="form4" action="updateconf" method="post">
	配置文件 <input id="filepath" type="text" name="filepath" value="" readonly="on"> <br />
	<textarea id="file" name="config" rows="30" cols="150" autocomplete="on"></textarea>
	<input type ="submit" value="更新配置" > 
	</form>
	%s
	</fieldset>
	</body>`

	option = `<option value="%s">%s</option>`

	conffile = "<script>document.getElementById(\"file\").innerHTML=`%s`;document.getElementById(\"filepath\").value=`%s`;</script>"
)

//AddOptions 加载配置文件App
func AddOptions(options *string, appName string) {
	*options = *options + fmt.Sprintf(option, appName, appName)
}

//AddConfFile 加载配置文件内容
func AddConfFile(confFile *string, content, filepath string) {
	*confFile = *confFile + fmt.Sprintf(conffile, content, filepath)
}
