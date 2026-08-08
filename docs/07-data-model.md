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
ProductViewed
```

## Recommendations Table
`PK = UserID`

## Document Structure
```
{
  "userId": "15",
  "products": [
    "P10",
    "P50",
    "P80"
  ],
  "updatedAt": "..."
}
```