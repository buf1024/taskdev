package main

import (
	"fmt"
	"os"
	"taskdev/util"
)

type builderInfo struct {
}

type builder struct {
	util.Client
}

func (b *builder) detectInfo() error {
	return nil
}

func (b *builder) builderInit() {
	b.ParseArgs()
	err := b.InitLog("builder")
	if err != nil {
		fmt.Printf("init log failed, err = %s\n", err)
		os.Exit(-1)
	}

	err = b.detectInfo()
	if err != nil {
		fmt.Printf("detect os info failed, err = %s\n", err)
		os.Exit(-1)
	}
	b.InitClient()

}

func newBuilder() *builder {
	b := &builder{}
	return b
}

func main() {
	b := newBuilder()
	b.builderInit()

	b.Start()
	b.Wait()

}
