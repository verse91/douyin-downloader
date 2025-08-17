# Document

This document explains how to extract the final video page from a Douyin (TikTok China) original link using `curl` to understand the code.  
<div align="center">
  <img width="2192" height="1128" alt="image" src="https://github.com/user-attachments/assets/790d14aa-13ff-4e96-9730-bc764a742419" />
  <p><i>Douyin downloader with reponse as JSON</i></p>
</div>

> [!WARNING]  
> **This is for educational purposes only.**  
> **I do not take responsibility for anything you do.**

## Steps

1. **Send a request to the Douyin short link**
  - Example: `"https://v.douyin.com/9N2HGwrYB70/"`

     ```bash
     curl -A "Mozilla/5.0 (Linux; Android 11; SAMSUNG SM-G973U) AppleWebKit/537.36 (KHTML, like Gecko) SamsungBrowser/14.2 Chrome/87.0.4280.141 Mobile Safari/537.36" "https://v.douyin.com/9N2HGwrYB70/"
     ```

   - This will return an HTML page containing a redirect link, for example:

     ```html
     <a href="https://www.iesdouyin.com/share/video/7525380501322911011/?region=VN&amp;mid=7525380482746403647&amp;u_code=-1&amp;did=MS4wLjABAAAAvZfZRVqhGh6ryU2DJAnJ7bhfhSWSy6xM9wpY4M1ZHFfLtlx3ij92SH1eSqQrmUwA&amp;iid=MS4wLjABAAAANwkJuWIRFOzg5uCpDRpMj4OX-QryoDgn-yYlXQnRwQQ&amp;with_sec_did=1&amp;video_share_track_ver=&amp;titleType=title&amp;share_sign=DLXFtSzzzFpSxbbCrKvRU1RDDaexZLkWJmOPOhlqn.g-&amp;share_version=190500&amp;ts=1754977875&amp;from_aid=6383&amp;from_ssr=1&amp;from=web_code_link">Found</a>
     ```

2. **Extract the redirect URL from the response**

   - Copy the URL from the `href` attribute in the `<a>` tag.

3. **Send a request to the extracted URL**

   ```bash
   curl -A "Mozilla/5.0 (Linux; Android 11; SAMSUNG SM-G973U) AppleWebKit/537.36 (KHTML, like Gecko) SamsungBrowser/14.2 Chrome/87.0.4280.141 Mobile Safari/537.36" "https://www.iesdouyin.com/share/video/7525380501322911011/?region=VN&mid=7525380482746403647&u_code=-1&did=MS4wLjABAAAAvZfZRVqhGh6ryU2DJAnJ7bhfhSWSy6xM9wpY4M1ZHFfLtlx3ij92SH1eSqQrmUwA&iid=MS4wLjABAAAANwkJuWIRFOzg5uCpDRpMj4OX-QryoDgn-yYlXQnRwQQ&with_sec_did=1&video_share_track_ver=&titleType=title&share_sign=DLXFtSzzzFpSxbbCrKvRU1RDDaexZLkWJmOPOhlqn.g-&share_version=190500&ts=1754977875&from_aid=6383&from_ssr=1&from=web_code_link" > resp.html
   ```

   - This will save the response HTML to `resp.html`.

4. **Process the response**

   - Open `resp.html` and extract the information you need (such as video URL, metadata, etc.).
   - Check mine [here](https://github.com/verse91/douyin-downloader/blob/main/resp.txt).
## Result
1. Video
```
{
  "media": {
    "type": "video",
    "video_id": "7525380501322911011",
    "like": 746,
    "comment": 52,
    "save": 75,
    "share": 76,
    "video_desc": "泸沽湖小落水银河星空  #跟着中国国家地理镜收山河 #一起去收集中国限量版地貌  #泸沽湖小落水村",      
    "create_at": "2025-07-10 16:19:57",
    "video_url": "https://www.iesdouyin.com/aweme/v1/play/?video_id=v0200fg10000d1no9qvog65ikp9k6fr0&ratio=1080p&line=0"
  },
  "user": {
    "name": "骑驴到西藏",
    "bio": "畅销书《星空摄影与后期》作者\\n8KRAW签约摄影师\\n微信：wzc554593239",
    "username": "qilv919"
  }
}
```
2. Image
```
{
  "media": {
    "type": "img",
    "video_id": "7536825496672701738",
    "like": 878,
    "comment": 175,
    "save": 65,
    "share": 63,
    "video_desc": "常山阴 我不恨你\\n#大爱仙尊#蛊真人#古月方源",
    "create_at": "2025-08-10 12:32:16",
    "img_url": [
      "https://p3-sign.douyinpic.com/tos-cn-i-0813c001/oUvY6uaWF3g8PIgAXgCnNILA5iAAJBtiIEcAd~tplv-dy-resize-walign-adapt-aq:540:q75.webp?lk3s=138a59ce&x-expires=1756206000&x-signature=D0Z3wpV2zhfhLykHVLTO0UFIwRk%3D&from=327834062&s=PackSourceEnum_DOUYIN_REFLOW&se=false&sc=cover&biz_tag=aweme_images&l=20250812191910A7A5CC547B14D58D1C2D",
      "https://p3-sign.douyinpic.com/tos-cn-i-0813c001/owEdYvAPIg7IC83nggNi6iytWvBaJJAdAAuAL~tplv-dy-lqen-new:1440:1920:q80.webp?lk3s=138a59ce&x-expires=1757588400&x-signature=lZcNdqmWhUCEOFRf0hYlUBIk3ak%3D&from=327834062&s=PackSourceEnum_DOUYIN_REFLOW&se=false&sc=image&biz_tag=aweme_images&l=20250812191910A7A5CC547B14D58D1C2D"
    ]
  },
  "user": {
    "username": "llllhhhcccc",
    "name": "Rainco",
    "bio": "魔丸@saber \\n原图都在群里可以进来玩"
  }
}
```
> [!TIP]
> Add more information like music, thumbnail, avatar, ... what you want from the response

## How to run

- Clone the repo
```bash
git clone https://github.com/verse91/douyin-downloader/
cd /douyin-downloader/
```

- Init project
  
```go
go mod init douyin
```

- Install packages

```go
go mod tidy
```

- Run

```go
go run .
```

- OR build

```go
go build .
```

> [!NOTE]  
> Always use a mobile User-Agent to avoid being blocked or served different content.
> You can automate the extraction of the redirect URL using tools like `grep`, `sed`, or scripting languages.

---
