# URL Shortener Service

A lightweight, self-hosted URL shortener service built with Go. This application allows users to create and manage short, easily shareable URLs that redirect to longer URLs. The service also includes periodic persistence to a CSV file for data durability.

## Features

- **Static file server**: Serves a web frontend embedded within the binary.
- **Lightweight and efficient**: Designed to be fast, simple, and easy to deploy.

## Prerequisites

- [Go 1.18+](https://golang.org/dl/)
- A basic understanding of how to run Go applications.

## Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/gabrielmajosi/url-shortener-go.git
   cd url-shortener-go
   ```
2. Build the frontend with Vite:
   ```bash
   cd frontend
   npm install
   npm run build
   cp -r dist/* ../web/
   cd ..
   ```
3. Build and run the Go project:
   ```bash
   go build .
   ./url-shortener-go
   ```
   
The server will start on port `8080` by default.

## Configuration
Server Port: Update the serverAddr constant to change the listening port.
Storage File: Modify the storeFile constant to use a different file for persistence.

## License

This project is licensed under the [GNU General Public License v3.0](https://www.gnu.org/licenses/gpl-3.0.html).

By using, modifying, or distributing this software, you agree to comply with the terms of the GPL. This license ensures that:

- You are free to use, modify, and distribute the software for any purpose.
- You must share any changes or derivative works under the same license.
- You must include the original license notice and copyright in any copies of the software.

For the full license text, see the [LICENSE](LICENSE) file in this repository.
