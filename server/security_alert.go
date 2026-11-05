package server

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/Aurox/artex/notify"
)

// 安全事件推送：把审计里的关键事件推到已配置的通知渠道（钉钉/飞书/企业微信/邮件/
// webhook/Telegram），让值守的人在演练或夜间也能第一时间收到。
//
// 设计要点：
//   - 复用现有 notify 渠道与 Send 接口，不新增投递链路（渠道配置、凭据、限流都在
//     既有系统里）；
//   - 只推"需要人看"的事件：撞库锁定、改密、范围门禁拦截；普通写操作只落审计表；
//   - 按 action+target 做 60 秒节流：演练时 Agent 可能连续触发几十次门禁拒绝，
//     不节流会把 IM 刷爆并撞上渠道限流，反而让关键告警丢失；
//   - 发送异步 + 超时，失败只打日志：告警不能拖慢请求，也不能成为可用性单点；
//   - 开关：settings.security_alert_enabled（默认开启），维护期可一键止血。
const settingSecurityAlert = "security_alert_enabled"

// securityAlertMinInterval 同一 action+target 的最小推送间隔。
const securityAlertMinInterval = 60 * time.Second

var (
	securityAlertMu   sync.Mutex
	securityAlertLast = map[string]time.Time{}
)

// securityAlertActions 是需要推送的审计动作白名单。
func securityAlertWanted(action string) bool {
	switch action {
	case "auth.login.locked", "auth.password_change", "scope.denied", "auth.init":
		return true
	}
	return false
}

// securityAlert 异步推送一条安全事件到所有启用的通知渠道。
func (s *Server) securityAlert(r *http.Request, action, target, detail string) {
	if s == nil || s.m == nil || s.m.pg == nil {
		return
	}
	if !securityAlertWanted(action) {
		return
	}
	// 节流：同一 action+target 60 秒内只推一次。
	key := action + "|" + target
	securityAlertMu.Lock()
	if last, ok := securityAlertLast[key]; ok && time.Since(last) < securityAlertMinInterval {
		securityAlertMu.Unlock()
		return
	}
	securityAlertLast[key] = time.Now()
	// 避免 map 无限增长：数量大时清理过期项。
	if len(securityAlertLast) > 512 {
		for k, v := range securityAlertLast {
			if time.Since(v) > 10*time.Minute {
				delete(securityAlertLast, k)
			}
		}
	}
	securityAlertMu.Unlock()

	ip := ""
	if r != nil {
		ip = clientIP(r)
	}
	pg := s.m.pg
	jwtKey := s.jwtKey
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), notifySendTimeout)
		defer cancel()
		enabled, _, err := pg.GetSetting(settingSecurityAlert)
		if err != nil {
			log.Printf("[security-alert] 读取开关失败: %v", err)
			return
		}
		if strings.TrimSpace(enabled) == "0" {
			return
		}
		channels, err := pg.ListNotificationChannels(ctx)
		if err != nil {
			log.Printf("[security-alert] 读取通知渠道失败: %v", err)
			return
		}
		base, _, _ := pg.GetSetting(settingNotifyPublicBaseURL)
		base = strings.TrimRight(strings.TrimSpace(base), "/")
		detailURL := ""
		if base != "" {
			detailURL = base + "/system/audit"
		}
		title := securityAlertTitle(action)
		summary := "目标: " + target + "；来源 IP: " + ip
		if detail != "" {
			summary += "；" + detail
		}
		msg := notify.Message{
			Items: []notify.Item{{
				Name:      title,
				VulnClass: "安全事件",
				Severity:  "high",
				Summary:   summary,
				Assets:    []string{ip},
				DetailURL: detailURL,
			}},
			HomeURL: base,
		}
		for _, ch := range channels {
			// Enabled 是 *bool（区分「没传」与「显式 false」），nil 视为启用。
			if ch.Enabled != nil && !*ch.Enabled {
				continue
			}
			impl, ok := notify.Get(ch.Kind)
			if !ok {
				continue
			}
			var cfg map[string]any
			if len(ch.Config) > 0 {
				if err := json.Unmarshal(ch.Config, &cfg); err != nil {
					log.Printf("[security-alert] 渠道 %d 配置解析失败: %v", ch.ID, err)
					continue
				}
			}
			sendCtx, cancelSend := context.WithTimeout(ctx, notifySendTimeout)
			_, err := impl.Send(sendCtx, cfg, msg)
			cancelSend()
			if err != nil && !notify.IsPermanent(err) {
				log.Printf("[security-alert] 渠道 %d(%s) 推送失败: %v", ch.ID, ch.Kind, err)
			}
		}
		_ = jwtKey
	}()
}

func securityAlertTitle(action string) string {
	switch action {
	case "auth.login.locked":
		return "【安全告警】登录失败次数过多，已临时锁定"
	case "auth.password_change":
		return "【安全事件】管理员口令已修改"
	case "scope.denied":
		return "【安全告警】网络范围门禁拦截了一次出站请求"
	case "auth.init":
		return "【安全事件】平台完成首次口令初始化"
	}
	return "【安全事件】" + action
}
