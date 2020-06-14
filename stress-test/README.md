## 負荷試験用ツール

### 準備

```
# ビルド
$ go build

# .envを用意
$ cp .env.example .env

# .envを編集(適当な方法で.env内の指示に従って編集してください)
```

### オプション
```
$ ./stress-test -h
Usage of ./vegeta:
  -duration int
    	how many seconds (default -1)
  -random-page-access
    	random page access mode
  -random-page-update
    	random page update mode
  -rate int
    	request per second (default -1)
```

### 使用例

ランダムなページアクセスを1秒間に10回, 5秒間行う

```
$ ./stress-test --rate 10 --duration 5 --random-page-access
## Vegeta Metrics

Requests      [total, rate, throughput]  50, 10.21, 10.15
Duration      [total, attack, wait]      4.927114394s, 4.895700049s, 31.414345ms
Latencies     [mean, 50, 95, 99, max]    36.513636ms, 36.058349ms, 42.011943ms, 59.931445ms, 59.931445ms
Bytes In      [total, mean]              2464193, 49283.86
Bytes Out     [total, mean]              0, 0.00
Success       [ratio]                    100.00%
Status Codes  [code:count]               200:50
Error Set:
```

ランダムなページ更新を1秒間に2回, 5秒間行う

```
$ ./stress-test --rate 2 --duration 5 --random-page-update
## Vegeta Metrics

Requests      [total, rate, throughput]  10, 2.22, 2.19
Duration      [total, attack, wait]      4.559234748s, 4.501258855s, 57.975893ms
Latencies     [mean, 50, 95, 99, max]    61.223259ms, 62.01974ms, 66.027239ms, 66.027239ms, 66.027239ms
Bytes In      [total, mean]              214135, 21413.50
Bytes Out     [total, mean]              200915, 20091.50
Success       [ratio]                    100.00%
Status Codes  [code:count]               200:10
Error Set:
```
