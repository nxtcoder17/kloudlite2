package app

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/kloudlite/api/apps/console/internal/domain"
	"github.com/kloudlite/operator/pkg/errors"
	"github.com/miekg/dns"
)

type dnsHandler struct {
	logger               *slog.Logger
	serviceBindingDomain domain.ServiceBindingDomain
	kloudliteDNSSuffix   string
}

const (
	DefaultDNSTTL = 5
)

func (h *dnsHandler) ServeDNS(w dns.ResponseWriter, r *dns.Msg) {
	logger := h.logger.With("query", r.Question[0].Name)
	logger.Debug("INCOMING dns request")
	start := time.Now()

	msg := new(dns.Msg)
	msg.SetReply(r)
	msg.Authoritative = true

	ctx, cf := context.WithCancel(context.TODO())
	defer cf()

	for _, question := range r.Question {
		answers, err := h.resolver(ctx, question.Name, question.Qtype)
		if err != nil {
			logger.Error("FAILED to resolve dns record, got", "err", err, "question", question.Name)
			msg.Rcode = dns.RcodeNameError
			continue
		}
		msg.Answer = append(msg.Answer, answers...)
	}

	w.WriteMsg(msg)
	if msg.Rcode != dns.RcodeNameError {
		logger.Info("RESOLVED dns request", "answers", msg.Answer, "took", fmt.Sprintf("%.2fs", time.Since(start).Seconds()))
	}
}

func (h *dnsHandler) newRR(domain string, ttl int, ip string) ([]dns.RR, error) {
	r, err := dns.NewRR(fmt.Sprintf("%s %d IN A %s", domain, ttl, ip))
	if err != nil {
		return nil, errors.NewEf(err, "failed to create dns record")
	}
	return []dns.RR{r}, nil
}

var errNoServiceBinding = errors.Newf("no service binding found")

func (h *dnsHandler) resolver(ctx context.Context, domain string, qtype uint16) ([]dns.RR, error) {
	m := new(dns.Msg)
	m.SetQuestion(dns.Fqdn(domain), qtype)
	m.RecursionDesired = true

	question := m.Question[0]
	sp := strings.SplitN(question.Name, fmt.Sprintf(".%s", h.kloudliteDNSSuffix), 2)
	if len(sp) < 2 {
		return nil, fmt.Errorf("failed to split into 2 over .%s", h.kloudliteDNSSuffix)
	}

	comps := strings.Split(sp[0], ".")
	accountName := comps[len(comps)-1]
	hostname := strings.Join(comps[:len(comps)-1], ".")

	sb, err := h.serviceBindingDomain.FindServiceBindingByHostname(ctx, accountName, hostname)
	if err != nil {
		return nil, errors.NewEf(err, "failed to find service binding")
	}

	if sb == nil {
		return nil, errNoServiceBinding
	}

	return h.newRR(question.Name, DefaultDNSTTL, sb.Spec.GlobalIP)
}
