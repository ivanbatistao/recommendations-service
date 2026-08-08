# API Design 

## Get Recommendations
GET `/recommendations/{userId}`

### Response
```json
{
  "userId": "123",
  "recommendations": [
    {
      "productId": "P20",
      "score": 0.95
    },
    {
      "productId": "P40",
      "score": 0.82
    }
  ]
}
```

## Health Check
GET `/health`


## Metrics
GET `/metrics`