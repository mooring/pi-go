package tui

import (
	"testing"
)

// /sidebar off hides the panel via the command path the mouse path cannot
// reach in SSH. The proof is the same one the click test uses: the sidebar
// column shrinks to its strip and the chat reclaims the columns, so a later
// frame is not wrapped for the old, wider chat.
func TestSidebarSlashOffHidesThePanel(t *testing.T) {
	m := historyModel(t, "")
	m.width, m.height = 120, 24
	m.applyResize()
	m.cfg.AppVersion = "1.4.2"
	m.chatModel.Messages = append(m.chatModel.Messages, realisticChat()...)
	m.syncPalette()
	m.View()

	if m.sidebarWidth() != SidebarWidth {
		t.Fatalf("expanded sidebarWidth = %d, want %d", m.sidebarWidth(), SidebarWidth)
	}

	next, _ := m.handleSlashCommand("/sidebar off")
	mm, ok := next.(*model)
	if !ok {
		t.Fatalf("handleSlashCommand returned %T, want *model", next)
	}
	if !mm.sidebarCollapsed {
		t.Fatal("/sidebar off did not set sidebarCollapsed")
	}
	if got := mm.sidebarWidth(); got != SidebarStripWidth {
		t.Errorf("collapsed sidebarWidth = %d, want %d", got, SidebarStripWidth)
	}
	if got, want := mm.mainWidth(), 120-SidebarStripWidth; got != want {
		t.Errorf("collapsed mainWidth = %d, want %d", got, want)
	}

	// A collapsed frame is still a rectangle.
	assertFrameWidth(t, mm, 120)
}

// /sidebar on (and its synonym, no arguments) brings the panel back, so the
// same command that hides it in SSH can also restore it.
func TestSidebarSlashOnRestoresThePanel(t *testing.T) {
	m := historyModel(t, "")
	m.width, m.height = 120, 24
	m.applyResize()
	m.cfg.AppVersion = "1.4.2"
	m.chatModel.Messages = append(m.chatModel.Messages, realisticChat()...)
	m.syncPalette()
	m.View()

	next, _ := m.handleSlashCommand("/sidebar off")
	mm, ok := next.(*model)
	if !ok {
		t.Fatalf("handleSlashCommand returned %T, want *model", next)
	}

	next, _ = mm.handleSlashCommand("/sidebar on")
	mm, ok = next.(*model)
	if !ok {
		t.Fatalf("handleSlashCommand returned %T, want *model", next)
	}
	if mm.sidebarCollapsed {
		t.Fatal("/sidebar on did not clear sidebarCollapsed")
	}
	if got, want := mm.mainWidth(), 120-SidebarWidth; got != want {
		t.Errorf("restored mainWidth = %d, want %d", got, want)
	}
}

// /sidebar with no arguments is treated as /sidebar on: a plain show request.
func TestSidebarSlashNoArgsShowsThePanel(t *testing.T) {
	m := historyModel(t, "")
	m.width, m.height = 120, 24
	m.applyResize()
	m.cfg.AppVersion = "1.4.2"
	m.chatModel.Messages = append(m.chatModel.Messages, realisticChat()...)
	m.syncPalette()
	m.View()

	next, _ := m.handleSlashCommand("/sidebar off")
	mm, ok := next.(*model)
	if !ok {
		t.Fatalf("handleSlashCommand returned %T, want *model", next)
	}

	next, _ = mm.handleSlashCommand("/sidebar")
	mm, ok = next.(*model)
	if !ok {
		t.Fatalf("handleSlashCommand returned %T, want *model", next)
	}
	if mm.sidebarCollapsed {
		t.Fatal("/sidebar with no args did not show the panel")
	}
	if got, want := mm.mainWidth(), 120-SidebarWidth; got != want {
		t.Errorf("no-arg sidebar mainWidth = %d, want %d", got, want)
	}
}

// The command is a documented, discoverable part of the slash set.
func TestSidebarSlashCommandIsRegistered(t *testing.T) {
	if _, ok := slashCommandByName["/sidebar"]; !ok {
		t.Fatal("/sidebar is not registered in slashCommandByName")
	}
}
