package plugin

import (
	"encoding/json"
	"errors"
	"sync"
	"time"

	"github.com/Jeffail/benthos/v3/lib/log"
	"github.com/Jeffail/benthos/v3/lib/metrics"
	"github.com/Jeffail/benthos/v3/lib/output"
	"github.com/Jeffail/benthos/v3/lib/response"
	"github.com/Jeffail/benthos/v3/lib/types"
	"golang.org/x/net/context"
	"google.golang.org/api/sheets/v4"
)

func init() {
	output.RegisterPlugin(
		"sheets",
		func() interface{} {
			return NewSheetsConfig()
		},
		func(iconf interface{}, mgr types.Manager, logger log.Modular, stats metrics.Type) (types.Output, error) {
			conf, ok := iconf.(*SheetsConfig)
			if !ok {
				return nil, errors.New("Failed to cast config")
			}
			return NewSheetsOut(*conf, mgr, logger, stats)
		},
	)

	output.DocumentPlugin(
		"sheets",
		`This plugin outputs to a google spreadsheet`,
		nil,
	)
}

// SheetsConfig defines the plugin parameters
type SheetsConfig struct {
	SheetID string `json:"sheetId" yaml:"sheetId"`
}

// NewSheetsConfig create a new SheetsConfig
func NewSheetsConfig() *SheetsConfig {
	return &SheetsConfig{}
}

//------------------------------------------------------------------------------

// SheetsOut is an example plugin that creates gibberish messages.
type SheetsOut struct {
	transactionsChan <-chan types.Transaction

	log           log.Modular
	stats         metrics.Type
	sheetsService *sheets.Service
	sheetID       string

	closeOnce  sync.Once
	closeChan  chan struct{}
	closedChan chan struct{}
}

// NewSheetsOut creates a new example plugin output type.
func NewSheetsOut(
	conf SheetsConfig,
	mgr types.Manager,
	log log.Modular,
	stats metrics.Type,
) (output.Type, error) {

	// Create sheets service
	ctx := context.Background()
	sheetsService, err := sheets.NewService(ctx)
	if err != nil {
		log.Errorf("failed initialize SheetsService")
		panic(err)
	}

	e := &SheetsOut{
		log:           log,
		stats:         stats,
		sheetsService: sheetsService,
		sheetID:       conf.SheetID,

		closeChan:  make(chan struct{}),
		closedChan: make(chan struct{}),
	}

	return e, nil
}

//------------------------------------------------------------------------------

func (e *SheetsOut) loop() {
	defer func() {
		close(e.closedChan)
	}()

	for {
		var tran types.Transaction
		var open bool
		select {
		case tran, open = <-e.transactionsChan:
			if !open {
				return
			}
		case <-e.closeChan:
			return
		}
		tran.Payload.Iter(func(i int, p types.Part) error {
			e.log.Infoln("parsing output to interface array")
			var v []interface{}
			json.Unmarshal(p.Get(), &v)

			var vr sheets.ValueRange
			vr.Values = append(vr.Values, v)

			e.log.Infof("sending output to sheet (%s)", e.sheetID)
			_, err := e.sheetsService.Spreadsheets.Values.
				Append(e.sheetID, "A1", &vr).
				ValueInputOption("RAW").
				Do()

			if err != nil {
				e.log.Errorf("failed to send output to sheet (%s)", e.sheetID)
				panic(err)
			}

			e.log.Infof("sucessfully sent output to sheet (%s)", e.sheetID)
			return nil
		})
		select {
		case tran.ResponseChan <- response.NewAck():
		case <-e.closeChan:
			return
		}
	}
}

// Connected returns true if this output is currently connected to its target.
func (e *SheetsOut) Connected() bool {
	return true // We're always connected
}

// Consume starts this output consuming from a transaction channel.
func (e *SheetsOut) Consume(tChan <-chan types.Transaction) error {
	e.transactionsChan = tChan
	go e.loop()
	return nil
}

// CloseAsync shuts down the output and stops processing requests.
func (e *SheetsOut) CloseAsync() {
	e.closeOnce.Do(func() {
		close(e.closeChan)
	})
}

// WaitForClose blocks until the output has closed down.
func (e *SheetsOut) WaitForClose(timeout time.Duration) error {
	select {
	case <-e.closedChan:
	case <-time.After(timeout):
		return types.ErrTimeout
	}
	return nil
}
