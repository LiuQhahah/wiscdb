package main

type txnCb struct {
	commit func() error
	user   func(err error)
	err    error
}

func runTxnCallback(cb *txnCb) {

}
