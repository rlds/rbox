/*
*  csbox.js
*  csbox
*
*  Created by 吴道睿 on 2018/4/26.
*  Copyright © 2018年 吴道睿. All rights reserved.
*/

$(document).ready(function(){
	if (sel_groupname != "<no value>" && sel_boxname != "<no value>"){
		$("#"+sel_groupname+"-"+sel_boxname).addClass( "active");
	}
});

$("#BoxSubmit").click(function(){
    if (sel_groupname.length > 0 && sel_boxname.length > 0) {
		var callpath = "/call/"+sel_groupname+"/"+sel_boxname ;
		var inputs = $("[id='boxInput']");;
		var params = {};
		for (var i=0;i<inputs.length;i++){
			params[inputs[i].name]  =	 inputs[i].value; 
		}
		var IsSync = false;
		DoingText = Doing;
		$.ajax({
			   url : callpath,
			   type: "POST",
			   data : JSON.stringify(params),
			   dataType : "json",
			   contentType:"application/json; charset=utf-8",
			   success: function(data, textStatus, jqXHR)
			   {
			       IsSync = data.IsSync ;
			       if (!IsSync) {
			           callpath = "/status/"+sel_groupname+"/"+sel_boxname + "?TaskId="+data.TaskId;
			           callstatus(callpath,params);
			       }else{
			           setResData(data);
			        }
				},
				error: function (jqXHR, textStatus, errorThrown){
			        $("#boxResData").html("请求错误")
				}
		});
	}
});

var Doing = "  执行中,请稍候.";
var DoingText = "";

function callstatus(callpath,params){
	$.ajax({
		   url : callpath,
		   type: "POST",
		   data : JSON.stringify(params),
		   dataType : "json",
		   contentType:"application/json; charset=utf-8",
		   success: function(data, textStatus, jqXHR)
		   {
		       setResData(data);
		       if (data.Status != "COMPLETE") {
		           setTimeout(function() {
					  callstatus(callpath,params)
				   },2000);
			   };
		   },
		   error: function (jqXHR, textStatus, errorThrown){
		        $("#boxResData").html("请求错误")
		   }
	});
}

function setResData(data){
	if (data.Status == "COMPLETE"){
		switch (data.Type) {
			case "markdown":{
				// 过长的数据不能转markdown
				var html_content = markdown.toHTML(data.Data);
				$("#boxResData").html(html_content);
				break;
			}
			default:{
				$("#boxResData").html(data.Data);
				break;
			}
		}
	}else{
		DoingText += ".";
		$("#boxResData").html(DoingText);
	}
}

