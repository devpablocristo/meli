package graphqltransaction

const queryPretty = `query payment($id: String!) {
	payment(id: $id) {
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
					rate
					decimal_digits
					date
					from
				  }
			  }
			  settlement {
				amount
				currency
				decimal_digits
				total_amount
				conversion {
					rate
					decimal_digits
					date
					from
				  }
			  }			 
			}
			reversal_ids
			status			
			card_acceptor {
			  terminal
			  name
			}
		  }
		  provider {
			id
		  }
		  type		  
		}
		client_id
		payer_id
	  }
  }
`
