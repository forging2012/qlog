$(document).ready(function(){
	 //添加日志分析配置
     $("#settings-form").submit(function(event){
        event.preventDefault();
        var error=$("#settings-form-error");
        var info=$("#settings-form-info");
        var settingsGroup=$("#settings-group");
        var bucket=$(this).find("[name='bucket']").val();
        var saveBucket=$(this).find("[name='saveBucket']").val();
        var saveBucketDomain=$(this).find("[name='saveBucketDomain']").val();
        var isSaveBucketPrivate=$(this).find("[name='isSaveBucketPrivate']:checked").val();

        if ($.trim(bucket)==""){
            error.text("请填写日志源空间!");
            error.show();
        }else if ($.trim(saveBucket)==""){
            error.text("请填写日志保存空间");
            error.show();
        }else if($.trim(saveBucketDomain)==""){
            error.text("请填写日志保存空间的域名");
            error.show();
        }else if($.trim(isSaveBucketPrivate)==""){
            error.text("请选择日志保存空间的类型");
            error.show();
        }else{
            error.hide();
            $.post("/settings/add",{
                bucket:bucket,
                saveBucket:saveBucket,
                saveBucketDomain:saveBucketDomain,
                isSaveBucketPrivate:isSaveBucketPrivate,
            },function(respData){
                if (respData.error!=undefined){
                    error.text(respData.error);
                    error.show();
                }else{
                    info.text(respData.data);
                    info.fadeIn(1000);
                    info.fadeOut(500);
                    //append one row
                    type="公开空间";
                    typeImg="";
                    if(isSaveBucketPrivate==1){
                        type="私有空间";
                        typeImg='<span class="glyphicon glyphicon-lock" aria-hidden="true"></span>';
                    }
                    settingsGroup.find("button").each(function(){
                        var parent=$(this).parent().parent();
                        var sid=$(this).attr("sid");
                        if (sid==bucket){
                            parent.remove();
                        }
                    });
                    settingsGroup.find("tbody").append("<tr><td class='col1'>"+bucket+"</td><td class='col2'>"+
                        saveBucket+"</td><td class='col3'>"+saveBucketDomain+
                        "</td><td class='col4'>"+type+"&nbsp;"+typeImg+"</td><td>"+
                        '<button class="edit-settings-btn btn btn-xs btn-primary" sid="'+bucket+'">修改</button>&nbsp;'+
                        '<button class="delete-settings-btn btn btn-xs btn-danger" sid="'+bucket+'">删除</button>'+
                        "</td></tr>");
                }
            });
        }
     });

    //修改日志分析配置
    $("#settings-group").on("click",".edit-settings-btn",function(){
        var parent=$(this).parent().parent();
        var bucket=$.trim(parent.find("td.col1").text());
        var saveBucket=$.trim(parent.find("td.col2").text());
        var saveBucketDomain=$.trim(parent.find("td.col3").text());
        var isPrivate=$.trim(parent.find("td.col4").text());
        var isSaveBucketPrivate=1;
        if(isPrivate=="公开空间"){
            isSaveBucketPrivate=0;
        }
        var settingsForm=$("#settings-form");
        settingsForm.find("[name='bucket']").val(bucket);
        settingsForm.find("[name='saveBucket']").val(saveBucket);
        settingsForm.find("[name='saveBucketDomain']").val(saveBucketDomain);
        settingsForm.find("[name='isSaveBucketPrivate']").each(function(){
            if ($(this).val()==isSaveBucketPrivate){
                $(this).prop("checked",true);
                return;
            }
        });
    });

    //删除日志分析配置
    $("#settings-group").on("click",".delete-settings-btn",function(){
        var objToDel=$(this).parent().parent();
        var bucket=$(this).attr("sid");
        if(confirm("确认删除么？")){
            $.post("/settings/delete",{
                bucket:bucket,
            },function(respData){
                if (respData.error!=undefined){
                    alert(respData.error);
                }else{
                    objToDel.remove();
                }
            });
        }
    });

    ///////////////////////////////////////

    $("#prepare-form").submit(function(event){
        event.preventDefault();
        var error=$("#prepare-form-error");
        var info=$("#prepare-form-info");
        var bucket=$(this).find("[name='bucket']").val();
        var date=$(this).find("[name='date']").val();
        if($.trim(bucket)==""){
            error.text("请选择一个要进行预处理的空间!");
            error.show();
        }else if($.trim(date)==""){
            error.text("请选择一个要进行预处理的日志日期!");
            error.show();
        }else{
            error.hide();
            $.post("/prepare",{
                bucket:bucket,
                date:date,
            },function(respData){
                if (respData.error!=undefined){
                    error.text(respData.error);
                    error.show();
                }else{
                    info.text(respData.data.msg);
                    info.fadeIn(1000);
                    info.fadeOut(500);
                }
            });
        }
    });
});