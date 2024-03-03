// ==UserScript==
// @name         colab script
// @namespace    http://tampermonkey.net/
// @version      2024-02-07
// @description  try to take over the world!
// @author       You
// @match        https://colab.research.google.com/drive/1Nb6pWbakc9nH7qYxnEjFz1mTyMR_bE-n
// @icon         https://www.google.com/s2/favicons?sz=64&domain=google.com
// @grant        none
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
    scrapySwitchElem.innerText = "持续开关";
    scrapySwitchElem.type = "range";
    scrapySwitchElem.name = "持续开关";
    scrapySwitchElem.min = 0;
    scrapySwitchElem.max = 1;
    scrapySwitchElem.value = disable_scrapy ? 0 : 1;
    setSwitchClick(scrapySwitchElem);
    divElem.appendChild(scrapySwitchElem);

    // scrapy switch desc
    let scrapySwitchDescElem = document.createElement('span');
    scrapySwitchDescElem.innerText = "持续开关";
    divElem.appendChild(scrapySwitchDescElem);
}

function setSwitchClick(scrapySwitchElem) {
    scrapySwitchElem.addEventListener('change', function(event) {
        if (localStorage[`enable-scrapy`] !== "1") {
            localStorage[`enable-scrapy`] = "1";
            inject_script();
        } else {
            localStorage[`enable-scrapy`] = "0";
        }
    });
}

async function inject_script() {
    let timer = setInterval(function() {
        console.log("Clicked on connect button");
        document.querySelector("colab-connect-button").shadowRoot.querySelector('#connect').click();

        let disable_scrapy = localStorage[`enable-scrapy`] !== "1";
        if (disable_scrapy) {
            console.log("cancel interval");
            clearInterval(timer);
        }
    },30000);
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