> [!warning]
> **THIS REPOSITORY IS NO LONGER MAINTAINED.**

# reCAPTCHA Server

A lightweight server implementation for handling reCAPTCHA verification in your applications.

## Features

- Verify reCAPTCHA tokens securely.
- Easy integration with backend services.
- Configurable for different reCAPTCHA versions.

## Development Setup

1. Clone the repository:
    ```bash
    git clone https://github.com/yourusername/recaptcha-server.git
    cd recaptcha-server
    ```

2. Install dependencies:
    ```bash
    make tidy
    ```

3. Configure environment variables:
    Create a `.env` file and set the following:
    ```
    RECAPTCHA_SECRET_KEY="" # Required
    RECAPTCHA_FRONTEND="" # Optional, but HIGHLY recommended for production  use,can be any valid URL (e.g., http://localhost:3000), defaults to "*"
    RECAPTCHA_TIMEZONE="" # Optional, can be any valid timezone string (See https://en.wikipedia.org/wiki/List_of_tz_database_time_zones)
    RECAPTCHA_PORT="" # Optional, defaults to 8080
    LOG_LEVEL="" # Optional, defaults to "info" ("debug", "info", "warn", "error", "fatal")
    ```

> [!warning]
> You should _only_ use .env files in development. In production, use environment variables directly.

### Running the Server

Start the server:
```bash
make dev
```

Send a POST request to verify a token:
```bash
POST /verify
Content-Type: application/json

{
  "token": "your-recaptcha-token"
}
```

## Production Setup

For production, you can build the server binary or use the Docker image. There is a publicly available Docker image on GitHub Container Registry. However, it is recommended to build the image yourself to ensure security and integrity. You can build the Docker image using the following command:

```bash
make build
```

### Required Environment Variables

- `RECAPTCHA_SECRET_KEY`: Your reCAPTCHA secret key.
- `RECAPTCHA_FRONTEND`: Optional, but highly recommended for production use. Can be any valid URL (e.g., http://localhost:3000), defaults to "*".

#### Optional Environment Variables

- `RECAPTCHA_TIMEZONE`: Can be any valid timezone string (See [List of tz database time zones](https://en.wikipedia.org/wiki/List_of_tz_database_time_zones)).
- `RECAPTCHA_PORT`: Defaults to `8080`.
- `LOG_LEVEL`: Defaults to `info`. Can be set to `debug`, `info`, `warn`, `error`, or `fatal`.


## Request Format

The only endpoint available is `/verify`, which accepts a POST request with the following JSON body:

```json
{
  "token": "your-recaptcha-token"
}
```

## Response
- **200 OK**: Verification successful.
- **400 Bad Request**: Invalid or missing token. (See response body for details)

## Contributing

Contributions are welcome! Please fork the repository and submit a pull request.

## License

This project is licensed under the [MIT License](LICENSE).
