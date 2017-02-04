bindir=bin
exe=builder runner scheduler

builder_go=builder/builder.go
runner_go=runner/runner.go
scheduler_go=scheduler/scheduler.go

all:$(exe)


builder: $(builder_go)
	@echo "building $@"
	go build -o $(bindir)/$@ $^

runner: $(runner_go)
	@echo "building $@"
	go build -o $(bindir)/$@ $^

scheduler: $(scheduler_go)
	@echo "building $@"
	go build -o $(bindir)/$@ $^

