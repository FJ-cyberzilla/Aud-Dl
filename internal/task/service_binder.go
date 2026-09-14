package task

import (
	"audio-command-center/internal/gate"
	"audio-command-center/internal/processor"
	"audio-command-center/internal/search"
	"audio-command-center/internal/strategy"
	"audio-command-center/internal/vault"
)

type ServiceBinder struct {
	Gate       *gate.EgressTLSClient
	Strategy   *strategy.AntibotConductor
	Search     *search.SearchManager
	Transcoder *processor.PipeTranscoder
	Vault      *vault.VaultAdmin
}

func NewServiceBinder(
	g *gate.EgressTLSClient,
	s *strategy.AntibotConductor,
	sm *search.SearchManager,
	pt *processor.PipeTranscoder,
	v *vault.VaultAdmin,
) *ServiceBinder {
	return &ServiceBinder{
		Gate:       g,
		Strategy:   s,
		Search:     sm,
		Transcoder: pt,
		Vault:      v,
	}
}
