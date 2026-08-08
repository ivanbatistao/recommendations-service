# Events
Here we will define exactly what will travel through Kinesis.

## Event Types
- `product_viewed`
- `search_performed`
- `product_purchased`
- `product_added_to_cart`

## Event Contract
```
{
  "eventId": "...",
  "eventType": "product_viewed",
  "userId": "123",
  "productId": "P101",
  "productCategory": "electronics",
  "productBrand": "Samsung",
  "metadata": {
    "device": "mobile",
    "country": "CO"
  },
  "occurredAt": "..."
}
```

This contract will be stable and versionable.

