# TelemetryGuard

TelemetryGuard 是工业设备遥测接入与告警管理平台：设备点位数据按序列号接入去重，
按已发布规则快照评估阈值，越限产生告警并走打开、确认、恢复、关闭的生命周期；
长时间未处理按 SLA 逐级升级并通过短信/企微渠道通知，支持维护窗口抑制、
渠道故障切换和全量操作审计。

## 功能

- 设备/点位注册与启停管理
- 遥测接入、序列号去重与批量评估
- 规则版本发布、快照评估与回滚
- 告警生命周期（打开/确认/恢复/关闭）
- SLA 升级计时与通知渠道故障切换
- 维护窗口抑制与操作审计
- 控制台页面与 HTTP API

## 运行

```bash
go build -mod=vendor -o telemetryguard ./cmd/server
./telemetryguard -addr :8080 -data ./data
```

打开 http://localhost:8080/ 查看总览，/console/alarms、/console/rules、
/console/points、/console/audit 分别查看告警、规则、点位与审计页面。

## 测试

```bash
go test -mod=vendor ./...
```
