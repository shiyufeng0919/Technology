//数据收集脚本maidian.js
(function () {
    var params = {};
    //Document对象数据
    if(document) {
        params.domain = document.domain || '';
        params.url = document.URL || '';
        params.title = document.title || ''; //页面 title
        params.referrer = document.referrer || ''; //上一跳 url
    }
    //Window对象数据
    if(window && window.screen) { //用户显示器分辨率
        params.sh = window.screen.height || 0;
        params.sw = window.screen.width || 0;
        params.cd = window.screen.colorDepth || 0;
    }
    //navigator对象数据
    if(navigator) {
        params.lang = navigator.language || '';
    }
    //解析_maq配置,收集配置信息。这里面可能会包括用户自定义的事件跟踪、业务数据(如电子商务网站的商品编号等)等
    if(_maq) {
        for(var i in _maq) {
            var execAction=_maq[i][0];//设置key值
            params[execAction]=_maq[i][1];//设置value值
        }
    }

    //拼接参数串
    var args = '';
    for(var i in params) {
        if(args != '') {
            args += '&';
        }
        args += i + '=' + encodeURIComponent(params[i]);
    }

    //通过Image对象请求后端脚本(不用ajax避免跨域)
    var img = new Image(1, 1);
    //log.gif 是后端脚本，是一个伪装成 gif 图片的脚本
    img.src = 'http://0.0.0.0:9999/log.gif?' + args; //前端是一个静态资源，Nginx后台实则为一个后台脚本
})();
