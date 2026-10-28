package pqcaddress

import (
	"sort"

	"github.com/post-quantumqoin/address"
	typescrypto "github.com/post-quantumqoin/core-types/crypto"
	"golang.org/x/xerrors"
)

type LegacyAddressCandidate struct {
	Address   address.Address
	Selection string
	Encoding  string
}

// LegacyWalletAddressCandidates recognizes both historical wallet constructors.
// It does not reorder the certificate or change new-address generation rules.
func LegacyWalletAddressCandidates(cert typescrypto.SignPQCCert, protocol address.Protocol) ([]LegacyAddressCandidate, error) {
	if Scheme(cert.Version) != SchemeLegacyCompatible {
		return nil, xerrors.New("expected legacy PQC wallet certificate")
	}
	if err := ValidateCertificate(cert); err != nil {
		return nil, err
	}
	sorted := append([]typescrypto.SignPqcCertPubkey(nil), cert.Pubkeys...)
	sort.Slice(sorted, func(i, j int) bool { return sorted[i].Typ < sorted[j].Typ })
	selections := []struct {
		name string
		key  typescrypto.SignPqcCertPubkey
	}{
		{"original_first", cert.Pubkeys[0]},
		{"sorted_first", sorted[0]},
	}
	var result []LegacyAddressCandidate
	seen := make(map[address.Address]bool)
	for _, selection := range selections {
		candidates, err := SignerAddressCandidates(typescrypto.SignPQCCert{
			Version: cert.Version, Pubkeys: []typescrypto.SignPqcCertPubkey{selection.key},
		}, protocol)
		if err != nil {
			return nil, err
		}
		for index, candidate := range candidates {
			if seen[candidate] {
				continue
			}
			seen[candidate] = true
			encoding := "legacy_cbor"
			if index == 1 {
				encoding = "raw"
			}
			result = append(result, LegacyAddressCandidate{Address: candidate, Selection: selection.name, Encoding: encoding})
		}
	}
	return result, nil
}
