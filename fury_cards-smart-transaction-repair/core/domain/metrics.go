package domain

// Spans.
const (
	MetricRoot = "cards.smart.transaction.repair."
)

// Metrics.
const (
	MetricReverse               = MetricRoot + "reverse"
	MetricUnable                = MetricRoot + "unable"
	MetricRule                  = MetricRoot + "rule_not_elegible"
	MetricBlockedUser           = MetricRoot + "blockedlist"
	MetricTimer                 = MetricRoot + "timer"
	MetricEventValidationResult = MetricRoot + "event_validation_result"
)

// Tags Reverse.
const (
	MetricTagReverseSuccess     = "reverse_success"
	MetricTagReverseFail        = "reverse_fail"
	MetricTagReverseAlready     = "reverse_already_done"
	MetricTagReverseEligible    = "reverse_eligible"
	MetricTagReverseNotEligible = "reverse_not_eligible"
	MetricTagReverseUnable      = "reverse_unable"
	MetricTagReverseWarning     = "reverse_warning"
)

// Tags Reverse Unable.
const (
	MetricTagUnablePersonIDEmpty                   = "person_id_empty"
	MetricTagUnableTypeDiffAuthorization           = "type_diff_authorization"
	MetricTagUnableUserNotAvailable                = "user_not_available"
	MetricTagAuthorizationAlreadyReversed          = "authorization_already_reversed"
	MetricTagInvalidPrerequisitesForRuleValidation = "invalid_prerequisites_for_rule_validation"
)

// Tags Not Eligible Validation Result.
const (
	MetricTagRuleMaxAmount                 MetricTagRule = "rule_max_amount"
	MetricTagRuleMaxQtyRepairPerPeriod     MetricTagRule = "rule_qty_repair_per_period"
	MetricTagRuleUserIsBlocked             MetricTagRule = "rule_blocked_user"
	MetricTagRuleStatusNotAllowed          MetricTagRule = "rule_status_not_allowed"
	MetricTagRuleUserWithRestriction       MetricTagRule = "rule_user_with_restriction"
	MetricTagRuleMaxQtyChargebackPerPeriod MetricTagRule = "rule_qty_chargeback_per_period"
	MetricTagRuleSubTypeNotAllowed         MetricTagRule = "rule_subtype_not_allowed"
	MetricTagMerchantBlockedlist           MetricTagRule = "rule_merchant_blockedlist"
	MetricTagUserLifetime                  MetricTagRule = "rule_user_lifetime"
	MetricTagMinTotalHistoricalAmount      MetricTagRule = "rule_min_total_historical_amount"
	MetricTagMinTransactionsQty            MetricTagRule = "rule_min_transactions_qty"
	MetricTagScore                         MetricTagRule = "rule_score"
)

// Tags BlockedUser.
const (
	MetricTagBlockedUserCheckRepairByCapture             = "blockedlist_check_repair_by_capture"
	MetricTagBlockedUserCheckRepairByCaptureInBulk       = "blockedlist_check_repair_by_capture_in_bulk"
	MetricTagBlockedUserBlocked                          = "blockedlist_blocked"
	MetricTagBlockedUserBlockedAgain                     = "blockedlist_blocked_again"
	MetricTagBlockedUserUnblocked                        = "blockedlist_unblocked"
	MetricTagBlockedUserUnblockedByReversalClearing      = "blockedlist_unblocked_by_reversal_clearing"
	MetricTagBlockedUserCheckBlockedUserByReversal       = "blockedlist_check_blocked_users_by_reversal"
	MetricTagBlockedUserCheckBlockedUserByReversalInBulk = "blockedlist_check_blocked_users_by_reversal_in_bulk"
	MetricTagBlockedUserFail                             = "blockedlist_fail"
)

// Tags Events in general.
const (
	MetricTagEventReceived = "event_received"
	MetricTagEventFail     = "event_fail"
	MetricTagEventDone     = "event_done"
)

// Tags Queue.
const (
	MetricTagSourceQCaptureTransactions  = "q_capture_transactions"
	MetricTagSourceQReversalTransactions = "q_reversal_transactions"
)

// Tags Timer.
const (
	MetricTagTimerReverseHandler                           = TimerHandlerReverse
	MetricTagTimerConsumerCaptures                         = TimerConsumerCaptures
	MetricTagTimerConsumerReverses                         = TimerConsumerReversals
	MetricTagTimerConsumerValidationResult                 = TimerConsumerValidationResult
	MetricTagTimerDsRepairGetByPaymentIDAndUserID          = TimerDsRepairGetByPaymentIDAndUserID
	MetricTagTimerCardsTransactionReverse                  = TimerCardsTransactionReverse
	MetricTagTimerConfigurationsLoadParametersValidation   = TimerConfigurationsLoadParametersValidation
	MetricTagTimerDsRepairGetListByUserIDAndCreationPeriod = TimerDsRepairGetListByUserIDAndCreationPeriod
	MetricTagTimerDsValidationSave                         = TimerDsValidationSave
	MetricTagTimerKvsRepairSave                            = TimerKvsRepairSave
	MetricTagTimerKvsBlockelistGet                         = TimerKvsBlockelistGet
	MetricTagTimerKvsBlockelistSave                        = TimerKvsBlockelistSave
	MetricTagTimerDsRepairGetList                          = TimerDsRepairGetList
	MetricTagTimerDsBlockedUserGetListByAuthorizationIDs   = TimerDsBlockedUserGetListByAuthorizationIDs
	MetricTagTimerPolicyAgentEvaluate                      = TimerPolicyAgentEvaluate
	MetricTagTimerSearchHubCountChargebackFromUser         = TimerSearchHubCountChargebackFromUser
	MetricTagTimerSearchKycGetUserKyc                      = TimerSearchKycGetUserKyc
	MetricTagTimerSearchHubGetTransactionHistory           = TimerSearchHubGetTransactionHistory
	MetricTagTimerGetScore                                 = TimerGetScore
	MetricTagTimerEventValidationResultPublish             = TimerEventValidationResultPublish
	MetricTagTimerGraphqlGetPaymentTransactionByPaymentID  = TimerGraphqlGetPaymentTransactionByPaymentID
)

type MetricTagRule string
