# Architecture Notes

Tài liệu này mô tả đúng trạng thái hiện tại của template `go-modular-monolith`.

## 1. Mục tiêu của repo này

Repo này là starter template cho:
- `modular monolith`
- `clean architecture`
- `DDD-lite`

Nó không nhắm tới:
- microservices ngay từ đầu
- enterprise framework quá nặng
- quá nhiều abstraction khi chưa có nhu cầu thật

Ý tưởng chính là:
- một app
- một process
- nhiều module business chạy cùng nhau
- mỗi module có boundary rõ để về sau không rơi vào kiểu `service/` và `repository/` phình to

## 2. Luồng phụ thuộc

Luồng chuẩn trong từng module là:

```text
delivery/http
-> application
-> domain
<- infrastructure
```

Giải thích:
- `delivery/http` nhận request HTTP, gọi use case, trả response
- `application` chứa use case và orchestration logic
- `domain` chứa entity, repository contract, domain error
- `infrastructure` implement các contract mà `application` hoặc `domain` cần

Hướng phụ thuộc phải đi vào trong:
- `delivery` có thể biết `application`
- `application` có thể biết `domain`
- `domain` không biết `gin`, `gorm`, `zap`, HTTP hay DB

## 3. Composition Root

File [app.go](/Users/vothanh/Documents/Playground/project-example/internal/app/app.go) là `composition root`.

Đây là nơi:
- load shared resources
- tạo router
- gắn middleware
- quyết định module nào dùng `memory` repo, module nào dùng `postgres`
- đăng ký routes của toàn app

Nguyên tắc đang dùng trong repo:
- shared resource như Postgres connection được mở ở `app`
- concrete implementation được chọn ở `app`
- module chỉ nhận dependency dưới dạng contract phù hợp

Ví dụ với `order`:
- nếu `DATABASE_DSN` rỗng thì dùng `memory repository`
- nếu `DATABASE_DSN` có giá trị thì dùng `postgres repository`

Điều này giúp:
- `application` không cần biết đang chạy bằng memory hay Postgres
- đổi persistence không làm ảnh hưởng `delivery` và `domain`

## 4. Cấu trúc chính của repo

```text
cmd/server
cmd/migrate

internal/app

internal/platform/
  config/
  database/
  httpserver/
  logger/

internal/shared/
  collection/
  httpx/
  middleware/
  ptr/

internal/modules/
  order/
  customer/
  discount/
```

### `cmd/server`

File [main.go](/Users/vothanh/Documents/Playground/project-example/cmd/server/main.go):
- load config
- bootstrap logger
- tạo app
- xử lý shutdown signal

### `cmd/migrate`

File [main.go](/Users/vothanh/Documents/Playground/project-example/cmd/migrate/main.go):
- load config
- bootstrap logger
- chạy `migrate up` hoặc `migrate down`

### `internal/platform`

Phần hạ tầng dùng chung cho toàn app:

- [config.go](/Users/vothanh/Documents/Playground/project-example/internal/platform/config/config.go)
  - load env config
  - có generic helper `getEnv[T]`
- [postgres.go](/Users/vothanh/Documents/Playground/project-example/internal/platform/database/postgres.go)
  - mở kết nối Postgres
- [migrate.go](/Users/vothanh/Documents/Playground/project-example/internal/platform/database/migrate.go)
  - chạy migration SQL versioned
- [server.go](/Users/vothanh/Documents/Playground/project-example/internal/platform/httpserver/server.go)
  - bọc `http.Server`
- [logger.go](/Users/vothanh/Documents/Playground/project-example/internal/platform/logger/logger.go)
  - bootstrap `zap`
  - hỗ trợ `stdout|file|both`
  - hỗ trợ `json|console`
  - hỗ trợ log rotation bằng `lumberjack`
- [context.go](/Users/vothanh/Documents/Playground/project-example/internal/platform/logger/context.go)
  - gắn request-scoped logger vào `context.Context`
- [close.go](/Users/vothanh/Documents/Playground/project-example/internal/platform/logger/close.go)
  - đóng logger output ở entrypoint

### `internal/shared`

Phần dùng chung nhưng không thuộc business domain cụ thể:

- [response.go](/Users/vothanh/Documents/Playground/project-example/internal/shared/httpx/response.go)
  - success response envelope
- [errors.go](/Users/vothanh/Documents/Playground/project-example/internal/shared/httpx/errors.go)
  - error mapping cho HTTP
