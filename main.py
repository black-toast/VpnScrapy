from crawl.novel.crawl import Crawl
from crawl.novel.WuxiaCrawl import WuxiaCrawl
import os
import sys

generate_novels_json = True

crawl_novels = (
    {
        "title": "Cultivation Online", "novelPath": "cultivation-online-novel",
        "start": 0, "end": 0, "crawl": False, "makeMp3": False, "makeMp4": False
    },
    {
        "title": "Dual Cultivation", "novelPath": "dual-cultivation-novel",
        "start": 0, "end": 0, "crawl": False, "makeMp3": False, "makeMp4": False
    },
    {
        "title": "Forty Millenniums of Cultivation", "novelPath": "forty-millenniums-of-cultivation-novel",
        "start": 0, "end": 0, "crawl": False, "makeMp3": False, "makeMp4": False
    },
    {
        "title": "Immortal and Martial Dual Cultivation", "novelPath": "immortal-and-martial-dual-cultivation",
        "start": 0, "end": 0, "crawl": False, "makeMp3": False, "makeMp4": False
    },
    {
        "title": "Invincible Divine Dragon's Cultivation System", "novelPath": "invincible-divine-dragons-cultivation-system",
        "start": 0, "end": 0, "crawl": False, "makeMp3": False, "makeMp4": False
    },
    {
        "title": "Legend of the Great Sage", "novelPath": "legend-of-the-great-sage",
        "start": 0, "end": 0, "crawl": False, "makeMp3": False, "makeMp4": False
    },
    {
        "title": "Lucifer's Descendant System", "novelPath": "lucifers-descendant-system",
        "start": 0, "end": 0, "crawl": False, "makeMp3": False, "makeMp4": False
    },
    {
        "title": "Necropolis Immortal", "novelPath": "necropolis-immortal",
        "start": 0, "end": 0, "crawl": False, "makeMp3": False, "makeMp4": False
    },
    {
        "title": "Shadow Slave", "novelPath": "shadow-slave",
        "start": 0, "end": 0, "crawl": False, "makeMp3": False, "makeMp4": False
    },
    {
        "title": "Supreme Lord I can extract everything", "novelPath": "supreme-lord-i-can-extract-everything",
        "start": 0, "end": 0, "crawl": False, "makeMp3": False, "makeMp4": False
    },
    {
        "title": "The Death Mage Who Doesn't Want a Fourth Time", "novelPath": "the-death-mage-who-doesnt-want-a-fourth-time",
        "start": 0, "end": 0, "crawl": False, "makeMp3": False, "makeMp4": False
    },
    {
        "title": "Why Should I Stop Being a Villain", "novelPath": "why-should-i-stop-being-a-villain",
        "start": 0, "end": 0, "crawl": False, "makeMp3": False, "makeMp4": False
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
            crawl.offlineMake(novel["title"], novel["novelPath"], novel["start"], novel["end"],
                              novel["makeMp3"], novel["makeMp4"])
        else:
            crawl.scrapy(novel["title"], novel["novelPath"], novel["start"], novel["end"],
                         novel["makeMp3"], novel["makeMp4"])
    
    for novel in wuxiacrawl_novels:
        if novel["crawl"] == False:
            wuxiaCrawl.offlineMake(novel["title"], novel["novelPath"], novel["start"], novel["end"],
                              novel["makeMp3"], novel["makeMp4"])
        else:
            wuxiaCrawl.scrapy(novel["title"], novel["novelPath"], novel["start"], novel["end"],
                         novel["makeMp3"], novel["makeMp4"])

    if generate_novels_json:
        crawl.generateNovelsJson()
