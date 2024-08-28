# go-dict-tagging
Annotate natural sentences based on a dictionary.

## Example
1. request 
``` sh
http://127.0.0.1:8080/tag?statement=我有糖尿病，在吃格华止，可以同时吃奥利司他来减肥吗？
```
2. response
``` json
{
    "code": 1,
    "msg": "",
    "result": {
        "奥利司他": {
            "keyword": "奥利司他",
            "positions": [
                {
                    "start": 17,
                    "end": 20
                }
            ],
            "dictwords": {
                "medicine": [
                    {
                        "dict": "medicine",
                        "word": "奥利司他胶囊",
                        "index": [
                            "奥利司他",
                            "奥利司他胶囊"
                        ],
                        "data": {
                            "category": "西药",
                            "categoryId": 0,
                            "isControl": true,
                            "isHighRisk": true,
                            "medicineType": 0
                        }
                    },
                    {
                        "dict": "medicine",
                        "word": "奥利司他片",
                        "index": [
                            "奥利司他",
                            "奥利司他片"
                        ],
                        "data": {
                            "category": "西药",
                            "categoryId": 0,
                            "isControl": true,
                            "isHighRisk": true,
                            "medicineType": 0
                        }
                    }
                ]
            }
        },
        "格华止": {
            "keyword": "格华止",
            "positions": [
                {
                    "start": 8,
                    "end": 10
                }
            ],
            "dictwords": {
                "brand": [
                    {
                        "dict": "brand",
                        "word": "格华止",
                        "index": [
                            "格华止"
                        ],
                        "data": {
                            "category": "西药"
                        }
                    }
                ]
            }
        },
        "糖尿": {
            "keyword": "糖尿",
            "positions": [
                {
                    "start": 2,
                    "end": 3
                }
            ],
            "dictwords": {
                "symptom": [
                    {
                        "dict": "symptom",
                        "word": "糖尿",
                        "index": [
                            "糖尿"
                        ],
                        "data": {
                            "symptomId": 0
                        }
                    }
                ]
            }
        },
        "糖尿病": {
            "keyword": "糖尿病",
            "positions": [
                {
                    "start": 2,
                    "end": 4
                }
            ],
            "dictwords": {
                "disease": [
                    {
                        "dict": "disease",
                        "word": "糖尿病",
                        "index": [
                            "糖尿病"
                        ],
                        "data": {
                            "diseaseId": 0
                        }
                    }
                ]
            }
        }
    },
    "micros": 672
}
```

## API
| uri | method | content-type | parameter | usage |
|-------|-------|-------|-------|-------|
| / | GET | | | HTML API Document|
| /tag | GET/POST | application/json | statement | Annotate for statement |
| /put | POST | multipart/form-data | file | Upload single {dictname}.json file |
| /reload | GET/POST | | | Rebuild the index tree based on the dictionary file |
| /info | GET/POST | | | Display the dictionary information that the current index tree depends on | 

## Usage
### Windows
1. Download the Windows program archive from the release page.
2. Extract the archive.
3. (Optional)Edit your own dictionary entries in the data directory.
4. Double-click to run dict_tagging.exe.
5. Access in the browser: http://localhost:8080/tag?statement=your_statement

### Linux
1. Download the linux program archive from the release page.
2. Extract the archive.
3. (Optional)Edit your own dictionary entries in the data directory.
4. ./dict_tagging
5. Access in the browser: http://your-server-ip:8080/tag?statement=your_statement

### Docker
```sh
cd /path/to/dict_tagging_dir
docker build -t dict_tagging:1.0.1 .
docker run -d -p 8080:8080 dict_tagging:1.0.1
```

