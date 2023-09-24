// ==UserScript==
// @name         Qin
// @namespace    http://tampermonkey.net/
// @version      1.0
// @description  Qin desc
// @author       You
// @match        https://studio.youtube.com/channel/*
// @icon         data:image/gif;base64,R0lGODlhAQABAAAAACH5BAEKAAEALAAAAAABAAEAAAICTAEAOw==
// @grant        unsafeWindow
// @run-at       document-end
// ==/UserScript==

let YTB_TOOL_CSS = `
#mock-ytb-tool {
    width: 200px;
    height: auto;
    position: fixed;
    right: 0;
    top: 100px;
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

.mock_ytb-icon {
    text-align: center;
}

#mock_video_num_on_day,
#mock_video_index,
#mock_day_upload_number,
#mock_video_date_schedule {
    width: 150px;
    margin-top: 5px;
    text-align: center;
}

#mock_import_config,
#mock_upload_video,
#mock_edit_video,
#mock-stop-notify {
    margin-top: 5px;
}`;

let stop_script = false;
let novels;
let video_index;
let day_upload_number = -1;
let video_number_every_day = 5;
let video_date_schedule = "";

// let audio_link = "https://webfs.ali.kugou.com/202307090851/3f3b73fa771392eb91a11db76fb2abd7/KGTX/CLTX003/7c9c160cd67cdf0acf68aac25a7ce2c2.mp3";
let audio_link = "https://webfs.ali.kugou.com/202307221610/e31e75a6b2f98001a8bff3f4720f5a93/v2/fc89e696a05a6c6b19c94aabe2739293/part/0/960117/G292/M00/29/C0/clip_BJUEAGSrsnyAGp_4AE0468SKETQ920.mp3";
let notify_audio;
let notify_audio_start = 21;

const async_await = function(wait_ms) {
    return new Promise((resolve, reject) => {
        return setTimeout(() => {
            resolve();
        }, wait_ms);
    });
}

const async_polling = function(wait_ms, executeFunc) {
    return new Promise((resolve, reject) => {
        const polling_func = function() {
            if (executeFunc()) {
                console.log('polling successful');
                resolve();
            } else {
                setTimeout(polling_func, wait_ms);
            }
        };
        return setTimeout(polling_func, wait_ms);
    });
}

const mock_input = function(elem, text) {
    // document.execCommand("insertText", !1, "sdsd\n书法和")
    let evt = document.createEvent('HTMLEvents');
    evt.initEvent('input', true, true);
    elem.value = text;
    elem.dispatchEvent(evt);
    elem.click();

    var event = document.createEvent('Event');
    // 注意这块触发的是keydown事件，在awx的ui源码中bind监控的是keypress事件，所以这块要改成keypress
    event.initEvent('keydown', true, false);
    event = Object.assign(event, {
        ctrlKey: false,
        metaKey: false,
        altKey: false,
        which: 13,
        keyCode: 13,
        key: 'Enter',
        code: 'Enter'
    });
    elem.focus();
    elem.dispatchEvent(event);
}

const mock_custom_input = function(elem, text) {
    let evt = document.createEvent('HTMLEvents');
    evt.initEvent('input', true, true);
    elem.innerText = "";
    elem.dispatchEvent(evt);
    elem.click();

    var event = document.createEvent('Event');
    // 注意这块触发的是keydown事件，在awx的ui源码中bind监控的是keypress事件，所以这块要改成keypress
    event.initEvent('keydown', true, false);
    event = Object.assign(event, {
        ctrlKey: false,
        metaKey: false,
        altKey: false,
        which: 13,
        keyCode: 13,
        key: 'Enter',
        code: 'Enter'
    });
    elem.focus();
    elem.dispatchEvent(event);

    document.execCommand("insertText", !1, text)
}

// 注入css样式
function inject_css() {
    let style = document.createElement("style");
    style.type = "text/css";
    style.textContent = YTB_TOOL_CSS;
    document.getElementsByTagName("head").item(0).appendChild(style);
}

