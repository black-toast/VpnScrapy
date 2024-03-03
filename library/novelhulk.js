// ==UserScript==
// @name         novelhulk script
// @namespace    http://tampermonkey.net/
// @version      2024-02-01
// @description  try to take over the world!
// @author       You
// @match        https://novelhulk.com/nb/*
// @match        https://allnovelnext.com/allnovelnext/*
// @match        https://novelnextz.com/novelnextz/*
// @icon         https://www.google.com/s2/favicons?sz=64&domain=novelhulk.com
// @grant        none
// @grant        unsafeWindow
// @run-at       document-end
// ==/UserScript==
let CSS_CONTENT = `
#scrapy-ui {
    width: 200px;
    height: auto;
    position: fixed;
    right: 0;
    top: 300px;
    text-align: center;
    display: inline-block;
    z-index: 450;
    border-radius: 40px 0 0 40px;
    box-shadow: -1px 4px 8px rgba(0,0,0,.06);
    background-color: white;
    padding: 10px 20px;
    box-sizing: border-box;
    z-index: 9999;
    opacity: 0.8;
}

#scrapy-ui .icon {
    text-align: center;
}
#scrapy-ui #scrapy-switch {
    display: inline-block;
    width: 25px;
    margin-right: 5px;
}

#scrapy-end-index,
#download-list-len {
    width: 150px;
    margin-top: 5px;
    text-align: center;
}
`;

// 注入css样式
function inject_css() {
    let style = document.createElement("style");
    style.type = "text/css";
    style.textContent = CSS_CONTENT;
    document.getElementsByTagName("head").item(0).appendChild(style);
}

// 注入窗口工具
function inject_window_tool(disable_scrapy) {
    let divElem = document.createElement('div');
    divElem.id = "scrapy-ui";
    document.querySelector('body').appendChild(divElem);

    // icon
    let iconElem = document.createElement('div');
    iconElem.className = "icon";
    iconElem.innerText = "👻";
    divElem.appendChild(iconElem);

    // scrapy switch
    let scrapySwitchElem = document.createElement('input');
    scrapySwitchElem.id = "scrapy-switch";
    scrapySwitchElem.innerText = "上传视频";
    scrapySwitchElem.type = "range";
    scrapySwitchElem.name = "爬虫开关";
    scrapySwitchElem.min = 0;
    scrapySwitchElem.max = 1;
    scrapySwitchElem.value = disable_scrapy ? 0 : 1;
    setSwitchClick(scrapySwitchElem);
    divElem.appendChild(scrapySwitchElem);

    // scrapy switch desc
    let scrapySwitchDescElem = document.createElement('span');
    scrapySwitchDescElem.innerText = "爬虫开关";
    divElem.appendChild(scrapySwitchDescElem);

    // scrapy end index
    let scrapyEndIndexElem = document.createElement('input');
    scrapyEndIndexElem.id = "scrapy-end-index";
    scrapyEndIndexElem.type = "text";
    scrapyEndIndexElem.placeholder = "爬取结束索引";
    let scrapyEndIndex = localStorage["scrapy-end-index"];
    if (typeof scrapyEndIndex === 'undefined') {
        scrapyEndIndexElem.value = "0";
    } else {
        scrapyEndIndexElem.value = scrapyEndIndex;
    }
    divElem.appendChild(scrapyEndIndexElem);

    // download list len
    let downloadListLenElem = document.createElement('input');
    downloadListLenElem.id = "download-list-len";
    downloadListLenElem.type = "text";
    downloadListLenElem.placeholder = "爬取n个章节后保存";
    let downloadListLen = localStorage["download-list-len"];
    if (typeof downloadListLen === 'undefined') {
        downloadListLenElem.value = "0";
    } else {
        downloadListLenElem.value = downloadListLen;
    }
    divElem.appendChild(downloadListLenElem);
}

function setSwitchClick(scrapySwitchElem) {
    scrapySwitchElem.addEventListener('change', function(event) {
        if (localStorage[`enable-scrapy`] !== "1") {
            localStorage[`enable-scrapy`] = "1";
            localStorage["scrapy-end-index"] = document.querySelector("#scrapy-end-index").value;
            localStorage["download-list-len"] = document.querySelector("#download-list-len").value;
            inject_script();
        } else {
            localStorage[`enable-scrapy`] = "0";
        }
    });
}

