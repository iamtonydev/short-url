### Run:

* `git clone https://github.com/iamtonydev/url-shortener.git`
* `docker-compose build`
* `docker-compose up -d`

### Api:

Shorten url
    
* url **POST** `http://localhost:8000/urls/v1/add`
* data example 
```
{
    "url": "https://github.com/iamtonydev/url-shortener"
}
```
* response

```
{
    "result": {
        "short_url": "2txt3zUSuf"
    }
}
```

Get url

* url **GET** `http://localhost:8000/`
* request example `http://localhost:8000/2txt3zUSuf`
* response 
```
{
    "result": {
        "url": "https://github.com/iamtonydev/url-shortener"
    }
}
```