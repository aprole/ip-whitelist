# IP Whitelist

A simple IP whitelist service that checks if an IP address is from a whitelisted country. It supports HTTP requests and gRPC requests, and uses an IP geolocation database from MaxMind to determine the location of an IP. The database is updated on a weekly basis.

## Installation

To run the server: ```docker compose up```

## Usage

HTTP:
```bash
curl -X POST http://localhost:8080/api/check-ip -d '{"ip":"75.127.6.164", "allowedCountries":["US", "RS","IT","JP"]}'
```

gRPC:
```bash
grpcurl -plaintext -d '{"ip":"75.127.6.164", "allowedCountries":["RS", "IT", "JP"]}' localhost:50051 IPWhitelist.CheckIP
```

## Contributing

Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

## License

[MIT](https://choosealicense.com/licenses/mit/)




