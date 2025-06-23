package CommandProcessing

type CommandArgs struct {
	DbNames             []string
	OperationType       string
	PostgresPasswordURL string
	HavingHelpFlag      bool // для флагов, которые выполняют некоторые действия, но не с СУБД
	HelperDump          bool
}