// 注入窗口工具
function inject_window_tool() {
    let divElem = document.createElement('div');
    divElem.id = "mock-ytb-tool";
    document.querySelector('body').appendChild(divElem);

    // icon
    let iconElem = document.createElement('div');
    iconElem.className = "mock-ytb-icon";
    iconElem.innerText = "👻";
    divElem.appendChild(iconElem);

    // video number every day
    let videoNumOnDayIndexElem = document.createElement('input');
    videoNumOnDayIndexElem.id = "mock_video_num_on_day";
    videoNumOnDayIndexElem.type = "text";
    videoNumOnDayIndexElem.placeholder = "视频每日展示数";
    videoNumOnDayIndexElem.value = "5";
    divElem.appendChild(videoNumOnDayIndexElem);

    // video index
    let videoIndexElem = document.createElement('input');
    videoIndexElem.id = "mock_video_index";
    videoIndexElem.type = "text";
    videoIndexElem.placeholder = "视频索引号";
    videoIndexElem.value = "0";
    divElem.appendChild(videoIndexElem);

    // day upload number
    let dayUploadNumberElem = document.createElement('input');
    dayUploadNumberElem.id = "mock_day_upload_number";
    dayUploadNumberElem.type = "text";
    dayUploadNumberElem.placeholder = "第一个视频为排定日期的第几个";
    dayUploadNumberElem.value = "0";
    divElem.appendChild(dayUploadNumberElem);

    // video date schedule
    let videoDateScheduleElem = document.createElement('input');
    videoDateScheduleElem.id = "mock_video_date_schedule";
    videoDateScheduleElem.type = "text";
    videoDateScheduleElem.placeholder = "第一个视频为排定日期的第几个";
    videoDateScheduleElem.value = "2023/8/8";
    divElem.appendChild(videoDateScheduleElem);

    // import config
    let importConfigElem = document.createElement('input');
    importConfigElem.id = "mock_import_config";
    importConfigElem.value = "导入配置";
    importConfigElem.type = "file";
    setNovelsChange(importConfigElem);
    divElem.appendChild(importConfigElem);

    // upload video
    let uploadElem = document.createElement('button');
    uploadElem.id = "mock_upload_video";
    uploadElem.innerText = "上传视频";
    uploadElem.type = "button";
    setUploadButtonClick(uploadElem);
    divElem.appendChild(uploadElem);
}

// 设置小说配置更改
function setNovelsChange(config_elem) {
    config_elem.addEventListener('change', () => {
		const reader = new FileReader();
		reader.readAsText(config_elem.files[0], 'utf-8');
		reader.onload = () => {
			var local_novels = JSON.parse(reader.result);
            if (local_novels) {
                novels = local_novels;
            }
		};
    });
}

// 设置上传按钮点击
async function setUploadButtonClick(uploadElem) {
    uploadElem.onclick = function() {
        let upload_elem = document.querySelector('#mock_upload_video');
        if (upload_elem.innerText == "停止脚本") {
            stop_script = true;
            upload_elem.innerText = "上传视频";
            return;
        }

        upload_elem.innerText = "停止脚本";
        try {
            let video_index_elem = document.querySelector('#mock_video_index');
            let parse_video_index = parseInt(video_index_elem.value);
            if (novels && parse_video_index >= 0 && parse_video_index < novels.length) {
                video_index = parse_video_index;
                if (day_upload_number == -1) {
                    let day_upload_number_elem = document.querySelector('#mock_day_upload_number');
                    let local_day_upload_number = parseInt(day_upload_number_elem.value);
                    if (local_day_upload_number >= 0 && local_day_upload_number < 5) {
                        day_upload_number = local_day_upload_number;
                    } else {
                        alert("day upload number incorret!(" + local_day_upload_number + ")");
                        return;
                    }
                }

                if (video_date_schedule == "") {
                    let video_date_schedule_elem = document.querySelector('#mock_video_date_schedule');
                    if (video_date_schedule_elem.value && video_date_schedule_elem.value != "") {
                        video_date_schedule = video_date_schedule_elem.value;
                    } else {
                        alert("video date schedule incorret!");
                        return;
                    }
                }
                mockClickUploadVideo();
                executeUploadScript(0);
            } else {
                console.log("not found novels json");
                alert("Please load config file first!");
            }
        } catch(e) {
            alert("video index incorret!");
            return;
        }
    }
}

