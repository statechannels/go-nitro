package protocols

import (
	"math/big"

	"github.com/statechannels/go-nitro/channel"
)

func SignPreFundEffect(c channel.Channel) string {
	return "sign Prefundsetup for" + c.Id()
}
func SignPostFundEffect(c channel.Channel) string {
	return "sign Postfundsetup for" + c.Id()
}
func FundOnChainEffect(c channel.Channel, asset string, amount big.Int) string {
	return "deposit" + amount.Text(64) + "into" + c.Id()
}

// Pause points. By Returning these, the imperative shell can detect a lack of progress after multiple cranks
type WaitingFor string

var CompletePrefund = WaitingFor("CompletePrefund")
var MyTurnToFund = WaitingFor("MyTurnToFund")
var CompleteFunding = WaitingFor("CompleteFunding")
var CompletePostfund = WaitingFor("CompletePostfund")

// Crank inspects the objective o, channel in scope c, and holdings h -- and declares a list of Effects to be executed
func Crank(o Objective, c channel.Channel, h Holding) (SideEffects, WaitingFor, error) {
	// TODO handle an array of Holdings

	// Input validation
	if o.Status != "approved" {
		return NoSideEffects, "", ErrNotApproved
	}

	if o.Scope[0] != c.Id() {
		return NoSideEffects, "", ErrNotInScope
	}

	if h.ChannelId() != c.Id() {
		return NoSideEffects, "", ErrIncorrectChannelId
	}

	// Prefunding
	if c.ShouldSignPreFund() {
		return []string{SignPreFundEffect(c)}, "", nil
	}
	if !c.IsPrefundComplete() {
		return NoSideEffects, CompletePrefund, nil
	}

	// Funding

	fundingComplete := c.IsFundingComplete(h.Asset(), h.Amount())
	amountToDeposit, safeToDeposit := c.AmountToDeposit(h.Asset(), h.Amount())

	if !fundingComplete && !safeToDeposit {
		return []string{}, MyTurnToFund, nil
	}

	if !fundingComplete && amountToDeposit.Cmp(zero) != 0 && safeToDeposit {
		var effects = make([]string, 0) // TODO loop over assets
		effects = append(effects, FundOnChainEffect(c, h.Asset(), amountToDeposit))
		if len(effects) > 0 {
			return effects, "", nil
		}
	}

	if !fundingComplete {
		return NoSideEffects, CompleteFunding, nil
	}

	// Postfunding
	if c.ShouldSignPostFund() {
		return []string{SignPostFundEffect(c)}, "", nil
	}

	if !c.IsPostFundComplete() {
		return NoSideEffects, CompletePostfund, nil
	}

	// Completion
	return []string{"Objective" + o.Id + "complete"}, "", nil
}

// mermaid diagram
// key:
// - effect!
// - waiting...
//
// https://mermaid-js.github.io/mermaid-live-editor/edit/#eyJjb2RlIjoiZ3JhcGggVERcbiAgICBTdGFydCAtLT4gQ3tJbnZhbGlkIElucHV0P31cbiAgICBDIC0tPnxZZXN8IEVbZXJyb3JdXG4gICAgQyAtLT58Tm98IEQwXG4gICAgXG4gICAgRDB7U2hvdWxkU2lnblByZUZ1bmR9XG4gICAgRDAgLS0-fFllc3wgUjFbU2lnblByZWZ1bmQhXVxuICAgIEQwIC0tPnxOb3wgRDFcbiAgICBcbiAgICBEMXtTYWZlVG9EZXBvc2l0ICY8YnI-ICFGdW5kaW5nQ29tcGxldGV9XG4gICAgRDEgLS0-IHxZZXN8IFIyW0Z1bmQgb24gY2hhaW4hXVxuICAgIEQxIC0tPiB8Tm98IEQyXG4gICAgXG4gICAgRDJ7IVNhZmVUb0RlcG9zaXQgJjxicj4gIUZ1bmRpbmdDb21wbGV0ZX1cbiAgICBEMiAtLT4gfFllc3wgUjNbXCJteSB0dXJuLi4uXCJdXG4gICAgRDIgLS0-IHxOb3wgRDNcblxuICAgIEQze1NhZmVUb0RlcG9zaXQgJjxicj4gIUZ1bmRpbmdDb21wbGV0ZX1cbiAgICBEMyAtLT4gfFllc3wgUjRbRGVwb3NpdCFdXG4gICAgRDMgLS0-IHxOb3wgRDRcblxuICAgIEQ0eyFGdW5kaW5nQ29tcGxldGV9XG4gICAgRDQgLS0-IHxZZXN8IFI1W1wiY29tcGxldGUgZnVuZGluZy4uLlwiXVxuICAgIEQ0IC0tPiB8Tm98IEQ1XG5cbiAgICBENXtTaG91bGRTaWduUHJlRnVuZH1cbiAgICBENSAtLT58WWVzfCBSNltTaWduUG9zdGZ1bmQhXVxuICAgIEQ1IC0tPnxOb3wgRDZcblxuICAgIEQ2eyFQb3N0RnVuZENvbXBsZXRlfVxuICAgIEQ2IC0tPnxZZXN8IFI3W1wiY29tcGxldGUgcG9zdGZ1bmQuLi5cIl1cbiAgICBENiAtLT58Tm98IFI4XG5cbiAgICBSOFtcImZpbmlzaFwiXVxuICAgIFxuXG5cbiIsIm1lcm1haWQiOiJ7fSIsInVwZGF0ZUVkaXRvciI6ZmFsc2UsImF1dG9TeW5jIjp0cnVlLCJ1cGRhdGVEaWFncmFtIjp0cnVlfQ//
// graph TD
//     Start --> C{Invalid Input?}
//     C -->|Yes| E[error]
//     C -->|No| D0

//     D0{ShouldSignPreFund}
//     D0 -->|Yes| R1[SignPrefund!]
//     D0 -->|No| D1

//     D1{SafeToDeposit &<br> !FundingComplete}
//     D1 --> |Yes| R2[Fund on chain!]
//     D1 --> |No| D2

//     D2{!SafeToDeposit &<br> !FundingComplete}
//     D2 --> |Yes| R3["my turn..."]
//     D2 --> |No| D3

//     D3{SafeToDeposit &<br> !FundingComplete}
//     D3 --> |Yes| R4[Deposit!]
//     D3 --> |No| D4

//     D4{!FundingComplete}
//     D4 --> |Yes| R5["complete funding..."]
//     D4 --> |No| D5

//     D5{ShouldSignPreFund}
//     D5 -->|Yes| R6[SignPostfund!]
//     D5 -->|No| D6

//     D6{!PostFundComplete}
//     D6 -->|Yes| R7["complete postfund..."]
//     D6 -->|No| R8

//     R8["finish"]