- [request_id.go](/Users/vothanh/Documents/Playground/project-example/internal/shared/httpx/request_id.go)
  - request id cho `gin.Context`
- [request_id.go](/Users/vothanh/Documents/Playground/project-example/internal/shared/middleware/request_id.go)
  - tạo hoặc lấy `X-Request-ID`
  - nhét request-scoped logger vào `context`
- [access_log.go](/Users/vothanh/Documents/Playground/project-example/internal/shared/middleware/access_log.go)
  - ghi log `[REQUEST] ...`
  - ghi log `[RESPONSE] ...`
  - có thể log request/response payload nếu bật config
- [index.go](/Users/vothanh/Documents/Playground/project-example/internal/shared/collection/index.go)
  - helper generic `IndexBy`
- [ptr.go](/Users/vothanh/Documents/Playground/project-example/internal/shared/ptr/ptr.go)
  - helper generic `ptr.Of`

## 5. Cấu trúc một module

Ví dụ module `order`:

```text
internal/modules/order/
  module.go
  domain/
  application/
  delivery/http/
  infrastructure/memory/
  infrastructure/postgres/
```

Ý nghĩa từng phần:

### `domain`

Các file:
- [order.go](/Users/vothanh/Documents/Playground/project-example/internal/modules/order/domain/order.go)
- [repository.go](/Users/vothanh/Documents/Playground/project-example/internal/modules/order/domain/repository.go)
- [errors.go](/Users/vothanh/Documents/Playground/project-example/internal/modules/order/domain/errors.go)

Chứa:
- entity
- repository contract
- domain error

Không chứa:
- HTTP
- GORM
- logger
- config

### `application`

Các file:
- [service.go](/Users/vothanh/Documents/Playground/project-example/internal/modules/order/application/service.go)
- [dto.go](/Users/vothanh/Documents/Playground/project-example/internal/modules/order/application/dto.go)
- [usecase.go](/Users/vothanh/Documents/Playground/project-example/internal/modules/order/application/usecase.go)

Chứa:
- use case
- input validation ở mức use case
- map domain entity sang DTO

`usecase.go` ở đây là inbound contract:
- `delivery/http` gọi vào `application` qua contract này
- implementation thật là `Service`

Với template hiện tại, nếu module A cần dùng một use case của module B thì ưu tiên cách đơn giản trước:
- inject thẳng `moduleB/application.UseCase` từ `internal/app`
- không gọi repo của module B
- không tạo thêm adapter hoặc outbound contract nếu chưa có pain point thật

Ví dụ hiện có trong repo:
- `order/application/service.go` dùng `discountapplication.UseCase`
- phần wiring nằm ở `internal/app/app.go`
- đây là dependency một chiều `order -> discount`
- nếu sau này xuất hiện gọi chéo hai chiều, khi đó mới nên tách orchestration ra ngoài

### `delivery/http`

File [handler.go](/Users/vothanh/Documents/Playground/project-example/internal/modules/order/delivery/http/handler.go):
- nhận request từ Gin
- gọi use case
- map lỗi sang HTTP response
- dùng `httpx.OK(...)` và `httpx.WriteError(...)`

### `infrastructure`

`order` hiện có 2 implementation:

- [repository.go](/Users/vothanh/Documents/Playground/project-example/internal/modules/order/infrastructure/memory/repository.go)
  - dùng để app chạy ngay khi chưa có DB
- [repository.go](/Users/vothanh/Documents/Playground/project-example/internal/modules/order/infrastructure/postgres/repository.go)
  - dùng Postgres thật
  - model và mapper tách riêng

### `module.go`

File [module.go](/Users/vothanh/Documents/Playground/project-example/internal/modules/order/module.go):
- nhận `Dependencies`
- tự lắp `repo -> usecase -> handler`
- expose `RegisterRoutes(...)`

Điểm quan trọng:
- module không tự mở DB
- module không tự đọc env
- module nhận dependency đã được quyết định từ `app`

## 6. Request Flow Hiện Tại

Luồng request hiện tại là:

```text
HTTP request
-> RequestID middleware
-> AccessLog middleware ([REQUEST] ...)
-> handler
-> application/service
-> repository
-> handler
-> AccessLog middleware ([RESPONSE] ...)
-> HTTP response
```

Chi tiết:

1. `RequestID` middleware:
- lấy `X-Request-ID` từ request nếu có
- nếu chưa có thì tự sinh
- set lại vào response header
- tạo request-scoped logger có field `request_id`
- nhét logger đó vào `context.Context`

2. `AccessLog` middleware:
- log một dòng mở đầu `[REQUEST] <path>`
- sau khi xử lý xong, log một dòng `[RESPONSE] <path>`
- response log có `status`, `latency_ms`, `response_size_bytes`
- nếu bật config, request/response payload cũng được log

3. `service` có thể log bằng:
- `logger.FromContext(ctx).Debug(...)`

Nhờ vậy log trong `service` cũng tự có `request_id`.

## 7. Nếu Sau Này Có Cross-Module Call

Template hiện tại không cài sẵn ví dụ module A gọi module B trong runtime, vì với starter template như vậy thường hơi over-engineering.

Cách khuyên dùng là:
- bắt đầu đơn giản
- wiring ở `internal/app`
- chỉ thêm outbound contract khi thật sự có áp lực về coupling hoặc cần thu hẹp surface của module kia

Rule thực dụng:
- đừng gọi repo của module khác
- đừng query thẳng bảng của module khác
- nếu cần logic của module kia, đi qua `application` contract của nó
- nếu flow bắt đầu lớn, lúc đó mới cân nhắc `application/dependencies.go`

## 8. Logging Hiện Tại

Repo hiện dùng `zap`.

Các đặc điểm:
- structured log
- JSON line mặc định
- có thể ghi ra file
- có thể rotate file
- có request-scoped logger

Các config hiện có trong [config.go](/Users/vothanh/Documents/Playground/project-example/internal/platform/config/config.go):
- `LOG_OUTPUT`
- `LOG_FILE_PATH`
- `LOG_LEVEL`
- `LOG_FORMAT`
- `LOG_INCLUDE_REQUEST_BODY`
- `LOG_INCLUDE_RESPONSE_BODY`
- `LOG_BODY_MAX_BYTES`
- `LOG_MAX_SIZE_MB`
- `LOG_MAX_BACKUPS`
- `LOG_MAX_AGE_DAYS`
- `LOG_COMPRESS`

Lưu ý:
- log body đang là `opt-in`
- nên chỉ bật khi debug hoặc môi trường kiểm soát tốt
- không nên bật mặc định ở production nếu request có dữ liệu nhạy cảm

## 9. Persistence Hiện Tại

`customer`:
- đang dùng memory repository

`order`:
- dùng memory repository nếu chưa có `DATABASE_DSN`
- dùng Postgres repository nếu có `DATABASE_DSN`

Schema của Postgres được quản lý bằng migration SQL versioned:
- [000001_create_orders_table.up.sql](/Users/vothanh/Documents/Playground/project-example/migrations/000001_create_orders_table.up.sql)
- [000001_create_orders_table.down.sql](/Users/vothanh/Documents/Playground/project-example/migrations/000001_create_orders_table.down.sql)

Không dùng `AutoMigrate`.

Điểm này quan trọng vì nó giữ:
- schema có version rõ ràng
- app boot không tự âm thầm sửa DB

## 10. Vì sao repo có thêm vài helper generic nhỏ?

Repo hiện có một số helper generic nhỏ như:
- `getEnv[T]`
- `collection.IndexBy`
- `ptr.Of`

Mục tiêu không phải là “generic hóa mọi thứ”.
Mục tiêu là:
- cắt lặp ở những pattern rất cơ học
- giữ call site ngắn, rõ, dễ đọc

Nếu một helper làm code business khó hiểu hơn, thì không nên thêm.

## 11. Những gì repo chưa cố làm

Repo này chưa cố thêm:
- auth framework
- event bus
- outbox
- CQRS tách sâu
- transaction abstraction lớn
- distributed tracing stack

Lý do là để tránh over-engineering ở giai đoạn làm starter template.

## 12. Khi nào dùng template này?

Template này hợp với:
- dự án monolith sống lâu
- có nhiều domain/module business
- muốn bắt đầu nhanh nhưng vẫn giữ boundary sạch

Template này không tối ưu cho:
- tool rất nhỏ
- CRUD rất đơn giản
- spike ngắn ngày

Nếu project đủ nhỏ, layer architecture đơn giản vẫn có thể hợp lý hơn.
Nếu project sẽ sống lâu, modular monolith như repo này thường là điểm bắt đầu tốt hơn.