// mock点击上传视频
async function mockClickUploadVideo() {
    // 点击创建
    document.querySelector('#create-icon').click();

    await async_await(500);

    // 点击上传视频
    document.querySelector('#text-item-0').click();

    await async_await(500);

    // 点击选择文件
    document.querySelector("[name='Filedata']").accept = ".mp4";
    await async_await(200);
    document.querySelector('#select-files-button').click();
}

// 统计上传
function statisticalUpload() {
    let video_number_every_day = document.querySelector('#mock_video_num_on_day');
    let video_number = parseInt(video_number_every_day.value);
    day_upload_number = (day_upload_number + 1) % video_number;
    document.querySelector('#mock_day_upload_number').value = day_upload_number;
    if (day_upload_number == 0) {
        let curUploadDate = new Date(video_date_schedule);
        let nextDayDate = new Date(curUploadDate.getTime() + 24 * 60 * 60 * 1000);
        // js计算的月份比实际小一个月
        video_date_schedule = nextDayDate.getFullYear() + "/" + (nextDayDate.getMonth() + 1) + "/" + nextDayDate.getDate();
        document.querySelector('#mock_video_date_schedule').value = video_date_schedule;
    }
}

function notifyAutoUploadFinish() {
    console.log("auto upload finish");
    if (notify_audio) {
        notify_audio.currentTime = 0;
    } else {
        notify_audio = new Audio(audio_link);
    }
    notify_audio.currentTime = notify_audio_start;
	notify_audio.play();

    let mock_container = document.querySelector('#mock-ytb-tool');
    let mock_stop_notify_elem = document.createElement('button');
    mock_stop_notify_elem.id = "mock-stop-notify";
    mock_stop_notify_elem.innerText = "停止通知";
    mock_stop_notify_elem.type = "button";
    mock_stop_notify_elem.onclick = stopAutoUploadNotify;
    mock_container.appendChild(mock_stop_notify_elem);

    let upload_elem = document.querySelector('#mock_upload_video');
    upload_elem.innerText = "上传视频";
    stop_script = false;
}

function stopAutoUploadNotify() {
    if (notify_audio) {
        notify_audio.pause();
        notify_audio.currentTime = 0;
    }

    let stop_notify_elem = document.querySelector('#mock-stop-notify');
    stop_notify_elem.remove();
}

