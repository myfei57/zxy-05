package console

import (
	"html/template"
	"net/http"
)

var templates = template.Must(template.New("pages").Parse(`
{{define "layout"}}<!DOCTYPE html>
<html><head><meta charset="utf-8"><title>{{.Title}}</title>
<style>
body{font-family:sans-serif;margin:24px;color:#222}
nav a{margin-right:12px}
table{border-collapse:collapse;margin-top:12px}
td,th{border:1px solid #ccc;padding:6px 10px;font-size:13px}
</style></head><body>
<nav><a href="/">总览</a><a href="/console/alarms">告警</a><a href="/console/rules">规则</a><a href="/console/points">点位</a><a href="/console/audit">审计</a></nav>
<h2>{{.Title}}</h2>{{.Body}}</body></html>{{end}}
{{define "index"}}{{template "layout" .}}{{end}}
{{define "alarms"}}{{template "layout" .}}{{end}}
{{define "rules"}}{{template "layout" .}}{{end}}
{{define "points"}}{{template "layout" .}}{{end}}
{{define "audit"}}{{template "layout" .}}{{end}}
`))

type pageData struct {
	Title string
	Body  template.HTML
}

func (a *API) render(w http.ResponseWriter, name, title string, body template.HTML) {
	data := pageData{Title: title, Body: body}
	if err := templates.ExecuteTemplate(w, name, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (a *API) IndexPage(w http.ResponseWriter, r *http.Request) {
	devices, points, disabled := pointCounts(a.state)
	body := template.HTML(`<p>设备 ` + itoa(devices) + ` 台，点位 ` + itoa(points) + ` 个，停用设备 ` + itoa(disabled) + ` 台。</p>`)
	a.render(w, "index", "总览", body)
}

func (a *API) AlarmsPage(w http.ResponseWriter, r *http.Request) {
	rows := ""
	for _, al := range a.state.ActiveAlarms() {
		rows += `<tr><td>` + al.PointName + `</td><td>` + al.Status + `</td><td>` + al.LastTriggerAt.Format("2006-01-02 15:04:05") + `</td></tr>`
	}
	body := template.HTML(`<table><tr><th>点位</th><th>状态</th><th>最近触发</th></tr>` + rows + `</table>`)
	a.render(w, "alarms", "告警工作台", body)
}

func (a *API) RulesPage(w http.ResponseWriter, r *http.Request) {
	rows := ""
	for _, rl := range a.state.Rules() {
		rows += `<tr><td>` + rl.Name + `</td><td>` + rl.PointName + `</td><td>` + rl.Op + `</td><td>` + ftoa(rl.Threshold) + `</td><td>` + rl.State + `</td></tr>`
	}
	body := template.HTML(`<p>生效版本 ` + itoa64(a.state.Effective()) + `</p><table><tr><th>规则</th><th>点位</th><th>操作</th><th>阈值</th><th>状态</th></tr>` + rows + `</table>`)
	a.render(w, "rules", "规则配置", body)
}

func (a *API) PointsPage(w http.ResponseWriter, r *http.Request) {
	rows := ""
	for _, p := range a.state.Points() {
		rows += `<tr><td>` + p.Name + `</td><td>` + p.Unit + `</td><td>` + ftoa(p.LastValue) + `</td></tr>`
	}
	body := template.HTML(`<table><tr><th>点位</th><th>单位</th><th>最近值</th></tr>` + rows + `</table>`)
	a.render(w, "points", "点位监控", body)
}

func (a *API) AuditPage(w http.ResponseWriter, r *http.Request) {
	rows := ""
	for _, e := range a.state.AuditEntries() {
		rows += `<tr><td>` + e.Action + `</td><td>` + e.Target + `</td><td>` + e.Detail + `</td><td>` + e.At.Format("2006-01-02 15:04:05") + `</td></tr>`
	}
	body := template.HTML(`<table><tr><th>动作</th><th>对象</th><th>详情</th><th>时间</th></tr>` + rows + `</table>`)
	a.render(w, "audit", "审计日志", body)
}

func itoa(v int) string { return formatInt(int64(v)) }
func itoa64(v int64) string { return formatInt(v) }
func ftoa(v float64) string { return formatFloat(v) }
