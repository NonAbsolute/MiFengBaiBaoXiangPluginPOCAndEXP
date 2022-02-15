Content-Type: application/json

## 保存规则

```http
/wp-json/wp-beebox/v1/crawler
```

```json
{
    "action": "save_single_rule",
    "data": [{
        "name": "baidu",
        "url": "baidu.com",
        "title": "2|h1",
        "content": "2|main",
        "inner_image": "2|img|src",
        "encode": "utf8"
    }],
    "target": "single"
}
```

```json
{
    "success": true,
    "data": {
        "target": "crawler_single_rule",
        "value": [{
            "name": "baidu",
            "url": "baidu.com",
            "title": "2|h1",
            "content": "2|main",
            "inner_image": "2|img|src",
            "encode": "utf8"
        }]
    }
}
```

## 查询规则

```http
/wp-json/wp-beebox/v1/crawler
```

```json
{"action":"get_rule"}
```

```json
{
    "success": true,
    "data": {
        "single": [],
        "list": []
    }
}
```

