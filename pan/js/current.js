
function check_del(form)
{var temp;
temp=true;

if(window.confirm("确定删除选中的数据？"))
  {
  temp=true;
  }else{
  return false;
  }
return temp;
}

function check_boot(form)
{var temp;
temp=true;

if(window.confirm("确定创建？"))
  {
  temp=true;
  }else{
  return false;
  }
return temp;
}


function check_upload(form)
{var temp;
temp=true;

if(window.confirm("确定上传？"))
  {
  temp=true;
  }else{
  return false;
  }
return temp;
}

