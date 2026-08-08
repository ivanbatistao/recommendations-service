# Data Model
Here we decide how the information will be stored.

## InteractionEvents Table
```
PK = UserID
SK = Timestamp
```

## Example
```
USER#15
2026-08-06T09:15
{
  "eventId": "...",
  "eventType": "product_viewed",
  "productId": "P101",
  "productCategory": "electronics",
  "productBrand": "Samsung",
  "metadata": {
    "device": "mobile",
    "country": "CO"
  },
  "occurredAt": "2026-08-06T09:15:00Z"
}
```

## Recommendations Table
`PK = UserID`

## Document Structure
```
{
  "userId": "15",
  "products": [
    {
      "productId": "P10",
      "score": 0.95
    },
    {
      "productId": "P50",
      "score": 0.82
    },
    {
      "productId": "P80",
      "score": 0.78
    }
  ],
  "updatedAt": "..."
}
```