package velafi

import (
	"context"
	"strconv"
)

func (c *Client) ListMerchantAccounts(ctx context.Context, params *ListMerchantAccountsParams) (*MerchantAccountList, error) {
	q := map[string]string{
		"fiat": params.Fiat,
	}
	if params.MerchantID > 0 {
		q["merchantId"] = strconv.Itoa(params.MerchantID)
	}
	if params.MerchantName != "" {
		q["merchantName"] = params.MerchantName
	}
	if params.Email != "" {
		q["email"] = params.Email
	}
	if params.Status > 0 {
		q["status"] = strconv.Itoa(params.Status)
	}
	if params.CurrentPage > 0 {
		q["currentPage"] = strconv.Itoa(params.CurrentPage)
	}
	if params.PageSize > 0 {
		q["pageSize"] = strconv.Itoa(params.PageSize)
	}
	var result MerchantAccountList
	err := c.get(ctx, "/v2/merchant/accounts"+buildQuery(q), &result)
	return &result, err
}

func (c *Client) ActivateMerchantAccount(ctx context.Context, params *ActivateMerchantAccountParams) (*ActivateMerchantAccountResult, error) {
	var result ActivateMerchantAccountResult
	err := c.post(ctx, "/v2/merchant/accounts", params, &result)
	return &result, err
}

func (c *Client) GetPendingFundAccount(ctx context.Context, params *GetPendingFundAccountParams) (*PendingFundAccount, error) {
	q := map[string]string{
		"fiat": params.Fiat,
	}
	if params.MerchantID > 0 {
		q["merchantId"] = strconv.Itoa(params.MerchantID)
	}
	if params.PaymentID > 0 {
		q["paymentId"] = strconv.Itoa(params.PaymentID)
	}
	if params.DepositAlias != "" {
		q["depositAlias"] = params.DepositAlias
	}
	if params.Amount != "" {
		q["amount"] = params.Amount
	}
	var result PendingFundAccount
	err := c.get(ctx, "/v2/merchant/addfunds"+buildQuery(q), &result)
	return &result, err
}

func (c *Client) ClaimPendingFund(ctx context.Context, params *ClaimPendingFundParams) (*ClaimPendingFundResult, error) {
	var result ClaimPendingFundResult
	err := c.post(ctx, "/v2/merchant/claimfunds", params, &result)
	return &result, err
}

func (c *Client) GetFundingRecords(ctx context.Context, params *GetFundingRecordsParams) (*FundingRecordList, error) {
	q := map[string]string{
		"fiat":     params.Fiat,
		"type":     params.Type,
		"clientId": params.ClientID,
	}
	if params.TxID > 0 {
		q["txId"] = strconv.FormatInt(params.TxID, 10)
	}
	if params.StartTime > 0 {
		q["startTime"] = strconv.FormatInt(params.StartTime, 10)
	}
	if params.EndTime > 0 {
		q["endTime"] = strconv.FormatInt(params.EndTime, 10)
	}
	if params.MerchantID > 0 {
		q["merchantId"] = strconv.Itoa(params.MerchantID)
	}
	if params.Status > 0 {
		q["status"] = strconv.Itoa(params.Status)
	}
	if params.CurrentPage > 0 {
		q["currentPage"] = strconv.Itoa(params.CurrentPage)
	}
	if params.PageSize > 0 {
		q["pageSize"] = strconv.Itoa(params.PageSize)
	}
	var result FundingRecordList
	err := c.get(ctx, "/v2/merchant/funding/records"+buildQuery(q), &result)
	return &result, err
}

func (c *Client) GetFiatFlow(ctx context.Context, params *GetFiatFlowParams) (*FiatFlowList, error) {
	q := map[string]string{
		"fiat":     params.Fiat,
		"flowType": params.FlowType,
	}
	if params.StartTime > 0 {
		q["startTime"] = strconv.FormatInt(params.StartTime, 10)
	}
	if params.EndTime > 0 {
		q["endTime"] = strconv.FormatInt(params.EndTime, 10)
	}
	if params.MerchantID > 0 {
		q["merchantId"] = strconv.Itoa(params.MerchantID)
	}
	if params.TxID > 0 {
		q["txId"] = strconv.FormatInt(params.TxID, 10)
	}
	if params.CurrentPage > 0 {
		q["currentPage"] = strconv.Itoa(params.CurrentPage)
	}
	if params.PageSize > 0 {
		q["pageSize"] = strconv.Itoa(params.PageSize)
	}
	var result FiatFlowList
	err := c.get(ctx, "/v2/merchant/fiat/flows"+buildQuery(q), &result)
	return &result, err
}

func (c *Client) GetCryptoBalance(ctx context.Context, merchantID int) (*CryptoBalanceResult, error) {
	q := map[string]string{
		"merchantId": strconv.Itoa(merchantID),
	}
	var result CryptoBalanceResult
	err := c.get(ctx, "/v2/assets/account/assets"+buildQuery(q), &result)
	return &result, err
}

func (c *Client) GetCryptoDepositAddress(ctx context.Context, params *GetCryptoDepositAddressParams) (*CryptoDepositAddress, error) {
	q := map[string]string{
		"tokenId": params.TokenID,
	}
	if params.MerchantID > 0 {
		q["merchantId"] = strconv.Itoa(params.MerchantID)
	}
	if params.ChainType != "" {
		q["chainType"] = params.ChainType
	}
	var result CryptoDepositAddress
	err := c.get(ctx, "/v2/assets/depositAddress"+buildQuery(q), &result)
	return &result, err
}
