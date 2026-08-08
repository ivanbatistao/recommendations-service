```
recommendation-service/
│
├── cmd/
│   ├── api/
│   └── event-generator/
│
├── internal/
│   ├── domain/
│   │   ├── recommendation/
│   │   ├── product/
│   │   ├── user/
│   │   └── event/
│   │
│   ├── application/
│   │   ├── commands/
│   │   ├── queries/
│   │   └── dto/
│   │
│   └── infrastructure/
│       ├── http/
│       │   └── gin/
│       ├── dynamodb/
│       ├── kinesis/
│       └── logger/
│
├── configs/
├── deploy/
├── scripts/
├── tests/
├── benchmarks/
├── docs/
└── go.mod
```