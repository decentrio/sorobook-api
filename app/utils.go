package app

const (
	PAGE_SIZE         = 10                               //nolint
	LEDGER_TABLE      = "ledgers"                        //nolint
	TRANSACTION_TABLE = "transactions"                   //nolint
	CONTRACT_TABLE    = "contracts_data"                 //nolint
	CONTRACT_CODES    = "contracts_codes"                //nolint
	INVOKE_TXS        = "invoke_transactions"            //nolint
	EVENT_TABLE       = "wasm_contract_events"           //nolint
	TRANSFER_TABLE    = "asset_contract_transfer_events" //nolint
	MINT_TABLE        = "asset_contract_mint_events"     //nolint
	BURN_TABLE        = "asset_contract_burn_events"     //nolint
	CLAWBACK_TABLE    = "asset_contract_clawback_events" //nolint
)
