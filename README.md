# System Information Web App

This is a Go-based web application that provides detailed system information such as OS, CPU, Memory, Disk, and GPU details. It uses `gopsutil` to gather the system metrics and displays them in an HTML template.

## Features

- Shows detailed system information (OS, CPU, Memory, Disk, and GPU)
- Supports multiple operating systems: Windows, Linux, and macOS
- Displays real-time data on a web page via a Go web server

## Technologies Used

1. **Go (Golang)**: The primary language used to build the application.
2. **gopsutil**: A Go library used for retrieving system information such as CPU, Memory, and Disk usage.
3. **HTML Templates**: Used for rendering system information in a user-friendly format on a web page.
4. **net/http**: Go's standard library package used to handle HTTP requests and serve the web pages.

## How to Run

### 2. Install dependencies

Make sure you have Go installed on your system. If not, you can download it from [Go's official website](https://golang.org/dl/).

Also, you'll need to install the `gopsutil` package:

```bash
go get github.com/shirou/gopsutil/v3
