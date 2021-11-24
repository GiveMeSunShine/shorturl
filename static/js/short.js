
! function($) {

    "use strict";

    var Typed = function(el, options) {

        // chosen element to manipulate text
        this.el = $(el);

        // options
        this.options = $.extend({}, $.fn.typed.defaults, options);

        // attribute to type into
        this.isInput = this.el.is('input');
        this.attr = this.options.attr;

        // show cursor
        this.showCursor = this.isInput ? false : this.options.showCursor;

        // text content of element
        this.elContent = this.attr ? this.el.attr(this.attr) : this.el.text()

        // html or plain text
        this.contentType = this.options.contentType;

        // typing speed
        this.typeSpeed = this.options.typeSpeed;

        // add a delay before typing starts
        this.startDelay = this.options.startDelay;

        // backspacing speed
        this.backSpeed = this.options.backSpeed;

        // amount of time to wait before backspacing
        this.backDelay = this.options.backDelay;

        // input strings of text
        this.strings = this.options.strings;

        // character number position of current string
        this.strPos = 0;

        // current array position
        this.arrayPos = 0;

        // number to stop backspacing on.
        // default 0, can change depending on how many chars
        // you want to remove at the time
        this.stopNum = 0;

        // Looping logic
        this.loop = this.options.loop;
        this.loopCount = this.options.loopCount;
        this.curLoop = 0;

        // for stopping
        this.stop = false;

        // custom cursor
        this.cursorChar = this.options.cursorChar;

        // All systems go!
        this.build();
    };

    Typed.prototype = {

        constructor: Typed

        ,
        init: function() {
            // begin the loop w/ first current string (global self.string)
            // current string will be passed as an argument each time after this
            var self = this;
            self.timeout = setTimeout(function() {
                // Start typing
                self.typewrite(self.strings[self.arrayPos], self.strPos);
            }, self.startDelay);
        }

        ,
        build: function() {
            // Insert cursor
            if (this.showCursor === true) {
                this.cursor = $("<span class=\"typed-cursor\">" + this.cursorChar + "</span>");
                this.el.after(this.cursor);
            }
            this.init();
        }

        // pass current string state to each function, types 1 char per call
        ,
        typewrite: function(curString, curStrPos) {
            // exit when stopped
            if (this.stop === true) {
                return;
            }

            // varying values for setTimeout during typing
            // can't be global since number changes each time loop is executed
            var humanize = Math.round(Math.random() * (100 - 30)) + this.typeSpeed;
            var self = this;

            // ------------- optional ------------- //
            // backpaces a certain string faster
            // ------------------------------------ //

            // contain typing function in a timeout humanize'd delay
            self.timeout = setTimeout(function() {
                // check for an escape character before a pause value
                // format: \^\d+ .. eg: ^1000 .. should be able to print the ^ too using ^^
                // single ^ are removed from string
                var charPause = 0;
                var substr = curString.substr(curStrPos);
                if (substr.charAt(0) === '^') {
                    var skip = 1; // skip atleast 1
                    if (/^\^\d+/.test(substr)) {
                        substr = /\d+/.exec(substr)[0];
                        skip += substr.length;
                        charPause = parseInt(substr);
                    }

                    // strip out the escape character and pause value so they're not printed
                    curString = curString.substring(0, curStrPos) + curString.substring(curStrPos + skip);
                }

                if (self.contentType === 'html') {
                    // skip over html tags while typing
                    var curChar = curString.substr(curStrPos).charAt(0)
                    if (curChar === '<' || curChar === '&') {
                        var tag = '';
                        var endTag = '';
                        if (curChar === '<') {
                            endTag = '>'
                        } else {
                            endTag = ';'
                        }
                        while (curString.substr(curStrPos).charAt(0) !== endTag) {
                            tag += curString.substr(curStrPos).charAt(0);
                            curStrPos++;
                        }
                        curStrPos++;
                        tag += endTag;
                    }
                }

                // timeout for any pause after a character
                self.timeout = setTimeout(function() {
                    if (curStrPos === curString.length) {
                        // fires callback function
                        self.options.onStringTyped(self.arrayPos);

                        // is this the final string
                        if (self.arrayPos === self.strings.length - 1) {
                            // animation that occurs on the last typed string
                            self.options.callback();

                            self.curLoop++;

                            // quit if we wont loop back
                            if (self.loop === false || self.curLoop === self.loopCount)
                                return;
                        }

                        self.timeout = setTimeout(function() {
                            self.backspace(curString, curStrPos);
                        }, self.backDelay);
                    } else {

                        /* call before functions if applicable */
                        if (curStrPos === 0)
                            self.options.preStringTyped(self.arrayPos);

                        // start typing each new char into existing string
                        // curString: arg, self.el.html: original text inside element
                        var nextString = curString.substr(0, curStrPos + 1);
                        if (self.attr) {
                            self.el.attr(self.attr, nextString);
                        } else {
                            if (self.isInput) {
                                self.el.val(nextString);
                            } else if (self.contentType === 'html') {
                                self.el.html(nextString);
                            } else {
                                self.el.text(nextString);
                            }
                        }

                        // add characters one by one
                        curStrPos++;
                        // loop the function
                        self.typewrite(curString, curStrPos);
                    }
                    // end of character pause
                }, charPause);

                // humanized value for typing
            }, humanize);

        }

        ,
        backspace: function(curString, curStrPos) {
            // exit when stopped
            if (this.stop === true) {
                return;
            }

            // varying values for setTimeout during typing
            // can't be global since number changes each time loop is executed
            var humanize = Math.round(Math.random() * (100 - 30)) + this.backSpeed;
            var self = this;

            self.timeout = setTimeout(function() {

                // ----- this part is optional ----- //
                // check string array position
                // on the first string, only delete one word
                // the stopNum actually represents the amount of chars to
                // keep in the current string. In my case it's 14.
                if (self.contentType === 'html') {
                    // skip over html tags while backspacing
                    if (curString.substr(curStrPos).charAt(0) === '>') {
                        var tag = '';
                        while (curString.substr(curStrPos).charAt(0) !== '<') {
                            tag -= curString.substr(curStrPos).charAt(0);
                            curStrPos--;
                        }
                        curStrPos--;
                        tag += '<';
                    }
                }

                // ----- continue important stuff ----- //
                // replace text with base text + typed characters
                var nextString = curString.substr(0, curStrPos);
                if (self.attr) {
                    self.el.attr(self.attr, nextString);
                } else {
                    if (self.isInput) {
                        self.el.val(nextString);
                    } else if (self.contentType === 'html') {
                        self.el.html(nextString);
                    } else {
                        self.el.text(nextString);
                    }
                }

                // if the number (id of character in current string) is
                // less than the stop number, keep going
                if (curStrPos > self.stopNum) {
                    // subtract characters one by one
                    curStrPos--;
                    // loop the function
                    self.backspace(curString, curStrPos);
                }
                    // if the stop number has been reached, increase
                // array position to next string
                else if (curStrPos <= self.stopNum) {
                    self.arrayPos++;

                    if (self.arrayPos === self.strings.length) {
                        self.arrayPos = 0;
                        self.init();
                    } else
                        self.typewrite(self.strings[self.arrayPos], curStrPos);
                }

                // humanized value for typing
            }, humanize);

        }
        ,
        reset: function() {
            var self = this;
            clearInterval(self.timeout);
            var id = this.el.attr('id');
            this.el.after('<span id="' + id + '"/>')
            this.el.remove();
            if (typeof this.cursor !== 'undefined') {
                this.cursor.remove();
            }
            // Send the callback
            self.options.resetCallback();
        }

    };

    $.fn.typed = function(option) {
        return this.each(function() {
            var $this = $(this),
                data = $this.data('typed'),
                options = typeof option == 'object' && option;
            if (!data) $this.data('typed', (data = new Typed(this, options)));
            if (typeof option == 'string') data[option]();
        });
    };

    $.fn.typed.defaults = {
        strings: ["These are the default values...", "You know what you should do?", "Use your own!", "Have a great day!"],
        // typing speed
        typeSpeed: 0,
        // time before typing starts
        startDelay: 0,
        // backspacing speed
        backSpeed: 0,
        // time before backspacing
        backDelay: 500,
        // loop
        loop: false,
        // false = infinite
        loopCount: false,
        // show cursor
        showCursor: true,
        // character for cursor
        cursorChar: "|",
        // attribute to type (null == text)
        attr: null,
        // either html or text
        contentType: 'html',
        // call when done callback function
        callback: function() {},
        // starting callback function before each string
        preStringTyped: function() {},
        //callback for every typed string
        onStringTyped: function() {},
        // callback for reset
        resetCallback: function() {}
    };


}(window.jQuery);

