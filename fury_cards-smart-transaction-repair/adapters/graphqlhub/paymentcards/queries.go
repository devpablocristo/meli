package paymentcards

const queryTransactionAndCards = `query paymentAndWallet($user_id: String! $payment_id: String!) {
	wallet(user_id: $user_id) {
		id
		cards_ids
		cards {
		  id
		  is_test
		  business_mode
		  holder {
			kyc_identification_id
			version_id
			user_id
		  }
		  issuer_accounts {
			id
		  }
		}
	}
	payment(id: $payment_id) {
		id
		site_id
		status_detail
		status
		debit_transaction {
		  id
		  ... on Authorization {
			id
			acquirer_code
			capture_datetime
			capture_id
			environment
			operation {
			  is_international
			  creation_datetime
			  transmission_datetime
			  expiration_date
			  installments
			  acquirer_code
			  stan
			  transaction {
				amount
				additional_amounts {
				  amount
				  type
				}
				total_amount
				reversed_amount
				original_amount
				increased_amount
				decimal_digits
				currency
				captured_amount
			  }
			  card {
				country
				number_id
			  }
			  subtype
			  is_advice
			  billing {
				amount
				currency
				decimal_digits
				total_amount
				conversion {					
					decimal_digits
					date
					from
				  }
			  }
			  settlement {
				amount
				currency
				decimal_digits
			  }			 
			}
			reversal_ids
			status
			type
			card_acceptor {
			  terminal
			  name
			}
		  }
		  provider {
			id
		  }		  
		}
		client_id
		payer_id
	  }
  }
`
