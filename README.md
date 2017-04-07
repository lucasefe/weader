# weader 

## Install

```
glide install
```

## Run server 

```
WEATHER_API_KEY=<<KEY>> go run main.go
```

## Consume

Using [httpie](https://github.com/jakubroztocil/httpie).

```
http  localhost:8081/lucasefe
```


```
HTTP/1.1 200 OK
Content-Length: 98
Content-Type: application/json
Date: Fri, 07 Apr 2017 01:55:24 GMT
X-Runtime: 20.448870275s

{
    "avg_temperature": 22,
    "location": "Buenos Aires, Argentina",
    "repos_count": 30,
    "username": "lucasefe"
}
```
