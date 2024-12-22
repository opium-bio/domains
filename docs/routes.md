```bash
curl -X POST http://localhost:8000/v1/add \
     -H "Content-Type: application/json" \
     -d '{"domain": "test.com"}'
```

Response

```bash
STATUS 200
{
  "success": true,
  "domain": "test.com"
}
```

yo can somebody write better docs for this please 🙏