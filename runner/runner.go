package main

import (
	"fmt"
	"os"
	"taskdev/util"
)

type runnerInfo struct {
}

type runner struct {
	util.Client
}

func (b *runner) detectInfo() error {
	return nil
}

func (b *runner) runnerInit() {
	b.ParseArgs()
	err := b.InitLog("runner")
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

func newBuilder() *runner {
	b := &runner{}
	return b
}

func main() {
	b := newBuilder()
	b.runnerInit()

	b.Start()
	b.Wait()

}