$(window).load(function(){
    $('.preloader').delay(1000) .fadeOut("slow"); // set duration in brackets
});

$(function(){
    /* start typed element */
    //http://stackoverflow.com/questions/24874797/select-div-title-text-and-make-array-with-jquery
    var subElementArray = $.map($('.sub-element'), function(el) { return $(el).text(); });
    $(".element").typed({
        strings: subElementArray,
        typeSpeed: 30,
        contentType: 'html',
        showCursor: false,
        loop: false,
        loopCount: true,
    });
    /* end typed element */

    new WOW().init();
});

var re = new RegExp(/^(http|https):\/\/([\w.]+\/?)\S*/);
function getShortUrl() {
    var long_url = document.getElementById("long_url").value
    var lowerCase = long_url.trim().toLowerCase()
    console.log(long_url)

    if (long_url == ""){
        /*document.getElementById("modal-body").setAttribute("style","")*/
        document.getElementById("pop-tab").innerHTML="<h3 class=\"modal-title\"  style=\"color: darkred\" id=\"pop-tab\">警告</h3>"
        document.getElementById("modal-body").innerHTML="<h2 style='color:red;letter-spacing: 2px;font-weight: 410;'>请填写长链接信息!!!!!!!</h2>"
    }else {
        if (re.test(lowerCase)){
            console.log(long_url)
            var parms = {"long_url":long_url}
            var url = "http://localhost:8080/"
            sendPost(url,parms)
        }else {
            console.log(long_url)
            document.getElementById("pop-tab").innerHTML="<h3 class=\"modal-title\"  style=\"color: darkred\" id=\"pop-tab\">警告</h3>"
            document.getElementById("modal-body").innerHTML="<h2 style='color:red;letter-spacing: 2px;font-weight: 410;'>请保证链接格式正确(必须以http/https 开头)!!!!!!!</h2>"
        }


    }
}

