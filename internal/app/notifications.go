package app

import (
	"log/slog"

	"github.com/lusoris/revenge/internal/config"
	"github.com/lusoris/revenge/internal/service/notification"
	"github.com/lusoris/revenge/internal/service/notification/agents"
)

// registerNotificationAgents reads notification agent configs and registers them
// with the dispatcher. This runs at startup after the dispatcher is created.
func registerNotificationAgents(cfg *config.Config, dispatcher *notification.Dispatcher) {
	nc := cfg.Notifications

	for _, wc := range nc.Webhooks {
		if !wc.Enabled {
			continue
		}
		agent, err := agents.NewWebhookAgent(agents.WebhookConfig{
			AgentConfig: toAgentConfig(wc.Enabled, wc.Name, wc.EventTypes, wc.EventCategories),
			URL:         wc.URL,
			Method:      wc.Method,
			Headers:     wc.Headers,
			ContentType: wc.ContentType,
		})
		if err != nil {
			slog.Error("failed to create webhook agent", "name", wc.Name, "error", err)
			continue
		}
		if err := dispatcher.RegisterAgent(agent); err != nil {
			slog.Error("failed to register webhook agent", "name", wc.Name, "error", err)
		} else {
			slog.Info("registered notification agent", "type", "webhook", "name", wc.Name)
		}
	}

	for _, dc := range nc.Discord {
		if !dc.Enabled {
			continue
		}
		agent, err := agents.NewDiscordAgent(agents.DiscordConfig{
			AgentConfig: toAgentConfig(dc.Enabled, dc.Name, dc.EventTypes, dc.EventCategories),
			WebhookURL:  dc.WebhookURL,
			Username:    dc.Username,
			AvatarURL:   dc.AvatarURL,
		})
		if err != nil {
			slog.Error("failed to create discord agent", "name", dc.Name, "error", err)
			continue
		}
		if err := dispatcher.RegisterAgent(agent); err != nil {
			slog.Error("failed to register discord agent", "name", dc.Name, "error", err)
		} else {
			slog.Info("registered notification agent", "type", "discord", "name", dc.Name)
		}
	}

	for _, gc := range nc.Gotify {
		if !gc.Enabled {
			continue
		}
		agent, err := agents.NewGotifyAgent(agents.GotifyConfig{
			AgentConfig:     toAgentConfig(gc.Enabled, gc.Name, gc.EventTypes, gc.EventCategories),
			ServerURL:       gc.ServerURL,
			AppToken:        gc.AppToken,
			DefaultPriority: gc.DefaultPriority,
		})
		if err != nil {
			slog.Error("failed to create gotify agent", "name", gc.Name, "error", err)
			continue
		}
		if err := dispatcher.RegisterAgent(agent); err != nil {
			slog.Error("failed to register gotify agent", "name", gc.Name, "error", err)
		} else {
			slog.Info("registered notification agent", "type", "gotify", "name", gc.Name)
		}
	}

	for _, ntfyCfg := range nc.Ntfy {
		if !ntfyCfg.Enabled {
			continue
		}
		agent, err := agents.NewNtfyAgent(agents.NtfyConfig{
			AgentConfig:     toAgentConfig(ntfyCfg.Enabled, ntfyCfg.Name, ntfyCfg.EventTypes, ntfyCfg.EventCategories),
			ServerURL:       ntfyCfg.ServerURL,
			Topic:           ntfyCfg.Topic,
			AccessToken:     ntfyCfg.AccessToken,
			DefaultPriority: ntfyCfg.DefaultPriority,
		})
		if err != nil {
			slog.Error("failed to create ntfy agent", "name", ntfyCfg.Name, "error", err)
			continue
		}
		if err := dispatcher.RegisterAgent(agent); err != nil {
			slog.Error("failed to register ntfy agent", "name", ntfyCfg.Name, "error", err)
		} else {
			slog.Info("registered notification agent", "type", "ntfy", "name", ntfyCfg.Name)
		}
	}
}

// toAgentConfig converts config string slices to typed event/category slices.
func toAgentConfig(enabled bool, name string, eventTypes, eventCategories []string) notification.AgentConfig {
	ac := notification.AgentConfig{
		Enabled: enabled,
		Name:    name,
	}
	for _, et := range eventTypes {
		ac.EventTypes = append(ac.EventTypes, notification.EventType(et))
	}
	for _, ec := range eventCategories {
		ac.EventCategories = append(ac.EventCategories, notification.EventCategory(ec))
	}
	return ac
}
