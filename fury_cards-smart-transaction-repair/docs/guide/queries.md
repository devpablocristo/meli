## Queries

### Result Validations [(Document Search)](https://web.furycloud.io/cards-smart-transaction-repair/ds/search/787468)

> Query example
```json
{
  "query": {
    "and": [
      {
        "eq": {
          "field": "site_id",
          "value": "MLM"
        }
      },
      {
        "date_range": {
          "field": "created_at",
          "format": "yyyy/MM/dd",
          "time_zone": "UTC",
          "gte": "2022/12/28",
          "lte": "2023/12/28"
        }
      }
    ]
  },
  "type": "query_and_fetch",
  "size": 1000,
  "sort": [
    {
      "field": "created_at",
      "field_type": "date",
      "order": "desc"
    }
  ]
}
```

<br/>

### Repairs [(Document Search)](https://web.furycloud.io/cards-smart-transaction-repair/ds/search/787308)

> Query example
```json
{
  "query": {
    "and": [
      {
        "eq": {
          "field": "site_id",
          "value": "MLM"
        }
      },
      {
        "date_range": {
          "field": "created_at",
          "format": "yyyy/MM/dd",
          "time_zone": "UTC",
          "gte": "2022/12/28",
          "lte": "2023/12/28"
        }
      }
    ]
  },
  "type": "query_and_fetch",
  "size": 1000,
  "sort": [
    {
      "field": "created_at",
      "field_type": "date",
      "order": "desc"
    }
  ]
}
```