function sendPost(url,parms) {
    console.log("url:"+url + " ==> parms : "+JSON.stringify(parms))
    var httpRequest;//第一步：创建需要的对象
    if (window.XMLHttpRequest)
    {
        httpRequest=new XMLHttpRequest();
    }
    else if (window.ActiveXObject)
    {
        // IE
        httpRequest = new ActiveXObject("Microsoft.XMLHTTP");
    }
    httpRequest.open('POST', url, false); //第二步：打开连接/***发送json格式文件必须设置请求头 ；如下 - */
    httpRequest.setRequestHeader("Content-type","application/json");//设置请求头 注：post方式必须设置请求头（在建立连接后设置请求头）var obj = { name: 'zhansgan', age: 18 };
    /**
     * 获取数据后的处理程序
     */
    httpRequest.onreadystatechange = function () {//请求后的回调接口，可将请求成功后要执行的程序写在其中
        if (httpRequest.readyState == 4 && httpRequest.status == 200) {//验证请求是否发送成功
            var response = httpRequest.responseText
            var parse = JSON.parse(response);
            console.log("222222222222")
            console.log(parse.data)
            var code = parse.data.code
            var short_url = parse.data.short_url
            var chainInfo =  parse.data.chain_info
            var txid = chainInfo.txid
            var block_height = chainInfo.block_height
            showResult()
            document.getElementById("long_url").value=""
            document.getElementById("code").innerHTML=code
            document.getElementById("short_url_href").innerHTML = short_url
            document.getElementById("short_url_href").setAttribute("href",short_url)
            document.getElementById("txid").innerHTML = txid
            document.getElementById("block_height").innerHTML = block_height

        }
    };
    httpRequest.send(JSON.stringify(parms));//发送请求 将json写入send中
    //console.log(httpRequest)
}

function showResult() {
    var table = '<table>\n' +
        '                                    <tr>\n' +
        '                                        <td><label>压缩码：</label></td>\n' +
        '                                        <td><span id="code"></span></td>\n' +
        '                                    </tr>\n' +
        '                                    <tr>\n' +
        '                                        <td><label>短连接：</label></td>\n' +
        '                                        <td><span><a id="short_url_href" href="" target="_blank"></a></span>\n' +
        '                                        </td>\n' +
        '                                    </tr>\n' +
        '                                    <tr>\n' +
        '                                        <td><label>所在区块：</label></td>\n' +
        '                                        <td><span id="block_height"></span></td>\n' +
        '                                    </tr>\n' +
        '                                    <tr>\n' +
        '                                        <td><label>交易ID：</label></td>\n' +
        '                                        <td><span id="txid" ></span></td>\n' +
        '                                    </tr>\n' +
        '                                </table>'
    document.getElementById("pop-tab").innerHTML="<h3 class=\"modal-title\"  style=\"color:darkslategray\" id=\"pop-tab\">短链信息</h3>"
    document.getElementById("modal-body").innerHTML=table
}



