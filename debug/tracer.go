package debug

import (
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/vm"
)

// NewTracer returns a [vm.EVMLogger] that prints debug calls.
func NewTracer(tb testing.TB) vm.EVMLogger {
	tb.Helper()
	return &tracer{tb: tb}
}

type tracer struct {
	tb testing.TB
}

func (*tracer) CaptureTxStart(uint64) {}
func (*tracer) CaptureTxEnd(uint64)   {}
func (*tracer) CaptureStart(*vm.EVM, common.Address, common.Address, bool, []byte, uint64, *big.Int) {
}
func (*tracer) CaptureEnd([]byte, uint64, time.Duration, error) {}

func (t *tracer) CaptureEnter(typ vm.OpCode, from common.Address, to common.Address, input []byte, gas uint64, value *big.Int) {
	if to != addr || len(input) < 4 {
		return
	}

	sel := *(*[4]byte)(input[:4])
	args, ok := funcs[sel]
	if ok {
		t.logConsole(args, input[4:])
	} else {
		t.logMemory(input)
	}
}

func (*tracer) CaptureExit([]byte, uint64, error) {}
func (*tracer) CaptureState(uint64, vm.OpCode, uint64, uint64, *vm.ScopeContext, []byte, int, error) {
}
func (*tracer) CaptureFault(uint64, vm.OpCode, uint64, uint64, *vm.ScopeContext, int, error) {}

func (t *tracer) logConsole(args abi.Arguments, data []byte) {
	params, err := args.Unpack(data)
	if err != nil {
		t.tb.Fatalf("malformed log: %v", err)
	}

	for i, p := range params {
		switch p := p.(type) {
		case []byte:
			params[i] = hexutil.Bytes(p)
		}
	}

	t.tb.Log(params...)
}

func (t *tracer) logMemory(data []byte) {
	var str string
	for i := 0; i < len(data); i += 32 {
		end := i + 32
		if end > len(data) {
			end = len(data)
		}
		str += fmt.Sprintf("\n%4x | %x", i, data[i:end])
	}

	t.tb.Log("Memory:", str)
}