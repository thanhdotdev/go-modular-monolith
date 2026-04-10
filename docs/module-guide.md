# Module Guide

Tài liệu này mô tả cách thêm module mới trong modular monolith hiện tại.

## 1. Cấu trúc chuẩn của một module

```text
internal/modules/<module>/
  module.go
  domain/
  application/
  delivery/http/
  infrastructure/memory/
```

Mỗi tầng có vai trò riêng:
- `domain`: entity, repository contract, business error
- `application`: use case, inbound contract, DTO, app-level validation
- `delivery/http`: adapter HTTP
- `infrastructure`: repository implementation
- `module.go`: chỗ module tự lắp dependencies -> usecase -> handler

Naming đang dùng trong `application`:
- `usecase.go`: inbound contract mà `delivery` sẽ gọi
- `service.go`: implementation của use case
- `dependencies.go`: outbound contracts nếu module bắt đầu cần thu hẹp dependency ra ngoài; không bắt buộc phải tạo từ đầu

## 2. Sinh module mới

```bash
./scripts/new_module.sh invoice
```

Hoặc:

```bash
make new-module name=invoice
```

Script sẽ tạo skeleton tối thiểu, kèm test mẫu cho `application` và `delivery/http`.
Sau đó bạn cần:
1. Đăng ký module trong [app.go](../internal/app/app.go)
2. Viết entity và use case thật
3. Thay memory repository bằng persistence thật khi cần

## 3. Quy tắc làm module

- Không để `domain` import `gin`, `gorm`, `redis`
- Không gọi repository của module khác trực tiếp
- Nếu module A cần logic của module B, bắt đầu đơn giản bằng cách inject `application` contract của B từ `internal/app`
- Chỉ thêm `application/dependencies.go` khi direct dependency bắt đầu gây khó đọc hoặc cần thu hẹp contract
- Wiring giữa module A và B nên nằm ở `internal/app`, không nhét vào `domain`
- Chỉ thêm file mới khi module bắt đầu đủ phức tạp để cần nó
