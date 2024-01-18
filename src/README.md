# Exchange Rates API Documentation

## Convert Currency

Endpoint for converting currency rates.

### Endpoint


### Parameters

| Parameter | Type    | Description                        |
|-----------|---------|------------------------------------|
| `from`    | string  | The currency code to convert from. |
| `to`      | string  | The currency code to convert to.   |

### Example

```http
GET http://my-rates.com/convert?from=cad&to=inr
```
# Response
- Status Code: 200 OK
- Content Type: application/json

| Field | Type   | Description                      |
|-------|--------|----------------------------------|
| date  | string | The date of conversion           |
| inr   | number | The converted currency rate(INR) |

# Error Responses
- Status Code: 400 Bad request
```json
{
  "error": "Bad Request",
  "message": "Missing required parameter: 'from' or 'to'"
}

```
- Status Code: 404 Not Found
```json
{
  "error": "Not Found",
  "message": "Conversion rate not found for the specified currency pair"
}

```

- Status Code: 500 Internal Server Error
```json
{
  "error": "Internal Server Error",
  "message": "An unexpected error occurred while processing the request"
}

```

# Cache
-The service maintains an in-memory LRU cache to cache the conversion rates with an expiration of 30 minutes.
- The cached value will be removed upon expiration or if it is least recently used.
- The key of cached value is fromCurrency-toCurrency and value is the conversion rate and the expiration.