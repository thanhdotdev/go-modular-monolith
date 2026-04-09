# Architecture Notes

Tài liệu này giải thích mục đích của từng tầng trong template modular monolith hiện tại.

## 1. Luồng phụ thuộc

Luồng chuẩn là:

```text
HTTP request
-> delivery
-> application
-> domain
-> infrastructure
```

Nguyên tắc:
- `delivery` nhận request và trả response
- `application` chứa use case
- `domain` chứa business model và business rule
- `infrastructure` giao tiếp với thế giới bên ngoài như DB, cache, queue

Monolith ở đây nghĩa là:
- một app
- một process
- nhiều module business cùng chạy bên trong

Tức là đây không phải microservices, mà là modular monolith.

## 2. Use case nằm ở đâu?

Trong template này, `usecase` được đặt ở tầng `application`.

Ví dụ:
- [service.go](/Users/vothanh/Documents/Playground/project-example/internal/modules/order/application/service.go)
- [port_in.go](/Users/vothanh/Documents/Playground/project-example/internal/modules/order/application/port_in.go)

Lý do:
- delivery chỉ nên gọi use case
- domain không nên biết HTTP hay request/response
- infrastructure không nên chứa business flow

## 3. Vai trò từng phần

### `cmd/server/main.go`
- Điểm vào của app
- Load config, tạo logger, chạy app

### `internal/app/app.go`
- Composition root
- Nơi lắp các thành phần lại với nhau
- Đăng ký route, khởi động HTTP server, xử lý shutdown
- Đây cũng là nơi đăng ký các module của monolith
- Shared resource như DB connection/pool nên được mở ở đây một lần rồi inject xuống các module cần dùng

### `internal/platform/config`
- Đọc cấu hình từ env
- Không chứa business logic

### `internal/platform/logger`
- Tạo logger dùng chung

### `internal/platform/httpserver`
- Gói `http.Server` lại để app dễ quản lý

### `internal/platform/database`
- Bootstrap Postgres connection
- Chạy migrations
- Không chứa business logic
- Không nên để từng module tự mở kết nối DB riêng nếu cùng dùng một Postgres

### `internal/shared/middleware`
- Chứa middleware dùng chung như request id và access log
- Đây là technical shared code, không chứa business rule

### `internal/shared/httpx`
- Chứa response envelope và error mapping chung cho HTTP layer
- Mục tiêu là cắt lặp `c.JSON(...)` và chuẩn hóa response shape

### `internal/modules/order/domain`
- Entity `Order`
- `Repository` là contract mà application cần
- `ErrOrderNotFound` là business error của module

### `internal/modules/order/application`
- Chứa use case
- Nhận input, gọi repository, map sang DTO trả về cho delivery
- Đây là nơi điều phối flow của module

### `internal/modules/order/delivery/http`
- Nhận request HTTP
- Gọi use case
- Map lỗi sang status code phù hợp
- Dùng shared response/error helper thay vì tự build JSON mỗi handler

### `internal/modules/order/infrastructure/memory`
- Repository implementation cho ví dụ hiện tại
- Được dùng để app chạy được ngay mà không cần DB thật

### `internal/modules/order/infrastructure/postgres`
- Repository implementation dùng Postgres thật
- Chứa `model`, `mapper`, `repository`
- Đây là ví dụ rõ nhất cho việc domain không biết gì về GORM/Postgres

### `internal/modules/order/module.go`
- Entry point của module `order`
- Lắp `repository -> usecase -> handler`
- Expose `RegisterRoutes` cho app

### `internal/modules/customer`
- Module mẫu thứ hai
- Có cấu trúc giống `order`
- Dùng để chứng minh module boundary trong cùng một monolith

## 4. Tại sao chưa tách nhiều file hơn?

Để tránh over-engineering.

Hiện tại tôi chỉ tách file khi nó giúp nhìn rõ vai trò:
- `domain/order.go`: entity
- `application/service.go`: use case
- `delivery/http/handler.go`: adapter HTTP

Các phần như `commands.go`, `queries.go`, `mapper.go`, `response.go` có thể thêm sau khi module bắt đầu lớn hơn.

## 5. Khi nào nên thêm abstraction mới?

Chỉ thêm khi có nhu cầu thật:
- Thêm `postgres` repository khi app cần DB thật
- Thêm `shared/response` khi có nhiều module cần format response giống nhau
- Thêm `commands.go` và `queries.go` khi use case bắt đầu nhiều
- Thêm abstraction cross-module khi bạn thực sự có nhu cầu trao đổi giữa các module

Hiện tại template đã thêm `internal/shared/httpx` vì đã có hơn một module HTTP và response/error bắt đầu lặp lại.

## 6. Thay infrastructure mà không đổi core

Module `order` hiện là ví dụ cho flow này:

```text
memory repository -> postgres repository
```

Khi đổi persistence:
- `domain` không đổi
- `application` không đổi
- `delivery/http` không đổi
- chỉ đổi implementation ở `infrastructure` và cách lắp ở `app`

Schema của Postgres không còn tạo bằng `AutoMigrate`.
Thay vào đó repo dùng migration SQL versioned trong thư mục `migrations/`.

## 7. Khi nào dùng template này?

Template này hợp với:
- dự án monolith sống lâu
- có nhiều domain/module business
- muốn bắt đầu nhanh nhưng vẫn giữ codebase sạch

Template này không tối ưu cho:
- tool nhỏ một vài CRUD đơn giản
- spike ngắn ngày
- project chỉ có 1-2 màn hình và không có business rule đáng kể
