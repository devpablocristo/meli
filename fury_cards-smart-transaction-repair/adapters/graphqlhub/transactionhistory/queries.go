package transactionhistory

const queryTransactionHistory = `
	query GetTransactionHistory($id: String!, $queryParams: SearchParams!) { 
		wallet(user_id: $id) {
			cards { 
				issuer_accounts { 
					id 
					transactions(queryParams: $queryParams) {
						match_total 
						transactions {... on Authorization {
							id 
							operation {
								settlement {
									currency 
									amount 
									decimal_digits
								} 
								billing {
									amount 
									currency 
									decimal_digits
								}
							} 
							capture_id
						} 
						type 
					}
				}
			}
		}
	}
}
`
