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
- `application`: use case, input port, DTO, app-level validation
- `delivery/http`: adapter HTTP
- `infrastructure`: repository implementation
- `module.go`: chỗ lắp repo -> usecase -> handler

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
1. Đăng ký module trong [app.go](/Users/vothanh/Documents/Playground/project-example/internal/app/app.go)
2. Viết entity và use case thật
3. Thay memory repository bằng persistence thật khi cần

## 3. Quy tắc làm module

- Không để `domain` import `gin`, `gorm`, `redis`
- Không gọi repository của module khác trực tiếp
- Nếu module A cần logic của module B, đi qua `application` contract của B
- Chỉ thêm file mới khi module bắt đầu đủ phức tạp để cần nó
