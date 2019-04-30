# [js清除当前站点cookie](https://blog.csdn.net/shua67/article/details/81535883)

    function clearCookie() {            
        var keys = document.cookie.match(/[^ =;]+(?=\=)/g);
        if (keys) {
            for (var i = keys.length; i--;) {
                document.cookie = keys[i] + '=0;path=/;expires=' + new Date(0).toUTCString();//清除当前域名下的,例如：m.kevis.com
                document.cookie = keys[i] + '=0;path=/;domain=' + document.domain + ';expires=' + new Date(0).toUTCString();//清除当前域名下的，例如 .m.kevis.com
                document.cookie = keys[i] + '=0;path=/;domain=kevis.com;expires=' + new Date(0).toUTCString();//清除一级域名下的或指定的，例如 .kevis.com
            }
        }
        $("#divcookie").html(document.cookie);
        alert('已清除');
    }

**注意：**

    1、设置 cookie 时明确指定 domain 域名，子域名可读取（子域共享该cookie），删除时则也必须明确指定域名，否则无法删除。(这种产生的是 .baidu.com)
    
    2、设置 cookie 时不指定域名，使用默认值，则表示只有当前域名可见（子域不可共享）。删除时也不需要指定域名，否则无法删除。(这种产生的是 www.baidu.com)

---------------------------------------