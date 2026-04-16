# eight-id-go

Generator sinh ID base62 dài đúng `8` ký tự.

Mỗi ID:
- Dài đúng `8` ký tự
- Chỉ chứa `0-9`, `A-Z`, `a-z`
- Có phân biệt chữ hoa chữ thường
- Sortable theo thứ tự chuỗi

## Thiết kế hiện tại

Cấu trúc ID:
- `6` ký tự đầu là thời gian theo đơn vị `10ms`
- `2` ký tự cuối là sequence trong cùng một `10ms`

Trade-off:
- Tối đa `62^2 = 3,844` ID trong mỗi `10ms`
- Tương đương khoảng `384,400 ID/giây` trên mỗi process
- Tuổi thọ khoảng `18 năm` kể từ mốc `2025-01-01T00:00:00Z`

Đây là điểm cân bằng tốt giữa:
- Sortability
- Kích thước cố định `8` ký tự
- Throughput cao
- Vùng thời gian dùng đủ dài

## Usage

```go
package main

import (
	"fmt"
	"time"

	eightid "github.com/chi07/eight-id-go"
)

func main() {
	id := eightid.New()
	fmt.Println(id)

	// Deterministic ID theo mốc thời gian cho test hoặc backfill
	fixed := eightid.NewWithTime(time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC))
	fmt.Println(fixed)
}
```

## API

```go
id := eightid.New()
ok := eightid.IsValid(id)
fixed := eightid.NewWithTime(time.Now())
```
