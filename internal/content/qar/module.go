// Package qar provides adult content domain models with "Queen Anne's Revenge" obfuscation.
// All adult content uses pirate-themed terminology for discretion.
//
// Obfuscation mapping:
//   - Performers → Crew
//   - Scenes → Voyages
//   - Movies → Expeditions
//   - Studios → Ports
//   - Tags → Flags
//   - Libraries → Fleets
package qar

import (
	"go.uber.org/fx"

	"github.com/lusoris/revenge/internal/content/qar/crew"
	"github.com/lusoris/revenge/internal/content/qar/expedition"
	"github.com/lusoris/revenge/internal/content/qar/flag"
	"github.com/lusoris/revenge/internal/content/qar/fleet"
	"github.com/lusoris/revenge/internal/content/qar/port"
	"github.com/lusoris/revenge/internal/content/qar/voyage"
)

// Module provides all QAR (adult content) dependencies for fx.
var Module = fx.Module("qar",
	fleet.Module,
	expedition.Module,
	voyage.Module,
	crew.Module,
	port.Module,
	flag.Module,
)
