from crawl.novel.crawl import Crawl
from crawl.novel.WuxiaCrawl import WuxiaCrawl
import os
import sys

generate_novels_json = True

crawl_novels = (
    {
        "novelPath": "cultivation-online-novel", "start": 0, "end": 0,
        "crawl": False, "makeMp3": False, "makeMp4": False
    },
    {
        "novelPath": "dual-cultivation-novel", "start": 0, "end": 0,
        "crawl": False, "makeMp3": False, "makeMp4": False
    },
    {
        "novelPath": "forty-millenniums-of-cultivation-novel", "start": 0, "end": 0,
        "crawl": False, "makeMp3": False, "makeMp4": False
    },
    {
        "novelPath": "immortal-and-martial-dual-cultivation", "start": 0, "end": 0,
        "crawl": False, "makeMp3": False, "makeMp4": False
    },
    {
        "novelPath": "invincible-divine-dragons-cultivation-system", "start": 0, "end": 0,
        "crawl": False, "makeMp3": False, "makeMp4": False
    },
    {
        "novelPath": "legend-of-the-great-sage", "start": 0, "end": 0,
        "crawl": False, "makeMp3": False, "makeMp4": False
    },
    {
        "novelPath": "lucifers-descendant-system", "start": 0, "end": 0,
        "crawl": False, "makeMp3": False, "makeMp4": False
    },
    {
        "novelPath": "MarvelsSuperman", "start": 0, "end": 0,
        "crawl": False, "makeMp3": False, "makeMp4": False
    },
    {
        "novelPath": "necropolis-immortal", "start": 0, "end": 0,
        "crawl": False, "makeMp3": False, "makeMp4": False
    },
    {
        "novelPath": "shadow-slave", "start": 0, "end": 0,
        "crawl": False, "makeMp3": False, "makeMp4": False
    },
    {
        "novelPath": "spy-mage-system", "start": 0, "end": 0,
        "crawl": False, "makeMp3": False, "makeMp4": False
    },
    {
        "novelPath": "supreme-lord-i-can-extract-everything", "start": 0, "end": 0,
        "crawl": False, "makeMp3": False, "makeMp4": False
    },
    {
        "novelPath": "the-death-mage-who-doesnt-want-a-fourth-time", "start": 0, "end": 0,
        "crawl": False, "makeMp3": False, "makeMp4": False
    },
    {
        "novelPath": "the-experimental-log-of-the-crazy-lich", "start": 0, "end": 0,
        "crawl": False, "makeMp3": False, "makeMp4": False
    },
    {
        "novelPath": "trial-marriage-husband-need-to-work-hard", "start": 0, "end": 0,
        "crawl": False, "makeMp3": False, "makeMp4": False
    },
    {
        "novelPath": "why-should-i-stop-being-a-villain", "start": 0, "end": 0,
        "crawl": False, "makeMp3": False, "makeMp4": False
    },
)

wuxiacrawl_novels = (
    {
        "title": "Lord of the Mysteries", "novelPath": "lord-of-the-mysteries",
        "start": 0, "end": 0, "crawl": False, "makeMp3": False, "makeMp4": False
    },
    {
        "title": "Ancient Strengthening Technique", "novelPath": "ancient-strengthening-technique",
        "start": 0, "end": 0, "crawl": False, "makeMp3": False, "makeMp4": False
    }
)

if __name__ == '__main__':
    crawl = Crawl()
    wuxiaCrawl = WuxiaCrawl()
    for novel in crawl_novels:
        if novel["crawl"] == False:
            crawl.offlineMake(novel["novelPath"], novel["start"], novel["end"],
                              novel["makeMp3"], novel["makeMp4"])
        else:
            crawl.scrapy(novel["novelPath"], novel["start"], novel["end"],
                         novel["makeMp3"], novel["makeMp4"])
    
    for novel in wuxiacrawl_novels:
        if novel["crawl"] == False:
            crawl.offlineMake(novel["title"], novel["novelPath"], novel["start"], novel["end"],
                              novel["makeMp3"], novel["makeMp4"])
        else:
            wuxiaCrawl.scrapy(novel["title"], novel["novelPath"], novel["start"], novel["end"],
                         novel["makeMp3"], novel["makeMp4"])

    if generate_novels_json:
        crawl.generateNovelsJson()
