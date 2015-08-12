//
//
//

package main

// 此工具用来发送IOU和ICC
// 使用如下：
// ./tx -logtostderr=true -sender=iN8sGowQCg1qptWcJG1WyTmymKX7y9cpmr -secret=ss1TCkz333t3t2J5eobcEMkMY4bXk -recipient=iwsZ7gxHFzu2xbj8YMf4UvK1PrDEMuxGkf -currency=ICC -amount=123.456

import (
	"errors"
	"flag"
	"fmt"
	// "strconv"

	"github.com/wangch/glog"
	"github.com/wangch/ripple/data"
	"github.com/wangch/ripple/websockets"
)

var server = flag.String("server", "wss://icloudcoin.org:19528", "服务器地址")
var sender = flag.String("sender", "", "发送者地址, 类似: iN8sGowQCg1qptWcJG1WyTmymKX7y9cpmr")
var secret = flag.String("secret", "", "发送者的secret key, 类似: ss9qoFiFNkokVfgrb3YkKHido7a1q")
var recipient = flag.String("recipient", "", "接收者地址, 类似: iN8sGowQCg1qptWcJG1WyTmymKX7y9cpmr")
var currency = flag.String("currency", "ICC", "货币种类 USD/CNY/HKY/ICC/EUR/GBP之一")
var amount = flag.Float64("amount", 0.0, "数额 比如: 123.456")

// var issuer = flag.String("issuer", "iN8sGowQCg1qptWcJG1WyTmymKX7y9cpmr", "发行者")

// var invoiceID = flag.String("invoiceID", "", "用于标识此次交易")

func main() {
	flag.Parse()
	if *server == "" {
		flag.Usage()
		return
	}
	if *sender == "" || len(*sender) != 34 || (*sender)[0] != 'i' {
		flag.Usage()
		return
	}
	if *secret == "" || len(*secret) != 29 || (*secret)[0] != 's' {
		flag.Usage()
		return
	}
	if *recipient == "" || len(*recipient) != 34 || (*recipient)[0] != 'i' {
		flag.Usage()
		return
	}
	if *amount <= 0 {
		glog.Errorln(*amount, "发送金额必须>0")
		return
	}

	ws, err := websockets.NewRemote(*server)
	if err != nil {
		glog.Fatal(err)
	}

	issuer := "iN8sGowQCg1qptWcJG1WyTmymKX7y9cpmr"
	err = payment(ws, *secret, *sender, issuer, *recipient, *currency, "", *amount)
	if err != nil {
		glog.Fatal(err)
	}
}

func payment(ws *websockets.Remote, secret, sender, issuer, recipient, currency, invoiceID string, amount float64) error {
	glog.Info("payment:", secret, sender, recipient, currency, amount, issuer, invoiceID)
	sam := ""
	if currency == "ICC" {
		sam = fmt.Sprintf("%d/ICC", uint64(amount))
	} else {
		sam = fmt.Sprintf("%f/%s/%s", amount, currency, issuer)
	}
	a, err := data.NewAmount(sam)
	if err != nil {
		err = errors.New("NewAmount error: " + sam + err.Error())
		glog.Error(err)
		return err
	}

	ptx := &websockets.PaymentTx{
		TransactionType: "Payment",
		Account:         sender,
		Destination:     recipient,
		Amount:          a,
		InvoiceID:       invoiceID,
	}

	// glog.Infof("payment: %+v", ptx)

	r, err := ws.SubmitWithSign(ptx, secret)
	if err != nil {
		glog.Error(err)
		return err
	}
	glog.Infof("pament result: %+v", r)
	if !r.EngineResult.Success() {
		return errors.New(r.EngineResultMessage)
	}
	return nil
}