// 执行上传脚本
async function executeUploadScript(offset_index) {
    console.log('start detect upload video');

    await async_polling(500, function() {
        let progress_dialog_title = document.querySelector('.count');
        return progress_dialog_title && progress_dialog_title.innerText == "上传完毕" &&
            document.querySelector('ytcp-uploads-mini-indicator').querySelector('#dialog').style.display != 'none';
    });

    let novel = novels[video_index];
    let progress_items = document.querySelectorAll('#progress-list li');
progress_item_for:
    for (var index = 0; index < progress_items.length; index++) {
        if (index < offset_index) {
            continue;
        }
        // 停止脚本
        if (stop_script) {
            stop_script = false;
            break;
        }

        console.log('click index', index);
        let progress_item = progress_items[index];
        await async_polling(500, function() {
            return progress_item.querySelector(".progress-title");
        });

        const progress_title_elem = progress_items[index].querySelector('.progress-title');
        const progress_title = progress_title_elem.innerText;
        const progress_index = progress_title.replace('ch-', '').replace('.mp4', '');
        // console.log("progress_index:", progress_index);

        for (var chapter of novel.chapters) {
            if (progress_index == chapter.index) {
                console.log("chapter:", chapter);

                // 点击上传列表标题
                await async_await(5000);
                progress_title_elem.click();

                // 输入标题
                await async_polling(100, function() {
                    let title_elem = document.querySelector('.input-container.title #textbox');
                    if (title_elem) {
                        return true;
                    } else {
                        return false;
                    }
                });
                await async_await(300);
                const title_input = document.querySelector('.input-container.title #textbox');
                mock_custom_input(title_input, chapter.title);

                // 输入描述
                await async_await(1000);
                const desc_input = document.querySelector('.input-container.description #textbox');
                mock_custom_input(desc_input, chapter.desc);

                // 点击播放列表
                await async_await(1000);
                document.querySelector('.compact-row .dropdown').click();

                // 输入搜索播放列表
                await async_await(2000);
                const playlist_search = document.querySelector('#search-input');
                mock_input(playlist_search, novel.name);

                // 选中播放列表
                await async_await(500);
                const playlist_items = document.querySelectorAll('ytcp-playlist-dialog ytcp-ve');
                for (var playlist_item of playlist_items) {
                    if (!playlist_item.attributes.hasOwnProperty('hidden')) {
                        playlist_item.querySelector('label').click();
                        await async_await(500);
                    }
                }

                // 点击播放列表的完成按钮
                document.querySelector('ytcp-playlist-dialog .done-button').click();

                // 点击继续按钮直到保存按钮出现
                await async_polling(2000, function() {
                    let continue_button = document.querySelector('#next-button');
                    if (continue_button) {
                        if (!continue_button.attributes.hasOwnProperty("hidden")) {
                            continue_button.click();
                            return false;
                        }
                    } else {
                        return false;
                    }
                    return true;
                });

                // 点击安排时间选项
                await async_await(1000);
                let time_schedule = document.querySelector('#second-container-expand-button');
                if (time_schedule) {
                    time_schedule.click();
                } else {
                    document.querySelector('#second-container #radioContainer').click();
                }

                // 点击安排时间栏
                await async_await(1000);
                document.querySelector('#datepicker-trigger .right-container').click();

                // 输入安排时间
                await async_await(1000);
                let date_input = document.querySelector('ytcp-date-picker input');
                await async_await(1000);
                mock_input(date_input, video_date_schedule);

                // 隐藏时间选择弹框
                await async_await(1000);
                document.querySelector('body>tp-yt-iron-overlay-backdrop').click();

                // 点击时刻栏
                await async_await(1000);
                document.querySelector('#time-of-day-container #textbox').click();

                // 设置时间
                await async_await(1000);
                // let time_input = document.querySelector('#labelAndInputContainer input');
                // mock_input(time_input, "08:00");
                document.querySelectorAll('tp-yt-paper-listbox .ytcp-time-of-day-picker')[32].click();

                // 隐藏时间选择弹框
                await async_await(1000);
                document.querySelector('body>tp-yt-iron-overlay-backdrop').click();

                // 点击预定按钮
                await async_await(1000);
                document.querySelector('#done-button').click();

                // 点击分享关闭按钮
                await async_polling(1000, function() {
                    return document.querySelector('ytcp-video-share-dialog #close-button .label') ||
                        document.querySelector('ytcp-uploads-still-processing-dialog #close-button .label');;
                });

                await async_await(2000);
                let share_dialog_colse = document.querySelector('ytcp-video-share-dialog #close-button .label');
                if (share_dialog_colse) {
                    share_dialog_colse.click();
                }
                let upload_process_dialog_colse = document.querySelector('ytcp-uploads-still-processing-dialog #close-button .label');
                if (upload_process_dialog_colse) {
                    upload_process_dialog_colse.click();
                }

                statisticalUpload();
                continue progress_item_for;
            }
        }
    }

    // 上传完毕关闭上传弹窗
    await async_await(2000);
    document.querySelector('ytcp-uploads-mini-indicator #close-button').click();
    notifyAutoUploadFinish();
}

(function() {
    'use strict';
    window.onload = function() {
        inject_css();
        inject_window_tool();
        window.mock_tool = function(offset_index) {
            console.log("offset_index:", offset_index);
            executeUploadScript(offset_index);
        };
    }
    // Your code here...
})();