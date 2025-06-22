package CommandProcessing

type CommandArgs struct {
	DbNames             []string
	OperationType       string
	PostgresPasswordURL string
	HelperDump          bool
	DebugInfo           bool
}
