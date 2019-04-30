export function CollectInfos(page,time,event,user){
    var _maq =[];
    _maq.push(['page',page]);
    _maq.push(['time', time]);
    _maq.push(['event',event]);
    _maq.push(['user',user]);
    window._maq = _maq;
    //匿名函数-重点(自调用匿名函数只会在运行时执行一次，一般用于初始化)
    (function(window) {
        //引入外部js文件,创建一个script,并根据协议(http或https)将src指向对应的ma.js,最后将这个元素插入页面的dom树上
        var ma = document.createElement('script');
        ma.type = 'text/javascript';
        ma.async = true; //异步调用外部 js 文件，即不阻塞浏览器的解 析，待外部 js 下载完成后异步执行。这个属性是 HTML5 新引入的。
        ma.src = 'http://127.0.0.1:9999/maidian.js';
        var s = document.getElementsByTagName('script')[0];
        s.parentNode.insertBefore(ma, s);
    })(window);
}
