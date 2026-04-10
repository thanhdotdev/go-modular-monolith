# go-modular-monolith

Starter template Go theo hướng `modular monolith + clean architecture + DDD-lite`.

Template này cố ý đi theo hướng thực dụng:
- Có boundary rõ giữa `delivery`, `application`, `domain`, `infrastructure`
- Không dùng quá nhiều abstraction ngay từ đầu
- Chạy được ngay mà chưa cần database thật
- Có thể bật Postgres thật cho module `order` khi cần
- Có request id, access log, response envelope và error mapping chung cho HTTP layer
- Có thể ghi log ra file

## Quick Start

```bash
go run ./cmd/server
```

Server mặc định chạy ở `http://localhost:8080`.

Endpoints mẫu:
- `GET /healthz`
- `GET /api/v1/orders/ord-001`
- `GET /api/v1/customers/cus-001`
- `GET /api/v1/discounts/WELCOME10`

Response thành công hiện theo envelope:

```json
{
  "data": {
    "id": "ord-001",
    "discountCode": "WELCOME10",
    "customerName": "Alice",
    "status": "pending",
    "totalAmount": 125000,
    "discount": {
      "code": "WELCOME10",
      "type": "fixed",
      "value": 10000,
      "active": true
    }
  },
  "meta": {
    "requestId": "..."
  }
}
```

Response lỗi hiện theo envelope:

```json
{
  "error": {
    "code": "order_not_found",
    "message": "order not found"
  },
  "meta": {
    "requestId": "..."
  }
}
```

Nếu `DATABASE_DSN` rỗng:
- `order` dùng memory repository

Nếu `DATABASE_DSN` có giá trị:
- `order` dùng Postgres repository
- cần chạy migration trước khi dùng Postgres

## Cấu trúc chính

```text
cmd/server
internal/app
internal/platform
internal/modules/order
internal/modules/customer
internal/modules/discount
docs/architecture.md
docs/module-guide.md
```

Giải thích nhanh:
- `cmd/server`: điểm vào của ứng dụng
- `internal/app`: composition root, lắp các module và khởi động server
- `internal/platform`: phần hạ tầng dùng chung như config, logger, http server
- `internal/shared`: phần dùng chung ở mức HTTP như middleware và response helper
- `internal/modules/order`: module business mẫu cho order
- `internal/modules/customer`: module business mẫu thứ hai để thấy boundary giữa các module
- `internal/modules/discount`: module business mẫu cho discount và ví dụ module-to-module call một chiều

Chi tiết hơn xem ở [docs/architecture.md](docs/architecture.md).
Hướng dẫn thêm module mới xem ở [docs/module-guide.md](docs/module-guide.md).

## Lệnh hữu ích

```bash
make run
make test
make test-integration
make fmt
make fmt-check
make vet
make lint
make new-module name=invoice
make migrate-up
make migrate-down
```

## Logging ra file

Logger hiện dùng `zap` và hỗ trợ rotation file qua `lumberjack`.
Mặc định log được ghi theo dạng JSON line để dễ đẩy vào log monitor hoặc log shipper.

Mặc định app ghi log ra cả terminal lẫn file:

```bash
LOG_OUTPUT=both
LOG_FILE_PATH=logs/app.log
```

Các chế độ hỗ trợ:
- `LOG_OUTPUT=stdout`: chỉ ghi ra terminal
- `LOG_OUTPUT=file`: chỉ ghi ra file
- `LOG_OUTPUT=both`: ghi ra cả terminal và file

Một số cấu hình hữu ích:
- `LOG_SERVICE_NAME=project-example`
- `LOG_LEVEL=debug|info|warn|error`
- `LOG_FORMAT=json|console`
- `LOG_INCLUDE_REQUEST_BODY=true|false`
- `LOG_INCLUDE_RESPONSE_BODY=true|false`
- `LOG_BODY_MAX_BYTES=4096`
- `LOG_MAX_SIZE_MB=100`
- `LOG_MAX_BACKUPS=3`
- `LOG_MAX_AGE_DAYS=7`
- `LOG_COMPRESS=false`

Ví dụ:

```bash
go run ./cmd/server
tail -f logs/app.log
```

Ví dụ log JSON:

```json
{"L":"INFO","timestamp":"2026-04-09T16:42:19+07:00","C":"server/main.go:25","M":"starting http server","service":"project-example","addr":":8080"}
```

Nếu cần log payload để debug request/response:

```bash
LOG_INCLUDE_REQUEST_BODY=true
LOG_INCLUDE_RESPONSE_BODY=true
LOG_BODY_MAX_BYTES=4096
```

Body chỉ được log cho content type text/json/form và sẽ bị cắt theo `LOG_BODY_MAX_BYTES`.

## Bật Postgres cho order

Ví dụ DSN:

```bash
export DATABASE_DSN="host=localhost user=postgres password=postgres dbname=project_example port=5432 sslmode=disable"
go run ./cmd/server
```

Khi đó module `order` sẽ đọc từ Postgres thay vì memory.

Nếu bạn muốn dựng Postgres local nhanh:

```bash
make up
export DATABASE_DSN="host=localhost user=postgres password=postgres dbname=project_example port=5432 sslmode=disable"
make migrate-up
go run ./cmd/server
```

Khi xong:

```bash
make down
```

## Module gọi module khác theo cách nào?

Template này mặc định ưu tiên cách đơn giản:
- module A không gọi repo của module B
- module A có thể nhận thẳng `application.UseCase` của module B từ `internal/app`
- chỉ tách thêm outbound contract hoặc adapter khi pain point thật sự xuất hiện

Ví dụ đang có trong repo:
- `order` gọi `discountapplication.UseCase`
- việc wiring nằm ở [app.go](internal/app/app.go)
- `order` chỉ dùng `discount` ở `application/service.go`, không chạm vào repo hay hạ tầng của `discount`

## Vì sao chưa có Postgres/Redis thật?

Để tránh over-engineering ở giai đoạn học khung:
- module `order` có thể chạy bằng memory hoặc Postgres
- module `customer` và `discount` hiện dùng memory repository
- `order` là ví dụ đầu tiên cho việc thay memory bằng persistence thật
- lúc đó `application` và `domain` gần như không phải đổi

## Migrations

Migration files nằm ở [migrations](migrations).

Chạy migration:

```bash
export DATABASE_DSN="host=localhost user=postgres password=postgres dbname=project_example port=5432 sslmode=disable"
make migrate-up
```

Rollback toàn bộ migration:

```bash
make migrate-down
```

Migration hiện tại tạo bảng `orders` và seed sẵn dữ liệu demo để endpoint mẫu hoạt động ngay.

## CI và lint

Repo đã có workflow CI ở [.github/workflows/ci.yml](.github/workflows/ci.yml).

CI hiện chạy:
- `make fmt-check`
- `go vet ./...`
- `make test`
- `make lint`
- `make test-integration`

Lint dùng [golangci-lint](.golangci.yml) với cấu hình nhẹ:
- giữ standard linters
- thêm `depguard` để khóa boundary của `domain` và `application`
- thêm `errorlint`, `misspell`, `nolintlint`