function async_await(wait_ms) {
    return new Promise((resolve, reject) => {
        return setTimeout(() => {
            resolve();
        }, wait_ms);
    });
}

function async_polling(wait_ms, executeFunc) {
    return new Promise((resolve, reject) => {
        const polling_func = function() {
            if (executeFunc()) {
                //console.log('polling successful');
                resolve();
            } else {
                setTimeout(polling_func, wait_ms);
            }
        };
        return setTimeout(polling_func, wait_ms);
    });
}

function download(downloadContent, fileName) {
    const downloadLink = document.createElement("a");
    downloadLink.download = fileName;
    downloadLink.style.display = 'none';

    let blob = new Blob([downloadContent]);
    downloadLink.href = URL.createObjectURL(blob);

    document.body.appendChild(downloadLink);
    downloadLink.click();
    document.body.removeChild(downloadLink);
}

async function inject_script() {
    // 点击加载章节列表
    document.querySelectorAll('button.chr-jump')[0].click();

    // 等待章节列表展示出来
    await async_polling(200, function() {
        return document.querySelectorAll('select.chr-jump').length > 0;
    });
    let chapterListNodes = document.querySelectorAll('select.chr-jump')[0].children;

    // 匹配章节列表和当前章节的序号
    var matchChapterIndex = -1;
    let chapterTitle = document.querySelectorAll('.chr-text')[0].innerText;
    for (var index = 0; index < chapterListNodes.length; index++) {
        let chapterTitleText = chapterListNodes[index].innerText;
        if (chapterTitleText == chapterTitle) {
            matchChapterIndex = index;
            break;
        }
    }

    // log match chapter index
    if (matchChapterIndex == -1) {
        console.log("not find match chapter index.");
        return;
    }
    console.log("match chapter index:" + (index + 1));

    // 提取章节内容
    let chapterLines = document.querySelectorAll('#chr-content p');
    if (chapterLines.length <= 0) {
        console.log("not found chapter");
        return;
    }
    var chapterContent = `Chapter ${index + 1}\n`;
    for (var lineIndex = 0; lineIndex < chapterLines.length; lineIndex++) {
        let chapterLine = chapterLines[lineIndex].innerText;
        chapterLine = chapterLine.replace("\ufeff", "").replace(/ magic$/, '').replace(/magic$/, '');
        if (chapterLine == "" || chapterLine == "---" || chapterLine == "…" || chapterLine == "……" || chapterLine == "***"
            || chapterLine == "------" || chapterLine == "t"
            || chapterLine.match(/^Chapter \d{1,}/) || chapterLine.match(/^Translator:/) || chapterLine.match(/^Editor:/)) {
            continue;
        }

        chapterContent += `${chapterLine}\n`;
    }


    let scrapyEndIndex = localStorage["scrapy-end-index"];
    if (typeof scrapyEndIndex === 'undefined' || index + 1 >= scrapyEndIndex) {
        localStorage[`enable-scrapy`] = "0";
        document.querySelector('#scrapy-switch').value = "0";
        console.log(`end scrapy chapter index > ${scrapyEndIndex}`);
        return;
    }
    let downloadListLen = localStorage["download-list-len"];
    let chapterIndex = index % downloadListLen;
    localStorage[`chapter-${chapterIndex}`] = chapterContent;
    if (chapterIndex == downloadListLen - 1) {
        var mergedChapter = "";
        for (var downloadIndex = 0; downloadIndex < downloadListLen; downloadIndex++) {
            mergedChapter += localStorage[`chapter-${downloadIndex}`];
            localStorage.removeItem(`chapter-${downloadIndex}`);
        }
        download(mergedChapter, `ch-${index + 1 - downloadListLen + 1}-${index + 1}.txt`);
    }

    let nextElem = document.querySelectorAll('#next_chap_top.btn[disabled]');
    if (nextElem.length > 0) {
        localStorage[`enable-scrapy`] = "0";
        document.querySelector('#scrapy-switch').value = "0";
        console.log("next chapter end.");
        return;
    }
    await async_await(2000);
    document.querySelectorAll('#next_chap_top.btn')[0].click();
}

(function() {
    'use strict';

    // Your code here...
    window.onload = async function() {
        let disable_scrapy = localStorage[`enable-scrapy`] !== "1";
        inject_css();
        inject_window_tool(disable_scrapy);
        if (disable_scrapy) {
            console.log("disable-scrapy");
            return;
        }
        inject_script();
    }
})();