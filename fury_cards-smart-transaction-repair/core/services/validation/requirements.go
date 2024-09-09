package validation

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/melisource/cards-smart-transaction-repair/core/domain"
)

// Method responsible for validating the data used in the rules.
func runPrerequisitesForRuleValidation(
	input domain.ValidationInput,
	rulesSite domain.ConfigRulesSite,
) error {

	err := confrontReceivedDataFromTransactional(input)
	if err != nil {
		return err
	}

	err = validateDataForRule(rulesSite, input.TransactionData, input.UserKycData)
	if err != nil {
		return err
	}

	return nil
}

func confrontReceivedDataFromTransactional(
	input domain.ValidationInput,
) error {

	transactionData := input.TransactionData

	if !strings.EqualFold(input.SiteID, transactionData.SiteID) {
		err := fmt.Errorf("received siteID(%s) different from database(%s)", input.SiteID, transactionData.SiteID)
		return buildErrorUnprocessableEntity(err)
	}

	return nil
}

func validateDataForRule(
	rulesSite domain.ConfigRulesSite,
	transactionData domain.TransactionData,
	userKycData domain.SearchKycOutput,
) error {

	err := runPrerequisiteToValidateRuleMaxRepairAmount(transactionData, rulesSite)
	if err != nil {
		return err
	}

	err = runPrerequisiteToValidateRuleAllowedStatus(transactionData.StatusDetail, rulesSite)
	if err != nil {
		return err
	}

	err = runPrerequisiteToValidateRuleAllowedSubType(transactionData.Operation.SubType, rulesSite)
	if err != nil {
		return err
	}

	err = runPrerequisiteToValidateRuleMerchantBlockedlist(transactionData.Operation.Authorization.CardAcceptor.Name, rulesSite)
	if err != nil {
		return err
	}

	err = runPrerequisiteToValidateRuleUserLifetime(userKycData.DateCreated, rulesSite)
	if err != nil {
		return err
	}

	return nil
}

func runPrerequisiteToValidateRuleMaxRepairAmount(
	transactionData domain.TransactionData,
	rulesSite domain.ConfigRulesSite,
) error {
	configMaxAmountReparation := rulesSite.MaxAmountReparation
	if configMaxAmountReparation == nil {
		return nil
	}

	var err error

	switch {
	case transactionData.Operation.Settlement.Amount == 0, len(strings.TrimSpace(transactionData.Operation.Settlement.Currency)) == 0:
		err = errors.New("settlement is empty in the transactional database")

	case transactionData.Operation.Billing.Amount == 0, len(strings.TrimSpace(transactionData.Operation.Billing.Currency)) == 0:
		err = errors.New("billing is empty in the transactional database")

	default:
		switch configMaxAmountReparation.Currency {
		case internationalCurrency:
			if configMaxAmountReparation.Currency != transactionData.Operation.Settlement.Currency {
				err = fmt.Errorf("transaction settlement currency(%s) different from the parameters(%s)",
					transactionData.Operation.Settlement.Currency, configMaxAmountReparation.Currency)
			}
		default:
			if configMaxAmountReparation.Currency != transactionData.Operation.Billing.Currency {
				err = fmt.Errorf("transaction billing currency(%s) different from the parameters(%s)",
					transactionData.Operation.Billing.Currency, configMaxAmountReparation.Currency)
			}
		}
	}

	if err != nil {
		return buildErrorUnprocessableEntity(err)
	}

	return nil
}

func runPrerequisiteToValidateRuleAllowedStatus(
	statusDetail string,
	rulesSite domain.ConfigRulesSite,
) error {

	if rulesSite.StatusDetailAllowed != nil && len(statusDetail) == 0 {
		err := errors.New("status detail transaction is empty")
		return buildErrorUnprocessableEntity(err)
	}

	return nil
}

func runPrerequisiteToValidateRuleAllowedSubType(
	subType string,
	rulesSite domain.ConfigRulesSite,
) error {

	if rulesSite.SubTypeAllowed != nil && len(subType) == 0 {
		err := errors.New("subtype transaction is empty")
		return buildErrorUnprocessableEntity(err)
	}

	return nil
}

func runPrerequisiteToValidateRuleMerchantBlockedlist(
	cardAcceptorName string,
	rulesSite domain.ConfigRulesSite,
) error {

	if rulesSite.MerchantBlockedlist != nil && len(cardAcceptorName) == 0 {
		err := errors.New("merchant name is empty")
		return buildErrorUnprocessableEntity(err)
	}

	return nil
}

func runPrerequisiteToValidateRuleUserLifetime(
	dateCreated time.Time,
	rulesSite domain.ConfigRulesSite,
) error {

	if rulesSite.UserLifetime != nil && dateCreated.IsZero() {
		err := errors.New("date created user is empty")
		return buildErrorUnprocessableEntity(err)
	}

	return nil
}
