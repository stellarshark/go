package processors

import (
	"context"
	"fmt"
	"os"

	"github.com/stellar/go/exp/ingest/io"
	"github.com/stellar/go/exp/ingest/pipeline"
	"github.com/stellar/go/xdr"
)

func (p *CSVPrinter) ProcessState(ctx context.Context, store *pipeline.Store, r io.StateReadCloser, w io.StateWriteCloser) error {
	defer r.Close()
	defer w.Close()

	f, err := os.Create(p.Filename)
	if err != nil {
		return err
	}

	defer f.Close()

	for {
		entry, err := r.Read()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return err
			}
		}

		switch entry.Data.Type {
		case xdr.LedgerEntryTypeAccount:
			fmt.Fprintf(
				f,
				"%s,%d,%d\n",
				entry.Data.Account.AccountId.Address(),
				entry.Data.Account.Balance,
				entry.Data.Account.SeqNum,
			)
		default:
			// Ignore for now
		}

		select {
		case <-ctx.Done():
			return nil
		default:
			continue
		}
	}

	return nil
}

func (n *CSVPrinter) IsConcurrent() bool {
	return false
}

func (p *CSVPrinter) Name() string {
	return "CSVPrinter"
}
