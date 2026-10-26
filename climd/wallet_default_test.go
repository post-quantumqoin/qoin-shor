package cli

import (
	"context"
	"errors"
	"flag"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/post-quantumqoin/address"
	"github.com/post-quantumqoin/qoin-shor/api"
	ucli "github.com/urfave/cli/v2"
)

func TestSetDefaultPrintsOnlyAfterSuccess(t *testing.T) {
	addr, err := address.NewPqcAddress([]byte("cli default fixture"))
	if err != nil {
		t.Fatal(err)
	}
	for _, failure := range []bool{false, true} {
		service := &api.FullNodeStruct{}
		called := false
		service.Internal.WalletSetDefault = func(_ context.Context, got address.Address) error {
			called = true
			if got != addr {
				return errors.New("wrong address")
			}
			if failure {
				return errors.New("address validation failed")
			}
			return nil
		}
		app := ucli.NewApp()
		app.Metadata = map[string]interface{}{"test-full-api": service}
		flags := flag.NewFlagSet("set-default", flag.ContinueOnError)
		if err := flags.Parse([]string{addr.String()}); err != nil {
			t.Fatal(err)
		}
		ctx := ucli.NewContext(app, flags, nil)
		reader, writer, err := os.Pipe()
		if err != nil {
			t.Fatal(err)
		}
		previous := os.Stdout
		os.Stdout = writer
		actionErr := walletSetDefault.Action(ctx)
		os.Stdout = previous
		writer.Close()
		output, readErr := io.ReadAll(reader)
		reader.Close()
		if readErr != nil {
			t.Fatal(readErr)
		}
		if !called || (actionErr != nil) != failure {
			t.Fatalf("unexpected result: %v", actionErr)
		}
		printed := strings.Contains(string(output), "Default address set to:")
		if printed == failure {
			t.Fatalf("incorrect success output on failure=%v: %q", failure, output)
		}
	}
}
